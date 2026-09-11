import { derived, type Readable } from 'svelte/store';
import { createQuery, type CreateQueryResult } from '@tanstack/svelte-query';
import { apiClient } from './client';
import type { PublicStatusResponse } from '$lib/types/status';

/**
 * Fetch public telemetry and status for a given time window.
 */
export async function fetchPublicStatus(range: '1d' | '7d' | '30d'): Promise<PublicStatusResponse> {
	const res = await apiClient.get<PublicStatusResponse>(`/status?range=${range}`);
	return res.data;
}

/**
 * TanStack Svelte Query store hook with 60s background polling.
 * Accepts a reactive Svelte store of the selected time range.
 */
export function createStatusQuery(
	rangeStore: Readable<'1d' | '7d' | '30d'>
): CreateQueryResult<PublicStatusResponse, Error> {
	return createQuery<PublicStatusResponse, Error>(
		derived(rangeStore, ($range) => ({
			queryKey: ['status', $range],
			queryFn: () => fetchPublicStatus($range),
			refetchInterval: 60000,
			staleTime: 30000
		}))
	);
}

