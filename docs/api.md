# API 文档

本项目提供了基于 Gin 框架的 RESTful API，主要包括网页内容解析和多轮对话两大功能。

## 接口总览

- [POST /api/parse](#post-apiparse)
- [POST /api/chat](#post-apichat)
- [GET /api/history/:conversation_id](#get-apihistoryconversation_id)
- [GET /api/conversations](#get-apiconversations)
- [DELETE /api/history/:conversation_id](#delete-apihistoryconversation_id)

---

## POST /api/parse

解析指定 URL 的网页内容。

### 请求
- 路径：`/api/parse`
- 方法：POST
- Content-Type: `application/json`

#### 请求体
```json
{
  "url": "https://example.com"
}
```

| 字段 | 类型   | 是否必填 | 说明         |
|------|--------|----------|--------------|
| url  | string | 是       | 目标网页URL  |

#### 响应体
```json
{
  "success": true,
  "title": "网页标题",
  "content": "网页正文内容",
  "url": "https://example.com"
}
```

| 字段    | 类型   | 说明           |
|---------|--------|----------------|
| success | bool   | 是否成功       |
| title   | string | 网页标题       |
| content | string | 网页正文内容   |
| url     | string | 原始URL        |
| error   | string | 错误信息（可选）|

### 错误响应示例
```json
{
  "success": false,
  "error": "抓取URL失败: ..."
}
```

---

## POST /api/chat

基于指定网页内容进行上下文多轮对话。

### 请求
- 路径：`/api/chat`
- 方法：POST
- Content-Type: `application/json`

#### 请求体
```json
{
  "url": "https://example.com",
  "message": "请总结这个网页的内容",
  "model": "azure_openai",  // 可选: azure_openai, deepseek
  "conversation_id": "uuid"  // 可选，用于多轮对话
}
```

| 字段           | 类型   | 是否必填 | 说明                      |
|----------------|--------|----------|---------------------------|
| url            | string | 是       | 目标网页URL               |
| message        | string | 是       | 用户输入的对话内容        |
| model          | string | 否       | LLM模型（azure_openai, deepseek）|
| conversation_id| string | 否       | 对话ID（多轮对话用）      |

#### 响应体
```json
{
  "success": true,
  "response": "助手回复内容",
  "conversation_id": "uuid",
  "model": "azure_openai"
}
```

| 字段           | 类型   | 说明                  |
|----------------|--------|-----------------------|
| success        | bool   | 是否成功              |
| response       | string | 助手回复内容          |
| conversation_id| string | 当前对话ID            |
| model          | string | 实际使用的LLM模型      |
| error          | string | 错误信息（可选）      |

### 错误响应示例
```json
{
  "success": false,
  "error": "无效的请求: ..."
}
```

---

## GET /api/conversations

获取所有有效的 conversation_id。

### 请求
- 路径：`/api/conversations`
- 方法：GET

#### 响应体
```json
{
  "success": true,
  "conversation_ids": ["uuid1", "uuid2", ...]
}
```

| 字段             | 类型     | 说明         |
|------------------|----------|--------------|
| success          | bool     | 是否成功     |
| conversation_ids | string[] | 会话ID数组   |
| error            | string   | 错误信息（可选）|

### 错误响应示例
```json
{
  "success": false,
  "error": "服务器内部错误"
}
```

---

## DELETE /api/history/:conversation_id

删除指定 conversation_id 及其历史消息。

### 请求
- 路径：`/api/history/:conversation_id`
- 方法：DELETE

#### 路径参数
| 参数名           | 类型   | 是否必填 | 说明             |
|------------------|--------|----------|------------------|
| conversation_id  | string | 是       | 对话ID           |

#### 响应体
```json
{
  "success": true,
  "conversation_id": "uuid"
}
```

| 字段           | 类型   | 说明         |
|----------------|--------|--------------|
| success        | bool   | 是否成功     |
| conversation_id| string | 被删除的ID   |
| error          | string | 错误信息（可选）|

### 错误响应示例
```json
{
  "success": false,
  "error": "会话不存在"
}
```

---

## GET /api/history/:conversation_id

查询指定 conversation_id 的历史消息。

### 请求
- 路径：`/api/history/:conversation_id`
- 方法：GET

#### 路径参数
| 参数名           | 类型   | 是否必填 | 说明             |
|------------------|--------|----------|------------------|
| conversation_id  | string | 是       | 对话ID           |

#### 响应体
```json
{
  "success": true,
  "conversation_id": "uuid",
  "messages": [
    { "role": "user", "content": "用户消息内容" },
    { "role": "assistant", "content": "助手回复内容" }
  ]
}
```

| 字段           | 类型         | 说明             |
|----------------|--------------|------------------|
| success        | bool         | 是否成功         |
| conversation_id| string       | 对话ID           |
| messages       | Message[]    | 历史消息数组     |
| error          | string       | 错误信息（可选） |

#### Message 结构
```json
{
  "role": "user | assistant | system",
  "content": "消息内容"
}
```

### 错误响应示例
```json
{
  "success": false,
  "error": "会话不存在"
}
```

---

## 相关数据结构

### ParseRequest
```go
type ParseRequest struct {
    URL string `json:"url" binding:"required"`
}
```

### ParseResponse
```go
type ParseResponse struct {
    Success bool   `json:"success"`
    Title   string `json:"title,omitempty"`
    Content string `json:"content,omitempty"`
    URL     string `json:"url,omitempty"`
    Error   string `json:"error,omitempty"`
}
```

### ChatRequest
```go
type ChatRequest struct {
    URL            string `json:"url" binding:"required"`
    Message        string `json:"message" binding:"required"`
    Model          string `json:"model,omitempty"`
    ConversationID string `json:"conversation_id,omitempty"`
}
```

### ChatResponse
```go
type ChatResponse struct {
    Success        bool   `json:"success"`
    Response       string `json:"response,omitempty"`
    ConversationID string `json:"conversation_id,omitempty"`
    Model          string `json:"model,omitempty"`
    Error          string `json:"error,omitempty"`
}
```

---

## 错误码说明
- 400 Bad Request：请求参数无效或缺失。
- 500 Internal Server Error：服务器内部错误，如抓取失败、LLM响应错误等。

---

## 联系方式
如有问题或建议，请提交 issue。
