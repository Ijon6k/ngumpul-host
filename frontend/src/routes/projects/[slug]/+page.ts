import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { Project, ProjectAvailability } from '$lib/types/project';
import type { ActivityEvent } from '$lib/types/activity';

export const ssr = false;

export const load: PageLoad = async ({ params, fetch }) => {
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
		return {
			project: data.project as Project,
			availability: data.availability as ProjectAvailability,
			activities: (data.activities || []) as ActivityEvent[]
		};
	} catch (e: any) {
		if (e?.status) throw e;
		throw error(404, 'Failed to load project details');
	}
};
