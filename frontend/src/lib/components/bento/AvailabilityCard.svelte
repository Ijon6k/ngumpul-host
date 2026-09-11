<script lang="ts">
	import type { AvailabilityBlock } from '$lib/types/status';
	import { ArrowRight } from 'phosphor-svelte';

	export let availabilityPercent: number | null = 100;
	export let status: string = 'operational';
	export let uptime: string = 'Running';
	export let recordedPeriod: string = 'Recorded';
	export let dailyBlocks: AvailabilityBlock[] = [];

	$: displayBlocks = dailyBlocks && dailyBlocks.length > 0
		? dailyBlocks
		: Array.from({ length: 30 }, (_, i) => ({
				index: i,
				label: `Day ${i + 1}`,
				date: `Day ${i + 1}`,
				status: (i === 29 ? 'operational' : 'no_data') as AvailabilityBlock['status'],
				uptimePercent: i === 29 ? 100 : null,
				details: i === 29 ? '100% operational' : 'No recorded telemetry'
		  }));

	function getBlockClass(item: AvailabilityBlock): string {
		switch (item.status) {
			case 'operational':
				return 'bg-emerald-400/90 hover:bg-emerald-300';
			case 'partial':
				return 'bg-amber-400/90 hover:bg-amber-300';
			case 'incident':
				return 'bg-rose-400/90 hover:bg-rose-300';
			case 'no_data':
			default:
				return 'bg-neutral-800/90 hover:bg-neutral-700';
		}
	}
</script>

<a
	href="/status"
	class="group block h-full w-full bg-[#0E1013] dark:bg-(--bg-surface) text-white rounded-3xl p-5 sm:p-6 flex flex-col justify-between relative shadow-xs min-h-[175px] border border-neutral-800/80 dark:border-(--border-hairline) hover:border-neutral-700 dark:hover:border-neutral-500 transition-colors duration-300"
>
	<div>
		<div class="flex items-center justify-between">
			<div class="text-base font-normal text-white">Availability</div>
			<div class="flex items-center gap-1.5 text-xs font-normal">
				<span class="w-1.5 h-1.5 rounded-full {status === 'operational' || status === 'online' ? 'bg-emerald-400' : status === 'degraded' ? 'bg-amber-400' : 'bg-rose-400'} animate-pulse"></span>
				<span class="capitalize text-neutral-300">{status}</span>
			</div>
		</div>
		<div class="text-xs text-neutral-400 font-normal mt-0.5">Continuous host session monitor</div>
		<div class="w-6 h-[1px] bg-neutral-700 my-2.5"></div>
	</div>

	<!-- Primary Metric & 30-day daily blocks -->
	<div class="my-auto py-1">
		<div class="flex items-baseline justify-between mb-2.5">
			<div class="text-2xl sm:text-3xl font-normal text-white tracking-tight">
				{#if availabilityPercent !== null}
					{availabilityPercent}%
				{:else}
					100%
				{/if}
			</div>
			<div class="text-xs text-neutral-400 font-normal">
				{recordedPeriod} · {uptime}
			</div>
		</div>

		<!-- 30-day compact GitHub-style daily blocks -->
		<div class="grid grid-cols-15 sm:grid-cols-30 gap-1 py-1">
			{#each displayBlocks as item (item.index)}
				<div
					class="h-2 rounded-xs {getBlockClass(item)} transition-all"
					title="{item.label}: {item.status === 'no_data' ? 'No data recorded' : item.details}"
				></div>
			{/each}
		</div>
	</div>

	<!-- Footer with link to status page -->
	<div class="pt-2 border-t border-neutral-800/80 flex items-center justify-between text-xs text-neutral-400 font-normal">
		<span>Zero heartbeat spam · /proc telemetry</span>
		<span class="text-neutral-500 group-hover:text-neutral-300 transition-colors inline-flex items-center gap-1">
			<span>Status page</span>
			<ArrowRight size={11} weight="bold" />
		</span>
	</div>
</a>
