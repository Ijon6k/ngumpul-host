<script context="module" lang="ts">
	export interface BreadcrumbItemData {
		label: string;
		href?: string;
		icon?: any;
		current?: boolean;
	}
</script>

<script lang="ts">
	/**
	 * List of breadcrumb items to render declaratively.
	 * If omitted or empty, children can be provided via default slot.
	 */
	export let items: BreadcrumbItemData[] = [];

	/**
	 * Accessible label for the navigation element.
	 */
	export let ariaLabel: string = 'Breadcrumb';

	let className: string = '';
	export { className as class };
</script>

<nav aria-label={ariaLabel} class="flex items-center flex-wrap gap-2 text-sm text-(--text-muted) {className}">
	{#if items.length > 0}
		{#each items as item, i}
			{#if i > 0}
				<span class="text-(--text-muted) opacity-40 select-none" aria-hidden="true">/</span>
			{/if}

			{#if item.href && !item.current && i < items.length - 1}
				<a
					href={item.href}
					class="hover:text-(--text-main) transition-colors flex items-center gap-1.5"
				>
					{#if item.icon}
						<svelte:component this={item.icon} size={14} weight="bold" class="shrink-0" />
					{/if}
					<span>{item.label}</span>
				</a>
			{:else}
				<span
					class="text-(--text-main) font-medium flex items-center gap-1.5"
					aria-current={item.current || i === items.length - 1 ? 'page' : undefined}
				>
					{#if item.icon}
						<svelte:component this={item.icon} size={14} weight="bold" class="shrink-0" />
					{/if}
					<span>{item.label}</span>
				</span>
			{/if}
		{/each}
	{/if}

	<slot />
</nav>
