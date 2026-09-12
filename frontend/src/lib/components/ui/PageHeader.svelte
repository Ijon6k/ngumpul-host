<script lang="ts">
	export let title: string;
	export let description: string = '';

	let className: string = '';
	export { className as class };
</script>

<div class="flex flex-col md:flex-row md:items-end justify-between gap-4 pb-2 border-b border-(--border-hairline) {className}">
	<div class="max-w-2xl flex flex-col gap-1.5">
		{#if $$slots.breadcrumb}
			<div class="mb-1 text-sm text-(--text-muted)">
				<slot name="breadcrumb" />
			</div>
		{/if}

		<h1 class="font-sans font-normal text-2xl sm:text-3xl lg:text-4xl text-(--text-main) tracking-[-0.03em] leading-tight">
			{title}
		</h1>

		{#if description}
			<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed font-normal">
				{description}
			</p>
		{/if}
	</div>

	{#if $$slots.actions}
		<div class="flex items-center gap-2.5 shrink-0 self-start md:self-end">
			<slot name="actions" />
		</div>
	{/if}
</div>
