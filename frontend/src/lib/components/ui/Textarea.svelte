<script lang="ts">
	import { cn } from '$lib/utils';

	export let value: string = '';
	export let placeholder: string = '';
	export let label: string = '';
	export let helperText: string = '';
	export let error: string = '';
	export let disabled: boolean = false;
	export let readonly: boolean = false;
	export let rows: number = 3;
	export let id: string = '';

	let className: string = '';
	export { className as class };

	$: inputId = id || (label ? `textarea-${label.toLowerCase().replace(/\s+/g, '-')}` : undefined);
</script>

<div class={cn('flex flex-col gap-1.5 w-full', className)}>
	{#if label || $$slots['label-extra']}
		<div class="flex items-center justify-between">
			{#if label}
				<label for={inputId} class="text-xs font-medium text-(--text-secondary) select-none">
					{label}
				</label>
			{/if}
			<slot name="label-extra" />
		</div>
	{/if}

	<textarea
		id={inputId}
		{placeholder}
		{disabled}
		{readonly}
		{rows}
		bind:value
		class={cn(
			'w-full px-3 py-2 text-xs sm:text-sm bg-(--bg-surface) text-(--text-main) placeholder:text-(--text-muted) border rounded-[4px] outline-none transition-colors duration-150 resize-y',
			error
				? 'border-rose-500 focus:border-rose-500 focus:ring-1 focus:ring-rose-500'
				: 'border-(--border-hairline) focus:border-(--accent-sky) focus:ring-1 focus:ring-(--accent-sky)',
			disabled && 'opacity-50 cursor-not-allowed bg-(--bg-muted)'
		)}
		on:input
		on:change
		on:focus
		on:blur
		on:keydown
		{...$$restProps}
	></textarea>

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
