import { redirect, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params, fetch }) => {
	const slug = params.slug;
	if (!slug) {
		throw error(400, 'Project slug is required');
	}

	try {
		const res = await fetch(`/api/projects/${slug}`);
		if (!res.ok) {
			throw error(res.status, 'Project not found');
		}
		const data = await res.json();
		const publicUrl = data?.project?.public_url;
		if (!publicUrl) {
			throw error(404, 'Project has no configured public URL');
		}
		throw redirect(302, publicUrl);
	} catch (e: any) {
		if (e?.status && e?.location) {
			throw e;
		}
		throw error(404, 'Failed to redirect to project');
	}
};
