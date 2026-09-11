<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';
	import { Lock, Ticket, CheckCircle, WarningCircle } from 'phosphor-svelte';

	let username = '';
	let displayName = '';
	let email = '';
	let password = '';
	let invitationToken = '';
	let registrationMode: 'OPEN' | 'INVITE_ONLY' | 'CLOSED' | 'LOADING' = 'LOADING';

	let loading = false;
	let error: string | null = null;
	let validatingToken = false;
	let tokenStatus: { valid: boolean; message?: string; invitedEmail?: string } | null = null;

	onMount(async () => {
		const queryToken = $page.url.searchParams.get('token');
		if (queryToken) {
			invitationToken = queryToken.trim();
		}

		try {
			const res = await api.get('/auth/mode');
			registrationMode = res.data?.registration_mode || 'INVITE_ONLY';
		} catch (err) {
			registrationMode = 'INVITE_ONLY';
		}

		if (invitationToken) {
			await validateToken(invitationToken);
		}
	});

	async function validateToken(token: string) {
		const clean = token.trim();
		if (!clean) {
			tokenStatus = null;
			return;
		}
		validatingToken = true;
		try {
			const res = await api.get(`/invitations/validate?token=${encodeURIComponent(clean)}`);
			if (res.data?.valid) {
				tokenStatus = {
					valid: true,
					invitedEmail: res.data.invited_email || undefined
				};
				if (res.data.invited_email && !email) {
					email = res.data.invited_email;
				}
			} else {
				tokenStatus = {
					valid: false,
					message: res.data?.message || 'The invitation token is invalid or has expired.'
				};
			}
		} catch (err) {
			tokenStatus = {
				valid: false,
				message: extractError(err)
			};
		} finally {
			validatingToken = false;
		}
	}

	async function handleRegister() {
		error = null;
		if (password.length < 8) {
			error = 'Password must be at least 8 characters long.';
			return;
		}

		if (registrationMode === 'INVITE_ONLY' && !invitationToken.trim()) {
			error = 'An invitation token is required to register on this host.';
			return;
		}

		loading = true;
		try {
			const res = await api.post('/auth/register', {
				username,
				display_name: displayName || username,
				email,
				password,
				invitation_token: invitationToken.trim() || undefined
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
	<div class="w-full max-w-md bg-(--bg-surface) border border-(--border-hairline) rounded-md p-7 sm:p-8 flex flex-col gap-6">
		{#if registrationMode === 'LOADING'}
			<div class="py-12 flex flex-col items-center justify-center gap-3 text-(--text-muted)">
				<div class="w-6 h-6 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
				<p class="text-sm">Checking server registration status...</p>
			</div>
		{:else if registrationMode === 'CLOSED'}
			<header class="text-left flex flex-col gap-2.5">
				<div class="w-10 h-10 rounded-md bg-(--bg-muted) border border-(--border-hairline) flex items-center justify-center text-(--text-secondary)">
					<Lock size={20} weight="regular" />
				</div>
				<h1 class="font-display font-bold text-2xl text-(--text-main) tracking-tight">Registrations Closed</h1>
				<p class="text-sm text-(--text-secondary) leading-relaxed">
					This server operates privately for a restricted community. New member registration is currently disabled by the host operator.
				</p>
			</header>

			<div class="p-4 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-secondary) flex flex-col gap-1">
				<span class="font-medium text-(--text-main)">Already have an account?</span>
				<span>Please sign in directly using your account credentials.</span>
			</div>

			<a href="/login" class="btn btn-primary w-full py-2.5 text-sm font-medium text-center">
				Sign In →
			</a>
		{:else}
			<header class="text-left">
				<h1 class="font-display font-bold text-2xl text-(--text-main) tracking-tight">Create Account</h1>
				<p class="text-sm text-(--text-secondary) mt-1">
					{#if registrationMode === 'INVITE_ONLY'}
						Join the community server with a verified invitation token.
					{:else}
						Join the self-hosted community server to showcase your software.
					{/if}
				</p>
			</header>

			{#if error}
				<div class="p-3.5 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-sm">
					{error}
				</div>
			{/if}

			<form on:submit|preventDefault={handleRegister} class="flex flex-col gap-4">
				{#if registrationMode === 'INVITE_ONLY' || invitationToken}
					<div class="flex flex-col gap-1.5">
						<label for="reg-token" class="text-sm font-medium text-(--text-main) flex items-center justify-between">
							<span class="flex items-center gap-1.5">
								<Ticket size={16} />
								Invitation Token {registrationMode === 'INVITE_ONLY' ? '*' : '(Optional)'}
							</span>
							{#if validatingToken}
								<span class="text-xs text-(--text-muted)">Verifying...</span>
							{:else if tokenStatus?.valid}
								<span class="text-xs text-emerald-500 font-medium flex items-center gap-1">
									<CheckCircle size={14} weight="bold" /> Verified
								</span>
							{:else if tokenStatus && !tokenStatus.valid}
								<span class="text-xs text-(--color-danger) font-medium flex items-center gap-1">
									<WarningCircle size={14} weight="bold" /> Invalid
								</span>
							{/if}
						</label>
						<input
							id="reg-token"
							type="text"
							class="w-full px-3.5 py-2.5 bg-(--bg-muted) border {tokenStatus && !tokenStatus.valid ? 'border-(--color-danger)' : 'border-(--border-hairline)'} rounded-md text-sm font-mono text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
							required={registrationMode === 'INVITE_ONLY'}
							placeholder="Enter 64-character token..."
							bind:value={invitationToken}
							on:blur={() => validateToken(invitationToken)}
						/>
						{#if tokenStatus && !tokenStatus.valid && tokenStatus.message}
							<p class="text-xs text-(--color-danger)">{tokenStatus.message}</p>
						{/if}
						{#if tokenStatus?.valid && tokenStatus.invitedEmail}
							<p class="text-xs text-emerald-600 dark:text-emerald-400">
								This invitation is restricted to email <strong>{tokenStatus.invitedEmail}</strong>
							</p>
						{/if}
					</div>
				{/if}

				<div class="flex flex-col gap-1.5">
					<label for="reg-username" class="text-sm font-medium text-(--text-main)">Username *</label>
					<input
						id="reg-username"
						type="text"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
						required
						placeholder="e.g. arya"
						pattern="[a-zA-Z0-9_-]+"
						title="Letters, numbers, underscores, and dashes only"
						bind:value={username}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="reg-display" class="text-sm font-medium text-(--text-main)">Display Name</label>
					<input
						id="reg-display"
						type="text"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
						placeholder="e.g. Arya Pratama"
						bind:value={displayName}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="reg-email" class="text-sm font-medium text-(--text-main)">Email Address *</label>
					<input
						id="reg-email"
						type="email"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
						required
						autocomplete="email"
						bind:value={email}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="reg-password" class="text-sm font-medium text-(--text-main)">Password (min. 8 characters) *</label>
					<input
						id="reg-password"
						type="password"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
						required
						minlength="8"
						autocomplete="new-password"
						bind:value={password}
					/>
				</div>

				<button type="submit" class="btn btn-primary w-full py-2.5 mt-2 text-sm font-medium" disabled={loading || (tokenStatus !== null && !tokenStatus.valid)}>
					{loading ? 'Creating account...' : 'Create Account →'}
				</button>
			</form>

			<footer class="flex items-center justify-between text-sm text-(--text-muted) pt-3 border-t border-(--border-hairline)">
				<span>Already have an account?</span>
				<a href="/login" class="text-(--accent-strong) hover:underline font-medium">Sign in</a>
			</footer>
		{/if}
	</div>
</div>
