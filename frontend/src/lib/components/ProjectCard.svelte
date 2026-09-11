<script lang="ts">
	import StatusDot from './StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import { ArrowUpRight } from 'phosphor-svelte';

	export let project: {
		id: string;
		name: string;
		slug: string;
		description: string;
		cover_image_url?: string;
		technology_stack?: string[];
		hosting_type?: string;
		public_url?: string;
		status: string;
		published_at?: string;
		created_at?: string;
		owner?: {
			username: string;
			display_name: string;
			avatar_url?: string;
		};
	};

	function formatHostedDate(dateStr?: string): string {
		if (!dateStr) return '';
		try {
			const d = new Date(dateStr);
			return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
		} catch {
			return '';
		}
	}

	$: visibleTechs = (project.technology_stack || []).slice(0, 4);
	$: hostedDate = formatHostedDate(project.published_at || project.created_at);
</script>

<article class="group relative flex flex-col bg-(--bg-surface) border border-(--border-hairline) hover:border-(--border-subtle) rounded-md overflow-hidden transition-all duration-200 hover:shadow-xs">
	<!-- Cover area (4:3 aspect ratio with status badge overlay) -->
	<a href="/projects/{project.slug}" class="block relative w-full overflow-hidden bg-(--bg-muted)" tabindex="-1">
		<ProjectCover
			src={project.cover_image_url}
			alt={project.name}
			name={project.name}
			aspectRatio="4/3"
			class="group-hover:scale-[1.015] transition-transform duration-500 ease-out"
		/>
		<div class="absolute top-3 right-3 z-10 pointer-events-none">
			<StatusDot status={project.status} />
		</div>
	</a>

	<!-- Card Body: Editorial Hierarchy -->
	<div class="flex flex-col grow p-5 gap-3.5">
		<div class="flex flex-col gap-1.5">
			<h3 class="font-sans font-semibold text-base sm:text-lg text-(--text-main) tracking-[-0.01em] leading-snug">
				<a href="/projects/{project.slug}" class="hover:text-(--accent-strong) transition-colors">
					{project.name}
				</a>
			</h3>

			{#if project.description}
				<p class="text-xs sm:text-sm text-(--text-secondary) line-clamp-2 leading-relaxed font-normal">
					{project.description}
				</p>
			{/if}
		</div>

		<!-- Creator & Tech Stack (Clean typographic hierarchy, no chip/pill spam) -->
		<div class="flex items-center flex-wrap gap-x-2 gap-y-1 text-xs text-(--text-muted) pt-1">
			{#if project.owner}
				<a href="/people/{project.owner.username}" class="font-medium text-(--text-main) hover:underline">
					{project.owner.display_name}
				</a>
			{/if}

			{#if visibleTechs.length > 0}
				{#if project.owner}<span class="text-(--border-subtle)">·</span>{/if}
				<span class="text-(--text-secondary)">
					{visibleTechs.join(' · ')}
				</span>
			{/if}
		</div>

		<!-- Footer Metadata: Hosted Since & Direct Link -->
		<div class="flex items-center justify-between mt-auto pt-3 border-t border-(--border-hairline) text-xs text-(--text-muted)">
			<span class="text-[11px]">
				{#if project.hosting_type === 'HOSTED_HERE'}
					{hostedDate ? `Hosted since ${hostedDate}` : 'Hosted on node'}
				{:else}
					{hostedDate ? `Registered ${hostedDate}` : 'External service'}
				{/if}
			</span>

			{#if project.public_url}
				<a
					href={project.public_url}
					target="_blank"
					rel="noopener noreferrer"
					class="font-medium text-(--accent-strong) hover:underline inline-flex items-center gap-1 shrink-0 text-xs"
					on:click|stopPropagation
				>
					<span>Visit</span>
					<ArrowUpRight size={13} weight="bold" />
				</a>
			{/if}
		</div>
	</div>
</article>
