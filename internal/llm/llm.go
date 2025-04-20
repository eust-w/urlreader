package llm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eust-w/urlreader/config"
	"github.com/eust-w/urlreader/internal/logger"
)

// LLMProvider 接口定义了所有LLM提供商必须实现的方法
type LLMProvider interface {
	Chat(messages []Message) (string, error)
	Name() string
}

// Message 表示聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMFactory 创建不同的LLM提供商实例
type LLMFactory struct {
	config *config.Config
}

// NewLLMFactory 创建一个新的LLM工厂
func NewLLMFactory(cfg *config.Config) *LLMFactory {
	logger.GetLogger().Infow("初始化 LLMFactory", "provider_keys", map[string]string{
		"AzureOpenAIKey": cfg.AzureOpenAIKey,
		"DeepseekAPIKey": cfg.DeepseekAPIKey,
	})
	return &LLMFactory{
		config: cfg,
	}
}

// GetProvider 根据名称返回相应的LLM提供商
func (f *LLMFactory) GetProvider(name string) (LLMProvider, error) {
	log := logger.GetLogger()
	log.Infow("请求 LLM Provider", "name", name)
	switch strings.ToLower(name) {
	case "azure_openai", "azure", "openai":
		if f.config.AzureOpenAIKey == "" || f.config.AzureOpenAIEndpoint == "" {
			log.Errorw("Azure OpenAI API密钥或端点未配置", "AzureOpenAIKey", f.config.AzureOpenAIKey, "AzureOpenAIEndpoint", f.config.AzureOpenAIEndpoint)
			return nil, errors.New("Azure OpenAI API密钥或端点未配置")
		}
		log.Infow("使用 AzureOpenAI Provider")
		return NewAzureOpenAIProvider(f.config), nil
	case "deepseek":
		if f.config.DeepseekAPIKey == "" {
			log.Errorw("DeepSeek API密钥未配置", "DeepseekAPIKey", f.config.DeepseekAPIKey)
			return nil, errors.New("DeepSeek API密钥未配置")
		}
		log.Infow("使用 Deepseek Provider")
		return NewDeepseekProvider(f.config), nil
	default:
		log.Errorw("不支持的LLM提供商", "name", name)
		return nil, fmt.Errorf("不支持的LLM提供商: %s", name)
	}
}

// CreateContextPrompt 创建包含网页内容的上下文提示
func CreateContextPrompt(content string, query string) string {
	return fmt.Sprintf(`以下是从网页抓取的内容:

%s

请基于上述网页内容回答以下问题:
%s

如果网页内容中没有相关信息，请明确指出。请保持回答简洁、准确，并直接基于提供的网页内容。`, content, query)
}
