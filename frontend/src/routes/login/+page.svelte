<script lang="ts">
	import { untrack } from 'svelte';
	import { superForm } from 'sveltekit-superforms';
	import { valibotClient } from 'sveltekit-superforms/adapters';
	import { loginSchema } from '$lib/schemas/auth';
	import { user } from '$lib/stores/auth';
	import { AuthSplitLayout } from '$lib/components/layout';
	import { Alert, Input } from '$lib/components/ui';
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const { form, errors, constraints, message, submitting, enhance } = superForm(
		untrack(() => data.form),
		{ validators: valibotClient(loginSchema), resetForm: false, onResult: handleResult }
	);

	function handleResult({ result }: { result: { type: string; data?: any } }) {
		if (result.type === 'success' && result.data?.user) {
			user.set(result.data.user);
			goto(result.data.user.role === 'ADMIN' ? '/admin' : '/me');
		}
	}

	function fieldError(errors: Record<string, any> | undefined, field: string): string {
		const val = errors?.[field];
		return Array.isArray(val) && val.length > 0 ? val[0] : '';
	}
</script>

<svelte:head>
	<title>Sign in · Ngumpul Host</title>
</svelte:head>

<AuthSplitLayout
	title="Sign in"
	description="Authenticate to your personal workspace or operator console."
>
	{#if $message}
		<Alert variant="danger">{$message}</Alert>
	{/if}

	{#if $errors._errors}
		<Alert variant="danger">{$errors._errors[0]}</Alert>
	{/if}

	<form method="POST" use:enhance class="flex flex-col gap-4">
		<Input
			id="login-id"
			name="emailOrUsername"
			label="Email or Username"
			type="text"
			placeholder="operator or handle"
			autocomplete="username"
			inputClass="h-10 px-3.5 bg-(--bg-muted)"
			error={fieldError($errors, 'emailOrUsername')}
			bind:value={$form.emailOrUsername}
			{...$constraints.emailOrUsername}
		/>

		<Input
			id="login-pass"
			name="password"
			label="Password"
			type="password"
			placeholder="••••••••"
			autocomplete="current-password"
			inputClass="h-10 px-3.5 bg-(--bg-muted)"
			error={fieldError($errors, 'password')}
			bind:value={$form.password}
			{...$constraints.password}
		/>

		<button
			type="submit"
			class="btn btn-primary w-full py-2.5 mt-2 text-sm font-medium"
			disabled={$submitting}
		>
			{$submitting ? 'Authenticating...' : 'Sign in'}
		</button>
	</form>

	<svelte:fragment slot="footer">
		<footer class="flex items-center justify-between text-sm text-(--text-muted) pt-4 border-t border-(--border-hairline)">
			<span>Need an account?</span>
			<a href="/register" class="text-(--accent-strong) hover:underline font-medium">Create profile</a>
		</footer>
	</svelte:fragment>
</AuthSplitLayout>