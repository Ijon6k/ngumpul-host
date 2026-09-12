import { http } from './client';
import type { AppNotification } from '$lib/types/notification';

export interface NotificationsListResponse {
	notifications: AppNotification[];
}

export const notificationsApi = {
	/**
	 * Lists current user's notifications.
	 */
	getNotifications: () =>
		http.get<NotificationsListResponse>('/me/notifications'),

	/**
	 * Marks a notification as read.
	 */
	markNotificationRead: (id: string) =>
		http.patch(`/me/notifications/${id}/read`, {}),

	/**
	 * Marks all user notifications as read.
	 */
	markAllNotificationsRead: () =>
		http.post('/me/notifications/read-all', {})
};
