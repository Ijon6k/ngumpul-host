import { http } from './client';
import type { HostingRequest, CreateHostingRequestInput } from '$lib/types/hosting';

export interface HostingRequestsResponse {
	requests: HostingRequest[];
}

export const hostingRequestsApi = {
	/**
	 * Lists current user's submitted hosting requests.
	 */
	getMyHostingRequests: () =>
		http.get<HostingRequestsResponse>('/me/hosting-requests'),

	/**
	 * Submits a new hosting request for a side project/service.
	 */
	submitHostingRequest: (input: CreateHostingRequestInput) =>
		http.post<{ request: HostingRequest }>('/hosting-requests', input),

	/**
	 * Checks whether a proposed subdomain is available.
	 */
	checkSubdomain: (subdomain: string) =>
		http.get<{ available: boolean; subdomain: string; message: string }>(
			`/hosting-requests/check-subdomain?subdomain=${encodeURIComponent(subdomain)}`
		),

	/**
	 * Requests a subdomain modification for an existing project.
	 */
	requestSubdomainChange: (projectId: string, payload: { new_subdomain: string; reason?: string }) =>
		http.post<{ message: string; request_id: string; subdomain: string }>(
			`/projects/${projectId}/request-subdomain-change`,
			payload
		)
};

