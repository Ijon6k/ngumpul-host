import { http } from './client';
import type { CreateReportInput } from '$lib/types/report';

export const reportsApi = {
	/**
	 * Submits an abuse/spam report for a project.
	 */
	reportProject: (slug: string, input: CreateReportInput) =>
		http.post(`/projects/${slug}/report`, input),

	/**
	 * Submits an abuse/harassment report for a comment.
	 */
	reportComment: (commentId: string, input: CreateReportInput) =>
		http.post(`/comments/${commentId}/report`, input)
};
