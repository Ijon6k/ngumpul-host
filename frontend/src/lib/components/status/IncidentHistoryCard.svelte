<script lang="ts">
	import type { IncidentRecord } from '$lib/types/status';
	import { formatDate } from '$lib/utils/format';

	export let incidents: IncidentRecord[] = [];
</script>

<section class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-6 sm:p-7 flex flex-col gap-4 shadow-xs">
	<!-- Header -->
	<div class="flex items-center justify-between pb-2 border-b border-(--border-hairline)">
		<div>
			<h3 class="text-base font-normal text-(--text-main)">
				Incident History
			</h3>
			<p class="text-xs text-(--text-muted) mt-0.5">
				Documented downtime gaps and state transitions
			</p>
		</div>
		<span class="text-xs text-(--text-muted)">
			Filtered to active range
		</span>
	</div>

	<!-- Content -->
	{#if incidents && incidents.length > 0}
		<div class="divide-y divide-(--border-hairline)">
			{#each incidents as incident (incident.id)}
				<div class="py-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs">
					<div class="flex flex-col">
						<div class="flex items-center gap-2">
							<span class="w-2 h-2 rounded-full bg-rose-500 shrink-0"></span>
							<span class="font-medium text-(--text-main)">{incident.cause}</span>
							<span class="text-(--text-muted)">· {incident.duration_formatted} duration</span>
						</div>
						<span class="text-(--text-secondary) mt-1">{incident.details}</span>
					</div>

					<div class="text-right text-(--text-muted) shrink-0 font-mono text-[11px]">
						{formatDate(incident.started_at)}
					</div>
				</div>
			{/each}
		</div>
	{:else}
		<!-- Reassuring Clean Empty State -->
		<div class="py-6 flex flex-col items-center justify-center text-center gap-2">
			<span class="w-2 h-2 rounded-full bg-emerald-500"></span>
			<p class="text-sm font-normal text-(--text-main)">
				No downtime incidents recorded
			</p>
			<p class="text-xs text-(--text-muted) max-w-sm">
				Zero downtime gaps or unplanned outages detected during this monitoring period.
			</p>
		</div>
	{/if}
</section>
