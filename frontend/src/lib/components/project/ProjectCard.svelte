<script lang="ts">
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import { ArrowUpRight } from 'phosphor-svelte';
	import { resolveProjectStatus } from '$lib/utils/projectStatus';

	interface Props {
		project: {
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
	}

	let { project }: Props = $props();

	function formatHostedDate(dateStr?: string): string {
		if (!dateStr) return '';
		try {
			const d = new Date(dateStr);
			return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
		} catch {
			return '';
		}
	}

	let projectStatus = $derived(resolveProjectStatus(project));
	let visibleTechs = $derived((project.technology_stack || []).slice(0, 4));
	let hostedDate = $derived(formatHostedDate(project.published_at || project.created_at));
</script>

<article class="group relative flex flex-col bg-(--bg-surface) border border-(--border-hairline) hover:border-(--border-subtle) rounded-md overflow-hidden transition-all duration-200 hover:shadow-xs cursor-pointer">
	<!-- Cover area (16:9 aspect ratio with status badge overlay) -->
	<div class="relative w-full aspect-[16/9] overflow-hidden bg-(--bg-muted)">
		<ProjectCover
			src={project.cover_image_url}
			alt={project.name}
			name={project.name}
			aspectRatio="full"
			class="w-full h-full !rounded-none"
		/>
		<div class="absolute top-0 right-0 z-10 pointer-events-none">
			<div class="px-2.5 py-1 rounded-bl-md bg-(--bg-surface)/90 dark:bg-neutral-900/90 backdrop-blur-xs border-b border-l border-(--border-hairline) shadow-2xs">
				<span class="text-xs font-medium {projectStatus.colorClass}">
					{projectStatus.label}
				</span>
			</div>
		</div>
	</div>

	<!-- Card Body: Editorial Hierarchy -->
	<div class="flex flex-col grow p-5 gap-3">
		<div class="flex flex-col gap-1.5">
			<h3 class="font-sans font-semibold text-base sm:text-lg text-(--text-main) tracking-[-0.01em] line-clamp-1 leading-snug">
				<a href="/projects/{project.slug}" class="hover:text-(--accent-strong) transition-colors after:absolute after:inset-0 after:z-0">
					{project.name}
				</a>
			</h3>

			<p class="text-xs sm:text-sm text-(--text-secondary) line-clamp-2 min-h-[2.5rem] sm:min-h-[2.85rem] leading-relaxed font-normal">
				{project.description || ''}
			</p>
		</div>

		<!-- Creator & Tech Stack (Consistent vertical slot, clean typographic hierarchy) -->
		<div class="flex items-center gap-x-2 text-xs text-(--text-muted) min-h-[1.35rem] relative z-10">
			{#if project.owner}
				<a
					href="/people/{project.owner.username}"
					class="font-medium text-(--text-main) hover:underline shrink-0"
					onclick={(e) => e.stopPropagation()}
				>
					{project.owner.display_name}
				</a>
			{/if}

			{#if visibleTechs.length > 0}
				{#if project.owner}<span class="text-(--border-subtle) shrink-0">·</span>{/if}
				<span class="text-(--text-secondary) truncate">
					{visibleTechs.join(' · ')}
				</span>
			{/if}
		</div>

		<!-- Footer Metadata: Hosted Since & Direct Link -->
		<div class="flex items-center justify-between mt-auto pt-3 border-t border-(--border-hairline) text-xs text-(--text-muted) relative z-10">
			<span class="text-xs">
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
					onclick={(e) => e.stopPropagation()}
				>
					<span>Visit</span>
					<ArrowUpRight size={13} weight="bold" />
				</a>
			{/if}
		</div>
	</div>
</article>
