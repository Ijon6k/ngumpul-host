import { http } from '../client';

export interface RegistrationSettings {
	registration_mode: 'OPEN' | 'INVITE_ONLY' | 'CLOSED';
	domain?: string;
}

export interface Invitation {
	id: string;
	token: string;
	created_by: string;
	max_uses: number;
	uses_count: number;
	expires_at: string;
	created_at: string;
	creator_username?: string;
}

export const adminSettingsApi = {
	getSettings: () =>
		http.get<RegistrationSettings>('/admin/settings'),

	updateSettings: (body: { registration_mode?: string; domain?: string }) =>
		http.patch<RegistrationSettings>('/admin/settings', body),

	listInvitations: () =>
		http.get<{ invitations: Invitation[] }>('/admin/invitations'),

	createInvitation: (body: { max_uses?: number; expires_in_days?: number; invited_email?: string }) =>
		http.post<{ raw_token?: string; invite_url?: string; invitation: Invitation }>('/admin/invitations', body),

	revokeInvitation: (id: string) =>
		http.post(`/admin/invitations/${id}/revoke`, {})
};
