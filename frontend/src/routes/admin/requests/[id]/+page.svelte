<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { adminApi, hostingRequestsApi, extractError } from '$lib/api';
	import { setAdminBreadcrumb } from '$lib/stores';
	import {
		ArrowUpRight,
		ArrowClockwise,
		CheckCircle,
		XCircle,
		Clock,
		Globe,
		GitBranch,
		User,
		FileText,
		Folder,
		Terminal,
		Check,
		WarningCircle,
		CloudArrowUp,
		Gear,
		Image
	} from 'phosphor-svelte';
	import { ProjectCover } from '$lib/components/ui';
	import { domainSuffix } from '$lib/stores/node';

	let requestId = $page.params.id || '';
	let reqData: any = null;
	let linkedProject: any = null;
	let loading = true;
	let error: string | null = null;
	let success: string | null = null;

	// Review form states
	let actionType: 'approve' | 'reject' | 'complete' | null = null;
	let adminNotes = '';
	let adminSubdomain = '';
	let subdomainChecking = false;
	let subdomainStatus: 'idle' | 'available' | 'taken' | 'invalid' = 'idle';
	let subdomainMessage = '';
	let checkTimer: any = null;
	let publicURL = '';
	let processing = false;

	async function loadRequest() {
		if (!requestId) {
			error = 'Invalid request ID';
			loading = false;
			return;
		}
		loading = true;
		error = null;
		try {
			const res = await adminApi.requests.getRequest(requestId);
			reqData = res.request;
			linkedProject = res.project;

			if (reqData) {
				setAdminBreadcrumb([
					{ label: 'Hosting Requests', href: '/admin/requests' },
					{ label: reqData.project_name }
				]);

				adminSubdomain = reqData.subdomain || reqData.project_name.toLowerCase().replace(/[^a-z0-9-]/g, '-').replace(/^-+|-+$/g, '');
				publicURL = reqData.public_url || `https://${adminSubdomain || 'app'}${$domainSuffix}`;
				adminNotes = reqData.admin_notes || '';
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadRequest);

	onDestroy(() => {
		setAdminBreadcrumb(null);
		clearTimeout(checkTimer);
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
					subdomainMessage = res.message || 'Subdomain is already assigned';
				}
			} catch (err) {
				subdomainStatus = 'invalid';
				subdomainMessage = extractError(err);
			} finally {
				subdomainChecking = false;
			}
		}, 300);
	}

	async function submitAction(type: 'approve' | 'reject' | 'complete') {
		if (!reqData) return;
		processing = true;
		error = null;
		success = null;

		try {
			if (type === 'approve') {
				await adminApi.requests.approveRequest(reqData.id, {
					admin_notes: adminNotes,
					subdomain: adminSubdomain.trim()
				});
				success = `Request approved with subdomain "${adminSubdomain}".`;
			} else if (type === 'reject') {
				await adminApi.requests.rejectRequest(reqData.id, {
					admin_notes: adminNotes
				});
				success = `Request marked as rejected.`;
			} else if (type === 'complete') {
				await adminApi.requests.completeRequest(reqData.id, {
					admin_notes: adminNotes,
					public_url: publicURL
				});
				success = `Project deployed and published at ${publicURL}.`;
			}
			actionType = null;
			await loadRequest();
		} catch (err) {
			error = extractError(err);
		} finally {
			processing = false;
		}
	}

	function formatDate(iso?: string) {
		if (!iso) return '—';
		try {
			const d = new Date(iso);
			return d.toLocaleDateString('en-US', {
				day: 'numeric',
				month: 'short',
				year: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		} catch {
			return iso;
		}
	}
</script>

<svelte:head>
	<title>{reqData?.project_name ? `${reqData.project_name} · Hosting Request Review` : 'Hosting Request Review · Admin Console'}</title>
</svelte:head>

<div class="flex flex-col gap-8">
	{#if loading}
		<div class="py-24 text-center text-sm text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-6 h-6 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading hosting request details...</span>
		</div>
	{:else if error && !reqData}
		<div class="p-4 rounded-md bg-rose-500/10 text-rose-700 dark:text-rose-400 text-sm border border-rose-500/20">
			{error}
		</div>
	{:else if reqData}
		{#if success}
			<div class="p-3.5 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 text-sm rounded-md border border-emerald-500/20 flex items-center justify-between">
				<span class="font-medium">{success}</span>
				<button type="button" on:click={() => (success = null)} class="font-bold text-base cursor-pointer px-1">×</button>
			</div>
		{/if}

		{#if error}
			<div class="p-3.5 bg-rose-500/10 text-rose-700 dark:text-rose-400 text-sm rounded-md border border-rose-500/20 flex items-center justify-between">
				<span>{error}</span>
				<button type="button" on:click={() => (error = null)} class="font-bold text-base cursor-pointer px-1">×</button>
			</div>
		{/if}

		<!-- ══════════════════════════════════════════════
		     1. Header & Overview
		     ══════════════════════════════════════════════ -->
		<header class="flex flex-col lg:flex-row lg:items-start justify-between gap-6 pb-6 border-b border-(--border-hairline)">
			<div class="flex flex-col gap-2 min-w-0">
				<div class="flex items-center gap-3 flex-wrap">
					<h1 class="font-display font-bold text-2xl sm:text-3xl text-(--text-main) tracking-tight">
						{reqData.project_name}
					</h1>

					<!-- Request Type Chip -->
					{#if reqData.request_type === 'SUBDOMAIN_CHANGE'}
						<span class="inline-flex items-center gap-1.5 text-xs px-2.5 py-1 rounded-[4px] bg-sky-500/10 text-sky-700 dark:text-sky-300 border border-sky-500/20 font-medium">
							<Globe size={13} weight="bold" />
							<span>Subdomain change</span>
						</span>
					{:else}
						<span class="inline-flex items-center gap-1.5 text-xs px-2.5 py-1 rounded-[4px] bg-(--bg-muted) text-(--text-secondary) border border-(--border-hairline) font-medium">
							<CloudArrowUp size={13} weight="bold" />
							<span>New hosting</span>
						</span>
					{/if}

					<!-- Status Chip -->
					{#if reqData.status === 'PENDING'}
						<span class="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-[4px] bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20">
							<Clock size={13} weight="bold" />
							<span>Pending review</span>
						</span>
					{:else if reqData.status === 'APPROVED' || reqData.status === 'SETUP'}
						<span class="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-[4px] bg-sky-500/10 text-sky-700 dark:text-sky-400 border border-sky-500/20">
							<Gear size={13} weight="bold" />
							<span>Setup in progress</span>
						</span>
					{:else if reqData.status === 'COMPLETED'}
						<span class="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-[4px] bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
							<CheckCircle size={13} weight="bold" />
							<span>Published</span>
						</span>
					{:else}
						<span class="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-[4px] bg-rose-500/10 text-rose-700 dark:text-rose-400 border border-rose-500/20">
							<XCircle size={13} weight="bold" />
							<span>Declined</span>
						</span>
					{/if}
				</div>

				<p class="text-sm text-(--text-secondary) max-w-3xl leading-relaxed">
					{reqData.description || 'No description provided.'}
				</p>

				<div class="flex items-center gap-4 text-xs font-mono text-(--text-muted) pt-1 flex-wrap">
					<span>Target: <strong class="text-(--text-main)">https://{reqData.subdomain || 'app'}{$domainSuffix}</strong></span>
					<span>•</span>
					<span>Submitted: {formatDate(reqData.created_at)}</span>
				</div>
			</div>

			<div class="flex items-center gap-2 shrink-0">
				<a
					href="/admin/requests"
					class="px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-surface) hover:bg-(--bg-muted) text-(--text-secondary) hover:text-(--text-main) text-sm font-medium transition-colors"
				>
					← Back to queue
				</a>

				<button
					type="button"
					class="p-2 rounded-sm border border-(--border-hairline) bg-(--bg-surface) hover:bg-(--bg-muted) text-(--text-muted) hover:text-(--text-main) transition-colors cursor-pointer flex items-center justify-center"
					on:click={loadRequest}
					title="Refresh"
					aria-label="Refresh"
				>
					<ArrowClockwise size={16} weight="regular" />
				</button>
			</div>
		</header>

		<!-- ══════════════════════════════════════════════
		     2. Main Body Grid: Specifications & Review
		     ══════════════════════════════════════════════ -->
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
			<!-- Technical Specifications & Details (2 Cols) -->
			<div class="lg:col-span-2 flex flex-col gap-6">
				<!-- Cover Photo Preview (if submitted) -->
				{#if reqData.cover_image_url}
					<div class="p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3">
						<div class="flex items-center gap-2">
							<Image size={18} weight="regular" class="text-(--text-muted)" />
							<h2 class="font-sans font-semibold text-base text-(--text-main)">Cover Photo</h2>
						</div>
						<div class="w-full rounded-md overflow-hidden border border-(--border-hairline) bg-(--bg-muted)">
							<ProjectCover src={reqData.cover_image_url} alt={reqData.project_name} name={reqData.project_name} aspectRatio="16/9" />
						</div>
					</div>
				{/if}

				<!-- Core Technical Metadata -->
				<div class="p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-4">
					<h2 class="font-sans font-semibold text-base text-(--text-main)">Technical Specifications</h2>

					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
						<div class="flex flex-col gap-1 p-3 rounded-sm bg-(--bg-muted)/40 border border-(--border-hairline)">
							<span class="text-xs text-(--text-muted)">Repository URL</span>
							<a
								href={reqData.repository_url}
								target="_blank"
								rel="noreferrer"
								class="text-(--accent-sky) hover:underline font-mono text-xs truncate flex items-center gap-1 mt-0.5"
							>
								<span>{reqData.repository_url}</span>
								<ArrowUpRight size={12} weight="bold" />
							</a>
						</div>

						{#if reqData.documentation_url}
							<div class="flex flex-col gap-1 p-3 rounded-sm bg-(--bg-muted)/40 border border-(--border-hairline)">
								<span class="text-xs text-(--text-muted)">Documentation</span>
								<a
									href={reqData.documentation_url}
									target="_blank"
									rel="noreferrer"
									class="text-(--accent-sky) hover:underline font-mono text-xs truncate flex items-center gap-1 mt-0.5"
								>
									<span>{reqData.documentation_url}</span>
									<ArrowUpRight size={12} weight="bold" />
								</a>
							</div>
						{/if}

						<div class="flex flex-col gap-1 p-3 rounded-sm bg-(--bg-muted)/40 border border-(--border-hairline)">
							<span class="text-xs text-(--text-muted)">Requested Subdomain</span>
							<span class="font-mono text-xs text-(--text-main) mt-0.5">
								{reqData.subdomain || '—'}{$domainSuffix}
							</span>
						</div>

						<div class="flex flex-col gap-1 p-3 rounded-sm bg-(--bg-muted)/40 border border-(--border-hairline)">
							<span class="text-xs text-(--text-muted)">Submission ID</span>
							<span class="font-mono text-xs text-(--text-muted) mt-0.5">{reqData.id}</span>
						</div>
					</div>

					<!-- Technology Stack -->
					{#if reqData.technology_stack && reqData.technology_stack.length > 0}
						<div class="pt-2">
							<span class="text-xs text-(--text-muted) block mb-2">Declared Stack</span>
							<div class="flex flex-wrap gap-1.5">
								{#each reqData.technology_stack as tech}
									<span class="text-xs font-mono px-2.5 py-1 rounded-sm bg-(--bg-muted) border border-(--border-hairline) text-(--text-secondary)">
										{tech}
									</span>
								{/each}
							</div>
						</div>
					{/if}
				</div>

				<!-- Deployment Notes from Requester -->
				{#if reqData.deployment_notes}
					<div class="p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3">
						<div class="flex items-center gap-2">
							<Terminal size={18} weight="regular" class="text-(--text-muted)" />
							<h2 class="font-sans font-semibold text-base text-(--text-main)">Deployment & Environment Notes</h2>
						</div>
						<pre class="p-4 rounded-sm bg-(--bg-muted)/60 border border-(--border-hairline) text-xs font-mono text-(--text-main) whitespace-pre-wrap overflow-x-auto leading-relaxed">{reqData.deployment_notes}</pre>
					</div>
				{/if}

				<!-- Submitted README.md -->
				{#if reqData.readme}
					<div class="p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3">
						<div class="flex items-center gap-2">
							<FileText size={18} weight="regular" class="text-(--text-muted)" />
							<h2 class="font-sans font-semibold text-base text-(--text-main)">Submitted README</h2>
						</div>
						<pre class="p-4 rounded-sm bg-(--bg-muted)/60 border border-(--border-hairline) text-xs font-mono text-(--text-main) whitespace-pre-wrap overflow-x-auto max-h-96 leading-relaxed">{reqData.readme}</pre>
					</div>
				{/if}

				<!-- Linked Existing Project (if attached) -->
				{#if linkedProject}
					<div class="p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-4">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2">
								<Folder size={18} weight="regular" class="text-(--text-muted)" />
								<h2 class="font-sans font-semibold text-base text-(--text-main)">Linked Application</h2>
							</div>
							<a
								href="/admin/projects/{linkedProject.id}"
								class="text-xs text-(--accent-sky) hover:underline flex items-center gap-1 font-medium"
							>
								<span>Manage Project</span>
								<ArrowUpRight size={12} weight="bold" />
							</a>
						</div>

						<div class="p-4 rounded-sm bg-(--bg-muted)/40 border border-(--border-hairline) flex items-center justify-between gap-4 flex-wrap text-sm">
							<div>
								<span class="font-semibold text-(--text-main)">{linkedProject.name}</span>
								<span class="text-xs font-mono text-(--text-muted) block">/{linkedProject.slug}</span>
							</div>

							<div class="flex items-center gap-2">
								<span class="text-xs px-2 py-0.5 rounded-sm font-mono border bg-(--bg-surface) text-(--text-muted)">
									{linkedProject.lifecycle_status || 'ACTIVE'}
								</span>
								{#if linkedProject.public_url}
									<a
										href={linkedProject.public_url}
										target="_blank"
										rel="noreferrer"
										class="text-xs text-(--accent-sky) hover:underline font-mono"
									>
										{linkedProject.public_url} ↗
									</a>
								{/if}
							</div>
						</div>
					</div>
				{/if}
			</div>

			<!-- Requester Profile & Review Actions (1 Col, Sticky on Desktop) -->
			<div class="lg:col-span-1 lg:sticky lg:top-20 self-start flex flex-col gap-6">
				<!-- Requester Profile Card -->
				<div class="p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-4">
					<h2 class="font-sans font-semibold text-base text-(--text-main)">Requester Profile</h2>

					<div class="flex items-center gap-3">
						{#if reqData.requester?.avatar_url}
							<img src={reqData.requester.avatar_url} alt="" class="w-12 h-12 rounded-full object-cover shrink-0" />
						{:else}
							<div class="w-12 h-12 rounded-full bg-(--accent-sky)/10 text-(--accent-sky) flex items-center justify-center text-base font-bold shrink-0">
								{reqData.requester?.display_name?.slice(0, 2).toUpperCase() || 'RQ'}
							</div>
						{/if}

						<div class="min-w-0 flex-1">
							<div class="font-medium text-sm text-(--text-main) truncate">
								{reqData.requester?.display_name || 'Member'}
							</div>
							<div class="text-xs font-mono text-(--text-muted) truncate">
								@{reqData.requester?.username || 'user'}
							</div>
							{#if reqData.requester?.role}
								<span class="inline-block text-xs font-mono px-2 py-0.5 rounded-sm bg-(--bg-muted) dark:bg-neutral-900 text-(--text-muted) mt-1">
									{reqData.requester.role}
								</span>
							{/if}
						</div>
					</div>

					<div class="pt-2 border-t border-(--border-hairline) text-xs text-(--text-muted) flex flex-col gap-1.5 font-mono">
						<span>Account ID: {reqData.requester_id}</span>
						{#if reqData.requester?.created_at}
							<span>Joined: {formatDate(reqData.requester.created_at)}</span>
						{/if}
					</div>
				</div>

				<!-- Operator Review Action Center -->
				<div class="p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-5">
					<div>
						<h2 class="font-sans font-semibold text-base text-(--text-main)">Operator Decision</h2>
						<p class="text-xs text-(--text-muted) mt-0.5">Review workflow and proxy routing assignment</p>
					</div>

					<!-- Existing Review Metadata -->
					{#if reqData.reviewed_at}
						<div class="p-3.5 rounded-sm bg-(--bg-muted)/40 border border-(--border-hairline) text-xs flex flex-col gap-1">
							<span class="text-(--text-muted)">Reviewed on {formatDate(reqData.reviewed_at)}</span>
							{#if reqData.admin_notes}
								<p class="text-(--text-main) font-mono text-xs mt-1 bg-(--bg-surface) p-2 rounded-xs border border-(--border-hairline)">
									{reqData.admin_notes}
								</p>
							{/if}
						</div>
					{/if}

					<!-- Action Controls -->
					{#if reqData.status === 'PENDING'}
						<div class="flex flex-col gap-4">
							<div class="flex items-center gap-2">
								<button
									type="button"
									class="flex-1 py-2 rounded-sm text-sm font-medium transition-colors cursor-pointer {actionType === 'approve' ? 'bg-emerald-600 text-white font-semibold' : 'border border-(--border-hairline) bg-(--bg-surface) hover:bg-(--bg-muted) text-(--text-main)'}"
									on:click={() => (actionType = 'approve')}
								>
									Approve
								</button>
								<button
									type="button"
									class="flex-1 py-2 rounded-sm text-sm font-medium transition-colors cursor-pointer {actionType === 'reject' ? 'bg-rose-600 text-white font-semibold' : 'border border-(--border-hairline) bg-(--bg-surface) hover:bg-rose-500/10 text-rose-600 dark:text-rose-400'}"
									on:click={() => (actionType = 'reject')}
								>
									Reject
								</button>
							</div>

							{#if actionType === 'approve'}
								<div class="flex flex-col gap-3 pt-2 border-t border-(--border-hairline)">
									<div class="flex flex-col gap-1.5">
										<label for="subdomain-input" class="text-xs font-medium text-(--text-secondary)">
											Assign Subdomain *
										</label>
										<div class="flex items-stretch rounded-sm border {subdomainStatus === 'taken' ? 'border-amber-500/50 bg-amber-500/5' : subdomainStatus === 'available' ? 'border-emerald-500/50 bg-emerald-500/5' : 'border-(--border-hairline) bg-(--bg-muted)'}">
											<input
												id="subdomain-input"
												type="text"
												class="flex-1 px-3 py-1.5 bg-transparent text-xs font-mono text-(--text-main) outline-none"
												bind:value={adminSubdomain}
												on:input={handleSubdomainInput}
												placeholder="subdomain"
											/>
											<div class="px-2.5 py-1.5 bg-(--bg-surface) border-l border-(--border-hairline) text-xs font-mono text-(--text-muted) flex items-center select-none">
												{$domainSuffix}
											</div>
										</div>
										{#if subdomainChecking}
											<span class="text-xs text-(--text-muted) font-mono">Checking availability...</span>
										{:else if subdomainStatus === 'available'}
											<span class="text-xs text-emerald-600 dark:text-emerald-400 font-medium">✓ Subdomain is available</span>
										{:else if subdomainStatus === 'taken'}
											<span class="text-xs text-amber-600 dark:text-amber-400 font-medium">⚠ {subdomainMessage}</span>
										{/if}
									</div>

									<div class="flex flex-col gap-1.5">
										<label for="notes-approve" class="text-xs font-medium text-(--text-secondary)">
											Operator Notes (communicated to user)
										</label>
										<textarea
											id="notes-approve"
											rows="3"
											class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-sm text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
											placeholder="Optional instructions for deployment..."
											bind:value={adminNotes}
										></textarea>
									</div>

									<button
										type="button"
										class="w-full py-2.5 rounded-sm bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold transition-colors cursor-pointer disabled:opacity-50"
										on:click={() => submitAction('approve')}
										disabled={processing || !adminSubdomain.trim() || subdomainStatus === 'taken'}
									>
										{processing ? 'Submitting...' : 'Confirm Approval'}
									</button>
								</div>
							{:else if actionType === 'reject'}
								<div class="flex flex-col gap-3 pt-2 border-t border-(--border-hairline)">
									<div class="flex flex-col gap-1.5">
										<label for="notes-reject" class="text-xs font-medium text-(--text-secondary)">
											Rejection Reason *
										</label>
										<textarea
											id="notes-reject"
											rows="3"
											class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-sm text-xs text-(--text-main) outline-none focus:border-rose-500"
											placeholder="Explain why this request was declined..."
											bind:value={adminNotes}
										></textarea>
									</div>

									<button
										type="button"
										class="w-full py-2.5 rounded-sm bg-rose-600 hover:bg-rose-700 text-white text-xs font-semibold transition-colors cursor-pointer disabled:opacity-50"
										on:click={() => submitAction('reject')}
										disabled={processing}
									>
										{processing ? 'Rejecting...' : 'Confirm Rejection'}
									</button>
								</div>
							{/if}
						</div>
					{:else if reqData.status === 'APPROVED' || reqData.status === 'SETUP'}
						<div class="flex flex-col gap-4">
							<div class="p-3 rounded-sm bg-sky-500/10 border border-sky-500/20 text-sky-700 dark:text-sky-400 text-xs">
								Request is approved. Once you deploy the container, publish its public proxy URL here.
							</div>

							<div class="flex flex-col gap-3">
								<div class="flex flex-col gap-1.5">
									<label for="public-url-input" class="text-xs font-medium text-(--text-secondary)">
										Live Public URL *
									</label>
									<input
										id="public-url-input"
										type="url"
										class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-sm text-xs font-mono text-(--text-main) outline-none focus:border-(--accent-sky)"
										bind:value={publicURL}
										placeholder="https://app{$domainSuffix}"
									/>
								</div>

								<div class="flex flex-col gap-1.5">
									<label for="notes-publish" class="text-xs font-medium text-(--text-secondary)">
										Operator Notes
									</label>
									<textarea
										id="notes-publish"
										rows="2"
										class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-sm text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
										placeholder="Optional notes..."
										bind:value={adminNotes}
									></textarea>
								</div>

								<button
									type="button"
									class="w-full py-2.5 rounded-sm bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold transition-colors cursor-pointer disabled:opacity-50 flex items-center justify-center gap-1.5"
									on:click={() => submitAction('complete')}
									disabled={processing || !publicURL.trim()}
								>
									<Check size={16} weight="bold" />
									<span>{processing ? 'Publishing...' : 'Publish & Complete Request'}</span>
								</button>
							</div>
						</div>
					{:else if reqData.status === 'COMPLETED'}
						<div class="p-3.5 rounded-sm bg-emerald-500/10 border border-emerald-500/20 text-emerald-700 dark:text-emerald-400 text-xs flex flex-col gap-1.5">
							<span class="font-semibold">Project is Live</span>
							<p>This request has been completed and provisioned.</p>
							{#if linkedProject}
								<a href="/admin/projects/{linkedProject.id}" class="text-emerald-800 dark:text-emerald-300 underline font-medium mt-1">
									View project in admin console →
								</a>
							{/if}
						</div>
					{:else}
						<div class="p-3.5 rounded-sm bg-rose-500/10 border border-rose-500/20 text-rose-700 dark:text-rose-400 text-xs">
							This request was rejected.
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>
