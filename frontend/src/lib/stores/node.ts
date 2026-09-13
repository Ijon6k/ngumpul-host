import { writable, derived } from 'svelte/store';
import { setupApi } from '$lib/api/setup';
import { browser } from '$app/environment';

export const nodeDomain = writable<string>('ngumpul.local');
export const domainSuffix = derived(nodeDomain, ($d) => ($d ? `.${$d}` : ''));

export async function refreshNodeDomain(): Promise<string> {
	if (!browser) return 'ngumpul.local';
	try {
		const res = await setupApi.getStatus();
		if (res?.domain) {
			const clean = res.domain
				.replace(/^https?:\/\//, '')
				.replace(/:\d+$/, '')
				.trim();
			if (clean) {
				nodeDomain.set(clean);
				return clean;
			}
		}
	} catch (_) {}
	return 'ngumpul.local';
}
