/**
 * Shared date and time formatting utilities.
 */

// Formats an ISO string into standard localized date and time (e.g. "Sep 11, 2026, 14:30")
export function formatDate(iso: string): string {
	try {
		const d = new Date(iso);
		return d.toLocaleDateString([], {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	} catch (_) {
		return iso;
	}
}

// Formats an ISO string into clean 24h clock time (e.g. "14:30:15")
export function formatTime(iso: string): string {
	try {
		const d = new Date(iso);
		return d.toLocaleTimeString([], {
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});
	} catch (_) {
		return '';
	}
}

// Formats a duration in seconds into a concise human string (e.g. "45s", "4m 12s", "1h 15m")
export function formatDuration(sec: number): string {
	if (sec < 60) {
		return `${sec}s`;
	}
	const mins = Math.floor(sec / 60);
	const remSec = sec % 60;
	if (mins < 60) {
		return remSec > 0 ? `${mins}m ${remSec}s` : `${mins}m`;
	}
	const hours = Math.floor(mins / 60);
	const remMins = mins % 60;
	return remMins > 0 ? `${hours}h ${remMins}m` : `${hours}h`;
}

/**
 * Resolves semantic color classes for network latency and ping response times.
 * - Low / Fast: < 150ms -> Green
 * - Medium / Fair: 150ms - 400ms -> Amber/Yellow
 * - High / Slow: > 400ms -> Red
 */
export function getPingColorClass(ms: number | null | undefined): string {
	if (ms == null || ms <= 0) {
		return 'text-(--text-muted)';
	}
	if (ms < 150) {
		return 'text-emerald-600 dark:text-emerald-400';
	}
	if (ms <= 400) {
		return 'text-amber-600 dark:text-amber-400';
	}
	return 'text-rose-600 dark:text-rose-400';
}

/**
 * Formats byte size into human readable string (e.g. "145 KB", "2.1 MB")
 */
export function formatBytes(bytes: number): string {
	if (bytes <= 0) return '0 B';
	const k = 1024;
	const sizes = ['B', 'KB', 'MB', 'GB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
}
