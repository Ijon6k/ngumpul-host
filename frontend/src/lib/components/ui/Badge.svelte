<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		variant?: 'success' | 'warning' | 'danger' | 'sky' | 'neutral';
		size?: 'sm' | 'md';
		dot?: boolean;
		class?: string;
		children?: Snippet;
		[key: string]: any;
	}

	let {
		variant = 'neutral',
		size = 'sm',
		dot = false,
		class: className = '',
		children,
		...restProps
	}: Props = $props();

	const variantStyles = {
		success: 'bg-(--bg-muted) dark:bg-neutral-900 text-emerald-600 dark:text-emerald-400',
		warning: 'bg-(--bg-muted) dark:bg-neutral-900 text-amber-600 dark:text-amber-400',
		danger: 'bg-(--bg-muted) dark:bg-neutral-900 text-rose-600 dark:text-rose-400',
		sky: 'bg-(--bg-muted) dark:bg-neutral-900 text-sky-600 dark:text-sky-400',
		neutral: 'bg-(--bg-muted) dark:bg-neutral-900 text-(--text-secondary)'
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

	let baseClass = $derived(
		`inline-flex items-center font-medium rounded-full ${variantStyles[variant]} ${sizeStyles[size]} ${className}`
	);
</script>

<span class={baseClass} {...restProps}>
	{#if dot}
		<span class="w-1.5 h-1.5 rounded-full {dotStyles[variant]} shrink-0"></span>
	{/if}
	{@render children?.()}
</span>
