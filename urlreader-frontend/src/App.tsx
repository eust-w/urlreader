import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, Navigate, useLocation } from 'react-router-dom';
import { ConfigProvider, Layout, Menu, Button } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import enUS from 'antd/locale/en_US';
import { useTranslation } from 'react-i18next';
import Home from './pages/Home';
import Chat from './pages/Chat';
import './i18n';

const { Header, Content } = Layout;

const AppContent: React.FC = () => {
  const { t, i18n } = useTranslation();
  const [lang, setLang] = React.useState(i18n.language);
  const location = useLocation();

  const handleSwitchLang = () => {
    const next = lang === 'zh' ? 'en' : 'zh';
    i18n.changeLanguage(next);
    setLang(next);
  };

  return (
    <ConfigProvider locale={lang === 'zh' ? zhCN : enUS}>
      <Layout style={{ minHeight: '100vh' }}>
        <Header style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <div style={{ color: 'white', fontWeight: 'bold', fontSize: 22 }}>{t('app.title')}</div>
          <Menu
            theme="dark"
            mode="horizontal"
            selectedKeys={[location.pathname]}
            items={[
              { key: '/', label: <Link to="/">{t('nav.home')}</Link> },
              { key: '/chat', label: <Link to="/chat">{t('nav.chat')}</Link> },
            ]}
            style={{ flex: 1, marginLeft: 30 }}
          />
          <Button onClick={handleSwitchLang}>{t('button.switchLang')}</Button>
        </Header>
        <Content style={{ padding: 24 }}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/chat" element={<Chat />} />
            <Route path="*" element={<Navigate to="/" />} />
          </Routes>
        </Content>
      </Layout>
    </ConfigProvider>
  );
};

const App: React.FC = () => (
  <Router>
    <AppContent />
  </Router>
);

export default App;
