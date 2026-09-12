import { http } from './client';
import type { User } from '$lib/types/user';
import type { Project } from '$lib/types/project';

export interface PublicMember extends User {
	project_count?: number;
	projects_count?: number;
}

export interface MembersResponse {
	members: PublicMember[];
}

export interface MemberDetailResponse {
	user: PublicMember;
	projects: Project[];
}

export const usersApi = {
	/**
	 * Lists public community members.
	 */
	getMembers: () => http.get<MembersResponse>('/users'),

	/**
	 * Fetches public user profile and their published projects.
	 */
	getMember: (username: string) => http.get<MemberDetailResponse>(`/users/${encodeURIComponent(username)}`)
};
