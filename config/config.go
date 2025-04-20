package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/eust-w/urlreader/internal/logger"
)

// Config 存储应用程序配置
type Config struct {
	Port                 string
	AzureOpenAIKey       string
	AzureOpenAIEndpoint  string
	AzureOpenAIDeployment string
	AzureOpenAIAPIVersion string
	DeepseekAPIKey       string
	DeepseekAPIEndpoint  string
	DeepseekModel        string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.GetLogger().Warnw(".env file not found, using environment variables")
	}

	config := &Config{
		Port:                 getEnv("PORT", "8080"),
		AzureOpenAIKey:       getEnv("AZURE_OPENAI_API_KEY", ""),
		AzureOpenAIEndpoint:  getEnv("AZURE_OPENAI_ENDPOINT", ""),
		AzureOpenAIDeployment: getEnv("AZURE_OPENAI_DEPLOYMENT", ""),
		AzureOpenAIAPIVersion: getEnv("AZURE_OPENAI_API_VERSION", "2023-05-15"),
		DeepseekAPIKey:       getEnv("DEEPSEEK_API_KEY", ""),
		DeepseekAPIEndpoint:  getEnv("DEEPSEEK_API_ENDPOINT", "https://api.deepseek.com"),
		DeepseekModel:        getEnv("DEEPSEEK_MODEL", "deepseek-chat"),
	}

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
