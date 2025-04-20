# URL Reader

一个基于Go和Gin框架的网页内容读取和对话系统，类似于JinaAI读取器。

## 功能特点

- 输入URL，读取网页内容
- 基于网页内容进行上下文多轮对话
- 支持多种LLM API（Azure OpenAI、DeepSeek等）
- 提供两个API接口：
  - URL解析接口：仅解析和返回网页内容
  - 对话接口：解析网页内容并进行上下文对话

## 安装与运行

### 环境要求

- Go 1.18+

### 安装依赖

```bash
go mod tidy
```

### 配置环境变量

创建`.env`文件并配置相应的API密钥：

```
AZURE_OPENAI_API_KEY=your_azure_openai_api_key
AZURE_OPENAI_ENDPOINT=your_azure_openai_endpoint
DEEPSEEK_API_KEY=your_deepseek_api_key
```

### 运行服务

```bash
go run main.go
```

## API使用说明

详细API文档请见：[API 文档](docs/api.md)


### 1. URL解析接口

```
POST /api/parse
```

请求体：
```json
{
  "url": "https://example.com"
}
```

### 2. 对话接口

```
POST /api/chat
```

请求体：
```json
{
  "url": "https://example.com",
  "message": "请总结这个网页的内容",
  "model": "azure_openai",  // 可选: azure_openai, deepseek
  "conversation_id": "uuid"  // 可选，用于多轮对话
}
```