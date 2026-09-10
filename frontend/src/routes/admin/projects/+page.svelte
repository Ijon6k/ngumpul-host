<script lang="ts">
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';

	let projects: any[] = [];
	let loading = true;
	let error: string | null = null;
	let success: string | null = null;
	let searchQuery = '';
	let statusFilter = 'ALL';

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
			<p class="text-xs text-(--text-secondary) mt-1">
				Configure publication states, runtime status overrides, and public endpoints for community projects.
			</p>
		</div>

		<div class="flex items-center gap-2">
			<button type="button" class="btn btn-secondary btn-sm font-mono text-xs" on:click={loadProjects}>
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
				{ id: 'ONLINE', label: 'Online' },
				{ id: 'SETUP', label: 'Setup' },
				{ id: 'OFFLINE', label: 'Offline' },
				{ id: 'ARCHIVED', label: 'Archived' }
			] as tab}
				<button
					type="button"
					class="px-2.5 py-1 text-xs font-medium rounded transition-colors {statusFilter === tab.id ? 'bg-(--accent-soft) text-(--accent-strong) font-semibold' : 'text-(--text-secondary) hover:text-(--text-main)'}"
					on:click={() => (statusFilter = tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</div>

		<div class="w-full sm:w-64">
			<input
				type="text"
				placeholder="Search projects, slugs, owners..."
				bind:value={searchQuery}
				class="w-full px-3 py-1.5 bg-(--bg-surface) border border-(--border-hairline) rounded text-xs text-(--text-main) placeholder:text-(--text-muted) outline-none focus:border-(--accent-sky) font-sans"
			/>
		</div>
	</div>

	<!-- High-Density Projects Table -->
	<div class="border border-(--border-hairline) rounded-md bg-(--bg-surface) overflow-hidden">
		{#if loading}
			<div class="p-10 text-center text-xs text-(--text-muted)">Loading registered projects...</div>
		{:else if filteredProjects.length === 0}
			<div class="p-10 text-center text-xs text-(--text-secondary)">
				No projects match the selected criteria.
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs border-collapse font-sans">
					<thead>
						<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-muted) text-xs">
							<th class="px-5 py-2.5 font-medium">Project / Slug</th>
							<th class="px-5 py-2.5 font-medium">Owner</th>
							<th class="px-5 py-2.5 font-medium">Public Endpoint</th>
							<th class="px-5 py-2.5 font-medium">Visibility</th>
							<th class="px-5 py-2.5 font-medium">Status</th>
							<th class="px-5 py-2.5 text-right font-medium">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-(--border-hairline)">
						{#each filteredProjects as proj (proj.id)}
							<tr class="hover:bg-(--bg-muted)/25 transition-colors">
								<td class="px-5 py-3.5 align-middle">
									<a href="/projects/{proj.slug}" target="_blank" class="font-semibold text-(--text-main) hover:text-(--accent-strong) transition-colors">
										{proj.name}
									</a>
									<span class="block text-[11px] font-mono text-(--text-muted)">/{proj.slug}</span>
								</td>
								<td class="px-5 py-3.5 align-middle">
									<span class="font-medium text-(--text-main)">{proj.owner?.display_name || 'Owner'}</span>
									<span class="block text-[11px] font-mono text-(--text-muted)">@{proj.owner?.username}</span>
								</td>
								<td class="px-5 py-3.5 align-middle font-mono text-[11px]">
									{#if proj.public_url}
										<a href={proj.public_url} target="_blank" rel="noreferrer" class="text-(--accent-strong) hover:underline block truncate max-w-xs">
											{proj.public_url} ↗
										</a>
									{:else}
										<span class="text-(--text-muted)">—</span>
									{/if}
								</td>
								<td class="px-5 py-3.5 align-middle">
									<button
										type="button"
										class="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-mono font-medium cursor-pointer border transition-colors {proj.visibility === 'PUBLIC' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/25 hover:bg-emerald-500/20' : 'bg-(--bg-muted) text-(--text-secondary) border-(--border-hairline) hover:bg-(--bg-hover)'}"
										on:click={() => toggleVisibility(proj.id, proj.visibility)}
										title="Toggle visibility"
									>
										{proj.visibility}
									</button>
								</td>
								<td class="px-5 py-3.5 align-middle">
									<StatusDot status={proj.status} />
								</td>
								<td class="px-5 py-3.5 align-middle text-right">
									<select
										class="px-2 py-1 text-[11px] font-mono bg-(--bg-surface) border border-(--border-hairline) rounded text-(--text-main) outline-none focus:border-(--accent-sky)"
										value={proj.status}
										on:change={(e) => updateProjectStatus(proj.id, e.currentTarget.value)}
									>
										<option value="ONLINE">ONLINE</option>
										<option value="SETUP">SETUP</option>
										<option value="OFFLINE">OFFLINE</option>
										<option value="ARCHIVED">ARCHIVED</option>
									</select>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
