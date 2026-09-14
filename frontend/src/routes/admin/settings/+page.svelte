<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi, extractError } from '$lib/api';
	import { refreshNodeDomain } from '$lib/stores/node';
	import {
		Gear,
		Ticket,
		ShieldCheck,
		Plus,
		Copy,
		Check,
		CheckCircle,
		Clock,
		Globe,
		WarningCircle
	} from 'phosphor-svelte';
	import { Table, TableRow, TableCell, type TableColumn, ConfirmModal } from '$lib/components/ui';

	const columns: TableColumn[] = [
		{ key: 'creator', label: 'Created By' },
		{ key: 'email', label: 'Email Restriction' },
		{ key: 'link', label: 'Invitation Link' },
		{ key: 'usage', label: 'Usage' },
		{ key: 'expires', label: 'Expires At' },
		{ key: 'status', label: 'Status' },
		{ key: 'actions', label: 'Action', align: 'right' }
	];

	let registrationMode = 'INVITE_ONLY';
	let savedRegistrationMode = 'INVITE_ONLY';
	let serverDomain = '';
	let savedServerDomain = '';
	let loadingSettings = true;
	let savingSettings = false;
	let settingsError: string | null = null;
	let settingsSuccess: string | null = null;
	let confirmCloseModalOpen = false;

	// Revoke invitation modal state
	let revokeTargetId: string | null = null;
	let revokeLoading = false;
	let revokeError = '';

	let invitations: any[] = [];
	let loadingInvitations = true;
	let invitationsError: string | null = null;

	// New invitation form state
	let maxUses = 1;
	let expiresInDays = 7;
	let invitedEmail = '';
	let creatingInvite = false;
	let inviteSuccessToken: { raw_token: string; invite_url: string } | null = null;
	let copied = false;
	let copiedTokenId: string | null = null;

	function copyInviteLink(inv: any) {
		if (!inv.token) return;
		const url = `${window.location.origin}/invite/${inv.token}`;
		if (navigator?.clipboard) {
			navigator.clipboard.writeText(url);
		}
		copiedTokenId = inv.id;
		setTimeout(() => {
			if (copiedTokenId === inv.id) {
				copiedTokenId = null;
			}
		}, 2000);
	}

	async function loadSettings() {
		loadingSettings = true;
		try {
			const res = await adminApi.settings.getSettings();
			registrationMode = res.registration_mode || 'INVITE_ONLY';
			savedRegistrationMode = registrationMode;
			serverDomain = res.domain || '';
			savedServerDomain = serverDomain;
		} catch (err) {
			settingsError = extractError(err);
		} finally {
			loadingSettings = false;
		}
	}

	function handleSaveSettingsRequest() {
		if (registrationMode === 'CLOSED' && savedRegistrationMode !== 'CLOSED') {
			confirmCloseModalOpen = true;
			return;
		}
		executeSaveSettings();
	}

	async function executeSaveSettings() {
		confirmCloseModalOpen = false;
		savingSettings = true;
		settingsError = null;
		settingsSuccess = null;
		try {
			await adminApi.settings.updateSettings({
				registration_mode: registrationMode,
				domain: serverDomain.trim()
			});
			savedRegistrationMode = registrationMode;
			savedServerDomain = serverDomain.trim();
			await refreshNodeDomain();
			settingsSuccess = 'Server configuration updated successfully.';
		} catch (err) {
			settingsError = extractError(err);
		} finally {
			savingSettings = false;
		}
	}

	async function loadInvitations() {
		loadingInvitations = true;
		try {
			const res = await adminApi.settings.listInvitations();
			invitations = res.invitations || [];
		} catch (err) {
			invitationsError = extractError(err);
		} finally {
			loadingInvitations = false;
		}
	}

	async function handleCreateInvitation() {
		creatingInvite = true;
		invitationsError = null;
		inviteSuccessToken = null;
		try {
			const payload: Record<string, any> = {
				max_uses: Number(maxUses) || 1,
				expires_in_days: Number(expiresInDays) || 7
			};
			if (invitedEmail.trim()) {
				payload.invited_email = invitedEmail.trim();
			}

			const res = await adminApi.settings.createInvitation(payload);
			if (res.raw_token) {
				inviteSuccessToken = {
					raw_token: res.raw_token,
					invite_url: res.invite_url || `${window.location.origin}/invite/${res.raw_token}`
				};
				invitedEmail = '';
				maxUses = 1;
				await loadInvitations();
			}
		} catch (err) {
			invitationsError = extractError(err);
		} finally {
			creatingInvite = false;
		}
	}

	function openRevokeModal(id: string) {
		revokeTargetId = id;
		revokeError = '';
	}

	async function confirmRevoke() {
		if (!revokeTargetId) return;
		revokeLoading = true;
		revokeError = '';
		try {
			await adminApi.settings.revokeInvitation(revokeTargetId);
			revokeTargetId = null;
			await loadInvitations();
		} catch (err) {
			revokeError = extractError(err);
		} finally {
			revokeLoading = false;
		}
	}

	function copyToClipboard(text: string) {
		if (!navigator?.clipboard) return;
		navigator.clipboard.writeText(text);
		copied = true;
		setTimeout(() => {
			copied = false;
		}, 2000);
	}

	function formatDate(iso?: string): string {
		if (!iso) return 'Permanent (No expiry)';
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

	function getStatusLabel(inv: any): { text: string; color: string } {
		if (inv.revoked_at) {
			return { text: 'Revoked', color: 'bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/20' };
		}
		if (inv.expires_at && new Date(inv.expires_at).getTime() < Date.now()) {
			return { text: 'Expired', color: 'bg-neutral-500/10 text-neutral-600 dark:text-neutral-400 border-neutral-500/20' };
		}
		if (inv.used_count >= inv.max_uses) {
			return { text: 'Exhausted', color: 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20' };
		}
		return { text: 'Active', color: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20' };
	}

	onMount(() => {
		loadSettings();
		loadInvitations();
	});
</script>

<svelte:head>
	<title>Server Settings & Invitations — Ngumpul Host</title>
</svelte:head>

<div class="flex flex-col gap-8 pb-16">
	<!-- Page Header -->
	<header class="flex flex-col gap-1.5">
		<h1 class="font-sans font-semibold text-2xl sm:text-3xl text-(--text-main) tracking-tight">
			Server Settings & Access Control
		</h1>
		<p class="text-sm sm:text-base text-(--text-secondary) leading-relaxed">
			Manage member registration policy and generate cryptographic invitation tokens for private community access.
		</p>
	</header>

	<!-- Section 1: Registration Mode Settings -->
	<section class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-6 sm:p-7 flex flex-col gap-5">
		<div class="flex items-center gap-2 pb-3 border-b border-(--border-hairline)">
			<ShieldCheck size={20} weight="bold" class="text-(--accent-sky)" />
			<h2 class="font-sans font-semibold text-lg text-(--text-main)">
				Member Registration Policy
			</h2>
		</div>

		{#if settingsError}
			<div class="p-3.5 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-sm">
				{settingsError}
			</div>
		{/if}
		{#if settingsSuccess}
			<div class="p-3.5 bg-(--accent-soft) text-(--accent-strong) border border-(--accent-sky)/30 rounded-md text-sm flex items-center gap-2">
				<CheckCircle size={17} weight="bold" />
				<span>{settingsSuccess}</span>
			</div>
		{/if}

		{#if loadingSettings}
			<p class="text-sm text-(--text-muted)">Loading server registration policy...</p>
		{:else}
			{#if registrationMode !== savedRegistrationMode || serverDomain !== savedServerDomain}
				<div class="p-3.5 bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/30 rounded-md text-sm flex items-center justify-between">
					<div class="flex items-center gap-2">
						<WarningCircle size={18} weight="bold" />
						<span>You have unsaved changes in server settings.</span>
					</div>
					<span class="text-xs underline font-medium">Click "Save Configuration" below</span>
				</div>
			{/if}

			<form on:submit|preventDefault={handleSaveSettingsRequest} class="flex flex-col gap-4">
				<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
					<!-- Option 1: INVITE_ONLY (Recommended) -->
					<label class="relative flex flex-col p-4 rounded-md border transition-all cursor-pointer {registrationMode === 'INVITE_ONLY' ? 'border-(--accent-sky) bg-(--accent-soft)/20' : 'border-(--border-hairline) bg-(--bg-muted)'}">
						<div class="flex items-center justify-between mb-2">
							<span class="text-sm font-semibold text-(--text-main)">Invite Only</span>
							<input type="radio" name="regMode" value="INVITE_ONLY" bind:group={registrationMode} class="accent-(--accent-sky)" />
						</div>
						<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">
							(Homelab Default) Registration strictly requires a valid cryptographic invitation token from an administrator.
						</p>
					</label>

					<!-- Option 2: OPEN -->
					<label class="relative flex flex-col p-4 rounded-md border transition-all cursor-pointer {registrationMode === 'OPEN' ? 'border-(--accent-sky) bg-(--accent-soft)/20' : 'border-(--border-hairline) bg-(--bg-muted)'}">
						<div class="flex items-center justify-between mb-2">
							<span class="text-sm font-semibold text-(--text-main)">Open Registration</span>
							<input type="radio" name="regMode" value="OPEN" bind:group={registrationMode} class="accent-(--accent-sky)" />
						</div>
						<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">
							Anyone can register directly without an invitation token.
						</p>
					</label>

					<!-- Option 3: CLOSED -->
					<label class="relative flex flex-col p-4 rounded-md border transition-all cursor-pointer {registrationMode === 'CLOSED' ? 'border-(--accent-sky) bg-(--accent-soft)/20' : 'border-(--border-hairline) bg-(--bg-muted)'}">
						<div class="flex items-center justify-between mb-2">
							<span class="text-sm font-semibold text-(--text-main)">Registrations Closed</span>
							<input type="radio" name="regMode" value="CLOSED" bind:group={registrationMode} class="accent-(--accent-sky)" />
						</div>
						<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">
							Registration page is completely disabled. No new member accounts can be registered.
						</p>
					</label>
				</div>

				<!-- Canonical Node Domain -->
				<div class="flex flex-col gap-1.5 pt-3 border-t border-(--border-hairline)">
					<label for="server-domain" class="text-sm font-semibold text-(--text-main) flex items-center justify-between">
						<span class="flex items-center gap-1.5">
							<Globe size={16} class="text-(--accent-sky)" />
							<span>Canonical Node Domain</span>
						</span>
						<span class="text-xs text-(--text-muted) font-mono">e.g. ngumpul.example.com</span>
					</label>
					<input
						id="server-domain"
						type="text"
						bind:value={serverDomain}
						placeholder="ngumpul.local"
						class="w-full px-3.5 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs sm:text-sm text-(--text-main) outline-none focus:border-(--accent-sky) font-mono transition-colors"
					/>
					<p class="text-xs text-(--text-muted)">
						Canonical host address used for project subdomain routing, invite links, and telemetry display.
					</p>
				</div>

				<div class="flex items-center justify-between pt-2">
					<span class="text-xs text-(--text-muted)">
						Changes take effect immediately across all public interfaces and auth endpoints.
					</span>
					<button
						type="submit"
						class="btn btn-primary btn-sm text-sm px-4 py-2 {registrationMode !== savedRegistrationMode || serverDomain !== savedServerDomain ? 'ring-2 ring-(--accent-sky)' : ''}"
						disabled={savingSettings}
					>
						{savingSettings ? 'Saving...' : 'Save Configuration'}
					</button>
				</div>
			</form>
		{/if}
	</section>

	<!-- Section 2: Create Invitation Generator -->
	<section class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-6 sm:p-7 flex flex-col gap-5">
		<div class="flex items-center gap-2 pb-3 border-b border-(--border-hairline)">
			<Ticket size={20} weight="bold" class="text-(--accent-sky)" />
			<h2 class="font-sans font-semibold text-lg text-(--text-main)">
				Generate Invitation Token
			</h2>
		</div>

		{#if savedRegistrationMode === 'CLOSED'}
			<div class="p-4 bg-amber-500/10 border border-amber-500/30 rounded-md flex items-start gap-2.5 text-sm text-amber-600 dark:text-amber-400">
				<WarningCircle size={20} weight="bold" class="shrink-0 mt-0.5" />
				<div class="flex flex-col gap-1">
					<strong class="font-semibold">Notice: Server registration is currently CLOSED</strong>
					<span class="text-xs sm:text-sm leading-relaxed">
						While registration is closed, newly generated invitations cannot be used to create accounts. You must switch registration policy to "Invite Only" or "Open" before recipients can redeem invitations.
					</span>
				</div>
			</div>
		{/if}

		{#if inviteSuccessToken}
			<div class="p-4 bg-emerald-500/10 border border-emerald-500/30 rounded-md flex flex-col gap-3">
				<div class="flex items-center gap-2 text-emerald-600 dark:text-emerald-400 text-sm font-semibold">
					<CheckCircle size={18} weight="bold" />
					<span>Invitation Token Created Successfully!</span>
				</div>
				<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">
					Copy this URL and share it with the intended member. The raw token is displayed only once for security.
				</p>
				<div class="flex items-center gap-2">
					<input
						type="text"
						readonly
						value={inviteSuccessToken.invite_url}
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm font-mono text-(--text-main) select-all outline-none"
					/>
					<button
						type="button"
						class="btn btn-primary btn-sm text-sm shrink-0 inline-flex items-center gap-1.5 px-3 py-2"
						on:click={() => copyToClipboard(inviteSuccessToken?.invite_url || '')}
					>
						{#if copied}
							<Check size={16} weight="bold" />
							<span>Copied</span>
						{:else}
							<Copy size={16} weight="regular" />
							<span>Copy Link</span>
						{/if}
					</button>
				</div>
			</div>
		{/if}

		<form on:submit|preventDefault={handleCreateInvitation} class="flex flex-col gap-4">
			<div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
				<div class="flex flex-col gap-1.5">
					<label for="inv-email" class="text-sm font-medium text-(--text-main)">
						Restrict to Email (Optional)
					</label>
					<input
						id="inv-email"
						type="email"
						placeholder="friend@example.com"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main)"
						bind:value={invitedEmail}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="inv-uses" class="text-sm font-medium text-(--text-main)">
						Max Usage Limit
					</label>
					<input
						id="inv-uses"
						type="number"
						min="1"
						max="100"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main)"
						bind:value={maxUses}
						required
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="inv-expiry" class="text-sm font-medium text-(--text-main)">
						Expiration (Days)
					</label>
					<input
						id="inv-expiry"
						type="number"
						min="1"
						max="365"
						class="w-full px-3.5 py-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-sm text-(--text-main) outline-none focus:border-(--text-main)"
						bind:value={expiresInDays}
						required
					/>
				</div>
			</div>

			<div class="flex justify-end pt-1">
				<button type="submit" class="btn btn-primary btn-sm text-sm px-4 py-2 inline-flex items-center gap-1.5" disabled={creatingInvite}>
					<Plus size={16} weight="bold" />
					<span>{creatingInvite ? 'Generating Token...' : 'Generate Invitation Token'}</span>
				</button>
			</div>
		</form>
	</section>	<!-- Section 3: Active & Historical Invitations Table -->
	<section class="flex flex-col gap-4">
		<div class="flex items-center justify-between flex-wrap gap-2">
			<div>
				<h2 class="font-sans font-semibold text-lg text-(--text-main)">
					Issued Invitation Tokens
				</h2>
				<p class="text-sm text-(--text-muted) mt-0.5">
					Historical record of access tokens issued by server administrators.
				</p>
			</div>
			<button type="button" class="btn btn-secondary btn-sm text-sm px-3.5 py-1.5" on:click={loadInvitations}>
				Refresh
			</button>
		</div>

		<Table
			{columns}
			items={invitations}
			loading={loadingInvitations}
			loadingMessage="Loading invitations..."
			emptyMessage="No invitation tokens have been issued yet."
			let:item={inv}
		>
			{@const status = getStatusLabel(inv)}
			<TableRow>
				<TableCell class="font-medium">
					{inv.creator_name || 'Admin'}
				</TableCell>
				<TableCell class="text-(--text-secondary)">
					{inv.invited_email || '— (Any email)'}
				</TableCell>
				<TableCell>
					{#if inv.token}
						<button
							type="button"
							class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded bg-(--bg-muted) hover:bg-(--accent-soft)/30 border border-(--border-hairline) hover:border-(--accent-sky)/40 text-xs font-mono transition-all cursor-pointer group"
							title="Click to copy invitation link"
							on:click={() => copyInviteLink(inv)}
						>
							{#if copiedTokenId === inv.id}
								<Check size={14} weight="bold" class="text-emerald-500" />
								<span class="text-emerald-500 font-sans font-medium text-xs">Copied!</span>
							{:else}
								<Copy size={14} class="text-(--text-muted) group-hover:text-(--accent-sky)" />
								<span class="text-(--text-secondary) group-hover:text-(--text-main)">
									...{inv.token.slice(-6)}
								</span>
							{/if}
						</button>
					{:else}
						<span class="text-(--text-muted) text-xs font-mono">—</span>
					{/if}
				</TableCell>
				<TableCell mono>
					{inv.used_count} / {inv.max_uses}
				</TableCell>
				<TableCell class="text-(--text-muted)">
					{formatDate(inv.expires_at)}
				</TableCell>
				<TableCell>
					<span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium border {status.color}">
						{status.text}
					</span>
				</TableCell>
				<TableCell align="right">
					{#if !inv.revoked_at && (!inv.expires_at || new Date(inv.expires_at).getTime() > Date.now()) && inv.used_count < inv.max_uses}
						<button
							type="button"
							class="text-(--color-danger) hover:underline text-sm cursor-pointer font-medium"
							on:click={() => openRevokeModal(inv.id)}
						>
							Revoke
						</button>
					{:else}
						<span class="text-(--text-muted) text-xs">—</span>
					{/if}
				</TableCell>
			</TableRow>
		</Table>
	</section>

	<!-- Confirm: Close Server Registrations -->
	<ConfirmModal
		bind:open={confirmCloseModalOpen}
		confirmLabel="Close registrations"
		destructive={true}
		on:confirm={executeSaveSettings}
		on:cancel={() => (confirmCloseModalOpen = false)}
	>
		<svelte:fragment slot="title">Close server registrations?</svelte:fragment>
		<p>
			While registrations are closed, <strong class="font-medium text-(--text-main)">no new member accounts can be created on this host</strong> —
			even if someone holds a valid invitation token. Existing members will still be able to sign in.
		</p>
	</ConfirmModal>

	<!-- Confirm: Revoke Invitation Token -->
	<ConfirmModal
		open={revokeTargetId !== null}
		confirmLabel="Revoke token"
		destructive={true}
		loading={revokeLoading}
		error={revokeError}
		on:confirm={confirmRevoke}
		on:cancel={() => { revokeTargetId = null; revokeError = ''; }}
	>
		<svelte:fragment slot="title">Revoke invitation token?</svelte:fragment>
		<p>
			This token will be <strong class="font-medium text-(--text-main)">immediately invalidated</strong>.
			Anyone who hasn't used it yet will no longer be able to register with this link.
		</p>
	</ConfirmModal>
</div>
