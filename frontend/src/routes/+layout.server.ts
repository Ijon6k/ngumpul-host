import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, url }) => {
	let status: { initialized?: boolean; domain?: string } | null = null;
	try {
		const statusRes = await fetch('/api/setup/status');
		if (statusRes.ok) {
			status = await statusRes.json().catch(() => null);
		}
	} catch {
		// non-blocking
	}

	const initialized = Boolean(status?.initialized);

	if (!initialized && (url.pathname.startsWith('/admin') || url.pathname.startsWith('/me') || url.pathname === '/login')) {
		throw redirect(303, '/setup');
	}

	if (initialized && url.pathname === '/setup') {
		throw redirect(303, '/login');
	}

	return {
		initialized,
		setupDomain: typeof status?.domain === 'string' ? status.domain : null
	};
};