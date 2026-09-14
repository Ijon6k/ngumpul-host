import { superValidate, setError, message } from 'sveltekit-superforms';
import { valibot } from 'sveltekit-superforms/adapters';
import { fail, redirect } from '@sveltejs/kit';
import { setupSchema } from '$lib/schemas/auth';
import { forwardSessionCookie } from '$lib/server/session';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const form = await superValidate(valibot(setupSchema));

	const statusRes = await fetch('/api/setup/status');
	if (statusRes.ok) {
		const status = (await statusRes.json()) as { initialized?: boolean };
		if (status.initialized) throw redirect(302, '/login');
	}

	return { form };
};

export const actions: Actions = {
	default: async ({ request, fetch, cookies }) => {
		const form = await superValidate(request, valibot(setupSchema));
		if (!form.valid) return fail(400, { form });

		if (form.data.password !== form.data.confirmPassword) {
			return setError(form, 'confirmPassword', 'Passwords do not match.');
		}

		const res = await fetch('/api/setup', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				domain: form.data.domain.trim(),
				name: form.data.name.trim(),
				username: form.data.username.trim().toLowerCase(),
				email: form.data.email.trim().toLowerCase(),
				password: form.data.password
			})
		});

		const body = await res.json().catch(() => null);
		if (!res.ok) {
			if (res.status === 403 && body?.error?.includes('already been completed')) {
				throw redirect(302, '/login');
			}
			return message(form, body?.error ?? 'Node setup failed.', {
				status: res.status as NonNullable<Parameters<typeof message>[2]>["status"]
			});
		}

		forwardSessionCookie(cookies, res);

		return { form, user: body.user };
	}
};
