import { http } from './client';
import type { ActivityEvent } from '$lib/types/activity';

export interface ActivityListResponse {
	activities: ActivityEvent[];
	total?: number;
	page?: number;
	limit?: number;
	total_pages?: number;
}

export const activityApi = {
	/**
	 * Lists public recent activity events on the server with pagination.
	 */
	getPublicActivity: (page = 1, limit = 20) =>
		http.get<ActivityListResponse>(`/activity?page=${page}&limit=${limit}`),

	/**
	 * Lists authenticated user's personal activity history with pagination.
	 */
	getMyActivity: (page = 1, limit = 15) =>
		http.get<ActivityListResponse>(`/me/activity?page=${page}&limit=${limit}`)
};
