import request from '@/utils/request';
export const getCaptchaApi = () => request.get('/v1/auth/captcha');
export const loginApi = (data) => request.post('/v1/auth/login', data);
export const logoutApi = () => request.post('/v1/auth/logout');
export const refreshApi = (refresh_token) => request.post('/v1/auth/refresh', { refresh_token });
export const getUserInfoApi = () => request.get('/v1/auth/userinfo');
