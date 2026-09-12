<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/auth';
	import { projectsApi, hostingRequestsApi, activityApi, notificationsApi } from '$lib/api';
	import type { Project } from '$lib/types/project';
	import type { HostingRequest } from '$lib/types/hosting';
	import type { ActivityEvent } from '$lib/types/activity';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import { Timeline, TimelineItem } from '$lib/components/ui';
	import { Plus, ArrowRight, ArrowUpRight } from 'phosphor-svelte';

	let myProjects: Project[] = [];
	let myRequests: HostingRequest[] = [];
	let myActivities: ActivityEvent[] = [];
	let unreadNotificationsCount = 0;
	let loading = true;

	onMount(async () => {
		try {
			const [projRes, reqRes, actRes, notifRes] = await Promise.all([
				projectsApi.getMyProjects(),
				hostingRequestsApi.getMyHostingRequests(),
				activityApi.getMyActivity(),
				notificationsApi.getNotifications()
			]);
			myProjects = projRes.projects || [];
			myRequests = reqRes.requests || [];
			myActivities = actRes.activities || [];
			const notifs = notifRes.notifications || [];
			unreadNotificationsCount = notifs.filter((n) => !n.is_read).length;
		} catch (e) {
			console.error('Failed to load workspace overview:', e);
		} finally {
			loading = false;
		}
	});

	function getGreeting(): string {
		const hour = new Date().getHours();
		if (hour < 12) return 'Good morning';
		if (hour < 18) return 'Good afternoon';
		return 'Good evening';
	}

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
			return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
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
</script>

<svelte:head>
	<title>Workspace Overview · Ngumpul Host</title>
</svelte:head>

<div class="w-full flex flex-col">
	<!-- Editorial Introduction Header (Typography-first, zero metric card clutter) -->
	<header class="flex flex-col gap-1 pb-6">
		<div class="flex items-start justify-between gap-4 flex-wrap">
			<div class="flex flex-col">
				<h1 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-tight">
					{getGreeting()}, {$user?.display_name || 'Site Administrator'}.
				</h1>
				<p class="text-sm sm:text-base text-(--text-secondary) font-normal mt-1">
					Here's what's happening with your projects.
				</p>
			</div>

			<a
				href="/me/requests"
				class="btn btn-primary btn-sm text-xs px-3.5 py-1.5 inline-flex items-center gap-1.5 self-start"
			>
				<Plus size={13} weight="bold" />
				<span>Request Hosting</span>
			</a>
		</div>

		<!-- Lightweight summary line (real data, pure typography) -->
		<p class="text-xs text-(--text-muted) pt-3">
			{myProjects.length} {myProjects.length === 1 ? 'hosted project' : 'hosted projects'} ·
			{myRequests.length} {myRequests.length === 1 ? 'hosting request' : 'hosting requests'} ·
			{unreadNotificationsCount} {unreadNotificationsCount === 1 ? 'unread notification' : 'unread notifications'}
		</p>
	</header>

	<div class="h-px bg-(--border-hairline) my-6"></div>

	<!-- Section: Your Projects -->
	<section class="flex flex-col gap-4">
		<div class="flex items-baseline justify-between gap-4">
			<h2 class="text-base sm:text-lg font-medium tracking-tight text-(--text-main)">
				Your projects
			</h2>
			{#if myProjects.length > 0}
				<a href="/me/projects" class="text-xs text-(--text-secondary) hover:text-(--text-main) transition-colors">
					View all ({myProjects.length}) →
				</a>
			{/if}
		</div>

		{#if loading}
			<div class="py-12 text-center text-xs text-(--text-muted)">
				Loading projects...
			</div>
		{:else if myProjects.length === 0}
			<div class="py-10 px-6 border border-dashed border-(--border-hairline) rounded-md text-center flex flex-col items-center gap-2">
				<p class="text-xs font-medium text-(--text-main)">No hosted projects yet</p>
				<p class="text-xs text-(--text-muted) max-w-sm">
					Submit an application to request reverse-proxy allocation on this community node.
				</p>
				<a href="/me/requests" class="btn btn-primary btn-sm text-xs px-3.5 py-1.5 mt-2">
					Submit Hosting Request
				</a>
			</div>
		{:else}
			<!-- Restrained horizontal list (no heavy card-in-card boxes) -->
			<div class="flex flex-col divide-y divide-(--border-hairline)">
				{#each myProjects as proj (proj.id)}
					<div class="py-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:bg-(--bg-muted)/20 px-2 -mx-2 rounded-sm transition-colors group">
						<div class="flex items-center gap-3.5 min-w-0">
							<div class="w-14 h-10 rounded-xs overflow-hidden shrink-0 border border-(--border-hairline) bg-(--bg-muted)">
								<ProjectCover
									src={proj.cover_image_url}
									alt={proj.name}
									name={proj.name}
									aspectRatio="4/3"
								/>
							</div>

							<div class="flex flex-col min-w-0 gap-0.5">
								<div class="flex items-center gap-2.5">
									<a
										href="/me/projects/{proj.id}"
										class="font-medium text-sm text-(--text-main) hover:text-(--accent-sky) transition-colors truncate"
									>
										{proj.name}
									</a>
									<StatusDot status={proj.status} />
								</div>
								<p class="text-xs text-(--text-secondary) truncate max-w-lg">
									{proj.description || 'No description provided.'}
								</p>
							</div>
						</div>

						<div class="flex items-center gap-3 self-end sm:self-center shrink-0 text-xs">
							<a
								href="/me/projects/{proj.id}"
								class="text-(--text-secondary) hover:text-(--text-main) font-medium transition-colors"
							>
								Manage project →
							</a>
							<a
								href="/projects/{proj.slug}"
								target="_blank"
								rel="noreferrer"
								class="text-(--text-muted) hover:text-(--text-main) transition-colors p-1"
								title="Public Showcase"
							>
								<ArrowUpRight size={13} />
							</a>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</section>

	<div class="h-px bg-(--border-hairline) my-10"></div>

	<!-- Section: Recent Activity (Vertical Dot Timeline, no card boxes) -->
	<section class="flex flex-col gap-4">
		<div class="flex items-baseline justify-between gap-4">
			<h2 class="text-base sm:text-lg font-medium tracking-tight text-(--text-main)">
				Recent activity
			</h2>
			{#if myActivities.length > 0}
				<a href="/me/activity" class="text-xs text-(--text-secondary) hover:text-(--text-main) transition-colors">
					All activity →
				</a>
			{/if}
		</div>

		{#if loading}
			<div class="py-8 text-center text-xs text-(--text-muted)">
				Loading activity...
			</div>
		{:else if myActivities.length === 0}
			<p class="text-xs text-(--text-muted) py-4">No recent activity recorded.</p>
		{:else}
			<div class="pt-2">
				<Timeline density="compact">
					{#each myActivities.slice(0, 5) as act (act.id)}
						<TimelineItem
							title={getActivityDescription(act)}
							timestamp={formatRelativeTime(act.created_at)}
							href={act.project_slug ? `/projects/${act.project_slug}` : ''}
							linkText={act.project_slug ? `/${act.project_slug}` : ''}
							density="compact"
						/>
					{/each}
				</Timeline>
			</div>
		{/if}
	</section>
</div>
