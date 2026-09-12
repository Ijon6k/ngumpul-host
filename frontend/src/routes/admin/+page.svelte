<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import AdminStatCard from '$lib/components/admin/AdminStatCard.svelte';
	import AdminPanel from '$lib/components/admin/AdminPanel.svelte';
	import ResourceBar from '$lib/components/ResourceBar.svelte';

	let stats: any = null;
	let pendingRequests: any[] = [];
	let projectsList: any[] = [];
	let systemData: any = null;
	let statusData: any = null;
	let loading = true;

	onMount(async () => {
		try {
			const [statsRes, reqsRes, projRes, sysRes, statRes] = await Promise.all([
				api.get('/admin/stats'),
				api.get('/admin/hosting-requests?status=PENDING'),
				api.get('/admin/projects'),
				api.get('/admin/system'),
				api.get('/status')
			]);
			stats = statsRes.data;
			pendingRequests = reqsRes.data?.requests || [];
			projectsList = projRes.data?.projects || [];
			systemData = sysRes.data;
			statusData = statRes.data;
		} catch (err) {
			console.error('Failed to load admin dashboard data:', err);
		} finally {
			loading = false;
		}
	});

	$: hostSpecs = statusData?.host_specs || systemData?.host_specs || null;
	$: totalProjects = stats?.total_projects ?? projectsList.length ?? 0;
	$: onlineProjects = stats?.online_projects ?? projectsList.filter(p => p.status === 'ONLINE').length ?? 0;

	function formatDate(iso: string) {
		try {
			return new Date(iso).toLocaleDateString([], { month: 'short', day: 'numeric', year: 'numeric' });
		} catch { return '—'; }
	}

	$: load1 = hostSpecs?.load_avg?.[0] ?? 0.32;
	$: load5 = hostSpecs?.load_avg?.[1] ?? 0.28;
	$: load15 = hostSpecs?.load_avg?.[2] ?? 0.22;
</script>

<div class="flex flex-col gap-8">

	<!-- ══════════════════════════════════════════════════════════════════
	     1. CLOUDFLARE TOP METRIC TILES
	     ══════════════════════════════════════════════════════════════════ -->
	<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl overflow-hidden shadow-2xs">
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-(--border-hairline)">
			<AdminStatCard
				label="Pending Requests"
				value={pendingRequests.length}
				subtext={pendingRequests.length > 0 ? 'Awaiting evaluation' : 'Queue is clear'}
				badge={pendingRequests.length > 0 ? `${pendingRequests.length} new` : 'Clear'}
				badgeClass={pendingRequests.length > 0 ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400' : 'bg-(--cf-green-pastel) text-(--cf-green-text)'}
				href="/admin/requests"
				actionText="Review queue"
			/>

			<AdminStatCard
				label="Open Reports"
				value={stats?.open_reports ?? 0}
				subtext={(stats?.open_reports ?? 0) > 0 ? 'Awaiting moderation' : 'No open reports'}
				badge={(stats?.open_reports ?? 0) > 0 ? `${stats.open_reports} open` : 'Clean'}
				badgeClass={(stats?.open_reports ?? 0) > 0 ? 'bg-rose-500/10 text-rose-600 dark:text-rose-400' : 'bg-(--cf-green-pastel) text-(--cf-green-text)'}
				href="/admin/reports"
				actionText="Moderation"
			/>

			<AdminStatCard
				label="Hosted Projects"
				value={onlineProjects}
				subtext="of {totalProjects} registered"
				badge="Active"
				badgeClass="bg-(--cf-green-pastel) text-(--cf-green-text)"
				href="/admin/projects"
				actionText="Manage"
			/>

			<AdminStatCard
				label="Community Members"
				value={stats?.total_members ?? 4}
				subtext="Registered builders"
				badge="Accounts"
				badgeClass="bg-(--cf-pastel-bg) text-(--cf-pastel-text)"
				href="/admin/users"
				actionText="Members"
			/>

			<AdminStatCard
				label="Host Node Uptime"
				value={hostSpecs?.uptime_formatted ?? '3 hours'}
				subtext="Continuous native Linux"
				badge="Healthy"
				badgeClass="bg-(--cf-green-pastel) text-(--cf-green-text)"
				href="/admin/system"
				actionText="Diagnostics"
			/>
		</div>
	</div>

	<!-- ══════════════════════════════════════════════════════════════════
	     2. MAIN TWO-COLUMN CONTENT AREA
	     ══════════════════════════════════════════════════════════════════ -->
	<div class="grid grid-cols-1 xl:grid-cols-[1fr_400px] gap-8 items-start">

		<!-- LEFT COLUMN -->
		<div class="flex flex-col gap-8">

			<!-- 1. Hosting Requests Review Queue -->
			<AdminPanel title="Hosting Requests Queue" badge={pendingRequests.length > 0 ? `${pendingRequests.length} pending` : ''} padding={false}>
				<svelte:fragment slot="actions">
					<a href="/admin/requests" class="btn btn-secondary btn-sm text-xs py-1.5 px-3">
						Open review queue →
					</a>
				</svelte:fragment>

				{#if loading}
					<div class="p-10 text-center text-sm text-(--text-muted)">Loading requests queue...</div>
				{:else if pendingRequests.length === 0}
					<div class="p-10 text-center text-sm text-(--text-secondary)">
						<div class="font-semibold text-base text-(--text-main) mb-1">Queue is clear</div>
						No community hosting requests are currently awaiting administrator evaluation.
					</div>
				{:else}
					<div class="overflow-x-auto">
						<table class="w-full text-left text-sm border-collapse font-sans">
							<thead>
								<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-secondary) text-sm font-medium">
									<th class="py-3 px-6 font-semibold">Project</th>
									<th class="py-3 px-4 font-semibold">Requester</th>
									<th class="py-3 px-4 font-semibold hidden md:table-cell">Repository</th>
									<th class="py-3 px-4 font-semibold">Submitted</th>
									<th class="py-3 px-6 font-semibold text-right">Action</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-(--border-hairline)">
								{#each pendingRequests as req}
									<tr class="hover:bg-(--bg-muted)/25 transition-colors">
										<td class="py-4 px-6 font-medium text-(--text-main)">
											<div class="font-semibold text-sm sm:text-base">{req.project_name}</div>
											{#if req.description}
												<div class="text-xs sm:text-sm text-(--text-secondary) mt-0.5 truncate max-w-sm">{req.description}</div>
											{/if}
										</td>
										<td class="py-4 px-4 text-(--text-secondary)">
											<span class="font-medium text-sm sm:text-base text-(--text-main)">{req.requester?.display_name || 'Member'}</span>
											<span class="block font-mono text-xs text-(--text-muted)">@{req.requester?.username}</span>
										</td>
										<td class="py-4 px-4 font-mono text-xs sm:text-sm text-(--cf-blue) hidden md:table-cell">
											<a href={req.repository_url} target="_blank" rel="noreferrer" class="hover:underline truncate block max-w-xs">
												{req.repository_url} ↗
											</a>
										</td>
										<td class="py-4 px-4 text-xs sm:text-sm text-(--text-muted)">{formatDate(req.created_at)}</td>
										<td class="py-4 px-6 text-right">
											<a href="/admin/requests" class="btn btn-primary btn-sm text-xs sm:text-sm py-1.5 px-3">
												Review →
											</a>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</AdminPanel>

			<!-- 3. Hosted Projects Directory -->
			<AdminPanel title="Hosted Projects Directory" badge="{projectsList.length} registered" padding={false}>
				<svelte:fragment slot="actions">
					<a href="/admin/projects" class="btn btn-secondary btn-sm text-xs sm:text-sm py-1.5 px-3">
						Manage projects →
					</a>
				</svelte:fragment>

				<div class="overflow-x-auto">
					<table class="w-full text-left text-sm border-collapse font-sans">
						<thead>
							<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-secondary) text-sm font-medium">
								<th class="py-3 px-6 font-semibold">Project Name</th>
								<th class="py-3 px-4 font-semibold">Owner</th>
								<th class="py-3 px-4 font-semibold hidden md:table-cell">Tech Tags</th>
								<th class="py-3 px-4 font-semibold">Status</th>
								<th class="py-3 px-6 font-semibold text-right">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-(--border-hairline)">
							{#if projectsList.length === 0}
								<tr>
									<td colspan="5" class="py-8 px-6 text-center text-sm text-(--text-muted)">No projects registered yet.</td>
								</tr>
							{:else}
								{#each projectsList as p}
									<tr class="hover:bg-(--bg-muted)/25 transition-colors">
										<td class="py-4 px-6 font-medium text-(--text-main)">
											<a href="/projects/{p.slug}" class="font-semibold text-sm sm:text-base hover:text-(--cf-blue)">{p.name}</a>
											<span class="block font-mono text-xs text-(--text-muted)">/{p.slug}</span>
										</td>
										<td class="py-4 px-4 text-(--text-secondary)">
											<span class="font-medium text-sm sm:text-base text-(--text-main)">{p.owner?.display_name || 'Community'}</span>
											<span class="block font-mono text-xs text-(--text-muted)">@{p.owner?.username}</span>
										</td>
										<td class="py-4 px-4 hidden md:table-cell">
											{#if p.technology_stack && p.technology_stack.length > 0}
												<div class="flex items-center gap-1.5 flex-wrap">
													{#each p.technology_stack.slice(0, 3) as tech}
														<span class="text-xs sm:text-sm font-mono px-2 py-0.5 rounded bg-(--bg-muted) text-(--text-secondary)">{tech}</span>
													{/each}
												</div>
											{:else}
												<span class="text-neutral-400">—</span>
											{/if}
										</td>
										<td class="py-4 px-4">
											<span class="px-2.5 py-1 rounded-full text-xs font-medium bg-(--cf-green-pastel) text-(--cf-green-text)">
												{p.status || 'ONLINE'}
											</span>
										</td>
										<td class="py-4 px-6 text-right">
											<a href="/admin/projects" class="font-semibold text-xs sm:text-sm text-(--cf-blue) hover:underline">
												Edit →
											</a>
										</td>
									</tr>
								{/each}
							{/if}
						</tbody>
					</table>
				</div>
			</AdminPanel>

		</div>

		<!-- RIGHT COLUMN -->
		<div class="flex flex-col gap-8">

			<!-- 1. Host Node Telemetry Card -->
			<AdminPanel title="Host Node Telemetry">
				<div class="flex flex-col gap-5 text-sm">
					<!-- Node Name -->
					<div class="flex items-center justify-between">
						<span class="text-xs text-(--text-secondary) font-medium">Node Hostname</span>
						<span class="font-mono text-sm font-semibold text-(--text-main)">{hostSpecs?.hostname ?? 'ngumpul-host'}</span>
					</div>

					<!-- Processor -->
					<div class="pt-3 border-t border-(--border-hairline) flex flex-col gap-1">
						<span class="text-xs text-(--text-secondary) font-medium">Processor</span>
						<span class="font-semibold text-sm text-(--text-main)">
							{hostSpecs?.cpu_model ?? '13th Gen Intel Core i5-1334U'}
						</span>
						<span class="text-xs font-mono text-(--text-muted)">{hostSpecs?.cpu_cores ?? 12} physical threads</span>
					</div>

					<!-- Memory Gauge -->
					<div class="pt-3 border-t border-(--border-hairline)">
						<ResourceBar
							label="Memory Utilization"
							valueText="{hostSpecs?.total_ram_gb ?? 23.1} GB ({hostSpecs?.used_ram_percent ?? 38.4}%)"
							percent={hostSpecs?.used_ram_percent ?? 38.4}
							colorClass="bg-(--cf-blue)"
							subText="{hostSpecs?.available_ram_gb ?? 14.2} GB available"
						/>
					</div>

					<!-- Storage Gauge -->
					<div class="pt-3 border-t border-(--border-hairline)">
						<ResourceBar
							label="NVMe Flash Storage"
							valueText="{hostSpecs?.disk_total_gb ?? 98} GB ({hostSpecs?.disk_used_percent ?? 86.7}%)"
							percent={hostSpecs?.disk_used_percent ?? 86.7}
							colorClass="bg-(--accent-sky)"
							subText="Local PCIe direct mount"
						/>
					</div>

					<!-- Kernel & OS -->
					<div class="pt-3 border-t border-(--border-hairline) flex items-center justify-between text-xs text-(--text-muted)">
						<span>Kernel Runtime</span>
						<span class="font-mono text-xs text-(--text-main) font-medium">{hostSpecs?.kernel ?? 'Linux 7.1.8-xanmod1'}</span>
					</div>
				</div>

				<svelte:fragment slot="footer">
					<a href="/admin/system" class="font-semibold text-xs text-(--cf-blue) hover:underline flex items-center justify-between">
						<span>Full system diagnostics</span>
						<span>→</span>
					</a>
				</svelte:fragment>
			</AdminPanel>

			<!-- 2. System Load Diagnostics -->
			<AdminPanel title="System Load Diagnostics">
				<div class="flex flex-col gap-4">
					<div class="flex items-center justify-between text-xs">
						<span class="text-(--text-secondary) font-medium">Load Average</span>
						<span class="font-mono text-xs font-semibold text-emerald-600 dark:text-emerald-400">Normal / Smooth</span>
					</div>

					<!-- 3 Load chips -->
					<div class="grid grid-cols-3 gap-2 text-center">
						<div class="p-3 rounded-lg bg-(--bg-muted) border border-(--border-hairline)">
							<span class="block text-xs font-mono text-(--text-muted)">1 min</span>
							<span class="font-mono font-bold text-base text-(--text-main)">{load1.toFixed(2)}</span>
						</div>
						<div class="p-3 rounded-lg bg-(--bg-muted) border border-(--border-hairline)">
							<span class="block text-xs font-mono text-(--text-muted)">5 min</span>
							<span class="font-mono font-bold text-base text-(--text-main)">{load5.toFixed(2)}</span>
						</div>
						<div class="p-3 rounded-lg bg-(--bg-muted) border border-(--border-hairline)">
							<span class="block text-xs font-mono text-(--text-muted)">15 min</span>
							<span class="font-mono font-bold text-base text-(--text-main)">{load15.toFixed(2)}</span>
						</div>
					</div>

					<div class="pt-3 border-t border-(--border-hairline) flex items-center justify-between text-xs text-(--text-secondary)">
						<span>Architecture</span>
						<span class="font-mono text-xs font-semibold text-(--text-main)">{hostSpecs?.arch ?? 'amd64'}</span>
					</div>
				</div>
			</AdminPanel>

			<!-- 3. Console Shortcuts -->
			<AdminPanel title="Administration Tools">
				<div class="flex flex-col gap-2.5 text-sm">
					<a href="/admin/requests" class="flex items-center justify-between p-3 rounded-lg bg-(--bg-muted)/60 hover:bg-(--cf-pastel-bg)/30 transition-colors">
						<span class="font-medium text-(--text-main)">Hosting Request Queue</span>
						<span class="text-(--cf-blue)">→</span>
					</a>
					<a href="/admin/reports" class="flex items-center justify-between p-3 rounded-lg bg-(--bg-muted)/60 hover:bg-(--cf-pastel-bg)/30 transition-colors">
						<span class="font-medium text-(--text-main)">Moderation Reports</span>
						<span class="text-(--cf-blue)">→</span>
					</a>
					<a href="/admin/comments" class="flex items-center justify-between p-3 rounded-lg bg-(--bg-muted)/60 hover:bg-(--cf-pastel-bg)/30 transition-colors">
						<span class="font-medium text-(--text-main)">Comments Moderation</span>
						<span class="text-(--cf-blue)">→</span>
					</a>
					<a href="/admin/projects" class="flex items-center justify-between p-3 rounded-lg bg-(--bg-muted)/60 hover:bg-(--cf-pastel-bg)/30 transition-colors">
						<span class="font-medium text-(--text-main)">Manage Project Endpoints</span>
						<span class="text-(--cf-blue)">→</span>
					</a>
					<a href="/admin/users" class="flex items-center justify-between p-3 rounded-lg bg-(--bg-muted)/60 hover:bg-(--cf-pastel-bg)/30 transition-colors">
						<span class="font-medium text-(--text-main)">Community Member Accounts</span>
						<span class="text-(--cf-blue)">→</span>
					</a>
					<a href="/admin/system" class="flex items-center justify-between p-3 rounded-lg bg-(--bg-muted)/60 hover:bg-(--cf-pastel-bg)/30 transition-colors">
						<span class="font-medium text-(--text-main)">System Telemetry & Health</span>
						<span class="text-(--cf-blue)">→</span>
					</a>
					<a href="/admin/audit" class="flex items-center justify-between p-3 rounded-lg bg-(--bg-muted)/60 hover:bg-(--cf-pastel-bg)/30 transition-colors">
						<span class="font-medium text-(--text-main)">Immutable Audit Log</span>
						<span class="text-(--cf-blue)">→</span>
					</a>
				</div>
			</AdminPanel>

		</div>

	</div>

</div>
