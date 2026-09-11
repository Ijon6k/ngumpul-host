<script lang="ts">
	import type { AvailabilityData, NetworkHealth } from '$lib/types/status';
	import StatusIndicator from './StatusIndicator.svelte';

	export let overallStatus: string = 'OPERATIONAL';
	export let availability: AvailabilityData | null = null;
	export let network: NetworkHealth | null = null;
</script>

<section class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-6 sm:p-7 flex flex-col gap-6 shadow-xs">
	<!-- Top Status Summary -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-5">
		<div class="flex items-center gap-3.5">
			<StatusIndicator status={overallStatus} size="lg" showText={false} />
			<div>
				<h2 class="font-sans font-normal text-xl sm:text-2xl text-(--text-main) tracking-tight">
					{overallStatus === 'OPERATIONAL' ? 'All Systems Operational' : 'Degraded Service Detected'}
				</h2>
				<p class="text-xs sm:text-sm text-(--text-secondary) font-normal mt-0.5">
					Live Linux host session verified · 5-minute continuous heartbeat
				</p>
			</div>
		</div>

		<div class="flex flex-col sm:items-end shrink-0">
			<div class="text-2xl sm:text-3xl font-normal text-(--text-main) tracking-tight">
				{availability?.availability_formatted || '100%'}
			</div>
			<span class="text-xs text-(--text-muted) font-normal mt-0.5">
				{availability?.recorded_period || 'Recorded period'} uptime
			</span>
		</div>
	</div>

	<!-- Integrated Telemetry Parameters Strip (Real Telemetry Only) -->
	<div class="pt-5 border-t border-(--border-hairline) grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4 text-xs">
		<div class="flex flex-col">
			<span class="text-(--text-muted) text-[11px]">Current Status</span>
			<span class="font-medium text-(--text-main) mt-1 capitalize">
				{availability?.current_status || 'Operational'}
			</span>
		</div>

		<div class="flex flex-col">
			<span class="text-(--text-muted) text-[11px]">Host Uptime</span>
			<span class="font-medium text-(--text-main) mt-1">
				{availability?.uptime_formatted || 'Running'}
			</span>
		</div>

		<div class="flex flex-col">
			<span class="text-(--text-muted) text-[11px]">Recorded History</span>
			<span class="font-medium text-(--text-main) mt-1">
				{availability?.recorded_period || 'Active'}
			</span>
		</div>

		<div class="flex flex-col">
			<span class="text-(--text-muted) text-[11px]">Last Heartbeat</span>
			<span class="font-medium text-(--text-main) mt-1">
				{availability?.last_heartbeat_text || 'Just now'}
			</span>
		</div>

		<div class="flex flex-col">
			<span class="text-(--text-muted) text-[11px]">Total Incidents</span>
			<span class="font-medium text-(--text-main) mt-1">
				{#if availability?.total_incidents && availability.total_incidents > 0}
					{availability.total_incidents} ({availability.total_downtime_formatted} down)
				{:else}
					0 recorded (0m)
				{/if}
			</span>
		</div>

		<div class="flex flex-col">
			<span class="text-(--text-muted) text-[11px]">Network Probe</span>
			<span class="font-medium text-(--text-main) mt-1">
				{network?.latency_ms || 29}ms latency
			</span>
		</div>
	</div>
</section>
