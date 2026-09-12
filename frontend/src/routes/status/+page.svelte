<script lang="ts">
	import { writable } from 'svelte/store';
	import { createStatusQuery } from '$lib/api/status';
	import StatusHeader from '$lib/components/status/StatusHeader.svelte';
	import StatusOverviewCard from '$lib/components/status/StatusOverviewCard.svelte';
	import AvailabilityGrid from '$lib/components/status/AvailabilityGrid.svelte';
	import ServiceHealthCard from '$lib/components/status/ServiceHealthCard.svelte';
	import IncidentHistoryCard from '$lib/components/status/IncidentHistoryCard.svelte';
	import StatusIndicator from '$lib/components/status/StatusIndicator.svelte';

	const selectedRange = writable<'1d' | '7d' | '30d'>('30d');

	// TanStack Query store hook with automated 60s background polling & query caching
	const statusQuery = createStatusQuery(selectedRange);

	$: timeUnitLabel =
		$selectedRange === '1d'
			? '24 cells · 1 cell / hour'
			: $selectedRange === '7d'
			? '7 cells · 1 cell / day'
			: '30 cells · 1 cell / day';

	function handleRangeChange(range: '1d' | '7d' | '30d') {
		$selectedRange = range;
	}
</script>

<svelte:head>
	<title>System Status · Ngumpul Host</title>
	<meta name="description" content="Data-driven availability history and measured service latency for Ngumpul Host infrastructure." />
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-12 pb-24 flex flex-col gap-8 sm:gap-9">

	<!-- 1. Header & Range Controls -->
	<StatusHeader
		selectedRange={$selectedRange}
		onRangeChange={handleRangeChange}
	/>

		<!-- 2. Loading State (when no cache exists) -->
		{#if $statusQuery.isPending && !$statusQuery.data}
			<div class="flex flex-col gap-6 animate-pulse" aria-busy="true">
				<div class="h-32 bg-(--bg-muted) rounded-xl border border-(--border-hairline)"></div>
				<div class="h-36 bg-(--bg-muted) rounded-xl border border-(--border-hairline)"></div>
				<div class="h-44 bg-(--bg-muted) rounded-xl border border-(--border-hairline)"></div>
			</div>

		<!-- 3. Error State -->
		{:else if $statusQuery.isError && !$statusQuery.data}
			<div class="p-8 bg-(--bg-surface) border border-rose-200 dark:border-rose-900/40 rounded-xl flex flex-col items-center justify-center text-center gap-4">
				<StatusIndicator status="DOWN" size="lg" />
				<div class="max-w-md">
					<h2 class="text-base font-medium text-(--text-main)">Telemetry Unavailable</h2>
					<p class="text-xs text-(--text-secondary) mt-1">
						{$statusQuery.error instanceof Error ? $statusQuery.error.message : 'Unable to connect to telemetry probes.'}
					</p>
				</div>
				<button
					type="button"
					class="px-4 py-2 text-xs rounded-lg border border-(--border-hairline) bg-(--bg-muted) hover:bg-(--bg-hover) text-(--text-main) transition-colors cursor-pointer"
					on:click={() => $statusQuery.refetch()}
				>
					Retry Connection
				</button>
			</div>

		<!-- 4. Data-Driven 4-Section Content -->
		{:else if $statusQuery.data}
			<!-- Section 1: Overall Status -->
			<StatusOverviewCard
				overallStatus={$statusQuery.data.overall_status}
				availability={$statusQuery.data.availability}
				network={$statusQuery.data.network}
			/>

			<!-- Section 2: Availability History (Activity Heatmap) -->
			<section class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-6 sm:p-7 flex flex-col gap-5 shadow-xs">
				<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
					<div>
						<h3 class="text-base font-normal text-(--text-main)">
							Availability History
						</h3>
						<p class="text-xs text-(--text-muted) mt-0.5">
							Aggregated real telemetry for {$selectedRange === '1d' ? 'the past 24 hours' : $selectedRange === '7d' ? 'the past 7 days' : 'the past 30 days'}. Gray indicates unmonitored days.
						</p>
					</div>

					<div class="text-xs font-mono text-(--text-secondary) self-start sm:self-auto">
						{#if $statusQuery.data.availability?.has_sufficient_data}
							{$statusQuery.data.availability.availability_formatted} uptime
						{:else}
							Collecting initial observations
						{/if}
					</div>
				</div>

				<AvailabilityGrid
					blocks={$statusQuery.data.availability?.blocks || []}
					{timeUnitLabel}
					compact={false}
				/>
			</section>

			<!-- Section 3: Service Health -->
			<ServiceHealthCard
				services={$statusQuery.data.services || []}
				checkedAt={$statusQuery.data.checked_at}
			/>

			<!-- Section 4: Incident History -->
			<IncidentHistoryCard
				incidents={$statusQuery.data.availability?.incidents || []}
			/>

			<!-- Architectural Footnote -->
			<footer class="text-xs text-(--text-muted) leading-relaxed px-1 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
				<p class="max-w-2xl">
					Telemetry data is sourced directly from Linux kernel sessions (<code class="font-mono text-xs">/proc/uptime</code> and <code class="font-mono text-xs">boot_id</code>). Availability percentages reflect real recorded observation periods only.
				</p>
				<span class="shrink-0 text-right">Ngumpul Host · Bare Metal</span>
			</footer>
		{/if}

</div>

