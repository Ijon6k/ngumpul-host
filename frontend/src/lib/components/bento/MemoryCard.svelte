<script lang="ts">
	export let totalGB: number = 23.1;
	export let usagePercent: number = 68;
	export let availableGB: number = 7.3;

	// Dynamic visual allocation bars
	const barThresholds = [20, 35, 50, 60, 70, 80, 90, 95];
</script>

<div class="h-full w-full bg-[#0E1013] dark:bg-(--bg-surface) text-white rounded-xl p-5 sm:p-6 flex flex-col justify-between relative shadow-xs min-h-[240px] border border-neutral-800/80 dark:border-(--border-hairline) transition-colors duration-300">
	<!-- Top: Title & Equalizer Bars -->
	<div class="flex items-start justify-between">
		<div>
			<div class="text-base font-normal text-white">Memory</div>
			<div class="text-xs text-neutral-400 font-normal mt-0.5">System RAM</div>
			<div class="w-6 h-[1px] bg-neutral-700 my-2.5"></div>
		</div>

		<!-- Minimal Equalizer Bar Visualization -->
		<div class="flex items-end gap-1 h-9 pr-1">
			{#each barThresholds as threshold, idx}
				{@const isFilled = usagePercent >= threshold}
				<div class="w-1 rounded-full bg-neutral-800 h-full flex items-end">
					<div
						class="w-full rounded-full transition-all duration-500 {isFilled ? (usagePercent > 85 ? 'bg-amber-400' : 'bg-sky-400') : 'bg-transparent'}"
						style="height: {Math.min(100, Math.max(15, (threshold / 100) * 100))}%"
					></div>
				</div>
			{/each}
		</div>
	</div>

	<!-- Center: Metric Display -->
	<div class="my-auto py-2">
		<div class="text-3xl sm:text-4xl font-normal text-white tracking-tight">
			{totalGB} <span class="text-lg text-neutral-400 font-light">GB</span>
		</div>
		<div class="text-xs font-mono text-neutral-300 mt-1">
			{Math.round(usagePercent)}% currently in use
		</div>
	</div>

	<!-- Footer: Real Available Memory -->
	<div class="pt-3 border-t border-neutral-800/80 flex items-center justify-between text-xs text-neutral-400">
		<span>Available</span>
		<span class="font-mono text-neutral-200 font-normal">{availableGB} GB free</span>
	</div>
</div>
