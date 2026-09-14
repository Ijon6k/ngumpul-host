import { writable, derived } from 'svelte/store';
import { setupApi } from '$lib/api/setup';
import { browser } from '$app/environment';

export const nodeDomain = writable<string>('ngumpul.local');
export const domainSuffix = derived(nodeDomain, ($d) => ($d ? `.${$d}` : ''));

let refreshPromise: Promise<string> | null = null;

export function refreshNodeDomain(): Promise<string> {
	if (!browser) return Promise.resolve('ngumpul.local');
	if (refreshPromise) return refreshPromise;

	refreshPromise = (async () => {
		try {
			const res = await setupApi.getStatus();
			if (res && typeof res.domain === 'string' && res.domain) {
				const clean = res.domain
					.replace(/^https?:\/\//, '')
					.replace(/:\d+$/, '')
					.trim();
				if (clean) {
					nodeDomain.set(clean);
					return clean;
				}
			}
		} catch (_) {
			// Network error during early init, retain default fallback
		} finally {
			refreshPromise = null;
		}
		return 'ngumpul.local';
	})();

	return refreshPromise;
}

if (browser) {
	refreshNodeDomain();
}
