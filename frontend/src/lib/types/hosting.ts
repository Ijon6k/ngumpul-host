export interface HostingRequest {
	id: string;
	user_id: string;
	project_id?: string;
	request_type?: 'NEW_PROJECT' | 'SUBDOMAIN_CHANGE';
	project_name: string;
	subdomain?: string;
	description: string;
	readme?: string;
	cover_image_url?: string;
	repository_url?: string;
	documentation_url?: string;
	demo_url?: string;
	deployment_notes?: string;
	technology_stack?: string[];
	status: 'PENDING' | 'APPROVED' | 'REJECTED' | 'COMPLETED' | 'SETUP';
	operator_notes?: string;
	admin_notes?: string;
	reviewed_by?: string;
	reviewed_at?: string;
	created_at: string;
	updated_at: string;
	applicant_username?: string;
	applicant_name?: string;
	requester?: {
		display_name?: string;
		username?: string;
		avatar_url?: string;
	};
	environment_specs?: Record<string, any>;
}

export interface CreateHostingRequestInput {
	project_name: string;
	subdomain?: string;
	description: string;
	readme?: string;
	cover_image_url?: string;
	repository_url?: string;
	documentation_url?: string;
	demo_url?: string;
	deployment_notes?: string;
	technology_stack?: string[];
}
