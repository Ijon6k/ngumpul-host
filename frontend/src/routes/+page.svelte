<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import HeroSection from '$lib/components/HeroSection.svelte';
	import BentoGrid from '$lib/components/BentoGrid.svelte';
	import ProjectCard from '$lib/components/ProjectCard.svelte';
	import ActivityTimeline from '$lib/components/ActivityTimeline.svelte';
	import { ArrowRight } from 'phosphor-svelte';

	let projects: any[] = [];
	let activities: any[] = [];
	let statusData: any = null;
	let publicServer: any = null;
	let loading = true;

	onMount(async () => {
		try {
			const [projRes, actRes, statRes, serverRes] = await Promise.all([
				api.get('/projects'),
				api.get('/activity'),
				api.get('/status'),
				api.get('/public/server').catch(() => null)
			]);
			projects = projRes.data?.projects || [];
			activities = actRes.data?.activities || [];
			statusData = statRes.data || null;
			publicServer = serverRes?.data || null;
		} catch (e) {
			console.error('Failed to load portal data:', e);
		} finally {
			loading = false;
		}
	});

	$: hostSpecs = publicServer ? { ...statusData?.host_specs, ...publicServer } : (statusData?.host_specs || null);
</script>

<div class="flex flex-col">
	<!-- 1. Hero Section: Corner of the internet & live metric pillars -->
	<HeroSection {hostSpecs} />

	<!-- 2. Bento Grid Section: Real host device telemetry & Pinterest 3D server visual -->
	<BentoGrid {hostSpecs} />

	<!-- 3. Projects Section: Casual showcase of hosted things -->
	<section class="py-16 sm:py-20 border-b border-(--border-hairline) bg-(--bg-canvas)">
		<div class="container mx-auto px-6 max-w-6xl flex flex-col gap-8">
			<div class="flex items-end justify-between flex-wrap gap-4">
				<div>
					<h2 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-[-0.02em]">
						Things we're hosting
					</h2>
					<p class="text-base text-(--text-secondary) mt-1 leading-relaxed font-normal">
						Side projects, hobby tools, and open-source experiments kept alive on this node.
					</p>
				</div>
				<a href="/projects" class="text-sm font-medium text-(--cf-blue) hover:underline inline-flex items-center gap-1">
					<span>View all ({projects.length})</span>
					<ArrowRight size={13} weight="bold" />
				</a>
			</div>

			{#if projects.length === 0 && !loading}
				<div class="text-center py-16 px-4 text-(--text-secondary) border border-dashed border-(--border-hairline) rounded-2xl">
					<p class="text-sm">No projects registered yet.</p>
				</div>
			{:else}
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each projects.slice(0, 6) as project}
						<ProjectCard {project} />
					{/each}
				</div>
			{/if}
		</div>
	</section>

	<!-- 4. Recent Activity Stream -->
	<section class="py-16 sm:py-20 bg-(--bg-canvas)">
		<div class="container mx-auto px-6 max-w-6xl flex flex-col gap-8">
			<div>
				<h2 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-[-0.02em]">
					Recent activity
				</h2>
				<p class="text-base text-(--text-secondary) mt-1 leading-relaxed font-normal">
					What's newly deployed, updated, or registered around the server.
				</p>
			</div>

			<ActivityTimeline {activities} />

			<div class="pt-2">
				<a href="/activity" class="text-sm font-medium text-(--cf-blue) hover:underline inline-flex items-center gap-1">
					<span>View complete activity timeline</span>
					<ArrowRight size={13} weight="bold" />
				</a>
			</div>
		</div>
	</section>
</div>
