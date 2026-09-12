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
		http.post<{ request: HostingRequest }>('/hosting-requests', input)
};
