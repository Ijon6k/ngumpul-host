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

<div class="container mx-auto px-6 max-w-4xl py-14 pb-24 flex flex-col gap-10">
	<header class="max-w-2xl">
		<h1 class="font-display font-bold text-3xl sm:text-4xl text-(--text-main) tracking-tight mb-2">Activity</h1>
		<p class="text-base sm:text-lg text-(--text-secondary) leading-relaxed">
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
