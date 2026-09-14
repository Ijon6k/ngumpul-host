import { superValidate, setError } from 'sveltekit-superforms';
import { valibot } from 'sveltekit-superforms/adapters';
import { fail } from '@sveltejs/kit';
import { registerSchema } from '$lib/schemas/auth';
import { forwardSessionCookie } from '$lib/server/session';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async () => {
	const form = await superValidate(valibot(registerSchema));
	return { form };
};

export const actions: Actions = {
	default: async ({ request, fetch, cookies }) => {
		const form = await superValidate(request, valibot(registerSchema));
		if (!form.valid) return fail(400, { form });

		if (form.data.password !== form.data.confirmPassword) {
			return fail(400, {
				form: setError(form, 'confirmPassword', 'Passwords do not match.')
			});
		}

		const res = await fetch('/api/auth/register', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				username: form.data.username.trim().toLowerCase(),
				display_name: form.data.displayName.trim() || form.data.username.trim(),
				email: form.data.email.trim().toLowerCase(),
				password: form.data.password,
				invitation_token: form.data.invitationToken.trim() || undefined
			})
		});

		const body = await res.json().catch(() => null);
		if (!res.ok) {
			return fail(res.status, {
				form: setError(form, '', body?.error ?? 'Registration failed.')
			});
		}

		forwardSessionCookie(cookies, res);

		return { form, user: body.user };
	}
};
