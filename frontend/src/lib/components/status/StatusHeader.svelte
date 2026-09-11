<script lang="ts">
	export let selectedRange: '1d' | '7d' | '30d' = '30d';
	export let onRangeChange: (range: '1d' | '7d' | '30d') => void = () => {};
</script>

<header class="flex flex-col sm:flex-row sm:items-end justify-between gap-5 pb-2 border-b border-(--border-hairline)">
	<div class="max-w-xl">
		<h1 class="font-sans font-normal text-3xl sm:text-4xl text-(--text-main) tracking-tight">
			System Status
		</h1>
		<p class="text-sm sm:text-base text-(--text-secondary) font-normal leading-relaxed mt-1.5">
			Real-time availability probes and response latency metrics for core infrastructure components.
		</p>
	</div>

	<!-- Filter Range Controls -->
	<div class="flex items-center gap-1 p-1 bg-(--bg-muted) rounded-lg border border-(--border-hairline) self-start sm:self-auto shrink-0">
		<button
			type="button"
			class="px-3 py-1.5 text-xs font-normal rounded-md transition-all cursor-pointer {selectedRange === '1d' ? 'bg-(--bg-surface) text-(--text-main) shadow-xs font-medium' : 'text-(--text-muted) hover:text-(--text-main)'}"
			on:click={() => onRangeChange('1d')}
		>
			24 Hours
		</button>
		<button
			type="button"
			class="px-3 py-1.5 text-xs font-normal rounded-md transition-all cursor-pointer {selectedRange === '7d' ? 'bg-(--bg-surface) text-(--text-main) shadow-xs font-medium' : 'text-(--text-muted) hover:text-(--text-main)'}"
			on:click={() => onRangeChange('7d')}
		>
			7 Days
		</button>
		<button
			type="button"
			class="px-3 py-1.5 text-xs font-normal rounded-md transition-all cursor-pointer {selectedRange === '30d' ? 'bg-(--bg-surface) text-(--text-main) shadow-xs font-medium' : 'text-(--text-muted) hover:text-(--text-main)'}"
			on:click={() => onRangeChange('30d')}
		>
			30 Days
		</button>
	</div>
</header>
