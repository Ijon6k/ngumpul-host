<script lang="ts">
	import { goto } from '$app/navigation';
	import { createForm } from '$lib/utils/form.svelte';
	import { setupSchema } from '$lib/schemas/auth';
	import { setupApi } from '$lib/api/setup';
	import { user } from '$lib/stores/auth';
	import { nodeDomain } from '$lib/stores/node';
	import { AuthSplitLayout } from '$lib/components/layout';
	import { Alert, Input, Button } from '$lib/components/ui';

	const form = createForm({
		schema: setupSchema,
		initialValues: {
			domain: '',
			name: '',
			username: '',
			email: '',
			password: '',
			confirmPassword: ''
		},
		onSubmit: async (values) => {
			const res = await setupApi.setup({
				name: values.name.trim(),
				username: values.username.trim().toLowerCase(),
				email: values.email.trim().toLowerCase(),
				password: values.password,
				domain: values.domain.trim()
			});

			if (res?.user) {
				user.set(res.user);
				if (values.domain) {
					const clean = values.domain
						.replace(/^https?:\/\//, '')
						.replace(/:\d+$/, '')
						.trim();
					if (clean) nodeDomain.set(clean);
				}
				await goto('/admin');
			}
		}
	});
</script>

<svelte:head>
	<title>Node Setup · Ngumpul Host</title>
</svelte:head>

<AuthSplitLayout
	title="Host Node Setup"
	description="Configure the canonical node domain and initialize the primary administrator account."
>
	{#if form.serverError}
		<Alert variant="danger">{form.serverError}</Alert>
	{/if}

	<form onsubmit={form.handleSubmit} class="flex flex-col gap-4">
		<Input
			id="setup-domain"
			name="domain"
			label="Node Domain"
			type="text"
			placeholder="ngumpul.example.com"
			helperText="Canonical host address used for system telemetry and hosted project subdomains."
			inputClass="h-10 px-3.5 bg-(--bg-muted) font-mono"
			error={form.errors.domain}
			bind:value={form.values.domain}
		>
			{#snippet labelExtra()}
				<span class="text-xs text-(--text-muted) font-mono">e.g. ngumpul.id</span>
			{/snippet}
		</Input>

		<div class="h-px bg-(--border-hairline) my-0.5"></div>

		<Input
			id="setup-name"
			name="name"
			label="Operator Name"
			type="text"
			placeholder="Site Operator"
			inputClass="h-10 px-3.5 bg-(--bg-muted)"
			error={form.errors.name}
			bind:value={form.values.name}
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
				error={form.errors.username}
				bind:value={form.values.username}
			/>

			<Input
				id="setup-email"
				name="email"
				label="Email"
				type="email"
				placeholder="admin@example.com"
				autocomplete="email"
				inputClass="h-10 px-3.5 bg-(--bg-muted)"
				error={form.errors.email}
				bind:value={form.values.email}
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
				error={form.errors.password}
				bind:value={form.values.password}
			/>

			<Input
				id="setup-confirm-password"
				name="confirmPassword"
				label="Confirm"
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
			{form.isSubmitting ? 'Configuring node...' : 'Complete Node Setup'}
		</Button>
	</form>

	<svelte:fragment slot="footer">
		<footer class="pt-2 text-center text-xs text-(--text-muted)">
			Ngumpul Host · Initial Node Provisioning
		</footer>
	</svelte:fragment>
</AuthSplitLayout>