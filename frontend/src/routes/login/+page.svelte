<script lang="ts">
	import { authApi, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';

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

<div class="container mx-auto px-6 py-20 flex justify-center items-center">
	<div class="w-full max-w-md bg-(--bg-surface) border border-(--border-hairline) rounded-md p-7 sm:p-8 flex flex-col gap-6">
		<header class="text-left">
			<h1 class="font-display font-bold text-2xl sm:text-3xl text-(--text-main) tracking-tight">Sign in</h1>
			<p class="text-sm text-(--text-secondary) mt-1">Authenticate to your personal workspace or operator console.</p>
		</header>

		{#if error}
			<div class="p-3.5 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-sm">
				{error}
			</div>
		{/if}

		<form on:submit|preventDefault={handleSubmit} class="flex flex-col gap-4">
			<div class="flex flex-col gap-1.5">
				<label for="login-id" class="text-sm font-medium text-(--text-main)">Email or Username</label>
				<input
					id="login-id"
					type="text"
					class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
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
					required
					autocomplete="current-password"
					bind:value={password}
				/>
			</div>

			<button type="submit" class="btn btn-primary w-full py-2.5 mt-2 text-sm font-medium" disabled={loading}>
				{loading ? 'Authenticating...' : 'Sign in →'}
			</button>
		</form>

		<footer class="flex items-center justify-between text-sm text-(--text-muted) pt-3 border-t border-(--border-hairline)">
			<span>Don't have an account?</span>
			<a href="/register" class="text-(--accent-strong) hover:underline font-medium">Create profile</a>
		</footer>
	</div>
</div>
