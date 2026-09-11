<script lang="ts">
	export let title: string = '';
	export let badge: string = '';
	export let badgeClass: string = 'text-(--text-muted)';
	export let padding: string = 'p-6 sm:p-7';
	let className: string = '';
	export { className as class };
</script>

<div
	class="bg-(--bg-surface) border border-(--border-hairline) rounded-2xl flex flex-col justify-between shadow-xs transition-all duration-200 {padding} {className}"
>
	{#if title || badge || $$slots.header}
		<div class="flex items-center justify-between pb-3.5 border-b border-(--border-hairline) gap-3">
			{#if $$slots.header}
				<slot name="header" />
			{:else}
				{#if title}
					<h3 class="font-semibold text-sm sm:text-base text-(--text-main) tracking-tight">{title}</h3>
				{/if}
				{#if badge}
					<span class="text-xs font-mono {badgeClass}">{badge}</span>
				{/if}
			{/if}
		</div>
	{/if}

	<div class="grow flex flex-col justify-between {title || badge || $$slots.header ? 'pt-4' : ''}">
		<slot />
	</div>

	{#if $$slots.footer}
		<div class="pt-3.5 mt-auto border-t border-(--border-hairline) text-xs">
			<slot name="footer" />
		</div>
	{/if}
</div>
