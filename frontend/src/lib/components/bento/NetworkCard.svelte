<script lang="ts">
	export let linkCapacity: string = '1 Gbps';
	export let latencyMs: number = 29;
</script>

<div class="h-full w-full bg-(--bg-surface) rounded-3xl p-5 sm:p-6 border border-(--border-hairline) flex flex-col justify-between shadow-xs min-h-[175px] transition-colors duration-300">
	<div>
		<div class="flex items-center justify-between">
			<div class="text-base font-normal text-(--text-main)">Network</div>
			<span class="text-xs text-(--text-muted) font-normal">
				Direct connection
			</span>
		</div>
		<div class="text-xs text-(--text-secondary) font-normal mt-0.5">Local interface link capacity</div>
		<div class="w-6 h-[1px] bg-(--border-hairline) my-2.5"></div>
	</div>

	<!-- Primary Metric -->
	<div class="my-auto py-1 flex items-baseline justify-between">
		<div>
			<div class="text-2xl sm:text-3xl font-normal text-(--text-main) tracking-tight">
				{linkCapacity}
			</div>
			<div class="text-xs text-(--text-muted) font-normal mt-0.5">Interface capacity</div>
		</div>

		<!-- Latency Pillar -->
		<div class="text-right">
			<div class="text-lg sm:text-xl font-normal text-(--accent-sky)">
				{latencyMs} <span class="text-xs text-(--text-muted)">ms</span>
			</div>
			<div class="text-xs text-(--text-muted) font-normal">External response</div>
		</div>
	</div>

	<!-- Footer -->
	<div class="pt-2 border-t border-(--border-hairline) text-xs text-(--text-muted) font-normal flex items-center gap-1.5">
		<span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
		<span>Low-latency routing from home network</span>
	</div>
</div>
