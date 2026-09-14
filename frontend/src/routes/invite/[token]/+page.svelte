<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { authApi, extractError } from '$lib/api';
	import { Ticket, CheckCircle, WarningCircle, ArrowRight } from 'phosphor-svelte';

	let token = '';
	let loading = true;
	let valid = false;
	let message = '';
	let invitedEmail: string | null = null;
	let expiresAt: string | null = null;

	onMount(async () => {
		token = $page.params.token || '';
		if (!token) {
			loading = false;
			valid = false;
			message = 'No invitation token was provided.';
			return;
		}

		try {
			const res = await authApi.validateInvitation(token);
			if (res?.valid) {
				valid = true;
				invitedEmail = res.invited_email || null;
				expiresAt = res.expires_at || null;
			} else {
				valid = false;
				message = res?.message || 'The invitation token is invalid or has expired.';
			}
		} catch (err) {
			valid = false;
			message = extractError(err);
		} finally {
			loading = false;
		}
	});

	function formatDate(iso: string): string {
		try {
			return new Date(iso).toLocaleDateString('en-US', {
				day: 'numeric',
				month: 'short',
				year: 'numeric'
			});
		} catch {
			return iso;
		}
	}
</script>

<svelte:head>
	<title>Invitation Accepted — Ngumpul Host</title>
</svelte:head>

<div class="container mx-auto px-6 py-20 flex justify-center items-center">
	<div class="w-full max-w-md bg-(--bg-surface) border border-(--border-hairline) rounded-md p-7 sm:p-8 flex flex-col gap-6">
		{#if loading}
			<div class="py-12 flex flex-col items-center justify-center gap-3 text-(--text-muted)">
				<div class="w-6 h-6 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
				<p class="text-sm">Verifying invitation link...</p>
			</div>
		{:else if valid}
			<header class="text-left flex flex-col gap-2.5">
				<div class="w-11 h-11 rounded-md bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-500">
					<Ticket size={22} weight="regular" />
				</div>
				<h1 class="font-display font-bold text-2xl sm:text-3xl text-(--text-main) tracking-tight">Invitation Accepted</h1>
				<p class="text-sm text-(--text-secondary) leading-relaxed">
					You have received an invitation to join the self-hosted community server <strong>Ngumpul Host</strong>.
				</p>
			</header>

			<div class="p-4 bg-(--bg-muted) border border-(--border-hairline) rounded-md flex flex-col gap-3 text-sm">
				<div class="flex items-center justify-between text-(--text-muted)">
					<span>Invitation Status</span>
					<span class="text-emerald-500 font-medium flex items-center gap-1.5">
						<CheckCircle size={15} weight="bold" /> Ready to Use
					</span>
				</div>

				{#if invitedEmail}
					<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1.5 pt-2.5 border-t border-(--border-hairline)/60">
						<span class="text-(--text-muted)">Restricted to Email</span>
						<span class="font-medium text-(--text-main) min-w-0 break-all">{invitedEmail}</span>
					</div>
				{/if}

				{#if expiresAt}
					<div class="flex items-center justify-between pt-2.5 border-t border-(--border-hairline)/60">
						<span class="text-(--text-muted)">Expires At</span>
						<span class="text-(--text-main)">{formatDate(expiresAt)}</span>
					</div>
				{/if}
			</div>

			<div class="flex flex-col gap-3">
				<a
					href={`/register?token=${encodeURIComponent(token)}`}
					class="btn btn-primary w-full py-2.5 text-sm font-medium text-center flex items-center justify-center gap-2"
				>
					<span>Continue to Registration</span>
					<ArrowRight size={15} weight="bold" />
				</a>
				<a href="/projects" class="text-center text-sm text-(--text-muted) hover:text-(--text-main) transition-colors py-1">
					Explore Community Projects
				</a>
			</div>
		{:else}
			<header class="text-left flex flex-col gap-2.5">
				<div class="w-11 h-11 rounded-md bg-red-500/10 border border-(--color-danger)/20 flex items-center justify-center text-(--color-danger)">
					<WarningCircle size={22} weight="regular" />
				</div>
				<h1 class="font-display font-bold text-2xl sm:text-3xl text-(--text-main) tracking-tight">
					{message.toLowerCase().includes('closed') ? 'Registrations Closed' : 'Invalid Invitation'}
				</h1>
				<p class="text-sm text-(--text-secondary) leading-relaxed">
					{message || 'The invitation link is invalid, has been revoked, or has expired.'}
				</p>
			</header>

			<div class="pt-2 flex flex-col gap-2.5">
				<a href="/" class="btn btn-secondary w-full py-2.5 text-sm font-medium text-center">
					Back to Homepage
				</a>
			</div>
		{/if}
	</div>
</div>
