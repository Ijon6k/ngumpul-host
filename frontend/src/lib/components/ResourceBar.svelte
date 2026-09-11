<script lang="ts">
	export let label: string;
	export let valueText: string;
	export let percent: number = 0;
	export let colorClass: string = 'bg-(--cf-blue)';
	export let trackClass: string = 'bg-(--bg-muted)';
	export let subText: string = '';

	$: clampedPercent = Math.min(100, Math.max(0, percent || 0));
</script>

<div class="flex flex-col gap-1.5 w-full">
	<div class="flex items-center justify-between text-xs">
		<span class="text-(--text-secondary) font-medium">{label}</span>
		<span class="font-semibold text-(--text-main) font-mono">{valueText}</span>
	</div>

	<div class="w-full h-2 {trackClass} rounded-full overflow-hidden">
		<div
			class="h-full {colorClass} rounded-full transition-all duration-700 ease-out"
			style="width: {clampedPercent}%"
			role="progressbar"
			aria-valuenow={clampedPercent}
			aria-valuemin={0}
			aria-valuemax={100}
			aria-label={label}
		></div>
	</div>

	{#if subText}
		<span class="text-[11px] text-(--text-muted)">{subText}</span>
	{/if}
</div>
