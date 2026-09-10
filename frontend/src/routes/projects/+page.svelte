<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import ProjectCard from '$lib/components/ProjectCard.svelte';

	let projects: any[] = [];
	let loading = true;
	let searchQuery = '';
	let selectedFilter = 'ALL';

	async function loadProjects() {
		loading = true;
		try {
			const res = await api.get('/projects');
			projects = res.data?.projects || [];
		} catch (err) {
			console.error('Failed to load projects:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadProjects);

	$: filteredProjects = projects.filter((p) => {
		const matchQuery =
			p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			p.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
			(p.technology_stack && p.technology_stack.some((t: string) => t.toLowerCase().includes(searchQuery.toLowerCase())));

		const matchFilter =
			selectedFilter === 'ALL' ||
			(selectedFilter === 'HOSTED' && p.hosting_type === 'HOSTED_HERE') ||
			(selectedFilter === 'EXTERNAL' && p.hosting_type === 'EXTERNAL');

		return matchQuery && matchFilter;
	});
</script>

<div class="container mx-auto px-6 max-w-6xl py-14 pb-24 flex flex-col gap-10">
	<header class="max-w-2xl">
		<h1 class="font-display font-bold text-3xl sm:text-4xl text-(--text-main) tracking-tight mb-2">Projects</h1>
		<p class="text-base sm:text-lg text-(--text-secondary) leading-relaxed">
			Independent applications, microservices, and experiments hosted or registered on our system.
		</p>
	</header>

	<!-- Refined Filter Navigation (Segmented text controls without pills) -->
	<div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4 pb-4 border-b border-(--border-hairline)">
		<div class="w-full sm:max-w-sm">
			<input
				type="text"
				placeholder="Filter by name, description, or stack..."
				class="w-full px-3.5 py-2 bg-(--bg-surface) border border-(--border-hairline) focus:border-(--accent-sky) rounded-md text-xs sm:text-sm text-(--text-main) placeholder:text-(--text-muted) outline-none transition-colors"
				bind:value={searchQuery}
			/>
		</div>

		<!-- Clean text tabs -->
		<div class="flex items-center gap-6 text-xs font-medium">
			<button
				type="button"
				class="pb-1 transition-colors cursor-pointer {selectedFilter === 'ALL' ? 'text-(--text-main) border-b-2 border-(--text-main) font-semibold' : 'text-(--text-muted) hover:text-(--text-main)'}"
				on:click={() => (selectedFilter = 'ALL')}
			>
				All ({projects.length})
			</button>
			<button
				type="button"
				class="pb-1 transition-colors cursor-pointer {selectedFilter === 'HOSTED' ? 'text-(--text-main) border-b-2 border-(--text-main) font-semibold' : 'text-(--text-muted) hover:text-(--text-main)'}"
				on:click={() => (selectedFilter = 'HOSTED')}
			>
				Hosted here
			</button>
			<button
				type="button"
				class="pb-1 transition-colors cursor-pointer {selectedFilter === 'EXTERNAL' ? 'text-(--text-main) border-b-2 border-(--text-main) font-semibold' : 'text-(--text-muted) hover:text-(--text-main)'}"
				on:click={() => (selectedFilter = 'EXTERNAL')}
			>
				External
			</button>
		</div>
	</div>

	{#if loading}
		<div class="py-20 text-center text-(--text-muted) text-xs">
			<p>Loading project directory...</p>
		</div>
	{:else if filteredProjects.length === 0}
		<div class="py-16 px-4 text-center border border-dashed border-(--border-hairline) rounded-md text-(--text-muted)">
			<p class="text-sm font-medium">No projects matched your criteria.</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each filteredProjects as project (project.id)}
				<ProjectCard {project} />
			{/each}
		</div>
	{/if}
</div>
