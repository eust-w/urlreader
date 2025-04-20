import React, { useState, useRef, useEffect } from 'react';
import { Input, Button, Card, Typography, Alert, Select, Space, message as antdMessage, Modal } from 'antd';
import { chatWithPage, getHistory, getConversations, deleteConversation } from '../api';
import { useTranslation } from 'react-i18next';
import ConversationList from '../components/ConversationList';
import ReactMarkdown from 'react-markdown';

const { Paragraph } = Typography;

const Chat: React.FC = () => {
  const { t } = useTranslation();
  const [url, setUrl] = useState('');
  const [message, setMessage] = useState('');
  const [model, setModel] = useState<'azure_openai' | 'deepseek'>('azure_openai');
  const [conversationId, setConversationId] = useState<string | undefined>();
  const [history, setHistory] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [conversations, setConversations] = useState<string[]>([]);
  const [convLoading, setConvLoading] = useState(false);
  const [showNewChat, setShowNewChat] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const fetchConversations = async () => {
    setConvLoading(true);
    try {
      const res = await getConversations();
      if (res.success) {
        setConversations(res.conversation_ids || []);
      }
    } finally {
      setConvLoading(false);
    }
  };

  const fetchHistory = async (cid: string) => {
    try {
      const res = await getHistory(cid);
      if (res.success) {
        setHistory(res.messages || []);
      } else {
        setError(res.error || '');
      }
    } catch (e: any) {
      setError(e.message || 'Unknown error');
    }
  };


  useEffect(() => {
    fetchConversations();
  }, []);

  useEffect(() => {
    if (conversationId) fetchHistory(conversationId);
  }, [conversationId]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [history]);

  const handleSelectConversation = (id: string) => {
    setConversationId(id);
    setUrl('');
    setMessage('');
    setError('');
  };

  const handleDeleteConversation = async (id: string) => {
    await deleteConversation(id);
    antdMessage.success(t('button.delete') + ' success');
    fetchConversations();
    if (conversationId === id) {
      setConversationId(undefined);
      setHistory([]);
    }
  };

  const handleStartNewChat = () => {
    setConversationId(undefined);
    setHistory([]);
    setUrl('');
    setMessage('');
    setError('');
    setShowNewChat(true);
  };

  const handleSend = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await chatWithPage({ url, message, model, conversation_id: conversationId });
      if (res.success) {
        setConversationId(res.conversation_id);
        fetchConversations();
        setHistory(prev => [...prev, { role: 'user', content: message }, { role: 'assistant', content: res.response }]);
        setMessage('');
        setShowNewChat(false);
      } else {
        setError(res.error || '');
      }
    } catch (e: any) {
      setError(e.message || 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', height: '80vh', maxWidth: 1100, margin: '40px auto', background: '#fff', borderRadius: 10, boxShadow: '0 2px 8px #0001' }}>
      <ConversationList
        conversations={conversations}
        loading={convLoading}
        currentId={conversationId}
        onSelect={handleSelectConversation}
        onDelete={handleDeleteConversation}
      />
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', padding: 24, minWidth: 0 }}>
        <div style={{ marginBottom: 12, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <span style={{ fontWeight: 500, fontSize: 18 }}>{conversationId ? t('chat.history') : t('button.newChat')}</span>
          <Button type="primary" onClick={handleStartNewChat}>{t('button.newChat')}</Button>
        </div>
        <Card style={{ flex: 1, display: 'flex', flexDirection: 'column', minHeight: 400 }} bodyStyle={{ padding: 0 }}>
          {/* 聊天历史展示区 */}
          <div style={{
            flex: 1,
            minHeight: 320,
            maxHeight: 400,
            overflowY: 'auto',
            background: '#f7f7f7',
            padding: 16,
            borderRadius: 8,
            margin: 16,
            display: 'flex',
            flexDirection: 'column',
          }}>
            {history.length === 0 && (
              <div style={{ color: '#aaa', textAlign: 'center' }}>{t('chat.history')}</div>
            )}
            {history.map((item, idx) => (
              <div
                key={idx}
                style={{
                  display: 'flex',
                  justifyContent: item.role === 'user' ? 'flex-end' : 'flex-start',
                  marginBottom: 12,
                }}
              >
                {item.role === 'assistant' && (
                  <span style={{ alignSelf: 'flex-end', marginRight: 8 }}>🤖</span>
                )}
                <div
                  style={{
                    maxWidth: '70%',
                    background: item.role === 'user' ? '#95ec69' : '#fff',
                    color: '#222',
                    borderRadius: 12,
                    padding: '8px 14px',
                    boxShadow: '0 1px 2px #0001',
                    wordBreak: 'break-word',
                    textAlign: 'left',
                  }}
                >
                  {item.role === 'assistant' ? (
                    <ReactMarkdown>{item.content}</ReactMarkdown>
                  ) : (
                    item.content
                  )}
                </div>
                {item.role === 'user' && (
                  <span style={{ alignSelf: 'flex-end', marginLeft: 8 }}>🧑</span>
                )}
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
          {/* 输入区 */}
          <div style={{ padding: 16, borderTop: '1px solid #eee', background: '#fff' }}>
            <Space direction="vertical" style={{ width: '100%' }}>
              <Input
                placeholder={t('input.url')}
                value={url}
                onChange={e => setUrl(e.target.value)}
                disabled={loading || (!!conversationId && !showNewChat)}
              />
              <Select
                value={model}
                onChange={setModel}
                style={{ width: 180 }}
                options={[
                  { value: 'azure_openai', label: t('model.azure') },
                  { value: 'deepseek', label: t('model.deepseek') },
                ]}
                disabled={loading}
              />
              <Input.TextArea
                rows={3}
                placeholder={t('chat.input')}
                value={message}
                onChange={e => setMessage(e.target.value)}
                onPressEnter={e => { if (!e.shiftKey) handleSend(); }}
                disabled={loading || (!url && !conversationId)}
              />
              <Button type="primary" loading={loading} onClick={handleSend} disabled={loading || (!url && !conversationId) || !message}>
                {t('button.send')}
              </Button>
              {error && <Alert type="error" message={t('chat.error', { error })} showIcon />}
            </Space>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Chat;
