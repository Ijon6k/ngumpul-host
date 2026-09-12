<script lang="ts">
	import { ArrowUpRight } from 'phosphor-svelte';

	export let title: string = '';
	export let timestamp: string = '';
	export let description: string = '';
	export let dotColor: string = '';
	export let status: string = '';
	export let density: 'comfortable' | 'compact' = 'comfortable';
	export let href: string = '';
	export let linkText: string = '';

	$: computedDotColor = (() => {
		if (dotColor) return dotColor;
		if (status) {
			const s = status.toUpperCase();
			if (s === 'ONLINE' || s === 'OPERATIONAL') return 'bg-emerald-500';
			if (s === 'OFFLINE' || s === 'DEGRADED') return 'bg-rose-500';
			if (s === 'SETUP' || s === 'PENDING') return 'bg-amber-500';
		}
		return 'bg-(--accent-sky)';
	})();
</script>

<div class="relative pl-6 {density === 'compact' ? 'pb-4' : 'pb-7'} last:pb-0 group">
	<!-- Subtle Connecting Line between items -->
	<div
		class="absolute left-[3.5px] top-2.5 bottom-0 w-px bg-(--border-hairline) group-last:hidden"
		aria-hidden="true"
	></div>

	<!-- Dot indicator -->
	<div
		class="absolute left-0 top-1.5 w-2 h-2 rounded-full {computedDotColor} ring-4 ring-(--bg-canvas)"
		aria-hidden="true"
	></div>

	<!-- Item Content -->
	<div class="flex flex-col gap-1">
		<div class="flex items-baseline justify-between gap-3 flex-wrap">
			<slot name="title">
				{#if title}
					<span class="font-medium text-xs sm:text-sm text-(--text-main) leading-snug">
						{title}
					</span>
				{/if}
			</slot>

			<slot name="timestamp">
				{#if timestamp}
					<span class="text-xs text-(--text-muted) font-mono shrink-0">
						{timestamp}
					</span>
				{/if}
			</slot>
		</div>

		<slot>
			{#if description}
				<p class="text-xs text-(--text-secondary) leading-relaxed">
					{description}
				</p>
			{/if}
		</slot>

		{#if href}
			<div class="pt-0.5">
				<a
					{href}
					class="text-xs font-medium text-(--accent-sky) hover:underline inline-flex items-center gap-1"
				>
					<span>{linkText || href}</span>
					<ArrowUpRight size={11} weight="bold" />
				</a>
			</div>
		{/if}
	</div>
</div>
