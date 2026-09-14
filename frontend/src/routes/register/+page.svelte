<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { authApi, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';
	import { AuthSplitLayout } from '$lib/components/layout';
	import { CheckCircle, WarningCircle, Lock } from 'phosphor-svelte';

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
			const res = await authApi.getRegistrationMode();
			registrationMode = (res?.registration_mode as any) || 'INVITE_ONLY';
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
			const res = await authApi.validateInvitation(clean);
			if (res?.valid) {
				tokenStatus = {
					valid: true,
					invitedEmail: res.invited_email || undefined
				};
				if (res.invited_email && !email) {
					email = res.invited_email;
				}
			} else {
				tokenStatus = {
					valid: false,
					message: res?.message || 'The invitation token is invalid or has expired.'
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
			error = 'An invitation token is required to register on this node.';
			return;
		}

		loading = true;
		try {
			const res = await authApi.register({
				username: username.trim().toLowerCase(),
				display_name: displayName.trim() || username.trim(),
				email: email.trim().toLowerCase(),
				password,
				invitation_token: invitationToken.trim() || undefined
			});
			if (res?.user) {
				user.set(res.user);
				await goto('/me');
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Create Profile · Ngumpul Host</title>
</svelte:head>

<AuthSplitLayout
	title={registrationMode === 'CLOSED' ? 'Registrations Closed' : 'Create profile'}
	description={registrationMode === 'CLOSED'
		? 'New member registration is currently disabled by the host operator.'
		: registrationMode === 'INVITE_ONLY'
			? 'Join the community server with a verified invitation token.'
			: 'Join the community server to share and showcase independent projects.'}
>
	{#if registrationMode === 'LOADING'}
		<div class="py-12 flex flex-col items-center justify-center gap-3 text-(--text-muted)">
			<div class="w-5 h-5 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
			<p class="text-xs">Checking registration mode...</p>
		</div>
	{:else if registrationMode === 'CLOSED'}
		<div class="p-4 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-secondary) flex flex-col gap-1.5">
			<div class="flex items-center gap-2 text-(--text-main) font-medium">
				<Lock size={16} weight="regular" />
				<span>Private Community</span>
			</div>
			<p class="text-xs text-(--text-muted) leading-relaxed">
				This node operates on a closed basis. Existing members can sign in directly.
			</p>
		</div>

		<a href="/login" class="btn btn-primary w-full py-2.5 text-sm font-medium text-center">
			Sign in
		</a>
	{:else}
		{#if error}
			<div class="p-3.5 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-sm leading-relaxed">
				{error}
			</div>
		{/if}

		<form on:submit|preventDefault={handleRegister} class="flex flex-col gap-4">
			{#if registrationMode === 'INVITE_ONLY'}
				<div class="flex flex-col gap-1.5">
					<label for="reg-token" class="text-sm font-medium text-(--text-main) flex items-center justify-between">
						<span>Invitation Token</span>
						{#if validatingToken}
							<span class="text-xs text-(--text-muted)">Validating...</span>
						{:else if tokenStatus?.valid}
							<span class="text-xs text-emerald-500 font-medium flex items-center gap-1">
								<CheckCircle size={13} weight="bold" /> Verified
							</span>
						{:else if tokenStatus && !tokenStatus.valid}
							<span class="text-xs text-(--color-danger) font-medium flex items-center gap-1">
								<WarningCircle size={13} weight="bold" /> Invalid
							</span>
						{/if}
					</label>
					<input
						id="reg-token"
						type="text"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm font-mono text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50"
						placeholder="Paste invite token"
						required
						bind:value={invitationToken}
						on:blur={() => validateToken(invitationToken)}
					/>
					{#if tokenStatus?.message && !tokenStatus.valid}
						<p class="text-xs text-(--color-danger)">{tokenStatus.message}</p>
					{/if}
				</div>
			{/if}

			<div class="flex flex-col gap-1.5">
				<label for="reg-name" class="text-sm font-medium text-(--text-main)">Display Name</label>
				<input
					id="reg-name"
					type="text"
					class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50"
					placeholder="Alex Rivera"
					required
					bind:value={displayName}
				/>
			</div>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
				<div class="flex flex-col gap-1.5">
					<label for="reg-username" class="text-sm font-medium text-(--text-main)">Username</label>
					<input
						id="reg-username"
						type="text"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50"
						placeholder="alex"
						required
						autocomplete="username"
						bind:value={username}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="reg-email" class="text-sm font-medium text-(--text-main)">Email</label>
					<input
						id="reg-email"
						type="email"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors placeholder:text-(--text-muted)/50"
						placeholder="alex@example.com"
						required
						autocomplete="email"
						bind:value={email}
					/>
				</div>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="reg-pass" class="text-sm font-medium text-(--text-main)">Password</label>
				<input
					id="reg-pass"
					type="password"
					class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main) transition-colors"
					placeholder="•••••••• (min 8 characters)"
					required
					autocomplete="new-password"
					bind:value={password}
				/>
			</div>

			<button type="submit" class="btn btn-primary w-full py-2.5 mt-2 text-sm font-medium" disabled={loading}>
				{loading ? 'Creating account...' : 'Create account'}
			</button>
		</form>
	{/if}

	<svelte:fragment slot="footer">
		<footer class="flex items-center justify-between text-sm text-(--text-muted) pt-4 border-t border-(--border-hairline)">
			<span>Already registered?</span>
			<a href="/login" class="text-(--accent-strong) hover:underline font-medium">Sign in</a>
		</footer>
	</svelte:fragment>
</AuthSplitLayout>
