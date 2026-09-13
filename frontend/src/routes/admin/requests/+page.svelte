<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi, hostingRequestsApi, extractError } from '$lib/api';
	import type { HostingRequest } from '$lib/types/hosting';
	import { Table, TableRow, TableCell, type TableColumn } from '$lib/components/ui';
	import {
		Globe,
		CloudArrowUp,
		Clock,
		Gear,
		CheckCircle,
		XCircle
	} from 'phosphor-svelte';
	import { domainSuffix } from '$lib/stores/node';

	let requests: HostingRequest[] = [];
	let loading = true;
	let activeFilter = 'ALL';
	let searchQuery = '';

	let activeRequest: any = null;
	let actionType: 'approve' | 'reject' | 'complete' | null = null;

	let adminNotes = '';
	let adminSubdomain = '';
	let subdomainChecking = false;
	let subdomainStatus: 'idle' | 'available' | 'taken' | 'invalid' = 'idle';
	let subdomainMessage = '';
	let checkTimer: any = null;

	let publicURL = '';
	let processing = false;
	let error: string | null = null;
	let success: string | null = null;

	async function loadRequests() {
		loading = true;
		try {
			const res = await adminApi.requests.listRequests();
			requests = res.requests || [];
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
			r.subdomain?.toLowerCase().includes(q) ||
			r.requester?.display_name?.toLowerCase().includes(q) ||
			r.requester?.username?.toLowerCase().includes(q);

		return matchesFilter && matchesSearch;
	});

	function handleSubdomainInput() {
		adminSubdomain = adminSubdomain.toLowerCase().replace(/[^a-z0-9-]/g, '');
		publicURL = `https://${adminSubdomain || 'app'}${$domainSuffix}`;
		if (!adminSubdomain) {
			subdomainStatus = 'idle';
			subdomainMessage = '';
			return;
		}

		clearTimeout(checkTimer);
		subdomainChecking = true;
		checkTimer = setTimeout(async () => {
			try {
				const res = await hostingRequestsApi.checkSubdomain(adminSubdomain);
				if (res.available) {
					subdomainStatus = 'available';
					subdomainMessage = res.message || 'Subdomain available';
				} else {
					subdomainStatus = 'taken';
					subdomainMessage = res.message || 'Subdomain is already in use';
				}
			} catch (err) {
				subdomainStatus = 'invalid';
				subdomainMessage = extractError(err);
			} finally {
				subdomainChecking = false;
			}
		}, 300);
	}

	function openModal(req: any, type: 'approve' | 'reject' | 'complete') {
		activeRequest = req;
		actionType = type;
		adminNotes = '';
		adminSubdomain = req.subdomain || req.project_name.toLowerCase().replace(/[^a-z0-9-]/g, '-').replace(/^-+|-+$/g, '');
		publicURL = req.public_url || `https://${adminSubdomain || 'app'}${$domainSuffix}`;
		subdomainStatus = 'idle';
		subdomainMessage = '';
		error = null;
		success = null;
	}

	function closeModal() {
		activeRequest = null;
		actionType = null;
		clearTimeout(checkTimer);
	}

	async function handleAction() {
		if (!activeRequest || !actionType) return;
		processing = true;
		error = null;

		try {
			if (actionType === 'approve') {
				await adminApi.requests.approveRequest(activeRequest.id, {
					admin_notes: adminNotes,
					subdomain: adminSubdomain.trim()
				});
				if (activeRequest.request_type === 'SUBDOMAIN_CHANGE') {
					success = `Subdomain change for "${activeRequest.project_name}" approved and updated to "${adminSubdomain}".`;
				} else {
					success = `Request for ${activeRequest.project_name} marked as APPROVED with subdomain "${adminSubdomain}".`;
				}
			} else if (actionType === 'reject') {
				await adminApi.requests.rejectRequest(activeRequest.id, {
					admin_notes: adminNotes
				});
				success = `Request for ${activeRequest.project_name} marked as REJECTED.`;
			} else if (actionType === 'complete') {
				await adminApi.requests.completeRequest(activeRequest.id, {
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
				<div class="flex items-center gap-2 flex-wrap">
					<span class="font-semibold text-sm sm:text-base text-(--text-main)">{req.project_name}</span>
					{#if req.request_type === 'SUBDOMAIN_CHANGE'}
						<span class="inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-[4px] bg-sky-500/10 text-sky-700 dark:text-sky-300 border border-sky-500/20 font-medium">
							<Globe size={12} weight="bold" />
							<span>Subdomain change</span>
						</span>
					{:else}
						<span class="inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-[4px] bg-(--bg-muted) text-(--text-secondary) border border-(--border-hairline) font-medium">
							<CloudArrowUp size={12} weight="bold" />
							<span>New hosting</span>
						</span>
					{/if}
				</div>
				{#if req.subdomain}
					<div class="text-xs font-mono text-(--accent-strong) mt-0.5">
						https://{req.subdomain}{$domainSuffix}
					</div>
				{/if}
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
					<span class="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-[4px] bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20">
						<Clock size={12} weight="bold" />
						<span>Pending</span>
					</span>
				{:else if req.status === 'APPROVED' || req.status === 'SETUP'}
					<span class="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-[4px] bg-sky-500/10 text-sky-700 dark:text-sky-400 border border-sky-500/20">
						<Gear size={12} weight="bold" />
						<span>In setup</span>
					</span>
				{:else if req.status === 'COMPLETED'}
					<span class="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-[4px] bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
						<CheckCircle size={12} weight="bold" />
						<span>Published</span>
					</span>
				{:else}
					<span class="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-[4px] bg-rose-500/10 text-rose-700 dark:text-rose-400 border border-rose-500/20">
						<XCircle size={12} weight="bold" />
						<span>Declined</span>
					</span>
				{/if}
			</TableCell>
			<TableCell alignTop mono class="text-xs text-(--text-muted)">
				{req.created_at ? new Date(req.created_at).toLocaleDateString([], { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) : '—'}
			</TableCell>
			<TableCell alignTop align="right" class="whitespace-nowrap">
				<div class="flex items-center justify-end gap-1.5">
					<a
						href="/admin/requests/{req.id}"
						class="btn btn-secondary btn-sm text-xs py-1 px-2.5 font-mono"
					>
						Review
					</a>

					{#if req.status === 'PENDING'}
						<button
							type="button"
							class="btn btn-primary btn-sm text-xs py-1 px-2 font-mono"
							on:click={() => openModal(req, 'approve')}
						>
							Approve
						</button>
						<button
							type="button"
							class="btn btn-secondary btn-sm text-xs py-1 px-2 font-mono text-(--color-danger)"
							on:click={() => openModal(req, 'reject')}
						>
							Reject
						</button>
					{:else if req.status === 'APPROVED' || req.status === 'SETUP'}
						<button
							type="button"
							class="btn btn-primary btn-sm text-xs py-1 px-2 font-mono"
							on:click={() => openModal(req, 'complete')}
						>
							Publish
						</button>
					{/if}
				</div>
			</TableCell>
		</TableRow>
	</Table>

	<!-- Operational Action Modal -->
	{#if activeRequest && actionType}
		<div class="fixed inset-0 bg-black/60 backdrop-blur-xs flex items-center justify-center z-50 p-4">
			<div class="w-full max-w-lg p-6 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-4 shadow-xl max-h-[90vh] overflow-y-auto">
				<div class="border-b border-(--border-hairline) pb-3">
					<div class="text-xs text-(--accent-strong) font-medium">Operator Action</div>
					<h2 class="font-display font-bold text-lg text-(--text-main) mt-0.5">
						{#if actionType === 'approve'}
							{#if activeRequest.request_type === 'SUBDOMAIN_CHANGE'}
								Approve Subdomain Change: {activeRequest.project_name}
							{:else}
								Approve Hosting: {activeRequest.project_name}
							{/if}
						{:else if actionType === 'reject'}
							Reject Request: {activeRequest.project_name}
						{:else}
							Publish Endpoint: {activeRequest.project_name}
						{/if}
					</h2>
				</div>

				{#if actionType === 'approve'}
					<div class="flex flex-col gap-1.5">
						<label for="admin-subdomain" class="text-sm font-medium text-(--text-secondary)">
							{activeRequest.request_type === 'SUBDOMAIN_CHANGE' ? 'Review / Override New Subdomain *' : 'Assign / Edit Subdomain *'}
						</label>
						<div class="flex items-stretch rounded-md border {subdomainStatus === 'taken' ? 'border-amber-500/50 bg-amber-500/5' : subdomainStatus === 'available' ? 'border-emerald-500/50 bg-emerald-500/5' : 'border-(--border-hairline) bg-(--bg-muted)'} transition-colors">
							<input
								id="admin-subdomain"
								type="text"
								class="flex-1 px-3.5 py-2 bg-transparent text-sm font-mono text-(--text-main) outline-none"
								bind:value={adminSubdomain}
								on:input={handleSubdomainInput}
								placeholder="subdomain"
							/>
							<div class="px-3 py-2 bg-(--bg-surface) border-l border-(--border-hairline) text-sm font-mono text-(--text-muted) flex items-center select-none">
								{$domainSuffix}
							</div>
						</div>
						{#if subdomainChecking}
							<span class="text-xs text-(--text-muted) font-mono animate-pulse">Checking availability...</span>
						{:else if subdomainStatus === 'available'}
							<span class="text-xs text-emerald-500 font-medium">✓ Subdomain is available for assignment.</span>
						{:else if subdomainStatus === 'taken'}
							<span class="text-xs text-amber-500 font-medium">⚠ {subdomainMessage}</span>
						{:else if subdomainStatus === 'invalid'}
							<span class="text-xs text-red-500 font-medium">{subdomainMessage}</span>
						{/if}
						<span class="text-xs text-(--text-muted)">
							Target domain: <span class="font-mono text-(--text-main)">https://{adminSubdomain || 'app'}{$domainSuffix}</span>
						</span>
					</div>

					{#if activeRequest.readme}
						<details class="text-xs border border-(--border-hairline) rounded p-2 bg-(--bg-muted)">
							<summary class="font-mono cursor-pointer text-(--accent-strong) font-medium">View Submitted README.md ({activeRequest.readme.length} chars)</summary>
							<pre class="mt-2 p-2 bg-(--bg-surface) rounded font-mono text-xs overflow-x-auto max-h-48 whitespace-pre-wrap">{activeRequest.readme}</pre>
						</details>
					{/if}
				{/if}

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
