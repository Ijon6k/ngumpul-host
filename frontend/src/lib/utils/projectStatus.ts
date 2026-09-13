export type ProjectDisplayStatusKey =
	| 'online'
	| 'unreachable'
	| 'checking'
	| 'suspended'
	| 'setup'
	| 'archived';

export interface ResolvedProjectStatus {
	key: ProjectDisplayStatusKey;
	label: string;
	badgeVariant: 'success' | 'warning' | 'error' | 'neutral';
	dotStatus: string;
	description: string;
	colorClass: string;
	bgClass: string;
	borderClass: string;
}

export interface ProjectStatusInput {
	lifecycle_status?: string | null;
	availability?: string | null;
	status?: string | null;
	visibility?: string | null;
}

/**
 * Centrally resolves primary display status according to authoritative hierarchy:
 * 1. SUSPENDED (intentional administrative suspension, never downtime)
 * 2. ARCHIVED (retired project)
 * 3. SETUP (provisioning/preparation)
 * 4. ACTIVE + availability:
 *    - REACHABLE -> Online
 *    - UNREACHABLE -> Unreachable
 *    - UNKNOWN -> Checking / Unknown
 */
export function resolveProjectStatus(input?: ProjectStatusInput | null): ResolvedProjectStatus {
	if (!input) {
		return {
			key: 'checking',
			label: 'Unknown',
			badgeVariant: 'neutral',
			dotStatus: 'checking',
			description: 'Status unavailable',
			colorClass: 'text-(--text-muted)',
			bgClass: 'bg-(--bg-muted)',
			borderClass: 'border-(--border-hairline)'
		};
	}

	const lifecycle = (input.lifecycle_status || '').toUpperCase();
	const availability = (input.availability || '').toUpperCase();
	const legacyStatus = (input.status || '').toUpperCase();

	// Priority 1: SUSPENDED
	if (lifecycle === 'SUSPENDED' || legacyStatus === 'SUSPENDED') {
		return {
			key: 'suspended',
			label: 'Suspended',
			badgeVariant: 'warning',
			dotStatus: 'suspended',
			description: 'Administrative suspension',
			colorClass: 'text-amber-600 dark:text-amber-400',
			bgClass: 'bg-amber-500/10',
			borderClass: 'border-amber-500/30'
		};
	}

	// Priority 2: ARCHIVED
	if (lifecycle === 'ARCHIVED' || legacyStatus === 'ARCHIVED') {
		return {
			key: 'archived',
			label: 'Archived',
			badgeVariant: 'neutral',
			dotStatus: 'archived',
			description: 'Archived project',
			colorClass: 'text-(--text-muted)',
			bgClass: 'bg-(--bg-muted)',
			borderClass: 'border-(--border-hairline)'
		};
	}

	// Priority 3: SETUP
	if (lifecycle === 'SETUP' || legacyStatus === 'SETUP' || legacyStatus === 'PENDING') {
		return {
			key: 'setup',
			label: 'Setting up',
			badgeVariant: 'neutral',
			dotStatus: 'setup',
			description: 'Infrastructure setup in progress',
			colorClass: 'text-(--accent-sky)',
			bgClass: 'bg-(--accent-sky)/10',
			borderClass: 'border-(--accent-sky)/30'
		};
	}

	// Priority 4: ACTIVE + availability
	if (availability === 'REACHABLE' || (availability === '' && legacyStatus === 'ONLINE')) {
		return {
			key: 'online',
			label: 'Online',
			badgeVariant: 'success',
			dotStatus: 'online',
			description: 'Operational and reachable',
			colorClass: 'text-emerald-600 dark:text-emerald-400',
			bgClass: 'bg-emerald-500/10',
			borderClass: 'border-emerald-500/30'
		};
	}

	if (availability === 'UNREACHABLE' || (availability === '' && legacyStatus === 'OFFLINE')) {
		return {
			key: 'unreachable',
			label: 'Unreachable',
			badgeVariant: 'error',
			dotStatus: 'unreachable',
			description: 'Endpoint cannot be reached',
			colorClass: 'text-rose-600 dark:text-rose-400',
			bgClass: 'bg-rose-500/10',
			borderClass: 'border-rose-500/30'
		};
	}

	// Default fallback: ACTIVE + UNKNOWN
	return {
		key: 'checking',
		label: 'Checking',
		badgeVariant: 'warning',
		dotStatus: 'checking',
		description: 'Awaiting availability probe',
		colorClass: 'text-amber-600 dark:text-amber-400',
		bgClass: 'bg-amber-500/10',
		borderClass: 'border-amber-500/30'
	};
}
