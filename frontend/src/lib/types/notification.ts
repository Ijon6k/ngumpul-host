export interface AppNotification {
	id: string;
	user_id: string;
	type: string;
	title: string;
	body: string;
	is_read: boolean;
	created_at: string;
}
