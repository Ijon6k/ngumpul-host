<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';

	let project: any = null;
	let loading = true;
	let error: string | null = null;

	onMount(async () => {
		const slug = $page.params.slug;
		try {
			const res = await api.get(`/projects/${slug}`);
			project = res.data?.project;
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	});
</script>

<div class="container mx-auto px-6 max-w-5xl py-12 pb-24 flex flex-col gap-9">
	{#if loading}
		<div class="py-24 text-center text-(--text-muted) text-xs">
			<p>Loading project details...</p>
		</div>
	{:else if error || !project}
		<div class="py-16 px-6 text-center border border-dashed border-(--border-hairline) rounded-md max-w-md mx-auto w-full">
			<h2 class="font-display text-xl font-bold text-(--text-main) mb-1.5">Project Not Found</h2>
			<p class="text-xs text-(--text-secondary)">{error || 'This project record is either unpublished or does not exist.'}</p>
			<a href="/projects" class="btn btn-secondary btn-sm mt-4">← Back to directory</a>
		</div>
	{:else}
		<!-- Breadcrumb -->
		<div class="flex items-center gap-2 text-xs text-(--text-muted)">
			<a href="/projects" class="hover:text-(--text-main) transition-colors">Projects</a>
			<span>/</span>
			<span class="text-(--text-main) font-medium">{project.name}</span>
		</div>

		<!-- Top Area: Headline, Status, and Overview -->
		<div class="flex flex-col gap-4">
			<div class="flex items-start justify-between flex-wrap gap-4">
				<div>
					<h1 class="font-display font-bold text-3xl sm:text-4xl md:text-5xl text-(--text-main) tracking-tight">{project.name}</h1>
					{#if project.owner}
						<a href="/people/{project.owner.username}" class="text-xs text-(--text-muted) hover:text-(--text-main) transition-colors mt-1 inline-block">
							Maintained by {project.owner.display_name} (@{project.owner.username})
						</a>
					{/if}
				</div>

				<div class="flex items-center gap-3 pt-1 text-xs">
					<StatusDot status={project.status} />
					<span class="text-(--text-muted)">·</span>
					<span class="text-(--text-secondary)">
						{project.hosting_type === 'HOSTED_HERE' ? 'Hosted on node' : 'External project'}
					</span>
				</div>
			</div>

			{#if project.description}
				<p class="text-base sm:text-lg text-(--text-secondary) max-w-2xl leading-relaxed mt-1">{project.description}</p>
			{/if}
		</div>

		<!-- Project Large Artwork Visual -->
		{#if project.cover_image_url}
			<div class="w-full aspect-16/9 max-h-[440px] bg-(--bg-muted) border border-(--border-hairline) rounded-md overflow-hidden">
				<img src={project.cover_image_url} alt={project.name} class="w-full h-full object-cover" />
			</div>
		{/if}

		<!-- Two-Column Specifications -->
		<div class="grid grid-cols-1 md:grid-cols-[1.5fr_1fr] gap-10 pt-4 border-t border-(--border-hairline)">
			<!-- Left: Details & Endpoints -->
			<div class="flex flex-col gap-8">
				<div class="flex flex-col gap-3">
					<span class="text-xs font-semibold text-(--text-main)">Endpoints</span>
					<div class="flex flex-wrap gap-3">
						{#if project.public_url}
							<a href={project.public_url} target="_blank" rel="noopener noreferrer" class="btn btn-primary">
								<span>Visit project instance</span>
								<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
									<path d="M7 7h10v10" />
									<path d="M7 17 17 7" />
								</svg>
							</a>
						{/if}
						{#if project.repository_url}
							<a href={project.repository_url} target="_blank" rel="noopener noreferrer" class="btn btn-secondary">
								<span>Source repository</span>
							</a>
						{/if}
						{#if project.documentation_url}
							<a href={project.documentation_url} target="_blank" rel="noopener noreferrer" class="btn btn-secondary">
								<span>Documentation</span>
							</a>
						{/if}
					</div>
				</div>

				{#if project.technology_stack && project.technology_stack.length > 0}
					<div class="flex flex-col gap-2.5">
						<span class="text-xs font-semibold text-(--text-main)">Technology Stack</span>
						<div class="flex flex-wrap gap-2 text-xs">
							{#each project.technology_stack as tech}
								<span class="px-2.5 py-1 bg-(--bg-muted) border border-(--border-hairline) text-(--text-main) rounded-md">{tech}</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>

			<!-- Right: Owner Info & System Specs -->
			<div class="flex flex-col gap-6">
				{#if project.owner}
					<div class="p-6 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-4">
						<span class="text-xs text-(--text-muted)">Maintainer</span>
						<div class="flex items-center gap-3">
							{#if project.owner.avatar_url}
								<img src={project.owner.avatar_url} alt={project.owner.display_name} class="w-10 h-10 rounded-full object-cover" />
							{:else}
								<div class="w-10 h-10 rounded-full bg-(--accent-soft) text-(--accent-strong) font-bold flex items-center justify-center text-sm">{project.owner.display_name.charAt(0)}</div>
							{/if}
							<div>
								<a href="/people/{project.owner.username}" class="font-bold text-sm text-(--text-main) hover:text-(--accent-strong) transition-colors">{project.owner.display_name}</a>
								<span class="block text-xs text-(--text-muted)">@{project.owner.username}</span>
							</div>
						</div>
						{#if project.owner.bio}
							<p class="text-xs text-(--text-secondary) leading-relaxed">{project.owner.bio}</p>
						{/if}
						<a href="/people/{project.owner.username}" class="text-xs font-semibold text-(--accent-strong) hover:underline mt-1">
							All projects by {project.owner.display_name} →
						</a>
					</div>
				{/if}

				<div class="p-6 bg-(--bg-muted) border border-(--border-hairline) rounded-md flex flex-col gap-2 text-xs">
					<span class="text-xs font-medium text-(--text-main)">Hosting Verification</span>
					<p class="text-(--text-secondary) leading-relaxed">
						Deployed through administrator verification on independent Linux infrastructure.
					</p>
				</div>
			</div>
		</div>
	{/if}
</div>
