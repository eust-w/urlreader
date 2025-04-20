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

// AzureOpenAIProvider 实现了Azure OpenAI API
type AzureOpenAIProvider struct {
	apiKey     string
	endpoint   string
	deployment string
	apiVersion string
	httpClient *http.Client
}

// NewAzureOpenAIProvider 创建一个新的Azure OpenAI提供商
func NewAzureOpenAIProvider(cfg *config.Config) *AzureOpenAIProvider {
	return &AzureOpenAIProvider{
		apiKey:     cfg.AzureOpenAIKey,
		endpoint:   cfg.AzureOpenAIEndpoint,
		deployment: cfg.AzureOpenAIDeployment,
		apiVersion: cfg.AzureOpenAIAPIVersion,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// Name 返回提供商名称
func (p *AzureOpenAIProvider) Name() string {
	return "Azure OpenAI"
}

// AzureOpenAIRequest Azure OpenAI API请求结构
type AzureOpenAIRequest struct {
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// AzureOpenAIResponse Azure OpenAI API响应结构
type AzureOpenAIResponse struct {
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

// Chat 使用Azure OpenAI进行聊天
func (p *AzureOpenAIProvider) Chat(messages []Message) (string, error) {
	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		p.endpoint, p.deployment, p.apiVersion)

	requestBody := AzureOpenAIRequest{
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
	req.Header.Set("api-key", p.apiKey)

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

	var response AzureOpenAIResponse
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
