import axios from 'axios';
import request from '@/utils/request';
const ACCESS_TOKEN_KEY = 'opsnexus_access_token';
export const listAuditLogsApi = (params) => request.get('/v1/audit/logs', { params });
export const listLoginLogsApi = (params) => request.get('/v1/audit/login-logs', { params });
export const exportAuditLogsApi = async (params) => {
    const token = localStorage.getItem(ACCESS_TOKEN_KEY) || '';
    const res = await axios.get('/api/v1/audit/logs/export', {
        params,
        responseType: 'blob',
        headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    return res.data;
};
