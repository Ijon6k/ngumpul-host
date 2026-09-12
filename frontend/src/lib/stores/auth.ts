import { writable } from 'svelte/store';
import { authApi } from '../api/auth';
import type { User } from '$lib/types/user';

export type { User };

export const user = writable<User | null>(null);
export const authLoading = writable<boolean>(true);

export async function initAuth() {
	try {
		authLoading.set(true);
		const res = await authApi.getMe();
		if (res?.user) {
			user.set(res.user);
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
		await authApi.logout();
	} catch (_) {}
	user.set(null);
	window.location.href = '/';
}
