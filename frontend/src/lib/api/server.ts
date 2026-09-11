import { createQuery, type CreateQueryResult } from '@tanstack/svelte-query';
import { apiClient } from './client';
import type { PublicServerResponse } from '$lib/types/server';

/**
 * Fetch dynamic host specifications and telemetry for Bento and public showcases.
 */
export async function fetchPublicServer(): Promise<PublicServerResponse> {
	const res = await apiClient.get<PublicServerResponse>('/public/server');
	return res.data;
}

/**
 * TanStack Svelte Query store hook for server telemetry.
 */
export function createPublicServerQuery(): CreateQueryResult<PublicServerResponse, Error> {
	return createQuery<PublicServerResponse, Error>({
		queryKey: ['publicServer'],
		queryFn: fetchPublicServer,
		staleTime: 60000,
		refetchInterval: 120000
	});
}

