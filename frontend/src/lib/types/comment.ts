import type { PublicUser } from '$lib/types/user';

export interface Comment {
	id: string;
	project_id: string;
	project_slug?: string;
	project_name?: string;
	author_id: string;
	content: string;
	created_at: string;
	updated_at: string;
	deleted_at?: string | null;
	is_deleted?: boolean;
	report_count?: number;
	author?: PublicUser;
}

export interface CreateCommentInput {
	content: string;
}
