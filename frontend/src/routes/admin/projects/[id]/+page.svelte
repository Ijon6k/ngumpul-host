<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { adminApi, commentsApi, extractError } from '$lib/api';
	import { setAdminBreadcrumb } from '$lib/stores';
	import { resolveProjectStatus } from '$lib/utils/projectStatus';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import UptimeHistory from '$lib/components/status/UptimeHistory.svelte';
	import { ConfirmModal } from '$lib/components/ui';
	import {
		ArrowUpRight,
		Trash,
		ShieldWarning,
		ArrowClockwise,
		CheckCircle,
		XCircle,
		Clock,
		Globe,
		GitBranch,
		User,
		WarningCircle,
		Archive,
		StopCircle,
		PlayCircle,
		Eye,
		EyeSlash
	} from 'phosphor-svelte';

	let projectId = $page.params.id || '';
	let project: any = null;
	let availabilityData: any = null;
	let activities: any[] = [];
	let comments: any[] = [];
	let loading = true;
	let error: string | null = null;
	let success: string | null = null;

	// Modal states
	let deleteModalOpen = false;
	let commentToDelete: any = null;
	let deletingComment = false;

	let suspendModalOpen = false;
	let archiveModalOpen = false;
	let actionLoading = false;

	async function loadData() {
		if (!projectId) {
			error = 'Invalid project ID';
			loading = false;
			return;
		}
		loading = true;
		error = null;
		try {
			const res = await adminApi.projects.getProject(projectId);
			project = res.project;
			availabilityData = res.availability;
			activities = res.activities || [];

			if (project) {
				setAdminBreadcrumb([
					{ label: 'Projects', href: '/admin/projects' },
					{ label: project.name }
				]);

				const commRes = await commentsApi.getComments(project.slug);
				comments = commRes.comments || [];
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadData);

	onDestroy(() => {
		setAdminBreadcrumb(null);
	});

	// Reactive resolved status
	$: resolvedStatus = resolveProjectStatus(project);

	async function handleSuspend() {
		if (!project) return;
		actionLoading = true;
		error = null;
		try {
			await adminApi.projects.suspendProject(project.id);
			success = `Project "${project.name}" has been administratively suspended.`;
			suspendModalOpen = false;
			await loadData();
		} catch (err) {
			error = extractError(err);
		} finally {
			actionLoading = false;
		}
	}

	async function handleRestore() {
		if (!project) return;
		actionLoading = true;
		error = null;
		try {
			await adminApi.projects.restoreProject(project.id);
			success = `Project "${project.name}" has been restored to active status.`;
			await loadData();
		} catch (err) {
			error = extractError(err);
		} finally {
			actionLoading = false;
		}
	}

	async function handleActivate() {
		if (!project) return;
		actionLoading = true;
		error = null;
		try {
			await adminApi.projects.activateProject(project.id);
			success = `Project "${project.name}" setup is complete and marked as active.`;
			await loadData();
		} catch (err) {
			error = extractError(err);
		} finally {
			actionLoading = false;
		}
	}

	async function handleArchive() {
		if (!project) return;
		actionLoading = true;
		error = null;
		try {
			await adminApi.projects.archiveProject(project.id);
			success = `Project "${project.name}" has been archived and unpublished.`;
			archiveModalOpen = false;
			await loadData();
		} catch (err) {
			error = extractError(err);
		} finally {
			actionLoading = false;
		}
	}

	async function toggleVisibility() {
		if (!project) return;
		actionLoading = true;
		error = null;
		const newVis = project.visibility === 'PUBLIC' ? 'UNPUBLISHED' : 'PUBLIC';
		try {
			await adminApi.projects.updateProject(project.id, { visibility: newVis });
			project.visibility = newVis;
			success = `Visibility updated to ${newVis}.`;
		} catch (err) {
			error = extractError(err);
		} finally {
			actionLoading = false;
		}
	}

	function openDeleteCommentModal(c: any) {
		commentToDelete = c;
		deleteModalOpen = true;
	}

	async function handleDeleteComment() {
		if (!commentToDelete) return;
		deletingComment = true;
		try {
			await adminApi.comments.deleteComment(commentToDelete.id);
			comments = comments.map((c) =>
				c.id === commentToDelete.id
					? { ...c, is_deleted: true, content: 'This comment was removed.' }
					: c
			);
			deleteModalOpen = false;
			commentToDelete = null;
			success = 'Comment removed successfully.';
		} catch (err) {
			error = extractError(err);
		} finally {
			deletingComment = false;
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

	function getDiagnosticExplanation(reason?: string) {
		switch (reason) {
			case 'DNS_ERROR':
				return 'Hostname resolution failed. Ensure DNS records point to this host node.';
			case 'CONNECTION_FAILED':
				return 'TCP connection refused or dropped by the target service container.';
			case 'TIMEOUT':
				return 'Probe request timed out (>5000ms). Container may be overloaded or hung.';
			case 'TLS_ERROR':
				return 'SSL/TLS handshake failed or certificate validation error.';
			case 'HTTP_ERROR':
				return 'Endpoint returned an HTTP 5xx server error response.';
			default:
				return 'Endpoint probe failed to verify responsiveness.';
		}
	}
</script>

<svelte:head>
	<title>{project?.name ? `${project.name} · Admin Console` : 'Project Management · Admin Console'}</title>
</svelte:head>

<div class="flex flex-col gap-8">
	{#if loading}
		<div class="py-24 text-center text-sm text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-6 h-6 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading project details...</span>
		</div>
	{:else if error && !project}
		<div class="p-4 rounded-md bg-rose-500/10 text-rose-700 dark:text-rose-400 text-sm border border-rose-500/20">
			{error}
		</div>
	{:else if project}
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
		     1. Project Header & Primary Actions
		     ══════════════════════════════════════════════ -->
		<header class="flex flex-col lg:flex-row lg:items-start justify-between gap-6 pb-6 border-b border-(--border-hairline)">
			<div class="flex flex-col gap-2 min-w-0">
				<div class="flex items-center gap-3 flex-wrap">
					<h1 class="font-display font-bold text-2xl sm:text-3xl text-(--text-main) tracking-tight">
						{project.name}
					</h1>
					<!-- Status Badge -->
					<div class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium border {resolvedStatus.bgClass} {resolvedStatus.colorClass} {resolvedStatus.borderClass}">
						<StatusDot status={resolvedStatus.dotStatus} variant="dot" />
						<span>{resolvedStatus.label}</span>
					</div>
					<span class="text-sm font-mono text-(--text-muted)">/{project.slug}</span>
				</div>

				<p class="text-sm text-(--text-secondary) max-w-3xl leading-relaxed">
					{project.description || 'No project description provided.'}
				</p>

				{#if project.public_url}
					<div class="flex items-center gap-2 text-sm text-(--text-muted) pt-1">
						<Globe size={16} weight="regular" class="shrink-0" />
						<a
							href={project.public_url}
							target="_blank"
							rel="noreferrer"
							class="text-(--accent-sky) hover:underline font-mono text-xs sm:text-sm truncate"
						>
							{project.public_url}
						</a>
						<ArrowUpRight size={14} weight="bold" class="text-(--text-muted) shrink-0" />
					</div>
				{/if}
			</div>

			<!-- Operator Lifecycle Controls -->
			<div class="flex items-center gap-2 flex-wrap shrink-0">
				<!-- Visibility Button -->
				<button
					type="button"
					class="px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-surface) hover:bg-(--bg-muted) text-(--text-main) text-sm font-medium transition-colors cursor-pointer flex items-center gap-1.5"
					on:click={toggleVisibility}
					disabled={actionLoading}
					title="Toggle public showcase visibility"
				>
					{#if project.visibility === 'PUBLIC'}
						<Eye size={16} weight="regular" />
						<span>Public</span>
					{:else}
						<EyeSlash size={16} weight="regular" class="text-amber-500" />
						<span>Unpublished</span>
					{/if}
				</button>

				<!-- Lifecycle Actions -->
				{#if project.lifecycle_status === 'ACTIVE'}
					<button
						type="button"
						class="px-3.5 py-2 rounded-sm border border-amber-500/30 bg-amber-500/10 hover:bg-amber-500/20 text-amber-700 dark:text-amber-400 text-sm font-medium transition-colors cursor-pointer flex items-center gap-1.5"
						on:click={() => (suspendModalOpen = true)}
						disabled={actionLoading}
					>
						<StopCircle size={16} weight="regular" />
						<span>Suspend project</span>
					</button>
				{:else if project.lifecycle_status === 'SUSPENDED'}
					<button
						type="button"
						class="px-3.5 py-2 rounded-sm border border-emerald-500/30 bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-700 dark:text-emerald-400 text-sm font-medium transition-colors cursor-pointer flex items-center gap-1.5"
						on:click={handleRestore}
						disabled={actionLoading}
					>
						<PlayCircle size={16} weight="regular" />
						<span>Restore project</span>
					</button>
				{:else if project.lifecycle_status === 'SETUP'}
					<button
						type="button"
						class="px-3.5 py-2 rounded-sm border border-sky-500/30 bg-sky-500/10 hover:bg-sky-500/20 text-sky-700 dark:text-sky-400 text-sm font-medium transition-colors cursor-pointer flex items-center gap-1.5"
						on:click={handleActivate}
						disabled={actionLoading}
					>
						<CheckCircle size={16} weight="regular" />
						<span>Activate project</span>
					</button>
				{/if}

				{#if project.lifecycle_status !== 'ARCHIVED'}
					<button
						type="button"
						class="px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-surface) hover:bg-rose-500/10 hover:text-rose-600 dark:hover:text-rose-400 text-(--text-secondary) text-sm font-medium transition-colors cursor-pointer flex items-center gap-1.5"
						on:click={() => (archiveModalOpen = true)}
						disabled={actionLoading}
						title="Archive project and retire from service"
					>
						<Archive size={16} weight="regular" />
						<span>Archive</span>
					</button>
				{/if}

				<a
					href="/projects/{project.slug}"
					target="_blank"
					rel="noreferrer"
					class="px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-surface) hover:bg-(--bg-muted) text-(--text-secondary) hover:text-(--text-main) text-sm font-medium transition-colors inline-flex items-center gap-1.5"
					title="View public showcase page"
				>
					<span>Showcase</span>
					<ArrowUpRight size={14} weight="bold" />
				</a>

				<button
					type="button"
					class="p-2 rounded-sm border border-(--border-hairline) bg-(--bg-surface) hover:bg-(--bg-muted) text-(--text-muted) hover:text-(--text-main) transition-colors cursor-pointer flex items-center justify-center"
					on:click={loadData}
					title="Refresh project details"
					aria-label="Refresh"
				>
					<ArrowClockwise size={16} weight="regular" />
				</button>
			</div>
		</header>

		<!-- ══════════════════════════════════════════════
		     2. Operational Telemetry & State Architecture
		     ══════════════════════════════════════════════ -->
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- Telemetry & Diagnostic Status Card (2 cols) -->
			<div class="lg:col-span-2 flex flex-col gap-5 p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline)">
				<div class="flex items-center justify-between">
					<div>
						<h2 class="font-sans font-semibold text-base text-(--text-main)">Operational Diagnostics</h2>
						<p class="text-xs text-(--text-muted) mt-0.5">Real-time telemetry and health monitoring status</p>
					</div>

					<div class="text-right">
						<span class="text-xs text-(--text-muted) block">Lifecycle</span>
						<span class="font-mono text-xs font-semibold text-(--text-main)">{project.lifecycle_status || 'ACTIVE'}</span>
					</div>
				</div>

				<!-- Diagnostic Alert (if unreachable) -->
				{#if project.availability === 'UNREACHABLE'}
					<div class="p-4 rounded-sm bg-rose-500/10 border border-rose-500/20 text-rose-700 dark:text-rose-400 flex items-start gap-3">
						<WarningCircle size={20} weight="bold" class="shrink-0 mt-0.5" />
						<div class="flex flex-col gap-1 text-xs sm:text-sm leading-relaxed">
							<div class="font-semibold flex items-center gap-2">
								<span>Endpoint Unreachable</span>
								{#if project.availability_reason}
									<span class="font-mono px-2 py-0.5 rounded-full bg-rose-500/20 text-xs font-bold">
										{project.availability_reason}
									</span>
								{/if}
							</div>
							<p>{getDiagnosticExplanation(project.availability_reason)}</p>
						</div>
					</div>
				{:else if project.lifecycle_status === 'SUSPENDED'}
					<div class="p-4 rounded-sm bg-amber-500/10 border border-amber-500/20 text-amber-700 dark:text-amber-400 flex items-start gap-3">
						<StopCircle size={20} weight="bold" class="shrink-0 mt-0.5" />
						<div class="flex flex-col gap-1 text-xs sm:text-sm leading-relaxed">
							<span class="font-semibold">Administratively Suspended</span>
							<p>This application is temporarily taken offline by operator mandate. Telemetry probes are paused.</p>
						</div>
					</div>
				{/if}

				<!-- Telemetry Metrics Grid -->
				<div class="grid grid-cols-2 sm:grid-cols-4 gap-4 p-4 rounded-sm bg-(--bg-muted)/40 border border-(--border-hairline)">
					<div>
						<span class="text-xs text-(--text-muted) block">Availability Probe</span>
						<span class="font-mono font-medium text-sm text-(--text-main) mt-1 block">
							{project.availability || 'UNKNOWN'}
						</span>
					</div>

					<div>
						<span class="text-xs text-(--text-muted) block">Response Latency</span>
						<span class="font-mono font-medium text-sm text-(--text-main) mt-1 block">
							{availabilityData?.latest_latency_ms != null && availabilityData.latest_latency_ms > 0
								? `${availabilityData.latest_latency_ms}ms`
								: '—'}
						</span>
					</div>

					<div>
						<span class="text-xs text-(--text-muted) block">30-Day Uptime</span>
						<span class="font-mono font-medium text-sm text-(--text-main) mt-1 block">
							{availabilityData?.uptime_percent != null
								? `${availabilityData.uptime_percent.toFixed(1)}%`
								: '—'}
						</span>
					</div>

					<div>
						<span class="text-xs text-(--text-muted) block">Total Probes</span>
						<span class="font-mono font-medium text-sm text-(--text-main) mt-1 block">
							{availabilityData?.total_checks ?? 0}
						</span>
					</div>
				</div>

				<!-- Availability History Graph -->
				<div class="pt-2">
					<UptimeHistory
						availability={availabilityData}
						status={project.availability === 'REACHABLE' ? 'ONLINE' : 'OFFLINE'}
						variant="detailed"
						showRangeSelector={true}
					/>
				</div>
			</div>

			<!-- Safe Hosting Configuration & Owner Metadata (1 col) -->
			<div class="flex flex-col gap-5 p-6 rounded-md bg-(--bg-surface) border border-(--border-hairline)">
				<div>
					<h2 class="font-sans font-semibold text-base text-(--text-main)">Hosting Specifications</h2>
					<p class="text-xs text-(--text-muted) mt-0.5">Deployment architecture and ownership details</p>
				</div>

				<!-- Owner Profile -->
				<div class="p-3.5 rounded-sm bg-(--bg-muted)/40 border border-(--border-hairline) flex items-center gap-3">
					{#if project.owner?.avatar_url}
						<img src={project.owner.avatar_url} alt="" class="w-10 h-10 rounded-full object-cover shrink-0" />
					{:else}
						<div class="w-10 h-10 rounded-full bg-(--accent-sky)/10 text-(--accent-sky) flex items-center justify-center text-sm font-bold shrink-0">
							{project.owner?.display_name?.slice(0, 2).toUpperCase() || 'US'}
						</div>
					{/if}

					<div class="min-w-0 flex-1">
						<div class="font-medium text-sm text-(--text-main) truncate">
							{project.owner?.display_name || 'Project Owner'}
						</div>
						<div class="text-xs font-mono text-(--text-muted) truncate">
							@{project.owner?.username || 'member'}
						</div>
					</div>

					<span class="text-xs px-2 py-0.5 rounded-full font-mono bg-(--bg-surface) border border-(--border-hairline) text-(--text-muted)">
						{project.owner?.role || 'USER'}
					</span>
				</div>

				<!-- Key-Value Specs -->
				<div class="flex flex-col gap-3 text-xs sm:text-sm divide-y divide-(--border-hairline)">
					<div class="flex items-center justify-between pt-2">
						<span class="text-(--text-muted)">Hosting Type</span>
						<span class="font-mono text-(--text-main) font-medium">{project.hosting_type || 'CONTAINER'}</span>
					</div>

					<div class="flex items-center justify-between pt-2">
						<span class="text-(--text-muted)">Visibility</span>
						<span class="font-mono text-(--text-main) font-medium">{project.visibility}</span>
					</div>

					{#if project.repository_url}
						<div class="flex items-center justify-between pt-2">
							<span class="text-(--text-muted)">Repository</span>
							<a
								href={project.repository_url}
								target="_blank"
								rel="noreferrer"
								class="text-(--accent-sky) hover:underline font-mono text-xs truncate max-w-[200px]"
							>
								{project.repository_url.replace(/^https?:\/\//, '')} ↗
							</a>
						</div>
					{/if}

					{#if project.documentation_url}
						<div class="flex items-center justify-between pt-2">
							<span class="text-(--text-muted)">Documentation</span>
							<a
								href={project.documentation_url}
								target="_blank"
								rel="noreferrer"
								class="text-(--accent-sky) hover:underline font-mono text-xs truncate max-w-[200px]"
							>
								{project.documentation_url.replace(/^https?:\/\//, '')} ↗
							</a>
						</div>
					{/if}

					<div class="flex items-center justify-between pt-2">
						<span class="text-(--text-muted)">Created</span>
						<span class="font-mono text-xs text-(--text-muted)">{formatDate(project.created_at)}</span>
					</div>

					<div class="flex items-center justify-between pt-2">
						<span class="text-(--text-muted)">Last Updated</span>
						<span class="font-mono text-xs text-(--text-muted)">{formatDate(project.updated_at)}</span>
					</div>
				</div>

				<!-- Tech Stack Tags -->
				{#if project.technology_stack && project.technology_stack.length > 0}
					<div class="pt-2 border-t border-(--border-hairline)">
						<span class="text-xs text-(--text-muted) block mb-2">Technology Stack</span>
						<div class="flex flex-wrap gap-1.5">
							{#each project.technology_stack as tech}
								<span class="text-xs font-mono px-2 py-0.5 rounded-sm bg-(--bg-muted) border border-(--border-hairline) text-(--text-secondary)">
									{tech}
								</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>

		<!-- ══════════════════════════════════════════════
		     3. Project Audit & Activity Trail
		     ══════════════════════════════════════════════ -->
		{#if activities && activities.length > 0}
			<section class="flex flex-col gap-4">
				<div class="flex items-baseline justify-between gap-4">
					<h2 class="font-sans font-semibold text-base text-(--text-main)">
						Activity Trail
					</h2>
					<span class="text-xs text-(--text-muted)">Recent administrative and project lifecycle events</span>
				</div>

				<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md divide-y divide-(--border-hairline) shadow-2xs">
					{#each activities as act (act.id)}
						<div class="p-4 flex items-center justify-between gap-4 text-xs sm:text-sm">
							<div class="flex items-center gap-3">
								<div class="w-2 h-2 rounded-full bg-(--accent-sky) shrink-0"></div>
								<div class="flex flex-col">
									<span class="font-medium text-(--text-main)">
										{act.metadata?.title || act.type}
									</span>
									{#if act.metadata?.message}
										<span class="text-xs text-(--text-secondary) mt-0.5">{act.metadata.message}</span>
									{/if}
								</div>
							</div>

							<div class="text-right shrink-0">
								<span class="block font-mono text-xs text-(--text-muted)">{formatDate(act.created_at)}</span>
								{#if act.actor_name}
									<span class="text-xs text-(--text-secondary)">by {act.actor_name}</span>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			</section>
		{/if}

		<!-- ══════════════════════════════════════════════
		     4. Contextual Comments Moderation
		     ══════════════════════════════════════════════ -->
		<section class="flex flex-col gap-4">
			<div class="flex items-baseline justify-between gap-4">
				<h2 class="font-sans font-semibold text-base text-(--text-main)">
					Comments ({comments.length})
				</h2>
				<span class="text-xs text-(--text-muted)">Contextual moderation for community discussions</span>
			</div>

			{#if comments.length === 0}
				<div class="py-12 text-center text-sm text-(--text-muted) border border-dashed border-(--border-hairline) rounded-md bg-(--bg-surface)">
					No comments posted on this project yet.
				</div>
			{:else}
				<div class="overflow-x-auto bg-(--bg-surface) border border-(--border-hairline) rounded-md shadow-2xs">
					<table class="w-full text-left text-xs sm:text-sm border-collapse font-sans">
						<thead>
							<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-secondary) font-medium">
								<th class="py-3 px-4 font-semibold">Comment</th>
								<th class="py-3 px-3 font-semibold">Author</th>
								<th class="py-3 px-3 font-semibold">Created</th>
								<th class="py-3 px-3 font-semibold">Reports</th>
								<th class="py-3 px-4 font-semibold text-right">Action</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-(--border-hairline)">
							{#each comments as c (c.id)}
								<tr class="hover:bg-(--bg-muted)/20 transition-colors">
									<!-- Comment text -->
									<td class="py-3.5 px-4 max-w-sm">
										{#if c.is_deleted}
											<span class="text-(--text-muted) italic text-xs">[Comment removed]</span>
										{:else}
											<p class="line-clamp-2 text-(--text-main) leading-relaxed">
												{c.content}
											</p>
										{/if}
									</td>

									<!-- Author -->
									<td class="py-3.5 px-3 whitespace-nowrap">
										<span class="font-medium text-(--text-main)">{c.author?.display_name || 'Member'}</span>
										{#if c.author?.username}
											<span class="block text-xs text-(--text-muted) font-mono">@{c.author.username}</span>
										{/if}
									</td>

									<!-- Created -->
									<td class="py-3.5 px-3 font-mono text-xs text-(--text-muted) whitespace-nowrap">
										{formatDate(c.created_at)}
									</td>

									<!-- Reports -->
									<td class="py-3.5 px-3 whitespace-nowrap">
										{#if c.report_count > 0}
											<span class="inline-flex items-center gap-1 text-xs font-medium text-amber-600 dark:text-amber-400">
												<ShieldWarning size={14} weight="bold" />
												<span>{c.report_count}</span>
											</span>
										{:else}
											<span class="text-(--text-muted) text-xs">0</span>
										{/if}
									</td>

									<!-- Delete Action -->
									<td class="py-3.5 px-4 text-right whitespace-nowrap">
										{#if !c.is_deleted}
											<button
												type="button"
												class="p-1.5 rounded-sm border border-(--border-hairline) text-rose-600 dark:text-rose-400 hover:bg-rose-500/10 transition-colors cursor-pointer"
												title="Delete comment"
												on:click={() => openDeleteCommentModal(c)}
											>
												<Trash size={14} />
											</button>
										{:else}
											<span class="text-(--text-muted) text-xs">—</span>
										{/if}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</section>
	{/if}
</div>

<!-- Confirm Delete Comment Modal -->
<ConfirmModal
	bind:open={deleteModalOpen}
	confirmLabel={deletingComment ? 'Deleting...' : 'Delete comment'}
	cancelLabel="Cancel"
	destructive={true}
	loading={deletingComment}
	on:confirm={handleDeleteComment}
>
	<svelte:fragment slot="title">Delete comment?</svelte:fragment>
	<p class="text-sm text-(--text-secondary)">
		The comment will be soft-deleted from this project's discussion. The text will be replaced with <em>"This comment was removed."</em>
	</p>
</ConfirmModal>

<!-- Confirm Suspend Project Modal -->
<ConfirmModal
	bind:open={suspendModalOpen}
	confirmLabel={actionLoading ? 'Suspending...' : 'Suspend project'}
	cancelLabel="Cancel"
	destructive={true}
	loading={actionLoading}
	on:confirm={handleSuspend}
>
	<svelte:fragment slot="title">Suspend Project: {project?.name}?</svelte:fragment>
	<p class="text-sm text-(--text-secondary) leading-relaxed">
		This administrative action immediately flags the project as <strong>SUSPENDED</strong>.
		Availability probes will be paused. This will be recorded in the audit trail.
	</p>
</ConfirmModal>

<!-- Confirm Archive Project Modal -->
<ConfirmModal
	bind:open={archiveModalOpen}
	confirmLabel={actionLoading ? 'Archiving...' : 'Archive project'}
	cancelLabel="Cancel"
	destructive={true}
	loading={actionLoading}
	on:confirm={handleArchive}
>
	<svelte:fragment slot="title">Archive Project: {project?.name}?</svelte:fragment>
	<p class="text-sm text-(--text-secondary) leading-relaxed">
		This will mark the project as <strong>ARCHIVED</strong> and switch visibility to <strong>UNPUBLISHED</strong>.
		Archived projects are removed from public showcases and catalog listings.
	</p>
</ConfirmModal>
