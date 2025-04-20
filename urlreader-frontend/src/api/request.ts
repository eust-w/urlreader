import axios from 'axios';

const baseURL = import.meta.env.VITE_API_BASE_URL || '/api';

const request = axios.create({
  baseURL,
  timeout: 10000,
  // 可添加 headers、拦截器等
});

export default request;
