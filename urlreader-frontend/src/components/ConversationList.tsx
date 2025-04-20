import React from 'react';
import { List, Button, Popconfirm, Typography, Spin, message as antdMessage } from 'antd';
import { useTranslation } from 'react-i18next';

interface Props {
  conversations: string[];
  loading: boolean;
  currentId?: string;
  onSelect: (id: string) => void;
  onDelete: (id: string) => void;
}

const { Text } = Typography;

const ConversationList: React.FC<Props> = ({ conversations, loading, currentId, onSelect, onDelete }) => {
  const { t } = useTranslation();

  return (
    <div style={{ width: 220, borderRight: '1px solid #eee', height: '100%', overflowY: 'auto', background: '#fafbfc' }}>
      <div style={{ padding: '12px 16px', fontWeight: 'bold', borderBottom: '1px solid #eee' }}>{t('chat.history')}</div>
      {loading ? (
        <Spin style={{ margin: 32, display: 'block' }} />
      ) : (
        <List
          size="small"
          dataSource={conversations}
          locale={{ emptyText: t('chat.noConversations') }}
          renderItem={id => (
            <List.Item
              style={{
                background: currentId === id ? '#e6f4ff' : undefined,
                cursor: 'pointer',
                padding: '8px 16px',
                display: 'flex',
                alignItems: 'center',
                borderLeft: currentId === id ? '3px solid #1677ff' : '3px solid transparent',
                marginBottom: 2,
              }}
              onClick={() => onSelect(id)}
            >
              <Text ellipsis style={{ flex: 1, color: currentId === id ? '#1677ff' : undefined }}>{id}</Text>
              <Popconfirm
                title={t('chat.deleteConfirm')}
                onConfirm={e => { e?.stopPropagation(); onDelete(id); }}
                onCancel={e => e?.stopPropagation()}
                okText={t('button.confirm')}
                cancelText={t('button.cancel')}
              >
                <Button size="small" type="link" danger style={{ marginLeft: 8 }} onClick={e => e.stopPropagation()}>
                  {t('button.delete')}
                </Button>
              </Popconfirm>
            </List.Item>
          )}
        />
      )}
    </div>
  );
};

export default ConversationList;
