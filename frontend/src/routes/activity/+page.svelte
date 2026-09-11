<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import ActivityTimeline from '$lib/components/ActivityTimeline.svelte';

	let activities: any[] = [];
	let loading = true;

	onMount(async () => {
		try {
			const res = await api.get('/activity');
			activities = res.data?.activities || [];
		} catch (err) {
			console.error('Failed to load activity:', err);
		} finally {
			loading = false;
		}
	});
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
		<div class="py-20 text-center text-(--text-muted) text-xs">
			<p>Retrieving activity feed...</p>
		</div>
	{:else}
		<div class="p-8 sm:p-10 bg-(--bg-surface) border border-(--border-hairline) rounded-md">
			<ActivityTimeline {activities} />
		</div>
	{/if}
</div>
