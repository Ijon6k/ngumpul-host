<script lang="ts">
	import { onMount } from 'svelte';
	import { activityApi, hostingRequestsApi, extractError } from '$lib/api';
	import type { ActivityEvent } from '$lib/types/activity';
	import type { HostingRequest } from '$lib/types/hosting';
	import { Timeline, TimelineItem, Pagination } from '$lib/components/ui';
	import { formatDate } from '$lib/utils/format';
	import { domainSuffix } from '$lib/stores/node';
	import {
		Globe,
		ArrowSquareOut,
		Spinner,
		Tray,
		ClockCounterClockwise,
		CheckCircle,
		Clock,
		Gear,
		XCircle
	} from 'phosphor-svelte';

	let activeTab: 'all' | 'hosting' = 'all';

	let activities: ActivityEvent[] = [];
	let hostingRequests: HostingRequest[] = [];
	let loading = true;
	let error: string | null = null;

	// Pagination states
	let currentPage = 1;
	const pageSize = 15;
	let totalActivities = 0;

	let hostingPage = 1;
	const hostingPageSize = 5;
	$: paginatedHosting = hostingRequests.slice((hostingPage - 1) * hostingPageSize, hostingPage * hostingPageSize);


	async function loadActivities(page: number) {
		loading = true;
		error = null;
		try {
			const res = await activityApi.getMyActivity(page, pageSize);
			activities = res.activities || [];
			totalActivities = res.total ?? activities.length;
			currentPage = page;
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	async function loadData() {
		loading = true;
		error = null;
		try {
			const [actRes, hostRes] = await Promise.all([
				activityApi.getMyActivity(currentPage, pageSize),
				hostingRequestsApi.getMyHostingRequests()
			]);
			activities = actRes.activities || [];
			totalActivities = actRes.total ?? activities.length;
			hostingRequests = hostRes.requests || [];
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadData);

	function formatRelativeTime(dateStr: string): string {
		try {
			const d = new Date(dateStr);
			const diffSec = Math.floor((Date.now() - d.getTime()) / 1000);
			if (diffSec < 60) return 'Just now';
			const diffMin = Math.floor(diffSec / 60);
			if (diffMin < 60) return `${diffMin}m ago`;
			const diffHours = Math.floor(diffMin / 60);
			if (diffHours < 24) return `${diffHours}h ago`;
			const diffDays = Math.floor(diffHours / 24);
			if (diffDays < 30) return `${diffDays}d ago`;
			return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
		} catch {
			return dateStr;
		}
	}

	function getActivityDescription(a: any): string {
		if (a.metadata?.title) return a.metadata.title;
		if (a.metadata?.message) return a.metadata.message;
		switch (a.type) {
			case 'PROJECT_PUBLISHED':
				return `Published ${a.project_name || 'project'} to catalog`;
			case 'PROJECT_UPDATED':
				return `Updated project settings for ${a.project_name || 'project'}`;
			case 'HOSTING_REQUEST_CREATED':
				return `Submitted hosting request for ${a.metadata?.project_name || 'new application'}`;
			case 'COMMENT_CREATED':
				return `Posted comment on ${a.project_name || 'project'}`;
			default:
				return `Activity recorded for ${a.project_name || 'account'}`;
		}
	}

	interface Step {
		name: string;
		status: 'completed' | 'current' | 'upcoming' | 'rejected';
	}

	function getLifecycleSteps(req: HostingRequest): Step[] {
		const status = req.status;

		if (status === 'REJECTED') {
			return [
				{ name: 'Submitted', status: 'completed' },
				{ name: 'Reviewing', status: 'completed' },
				{ name: 'Rejected', status: 'rejected' },
				{ name: 'Setup', status: 'upcoming' },
				{ name: 'Published', status: 'upcoming' }
			];
		}

		if (status === 'COMPLETED') {
			return [
				{ name: 'Submitted', status: 'completed' },
				{ name: 'Reviewing', status: 'completed' },
				{ name: 'Approved', status: 'completed' },
				{ name: 'Setup', status: 'completed' },
				{ name: 'Published', status: 'completed' }
			];
		}

		if (status === 'SETUP') {
			return [
				{ name: 'Submitted', status: 'completed' },
				{ name: 'Reviewing', status: 'completed' },
				{ name: 'Approved', status: 'completed' },
				{ name: 'Setup', status: 'completed' },
				{ name: 'Published', status: 'current' }
			];
		}

		if (status === 'APPROVED') {
			return [
				{ name: 'Submitted', status: 'completed' },
				{ name: 'Reviewing', status: 'completed' },
				{ name: 'Approved', status: 'completed' },
				{ name: 'Setup', status: 'current' },
				{ name: 'Published', status: 'upcoming' }
			];
		}

		// PENDING or default
		return [
			{ name: 'Submitted', status: 'completed' },
			{ name: 'Reviewing', status: 'current' },
			{ name: 'Approved', status: 'upcoming' },
			{ name: 'Setup', status: 'upcoming' },
			{ name: 'Published', status: 'upcoming' }
		];
	}

	function getStatusChip(status: string): { label: string; class: string; icon: any } {
		switch (status) {
			case 'COMPLETED':
				return {
					label: 'Published',
					class: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20',
					icon: CheckCircle
				};
			case 'PENDING':
				return {
					label: 'Reviewing',
					class: 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20',
					icon: Clock
				};
			case 'SETUP':
			case 'APPROVED':
				return {
					label: 'Setup',
					class: 'bg-sky-500/10 text-sky-600 dark:text-sky-400 border border-sky-500/20',
					icon: Gear
				};
			case 'REJECTED':
				return {
					label: 'Rejected',
					class: 'bg-rose-500/10 text-rose-600 dark:text-rose-400 border border-rose-500/20',
					icon: XCircle
				};
			default:
				return {
					label: status,
					class: 'bg-(--bg-muted) text-(--text-secondary) border border-(--border-hairline)',
					icon: Clock
				};
		}
	}
</script>

<svelte:head>
	<title>Activity · Ngumpul Host</title>
</svelte:head>

<div class="w-full flex flex-col max-w-3xl">
	<!-- Page Header -->
	<header class="flex flex-col gap-1 pb-6">
		<h1 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-tight">
			Activity
		</h1>
		<p class="text-xs text-(--text-secondary) leading-relaxed">
			Audit log of actions, project changes, and submissions performed under your account.
		</p>
	</header>

	<div class="h-px bg-(--border-hairline) mb-6"></div>

	<!-- Dual-view switcher -->
	<div class="inline-flex items-center p-0.5 rounded-md bg-(--bg-muted) border border-(--border-hairline) self-start mb-8">
		<button
			type="button"
			class="flex items-center gap-1.5 px-3 py-1.5 rounded-sm text-xs font-medium transition-colors cursor-pointer {activeTab === 'all' ? 'bg-(--bg-surface) text-(--text-main) shadow-xs' : 'text-(--text-secondary) hover:text-(--text-main)'}"
			on:click={() => (activeTab = 'all')}
		>
			<ClockCounterClockwise size={14} />
			<span>All activity</span>
		</button>
		<button
			type="button"
			class="flex items-center gap-1.5 px-3 py-1.5 rounded-sm text-xs font-medium transition-colors cursor-pointer {activeTab === 'hosting' ? 'bg-(--bg-surface) text-(--text-main) shadow-xs' : 'text-(--text-secondary) hover:text-(--text-main)'}"
			on:click={() => (activeTab = 'hosting')}
		>
			<Tray size={14} />
			<span>Hosting history</span>
			{#if hostingRequests.length > 0}
				<span class="ml-1 px-1.5 py-0.2 rounded-full text-xs font-medium bg-(--border-hairline) text-(--text-muted)">
					{hostingRequests.length}
				</span>
			{/if}
		</button>
	</div>

	<!-- Content based on active view -->
	{#if loading}
		<div class="py-20 text-center text-xs text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<Spinner size={20} class="animate-spin text-(--text-muted)" />
			<span>Loading records...</span>
		</div>
	{:else if error}
		<div class="p-3 rounded-sm bg-rose-500/10 text-rose-700 dark:text-rose-400 text-xs border border-rose-500/20">
			{error}
		</div>
	{:else if activeTab === 'all'}
		<!-- All Activity Timeline -->
		{#if activities.length === 0}
			<p class="text-xs text-(--text-muted) py-8">
				No activity recorded yet.
			</p>
		{:else}
			<Timeline density="comfortable">
				{#each activities as act (act.id)}
					<TimelineItem
						title={getActivityDescription(act)}
						timestamp={formatRelativeTime(act.created_at)}
						description={act.metadata?.message || ''}
						href={act.project_slug ? `/projects/${act.project_slug}` : ''}
						linkText={act.project_slug ? `/${act.project_slug}` : ''}
						density="comfortable"
					/>
				{/each}
			</Timeline>

			<Pagination
				{currentPage}
				totalItems={totalActivities}
				{pageSize}
				onPageChange={loadActivities}
			/>
		{/if}
	{:else}
		<!-- Hosting Submissions Lifecycle History -->
		{#if hostingRequests.length === 0}
			<div class="p-8 text-center border border-(--border-hairline) rounded-md text-xs text-(--text-muted) bg-(--bg-surface)/30">
				No hosting requests submitted under this account.
			</div>
		{:else}
			<div class="flex flex-col gap-6">
				{#each paginatedHosting as req (req.id)}
					{@const steps = getLifecycleSteps(req)}
					{@const chip = getStatusChip(req.status)}
					<div class="border border-(--border-hairline) rounded-md p-5 bg-(--bg-surface)/30 flex flex-col gap-4">
						<!-- Top summary line -->
						<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-(--border-hairline)">
							<div class="flex flex-col gap-0.5 min-w-0">
								<div class="flex items-center gap-2 flex-wrap">
									<h2 class="font-medium text-sm text-(--text-main)">{req.project_name}</h2>
									{#if req.subdomain}
										<span class="text-xs font-mono text-(--text-muted) flex items-center gap-1">
											<Globe size={13} />
											<span>{req.subdomain}{$domainSuffix}</span>
										</span>
									{/if}
								</div>
								{#if req.description}
									<p class="text-xs text-(--text-secondary) line-clamp-1 truncate max-w-xl">
										{req.description}
									</p>
								{/if}
							</div>

							<div class="flex items-center gap-3 shrink-0 self-start sm:self-center">
								<span class="text-xs text-(--text-muted)">
									Submitted {formatDate(req.created_at)}
								</span>
								<span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-[4px] text-xs font-medium {chip.class}">
									<svelte:component this={chip.icon} size={12} weight="bold" />
									<span>{chip.label}</span>
								</span>
							</div>
						</div>

						<!-- Minimal Lifecycle Timeline -->
						<div class="pt-2 pb-1">
							<div class="relative flex items-start justify-between w-full">
								<!-- Connecting track line behind nodes -->
								<div class="absolute top-2 left-6 right-6 h-px bg-(--border-hairline) z-0" aria-hidden="true"></div>

								{#each steps as step, idx}
									<div class="relative z-10 flex flex-col items-center text-center flex-1 min-w-0 px-1">
										<!-- Status Node Indicator -->
										<div class="w-4 h-4 rounded-full flex items-center justify-center bg-(--bg-surface) transition-colors">
											{#if step.status === 'completed'}
												<div class="w-2.5 h-2.5 rounded-full bg-[#2E8555] dark:bg-emerald-500 ring-2 ring-(--bg-surface)"></div>
											{:else if step.status === 'current'}
												<div class="w-2.5 h-2.5 rounded-full bg-amber-500 ring-4 ring-amber-500/20 animate-pulse"></div>
											{:else if step.status === 'rejected'}
												<div class="w-2.5 h-2.5 rounded-full bg-rose-500 ring-2 ring-(--bg-surface)"></div>
											{:else}
												<div class="w-2 h-2 rounded-full bg-(--border-subtle) ring-2 ring-(--bg-surface)"></div>
											{/if}
										</div>

										<!-- Step Typography -->
										<div class="flex flex-col gap-0.5 mt-2">
											<span class="text-xs font-medium {step.status === 'completed' ? 'text-(--text-main)' : step.status === 'current' ? 'text-amber-600 dark:text-amber-400 font-semibold' : step.status === 'rejected' ? 'text-rose-600 dark:text-rose-400' : 'text-(--text-muted)'}">
												{step.name}
											</span>
											<span class="text-xs text-(--text-muted)/80 font-normal hidden sm:inline">
												{step.status === 'completed' ? 'Complete' : step.status === 'current' ? 'In progress' : step.status === 'rejected' ? 'Declined' : 'Pending'}
											</span>
										</div>
									</div>
								{/each}
							</div>
						</div>

						<!-- Operator notes / feedback if present -->
						{#if req.admin_notes || req.operator_notes}
							<div class="p-3 bg-amber-500/5 border border-amber-500/20 rounded-sm text-xs text-amber-700 dark:text-amber-300">
								<span class="font-medium">Operator feedback:</span> {req.admin_notes || req.operator_notes}
							</div>
						{/if}

						<!-- Quick links if project is published -->
						{#if req.status === 'COMPLETED'}
							<div class="flex items-center gap-3 pt-1">
								<a
									href="https://{req.subdomain}{$domainSuffix}"
									target="_blank"
									rel="noreferrer"
									class="inline-flex items-center gap-1.5 text-xs text-(--accent-sky) hover:underline font-medium"
								>
									<span>Visit live project</span>
									<ArrowSquareOut size={13} />
								</a>
							</div>
						{/if}
					</div>
				{/each}

				<Pagination
					currentPage={hostingPage}
					totalItems={hostingRequests.length}
					pageSize={hostingPageSize}
					onPageChange={(p) => (hostingPage = p)}
				/>
			</div>
		{/if}
	{/if}
</div>
