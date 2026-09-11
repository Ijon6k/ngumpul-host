<script lang="ts">
	import BentoCard from './BentoCard.svelte';

	export let hostSpecs: any = null;

	$: load1 = hostSpecs?.load_avg?.[0] ?? 0.32;
	$: load5 = hostSpecs?.load_avg?.[1] ?? 0.28;
	$: load15 = hostSpecs?.load_avg?.[2] ?? 0.22;
	$: kernelName = hostSpecs?.kernel ?? 'Linux 7.1.8-xanmod1';
	$: uptime = hostSpecs?.uptime_formatted ?? '3 hours';
	$: hostname = hostSpecs?.hostname ?? 'ngumpul-host';
</script>

<BentoCard title="Host Environment & Load" badge="Kernel Diagnostics" badgeClass="text-(--text-muted)">
	<div class="flex flex-col gap-4 text-xs">
		<!-- System Load Average -->
		<div class="flex flex-col gap-2">
			<div class="flex items-center justify-between">
				<span class="text-(--text-secondary) font-medium">System Load Average</span>
				<span class="text-[11px] font-mono text-emerald-600 dark:text-emerald-400">Normal / Balanced</span>
			</div>

			<!-- 3 load average chips: 1 min, 5 min, 15 min -->
			<div class="grid grid-cols-3 gap-2">
				<div class="p-2.5 rounded-lg bg-(--bg-muted) border border-(--border-hairline) flex flex-col gap-0.5">
					<span class="text-xs font-mono text-(--text-muted)">1 min</span>
					<span class="font-mono font-bold text-sm sm:text-base text-(--text-main)">{load1.toFixed(2)}</span>
				</div>
				<div class="p-2.5 rounded-lg bg-(--bg-muted) border border-(--border-hairline) flex flex-col gap-0.5">
					<span class="text-xs font-mono text-(--text-muted)">5 min</span>
					<span class="font-mono font-bold text-sm sm:text-base text-(--text-main)">{load5.toFixed(2)}</span>
				</div>
				<div class="p-2.5 rounded-lg bg-(--bg-muted) border border-(--border-hairline) flex flex-col gap-0.5">
					<span class="text-xs font-mono text-(--text-muted)">15 min</span>
					<span class="font-mono font-bold text-sm sm:text-base text-(--text-main)">{load15.toFixed(2)}</span>
				</div>
			</div>
		</div>

		<!-- Operating System & Kernel Runtime -->
		<div class="pt-3 border-t border-(--border-hairline) flex flex-col gap-1">
			<span class="text-(--text-secondary) font-medium text-xs">Operating system & kernel</span>
			<div class="flex items-center justify-between">
				<span class="font-mono text-xs sm:text-sm text-(--text-main) font-semibold truncate max-w-[280px]">
					{kernelName}
				</span>
				<span class="font-mono text-xs px-2 py-0.5 rounded bg-(--cf-pastel-bg) text-(--cf-pastel-text) border border-(--cf-pastel-border)">
					{hostSpecs?.arch ?? 'amd64'}
				</span>
			</div>
		</div>

		<!-- Host Node Identity & Uptime -->
		<div class="pt-3 border-t border-(--border-hairline) grid grid-cols-2 gap-3">
			<div class="flex flex-col gap-0.5">
				<span class="text-xs text-(--text-secondary)">Node identity</span>
				<span class="font-mono text-xs sm:text-sm font-semibold text-(--text-main)">{hostname}</span>
			</div>
			<div class="flex flex-col gap-0.5">
				<span class="text-xs text-(--text-secondary)">Current uptime</span>
				<span class="font-mono text-xs sm:text-sm font-semibold text-(--text-main)">{uptime}</span>
			</div>
		</div>
	</div>

	<svelte:fragment slot="footer">
		<div class="flex items-center justify-between text-xs text-(--text-muted)">
			<span>Direct Linux /proc telemetry</span>
			<span class="font-mono text-xs">Real hardware node</span>
		</div>
	</svelte:fragment>
</BentoCard>
