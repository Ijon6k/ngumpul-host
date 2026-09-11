<script lang="ts">
	export let title: string = '';
	export let description: string = '';
	export let badge: string = '';
	export let padding: boolean = true;
	let className: string = '';
	export { className as class };
</script>

<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl overflow-hidden shadow-2xs flex flex-col {className}">
	{#if title || $$slots.actions}
		<div class="px-6 py-4 border-b border-(--border-hairline) flex items-center justify-between flex-wrap gap-3">
			<div>
				<div class="flex items-center gap-2">
					<h2 class="font-semibold text-base text-(--text-main)">{title}</h2>
					{#if badge}
						<span class="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-(--cf-pastel-bg) text-(--cf-pastel-text)">
							{badge}
						</span>
					{/if}
				</div>
				{#if description}
					<p class="text-xs text-(--text-secondary) mt-0.5">{description}</p>
				{/if}
			</div>

			{#if $$slots.actions}
				<div class="flex items-center gap-2">
					<slot name="actions" />
				</div>
			{/if}
		</div>
	{/if}

	<div class={padding ? 'p-6' : ''}>
		<slot />
	</div>

	{#if $$slots.footer}
		<div class="px-6 py-3.5 border-t border-(--border-hairline) bg-(--bg-muted)/30 text-xs">
			<slot name="footer" />
		</div>
	{/if}
</div>
