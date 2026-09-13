<script lang="ts">
	import { authApi, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';
	import { AuthSplitLayout } from '$lib/components/layout';

	let emailOrUsername = '';
	let password = '';
	let loading = false;
	let error: string | null = null;

	async function handleSubmit() {
		error = null;
		loading = true;
		try {
			const res = await authApi.login({
				email_or_username: emailOrUsername,
				password
			});
			if (res?.user) {
				user.set(res.user);
				window.location.href = res.user.role === 'ADMIN' ? '/admin' : '/me';
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Sign in · Ngumpul Host</title>
</svelte:head>

<AuthSplitLayout
	title="Sign in"
	description="Authenticate to your personal workspace or operator console."
>
	{#if error}
		<div class="p-3.5 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-sm leading-relaxed">
			{error}
		</div>
	{/if}

	<form on:submit|preventDefault={handleSubmit} class="flex flex-col gap-4">
		<div class="flex flex-col gap-1.5">
			<label for="login-id" class="text-sm font-medium text-(--text-main)">Email or Username</label>
			<input
				id="login-id"
				type="text"
				class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50"
				placeholder="operator or handle"
				required
				autocomplete="username"
				bind:value={emailOrUsername}
			/>
		</div>

		<div class="flex flex-col gap-1.5">
			<label for="login-pass" class="text-sm font-medium text-(--text-main)">Password</label>
			<input
				id="login-pass"
				type="password"
				class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
				placeholder="••••••••"
				required
				autocomplete="current-password"
				bind:value={password}
			/>
		</div>

		<button type="submit" class="btn btn-primary w-full py-2.5 mt-2 text-sm font-medium" disabled={loading}>
			{loading ? 'Authenticating...' : 'Sign in'}
		</button>
	</form>

	<svelte:fragment slot="footer">
		<footer class="flex items-center justify-between text-sm text-(--text-muted) pt-4 border-t border-(--border-hairline)">
			<span>Need an account?</span>
			<a href="/register" class="text-(--accent-strong) hover:underline font-medium">Create profile</a>
		</footer>
	</svelte:fragment>
</AuthSplitLayout>
