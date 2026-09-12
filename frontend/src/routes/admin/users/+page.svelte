<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi, extractError } from '$lib/api';
	import type { User } from '$lib/types/user';
	import { user } from '$lib/stores/auth';
	import { ShieldCheck, WarningCircle, CheckCircle } from 'phosphor-svelte';
	import { Table, TableRow, TableCell, type TableColumn, ConfirmModal } from '$lib/components/ui';

	const columns: TableColumn[] = [
		{ key: 'member', label: 'Member' },
		{ key: 'email', label: 'Email Address' },
		{ key: 'projects', label: 'Projects' },
		{ key: 'role', label: 'Role' },
		{ key: 'status', label: 'Status' },
		{ key: 'actions', label: 'Actions', align: 'right' }
	];

	type ModalAction = 'PROMOTE' | 'DEMOTE' | 'SUSPEND' | 'RESTORE';

	let users: User[] = [];
	let loading = true;
	let error: string | null = null;
	let success: string | null = null;
	let searchQuery = '';
	let roleFilter = 'ALL';

	// Modal State
	let activeModal: {
		action: ModalAction;
		user: any;
	} | null = null;
	let confirmInput = '';
	let isSubmitting = false;

	async function loadUsers() {
		loading = true;
		try {
			const res = await adminApi.users.listUsers();
			users = res.users || [];
		} catch (err) {
			console.error('Failed to load users:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadUsers);

	function openModal(action: ModalAction, targetUser: any) {
		activeModal = { action, user: targetUser };
		confirmInput = '';
		error = null;
	}

	function closeModal() {
		activeModal = null;
		confirmInput = '';
	}

	$: expectedPhrase =
		activeModal?.action === 'PROMOTE' && activeModal?.user
			? `give admin role to ${activeModal.user.username}`
			: '';

	$: isPhraseValid =
		confirmInput.trim().toLowerCase() === expectedPhrase.toLowerCase();

	async function executeAction() {
		if (!activeModal || isSubmitting) return;

		if (activeModal.action === 'PROMOTE' && !isPhraseValid) {
			return;
		}

		const { action, user: target } = activeModal;
		isSubmitting = true;
		error = null;
		success = null;

		try {
			if (action === 'PROMOTE') {
				await adminApi.users.updateUserRole(target.id, 'ADMIN');
				success = `Administrator role successfully granted to ${target.display_name} (@${target.username}).`;
			} else if (action === 'DEMOTE') {
				await adminApi.users.updateUserRole(target.id, 'USER');
				success = `Administrator role successfully revoked from ${target.display_name} (@${target.username}).`;
			} else if (action === 'SUSPEND') {
				await adminApi.users.updateUserStatus(target.id, 'SUSPENDED');
				success = `Account ${target.display_name} (@${target.username}) has been suspended and active sessions terminated.`;
			} else if (action === 'RESTORE') {
				await adminApi.users.updateUserStatus(target.id, 'ACTIVE');
				success = `Account ${target.display_name} (@${target.username}) has been restored to active status.`;
			}
			closeModal();
			await loadUsers();
		} catch (err) {
			error = extractError(err);
			closeModal();
		} finally {
			isSubmitting = false;
		}
	}

	$: filteredUsers = users.filter((u) => {
		const matchesRole =
			roleFilter === 'ALL' ||
			u.role === roleFilter ||
			(roleFilter === 'SUSPENDED' && u.status === 'SUSPENDED');
		const q = searchQuery.toLowerCase();
		const matchesSearch =
			!searchQuery ||
			u.display_name?.toLowerCase().includes(q) ||
			u.username?.toLowerCase().includes(q) ||
			u.email?.toLowerCase().includes(q);
		return matchesRole && matchesSearch;
	});
</script>


<div class="flex flex-col gap-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-(--border-hairline)">
		<div>
			<h1 class="font-display font-bold text-2xl text-(--text-main)">Member Moderation</h1>
			<p class="text-xs text-(--text-secondary) mt-1">
				Directory of community accounts, authorization levels, and active participation permissions.
			</p>
		</div>

		<div class="flex items-center gap-2">
			<button type="button" class="btn btn-secondary btn-sm font-mono text-xs" on:click={loadUsers}>
				Refresh
			</button>
		</div>
	</div>

	{#if success}
		<div class="p-3 bg-(--accent-soft) text-(--accent-strong) border border-(--accent-sky)/30 rounded text-xs font-mono flex items-center justify-between">
			<span>{success}</span>
			<button type="button" on:click={() => (success = null)} class="text-(--accent-strong) font-bold">×</button>
		</div>
	{/if}
	{#if error}
		<div class="p-3 bg-red-950/20 text-(--color-danger) border border-(--color-danger)/30 rounded text-xs font-mono flex items-center justify-between">
			<span>{error}</span>
			<button type="button" on:click={() => (error = null)} class="text-(--color-danger) font-bold">×</button>
		</div>
	{/if}

	<!-- Search & Filter Controls -->
	<div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
		<div class="flex items-center gap-1 border border-(--border-hairline) rounded p-1 bg-(--bg-surface)">
			{#each [
				{ id: 'ALL', label: 'All' },
				{ id: 'ADMIN', label: 'Admins' },
				{ id: 'USER', label: 'Members' },
				{ id: 'SUSPENDED', label: 'Suspended' }
			] as tab}
				<button
					type="button"
					class="px-3 py-1.5 text-sm font-medium rounded transition-colors {roleFilter === tab.id ? 'bg-(--accent-soft) text-(--accent-strong) font-semibold' : 'text-(--text-secondary) hover:text-(--text-main)'}"
					on:click={() => (roleFilter = tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</div>

		<div class="w-full sm:w-72">
			<input
				type="text"
				placeholder="Search by name, handle, email..."
				bind:value={searchQuery}
				class="w-full px-3.5 py-2 bg-(--bg-surface) border border-(--border-hairline) rounded text-sm text-(--text-main) placeholder:text-(--text-muted) outline-none focus:border-(--accent-sky) font-sans"
			/>
		</div>
	</div>

	<!-- High-Density Members Reusable Table -->
	<Table
		{columns}
		items={filteredUsers}
		{loading}
		loadingMessage="Loading registered members..."
		emptyMessage="No member accounts match current criteria."
		let:item={u}
	>
		<TableRow>
			<TableCell>
				<div class="flex items-center gap-2.5">
					<div class="w-8 h-8 rounded bg-(--bg-muted) text-(--text-main) font-mono text-xs font-bold flex items-center justify-center border border-(--border-hairline) shrink-0">
						{u.username ? u.username.slice(0, 2).toUpperCase() : 'U'}
					</div>
					<div>
						<div class="font-semibold text-sm sm:text-base text-(--text-main) flex items-center gap-1.5">
							<span>{u.display_name}</span>
							{#if u.id === $user?.id}
								<span class="text-xs font-mono px-2 py-0.5 bg-(--accent-soft) text-(--accent-strong) rounded border border-(--accent-sky)/30">
									You
								</span>
							{/if}
						</div>
						<div class="text-xs font-mono text-(--text-muted)">@{u.username}</div>
					</div>
				</div>
			</TableCell>
			<TableCell mono class="text-sm text-(--text-secondary)">
				{u.email}
			</TableCell>
			<TableCell mono class="text-sm font-medium">
				{u.projects_count || 0}
			</TableCell>
			<TableCell>
				<span class="inline-flex items-center gap-1.5 text-xs font-mono font-medium px-2.5 py-1 rounded border {u.role === 'ADMIN' ? 'bg-(--accent-soft) text-(--accent-strong) border-(--accent-sky)/40' : 'bg-(--bg-muted) text-(--text-secondary) border-(--border-hairline)'}">
					{u.role}
				</span>
			</TableCell>
			<TableCell>
				<span class="inline-flex items-center gap-1.5 text-xs font-mono font-medium px-2.5 py-1 rounded border {u.status === 'ACTIVE' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/25' : 'bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/25'}">
					<span class="w-2 h-2 rounded-full {u.status === 'ACTIVE' ? 'bg-emerald-500' : 'bg-red-500'}"></span>
					{u.status}
				</span>
			</TableCell>
			<TableCell align="right" class="whitespace-nowrap">
				<div class="flex items-center justify-end gap-2">
					{#if u.id === $user?.id}
						<span class="inline-flex items-center text-xs font-mono text-(--text-muted) px-2.5 py-1 rounded bg-(--bg-muted) border border-(--border-hairline)">
							Your Account
						</span>
					{:else}
						<button
							type="button"
							class="btn btn-secondary btn-sm text-xs py-1 px-2.5 font-mono"
							on:click={() => openModal(u.role === 'ADMIN' ? 'DEMOTE' : 'PROMOTE', u)}
						>
							{u.role === 'ADMIN' ? 'Demote' : 'Make Admin'}
						</button>
						<button
							type="button"
							class="btn btn-secondary btn-sm text-xs py-1 px-2.5 font-mono {u.status === 'ACTIVE' ? 'text-(--color-danger)' : 'text-emerald-600 dark:text-emerald-400'}"
							on:click={() => openModal(u.status === 'ACTIVE' ? 'SUSPEND' : 'RESTORE', u)}
						>
							{u.status === 'ACTIVE' ? 'Suspend' : 'Restore'}
						</button>
					{/if}
				</div>
			</TableCell>
		</TableRow>
	</Table>

	<!-- Administrative Action Modals (using ConfirmModal foundation) -->

	<!-- PROMOTE: requires typed confirmation phrase -->
	<ConfirmModal
		open={activeModal?.action === 'PROMOTE'}
		confirmLabel="Grant admin role"
		destructive={false}
		loading={isSubmitting}
		disabled={!isPhraseValid}
		error={activeModal?.action === 'PROMOTE' ? (error ?? '') : ''}
		on:confirm={executeAction}
		on:cancel={closeModal}
	>
		<svelte:fragment slot="title">Promote to administrator?</svelte:fragment>
		<p>
			Granting administrator role to <strong class="text-(--text-main)">{activeModal?.user.display_name}</strong>
			(<span class="font-mono">@{activeModal?.user.username}</span>) gives full access to server
			configurations, invitation token generation, audit logs, and member moderation.
		</p>
		<svelte:fragment slot="extra">
			<div class="flex flex-col gap-1.5 p-3.5 bg-(--bg-muted) rounded-[6px] border border-(--border-hairline)">
				<span class="text-xs text-(--text-muted)">Type the following phrase to proceed:</span>
				<code class="px-2.5 py-1.5 bg-(--bg-surface) border border-(--border-hairline) rounded-[4px] font-mono text-sm text-(--text-main) font-semibold select-all">
					{expectedPhrase}
				</code>
			</div>
			<input
				type="text"
				bind:value={confirmInput}
				placeholder={expectedPhrase}
				class="mt-2 w-full px-3.5 py-2.5 bg-(--bg-surface) border border-(--border-hairline) rounded-[4px] text-sm font-mono text-(--text-main) placeholder:text-(--text-muted)/40 outline-none focus:border-(--accent-sky)"
				on:keydown={(e) => { if (e.key === 'Enter' && isPhraseValid) executeAction(); }}
			/>
		</svelte:fragment>
	</ConfirmModal>

	<!-- DEMOTE: destructive -->
	<ConfirmModal
		open={activeModal?.action === 'DEMOTE'}
		confirmLabel="Revoke admin role"
		destructive={true}
		loading={isSubmitting}
		error={activeModal?.action === 'DEMOTE' ? (error ?? '') : ''}
		on:confirm={executeAction}
		on:cancel={closeModal}
	>
		<svelte:fragment slot="title">Revoke administrator role?</svelte:fragment>
		<p>
			<strong class="text-(--text-main)">{activeModal?.user.display_name}</strong>
			(<span class="font-mono">@{activeModal?.user.username}</span>) will be demoted to standard
			member status and lose access to the management console, hosting approvals, and server
			configurations.
		</p>
	</ConfirmModal>

	<!-- SUSPEND: destructive -->
	<ConfirmModal
		open={activeModal?.action === 'SUSPEND'}
		confirmLabel="Suspend account"
		destructive={true}
		loading={isSubmitting}
		error={activeModal?.action === 'SUSPEND' ? (error ?? '') : ''}
		on:confirm={executeAction}
		on:cancel={closeModal}
	>
		<svelte:fragment slot="title">Suspend user account?</svelte:fragment>
		<p>
			All active sessions for <strong class="text-(--text-main)">{activeModal?.user.display_name}</strong>
			(<span class="font-mono">@{activeModal?.user.username}</span>) will be immediately terminated.
			Access will be blocked until restored by an administrator.
		</p>
	</ConfirmModal>

	<!-- RESTORE: non-destructive -->
	<ConfirmModal
		open={activeModal?.action === 'RESTORE'}
		confirmLabel="Restore account"
		destructive={false}
		loading={isSubmitting}
		error={activeModal?.action === 'RESTORE' ? (error ?? '') : ''}
		on:confirm={executeAction}
		on:cancel={closeModal}
	>
		<svelte:fragment slot="title">Restore user account?</svelte:fragment>
		<p>
			<strong class="text-(--text-main)">{activeModal?.user.display_name}</strong>
			(<span class="font-mono">@{activeModal?.user.username}</span>) will be able to sign in again
			and manage their projects and hosting applications.
		</p>
	</ConfirmModal>
</div>

