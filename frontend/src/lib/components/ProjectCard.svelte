<script lang="ts">
	import StatusDot from './StatusDot.svelte';

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
		owner?: {
			username: string;
			display_name: string;
			avatar_url?: string;
		};
	};
</script>

<article class="group flex flex-col bg-(--bg-surface) border border-(--border-hairline) hover:border-(--accent-sky)/25 rounded-2xl overflow-hidden transition-all duration-200 hover:shadow-sm">
	<!-- Cover image -->
	<a href="/projects/{project.slug}" class="block relative w-full aspect-video bg-(--bg-muted) overflow-hidden">
		{#if project.cover_image_url}
			<img
				src={project.cover_image_url}
				alt={project.name}
				class="w-full h-full object-cover group-hover:scale-[1.03] transition-transform duration-400"
				loading="lazy"
			/>
		{:else}
			<div class="w-full h-full flex items-center justify-center bg-(--bg-muted)">
				<span class="font-display text-5xl font-black text-(--text-muted) opacity-20">{project.name.charAt(0)}</span>
			</div>
		{/if}
		<!-- Status badge overlay -->
		<div class="absolute top-3 right-3">
			<StatusDot status={project.status} />
		</div>
	</a>

	<!-- Card body -->
	<div class="flex flex-col grow p-5 gap-3">
		<div class="flex flex-col gap-1">
			<h3 class="font-display font-bold text-[1.0625rem] text-(--text-main) leading-snug">
				<a href="/projects/{project.slug}" class="hover:text-(--accent-strong) transition-colors">
					{project.name}
				</a>
			</h3>
			{#if project.owner}
				<a href="/people/{project.owner.username}" class="text-xs text-(--text-muted) hover:text-(--text-main) transition-colors">
					by {project.owner.display_name}
				</a>
			{/if}
		</div>

		{#if project.description}
			<p class="text-sm text-(--text-secondary) line-clamp-2 leading-relaxed">{project.description}</p>
		{/if}

		<!-- Footer row -->
		<div class="flex items-center justify-between mt-auto pt-3 border-t border-(--border-subtle) text-xs gap-2">
			{#if project.technology_stack && project.technology_stack.length > 0}
				<div class="flex items-center gap-1.5 flex-wrap">
					{#each project.technology_stack.slice(0, 3) as tech}
						<span class="text-[10px] font-mono bg-(--bg-muted) text-(--text-secondary) px-2 py-0.5 rounded-full border border-(--border-hairline)">{tech}</span>
					{/each}
				</div>
			{:else}
				<div></div>
			{/if}

			{#if project.public_url}
				<a
					href={project.public_url}
					target="_blank"
					rel="noopener noreferrer"
					class="font-medium text-(--accent-strong) hover:underline inline-flex items-center gap-1 shrink-0 text-xs"
				>
					Visit ↗
				</a>
			{/if}
		</div>
	</div>
</article>
