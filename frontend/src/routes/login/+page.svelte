<script lang="ts">
	import { createForm } from '$lib/utils/form.svelte';
	import { loginSchema } from '$lib/schemas/auth';
	import { authApi } from '$lib/api';
	import { user } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { AuthSplitLayout } from '$lib/components/layout';
	import { Alert, Input, Button } from '$lib/components/ui';

	const form = createForm({
		schema: loginSchema,
		initialValues: { emailOrUsername: '', password: '' },
		onSubmit: async (values) => {
			const res = await authApi.login({
				email_or_username: values.emailOrUsername,
				password: values.password
			});
			if (res?.user) {
				user.set(res.user);
				await goto(res.user.role === 'ADMIN' ? '/admin' : '/me');
			}
		}
	});
</script>

<svelte:head>
	<title>Sign in · Ngumpul Host</title>
</svelte:head>

<AuthSplitLayout
	title="Sign in"
	description="Authenticate to your personal workspace or operator console."
>
	{#if form.serverError}
		<Alert variant="danger">{form.serverError}</Alert>
	{/if}

	<form onsubmit={form.handleSubmit} class="flex flex-col gap-4">
		<Input
			id="login-id"
			name="emailOrUsername"
			label="Email or Username"
			type="text"
			placeholder="operator or handle"
			autocomplete="username"
			inputClass="h-10 px-3.5 bg-(--bg-muted)"
			error={form.errors.emailOrUsername}
			bind:value={form.values.emailOrUsername}
		/>

		<Input
			id="login-pass"
			name="password"
			label="Password"
			type="password"
			placeholder="••••••••"
			autocomplete="current-password"
			inputClass="h-10 px-3.5 bg-(--bg-muted)"
			error={form.errors.password}
			bind:value={form.values.password}
		/>

		<Button
			type="submit"
			variant="primary"
			size="md"
			class="w-full mt-2 font-medium"
			loading={form.isSubmitting}
			disabled={form.isSubmitting}
		>
			{form.isSubmitting ? 'Authenticating...' : 'Sign in'}
		</Button>
	</form>

	<svelte:fragment slot="footer">
		<footer class="flex items-center justify-between text-sm text-(--text-muted) pt-4 border-t border-(--border-hairline)">
			<span>Need an account?</span>
			<a href="/register" class="text-(--accent-strong) hover:underline font-medium">Create profile</a>
		</footer>
	</svelte:fragment>
</AuthSplitLayout>