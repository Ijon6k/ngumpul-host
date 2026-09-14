<script lang="ts">
	import { cn } from '$lib/utils';
	import type { HTMLInputAttributes } from 'svelte/elements';

	export let value: string | number = '';
	export let type: string = 'text';
	export let placeholder: string = '';
	export let label: string = '';
	export let helperText: string = '';
	export let error: string = '';
	export let disabled: boolean = false;
	export let readonly: boolean = false;
	export let required: boolean = false;
	export let autocomplete: HTMLInputAttributes['autocomplete'] = undefined;
	export let id: string = '';
	export let inputClass: string = '';
	export let labelClass: string = '';

	let className: string = '';
	export { className as class };

	$: inputId = id || (label ? `input-${label.toLowerCase().replace(/\s+/g, '-')}` : undefined);
</script>

<div class={cn('flex flex-col gap-1.5 w-full', className)}>
	{#if label || $$slots['label-extra']}
		<div class="flex items-center justify-between">
			{#if label}
				<label for={inputId} class={cn('text-xs sm:text-sm font-medium text-(--text-main) select-none', labelClass)}>
					{label}
				</label>
			{/if}
			<slot name="label-extra" />
		</div>
	{/if}

	<div class="relative flex items-center w-full">
		<input
			id={inputId}
			{type}
			{placeholder}
			{disabled}
			{readonly}
			{required}
			{autocomplete}
			bind:value
			class={cn(
				'w-full h-9 px-3 text-xs sm:text-sm bg-(--bg-surface) text-(--text-main) placeholder:text-(--text-muted) border rounded-[4px] outline-none transition-colors duration-150',
				error
					? 'border-rose-500 focus:border-rose-500 focus:ring-1 focus:ring-rose-500'
					: 'border-(--border-hairline) focus:border-(--accent-sky) focus:ring-1 focus:ring-(--accent-sky)',
				disabled && 'opacity-50 cursor-not-allowed bg-(--bg-muted)',
				inputClass
			)}
			on:input
			on:change
			on:focus
			on:blur
			on:keydown
			{...$$restProps}
		/>
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
