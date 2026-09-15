<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { createForm } from '$lib/utils/form.svelte';
	import { registerSchema } from '$lib/schemas/auth';
	import { authApi, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';
	import { AuthSplitLayout } from '$lib/components/layout';
	import { Alert, Input, Button } from '$lib/components/ui';
	import { CheckCircle, WarningCircle, Lock } from 'phosphor-svelte';

	let registrationMode: 'OPEN' | 'INVITE_ONLY' | 'CLOSED' | 'LOADING' = $state('LOADING');
	let validatingToken = $state(false);
	let tokenStatus: { valid: boolean; message?: string; invitedEmail?: string } | null = $state(null);

	const form = createForm({
		schema: registerSchema,
		initialValues: {
			invitationToken: '',
			displayName: '',
			username: '',
			email: '',
			password: '',
			confirmPassword: ''
		},
		onSubmit: async (values) => {
			const res = await authApi.register({
				username: values.username.trim().toLowerCase(),
				email: values.email.trim().toLowerCase(),
				password: values.password,
				display_name: values.displayName.trim(),
				invitation_token: values.invitationToken ? values.invitationToken.trim() : undefined
			});
			if (res?.user) {
				user.set(res.user);
				await goto('/me');
			}
		}
	});

	async function validateToken(raw: string) {
		const clean = raw.trim();
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
				if (res.invited_email && !form.values.email) {
					form.values.email = res.invited_email;
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

	onMount(async () => {
		const queryToken = $page.url.searchParams.get('token');
		if (queryToken) {
			form.values.invitationToken = queryToken.trim();
		}

		try {
			const res = await authApi.getRegistrationMode();
			registrationMode = (res?.registration_mode as any) || 'INVITE_ONLY';
		} catch {
			registrationMode = 'INVITE_ONLY';
		}

		if (queryToken) {
			await validateToken(queryToken);
		}
	});
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
		{#if form.serverError}
			<Alert variant="danger">{form.serverError}</Alert>
		{/if}

		<form onsubmit={form.handleSubmit} class="flex flex-col gap-4">
			{#if registrationMode === 'INVITE_ONLY'}
				<Input
					id="reg-token"
					name="invitationToken"
					label="Invitation Token"
					type="text"
					placeholder="Paste invite token"
					inputClass="h-10 px-3.5 bg-(--bg-muted) font-mono"
					error={form.errors.invitationToken ||
						(tokenStatus && !tokenStatus.valid ? tokenStatus.message : '')}
					bind:value={form.values.invitationToken}
					onblur={() => validateToken(form.values.invitationToken)}
				>
					{#snippet labelExtra()}
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
					{/snippet}
				</Input>
			{/if}

			<Input
				id="reg-name"
				name="displayName"
				label="Display Name"
				type="text"
				placeholder="Alex Rivera"
				inputClass="h-10 px-3.5 bg-(--bg-muted)"
				error={form.errors.displayName}
				bind:value={form.values.displayName}
			/>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
				<Input
					id="reg-username"
					name="username"
					label="Username"
					type="text"
					placeholder="alex"
					autocomplete="username"
					inputClass="h-10 px-3.5 bg-(--bg-muted)"
					error={form.errors.username}
					bind:value={form.values.username}
				/>

				<Input
					id="reg-email"
					name="email"
					label="Email"
					type="email"
					placeholder="alex@example.com"
					autocomplete="email"
					inputClass="h-10 px-3.5 bg-(--bg-muted)"
					error={form.errors.email}
					bind:value={form.values.email}
				/>
			</div>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
				<Input
					id="reg-pass"
					name="password"
					label="Password"
					type="password"
					placeholder="•••••••• (min 8 characters)"
					autocomplete="new-password"
					inputClass="h-10 px-3.5 bg-(--bg-muted)"
					error={form.errors.password}
					bind:value={form.values.password}
				/>

				<Input
					id="reg-confirm-pass"
					name="confirmPassword"
					label="Confirm Password"
					type="password"
					placeholder="••••••••"
					autocomplete="new-password"
					inputClass="h-10 px-3.5 bg-(--bg-muted)"
					error={form.errors.confirmPassword}
					bind:value={form.values.confirmPassword}
				/>
			</div>

			<Button
				type="submit"
				variant="primary"
				size="md"
				class="w-full mt-2 font-medium"
				loading={form.isSubmitting}
				disabled={form.isSubmitting}
			>
				{form.isSubmitting ? 'Creating account...' : 'Create account'}
			</Button>
		</form>
	{/if}

	<svelte:fragment slot="footer">
		<footer class="flex items-center justify-between text-sm text-(--text-muted) pt-4 border-t border-(--border-hairline)">
			<span>Already registered?</span>
			<a href="/login" class="text-(--accent-strong) hover:underline font-medium">Sign in</a>
		</footer>
	</svelte:fragment>
</AuthSplitLayout>