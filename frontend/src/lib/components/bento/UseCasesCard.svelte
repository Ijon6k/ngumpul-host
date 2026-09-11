<script lang="ts">
	export let projectsTotal: number = 3;
	export let projectsOnline: number = 3;

	const useCases = [
		'Side projects',
		'Telegram & Discord bots',
		'Webhooks & scrapers',
		'Learning & staging',
		'Lightweight APIs'
	];
</script>

<div class="h-full w-full bg-white rounded-3xl p-5 sm:p-6 border border-neutral-200/80 flex flex-col justify-between shadow-xs min-h-[140px]">
	<div class="flex items-start justify-between gap-4 flex-wrap">
		<div>
			<div class="text-xs uppercase tracking-widest text-neutral-400 font-mono">Built for</div>
			<h3 class="font-sans font-normal text-base sm:text-lg text-neutral-900 mt-1">
				Experiments, hobby software, and side projects
			</h3>
		</div>

		<!-- Live Hosted Badge -->
		<div class="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-neutral-100 border border-neutral-200/60 text-xs text-neutral-700 font-normal">
			<span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
			<span>{projectsOnline} of {projectsTotal} projects online</span>
		</div>
	</div>

	<!-- Use cases tag pills -->
	<div class="flex items-center gap-2 flex-wrap pt-3">
		{#each useCases as tag}
			<span class="px-2.5 py-1 rounded-lg text-xs font-normal bg-neutral-50 text-neutral-600 border border-neutral-200/60 hover:bg-neutral-100 transition-colors">
				{tag}
			</span>
		{/each}
	</div>
</div>
