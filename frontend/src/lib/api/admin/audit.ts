import { http } from '../client';

export interface AuditEntry {
	id: string;
	actor_id: string;
	action: string;
	target_type: string;
	target_id: string;
	metadata: Record<string, any>;
	created_at: string;
	actor_name: string;
	actor_username: string;
}

export const adminAuditApi = {
	listAuditLogs: (query = '') =>
		http.get<{ audit_logs: AuditEntry[] }>(`/admin/audit${query ? `?${query}` : ''}`)
};
