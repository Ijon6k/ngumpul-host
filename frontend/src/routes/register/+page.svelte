<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { superForm } from 'sveltekit-superforms';
	import { valibotClient } from 'sveltekit-superforms/adapters';
	import { registerSchema } from '$lib/schemas/auth';
	import { authApi, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';
	import { AuthSplitLayout } from '$lib/components/layout';
	import { Alert, Input } from '$lib/components/ui';
	import { CheckCircle, WarningCircle, Lock } from 'phosphor-svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let registrationMode: 'OPEN' | 'INVITE_ONLY' | 'CLOSED' | 'LOADING' = $state('LOADING');
	let validatingToken = $state(false);
	let tokenStatus: { valid: boolean; message?: string; invitedEmail?: string } | null = $state(null);

	const formState = $derived(data.form);

	// svelte-ignore state_referenced_locally — superform reads the initial form payload at setup
	const { form, errors, constraints, message, submitting, enhance } = superForm(
		formState,
		{ validators: valibotClient(registerSchema), resetForm: false }
	);

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
				if (res.invited_email && !$form.email) {
					form.update(($f) => ({ ...$f, email: res.invited_email! }), { taint: false });
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
			form.update(($f) => ({ ...$f, invitationToken: queryToken.trim() }), { taint: false });
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

	function handleResult({ result }: { result: { type: string; data?: any } }) {
		if (result.type === 'success' && result.data?.user) {
			user.set(result.data.user);
			goto('/me');
		}
	}

	function fieldError(errors: Record<string, any> | undefined, field: string): string {
		const val = errors?.[field];
		return Array.isArray(val) && val.length > 0 ? val[0] : '';
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
		{#if $message}
			<Alert variant="danger">{$message}</Alert>
		{/if}

		{#if $errors._errors}
			<Alert variant="danger">{$errors._errors[0]}</Alert>
		{/if}

		<form method="POST" use:enhance class="flex flex-col gap-4">
			{#if registrationMode === 'INVITE_ONLY'}
				<Input
					id="reg-token"
					label="Invitation Token"
					type="text"
					placeholder="Paste invite token"
					inputClass="h-10 px-3.5 bg-(--bg-muted) font-mono"
					error={fieldError($errors, 'invitationToken') ||
						(tokenStatus && !tokenStatus.valid ? tokenStatus.message : '')}
					bind:value={$form.invitationToken}
					onblur={() => validateToken($form.invitationToken)}
					{...$constraints.invitationToken}
				>
					<svelte:fragment slot="label-extra">
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
					</svelte:fragment>
				</Input>
			{/if}

			<Input
				id="reg-name"
				label="Display Name"
				type="text"
				placeholder="Alex Rivera"
				inputClass="h-10 px-3.5 bg-(--bg-muted)"
				error={fieldError($errors, 'displayName')}
				bind:value={$form.displayName}
				{...$constraints.displayName}
			/>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
				<Input
					id="reg-username"
					label="Username"
					type="text"
					placeholder="alex"
					autocomplete="username"
					inputClass="h-10 px-3.5 bg-(--bg-muted)"
					error={fieldError($errors, 'username')}
					bind:value={$form.username}
					{...$constraints.username}
				/>

				<Input
					id="reg-email"
					label="Email"
					type="email"
					placeholder="alex@example.com"
					autocomplete="email"
					inputClass="h-10 px-3.5 bg-(--bg-muted)"
					error={fieldError($errors, 'email')}
					bind:value={$form.email}
					{...$constraints.email}
				/>
			</div>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
				<Input
					id="reg-pass"
					label="Password"
					type="password"
					placeholder="•••••••• (min 8 characters)"
					autocomplete="new-password"
					inputClass="h-10 px-3.5 bg-(--bg-muted)"
					error={fieldError($errors, 'password')}
					bind:value={$form.password}
					{...$constraints.password}
				/>

				<Input
					id="reg-confirm-pass"
					label="Confirm Password"
					type="password"
					placeholder="••••••••"
					autocomplete="new-password"
					inputClass="h-10 px-3.5 bg-(--bg-muted)"
					error={fieldError($errors, 'confirmPassword')}
					bind:value={$form.confirmPassword}
					{...$constraints.confirmPassword}
				/>
			</div>

			<button
				type="submit"
				class="btn btn-primary w-full py-2.5 mt-2 text-sm font-medium"
				disabled={$submitting}
			>
				{$submitting ? 'Creating account...' : 'Create account'}
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