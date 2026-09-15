<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		variant?: 'surface' | 'muted' | 'interactive';
		padding?: 'none' | 'sm' | 'md' | 'lg';
		href?: string;
		class?: string;
		children?: Snippet;
		[key: string]: any;
	}

	let {
		variant = 'surface',
		padding = 'md',
		href = undefined,
		class: className = '',
		children,
		...restProps
	}: Props = $props();

	const variantStyles = {
		surface: 'bg-(--bg-surface) border-(--border-hairline) shadow-xs',
		muted: 'bg-(--bg-muted) border-(--border-hairline)',
		interactive: 'bg-(--bg-surface) border-(--border-hairline) hover:border-(--text-muted)/40 transition-colors shadow-xs cursor-pointer'
	};

	const paddingStyles = {
		none: 'p-0',
		sm: 'p-4',
		md: 'p-5 sm:p-6',
		lg: 'p-6 sm:p-8'
	};

	let baseClass = $derived(
		`rounded-[8px] border transition-colors duration-200 ${variantStyles[variant]} ${paddingStyles[padding]} ${className}`
	);
</script>

{#if href}
	<a {href} class="block {baseClass}" {...restProps}>
		{@render children?.()}
	</a>
{:else}
	<div class={baseClass} {...restProps}>
		{@render children?.()}
	</div>
{/if}
