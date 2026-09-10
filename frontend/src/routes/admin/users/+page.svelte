<script lang="ts">
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import { user } from '$lib/stores/auth';

	let users: any[] = [];
	let loading = true;
	let error: string | null = null;
	let success: string | null = null;
	let searchQuery = '';
	let roleFilter = 'ALL';

	async function loadUsers() {
		loading = true;
		try {
			const res = await api.get('/admin/users');
			users = res.data?.users || [];
		} catch (err) {
			console.error('Failed to load users:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadUsers);

	async function toggleRole(u: any) {
		const newRole = u.role === 'ADMIN' ? 'USER' : 'ADMIN';
		error = null;
		success = null;
		try {
			await api.patch(`/admin/users/${u.id}/role`, { role: newRole });
			success = `Updated role for ${u.display_name} to ${newRole}.`;
			await loadUsers();
		} catch (err) {
			error = extractError(err);
		}
	}

	async function toggleStatus(u: any) {
		if (u.id === $user?.id) {
			error = 'You cannot suspend your own active administrative account.';
			return;
		}
		const newStatus = u.status === 'ACTIVE' ? 'SUSPENDED' : 'ACTIVE';
		error = null;
		success = null;
		try {
			await api.patch(`/admin/users/${u.id}/status`, { status: newStatus });
			success = `Account ${u.display_name} is now ${newStatus}.`;
			await loadUsers();
		} catch (err) {
			error = extractError(err);
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
		<div class="flex items-center gap-1 border border-(--border-hairline) rounded p-0.5 bg-(--bg-surface)">
			{#each [
				{ id: 'ALL', label: 'All' },
				{ id: 'ADMIN', label: 'Admins' },
				{ id: 'USER', label: 'Members' },
				{ id: 'SUSPENDED', label: 'Suspended' }
			] as tab}
				<button
					type="button"
					class="px-2.5 py-1 text-xs font-medium rounded transition-colors {roleFilter === tab.id ? 'bg-(--accent-soft) text-(--accent-strong) font-semibold' : 'text-(--text-secondary) hover:text-(--text-main)'}"
					on:click={() => (roleFilter = tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</div>

		<div class="w-full sm:w-64">
			<input
				type="text"
				placeholder="Search by name, handle, email..."
				bind:value={searchQuery}
				class="w-full px-3 py-1.5 bg-(--bg-surface) border border-(--border-hairline) rounded text-xs text-(--text-main) placeholder:text-(--text-muted) outline-none focus:border-(--accent-sky) font-sans"
			/>
		</div>
	</div>

	<!-- High-Density Users Table -->
	<div class="border border-(--border-hairline) rounded-md bg-(--bg-surface) overflow-hidden">
		{#if loading}
			<div class="p-10 text-center text-xs text-(--text-muted)">Loading registered members...</div>
		{:else if filteredUsers.length === 0}
			<div class="p-10 text-center text-xs text-(--text-secondary)">
				No member accounts match current criteria.
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs border-collapse font-sans">
					<thead>
						<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-muted) text-xs">
							<th class="px-5 py-2.5 font-medium">Member</th>
							<th class="px-5 py-2.5 font-medium">Email Address</th>
							<th class="px-5 py-2.5 font-medium">Projects</th>
							<th class="px-5 py-2.5 font-medium">Role</th>
							<th class="px-5 py-2.5 font-medium">Status</th>
							<th class="px-5 py-2.5 text-right font-medium">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-(--border-hairline)">
						{#each filteredUsers as u (u.id)}
							<tr class="hover:bg-(--bg-muted)/25 transition-colors">
								<td class="px-5 py-3.5 align-middle">
									<div class="flex items-center gap-2.5">
										<div class="w-7 h-7 rounded bg-(--bg-muted) text-(--text-main) font-mono text-xs font-bold flex items-center justify-center border border-(--border-hairline) shrink-0">
											{u.username ? u.username.slice(0, 2).toUpperCase() : 'U'}
										</div>
										<div>
											<div class="font-semibold text-(--text-main)">{u.display_name}</div>
											<div class="text-[11px] font-mono text-(--text-muted)">@{u.username}</div>
										</div>
									</div>
								</td>
								<td class="px-5 py-3.5 align-middle font-mono text-[11px] text-(--text-secondary)">
									{u.email}
								</td>
								<td class="px-5 py-3.5 align-middle font-mono text-xs font-medium">
									{u.projects_count || 0}
								</td>
								<td class="px-5 py-3.5 align-middle">
									<span class="inline-flex items-center gap-1.5 text-[11px] font-mono font-medium px-2 py-0.5 rounded border {u.role === 'ADMIN' ? 'bg-(--accent-soft) text-(--accent-strong) border-(--accent-sky)/40' : 'bg-(--bg-muted) text-(--text-secondary) border-(--border-hairline)'}">
										{u.role}
									</span>
								</td>
								<td class="px-5 py-3.5 align-middle">
									<span class="inline-flex items-center gap-1.5 text-[11px] font-mono font-medium px-2 py-0.5 rounded border {u.status === 'ACTIVE' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/25' : 'bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/25'}">
										<span class="w-1.5 h-1.5 rounded-full {u.status === 'ACTIVE' ? 'bg-emerald-500' : 'bg-red-500'}"></span>
										{u.status}
									</span>
								</td>
								<td class="px-5 py-3.5 align-middle text-right whitespace-nowrap">
									<div class="flex items-center justify-end gap-1.5">
										<button
											type="button"
											class="btn btn-secondary btn-sm text-[11px] py-0.5 px-2 font-mono"
											on:click={() => toggleRole(u)}
										>
											{u.role === 'ADMIN' ? 'Demote' : 'Make Admin'}
										</button>
										{#if u.id !== $user?.id}
											<button
												type="button"
												class="btn btn-secondary btn-sm text-[11px] py-0.5 px-2 font-mono {u.status === 'ACTIVE' ? 'text-(--color-danger)' : 'text-emerald-600'}"
												on:click={() => toggleStatus(u)}
											>
												{u.status === 'ACTIVE' ? 'Suspend' : 'Restore'}
											</button>
										{/if}
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
