import { http } from '../client';
import type { User } from '$lib/types/user';

export const adminUsersApi = {
	listUsers: (query = '') =>
		http.get<{ users: User[] }>(`/admin/users${query ? `?${query}` : ''}`),

	updateUserRole: (id: string, role: string) =>
		http.patch(`/admin/users/${id}/role`, { role }),

	updateUserStatus: (id: string, status: string) =>
		http.patch(`/admin/users/${id}/status`, { status })
};
