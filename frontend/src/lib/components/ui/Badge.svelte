<script lang="ts">
	export let variant: 'success' | 'warning' | 'danger' | 'sky' | 'neutral' = 'neutral';
	export let size: 'sm' | 'md' = 'sm';
	export let dot: boolean = false;

	let className: string = '';
	export { className as class };

	const variantStyles = {
		success: 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border-emerald-500/20',
		warning: 'bg-amber-500/10 text-amber-700 dark:text-amber-400 border-amber-500/20',
		danger: 'bg-rose-500/10 text-rose-700 dark:text-rose-400 border-rose-500/20',
		sky: 'bg-sky-500/10 text-sky-700 dark:text-sky-400 border-sky-500/20',
		neutral: 'bg-(--bg-muted) text-(--text-secondary) border-(--border-hairline)'
	};

	const dotStyles = {
		success: 'bg-emerald-500',
		warning: 'bg-amber-500',
		danger: 'bg-rose-500',
		sky: 'bg-sky-500',
		neutral: 'bg-(--text-muted)'
	};

	const sizeStyles = {
		sm: 'px-2 py-0.5 text-xs gap-1.5',
		md: 'px-2.5 py-1 text-xs gap-1.5'
	};

	$: baseClass = `inline-flex items-center font-normal rounded-full border ${variantStyles[variant]} ${sizeStyles[size]} ${className}`;
</script>

<span class={baseClass} {...$$restProps}>
	{#if dot}
		<span class="w-1.5 h-1.5 rounded-full {dotStyles[variant]} shrink-0"></span>
	{/if}
	<slot />
</span>
