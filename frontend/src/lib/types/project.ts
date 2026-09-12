import type { PublicUser } from '$lib/types/user';

export interface Project {
	id: string;
	owner_id: string;
	name: string;
	slug: string;
	description: string;
	readme?: string;
	cover_image_url?: string;
	repository_url?: string;
	documentation_url?: string;
	demo_url?: string;
	technology_stack?: string[];
	hosting_type: string;
	public_url?: string;
	status: string;
	visibility: string;
	created_at: string;
	updated_at: string;
	published_at?: string;
	owner?: PublicUser;
}

export interface ProjectAvailability {
	project_id: string;
	uptime_percent: number;
	total_checks: number;
	successful_checks: number;
	latest_response_time_ms?: number;
	last_checked_at?: string;
	last_status?: string;
}

export interface ProjectVisits {
	total_page_views: number;
	total_visits: number;
	daily: Array<{
		date: string;
		page_views: number;
		visits: number;
	}>;
}

export interface UpdateProjectInput {
	name?: string;
	description?: string;
	readme?: string;
	cover_image_url?: string;
	repository_url?: string;
	documentation_url?: string;
	demo_url?: string;
	technology_stack?: string[];
}
