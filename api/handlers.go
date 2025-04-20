package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eust-w/urlreader/config"
	"github.com/eust-w/urlreader/internal/llm"
	"github.com/eust-w/urlreader/internal/models"
	"github.com/eust-w/urlreader/internal/scraper"
	"github.com/eust-w/urlreader/internal/storage"
	"github.com/eust-w/urlreader/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler 处理API请求
type Handler struct {
	config        *config.Config
	scraper       *scraper.Scraper
	llmFactory    *llm.LLMFactory
	conversations *storage.ConversationStore
}

// NewHandler 创建一个新的API处理程序
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config:        cfg,
		scraper:       scraper.NewScraper(),
		llmFactory:    llm.NewLLMFactory(cfg),
		conversations: storage.NewConversationStore(),
	}
}

// SetupRoutes 设置API路由
func (h *Handler) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/parse", h.ParseURL)
		api.POST("/chat", h.Chat)
		api.GET("/history/:conversation_id", h.GetHistory)
		api.GET("/conversations", h.ListConversations)
		api.DELETE("/history/:conversation_id", h.DeleteConversation)
	}
}

// ParseURL 处理URL解析请求
func (h *Handler) ParseURL(c *gin.Context) {
	var req models.ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "无效的请求: " + err.Error(),
		})
		return
	}

	content, err := h.scraper.ScrapeURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "抓取URL失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ParseResponse{
		Success: true,
		Title:   content.Title,
		Content: content.Content,
		URL:     content.URL,
	})
}

// Chat 处理聊天请求
func (h *Handler) Chat(c *gin.Context) {
	log := logger.GetLogger()
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorw("/api/chat 参数解析失败", "error", err, "raw_body", c.Request.Body)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "无效的请求: " + err.Error(),
		})
		return
	}
	log.Infow("/api/chat 收到请求", "req", req)

	// 首次对话必须提供URL
	if req.ConversationID == "" && req.URL == "" {
		log.Errorw("/api/chat 缺少URL", "req", req)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "首次对话必须提供URL",
		})
		return
	}

	// 如果未指定模型，使用默认模型
	if req.Model == "" {
		req.Model = "azure_openai"
	}

	// 获取LLM提供商
	provider, err := h.llmFactory.GetProvider(req.Model)
	if err != nil {
		log.Errorw("获取LLM Provider失败", "model", req.Model, "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "LLM提供商错误: " + err.Error(),
		})
		return
	}

	var conversation *storage.Conversation
	var exists bool

	// 处理会话ID
	if req.ConversationID != "" {
		// 使用现有会话
		conversation, exists = h.conversations.Get(req.ConversationID)
		if !exists {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   "会话不存在",
			})
			return
		}
	} else {
		// 创建新会话，首先抓取URL内容
		content, err := h.scraper.ScrapeURL(req.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "抓取URL失败: " + err.Error(),
			})
			return
		}

		// 创建新会话
		conversationID := uuid.New().String()
		conversation = h.conversations.Create(conversationID, content.URL, content.Content)
		req.ConversationID = conversationID
	}

	// 准备消息
	var messages []llm.Message

	if len(conversation.Messages) == 0 {
		// 新会话，添加系统消息
		systemPrompt := fmt.Sprintf("你是一个网页内容助手。你将基于从URL %s 抓取的内容回答问题。请保持回答简洁、准确，并直接基于提供的网页内容。", conversation.URL)
		messages = append(messages, llm.Message{
			Role:    "system",
			Content: systemPrompt,
		})

		// 添加包含网页内容的第一条消息
		contentPrompt := fmt.Sprintf("以下是从网页抓取的内容:\n\n%s", conversation.Content)
		messages = append(messages, llm.Message{
			Role:    "user",
			Content: contentPrompt,
		})

		// 添加助手确认消息
		messages = append(messages, llm.Message{
			Role:    "assistant",
			Content: "我已经阅读了网页内容，请问有什么我可以帮助你的？",
		})

		// 保存这些初始消息到会话
		for _, msg := range messages {
			h.conversations.AddMessage(req.ConversationID, msg)
		}
	} else {
		// 使用现有会话的消息历史
		messages, _ = h.conversations.GetMessages(req.ConversationID)
	}

	// 添加用户的新消息
	userMessage := llm.Message{
		Role:    "user",
		Content: req.Message,
	}
	messages = append(messages, userMessage)
	h.conversations.AddMessage(req.ConversationID, userMessage)

	// 调用LLM获取响应
	response, err := provider.Chat(messages)
	if err != nil {
		// 检查是否是 Azure OpenAI 的速率限制错误
		if req.Model == "azure_openai" && (err.Error() == "API错误: {\"error\":{\"code\":\"429\"" || err.Error() == "状态码: 429") {
			// 尝试切换到 DeepSeek 模型
			fmt.Println("Azure OpenAI 速率限制，切换到 DeepSeek 模型")
			deepseekProvider, deepseekErr := h.llmFactory.GetProvider("deepseek")
			if deepseekErr != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Success: false,
					Error:   "LLM响应错误: " + err.Error() + "，切换到备用模型失败: " + deepseekErr.Error(),
				})
				return
			}
			
			// 使用 DeepSeek 模型重试
			response, err = deepseekProvider.Chat(messages)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Success: false,
					Error:   "备用模型响应错误: " + err.Error(),
				})
				return
			}
			// 更新当前使用的模型为 DeepSeek
			provider = deepseekProvider
		} else {
			log.Errorw("LLM 响应错误", "model", req.Model, "error", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "LLM响应错误: " + err.Error(),
			})
			return
		}
	}

	// 保存助手响应到会话
	assistantMessage := llm.Message{
		Role:    "assistant",
		Content: response,
	}
	h.conversations.AddMessage(req.ConversationID, assistantMessage)

	// 返回响应
	c.JSON(http.StatusOK, models.ChatResponse{
		Success:        true,
		Response:       response,
		ConversationID: req.ConversationID,
		Model:          provider.Name(),
	})
}

// ListConversations 获取所有会话ID
func (h *Handler) ListConversations(c *gin.Context) {
	ids := h.conversations.ListIDs()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"conversation_ids": ids,
	})
}

// DeleteConversation 删除指定conversation_id及其历史
func (h *Handler) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "conversation_id不能为空"})
		return
	}
	ok := h.conversations.Delete(conversationID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "会话不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "conversation_id": conversationID})
}

// GetHistory 查询指定conversation_id的历史消息
func (h *Handler) GetHistory(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "conversation_id不能为空"})
		return
	}
	messages, ok := h.conversations.GetMessages(conversationID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "会话不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"conversation_id": conversationID,
		"messages": messages,
	})
}

// StartCleanupTask 启动定期清理旧会话的任务
func (h *Handler) StartCleanupTask() {
	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			// 清理24小时未活动的会话
			count := h.conversations.CleanupOldConversations(24 * time.Hour)
			if count > 0 {
				fmt.Printf("已清理 %d 个过期会话\n", count)
			}
		}
	}()
}
