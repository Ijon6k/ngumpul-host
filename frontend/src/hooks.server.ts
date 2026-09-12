import type { HandleFetch } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const handleFetch: HandleFetch = async ({ request, fetch }) => {
	const url = new URL(request.url);
	if (url.pathname.startsWith('/api/')) {
		const backendHost = env.BACKEND_URL || 'http://backend:8080';
		const targetUrl = new URL(url.pathname + url.search, backendHost);
		const newRequest = new Request(targetUrl.toString(), request);
		return fetch(newRequest);
	}
	return fetch(request);
};
