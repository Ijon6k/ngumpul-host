<script lang="ts">
	import { onMount } from 'svelte';
	import { projectsApi, extractError } from '$lib/api';
	import type { Project } from '$lib/types/project';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import { resolveProjectStatus } from '$lib/utils/projectStatus';
	import { Plus, Folder, ArrowUpRight } from 'phosphor-svelte';

	let projects: Project[] = [];
	let loading = true;
	let error: string | null = null;

	async function loadProjects() {
		loading = true;
		error = null;
		try {
			const res = await projectsApi.getMyProjects();
			projects = res.projects || [];
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadProjects);
</script>

<svelte:head>
	<title>Projects · Ngumpul Host</title>
</svelte:head>

<div class="w-full flex flex-col">
	<!-- Page Header -->
	<header class="flex flex-col sm:flex-row sm:items-baseline justify-between gap-4 pb-6">
		<div class="flex flex-col gap-1">
			<h1 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-tight">
				Projects
			</h1>
			<p class="text-xs text-(--text-secondary) leading-relaxed">
				All applications and software services linked to your account.
			</p>
		</div>

		<a
			href="/me/requests"
			class="btn btn-primary btn-sm text-xs px-3.5 py-1.5 inline-flex items-center gap-1.5 self-start sm:self-auto"
		>
			<Plus size={13} weight="bold" />
			<span>Request Project</span>
		</a>
	</header>

	<div class="h-px bg-(--border-hairline) mb-8"></div>

	<!-- Projects 2-Column Editorial Grid (Desktop: 2 cols, Mobile: 1 col) -->
	{#if loading}
		<div class="py-20 text-center text-xs text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-5 h-5 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading projects...</span>
		</div>
	{:else if error}
		<div class="p-4 rounded-md bg-rose-500/10 border border-rose-500/30 text-rose-700 dark:text-rose-400 text-xs">
			{error}
		</div>
	{:else if projects.length === 0}
		<div class="py-16 px-6 text-center border border-dashed border-(--border-hairline) rounded-md flex flex-col items-center gap-3 max-w-md mx-auto w-full">
			<div class="w-10 h-10 rounded-full bg-(--bg-muted) flex items-center justify-center text-(--text-muted)">
				<Folder size={20} />
			</div>
			<div class="flex flex-col gap-1">
				<h3 class="font-medium text-sm text-(--text-main)">No projects hosted</h3>
				<p class="text-xs text-(--text-muted) leading-relaxed">
					You haven't published any projects on this server yet. Submit a request to deploy your application.
				</p>
			</div>
			<a href="/me/requests" class="btn btn-primary btn-sm text-xs px-4 py-1.5 mt-2">
				Submit Hosting Request
			</a>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8">
			{#each projects as proj (proj.id)}
				<article class="bg-(--bg-surface) border border-(--border-hairline) hover:border-(--border-subtle) rounded-md overflow-hidden flex flex-col transition-all duration-200 group">
					<!-- Visual Anchor: Large Cover Area -->
					<a href="/me/projects/{proj.id}" class="block relative w-full aspect-video bg-(--bg-muted) overflow-hidden" tabindex="-1">
						<ProjectCover
							src={proj.cover_image_url}
							alt={proj.name}
							name={proj.name}
							aspectRatio="16/9"
						/>
					</a>

					<!-- Content Below Image (Restrained, Editorial) -->
					<div class="p-5 flex-1 flex flex-col justify-between gap-4">
						<div class="flex flex-col gap-2">
							<div class="flex items-center justify-between gap-3">
								<h2 class="font-medium text-base text-(--text-main) group-hover:text-(--accent-sky) transition-colors">
									<a href="/me/projects/{proj.id}">
										{proj.name}
									</a>
								</h2>
								<StatusDot status={resolveProjectStatus(proj).dotStatus} />
							</div>

							<p class="text-xs text-(--text-secondary) line-clamp-2 leading-relaxed font-normal">
								{proj.description || 'No description provided for this project.'}
							</p>

							{#if proj.technology_stack && proj.technology_stack.length > 0}
								<p class="text-xs text-(--text-muted) font-mono pt-1">
									{proj.technology_stack.slice(0, 4).join(' · ')}
								</p>
							{/if}
						</div>

						<div class="pt-3 border-t border-(--border-hairline)/60 flex items-center justify-between text-xs">
							<a
								href="/me/projects/{proj.id}"
								class="text-(--text-secondary) hover:text-(--text-main) font-medium transition-colors"
							>
								Manage project →
							</a>

							{#if proj.slug}
								<a
									href="/projects/{proj.slug}"
									target="_blank"
									rel="noreferrer"
									class="text-(--text-muted) hover:text-(--text-main) transition-colors p-1"
									title="Public Showcase"
								>
									<ArrowUpRight size={13} />
								</a>
							{/if}
						</div>
					</div>
				</article>
			{/each}
		</div>
	{/if}
</div>
