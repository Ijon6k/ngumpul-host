<script lang="ts">
	/**
	 * Destination URL for link breadcrumbs.
	 */
	export let href: string | undefined = undefined;

	/**
	 * Whether this breadcrumb item represents the current active page.
	 */
	export let current: boolean = false;

	/**
	 * Optional Phosphor icon component.
	 */
	export let icon: any = undefined;

	/**
	 * Whether to render a separator slash before this item.
	 * Set to false for the first item in custom slotted breadcrumbs.
	 */
	export let separator: boolean = true;

	/**
	 * Optional text label if not using slot content.
	 */
	export let label: string = '';

	let className: string = '';
	export { className as class };
</script>

{#if separator}
	<span class="text-(--text-muted) opacity-40 select-none" aria-hidden="true">/</span>
{/if}

{#if href && !current}
	<a
		{href}
		class="hover:text-(--text-main) transition-colors flex items-center gap-1.5 {className}"
	>
		{#if icon}
			<svelte:component this={icon} size={14} weight="bold" class="shrink-0" />
		{/if}
		<slot>{label}</slot>
	</a>
{:else}
	<span
		class="text-(--text-main) font-medium flex items-center gap-1.5 {className}"
		aria-current={current ? 'page' : undefined}
	>
		{#if icon}
			<svelte:component this={icon} size={14} weight="bold" class="shrink-0" />
		{/if}
		<slot>{label}</slot>
	</span>
{/if}
