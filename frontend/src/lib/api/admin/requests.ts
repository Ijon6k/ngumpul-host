import { http } from '../client';
import type { HostingRequest } from '$lib/types/hosting';

export const adminRequestsApi = {
	listRequests: (status = 'ALL') =>
		http.get<{ requests: HostingRequest[] }>(`/admin/hosting-requests${status !== 'ALL' ? `?status=${status}` : ''}`),

	approveRequest: (id: string, body?: Record<string, any>) =>
		http.post(`/admin/hosting-requests/${id}/approve`, body || {}),

	rejectRequest: (id: string, body?: Record<string, any>) =>
		http.post(`/admin/hosting-requests/${id}/reject`, body || {}),

	completeRequest: (id: string, body?: Record<string, any>) =>
		http.post(`/admin/hosting-requests/${id}/complete`, body || {})
};
