import { http } from '../client';
import type { Report } from '$lib/types/report';

export const adminReportsApi = {
	listReports: (status = 'ALL') =>
		http.get<{ reports: Report[] }>(`/admin/reports${status !== 'ALL' ? `?status=${status}` : ''}`),

	resolveReport: (id: string, notes?: string) =>
		http.post(`/admin/reports/${id}/resolve`, { resolution_notes: notes }),

	dismissReport: (id: string, notes?: string) =>
		http.post(`/admin/reports/${id}/dismiss`, { resolution_notes: notes }),

	patchReport: (id: string, body: Record<string, any>) =>
		http.patch(`/admin/reports/${id}`, body)
};
