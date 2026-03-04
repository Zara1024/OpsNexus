import request from '@/utils/request';
export const getDashboardOverviewApi = () => request.get('/v1/dashboard/overview');
export const getDashboardTrendsApi = (days = 7) => request.get('/v1/dashboard/trends', { params: { days } });
