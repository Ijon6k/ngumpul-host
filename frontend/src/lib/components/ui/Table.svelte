<script lang="ts" context="module">
	export interface TableColumn {
		key: string;
		label: string;
		align?: 'left' | 'center' | 'right';
		width?: string;
		class?: string;
		headerClass?: string;
		cellClass?: string;
	}
</script>

<script lang="ts" generics="T = any">
	import TableRow from './TableRow.svelte';
	import TableCell from './TableCell.svelte';

	export let columns: TableColumn[] = [];
	export let items: T[] | undefined = undefined;
	export let loading: boolean = false;
	export let loadingMessage: string = 'Loading records...';
	export let emptyMessage: string = 'No records found.';
	export let keyField: string = 'id';
	export let alignTop: boolean = false;
	let className: string = '';
	export { className as class };
</script>

<div class="w-full overflow-hidden rounded-md border border-(--border-hairline) bg-(--bg-surface) transition-colors {className}">
	{#if loading}
		<div class="p-10 text-center text-sm text-(--text-muted) flex flex-col items-center justify-center gap-2.5">
			<div class="w-5 h-5 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>{loadingMessage}</span>
		</div>
	{:else if items !== undefined && items.length === 0}
		<div class="p-10 text-center text-sm text-(--text-secondary)">
			<slot name="empty">
				{emptyMessage}
			</slot>
		</div>
	{:else}
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm border-collapse font-sans">
				{#if columns && columns.length > 0}
					<thead>
						<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-secondary) text-sm font-medium">
							{#each columns as col}
								<th
									style={col.width ? `width: ${col.width}` : undefined}
									class="px-5 py-3 font-semibold {col.align === 'right' ? 'text-right' : col.align === 'center' ? 'text-center' : 'text-left'} {col.headerClass || col.class || ''}"
								>
									{col.label}
								</th>
							{/each}
						</tr>
					</thead>
				{:else if $$slots.header}
					<thead>
						<slot name="header" />
					</thead>
				{/if}

				<tbody class="divide-y divide-(--border-hairline)">
					{#if items !== undefined && items.length > 0}
						{#each items as item, index ((item as any)?.[keyField] ?? index)}
							<slot {item} {index}>
								<TableRow>
									{#each columns as col}
										<TableCell align={col.align || 'left'} {alignTop} class={col.cellClass || ''}>
											{(item as any)?.[col.key] ?? '—'}
										</TableCell>
									{/each}
								</TableRow>
							</slot>
						{/each}
					{:else}
						<slot name="content" />
					{/if}
				</tbody>
			</table>
		</div>
	{/if}
</div>
