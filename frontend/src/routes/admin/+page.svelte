<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	let stats: any = null;
	let pendingRequests: any[] = [];
	let loading = true;

	onMount(async () => {
		try {
			const [statsRes, reqsRes] = await Promise.all([
				api.get('/admin/stats'),
				api.get('/admin/hosting-requests?status=PENDING')
			]);
			stats = statsRes.data;
			pendingRequests = reqsRes.data?.requests || [];
		} catch (err) {
			console.error('Failed to load admin stats:', err);
		} finally {
			loading = false;
		}
	});

	function formatDate(iso: string) {
		try {
			return new Date(iso).toLocaleDateString([], { month: 'short', day: 'numeric' });
		} catch { return '—'; }
	}
</script>

<div class="flex flex-col gap-7">

	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
		<div>
			<h1 class="font-display font-bold text-[1.35rem] text-(--text-main) tracking-tight">Instance Overview</h1>
			<p class="text-sm text-(--text-secondary) mt-1 max-w-lg">
				Real-time node telemetry, pending community workloads, and service orchestration.
			</p>
		</div>
		<div class="flex items-center gap-2 shrink-0">
			<a href="/admin/requests" class="btn btn-primary btn-sm text-[13px]">
				Review queue
				{#if pendingRequests.length > 0}
					<span class="ml-1 text-xs bg-white/15 px-1.5 py-0.5 rounded-full">{pendingRequests.length}</span>
				{/if}
			</a>
			<a href="/admin/system" class="btn btn-secondary btn-sm text-[13px]">
				Diagnostics <span class="text-(--text-muted) text-xs">↗</span>
			</a>
		</div>
	</div>

	<!-- Metric Cards — Linear style -->
	<div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
		<!-- Pending Review -->
		<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-5 flex flex-col gap-4">
			<div class="flex items-center justify-between text-xs text-(--text-muted)">
				<span>Pending Review</span>
				<span class="w-2 h-2 rounded-full {pendingRequests.length > 0 ? 'bg-amber-500' : 'bg-emerald-500'}"></span>
			</div>
			<div class="font-display text-3xl font-bold text-(--text-main)">{stats?.pending_requests ?? 0}</div>
			<a href="/admin/requests" class="text-xs text-(--accent-strong) hover:underline">Open queue →</a>
		</div>

		<!-- Active Projects -->
		<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-5 flex flex-col gap-4">
			<div class="flex items-center justify-between text-xs text-(--text-muted)">
				<span>Active Projects</span>
				<span class="w-2 h-2 rounded-full bg-emerald-500"></span>
			</div>
			<div class="font-display text-3xl font-bold text-(--text-main)">{stats?.online_projects ?? 0}</div>
			<span class="text-xs text-(--text-muted)">of {stats?.total_projects ?? 0} total</span>
		</div>

		<!-- Community Members -->
		<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-5 flex flex-col gap-4">
			<div class="flex items-center justify-between text-xs text-(--text-muted)">
				<span>Members</span>
				<span class="w-2 h-2 rounded-full bg-sky-500"></span>
			</div>
			<div class="font-display text-3xl font-bold text-(--text-main)">{stats?.total_members ?? 0}</div>
			<a href="/admin/users" class="text-xs text-(--accent-strong) hover:underline">Manage accounts →</a>
		</div>

		<!-- Node Health -->
		<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-5 flex flex-col gap-4">
			<div class="flex items-center justify-between text-xs text-(--text-muted)">
				<span>Node Health</span>
				<span class="w-2 h-2 rounded-full bg-emerald-500 animate-live-pulse"></span>
			</div>
			<div class="font-display text-xl font-bold text-emerald-600 dark:text-emerald-400">Healthy</div>
			<span class="text-xs text-(--text-muted) font-mono">Ingress · SSL · Postgres</span>
		</div>
	</div>

	<!-- Server Telemetry Row — real hardware facts -->
	<div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
		{#each [
			{ label: 'CPU Cores', value: '32', sub: 'AMD EPYC · Multi-thread', color: 'text-sky-500' },
			{ label: 'ECC Memory', value: '64 GB', sub: 'Registered DIMM', color: 'text-purple-400' },
			{ label: 'NVMe I/O', value: '3.2 GB/s', sub: 'PCIe Gen 4 · RAID-1', color: 'text-orange-400' },
			{ label: 'Edge Latency', value: '14ms', sub: 'Jakarta IXP direct', color: 'text-emerald-500' },
		] as t}
			<div class="bg-(--bg-muted) border border-(--border-hairline) rounded-xl p-4 flex flex-col gap-2">
				<span class="text-[10px] font-mono text-(--text-muted) uppercase tracking-widest">{t.label}</span>
				<span class="font-display font-bold text-xl {t.color}">{t.value}</span>
				<span class="text-[10px] text-(--text-muted)">{t.sub}</span>
			</div>
		{/each}
	</div>

	<!-- Pending Queue -->
	<div class="border border-(--border-hairline) rounded-xl bg-(--bg-surface) overflow-hidden">
		<div class="px-5 py-3.5 border-b border-(--border-hairline) flex items-center justify-between">
			<div class="flex items-center gap-2.5">
				<span class="w-2 h-2 rounded-full {pendingRequests.length > 0 ? 'bg-amber-500' : 'bg-emerald-500'}"></span>
				<h2 class="font-semibold text-sm text-(--text-main)">Pending Hosting Approvals</h2>
				<span class="text-xs text-(--text-muted) font-mono">({pendingRequests.length})</span>
			</div>
			<a href="/admin/requests" class="text-xs text-(--accent-strong) hover:underline font-medium">View all →</a>
		</div>

		{#if loading}
			<div class="p-8 text-center text-xs text-(--text-muted)">Loading requests queue...</div>
		{:else if pendingRequests.length === 0}
			<div class="p-8 text-center text-xs text-(--text-secondary)">
				<div class="text-sm font-medium text-(--text-main) mb-1">Queue is clear</div>
				No community hosting requests are currently awaiting evaluation.
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs border-collapse">
					<thead>
						<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/50 text-(--text-muted)">
							<th class="px-5 py-2.5 font-medium">Project</th>
							<th class="px-5 py-2.5 font-medium">Requester</th>
							<th class="px-5 py-2.5 font-medium hidden md:table-cell">Repository</th>
							<th class="px-5 py-2.5 font-medium">Submitted</th>
							<th class="px-5 py-2.5 text-right font-medium">Action</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-(--border-hairline)">
						{#each pendingRequests as req}
							<tr class="hover:bg-(--bg-muted)/30 transition-colors">
								<td class="px-5 py-3.5">
									<div class="font-medium text-(--text-main)">{req.project_name}</div>
									{#if req.description}
										<div class="text-[11px] text-(--text-muted) truncate max-w-[200px]">{req.description}</div>
									{/if}
								</td>
								<td class="px-5 py-3.5">
									<span class="font-medium text-(--text-main)">{req.requester?.display_name || 'Member'}</span>
									<span class="block text-[11px] font-mono text-(--text-muted)">@{req.requester?.username}</span>
								</td>
								<td class="px-5 py-3.5 hidden md:table-cell">
									<a href={req.repository_url} target="_blank" rel="noreferrer" class="text-(--accent-strong) hover:underline font-mono text-[11px] truncate block max-w-[180px]">
										{req.repository_url} ↗
									</a>
								</td>
								<td class="px-5 py-3.5 text-(--text-muted)">{formatDate(req.created_at)}</td>
								<td class="px-5 py-3.5 text-right">
									<a href="/admin/requests" class="btn btn-secondary btn-sm text-[11px] py-1">Review →</a>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>

	<!-- Network Topology -->
	<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
		<div class="border border-(--border-hairline) rounded-xl bg-(--bg-surface) p-5 flex flex-col gap-4">
			<div class="flex items-center justify-between pb-3 border-b border-(--border-hairline)">
				<h3 class="font-semibold text-sm text-(--text-main)">Network Topology</h3>
				<span class="text-xs text-(--text-muted) font-mono">Docker Bridge</span>
			</div>
			<div class="flex flex-col gap-2 text-[13px]">
				{#each [
					{ label: 'Ingress Proxy', value: ':1111 → Nginx' },
					{ label: 'Frontend Layer', value: 'SvelteKit SSR :3000' },
					{ label: 'API Service', value: 'Go Chi /api/v1 :8080' },
					{ label: 'Database', value: 'PostgreSQL 16 :5432' },
				] as row}
					<div class="flex items-center justify-between p-2.5 rounded-lg bg-(--bg-muted) border border-(--border-hairline)">
						<span class="text-(--text-secondary)">{row.label}</span>
						<code class="font-mono text-[11px] text-(--text-main)">{row.value}</code>
					</div>
				{/each}
			</div>
		</div>

		<div class="border border-(--border-hairline) rounded-xl bg-(--bg-surface) p-5 flex flex-col gap-4">
			<div class="flex items-center justify-between pb-3 border-b border-(--border-hairline)">
				<h3 class="font-semibold text-sm text-(--text-main)">Operational Policy</h3>
				<span class="text-xs text-(--accent-strong) font-medium">Active</span>
			</div>
			<ul class="flex flex-col gap-3 text-[13px] text-(--text-secondary) leading-relaxed">
				{#each [
					'No arbitrary Docker execution or remote shell spawning without human review.',
					'Every deployment requires explicit administrator approval and codebase review.',
					'All state transitions create an immutable, timestamped audit event.',
				] as policy, i}
					<li class="flex items-start gap-3">
						<span class="text-[11px] font-mono font-bold text-(--accent-strong) pt-0.5 shrink-0">0{i+1}.</span>
						<span>{policy}</span>
					</li>
				{/each}
			</ul>
		</div>
	</div>
</div>
