<script lang="ts">
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import { Timeline, TimelineItem } from '$lib/components/ui';

	let activities: any[] = [];
	let loading = true;
	let error: string | null = null;

	async function loadActivities() {
		loading = true;
		error = null;
		try {
			const res = await api.get('/me/activity');
			activities = res.data?.activities || [];
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadActivities);

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
</script>

<svelte:head>
	<title>Activity · Ngumpul Host</title>
</svelte:head>

<div class="w-full flex flex-col max-w-2xl">
	<!-- Page Header -->
	<header class="flex flex-col gap-1 pb-6">
		<h1 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-tight">
			Activity
		</h1>
		<p class="text-xs text-(--text-secondary) leading-relaxed">
			Audit log of actions, project changes, and submissions performed under your account.
		</p>
	</header>

	<div class="h-px bg-(--border-hairline) mb-8"></div>

	<!-- Activity Timeline -->
	{#if loading}
		<div class="py-20 text-center text-xs text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-5 h-5 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading activity history...</span>
		</div>
	{:else if error}
		<div class="p-3 rounded-sm bg-rose-500/10 text-rose-700 dark:text-rose-400 text-xs">
			{error}
		</div>
	{:else if activities.length === 0}
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
	{/if}
</div>
