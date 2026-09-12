import { http } from './client';
import type { Comment, CreateCommentInput } from '$lib/types/comment';

export interface CommentsListResponse {
	comments: Comment[];
}

export const commentsApi = {
	/**
	 * Lists public comments for a project by slug.
	 */
	getProjectComments: (slug: string) =>
		http.get<CommentsListResponse>(`/projects/${slug}/comments`),

	getComments: (slug: string) =>
		http.get<CommentsListResponse>(`/projects/${slug}/comments`),

	/**
	 * Adds a comment to a project.
	 */
	createComment: (slug: string, input: CreateCommentInput) =>
		http.post<{ comment: Comment }>(`/projects/${slug}/comments`, input),

	/**
	 * Deletes a comment owned by user or moderated.
	 */
	deleteComment: (commentId: string) =>
		http.delete(`/comments/${commentId}`)
};
