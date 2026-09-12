<script lang="ts">
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import { Table, TableRow, TableCell, type TableColumn } from '$lib/components/ui';

	let projects: any[] = [];
	let loading = true;
	let error: string | null = null;
	let success: string | null = null;
	let searchQuery = '';
	let statusFilter = 'ALL';

	const columns: TableColumn[] = [
		{ key: 'project', label: 'Project / Slug' },
		{ key: 'owner', label: 'Owner' },
		{ key: 'endpoint', label: 'Public Endpoint' },
		{ key: 'visibility', label: 'Visibility' },
		{ key: 'status', label: 'Status' },
		{ key: 'actions', label: 'Actions', align: 'right' }
	];

	async function loadProjects() {
		loading = true;
		try {
			const res = await api.get('/admin/projects');
			projects = res.data?.projects || [];
		} catch (err) {
			console.error('Failed to load admin projects:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadProjects);

	async function updateProjectStatus(projId: string, newStatus: string) {
		error = null;
		success = null;
		try {
			await api.patch(`/admin/projects/${projId}`, {
				status: newStatus
			});
			success = 'Project status successfully updated.';
			await loadProjects();
		} catch (err) {
			error = extractError(err);
		}
	}

	async function toggleVisibility(projId: string, currentVis: string) {
		error = null;
		success = null;
		const newVis = currentVis === 'PUBLIC' ? 'UNPUBLISHED' : 'PUBLIC';
		try {
			await api.patch(`/admin/projects/${projId}`, {
				visibility: newVis
			});
			success = `Visibility changed to ${newVis}.`;
			await loadProjects();
		} catch (err) {
			error = extractError(err);
		}
	}

	$: filteredProjects = projects.filter((p) => {
		const matchesStatus = statusFilter === 'ALL' || p.status === statusFilter;
		const q = searchQuery.toLowerCase();
		const matchesSearch =
			!searchQuery ||
			p.name?.toLowerCase().includes(q) ||
			p.slug?.toLowerCase().includes(q) ||
			p.owner?.display_name?.toLowerCase().includes(q) ||
			p.owner?.username?.toLowerCase().includes(q) ||
			p.public_url?.toLowerCase().includes(q);
		return matchesStatus && matchesSearch;
	});
</script>

<div class="flex flex-col gap-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-(--border-hairline)">
		<div>
			<h1 class="font-display font-bold text-2xl text-(--text-main)">Project Management</h1>
			<p class="text-sm text-(--text-secondary) mt-1">
				Configure publication states, runtime status overrides, and public endpoints for community projects.
			</p>
		</div>

		<div class="flex items-center gap-2">
			<button type="button" class="btn btn-secondary btn-sm font-mono text-xs sm:text-sm" on:click={loadProjects}>
				Refresh
			</button>
		</div>
	</div>

	{#if success}
		<div class="p-3 bg-(--accent-soft) text-(--accent-strong) border border-(--accent-sky)/30 rounded text-sm font-mono flex items-center justify-between">
			<span>{success}</span>
			<button type="button" on:click={() => (success = null)} class="text-(--accent-strong) font-bold">×</button>
		</div>
	{/if}
	{#if error}
		<div class="p-3 bg-red-950/20 text-(--color-danger) border border-(--color-danger)/30 rounded text-sm font-mono flex items-center justify-between">
			<span>{error}</span>
			<button type="button" on:click={() => (error = null)} class="text-(--color-danger) font-bold">×</button>
		</div>
	{/if}

	<!-- Search & Filter Controls -->
	<div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
		<div class="flex items-center gap-1 border border-(--border-hairline) rounded p-1 bg-(--bg-surface)">
			{#each [
				{ id: 'ALL', label: 'All' },
				{ id: 'ONLINE', label: 'Online' },
				{ id: 'SETUP', label: 'Setup' },
				{ id: 'OFFLINE', label: 'Offline' },
				{ id: 'ARCHIVED', label: 'Archived' }
			] as tab}
				<button
					type="button"
					class="px-3 py-1.5 text-sm font-medium rounded transition-colors {statusFilter === tab.id ? 'bg-(--accent-soft) text-(--accent-strong) font-semibold' : 'text-(--text-secondary) hover:text-(--text-main)'}"
					on:click={() => (statusFilter = tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</div>

		<div class="w-full sm:w-72">
			<input
				type="text"
				placeholder="Search projects, slugs, owners..."
				bind:value={searchQuery}
				class="w-full px-3.5 py-2 bg-(--bg-surface) border border-(--border-hairline) rounded text-sm text-(--text-main) placeholder:text-(--text-muted) outline-none focus:border-(--accent-sky) font-sans"
			/>
		</div>
	</div>

	<!-- High-Density Projects Reusable Table -->
	<Table
		{columns}
		items={filteredProjects}
		{loading}
		loadingMessage="Loading registered projects..."
		emptyMessage="No projects match the selected criteria."
		let:item={proj}
	>
		<TableRow>
			<TableCell>
				<a href="/admin/projects/{proj.id}" class="font-semibold text-sm text-(--text-main) hover:text-(--accent-sky) transition-colors">
					{proj.name}
				</a>
				<span class="block text-xs font-mono text-(--text-muted)">/{proj.slug}</span>
			</TableCell>
			<TableCell>
				<span class="font-medium text-sm text-(--text-main)">{proj.owner?.display_name || 'Owner'}</span>
				<span class="block text-xs font-mono text-(--text-muted)">@{proj.owner?.username}</span>
			</TableCell>
			<TableCell mono class="text-xs">
				{#if proj.public_url}
					<a href={proj.public_url} target="_blank" rel="noreferrer" class="text-(--accent-sky) hover:underline block truncate max-w-xs">
						{proj.public_url} ↗
					</a>
				{:else}
					<span class="text-(--text-muted)">—</span>
				{/if}
			</TableCell>
			<TableCell>
				<button
					type="button"
					class="inline-flex items-center px-2.5 py-1 rounded-sm text-xs font-mono font-medium cursor-pointer border transition-colors {proj.visibility === 'PUBLIC' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/25 hover:bg-emerald-500/20' : 'bg-(--bg-muted) text-(--text-secondary) border-(--border-hairline) hover:bg-(--bg-hover)'}"
					on:click={() => toggleVisibility(proj.id, proj.visibility)}
					title="Toggle visibility"
				>
					{proj.visibility}
				</button>
			</TableCell>
			<TableCell>
				<StatusDot status={proj.status} />
			</TableCell>
			<TableCell align="right">
				<div class="flex items-center justify-end gap-2">
					<a
						href="/admin/projects/{proj.id}"
						class="btn btn-secondary btn-sm text-xs px-2.5 py-1"
					>
						Moderate
					</a>
					<select
						class="px-2 py-1 text-xs font-mono bg-(--bg-surface) border border-(--border-hairline) rounded-sm text-(--text-main) outline-none focus:border-(--accent-sky)"
						value={proj.status}
						on:change={(e) => updateProjectStatus(proj.id, e.currentTarget.value)}
					>
						<option value="ONLINE">ONLINE</option>
						<option value="SETUP">SETUP</option>
						<option value="OFFLINE">OFFLINE</option>
						<option value="ARCHIVED">ARCHIVED</option>
					</select>
				</div>
			</TableCell>
		</TableRow>
	</Table>
</div>
