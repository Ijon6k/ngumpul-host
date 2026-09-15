<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let tabs: Array<{
		id: string;
		label: string;
		icon?: any;
		badge?: number;
	}> = [];
	export let active: string = tabs[0]?.id || '';
	export let size: 'sm' | 'default' = 'default';

	const dispatch = createEventDispatcher<{
		change: string;
	}>();

	function selectTab(id: string) {
		active = id;
		dispatch('change', id);
	}
</script>

<div
	role="tablist"
	class="flex items-center gap-1 sm:gap-2 border-b border-(--border-hairline) overflow-x-auto no-scrollbar"
>
	{#each tabs as tab}
		{@const isActive = active === tab.id}
		<button
			type="button"
			role="tab"
			aria-selected={isActive}
			class="inline-flex items-center gap-2 border-b-2 transition-colors cursor-pointer whitespace-nowrap {size === 'sm' ? 'px-3 py-2 text-xs' : 'px-4 py-2.5 text-xs sm:text-sm'} {isActive ? 'border-(--accent-sky) text-(--text-main) font-semibold' : 'border-transparent text-(--text-secondary) hover:text-(--text-main) font-normal'}"
			on:click={() => selectTab(tab.id)}
		>
			{#if tab.icon}
				<svelte:component this={tab.icon} size={size === 'sm' ? 14 : 16} weight={isActive ? 'bold' : 'regular'} />
			{/if}
			<span>{tab.label}</span>
			{#if tab.badge !== undefined && tab.badge > 0}
				<span class="px-2 py-0.5 rounded-full text-xs font-mono font-medium bg-(--bg-muted) dark:bg-neutral-900 text-(--accent-sky)">
					{tab.badge}
				</span>
			{/if}
		</button>
	{/each}
</div>
