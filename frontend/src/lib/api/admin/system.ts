import { http } from '../client';

export const adminSystemApi = {
	getSystemHealth: () => http.get('/admin/system')
};
