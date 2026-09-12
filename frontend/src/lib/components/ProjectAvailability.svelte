<script lang="ts">
	import { Pulse, Clock } from 'phosphor-svelte';
	import StatusDot from './StatusDot.svelte';

	export let availability: any = null;
	export let status: string = 'ONLINE';

	function formatUptime(val?: number | null): string {
		if (val == null) return '—';
		return `${val.toFixed(1)}%`;
	}

	function formatTime(dateStr?: string | null): string {
		if (!dateStr) return '—';
		try {
			const d = new Date(dateStr);
			return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
		} catch {
			return '—';
		}
	}
</script>

<div class="p-5 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-4">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-2">
			<Pulse size={16} class="text-(--text-muted)" />
			<h3 class="text-xs font-semibold text-(--text-main)">Availability & Health</h3>
		</div>
		<div class="flex items-center gap-1.5 text-xs">
			<StatusDot {status} />
			<span class="font-medium text-(--text-main)">{status === 'ONLINE' ? 'Operational' : 'Unavailable'}</span>
		</div>
	</div>

	<div class="grid grid-cols-2 gap-3 pt-1 border-t border-(--border-hairline)">
		<div class="flex flex-col gap-0.5">
			<span class="text-[11px] text-(--text-muted)">30-day Uptime</span>
			<span class="font-mono text-sm font-semibold text-(--text-main)">
				{formatUptime(availability?.uptime_percent)}
			</span>
		</div>

		<div class="flex flex-col gap-0.5">
			<span class="text-[11px] text-(--text-muted)">Latest Latency</span>
			<span class="font-mono text-sm font-semibold text-(--text-main)">
				{availability?.latest_response_time_ms ? `${availability.latest_response_time_ms} ms` : '—'}
			</span>
		</div>
	</div>

	{#if availability?.last_checked_at}
		<div class="flex items-center gap-1.5 text-[11px] text-(--text-muted) pt-1 border-t border-(--border-hairline)">
			<Clock size={12} />
			<span>Checked at {formatTime(availability.last_checked_at)} (periodic 5m probe)</span>
		</div>
	{/if}
</div>
