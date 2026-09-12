import { http } from '../client';

export interface AdminStats {
	total_projects: number;
	online_projects: number;
	total_users: number;
	total_members?: number;
	pending_requests: number;
	open_reports: number;
}

export const adminStatsApi = {
	getStats: () => http.get<AdminStats>('/admin/stats')
};
