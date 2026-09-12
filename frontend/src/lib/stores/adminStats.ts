import { writable } from 'svelte/store';
import { adminStatsApi } from '$lib/api/admin/stats';

export const adminPendingCount = writable<number>(0);
export const adminOpenReportsCount = writable<number>(0);

/**
 * Refreshes admin statistics (pending hosting requests, open reports)
 */
export async function refreshAdminStats(): Promise<void> {
	try {
		const res = await adminStatsApi.getStats();
		if (res?.pending_requests != null) {
			adminPendingCount.set(res.pending_requests);
		}
		if (res?.open_reports != null) {
			adminOpenReportsCount.set(res.open_reports);
		}
	} catch {
		// Silently ignore if admin stats cannot be fetched
	}
}
