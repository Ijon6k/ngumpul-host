<script lang="ts">
	import { CaretLeft, CaretRight } from 'phosphor-svelte';

	export let currentPage: number = 1;
	export let totalItems: number = 0;
	export let pageSize: number = 12;
	export let onPageChange: (page: number) => void;

	$: totalPages = Math.max(1, Math.ceil(totalItems / pageSize));
	$: startIndex = totalItems === 0 ? 0 : (currentPage - 1) * pageSize + 1;
	$: endIndex = Math.min(totalItems, currentPage * pageSize);

	function goToPage(page: number) {
		if (page >= 1 && page <= totalPages && page !== currentPage) {
			onPageChange(page);
		}
	}

	$: visiblePages = (() => {
		const pages: (number | 'ellipsis')[] = [];
		if (totalPages <= 7) {
			for (let i = 1; i <= totalPages; i++) pages.push(i);
		} else {
			pages.push(1);
			if (currentPage > 3) {
				pages.push('ellipsis');
			}
			const start = Math.max(2, currentPage - 1);
			const end = Math.min(totalPages - 1, currentPage + 1);
			for (let i = start; i <= end; i++) {
				pages.push(i);
			}
			if (currentPage < totalPages - 2) {
				pages.push('ellipsis');
			}
			pages.push(totalPages);
		}
		return pages;
	})();
</script>

{#if totalPages > 1}
	<nav
		aria-label="Pagination Navigation"
		class="flex flex-col sm:flex-row items-center justify-between gap-4 py-4 text-xs select-none border-t border-(--border-hairline)"
	>
		<!-- Counter info -->
		<div class="text-(--text-muted) font-mono">
			Showing <span class="text-(--text-main) font-medium">{startIndex}–{endIndex}</span> of
			<span class="text-(--text-main) font-medium">{totalItems}</span>
		</div>

		<!-- Page navigation buttons -->
		<div class="flex items-center gap-1">
			<button
				type="button"
				disabled={currentPage <= 1}
				class="inline-flex items-center justify-center w-8 h-8 rounded-[4px] border border-(--border-hairline) bg-(--bg-surface) text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-hover) disabled:opacity-40 disabled:pointer-events-none transition-colors"
				aria-label="Previous page"
				on:click={() => goToPage(currentPage - 1)}
			>
				<CaretLeft size={14} weight="bold" />
			</button>

			{#each visiblePages as p, i}
				{#if p === 'ellipsis'}
					<span class="w-8 h-8 flex items-center justify-center text-(--text-muted) font-mono">…</span>
				{:else}
					<button
						type="button"
						class="w-8 h-8 rounded-[4px] border text-xs font-mono transition-colors {p === currentPage
							? 'bg-(--accent-sky) text-white border-(--accent-sky) font-semibold shadow-xs'
							: 'bg-(--bg-surface) text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-hover) border-(--border-hairline)'}"
						aria-current={p === currentPage ? 'page' : undefined}
						on:click={() => goToPage(p)}
					>
						{p}
					</button>
				{/if}
			{/each}

			<button
				type="button"
				disabled={currentPage >= totalPages}
				class="inline-flex items-center justify-center w-8 h-8 rounded-[4px] border border-(--border-hairline) bg-(--bg-surface) text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-hover) disabled:opacity-40 disabled:pointer-events-none transition-colors"
				aria-label="Next page"
				on:click={() => goToPage(currentPage + 1)}
			>
				<CaretRight size={14} weight="bold" />
			</button>
		</div>
	</nav>
{/if}
