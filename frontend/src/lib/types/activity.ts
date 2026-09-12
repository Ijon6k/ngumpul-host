export interface ActivityEvent {
	id: string;
	type: string;
	actor_id?: string;
	target_type?: string;
	target_id?: string;
	metadata: Record<string, any>;
	created_at: string;
	actor_username?: string;
	actor_name?: string;
	project_name?: string;
	project_slug?: string;
}
