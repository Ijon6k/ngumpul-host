<script lang="ts">
	import DotMatrix from './DotMatrix.svelte';
	import StatusDot from './StatusDot.svelte';

	export let counts = { projects: 4, members: 12 };
	export let systemStatus = 'OPERATIONAL';
	export let projects: Array<{ name: string; public_url?: string; status: string; slug: string }> = [];
</script>

<div class="w-full bg-(--bg-surface) border border-(--border-hairline) rounded-md p-6 sm:p-7 text-xs text-(--text-secondary) shadow-xs relative overflow-hidden">
	<!-- Top Bar -->
	<div class="flex items-center justify-between pb-4 border-b border-(--border-hairline)">
		<div class="flex items-center gap-2.5">
			<span class="w-2 h-2 rounded-full bg-(--color-success) inline-block"></span>
			<span class="font-semibold text-sm text-(--text-main)">Node 01</span>
			<span class="text-(--text-muted)">·</span>
			<span class="text-xs text-(--text-secondary)">Jakarta Gateway</span>
		</div>
		<div class="flex items-center gap-2">
			<span class="text-[11px] text-(--text-muted)">Cluster</span>
			<div class="text-(--text-main)">
				<DotMatrix text=":1111" dotSize={2} gap={2} />
			</div>
		</div>
	</div>

	<!-- Core Metrics Grid -->
	<div class="grid grid-cols-2 sm:grid-cols-3 gap-4 py-5 border-b border-(--border-hairline)">
		<div>
			<span class="text-xs text-(--text-muted) block mb-1">Active Projects</span>
			<div class="flex items-baseline gap-2">
				<span class="text-2xl font-bold text-(--text-main) font-display">{counts.projects}</span>
				<span class="text-[11px] text-(--color-success) font-medium">online</span>
			</div>
		</div>
		<div>
			<span class="text-xs text-(--text-muted) block mb-1">Community Members</span>
			<div class="flex items-baseline gap-2">
				<span class="text-2xl font-bold text-(--text-main) font-display">{counts.members}</span>
				<span class="text-[11px] text-(--text-muted)">makers</span>
			</div>
		</div>
		<div class="col-span-2 sm:col-span-1">
			<span class="text-xs text-(--text-muted) block mb-1">System Status</span>
			<div class="flex items-center gap-2 mt-1">
				<StatusDot status="ONLINE" showLabel={false} />
				<span class="text-xs font-semibold text-(--text-main)">{systemStatus}</span>
			</div>
		</div>
	</div>

	<!-- Active Services List -->
	<div class="pt-4 flex flex-col gap-2.5">
		<div class="flex items-center justify-between text-xs font-medium text-(--text-muted)">
			<span>Hosted Services</span>
			<a href="/projects" class="text-(--accent-strong) hover:underline text-[11px]">View all →</a>
		</div>

		<div class="flex flex-col divide-y divide-(--border-subtle)">
			{#if projects.length === 0}
				<div class="text-(--text-muted) py-3 text-xs">No active services registered yet.</div>
			{:else}
				{#each projects.slice(0, 4) as proj}
					<div class="flex items-center justify-between text-xs py-2.5 hover:text-(--text-main) transition-colors">
						<div class="flex items-center gap-2 truncate max-w-[65%]">
							<span class="w-1.5 h-1.5 rounded-full bg-(--color-success) shrink-0"></span>
							<a href="/projects/{proj.slug}" class="font-medium hover:text-(--accent-strong) truncate text-(--text-main)">
								{proj.name}
							</a>
						</div>
						<div class="flex items-center gap-3 shrink-0">
							{#if proj.public_url}
								<span class="text-(--text-muted) hidden sm:inline truncate text-[11px] max-w-[140px]">
									{proj.public_url.replace('https://', '')}
								</span>
							{/if}
							<span class="text-[11px] text-(--color-success) font-medium">Ready</span>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
</div>
