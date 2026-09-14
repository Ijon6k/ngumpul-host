<script lang="ts">
	import { cn } from '$lib/utils';
	import { CaretDown } from 'phosphor-svelte';

	export let value: string | number = '';
	export let label: string = '';
	export let helperText: string = '';
	export let error: string = '';
	export let disabled: boolean = false;
	export let id: string = '';

	let className: string = '';
	export { className as class };

	$: selectId = id || (label ? `select-${label.toLowerCase().replace(/\s+/g, '-')}` : undefined);
</script>

<div class={cn('flex flex-col gap-1.5 w-full', className)}>
	{#if label}
		<label for={selectId} class="text-xs font-medium text-(--text-secondary) select-none">
			{label}
		</label>
	{/if}

	<div class="relative flex items-center w-full">
		<select
			id={selectId}
			{disabled}
			bind:value
			class={cn(
				'w-full h-9 pl-3 pr-8 text-xs sm:text-sm bg-(--bg-surface) text-(--text-main) border rounded-[4px] outline-none transition-colors duration-150 appearance-none cursor-pointer',
				error
					? 'border-rose-500 focus:border-rose-500 focus:ring-1 focus:ring-rose-500'
					: 'border-(--border-hairline) focus:border-(--accent-sky) focus:ring-1 focus:ring-(--accent-sky)',
				disabled && 'opacity-50 cursor-not-allowed bg-(--bg-muted)'
			)}
			on:change
			on:input
			on:focus
			on:blur
			{...$$restProps}
		>
			<slot />
		</select>

		<div class="absolute right-2.5 pointer-events-none text-(--text-muted)">
			<CaretDown size={14} weight="bold" />
		</div>
	</div>

	{#if error}
		<p class="text-xs text-rose-600 dark:text-rose-400 font-normal">
			{error}
		</p>
	{:else if helperText}
		<p class="text-xs text-(--text-muted) font-normal">
			{helperText}
		</p>
	{/if}
</div>
