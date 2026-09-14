import { superValidate, message } from 'sveltekit-superforms';
import { valibot } from 'sveltekit-superforms/adapters';
import { fail } from '@sveltejs/kit';
import { loginSchema } from '$lib/schemas/auth';
import { forwardSessionCookie } from '$lib/server/session';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async () => {
	const form = await superValidate(valibot(loginSchema));
	return { form };
};

export const actions: Actions = {
	default: async ({ request, fetch, cookies }) => {
		const form = await superValidate(request, valibot(loginSchema));
		if (!form.valid) return fail(400, { form });

		const res = await fetch('/api/auth/login', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				email_or_username: form.data.emailOrUsername,
				password: form.data.password
			})
		});

		const body = await res.json().catch(() => null);
		if (!res.ok) {
			return message(form, body?.error ?? 'Invalid credentials.');
		}

		forwardSessionCookie(cookies, res);

		return { form, user: body.user };
	}
};
