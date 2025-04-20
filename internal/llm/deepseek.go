package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eust-w/urlreader/config"
)

// DeepseekProvider 实现了DeepSeek API
type DeepseekProvider struct {
	apiKey     string
	endpoint   string
	model      string
	httpClient *http.Client
}

// NewDeepseekProvider 创建一个新的DeepSeek提供商
func NewDeepseekProvider(cfg *config.Config) *DeepseekProvider {
	return &DeepseekProvider{
		apiKey:     cfg.DeepseekAPIKey,
		endpoint:   cfg.DeepseekAPIEndpoint,
		model:      cfg.DeepseekModel,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// Name 返回提供商名称
func (p *DeepseekProvider) Name() string {
	return "DeepSeek"
}

// DeepseekRequest DeepSeek API请求结构
type DeepseekRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// DeepseekResponse DeepSeek API响应结构
type DeepseekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// Chat 使用DeepSeek进行聊天
func (p *DeepseekProvider) Chat(messages []Message) (string, error) {
	url := fmt.Sprintf("%s/v1/chat/completions", p.endpoint)

	requestBody := DeepseekRequest{
		Model:       p.model,
		Messages:    messages,
		MaxTokens:   2000,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	var response DeepseekResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("API错误: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("API没有返回任何选择")
	}

	return response.Choices[0].Message.Content, nil
}
