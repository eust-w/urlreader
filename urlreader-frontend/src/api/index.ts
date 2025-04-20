import request from './request';

export const parseUrl = async (url: string) => {
  const res = await request.post('/parse', { url });
  return res.data;
};

export const chatWithPage = async (params: {
  url: string;
  message: string;
  model?: string;
  conversation_id?: string;
}) => {
  const res = await request.post('/chat', params);
  return res.data;
};

export const getHistory = async (conversation_id: string) => {
  const res = await request.get(`/history/${conversation_id}`);
  return res.data;
};

export const getConversations = async () => {
  const res = await request.get('/conversations');
  return res.data;
};

export const deleteConversation = async (conversation_id: string) => {
  const res = await request.delete(`/history/${conversation_id}`);
  return res.data;
};
