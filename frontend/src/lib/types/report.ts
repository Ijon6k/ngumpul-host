export interface Report {
	id: string;
	reporter_id: string;
	target_type: 'PROJECT' | 'COMMENT';
	target_id: string;
	target_slug?: string;
	target_name?: string;
	target_content?: string;
	comment_content?: string;
	reason: string;
	details?: string;
	status: 'OPEN' | 'RESOLVED' | 'DISMISSED' | 'REVIEWED';
	resolution_notes?: string;
	resolved_by?: string;
	resolved_by_name?: string;
	resolved_at?: string;
	created_at: string;
	reporter_username?: string;
	reporter_name?: string;
}

export interface CreateReportInput {
	reason: string;
	details?: string;
}
