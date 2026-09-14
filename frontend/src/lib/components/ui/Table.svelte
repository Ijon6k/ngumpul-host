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
	import Skeleton from './Skeleton.svelte';

	export let columns: TableColumn[] = [];
	export let items: T[] | undefined = undefined;
	export let loading: boolean = false;
	export let loadingMessage: string = 'Loading records...';
	export let emptyMessage: string = 'No records found.';
	export let keyField: string = 'id';
	export let alignTop: boolean = false;
	export let skeletonRows: number = 5;
	let className: string = '';
	export { className as class };
</script>

<div class="w-full overflow-hidden rounded-md border border-(--border-hairline) bg-(--bg-surface) transition-colors {className}">
	{#if loading}
		<div class="aria-busy:overflow-hidden" role="status" aria-live="polite">
			<div class="overflow-x-auto">
				<table class="w-full text-left text-sm border-collapse font-sans">
					<thead>
						<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-secondary) text-sm font-medium">
							{#if columns && columns.length > 0}
								{#each columns as col}
									<th
										style={col.width ? `width: ${col.width}` : undefined}
										class="px-5 py-3 font-semibold {col.align === 'right' ? 'text-right' : col.align === 'center' ? 'text-center' : 'text-left'}"
									>
										{col.label}
									</th>
								{/each}
							{:else}
								{@const cols = skeletonRows > 0 ? 4 : 1}
								{#each Array(cols) as _}
									<th class="px-5 py-3"></th>
								{/each}
							{/if}
						</tr>
					</thead>
					<tbody class="divide-y divide-(--border-hairline)">
						{#each Array(Math.max(skeletonRows, 1)) as _}
							<tr>
								{#if columns && columns.length > 0}
									{#each columns as col}
										<TableCell align={col.align || 'left'} {alignTop} class={col.cellClass || ''}>
											<Skeleton variant="line" class={col.align === 'right' ? 'ml-auto' : ''} />
										</TableCell>
									{/each}
								{:else}
									{#each Array(4) as _}
										<TableCell>
											<Skeleton variant="line" />
										</TableCell>
									{/each}
								{/if}
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			<span class="sr-only">{loadingMessage}</span>
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
