<script lang="ts">
	export let value: string | number = '';
	export let type: string = 'text';
	export let placeholder: string = '';
	export let label: string = '';
	export let helperText: string = '';
	export let error: string = '';
	export let disabled: boolean = false;
	export let readonly: boolean = false;
	export let id: string = '';

	let className: string = '';
	export { className as class };

	$: inputId = id || (label ? `input-${label.toLowerCase().replace(/\s+/g, '-')}` : undefined);
</script>

<div class="flex flex-col gap-1.5 w-full {className}">
	{#if label}
		<label for={inputId} class="text-xs font-medium text-(--text-secondary) select-none">
			{label}
		</label>
	{/if}

	<div class="relative flex items-center w-full">
		<input
			{id}
			{type}
			{placeholder}
			{disabled}
			{readonly}
			bind:value
			class="w-full h-9 px-3 text-xs sm:text-sm bg-(--bg-surface) text-(--text-main) placeholder:text-(--text-muted) border rounded-[4px] outline-none transition-colors duration-150 {error
				? 'border-rose-500 focus:border-rose-500 focus:ring-1 focus:ring-rose-500'
				: 'border-(--border-hairline) focus:border-(--accent-sky) focus:ring-1 focus:ring-(--accent-sky)'} {disabled ? 'opacity-50 cursor-not-allowed bg-(--bg-muted)' : ''}"
			on:input
			on:change
			on:focus
			on:blur
			on:keydown
			{...$$restProps}
		/>
	</div>

	{#if error}
		<p class="text-[11px] text-rose-600 dark:text-rose-400 font-normal">
			{error}
		</p>
	{:else if helperText}
		<p class="text-[11px] text-(--text-muted) font-normal">
			{helperText}
		</p>
	{/if}
</div>
