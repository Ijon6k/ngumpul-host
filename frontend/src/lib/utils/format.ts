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
