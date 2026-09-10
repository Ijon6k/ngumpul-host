<script lang="ts">
	import StatusDot from './StatusDot.svelte';
	import InfrastructureDiagram from './InfrastructureDiagram.svelte';
	import DotMatrix from './DotMatrix.svelte';

	export let featuredProject: {
		name: string;
		slug: string;
		description: string;
		cover_image_url?: string;
		technology_stack?: string[];
		status: string;
		public_url?: string;
		owner?: {
			display_name: string;
		};
	} | null = null;

	export let counts = {
		projects: 4,
		members: 12
	};

	export let systemStatus = 'OPERATIONAL';
</script>

<div class="flex flex-col gap-6 w-full">
	<!-- Asymmetric Top Row -->
	<div class="grid grid-cols-1 lg:grid-cols-[1.4fr_1fr] gap-6">
		<!-- Featured Project Editorial Block -->
		{#if featuredProject}
			<div class="flex flex-col bg-(--bg-surface) border border-(--border-hairline) rounded-md overflow-hidden group">
				<a href="/projects/{featuredProject.slug}" class="block relative w-full h-[240px] sm:h-[300px] bg-(--bg-muted) overflow-hidden border-b border-(--border-hairline)">
					{#if featuredProject.cover_image_url}
						<img
							src={featuredProject.cover_image_url}
							alt={featuredProject.name}
							class="w-full h-full object-cover group-hover:scale-[1.01] transition-transform duration-300"
						/>
					{:else}
						<img
							src="/assets/infra_abstract_cover.jpg"
							alt="Infrastructure art"
							class="w-full h-full object-cover group-hover:scale-[1.01] transition-transform duration-300"
						/>
					{/if}
				</a>

				<div class="p-6 sm:p-7 flex flex-col gap-3 grow">
					<div class="flex items-baseline justify-between gap-4">
						<div>
							<h3 class="font-display font-bold text-xl sm:text-2xl text-(--text-main)">
								<a href="/projects/{featuredProject.slug}" class="hover:text-(--accent-strong) transition-colors">
									{featuredProject.name}
								</a>
							</h3>
							{#if featuredProject.owner}
								<span class="text-xs text-(--text-muted) block mt-0.5">by {featuredProject.owner.display_name}</span>
							{/if}
						</div>
						<StatusDot status={featuredProject.status} />
					</div>

					<p class="text-sm text-(--text-secondary) leading-relaxed">{featuredProject.description}</p>

					<div class="mt-auto pt-4 border-t border-(--border-subtle) flex items-center justify-between text-xs">
						{#if featuredProject.technology_stack && featuredProject.technology_stack.length > 0}
							<span class="text-xs text-(--text-muted)">
								{featuredProject.technology_stack.join(' · ')}
							</span>
						{:else}
							<span></span>
						{/if}

						{#if featuredProject.public_url}
							<a
								href={featuredProject.public_url}
								target="_blank"
								rel="noopener noreferrer"
								class="font-medium text-(--accent-strong) hover:underline inline-flex items-center gap-1"
							>
								<span>Visit instance</span>
								<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
									<path d="M7 7h10v10" />
									<path d="M7 17 17 7" />
								</svg>
							</a>
						{/if}
					</div>
				</div>
			</div>
		{/if}

		<!-- Metric Counters Column -->
		<div class="flex flex-col gap-6">
			<!-- Projects Metric Card with DotMatrix -->
			<div class="p-6 sm:p-7 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col justify-between grow">
				<div class="flex items-center justify-between">
					<span class="text-xs text-(--text-muted)">Active Deployments</span>
					<div class="text-(--accent-orange)">
						<DotMatrix text={String(counts.projects).padStart(2, '0')} dotSize={2.5} gap={2} />
					</div>
				</div>
				<div class="my-4">
					<div class="font-display text-4xl sm:text-5xl font-bold text-(--text-main) leading-none">
						{counts.projects}
					</div>
				</div>
				<p class="text-xs text-(--text-secondary) leading-relaxed">
					Independent web applications, APIs, and micro-tools hosted and running on our shared node.
				</p>
			</div>

			<!-- Members Metric Card with DotMatrix -->
			<div class="p-6 sm:p-7 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col justify-between grow">
				<div class="flex items-center justify-between">
					<span class="text-xs text-(--text-muted)">Community Members</span>
					<div class="text-(--text-main)">
						<DotMatrix text={String(counts.members).padStart(2, '0')} dotSize={2.5} gap={2} />
					</div>
				</div>
				<div class="my-4">
					<div class="font-display text-4xl sm:text-5xl font-bold text-(--text-main) leading-none">
						{counts.members}
					</div>
				</div>
				<p class="text-xs text-(--text-secondary) leading-relaxed">
					Engineers, designers, and builders with direct access to deploy, test, and showcase their craft.
				</p>
			</div>
		</div>
	</div>

	<!-- Bottom Row: Infrastructure Topology Component -->
	<InfrastructureDiagram {systemStatus} />
</div>
