<script lang="ts">
	import StatusIndicator from './StatusIndicator.svelte';

	export let name: string = '';
	export let role: string = '';
	export let status: string = 'OPERATIONAL';
	export let latencyMs: number = 0;
	export let lastCheckedAt: string = '';

	function formatTime(iso: string): string {
		try {
			const d = new Date(iso);
			return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
		} catch (_) {
			return '';
		}
	}
</script>

<div class="px-5 sm:px-6 py-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:bg-(--bg-muted)/40 transition-colors">
	<div class="flex flex-col">
		<span class="text-sm font-medium text-(--text-main)">
			{name}
		</span>
		{#if role}
			<span class="text-xs text-(--text-muted) mt-0.5">
				{role}
			</span>
		{/if}
	</div>

	<div class="flex items-center gap-5 sm:gap-6 self-start sm:self-auto">
		{#if latencyMs > 0}
			<div class="flex items-center gap-1 text-xs text-(--text-secondary)">
				<span class="font-mono text-[11px] text-(--text-muted)">{latencyMs}ms</span>
				<span class="text-[10px] text-(--text-muted)/70">latency</span>
			</div>
		{/if}

		<div class="shrink-0" title={lastCheckedAt ? `Last verified: ${formatTime(lastCheckedAt)}` : ''}>
			<StatusIndicator {status} size="md" />
		</div>
	</div>
</div>
