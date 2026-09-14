<script lang="ts">
	import { untrack } from 'svelte';
	import { superForm } from 'sveltekit-superforms';
	import { valibotClient } from 'sveltekit-superforms/adapters';
	import { setupSchema } from '$lib/schemas/auth';
	import { user } from '$lib/stores/auth';
	import { AuthSplitLayout } from '$lib/components/layout';
	import { Alert, Input } from '$lib/components/ui';
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const { form, errors, constraints, message, submitting, enhance } = superForm(
		untrack(() => data.form),
		{
			validators: valibotClient(setupSchema),
			resetForm: false,
			onResult: handleResult
		}
	);

	function handleResult({ result }: { result: { type: string; data?: any } }) {
		if (result.type === 'success' && result.data?.user) {
			user.set(result.data.user);
			goto('/admin');
		}
	}

	function fieldError(errors: Record<string, any> | undefined, field: string): string {
		const val = errors?.[field];
		return Array.isArray(val) && val.length > 0 ? val[0] : '';
	}
</script>

<svelte:head>
	<title>Node Setup · Ngumpul Host</title>
</svelte:head>

<AuthSplitLayout
	title="Host Node Setup"
	description="Configure the canonical node domain and initialize the primary administrator account."
>
	{#if $message}
		<Alert variant="danger">{$message}</Alert>
	{/if}

	{#if $errors._errors}
		<Alert variant="danger">{$errors._errors[0]}</Alert>
	{/if}

	<form method="POST" use:enhance class="flex flex-col gap-4">
		<Input
			id="setup-domain"
			name="domain"
			label="Node Domain"
			type="text"
			placeholder="ngumpul.example.com"
			helperText="Canonical host address used for system telemetry and hosted project subdomains."
			inputClass="h-10 px-3.5 bg-(--bg-muted) font-mono"
			error={fieldError($errors, 'domain')}
			bind:value={$form.domain}
			{...$constraints.domain}
		>
			<svelte:fragment slot="label-extra">
				<span class="text-xs text-(--text-muted) font-mono">e.g. ngumpul.id</span>
			</svelte:fragment>
		</Input>

		<div class="h-px bg-(--border-hairline) my-0.5"></div>

		<Input
			id="setup-name"
			name="name"
			label="Operator Name"
			type="text"
			placeholder="Site Operator"
			inputClass="h-10 px-3.5 bg-(--bg-muted)"
			error={fieldError($errors, 'name')}
			bind:value={$form.name}
			{...$constraints.name}
		/>

		<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
			<Input
				id="setup-username"
				name="username"
				label="Username"
				type="text"
				placeholder="admin"
				autocomplete="username"
				inputClass="h-10 px-3.5 bg-(--bg-muted)"
				error={fieldError($errors, 'username')}
				bind:value={$form.username}
				{...$constraints.username}
			/>

			<Input
				id="setup-email"
				name="email"
				label="Email"
				type="email"
				placeholder="admin@example.com"
				autocomplete="email"
				inputClass="h-10 px-3.5 bg-(--bg-muted)"
				error={fieldError($errors, 'email')}
				bind:value={$form.email}
				{...$constraints.email}
			/>
		</div>

		<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
			<Input
				id="setup-password"
				name="password"
				label="Password"
				type="password"
				placeholder="••••••••"
				autocomplete="new-password"
				inputClass="h-10 px-3.5 bg-(--bg-muted)"
				error={fieldError($errors, 'password')}
				bind:value={$form.password}
				{...$constraints.password}
			/>

			<Input
				id="setup-confirm-password"
				name="confirmPassword"
				label="Confirm"
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
			{$submitting ? 'Configuring node...' : 'Complete Node Setup'}
		</button>
	</form>

	<svelte:fragment slot="footer">
		<footer class="pt-2 text-center text-xs text-(--text-muted)">
			Ngumpul Host · Initial Node Provisioning
		</footer>
	</svelte:fragment>
</AuthSplitLayout>