import { apiClient, http } from './client';
import type { Project, ProjectAvailability, ProjectVisits, UpdateProjectInput } from '$lib/types/project';

export interface ProjectDetailResponse {
	project: Project;
	availability: ProjectAvailability;
	activities: any[];
}

export interface ProjectsListResponse {
	projects: Project[];
	total?: number;
	page?: number;
	limit?: number;
	total_pages?: number;
}

export const projectsApi = {
	/**
	 * Lists public projects for the showcase catalog.
	 */
	getProjects: (params?: { page?: number; limit?: number; q?: string; type?: string }) =>
		http.get<ProjectsListResponse>('/projects', { params }),

	/**
	 * Fetches public project detail by slug with availability and timeline.
	 */
	getProject: (slug: string) => http.get<ProjectDetailResponse>(`/projects/${slug}`),

	/**
	 * Lists current user's projects.
	 */
	getMyProjects: () => http.get<{ projects: Project[] }>('/me/projects'),

	/**
	 * Gets project detail owned by current user.
	 */
	getMyProject: (id: string) => http.get<ProjectDetailResponse>(`/me/projects/${id}`),

	/**
	 * Gets 30-day visit and page view analytics for a project.
	 */
	getMyProjectVisits: (id: string) => http.get<ProjectVisits>(`/me/projects/${id}/visits`),

	/**
	 * Updates personal project metadata and links.
	 */
	updateMyProject: (id: string, input: UpdateProjectInput) =>
		http.patch<{ project: Project }>(`/me/projects/${id}`, input),

	/**
	 * Uploads cover image artwork (JPG, PNG, WebP up to 10MB).
	 */
	uploadCover: async (file: File): Promise<string> => {
		const formData = new FormData();
		formData.append('file', file);
		formData.append('purpose', 'cover');
		const res = await apiClient.post<{ url: string }>('/upload', formData, {
			headers: { 'Content-Type': 'multipart/form-data' }
		});
		return res.data.url;
	}
};
