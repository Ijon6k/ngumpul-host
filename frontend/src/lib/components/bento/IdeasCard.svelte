<script lang="ts">
	import { Sparkle, ArrowRight } from 'phosphor-svelte';

	export let projectsOnline: number = 3;
	export let projectsTotal: number = 3;
</script>

<div class="h-full w-full bg-[#406882] dark:bg-[#1E3A4C] text-white rounded-3xl p-5 sm:p-6 flex flex-col justify-between relative overflow-hidden shadow-xs border border-white/10 dark:border-(--border-hairline) min-h-[175px] transition-colors duration-300">
	<!-- Background subtle decorative circles -->
	<div class="absolute -right-6 -bottom-6 w-32 h-32 rounded-full border border-white/10 pointer-events-none"></div>

	<!-- Top Glyph & Action Button -->
	<div class="relative z-10 flex items-center justify-between">
		<Sparkle size={22} weight="regular" class="text-white opacity-90" />

		<a
			href="/register"
			class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-white text-neutral-900 text-xs font-medium hover:bg-neutral-100 transition-all active:scale-[0.98] shadow-xs"
		>
			<span>Request a spot</span>
			<ArrowRight size={12} weight="bold" />
		</a>
	</div>

	<!-- Main Headline -->
	<div class="relative z-10 my-auto py-1">
		<h3 class="font-sans font-normal text-lg sm:text-xl text-white leading-snug tracking-tight">
			A small server for bigger ideas.
		</h3>
		<p class="text-xs text-white/85 font-normal mt-1 leading-relaxed">
			Side projects, bots, and experiments—kept online without cloud bills.
		</p>
	</div>

	<!-- Footer Info -->
	<div class="relative z-10 pt-2 border-t border-white/15 flex items-center justify-between text-xs text-white/80 font-normal">
		<span>{projectsOnline} of {projectsTotal} projects online</span>
		<span>3-step setup</span>
	</div>
</div>
