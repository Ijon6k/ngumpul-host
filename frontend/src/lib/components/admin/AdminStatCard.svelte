<script lang="ts">
	export let label: string;
	export let value: string | number;
	export let subtext: string = '';
	export let badge: string = '';
	export let badgeClass: string = 'bg-(--cf-pastel-bg) text-(--cf-pastel-text)';
	export let href: string = '';
	export let actionText: string = '';
</script>

<div class="p-6 bg-(--bg-surface) flex flex-col justify-between gap-3">
	<div class="flex items-center justify-between text-xs text-(--text-secondary)">
		<span class="font-medium text-sm text-(--text-secondary)">{label}</span>
		{#if badge}
			<span class="px-2.5 py-0.5 rounded-full text-xs font-medium {badgeClass}">
				{badge}
			</span>
		{/if}
	</div>

	<div class="font-display text-3xl sm:text-4xl font-semibold text-(--text-main) tracking-tight">
		{value}
	</div>

	<div class="flex items-center justify-between text-xs pt-1">
		{#if subtext}
			<span class="text-(--text-muted)">{subtext}</span>
		{:else}
			<span></span>
		{/if}

		{#if href && actionText}
			<a {href} class="font-semibold text-(--cf-blue) hover:underline flex items-center gap-1">
				<span>{actionText}</span>
				<span>→</span>
			</a>
		{/if}
	</div>
</div>
