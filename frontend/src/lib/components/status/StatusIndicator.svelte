<script lang="ts">
	export let status: string = 'OPERATIONAL';
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let showText: boolean = true;
	export let customLabel: string = '';

	$: normalized = (status || '').toLowerCase();

	$: isOperational = normalized === 'operational' || normalized === 'online' || normalized === 'up';
	$: isDegraded = normalized === 'degraded' || normalized === 'partial' || normalized === 'warning';
	$: isDown = normalized === 'down' || normalized === 'incident' || normalized === 'offline';
	$: isNoData = normalized === 'no_data' || normalized === 'unknown';

	$: dotClass = isOperational
		? 'bg-emerald-500'
		: isDegraded
		? 'bg-amber-500'
		: isDown
		? 'bg-rose-500'
		: 'bg-neutral-400 dark:bg-neutral-600';

	$: textClass = isOperational
		? 'text-emerald-700 dark:text-emerald-400'
		: isDegraded
		? 'text-amber-700 dark:text-amber-400'
		: isDown
		? 'text-rose-700 dark:text-rose-400'
		: 'text-neutral-500 dark:text-neutral-400';

	$: label = customLabel
		? customLabel
		: isOperational
		? 'Operational'
		: isDegraded
		? 'Degraded'
		: isDown
		? 'Outage'
		: 'No data';

	$: dotSizeClass = size === 'sm' ? 'w-1.5 h-1.5' : size === 'lg' ? 'w-2.5 h-2.5' : 'w-2 h-2';
	$: textSizeClass = size === 'sm' ? 'text-xs' : size === 'lg' ? 'text-sm' : 'text-xs sm:text-sm';
</script>

<div class="inline-flex items-center gap-2 font-normal select-none">
	<span class="relative flex items-center justify-center shrink-0">
		{#if isOperational}
			<span class="absolute inline-flex h-full w-full rounded-full {dotClass} opacity-40 animate-ping"></span>
		{/if}
		<span class="relative inline-flex rounded-full {dotSizeClass} {dotClass}"></span>
	</span>

	{#if showText}
		<span class="font-normal {textSizeClass} {textClass} tracking-normal">
			{label}
		</span>
	{/if}
</div>
