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
			subdomain: '',
			domain: '',
			name: '',
			username: '',
			email: '',
			password: '',
			confirmPassword: ''
		},
		onSubmit: async (values) => {
			const sub = values.subdomain?.trim().toLowerCase()
				.replace(/^https?:\/\//, '')
				.replace(/:\d+$/, '')
				.replace(/^\.+/, '')
				.replace(/\.+$/, '');
			const baseDom = values.domain.trim().toLowerCase()
				.replace(/^https?:\/\//, '')
				.replace(/:\d+$/, '')
				.replace(/^\.+/, '')
				.replace(/\/.*$/, '');

			const res = await setupApi.setup({
				name: values.name.trim(),
				username: values.username.trim().toLowerCase(),
				email: values.email.trim().toLowerCase(),
				password: values.password,
				domain: baseDom,
				subdomain: sub || undefined
			});

			if (res?.user) {
				user.set(res.user);
				if (baseDom) {
					nodeDomain.set(baseDom);
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
		<!-- Domain & Optional Subdomain -->
		<div class="flex flex-col gap-1.5">
			<span class="text-xs sm:text-sm font-medium text-(--text-main) select-none">
				Host Domain
			</span>
			<div class="flex items-start gap-2">
				<div class="w-2/5 min-w-[110px]">
					<Input
						id="setup-subdomain"
						name="subdomain"
						type="text"
						placeholder="subdomain (opt)"
						inputClass="h-10 px-3 bg-(--bg-muted) font-mono text-xs"
						error={form.errors.subdomain}
						bind:value={form.values.subdomain}
					/>
				</div>
				<div class="flex-1">
					<Input
						id="setup-domain"
						name="domain"
						type="text"
						placeholder="example.com"
						inputClass="h-10 px-3.5 bg-(--bg-muted) font-mono"
						error={form.errors.domain}
						bind:value={form.values.domain}
					/>
				</div>
			</div>
			<p class="text-xs text-(--text-muted) leading-relaxed">
				{#if form.values.domain}
					Canonical address: <span class="font-mono text-(--accent-sky)">https://{form.values.subdomain ? `${form.values.subdomain.trim()}.${form.values.domain.trim()}` : form.values.domain.trim()}</span> (used for invite links and discovery)
				{:else}
					Subdomain is optional (e.g. <span class="font-mono">host</span>.<span class="font-mono">example.com</span>).
				{/if}
			</p>
		</div>

		<div class="h-px bg-(--border-hairline) my-0.5"></div>

		<!-- Vertical Stack: Name, Username, Email, Password, Confirm Password -->
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
			label="Confirm Password"
			type="password"
			placeholder="••••••••"
			autocomplete="new-password"
			inputClass="h-10 px-3.5 bg-(--bg-muted)"
			error={form.errors.confirmPassword}
			bind:value={form.values.confirmPassword}
		/>

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