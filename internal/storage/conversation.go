package storage

import (
	"sync"
	"time"

	"github.com/eust-w/urlreader/internal/llm"
)

// Conversation 表示一个对话会话
type Conversation struct {
	ID        string        `json:"id"`
	URL       string        `json:"url"`
	Content   string        `json:"content"`
	Messages  []llm.Message `json:"messages"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// ConversationStore 管理对话会话
type ConversationStore struct {
	conversations map[string]*Conversation
	mu            sync.RWMutex
}

// ListIDs 返回所有会话ID
func (s *ConversationStore) ListIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]string, 0, len(s.conversations))
	for id := range s.conversations {
		ids = append(ids, id)
	}
	return ids
}

// NewConversationStore 创建一个新的对话存储
func NewConversationStore() *ConversationStore {
	return &ConversationStore{
		conversations: make(map[string]*Conversation),
	}
}

// Get 获取指定ID的对话
func (s *ConversationStore) Get(id string) (*Conversation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, exists := s.conversations[id]
	return conv, exists
}

// Create 创建一个新的对话
func (s *ConversationStore) Create(id, url, content string) *Conversation {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	conv := &Conversation{
		ID:        id,
		URL:       url,
		Content:   content,
		Messages:  []llm.Message{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.conversations[id] = conv
	return conv
}

// AddMessage 向对话添加一条消息
func (s *ConversationStore) AddMessage(id string, message llm.Message) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, exists := s.conversations[id]
	if !exists {
		return false
	}

	conv.Messages = append(conv.Messages, message)
	conv.UpdatedAt = time.Now()
	return true
}

// GetMessages 获取对话的所有消息
func (s *ConversationStore) GetMessages(id string) ([]llm.Message, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, exists := s.conversations[id]
	if !exists {
		return nil, false
	}

	// 返回消息的副本以避免并发修改
	messages := make([]llm.Message, len(conv.Messages))
	copy(messages, conv.Messages)

	return messages, true
}

// Delete 删除指定ID的对话
func (s *ConversationStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.conversations[id]
	if !exists {
		return false
	}

	delete(s.conversations, id)
	return true
}

// CleanupOldConversations 清理超过指定时间的旧对话
func (s *ConversationStore) CleanupOldConversations(maxAge time.Duration) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	count := 0

	for id, conv := range s.conversations {
		if now.Sub(conv.UpdatedAt) > maxAge {
			delete(s.conversations, id)
			count++
		}
	}

	return count
}
