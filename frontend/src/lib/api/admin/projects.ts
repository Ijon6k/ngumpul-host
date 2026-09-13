import { http } from '../client';
import type { Project } from '$lib/types/project';

export const adminProjectsApi = {
	listProjects: (query = '') =>
		http.get<{ projects: Project[] }>(`/admin/projects${query ? `?${query}` : ''}`),

	getProject: (id: string) =>
		http.get<{ project: Project; availability?: any; activities?: any[]; comments_count?: number }>(`/admin/projects/${id}`),

	createProject: (body: Record<string, any>) =>
		http.post<{ project: Project }>('/admin/projects', body),

	updateProject: (id: string, body: Record<string, any>) =>
		http.patch<{ project: Project }>(`/admin/projects/${id}`, body),

	suspendProject: (id: string) =>
		http.post<{ message: string }>(`/admin/projects/${id}/suspend`, {}),

	restoreProject: (id: string) =>
		http.post<{ message: string }>(`/admin/projects/${id}/restore`, {}),

	archiveProject: (id: string) =>
		http.post<{ message: string }>(`/admin/projects/${id}/archive`, {}),

	activateProject: (id: string) =>
		http.post<{ message: string }>(`/admin/projects/${id}/activate`, {})
};
