import { http } from '../client';
import type { Comment } from '$lib/types/comment';

export const adminCommentsApi = {
	listComments: (query = '') =>
		http.get<{ comments: Comment[] }>(`/admin/comments${query ? `?${query}` : ''}`),

	deleteComment: (id: string) =>
		http.delete(`/admin/comments/${id}`)
};
