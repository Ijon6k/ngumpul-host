<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes, HTMLAnchorAttributes } from 'svelte/elements';

	export let variant: 'primary' | 'secondary' | 'ghost' | 'danger' = 'secondary';
	export let size: 'sm' | 'md' | 'lg' = 'sm';
	export let href: string | undefined = undefined;
	export let disabled: boolean = false;
	export let loading: boolean = false;
	export let type: 'button' | 'submit' | 'reset' = 'button';

	let className: string = '';
	export { className as class };

	// Variant styling strictly respecting design.md tokens and radius-sm (2px - 4px)
	const variantStyles = {
		primary: 'bg-(--accent-sky) text-white hover:bg-(--accent-strong) active:scale-[0.99] border-transparent shadow-xs font-medium',
		secondary: 'bg-(--bg-surface) text-(--text-main) border-(--border-hairline) hover:bg-(--bg-hover) hover:border-(--text-muted)/40 active:scale-[0.99] shadow-xs font-normal',
		ghost: 'bg-transparent text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-hover) border-transparent font-normal',
		danger: 'bg-rose-600/10 text-rose-700 dark:text-rose-400 border border-rose-300 dark:border-rose-900/50 hover:bg-rose-600/20 active:scale-[0.99] font-medium'
	};

	const sizeStyles = {
		sm: 'h-8 px-3 text-xs gap-1.5',
		md: 'h-9 px-4 text-xs sm:text-sm gap-2',
		lg: 'h-10 px-5 text-sm gap-2.5'
	};

	$: baseClass = `inline-flex items-center justify-center rounded-[4px] transition-all duration-150 select-none cursor-pointer border ${variantStyles[variant]} ${sizeStyles[size]} ${disabled || loading ? 'opacity-50 pointer-events-none' : ''} ${className}`;
</script>

{#if href && !disabled}
	<a {href} class={baseClass} {...$$restProps}>
		{#if loading}
			<span class="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin"></span>
		{/if}
		<slot />
	</a>
{:else}
	<button {type} {disabled} class={baseClass} on:click {...$$restProps}>
		{#if loading}
			<span class="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin"></span>
		{/if}
		<slot />
	</button>
{/if}
