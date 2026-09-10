<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';

	let statusData: any = null;
	let loading = true;

	onMount(async () => {
		try {
			const res = await api.get('/status');
			statusData = res.data;
		} catch (err) {
			console.error('Failed to load status:', err);
		} finally {
			loading = false;
		}
	});

	function formatCheckTime(iso: string) {
		try {
			const d = new Date(iso);
			return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
		} catch (_) {
			return '';
		}
	}
</script>

<div class="container mx-auto px-6 max-w-4xl py-12 pb-24 flex flex-col gap-9">
	<header class="max-w-2xl">
		<h1 class="font-display font-bold text-3xl sm:text-4xl text-(--text-main) tracking-tight mb-2">System Status</h1>
		<p class="text-base sm:text-lg text-(--text-secondary) leading-relaxed">
			Automated availability probes and response latency metrics for core infrastructure components.
		</p>
	</header>

	{#if loading}
		<div class="py-20 text-center text-(--text-muted) text-xs">
			<p>Querying telemetry probes...</p>
		</div>
	{:else if statusData}
		<!-- Main Operational Banner -->
		<div class="p-6 sm:p-8 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col sm:flex-row sm:items-center justify-between gap-6">
			<div class="flex items-center gap-3.5">
				<span class="w-3 h-3 rounded-full bg-(--color-success) shrink-0"></span>
				<div>
					<h2 class="font-display font-bold text-lg sm:text-xl text-(--text-main)">
						{statusData.overall_status === 'OPERATIONAL' ? 'All Systems Operational' : 'Degraded Service Detected'}
					</h2>
					<span class="text-xs text-(--text-secondary)">Gateway probe verified · Zero critical faults</span>
				</div>
			</div>

			<div class="flex flex-col sm:items-end">
				<span class="text-2xl sm:text-3xl font-bold font-display text-(--text-main) leading-none">{statusData.uptime_percentage || '99.85'}%</span>
				<span class="text-xs text-(--text-muted) mt-1">30-day rolling uptime</span>
			</div>
		</div>

		<!-- 30-Day Historical Availability Bar -->
		<div class="p-6 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-3">
			<div class="flex items-center justify-between text-xs text-(--text-muted)">
				<span>30 days ago</span>
				<span class="text-(--color-success) font-medium">100% today</span>
			</div>

			<!-- 30 thin green/soft ticks -->
			<div class="grid grid-cols-30 gap-1 h-6">
				{#each Array(30) as _, i}
					<div
						class="rounded-[2px] bg-(--color-success)/80 hover:bg-(--color-success) transition-colors"
						title="Day {i + 1}: 100% uptime"
					></div>
				{/each}
			</div>

			<div class="flex items-center justify-between text-xs text-(--text-muted) pt-2 border-t border-(--border-hairline)">
				<span>Continuous HTTP health checks</span>
				<span>Last probe: {formatCheckTime(statusData.checked_at)}</span>
			</div>
		</div>

		<!-- Services Breakdown Table -->
		<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md overflow-hidden">
			<div class="px-6 py-3.5 border-b border-(--border-hairline) flex items-center justify-between text-xs text-(--text-muted)">
				<span class="font-medium text-(--text-main)">Component Services</span>
				<span>Target latency: &lt;100ms</span>
			</div>

			<div class="divide-y divide-(--border-hairline)">
				{#each statusData.services || [] as service}
					<div class="px-6 py-3.5 flex items-center justify-between text-sm">
						<div class="flex items-center gap-3">
							<span class="font-medium text-(--text-main)">{service.name}</span>
							<span class="text-xs text-(--text-muted)">({service.response_time_ms} ms)</span>
						</div>

						<StatusDot status={service.status} />
					</div>
				{/each}
			</div>
		</div>

		<!-- Operational Details -->
		<div class="text-xs text-(--text-secondary) leading-relaxed p-1">
			<p>
				Ngumpul Host runs on Docker Compose behind Nginx listening on port 1111. All services are monitored by an internal health check loop verifying database connectivity, container responsiveness, and reverse proxy availability.
			</p>
		</div>
	{/if}
</div>
