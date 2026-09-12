import { http } from './client';
import type { ActivityEvent } from '$lib/types/activity';

export interface ActivityListResponse {
	activities: ActivityEvent[];
}

export const activityApi = {
	/**
	 * Lists public recent activity events on the server.
	 */
	getPublicActivity: () => http.get<ActivityListResponse>('/activity'),

	/**
	 * Lists authenticated user's personal activity history.
	 */
	getMyActivity: () => http.get<ActivityListResponse>('/me/activity')
};
