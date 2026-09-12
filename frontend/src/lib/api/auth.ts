import { apiClient, http } from './client';
import type { User, UpdateProfileInput } from '$lib/types/user';

export interface LoginResponse {
	user: User;
	session: {
		token: string;
		expires_at: string;
	};
}

export const authApi = {
	/**
	 * Logs in with username/email and password.
	 */
	login: (body: Record<string, any>) =>
		http.post<LoginResponse>('/auth/login', body),

	/**
	 * Registers a new user account.
	 */
	register: (body: Record<string, any>) =>
		http.post<LoginResponse>('/auth/register', body),

	/**
	 * Logs out current session.
	 */
	logout: () => http.post('/auth/logout', {}),

	/**
	 * Fetches current authenticated session and profile.
	 */
	getMe: () => http.get<{ user: User }>('/me'),

	/**
	 * Updates personal profile metadata.
	 */
	updateProfile: (input: UpdateProfileInput) =>
		http.patch<{ user: User }>('/me', input),

	/**
	 * Gets server registration mode (OPEN, INVITE_ONLY, CLOSED).
	 */
	getRegistrationMode: () =>
		http.get<{ registration_mode: string }>('/auth/mode'),

	/**
	 * Validates an invitation token.
	 */
	validateInvitation: (token: string) =>
		http.get<{ valid: boolean; message?: string; invited_email?: string; expires_at?: string }>(`/invitations/validate?token=${encodeURIComponent(token)}`),

	/**
	 * Uploads profile avatar artwork (JPG, PNG, WebP up to 10MB).
	 */
	uploadAvatar: async (file: File): Promise<string> => {
		const formData = new FormData();
		formData.append('file', file);
		formData.append('purpose', 'avatar');
		const res = await apiClient.post<{ url: string }>('/upload', formData, {
			headers: { 'Content-Type': 'multipart/form-data' }
		});
		return res.data.url;
	}
};
