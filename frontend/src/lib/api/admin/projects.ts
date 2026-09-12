import { http } from '../client';
import type { Project } from '$lib/types/project';

export const adminProjectsApi = {
	listProjects: (query = '') =>
		http.get<{ projects: Project[] }>(`/admin/projects${query ? `?${query}` : ''}`),

	createProject: (body: Record<string, any>) =>
		http.post<{ project: Project }>('/admin/projects', body),

	updateProject: (id: string, body: Record<string, any>) =>
		http.patch<{ project: Project }>(`/admin/projects/${id}`, body)
};
