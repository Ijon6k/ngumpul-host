<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	let logs: any[] = [];
	let loading = true;
	let searchQuery = '';
	let actionFilter = 'ALL';

	async function loadLogs() {
		loading = true;
		try {
			const res = await api.get('/admin/audit');
			logs = res.data?.audit_logs || [];
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
			<p class="text-xs text-(--text-secondary) mt-1">
				Cryptographic record of administrative operations, project state mutations, approvals, and security events.
			</p>
		</div>

		<div class="flex items-center gap-2">
			<button type="button" class="btn btn-secondary btn-sm font-mono text-xs" on:click={loadLogs}>
				Refresh
			</button>
		</div>
	</div>

	<!-- Filter & Search Controls -->
	<div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
		<div class="flex items-center gap-2">
			<label for="audit-action-filter" class="text-xs font-mono text-(--text-muted)">Action:</label>
			<select
				id="audit-action-filter"
				class="px-2.5 py-1 text-xs font-mono bg-(--bg-surface) border border-(--border-hairline) rounded text-(--text-main) outline-none focus:border-(--accent-sky)"
				bind:value={actionFilter}
			>
				<option value="ALL">ALL ACTIONS</option>
				{#each actionsList as act}
					<option value={act}>{act}</option>
				{/each}
			</select>
		</div>

		<div class="w-full sm:w-64">
			<input
				type="text"
				placeholder="Filter audit ledger..."
				bind:value={searchQuery}
				class="w-full px-3 py-1.5 bg-(--bg-surface) border border-(--border-hairline) rounded text-xs text-(--text-main) placeholder:text-(--text-muted) outline-none focus:border-(--accent-sky) font-sans"
			/>
		</div>
	</div>

	<!-- High-Density Audit Ledger Table -->
	<div class="border border-(--border-hairline) rounded-md bg-(--bg-surface) overflow-hidden">
		{#if loading}
			<div class="p-10 text-center text-xs text-(--text-muted)">Reading immutable ledger...</div>
		{:else if filteredLogs.length === 0}
			<div class="p-10 text-center text-xs text-(--text-secondary)">
				No audit entries match current filter.
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs border-collapse font-sans">
					<thead>
						<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-muted) text-xs">
							<th class="px-5 py-2.5 font-medium">Timestamp</th>
							<th class="px-5 py-2.5 font-medium">Actor</th>
							<th class="px-5 py-2.5 font-medium">Action</th>
							<th class="px-5 py-2.5 font-medium">Target</th>
							<th class="px-5 py-2.5 font-medium">Metadata Context</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-(--border-hairline)">
						{#each filteredLogs as log (log.id)}
							<tr class="hover:bg-(--bg-muted)/25 transition-colors">
								<td class="px-5 py-3.5 align-top font-mono text-[11px] text-(--text-muted) whitespace-nowrap">
									{formatDate(log.created_at)}
								</td>
								<td class="px-5 py-3.5 align-top">
									<strong class="text-(--text-main) font-medium">{log.actor_name || 'Administrator'}</strong>
									{#if log.actor_username}
										<span class="block text-[11px] font-mono text-(--text-muted)">@{log.actor_username}</span>
									{/if}
								</td>
								<td class="px-5 py-3.5 align-top">
									<span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-(--bg-muted) text-(--text-main) border border-(--border-hairline)">
										{log.action}
									</span>
								</td>
								<td class="px-5 py-3.5 align-top font-mono text-[11px]">
									<span class="text-(--text-main)">{log.target_type}</span>
									{#if log.target_id}
										<span class="text-(--text-muted) block text-[10px]">#{log.target_id.slice(0, 8)}</span>
									{/if}
								</td>
								<td class="px-5 py-3.5 align-top">
									{#if log.metadata && Object.keys(log.metadata).length > 0}
										<pre class="text-[10px] font-mono text-(--text-secondary) bg-(--bg-muted)/60 p-2 rounded max-w-sm overflow-x-auto border border-(--border-subtle)">{JSON.stringify(log.metadata, null, 2)}</pre>
									{:else}
										<span class="text-[11px] font-mono text-(--text-muted)">—</span>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
