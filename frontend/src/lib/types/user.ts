export interface User {
	id: string;
	username: string;
	email: string;
	display_name: string;
	bio?: string;
	avatar_url?: string;
	role: 'ADMIN' | 'USER' | 'MEMBER';
	status: 'ACTIVE' | 'SUSPENDED';
	email_verified?: boolean;
	projects_count?: number;
	created_at: string;
	updated_at?: string;
}

export interface PublicUser {
	id: string;
	username: string;
	display_name: string;
	bio?: string;
	avatar_url?: string;
	created_at: string;
}

export interface UpdateProfileInput {
	display_name?: string;
	bio?: string;
	avatar_url?: string;
}
