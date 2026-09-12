<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi } from '$lib/api';
	import type { AuditEntry } from '$lib/api/admin/audit';
	import { Table, TableRow, TableCell, type TableColumn } from '$lib/components/ui';

	let logs: AuditEntry[] = [];
	let loading = true;
	let searchQuery = '';
	let actionFilter = 'ALL';

	const columns: TableColumn[] = [
		{ key: 'timestamp', label: 'Timestamp' },
		{ key: 'actor', label: 'Actor' },
		{ key: 'action', label: 'Action' },
		{ key: 'target', label: 'Target' },
		{ key: 'metadata', label: 'Metadata Context' }
	];

	async function loadLogs() {
		loading = true;
		try {
			const res = await adminApi.audit.listAuditLogs();
			logs = res.audit_logs || [];
		} catch (err) {
			console.error('Failed to load audit logs:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadLogs);

	function formatDate(iso: string) {
		try {
			const d = new Date(iso);
			return `${d.toLocaleDateString([], { month: 'short', day: 'numeric', year: 'numeric' })} ${d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}`;
		} catch (_) {
			return '—';
		}
	}

	$: actionsList = Array.from(new Set(logs.map((l) => l.action).filter(Boolean)));

	$: filteredLogs = logs.filter((log) => {
		const matchesAction = actionFilter === 'ALL' || log.action === actionFilter;
		const q = searchQuery.toLowerCase();
		const matchesSearch =
			!searchQuery ||
			log.action?.toLowerCase().includes(q) ||
			log.actor_name?.toLowerCase().includes(q) ||
			log.actor_username?.toLowerCase().includes(q) ||
			log.target_type?.toLowerCase().includes(q) ||
			JSON.stringify(log.metadata || {}).toLowerCase().includes(q);
		return matchesAction && matchesSearch;
	});
</script>

<div class="flex flex-col gap-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-(--border-hairline)">
		<div>
			<h1 class="font-display font-bold text-2xl text-(--text-main)">Immutable Audit Log</h1>
			<p class="text-sm text-(--text-secondary) mt-1">
				Cryptographic record of administrative operations, project state mutations, approvals, and security events.
			</p>
		</div>

		<div class="flex items-center gap-2">
			<button type="button" class="btn btn-secondary btn-sm font-mono text-xs sm:text-sm" on:click={loadLogs}>
				Refresh
			</button>
		</div>
	</div>

	<!-- Filter & Search Controls -->
	<div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
		<div class="flex items-center gap-2">
			<label for="audit-action-filter" class="text-sm font-medium text-(--text-secondary)">Action:</label>
			<select
				id="audit-action-filter"
				class="px-3 py-1.5 text-sm font-mono bg-(--bg-surface) border border-(--border-hairline) rounded text-(--text-main) outline-none focus:border-(--accent-sky)"
				bind:value={actionFilter}
			>
				<option value="ALL">ALL ACTIONS</option>
				{#each actionsList as act}
					<option value={act}>{act}</option>
				{/each}
			</select>
		</div>

		<div class="w-full sm:w-72">
			<input
				type="text"
				placeholder="Filter audit ledger..."
				bind:value={searchQuery}
				class="w-full px-3.5 py-2 bg-(--bg-surface) border border-(--border-hairline) rounded text-sm text-(--text-main) placeholder:text-(--text-muted) outline-none focus:border-(--accent-sky) font-sans"
			/>
		</div>
	</div>

	<!-- High-Density Audit Ledger Reusable Table -->
	<Table
		{columns}
		items={filteredLogs}
		{loading}
		alignTop
		loadingMessage="Reading immutable ledger..."
		emptyMessage="No audit entries match current filter."
		let:item={log}
	>
		<TableRow>
			<TableCell alignTop mono class="text-xs text-(--text-muted) whitespace-nowrap">
				{formatDate(log.created_at)}
			</TableCell>
			<TableCell alignTop>
				<strong class="text-(--text-main) font-medium text-sm sm:text-base">{log.actor_name || 'Administrator'}</strong>
				{#if log.actor_username}
					<span class="block text-xs font-mono text-(--text-muted)">@{log.actor_username}</span>
				{/if}
			</TableCell>
			<TableCell alignTop>
				<span class="inline-flex items-center px-2.5 py-1 rounded text-xs font-mono font-medium bg-(--bg-muted) text-(--text-main) border border-(--border-hairline)">
					{log.action}
				</span>
			</TableCell>
			<TableCell alignTop mono class="text-xs sm:text-sm">
				<span class="text-(--text-main)">{log.target_type}</span>
				{#if log.target_id}
					<span class="text-(--text-muted) block text-xs">#{log.target_id.slice(0, 8)}</span>
				{/if}
			</TableCell>
			<TableCell alignTop>
				{#if log.metadata && Object.keys(log.metadata).length > 0}
					<pre class="text-xs font-mono text-(--text-secondary) bg-(--bg-muted)/60 p-2.5 rounded max-w-md overflow-x-auto border border-(--border-subtle)">{JSON.stringify(log.metadata, null, 2)}</pre>
				{:else}
					<span class="text-xs font-mono text-(--text-muted)">—</span>
				{/if}
			</TableCell>
		</TableRow>
	</Table>
</div>
