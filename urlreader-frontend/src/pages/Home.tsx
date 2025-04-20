import React, { useState } from 'react';
import { Input, Button, Card, Typography, Alert, Space } from 'antd';
import { parseUrl } from '../api';
import { useTranslation } from 'react-i18next';
import ReactMarkdown from 'react-markdown';

const { Title } = Typography;

const Home: React.FC = () => {
  const { t } = useTranslation();
  const [url, setUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState('');

  const handleParse = async () => {
    setLoading(true);
    setResult(null);
    setError('');
    try {
      const res = await parseUrl(url);
      if (res.success) {
        setResult(res);
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
    <Card style={{ maxWidth: 600, margin: '40px auto' }}>
      <Space direction="vertical" style={{ width: '100%' }}>
        <Input
          placeholder={t('input.url')}
          value={url}
          onChange={e => setUrl(e.target.value)}
          onPressEnter={handleParse}
          disabled={loading}
        />
        <Button type="primary" loading={loading} onClick={handleParse} disabled={!url}>
          {t('button.parse')}
        </Button>
        {error && <Alert type="error" message={t('parse.error', { error })} showIcon />}
        {result && (
          <div>
            <Title level={4}>{t('parse.result.title')}</Title>
            <ReactMarkdown>{result.title}</ReactMarkdown>
            <Title level={5}>{t('parse.result.content')}</Title>
            <ReactMarkdown>{result.content}</ReactMarkdown>
          </div>
        )}
      </Space>
    </Card>
  );
};

export default Home;
