<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import ProjectCard from '$lib/components/ProjectCard.svelte';
	import Pagination from '$lib/components/ui/Pagination.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import { MagnifyingGlass, FolderSimpleDashed } from 'phosphor-svelte';

	let projects: any[] = [];
	let loading = true;
	let searchQuery = '';
	let searchTimeout: any = null;
	let selectedFilter = 'ALL';
	let currentPage = 1;
	let totalPages = 1;
	let totalProjects = 0;
	const pageSize = 12;

	async function loadProjects(pageNum = 1) {
		loading = true;
		try {
			const params: Record<string, any> = {
				page: pageNum,
				limit: pageSize
			};

			if (searchQuery.trim()) {
				params.q = searchQuery.trim();
			}

			if (selectedFilter === 'HOSTED') {
				params.type = 'HOSTED_HERE';
			} else if (selectedFilter === 'EXTERNAL') {
				params.type = 'EXTERNAL';
			}

			const res = await api.get('/projects', { params });
			projects = res.data?.projects || [];
			totalProjects = res.data?.total ?? projects.length;
			currentPage = res.data?.page ?? pageNum;
			totalPages = res.data?.total_pages ?? Math.max(1, Math.ceil(totalProjects / pageSize));
		} catch (err) {
			console.error('Failed to load projects:', err);
			projects = [];
		} finally {
			loading = false;
		}
	}

	function handleSearchInput() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			currentPage = 1;
			loadProjects(1);
		}, 300);
	}

	function handleFilterChange(filter: string) {
		if (selectedFilter === filter) return;
		selectedFilter = filter;
		currentPage = 1;
		loadProjects(1);
	}

	function handlePageChange(newPage: number) {
		currentPage = newPage;
		loadProjects(newPage);
		if (typeof window !== 'undefined') {
			window.scrollTo({ top: 0, behavior: 'smooth' });
		}
	}

	onMount(() => {
		loadProjects(1);
	});
</script>

<svelte:head>
	<title>Project Catalog — Ngumpul Host</title>
	<meta name="description" content="Collection of independent projects, experiments, and tools built and hosted by the Ngumpul Host community." />
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-14 pb-28 flex flex-col gap-8 sm:gap-12">
	<!-- Editorial Page Header -->
	<header class="max-w-2xl flex flex-col gap-2.5">
		<h1 class="font-sans font-normal text-3xl sm:text-4xl text-(--text-main) tracking-[-0.03em]">
			Community Projects
		</h1>
		<p class="text-sm sm:text-base text-(--text-secondary) leading-relaxed font-normal">
			A quiet collection of side projects, experimental software, and bots hosted and shared on this server.
		</p>
	</header>

	<!-- Refined Filter & Search Bar -->
	<div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4 pb-4 border-b border-(--border-hairline)">
		<div class="relative w-full sm:max-w-xs">
			<MagnifyingGlass size={15} weight="regular" class="absolute left-3 top-1/2 -translate-y-1/2 text-(--text-muted) pointer-events-none" />
			<input
				type="text"
				placeholder="Search by name, creator, or technology..."
				class="w-full pl-9 pr-3.5 py-2 bg-(--bg-surface) border border-(--border-hairline) focus:border-(--accent-sky) rounded-md text-xs sm:text-sm text-(--text-main) placeholder:text-(--text-muted) outline-none transition-colors"
				bind:value={searchQuery}
				on:input={handleSearchInput}
			/>
		</div>

		<!-- Clean Segmented Text Filter Tabs -->
		<div class="flex items-center gap-6 text-xs font-medium">
			<button
				type="button"
				class="pb-1 transition-colors cursor-pointer {selectedFilter === 'ALL' ? 'text-(--text-main) border-b-2 border-(--text-main) font-semibold' : 'text-(--text-muted) hover:text-(--text-main)'}"
				on:click={() => handleFilterChange('ALL')}
			>
				All ({totalProjects})
			</button>
			<button
				type="button"
				class="pb-1 transition-colors cursor-pointer {selectedFilter === 'HOSTED' ? 'text-(--text-main) border-b-2 border-(--text-main) font-semibold' : 'text-(--text-muted) hover:text-(--text-main)'}"
				on:click={() => handleFilterChange('HOSTED')}
			>
				Hosted Here
			</button>
			<button
				type="button"
				class="pb-1 transition-colors cursor-pointer {selectedFilter === 'EXTERNAL' ? 'text-(--text-main) border-b-2 border-(--text-main) font-semibold' : 'text-(--text-muted) hover:text-(--text-main)'}"
				on:click={() => handleFilterChange('EXTERNAL')}
			>
				External
			</button>
		</div>
	</div>

	<!-- Main Project Grid: 2-column desktop, 1-column mobile -->
	{#if loading}
		<div class="py-24 text-center text-(--text-muted) text-xs flex flex-col items-center justify-center gap-3">
			<div class="w-5 h-5 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
			<p>Loading project catalog...</p>
		</div>
	{:else if projects.length === 0}
		<EmptyState
			title={searchQuery ? 'No Matching Projects' : 'No Projects Published Yet'}
			description={searchQuery
				? 'Try using different search keywords or clearing active filters.'
				: 'Once projects are reviewed and published by administrators, they will appear here.'}
		>
			<div slot="action">
				{#if searchQuery}
					<button
						type="button"
						class="btn btn-secondary text-xs px-3.5 py-1.5"
						on:click={() => { searchQuery = ''; loadProjects(1); }}
					>
						Reset Search
					</button>
				{/if}
			</div>
		</EmptyState>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8">
			{#each projects as project (project.id)}
				<ProjectCard {project} />
			{/each}
		</div>

		<!-- Pagination Control -->
		{#if totalProjects > pageSize}
			<div class="pt-6 border-t border-(--border-hairline)">
				<Pagination
					{currentPage}
					totalItems={totalProjects}
					{pageSize}
					onPageChange={handlePageChange}
				/>
			</div>
		{/if}
	{/if}
</div>
