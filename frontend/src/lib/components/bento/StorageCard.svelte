<script lang="ts">
	export let physicalGB: number = 0;
	export let linuxTotalGB: number = 0;
	export let totalGB: number = 0;
	export let usagePercent: number = 0;
	export let availableGB: number = 0;
	export let type: string = 'NVMe SSD';
	export let model: string = '';
	export let summary: string = '';

	$: displayCapacity = physicalGB > 0 ? physicalGB : (linuxTotalGB || totalGB || 0);
	$: displaySubtext = summary ? summary : `${model ? model + ' · ' : ''}${linuxTotalGB || totalGB} GB Linux Environment`;
</script>

<div class="h-full w-full bg-(--bg-surface) rounded-3xl p-5 sm:p-6 border border-(--border-hairline) flex flex-col justify-between shadow-xs min-h-[175px] transition-colors duration-300">
	<div>
		<div class="flex items-center justify-between">
			<div class="text-base font-normal text-(--text-main)">Storage</div>
			<span class="text-xs text-(--text-muted) font-normal">
				{type}
			</span>
		</div>
		<div class="text-xs text-(--text-secondary) font-normal mt-0.5 truncate" title={displaySubtext}>
			{displaySubtext}
		</div>
		<div class="w-6 h-[1px] bg-(--border-hairline) my-2.5"></div>
	</div>

	<!-- Metric Row -->
	<div class="my-auto py-1">
		<div class="flex items-baseline justify-between">
			<div class="text-3xl sm:text-4xl font-normal text-(--text-main) tracking-tight">
				{displayCapacity} <span class="text-lg text-(--text-muted) font-light">GB</span>
			</div>
			<div class="text-xs text-(--text-secondary) font-normal">
				{availableGB} GB available for projects
			</div>
		</div>

		<!-- Progress Bar -->
		<div class="w-full h-2 rounded-full bg-(--bg-hover) overflow-hidden mt-2 p-[1px]">
			<div
				class="h-full rounded-full bg-(--accent-sky) transition-all duration-700"
				style="width: {usagePercent}%"
			></div>
		</div>
	</div>

	<!-- Footer note -->
	<div class="pt-2 border-t border-(--border-hairline) flex items-center justify-between text-xs text-(--text-muted) font-normal">
		<span>Fast flash disk</span>
		<span>Dedicated for apps & DB</span>
	</div>
</div>
