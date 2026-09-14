<script lang="ts">
	import { onMount } from 'svelte';
	import { activityApi } from '$lib/api';
	import type { ActivityEvent } from '$lib/types/activity';
	import ActivityTimeline from '$lib/components/ActivityTimeline.svelte';
	import { Pagination } from '$lib/components/ui';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';

	let activities: ActivityEvent[] = [];
	let loading = true;
	let currentPage = 1;
	const pageSize = 20;
	let totalActivities = 0;

	async function loadActivities(page: number) {
		loading = true;
		try {
			const res = await activityApi.getPublicActivity(page, pageSize);
			activities = res.activities || [];
			totalActivities = res.total ?? activities.length;
			currentPage = page;
		} catch (err) {
			console.error('Failed to load activity:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => loadActivities(1));
</script>

<svelte:head>
	<title>Activity · Ngumpul Host</title>
	<meta name="description" content="A public chronological record of deployments, project publications, and community milestones." />
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-4xl py-8 sm:py-12 pb-24 flex flex-col gap-8 sm:gap-10">
	<header class="max-w-2xl">
		<h1 class="font-sans font-normal text-3xl sm:text-4xl text-(--text-main) tracking-[-0.03em] mb-2">Activity</h1>
		<p class="text-sm sm:text-base text-(--text-secondary) leading-relaxed font-normal">
			A public chronological record of deployments, project publications, and community milestones.
		</p>
	</header>

	{#if loading}
		<div class="p-8 sm:p-10 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-6" role="status" aria-live="polite">
			{#each Array(6) as _}
				<div class="flex items-start gap-4" aria-hidden="true">
					<Skeleton variant="avatar" class="!h-6 !w-6" />
					<div class="flex flex-col gap-2 flex-1 min-w-0">
						<Skeleton variant="line" class="w-1/2" />
						<Skeleton variant="line" class="w-2/3" />
						<Skeleton variant="line" class="w-1/4" />
					</div>
				</div>
			{/each}
		</div>
		<span class="sr-only">Retrieving activity feed...</span>
	{:else}
		<div class="p-8 sm:p-10 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-6">
			<ActivityTimeline {activities} />
			<Pagination
				{currentPage}
				totalItems={totalActivities}
				{pageSize}
				onPageChange={loadActivities}
			/>
		</div>
	{/if}
</div>
