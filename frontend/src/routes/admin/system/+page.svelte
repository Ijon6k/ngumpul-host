<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	let systemData: any = null;
	let loading = true;
	let refreshing = false;

	async function loadSystemStats() {
		try {
			const res = await api.get('/admin/system');
			systemData = res.data;
		} catch (err) {
			console.error('Failed to load system diagnostics:', err);
		} finally {
			loading = false;
			refreshing = false;
		}
	}

	onMount(loadSystemStats);

	function refresh() {
		refreshing = true;
		loadSystemStats();
	}

	function formatUptime(seconds?: number) {
		if (!seconds) return '—';
		const hrs = Math.floor(seconds / 3600);
		const mins = Math.floor((seconds % 3600) / 60);
		const secs = seconds % 60;
		return `${hrs}h ${mins}m ${secs}s`;
	}
</script>

<div class="flex flex-col gap-8">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-(--border-hairline)">
		<div>
			<h1 class="font-display font-bold text-2xl text-(--text-main) tracking-tight">System Diagnostics</h1>
			<p class="text-sm text-(--text-secondary) mt-1">
				Real-time runtime memory allocation, database connection pools, and container isolation.
			</p>
		</div>

		<div>
			<button
				type="button"
				class="btn btn-secondary btn-sm text-xs sm:text-sm font-medium"
				on:click={refresh}
				disabled={refreshing}
			>
				{refreshing ? 'Probing...' : 'Refresh telemetry'}
			</button>
		</div>
	</div>

	{#if loading && !systemData}
		<div class="p-16 text-center text-sm text-(--text-muted)">Querying runtime diagnostics...</div>
	{:else if systemData}
		<!-- Status Strip (Linear Cards) -->
		<div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
			<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-5 flex flex-col gap-1.5">
				<span class="text-sm text-(--text-secondary) font-medium">Go Service Runtime</span>
				<div class="flex items-center gap-2 mt-1">
					<span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
					<span class="font-display text-2xl font-bold text-emerald-600 dark:text-emerald-400">
						{systemData.application?.status || 'Healthy'}
					</span>
				</div>
				<span class="text-sm text-(--text-muted) mt-2">
					Uptime: {formatUptime(systemData.application?.uptime_seconds)}
				</span>
			</div>

			<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-5 flex flex-col gap-1.5">
				<span class="text-sm text-(--text-secondary) font-medium">PostgreSQL Engine</span>
				<div class="flex items-center gap-2 mt-1">
					<span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
					<span class="font-display text-2xl font-bold text-emerald-600 dark:text-emerald-400">
						{systemData.database?.status || 'Healthy'}
					</span>
				</div>
				<span class="text-sm text-(--text-muted) mt-2">
					Active Pool: {systemData.database?.acquired_conns || 0} / {systemData.database?.max_conns || 25}
				</span>
			</div>

			<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-5 flex flex-col gap-1.5">
				<span class="text-sm text-(--text-secondary) font-medium">Nginx Ingress Proxy</span>
				<div class="flex items-center gap-2 mt-1">
					<span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
					<span class="font-display text-2xl font-bold text-emerald-600 dark:text-emerald-400">Port :1111</span>
				</div>
				<span class="text-sm text-(--text-muted) mt-2">
					Docker Bridge Isolation Active
				</span>
			</div>
		</div>

		<!-- High-Density Telemetry Sections -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Application Service Telemetry -->
			<div class="border border-(--border-hairline) rounded-md bg-(--bg-surface) overflow-hidden flex flex-col">
				<div class="px-5 py-3.5 border-b border-(--border-hairline) flex items-center justify-between">
					<h2 class="font-semibold text-sm sm:text-base text-(--text-main)">
						Go Runtime Telemetry
					</h2>
					<span class="text-xs sm:text-sm font-mono text-(--accent-strong)">PID: Internal</span>
				</div>

				<div class="p-5 flex flex-col gap-3 text-sm">
					<div class="flex items-center justify-between py-2 border-b border-(--border-hairline)">
						<span class="text-(--text-secondary)">Active Goroutines</span>
						<span class="font-bold text-(--text-main)">{systemData.application?.goroutines}</span>
					</div>
					<div class="flex items-center justify-between py-2 border-b border-(--border-hairline)">
						<span class="text-(--text-secondary)">Allocated Heap (Alloc)</span>
						<span class="font-bold text-(--text-main)">{systemData.application?.memory_alloc_mb} MB</span>
					</div>
					<div class="flex items-center justify-between py-2 border-b border-(--border-hairline)">
						<span class="text-(--text-secondary)">System Reserved Memory (Sys)</span>
						<span class="font-bold text-(--text-main)">{systemData.application?.memory_sys_mb} MB</span>
					</div>
					<div class="flex items-center justify-between py-2">
						<span class="text-(--text-secondary)">GC Strategy</span>
						<span class="text-(--text-muted)">Default Go GC (GOGC=100)</span>
					</div>
				</div>
			</div>

			<!-- Database Connection Pool -->
			<div class="border border-(--border-hairline) rounded-md bg-(--bg-surface) overflow-hidden flex flex-col">
				<div class="px-5 py-3.5 border-b border-(--border-hairline) flex items-center justify-between">
					<h2 class="font-semibold text-sm sm:text-base text-(--text-main)">
						PostgreSQL Pool Telemetry
					</h2>
					<span class="text-xs sm:text-sm font-mono text-(--accent-strong)">pgx / db-pool</span>
				</div>

				<div class="p-5 flex flex-col gap-3 text-sm">
					<div class="flex items-center justify-between py-2 border-b border-(--border-hairline)">
						<span class="text-(--text-secondary)">Acquired (Active In-Flight)</span>
						<span class="font-bold text-(--text-main)">{systemData.database?.acquired_conns}</span>
					</div>
					<div class="flex items-center justify-between py-2 border-b border-(--border-hairline)">
						<span class="text-(--text-secondary)">Idle In Pool</span>
						<span class="font-bold text-(--text-main)">{systemData.database?.idle_conns}</span>
					</div>
					<div class="flex items-center justify-between py-2 border-b border-(--border-hairline)">
						<span class="text-(--text-secondary)">Total Pool Allocated</span>
						<span class="font-bold text-(--text-main)">{systemData.database?.total_conns}</span>
					</div>
					<div class="flex items-center justify-between py-2">
						<span class="text-(--text-secondary)">Configured Max Pool</span>
						<span class="font-bold text-(--text-main)">{systemData.database?.max_conns}</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Container Environment Diagnostics -->
		<div class="border border-(--border-hairline) rounded-md bg-(--bg-surface) p-5 flex flex-col gap-4">
			<div class="flex items-center justify-between border-b border-(--border-hairline) pb-3">
				<h3 class="font-semibold text-base text-(--text-main)">Container Network Specification</h3>
				<span class="text-xs sm:text-sm font-mono text-(--text-muted)">Compose Isolation</span>
			</div>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
				<div class="p-3.5 bg-(--bg-muted) rounded-md border border-(--border-hairline) flex flex-col gap-1">
					<span class="text-xs text-(--text-muted)">Proxy Container</span>
					<span class="font-semibold text-sm sm:text-base text-(--text-main)">ngumpul_nginx</span>
					<span class="text-xs sm:text-sm font-mono text-(--text-secondary)">Exposed: 0.0.0.0:1111 → 80</span>
				</div>
				<div class="p-3.5 bg-(--bg-muted) rounded-md border border-(--border-hairline) flex flex-col gap-1">
					<span class="text-xs text-(--text-muted)">Backend Container</span>
					<span class="font-semibold text-sm sm:text-base text-(--text-main)">ngumpul_backend</span>
					<span class="text-xs sm:text-sm font-mono text-(--text-secondary)">Internal: 8080 (No Host Bind)</span>
				</div>
				<div class="p-3.5 bg-(--bg-muted) rounded-md border border-(--border-hairline) flex flex-col gap-1">
					<span class="text-xs text-(--text-muted)">Frontend Container</span>
					<span class="font-semibold text-sm sm:text-base text-(--text-main)">ngumpul_frontend</span>
					<span class="text-xs sm:text-sm font-mono text-(--text-secondary)">Internal: 3000 (No Host Bind)</span>
				</div>
			</div>
		</div>
	{/if}
</div>
