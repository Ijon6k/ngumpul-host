<script lang="ts">
	import { setupApi, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';
	import { AuthSplitLayout } from '$lib/components/layout';
	import { goto } from '$app/navigation';

	let domain = '';
	let name = '';
	let username = '';
	let email = '';
	let password = '';
	let confirmPassword = '';

	let loading = false;
	let error: string | null = null;

	async function handleSetup() {
		error = null;

		if (password.length < 8) {
			error = 'Password must be at least 8 characters long.';
			return;
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match.';
			return;
		}

		loading = true;
		try {
			const res = await setupApi.setup({
				name: name.trim(),
				username: username.trim().toLowerCase(),
				email: email.trim().toLowerCase(),
				password,
				domain: domain.trim()
			});

			if (res?.user) {
				user.set(res.user);
				await goto('/admin');
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Node Setup · Ngumpul Host</title>
</svelte:head>

<AuthSplitLayout
	title="Host Node Setup"
	description="Configure the canonical node domain and initialize the primary administrator account."
>
	{#if error}
		<div class="p-3.5 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-sm leading-relaxed">
			{error}
		</div>
	{/if}

	<form on:submit|preventDefault={handleSetup} class="flex flex-col gap-4">
		<!-- Domain Section -->
		<div class="flex flex-col gap-1.5">
			<label for="setup-domain" class="text-sm font-medium text-(--text-main) flex items-center justify-between">
				<span>Node Domain</span>
				<span class="text-xs text-(--text-muted) font-mono">e.g. ngumpul.id</span>
			</label>
			<input
				id="setup-domain"
				type="text"
				class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50 font-mono text-xs sm:text-sm"
				placeholder="ngumpul.example.com"
				required
				bind:value={domain}
			/>
			<p class="text-xs text-(--text-muted)">Canonical host address used for system telemetry and hosted project subdomains.</p>
		</div>

		<div class="h-px bg-(--border-hairline) my-0.5"></div>

		<!-- Operator Section -->
		<div class="flex flex-col gap-1.5">
			<label for="setup-name" class="text-sm font-medium text-(--text-main)">Operator Name</label>
			<input
				id="setup-name"
				type="text"
				class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50"
				placeholder="Site Operator"
				required
				bind:value={name}
			/>
		</div>

		<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
			<div class="flex flex-col gap-1.5">
				<label for="setup-username" class="text-sm font-medium text-(--text-main)">Username</label>
				<input
					id="setup-username"
					type="text"
					class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50"
					placeholder="admin"
					required
					autocomplete="username"
					bind:value={username}
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="setup-email" class="text-sm font-medium text-(--text-main)">Email</label>
				<input
					id="setup-email"
					type="email"
					class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50"
					placeholder="admin@example.com"
					required
					autocomplete="email"
					bind:value={email}
				/>
			</div>
		</div>

		<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
			<div class="flex flex-col gap-1.5">
				<label for="setup-password" class="text-sm font-medium text-(--text-main)">Password</label>
				<input
					id="setup-password"
					type="password"
					class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
					placeholder="••••••••"
					required
					autocomplete="new-password"
					bind:value={password}
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="setup-confirm-password" class="text-sm font-medium text-(--text-main)">Confirm</label>
				<input
					id="setup-confirm-password"
					type="password"
					class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
					placeholder="••••••••"
					required
					autocomplete="new-password"
					bind:value={confirmPassword}
				/>
			</div>
		</div>

		<button
			type="submit"
			class="btn btn-primary w-full py-2.5 mt-2 text-sm font-medium"
			disabled={loading}
		>
			{loading ? 'Configuring node...' : 'Complete Node Setup'}
		</button>
	</form>

	<svelte:fragment slot="footer">
		<footer class="pt-2 text-center text-xs text-(--text-muted)">
			Ngumpul Host · Initial Node Provisioning
		</footer>
	</svelte:fragment>
</AuthSplitLayout>
