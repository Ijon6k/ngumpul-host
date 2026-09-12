import { writable } from 'svelte/store';
import { notificationsApi } from '$lib/api/notifications';

export const unreadNotificationsCount = writable<number>(0);

/**
 * Fetches the user's notifications and updates the unread badge store.
 */
export async function refreshUnreadNotifications(): Promise<number> {
	try {
		const res = await notificationsApi.getNotifications();
		const list = res?.notifications || [];
		const unread = list.filter((n) => !n.is_read).length;
		unreadNotificationsCount.set(unread);
		return unread;
	} catch {
		return 0;
	}
}

/**
 * Marks a single notification as read and updates the unread store.
 */
export async function markNotificationAsRead(id: string): Promise<boolean> {
	try {
		await notificationsApi.markNotificationRead(id);
		unreadNotificationsCount.update((count) => Math.max(0, count - 1));
		return true;
	} catch (err) {
		console.error('Failed to mark notification as read:', err);
		return false;
	}
}

/**
 * Marks all notifications as read and resets the unread store to 0.
 */
export async function markAllNotificationsAsRead(): Promise<boolean> {
	try {
		await notificationsApi.markAllNotificationsRead();
		unreadNotificationsCount.set(0);
		return true;
	} catch (err) {
		console.error('Failed to mark all notifications as read:', err);
		return false;
	}
}
