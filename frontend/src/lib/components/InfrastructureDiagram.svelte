<script lang="ts">
	export let systemStatus = 'OPERATIONAL';
	$: isHealthy = systemStatus === 'OPERATIONAL';
	$: nodeStatus = isHealthy ? 'Online' : 'Degraded';
	$: dbStatus = isHealthy ? 'Healthy' : 'Degraded';
	$: nodeColor = isHealthy ? 'var(--color-success)' : 'var(--color-warning)';
</script>

<div class="w-full bg-(--bg-surface) border border-(--border-hairline) rounded-md p-6 sm:p-7 flex flex-col gap-5">
	<div class="flex items-center justify-between flex-wrap gap-2 pb-4 border-b border-(--border-hairline)">
		<div class="flex items-center gap-2.5">
			<span class="w-2 h-2 rounded-full bg-(--color-success)"></span>
			<span class="font-display font-semibold text-sm text-(--text-main)">Infrastructure Topology</span>
		</div>
		<span class="text-xs text-(--text-muted)">Port 1111 Isolation · Single Public Entrypoint</span>
	</div>

	<!-- Flow Diagram with Nodes -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 text-xs">
		<!-- Node 1: Edge & Ingress -->
		<div class="flex flex-col justify-between p-4 min-h-[115px] bg-(--bg-muted) border border-(--border-hairline) rounded-md gap-3">
			<div class="flex flex-col gap-1">
				<span class="text-[11px] text-(--text-muted)">01 · Edge Ingress</span>
				<span class="font-semibold text-sm text-(--text-main)">Nginx Gateway</span>
				<span class="font-mono text-xs text-(--accent-strong)">:1111</span>
			</div>
			<div class="pt-2 text-[11px] text-(--text-muted) border-t border-(--border-hairline) flex items-center justify-between">
				<span>Proxy & SSL</span>
				<span class="font-medium" style="color: {nodeColor}">{nodeStatus}</span>
			</div>
		</div>

		<!-- Node 2: App Tier (SvelteKit) -->
		<div class="flex flex-col justify-between p-4 min-h-[115px] bg-(--bg-muted) border border-(--border-hairline) rounded-md gap-3">
			<div class="flex flex-col gap-1">
				<span class="text-[11px] text-(--text-muted)">02 · Web Tier</span>
				<span class="font-semibold text-sm text-(--text-main)">SvelteKit Node</span>
				<span class="font-mono text-xs text-(--accent-strong)">:3000</span>
			</div>
			<div class="pt-2 text-[11px] text-(--text-muted) border-t border-(--border-hairline) flex items-center justify-between">
				<span>SSR & Hydration</span>
				<span class="font-medium" style="color: {nodeColor}">{nodeStatus}</span>
			</div>
		</div>

		<!-- Node 3: Core API (Go) -->
		<div class="flex flex-col justify-between p-4 min-h-[115px] bg-(--bg-muted) border border-(--border-hairline) rounded-md gap-3">
			<div class="flex flex-col gap-1">
				<span class="text-[11px] text-(--text-muted)">03 · Core Engine</span>
				<span class="font-semibold text-sm text-(--text-main)">Go Chi Monolith</span>
				<span class="font-mono text-xs text-(--accent-strong)">:8080</span>
			</div>
			<div class="pt-2 text-[11px] text-(--text-muted) border-t border-(--border-hairline) flex items-center justify-between">
				<span>REST /api/v1</span>
				<span class="font-medium" style="color: {nodeColor}">{nodeStatus}</span>
			</div>
		</div>

		<!-- Node 4: Persistence (Postgres) -->
		<div class="flex flex-col justify-between p-4 min-h-[115px] bg-(--bg-muted) border border-(--border-hairline) rounded-md gap-3">
			<div class="flex flex-col gap-1">
				<span class="text-[11px] text-(--text-muted)">04 · Database</span>
				<span class="font-semibold text-sm text-(--text-main)">PostgreSQL 16</span>
				<span class="font-mono text-xs text-(--accent-strong)">:5432</span>
			</div>
			<div class="pt-2 text-[11px] text-(--text-muted) border-t border-(--border-hairline) flex items-center justify-between">
				<span>ACID Persistence</span>
				<span class="font-medium" style="color: {nodeColor}">{dbStatus}</span>
			</div>
		</div>
	</div>
</div>
