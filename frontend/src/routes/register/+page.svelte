<script lang="ts">
	import { api, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';

	let username = '';
	let displayName = '';
	let email = '';
	let password = '';
	let loading = false;
	let error: string | null = null;

	async function handleRegister() {
		error = null;
		if (password.length < 8) {
			error = 'Password must be at least 8 characters long.';
			return;
		}

		loading = true;
		try {
			const res = await api.post('/auth/register', {
				username,
				display_name: displayName || username,
				email,
				password
			});
			if (res.data?.user) {
				user.set(res.data.user);
				window.location.href = '/me';
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto px-6 py-16 flex justify-center items-center">
	<div class="w-full max-w-sm bg-(--bg-surface) border border-(--border-hairline) rounded-md p-7 sm:p-8 flex flex-col gap-6">
		<header class="text-left">
			<h1 class="font-display font-bold text-2xl text-(--text-main) tracking-tight">Create Profile</h1>
			<p class="text-xs text-(--text-secondary) mt-1">Join the community node to showcase projects and request hosting.</p>
		</header>

		{#if error}
			<div class="p-3 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-xs">
				{error}
			</div>
		{/if}

		<form on:submit|preventDefault={handleRegister} class="flex flex-col gap-3.5">
			<div class="flex flex-col gap-1.5">
				<label for="reg-username" class="text-xs font-medium text-(--text-main)">Username *</label>
				<input
					id="reg-username"
					type="text"
					class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
					required
					placeholder="e.g. dev_arya"
					pattern="[a-zA-Z0-9_-]+"
					title="Letters, numbers, underscores, or hyphens only"
					bind:value={username}
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="reg-display" class="text-xs font-medium text-(--text-main)">Display Name</label>
				<input
					id="reg-display"
					type="text"
					class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
					placeholder="e.g. Arya Pratama"
					bind:value={displayName}
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="reg-email" class="text-xs font-medium text-(--text-main)">Email Address *</label>
				<input
					id="reg-email"
					type="email"
					class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
					required
					autocomplete="email"
					bind:value={email}
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="reg-password" class="text-xs font-medium text-(--text-main)">Password (min 8 chars) *</label>
				<input
					id="reg-password"
					type="password"
					class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
					required
					minlength="8"
					autocomplete="new-password"
					bind:value={password}
				/>
			</div>

			<button type="submit" class="btn btn-primary w-full py-2.5 mt-2 text-xs" disabled={loading}>
				{loading ? 'Registering profile...' : 'Register Profile →'}
			</button>
		</form>

		<footer class="flex items-center justify-between text-xs text-(--text-muted) pt-3 border-t border-(--border-hairline)">
			<span>Already have a profile?</span>
			<a href="/login" class="text-(--accent-strong) hover:underline font-medium">Sign in</a>
		</footer>
	</div>
</div>
