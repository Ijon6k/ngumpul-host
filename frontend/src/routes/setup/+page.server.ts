import { superValidate, setError } from 'sveltekit-superforms';
import { valibot } from 'sveltekit-superforms/adapters';
import { fail } from '@sveltejs/kit';
import { setupSchema } from '$lib/schemas/auth';
import { forwardSessionCookie } from '$lib/server/session';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async () => {
	const form = await superValidate(valibot(setupSchema));
	return { form };
};

export const actions: Actions = {
	default: async ({ request, fetch, cookies }) => {
		const form = await superValidate(request, valibot(setupSchema));
		if (!form.valid) return fail(400, { form });

		if (form.data.password !== form.data.confirmPassword) {
			return fail(400, {
				form: setError(form, 'confirmPassword', 'Passwords do not match.')
			});
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
			return fail(res.status, {
				form: setError(form, '', body?.error ?? 'Node setup failed.')
			});
		}

		forwardSessionCookie(cookies, res);

		return { form, user: body.user };
	}
};
