export interface AppNotification {
	id: string;
	user_id: string;
	type: string;
	title: string;
	body: string;
	read_at?: string | null;
	is_read: boolean;
	created_at: string;
}
