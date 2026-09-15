import { http } from './client';
import type { User } from '$lib/types/user';

export interface SetupStatusResponse {
	initialized: boolean;
	domain: string;
}

export interface SetupPayload {
	name: string;
	username: string;
	email: string;
	password: string;
	domain: string;
	subdomain?: string;
}

export interface SetupResponse {
	message: string;
	user: User;
}

export const setupApi = {
	getStatus: () => http.get<SetupStatusResponse>('/setup/status'),
	setup: (payload: SetupPayload) => http.post<SetupResponse>('/setup', payload)
};
