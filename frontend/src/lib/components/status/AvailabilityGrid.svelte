<script lang="ts">
	import type { AvailabilityBlock } from '$lib/types/status';

	export let blocks: AvailabilityBlock[] = [];
	export let compact: boolean = false;
	export let dark: boolean = false;
	export let timeUnitLabel: string = '1 cell / day';

	let activeHover: AvailabilityBlock | null = null;
	let hoverPos = { x: 0, y: 0 };

	function getBlockColor(block: AvailabilityBlock): string {
		switch (block.status) {
			case 'operational':
				// Green for fully available
				return dark
					? 'bg-emerald-500 hover:bg-emerald-400'
					: 'bg-[#2E8555] hover:bg-emerald-700';
			case 'partial':
				// Muted lighter green for partial availability
				return dark
					? 'bg-emerald-400/50 hover:bg-emerald-400/70 border border-emerald-400/30'
					: 'bg-emerald-200 hover:bg-emerald-300 border border-emerald-300';
			case 'incident':
				// Restrained red only for actual detected/estimated downtime
				return dark
					? 'bg-rose-500 hover:bg-rose-400'
					: 'bg-rose-600 hover:bg-rose-700';
			case 'no_data':
			default:
				// Neutral gray for no recorded data
				return dark
					? 'bg-neutral-800 hover:bg-neutral-700'
					: 'bg-neutral-200 dark:bg-neutral-800 hover:bg-neutral-300 dark:hover:bg-neutral-700';
		}
	}

	function handleMouseEnter(e: MouseEvent, block: AvailabilityBlock) {
		activeHover = block;
		const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
		hoverPos = {
			x: rect.left + rect.width / 2,
			y: rect.top - 8
		};
	}

	function handleMouseLeave() {
		activeHover = null;
	}
</script>

<div class="flex flex-col gap-3 w-full select-none">
	<!-- High-resolution Vertical Bar Sequence -->
	<div
		class="flex items-center gap-[2px] sm:gap-[3px] w-full py-1"
		role="region"
		aria-label="Availability activity heatmap"
	>
		{#each blocks as block (block.index)}
			<div
				class="flex-1 min-w-[2px] rounded-[1.5px] transition-all duration-150 cursor-pointer {getBlockColor(block)} {compact ? 'h-4 sm:h-5' : 'h-6 sm:h-8'}"
				role="button"
				tabindex="0"
				aria-label="{block.label}: {block.details}"
				on:mouseenter={(e) => handleMouseEnter(e, block)}
				on:mouseleave={handleMouseLeave}
			></div>
		{/each}
	</div>

	<!-- Time Bounds & Resolution Indicator -->
	{#if blocks.length > 0}
		<div class="flex items-center justify-between text-xs {dark ? 'text-neutral-400' : 'text-(--text-muted)'} font-normal select-none">
			<span>{blocks[0]?.label || ''}</span>
			<span class="text-xs tracking-wider uppercase opacity-75">{timeUnitLabel}</span>
			<span>{blocks[blocks.length - 1]?.label || ''}</span>
		</div>
	{/if}

	<!-- Legend (shown when not compact) -->
	{#if !compact}
		<div class="flex items-center justify-between pt-3 border-t {dark ? 'border-neutral-800' : 'border-(--border-hairline)'} text-xs {dark ? 'text-neutral-400' : 'text-(--text-secondary)'} flex-wrap gap-3">
			<div class="flex items-center gap-3.5 sm:gap-5 flex-wrap">
				<div class="flex items-center gap-1.5">
					<span class="w-2.5 h-2.5 rounded-[2px] {dark ? 'bg-neutral-800' : 'bg-neutral-200 dark:bg-neutral-800'}"></span>
					<span class="text-xs">No recorded data</span>
				</div>
				<div class="flex items-center gap-1.5">
					<span class="w-2.5 h-2.5 rounded-[2px] {dark ? 'bg-emerald-500' : 'bg-[#2E8555]'}"></span>
					<span class="text-xs">Fully available</span>
				</div>
				<div class="flex items-center gap-1.5">
					<span class="w-2.5 h-2.5 rounded-[2px] {dark ? 'bg-emerald-400/50 border border-emerald-400/30' : 'bg-emerald-200 border border-emerald-300'}"></span>
					<span class="text-xs">Partial availability</span>
				</div>
				<div class="flex items-center gap-1.5">
					<span class="w-2.5 h-2.5 rounded-[2px] {dark ? 'bg-rose-500' : 'bg-rose-600'}"></span>
					<span class="text-xs">Downtime</span>
				</div>
			</div>

			<span class="text-xs {dark ? 'text-neutral-500' : 'text-(--text-muted)'}">
				5-min session heartbeat
			</span>
		</div>
	{/if}

	<!-- Floating Calm Tooltip -->
	{#if activeHover}
		<div
			class="fixed z-50 pointer-events-none transform -translate-x-1/2 -translate-y-full px-3 py-1.5 rounded-md text-xs shadow-md border backdrop-blur-md transition-opacity duration-100 {dark ? 'bg-neutral-900/95 text-white border-neutral-700' : 'bg-neutral-900/95 text-white border-neutral-800'}"
			style="left: {hoverPos.x}px; top: {hoverPos.y}px;"
		>
			<div class="font-medium text-neutral-100">{activeHover.label}</div>
			<div class="text-xs text-neutral-300 mt-0.5">
				{#if activeHover.status === 'no_data'}
					No telemetry recorded
				{:else}
					{activeHover.details}
				{/if}
			</div>
		</div>
	{/if}
</div>
