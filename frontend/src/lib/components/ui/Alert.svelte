<script lang="ts">
	import { WarningCircle, CheckCircle, Info, ShieldWarning } from 'phosphor-svelte';
	import { cn } from '$lib/utils';

	export let variant: 'danger' | 'warning' | 'success' | 'info' = 'danger';
	export let title: string = '';
	export let hideIcon: boolean = false;

	let className: string = '';
	export { className as class };

	const variantStyles = {
		danger: 'bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30',
		warning: 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/30',
		success: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30',
		info: 'bg-(--accent-sky)/10 text-(--accent-sky) border border-(--accent-sky)/30'
	};
</script>

<div
	role="alert"
	class={cn(
		'p-3.5 rounded-md text-xs sm:text-sm leading-relaxed flex items-start gap-2.5 transition-colors',
		variantStyles[variant],
		className
	)}
>
	{#if !hideIcon}
		<div class="shrink-0 mt-0.5">
			{#if variant === 'danger'}
				<ShieldWarning size={16} weight="bold" />
			{:else if variant === 'warning'}
				<WarningCircle size={16} weight="bold" />
			{:else if variant === 'success'}
				<CheckCircle size={16} weight="bold" />
			{:else}
				<Info size={16} weight="bold" />
			{/if}
		</div>
	{/if}

	<div class="flex-1 min-w-0">
		{#if title}
			<p class="font-medium mb-0.5">{title}</p>
		{/if}
		<slot />
	</div>
</div>
