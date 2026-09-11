<script lang="ts">
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import { Table, TableRow, TableCell, type TableColumn } from '$lib/components/ui';

	let requests: any[] = [];
	let loading = true;
	let activeFilter = 'ALL';
	let searchQuery = '';

	let activeRequest: any = null;
	let actionType: 'approve' | 'reject' | 'complete' | null = null;

	let adminNotes = '';
	let publicURL = '';
	let processing = false;
	let error: string | null = null;
	let success: string | null = null;

	async function loadRequests() {
		loading = true;
		try {
			const res = await api.get('/admin/hosting-requests');
			requests = res.data?.requests || [];
		} catch (err) {
			console.error('Failed to load requests:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadRequests);

	const columns: TableColumn[] = [
		{ key: 'project', label: 'Project & Specs' },
		{ key: 'requester', label: 'Requester' },
		{ key: 'repository', label: 'Repository' },
		{ key: 'status', label: 'Status' },
		{ key: 'submitted', label: 'Submitted' },
		{ key: 'actions', label: 'Actions', align: 'right' }
	];

	$: filteredRequests = requests.filter((r) => {
		const matchesFilter =
			activeFilter === 'ALL' ||
			r.status === activeFilter ||
			(activeFilter === 'IN_PROGRESS' && (r.status === 'APPROVED' || r.status === 'SETUP'));

		const q = searchQuery.toLowerCase();
		const matchesSearch =
			!searchQuery ||
			r.project_name?.toLowerCase().includes(q) ||
			r.repository_url?.toLowerCase().includes(q) ||
			r.requester?.display_name?.toLowerCase().includes(q) ||
			r.requester?.username?.toLowerCase().includes(q);

		return matchesFilter && matchesSearch;
	});

	function openModal(req: any, type: 'approve' | 'reject' | 'complete') {
		activeRequest = req;
		actionType = type;
		adminNotes = '';
		publicURL = `https://${req.project_name.toLowerCase().replace(/[^a-z0-9]/g, '')}.ngumpul.local`;
		error = null;
		success = null;
	}

	function closeModal() {
		activeRequest = null;
		actionType = null;
	}

	async function handleAction() {
		if (!activeRequest || !actionType) return;
		processing = true;
		error = null;

		try {
			if (actionType === 'approve') {
				await api.post(`/admin/hosting-requests/${activeRequest.id}/approve`, {
					admin_notes: adminNotes
				});
				success = `Request for ${activeRequest.project_name} marked as APPROVED.`;
			} else if (actionType === 'reject') {
				await api.post(`/admin/hosting-requests/${activeRequest.id}/reject`, {
					admin_notes: adminNotes
				});
				success = `Request for ${activeRequest.project_name} marked as REJECTED.`;
			} else if (actionType === 'complete') {
				await api.post(`/admin/hosting-requests/${activeRequest.id}/complete`, {
					admin_notes: adminNotes,
					public_url: publicURL
				});
				success = `Project ${activeRequest.project_name} published successfully at ${publicURL}.`;
			}
			closeModal();
			await loadRequests();
		} catch (err) {
			error = extractError(err);
		} finally {
			processing = false;
		}
	}
</script>

<div class="flex flex-col gap-6">
	<!-- Page Header & Action Controls -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-(--border-hairline)">
		<div>
			<h1 class="font-display font-bold text-2xl text-(--text-main)">Hosting Requests</h1>
			<p class="text-sm text-(--text-secondary) mt-1">
				Operational queue for community applications requesting server allocation and public proxy routes.
			</p>
		</div>

		<div class="flex items-center gap-2">
			<button type="button" class="btn btn-secondary btn-sm font-mono text-xs sm:text-sm" on:click={loadRequests}>
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
				{ id: 'PENDING', label: 'Pending' },
				{ id: 'IN_PROGRESS', label: 'Approved' },
				{ id: 'COMPLETED', label: 'Published' },
				{ id: 'REJECTED', label: 'Rejected' }
			] as tab}
				<button
					type="button"
					class="px-3 py-1.5 text-sm font-medium rounded transition-colors {activeFilter === tab.id ? 'bg-(--accent-soft) text-(--accent-strong) font-semibold' : 'text-(--text-secondary) hover:text-(--text-main)'}"
					on:click={() => (activeFilter = tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</div>

		<div class="w-full sm:w-72">
			<input
				type="text"
				placeholder="Filter by project or owner..."
				bind:value={searchQuery}
				class="w-full px-3.5 py-2 bg-(--bg-surface) border border-(--border-hairline) rounded text-sm text-(--text-main) placeholder:text-(--text-muted) outline-none focus:border-(--accent-sky) font-sans"
			/>
		</div>
	</div>

	<!-- High-Density Requests Reusable Table -->
	<Table
		{columns}
		items={filteredRequests}
		{loading}
		alignTop
		loadingMessage="Querying operational queue..."
		emptyMessage="No requests matching current criteria."
		let:item={req}
	>
		<TableRow>
			<TableCell alignTop>
				<div class="font-semibold text-sm sm:text-base text-(--text-main)">{req.project_name}</div>
				{#if req.description}
					<div class="text-xs text-(--text-secondary) mt-0.5 max-w-sm line-clamp-1">{req.description}</div>
				{/if}
				{#if req.environment_specs && Object.keys(req.environment_specs).length > 0}
					<div class="text-xs font-mono text-(--text-muted) mt-1">
						{JSON.stringify(req.environment_specs)}
					</div>
				{/if}
			</TableCell>
			<TableCell alignTop>
				<span class="font-medium text-sm sm:text-base text-(--text-main)">{req.requester?.display_name || 'Member'}</span>
				<span class="block text-xs font-mono text-(--text-muted)">@{req.requester?.username}</span>
			</TableCell>
			<TableCell alignTop mono class="text-xs">
				<a href={req.repository_url} target="_blank" rel="noreferrer" class="text-(--accent-strong) hover:underline block truncate max-w-xs">
					{req.repository_url} ↗
				</a>
			</TableCell>
			<TableCell alignTop>
				{#if req.status === 'PENDING'}
					<span class="inline-flex items-center gap-1.5 text-xs font-mono font-medium px-2.5 py-1 rounded bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20">
						<span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span>
						PENDING
					</span>
				{:else if req.status === 'APPROVED' || req.status === 'SETUP'}
					<span class="inline-flex items-center gap-1.5 text-xs font-mono font-medium px-2.5 py-1 rounded bg-blue-500/10 text-blue-600 dark:text-blue-400 border border-blue-500/20">
						<span class="w-1.5 h-1.5 rounded-full bg-blue-500"></span>
						{req.status}
					</span>
				{:else if req.status === 'COMPLETED'}
					<span class="inline-flex items-center gap-1.5 text-xs font-mono font-medium px-2.5 py-1 rounded bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20">
						<span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
						PUBLISHED
					</span>
				{:else}
					<span class="inline-flex items-center gap-1.5 text-xs font-mono font-medium px-2.5 py-1 rounded bg-red-500/10 text-red-600 dark:text-red-400 border border-red-500/20">
						<span class="w-1.5 h-1.5 rounded-full bg-red-500"></span>
						REJECTED
					</span>
				{/if}
			</TableCell>
			<TableCell alignTop mono class="text-xs text-(--text-muted)">
				{req.created_at ? new Date(req.created_at).toLocaleDateString([], { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) : '—'}
			</TableCell>
			<TableCell alignTop align="right" class="whitespace-nowrap">
				{#if req.status === 'PENDING'}
					<div class="flex items-center justify-end gap-1.5">
						<button
							type="button"
							class="btn btn-primary btn-sm text-xs py-1 px-2.5 font-mono"
							on:click={() => openModal(req, 'approve')}
						>
							Approve
						</button>
						<button
							type="button"
							class="btn btn-secondary btn-sm text-xs py-1 px-2.5 font-mono text-(--color-danger)"
							on:click={() => openModal(req, 'reject')}
						>
							Reject
						</button>
					</div>
				{:else if req.status === 'APPROVED' || req.status === 'SETUP'}
					<button
						type="button"
						class="btn btn-primary btn-sm text-xs py-1 px-2.5 font-mono"
						on:click={() => openModal(req, 'complete')}
					>
						Publish URL ➔
					</button>
				{:else}
					<span class="text-xs font-mono text-(--text-muted)">Closed</span>
				{/if}
			</TableCell>
		</TableRow>
	</Table>

	<!-- Operational Action Modal -->
	{#if activeRequest && actionType}
		<div class="fixed inset-0 bg-black/60 backdrop-blur-xs flex items-center justify-center z-50 p-4">
			<div class="w-full max-w-md p-6 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-4 shadow-xl">
				<div class="border-b border-(--border-hairline) pb-3">
					<div class="text-xs text-(--accent-strong) font-medium">Operator Action</div>
					<h2 class="font-display font-bold text-lg text-(--text-main) mt-0.5">
						{#if actionType === 'approve'}
							Approve Hosting: {activeRequest.project_name}
						{:else if actionType === 'reject'}
							Reject Hosting: {activeRequest.project_name}
						{:else}
							Publish Endpoint: {activeRequest.project_name}
						{/if}
					</h2>
				</div>

				{#if actionType === 'complete'}
					<div class="flex flex-col gap-1.5">
						<label for="admin-public-url" class="text-sm font-medium text-(--text-secondary)">Assigned Public URL *</label>
						<input
							id="admin-public-url"
							type="url"
							class="w-full px-3.5 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded text-sm font-mono text-(--text-main) outline-none focus:border-(--accent-sky)"
							required
							bind:value={publicURL}
						/>
						<span class="text-xs text-(--text-muted)">Must point to the configured reverse-proxy location.</span>
					</div>
				{/if}

				<div class="flex flex-col gap-1.5">
					<label for="admin-notes" class="text-sm font-medium text-(--text-secondary)">Operator Notes / Instructions</label>
					<textarea
						id="admin-notes"
						rows="3"
						class="w-full px-3.5 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded text-sm text-(--text-main) outline-none focus:border-(--accent-sky)"
						placeholder="Optional feedback communicated to requester..."
						bind:value={adminNotes}
					></textarea>
				</div>

				<div class="flex items-center justify-end gap-2 pt-3 border-t border-(--border-hairline)">
					<button type="button" class="btn btn-secondary btn-sm text-sm font-medium px-4 py-2" on:click={closeModal}>
						Cancel
					</button>
					<button
						type="button"
						class="btn btn-primary btn-sm text-sm font-medium px-4 py-2"
						on:click={handleAction}
						disabled={processing}
					>
						{processing ? 'Submitting...' : 'Confirm'}
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>
