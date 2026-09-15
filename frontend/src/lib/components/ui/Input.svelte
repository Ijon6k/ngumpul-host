<script lang="ts">
	import { cn } from '$lib/utils';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { Snippet } from 'svelte';

	import { Eye, EyeSlash } from 'phosphor-svelte';

	interface Props extends HTMLInputAttributes {
		value?: string | number;
		type?: string;
		placeholder?: string;
		label?: string;
		helperText?: string;
		error?: string;
		disabled?: boolean;
		readonly?: boolean;
		required?: boolean;
		autocomplete?: HTMLInputAttributes['autocomplete'];
		id?: string;
		inputClass?: string;
		labelClass?: string;
		class?: string;
		labelExtra?: Snippet;
		[key: string]: any;
	}

	let {
		value = $bindable(''),
		type = 'text',
		placeholder = '',
		label = '',
		helperText = '',
		error = '',
		disabled = false,
		readonly = false,
		required = false,
		autocomplete = undefined,
		id = '',
		inputClass = '',
		labelClass = '',
		class: className = '',
		labelExtra,
		...restProps
	}: Props = $props();

	let showPassword = $state(false);
	let effectiveType = $derived(type === 'password' ? (showPassword ? 'text' : 'password') : type);
	let inputId = $derived(id || (label ? `input-${label.toLowerCase().replace(/\s+/g, '-')}` : undefined));
</script>

<div class={cn('flex flex-col gap-1.5 w-full', className)}>
	{#if label || labelExtra}
		<div class="flex items-center justify-between">
			{#if label}
				<label for={inputId} class={cn('text-xs sm:text-sm font-medium text-(--text-main) select-none', labelClass)}>
					{label}
				</label>
			{/if}
			{@render labelExtra?.()}
		</div>
	{/if}

	<div class="relative flex items-center w-full">
		<input
			id={inputId}
			type={effectiveType}
			{placeholder}
			{disabled}
			{readonly}
			{required}
			{autocomplete}
			bind:value
			class={cn(
				'w-full h-9 px-3 text-xs sm:text-sm bg-(--bg-surface) text-(--text-main) placeholder:text-(--text-muted) border rounded-[4px] outline-none transition-colors duration-150',
				type === 'password' && 'pr-9',
				error
					? 'border-rose-500 focus:border-rose-500 focus:ring-1 focus:ring-rose-500'
					: 'border-(--border-hairline) focus:border-(--accent-sky) focus:ring-1 focus:ring-(--accent-sky)',
				disabled && 'opacity-50 cursor-not-allowed bg-(--bg-muted)',
				inputClass
			)}
			{...restProps}
		/>
		{#if type === 'password'}
			<button
				type="button"
				tabindex="-1"
				onclick={() => (showPassword = !showPassword)}
				class="absolute right-2 text-(--text-muted) hover:text-(--text-main) p-1 rounded-sm cursor-pointer transition-colors"
				aria-label={showPassword ? 'Hide password' : 'Show password'}
			>
				{#if showPassword}
					<EyeSlash size={16} />
				{:else}
					<Eye size={16} />
				{/if}
			</button>
		{/if}
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
