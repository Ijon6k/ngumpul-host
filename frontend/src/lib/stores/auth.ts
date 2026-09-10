import { writable } from 'svelte/store';
import { api } from '../api';

export interface User {
	id: string;
	username: string;
	email: string;
	display_name: string;
	avatar_url: string;
	bio: string;
	role: 'USER' | 'ADMIN';
	status: string;
	email_verified: boolean;
	created_at: string;
}

export const user = writable<User | null>(null);
export const authLoading = writable<boolean>(true);

export async function initAuth() {
	try {
		authLoading.set(true);
		const res = await api.get('/me');
		if (res.data?.user) {
			user.set(res.data.user);
		} else {
			user.set(null);
		}
	} catch (_) {
		user.set(null);
	} finally {
		authLoading.set(false);
	}
}

export async function logout() {
	try {
		await api.post('/auth/logout');
	} catch (_) {}
	user.set(null);
	window.location.href = '/';
}
