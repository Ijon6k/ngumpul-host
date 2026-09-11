<script lang="ts">
	import type { ServiceHealth } from '$lib/types/status';
	import ServiceStatusRow from './ServiceStatusRow.svelte';
	import { formatTime } from '$lib/utils/format';

	export let services: ServiceHealth[] = [];
	export let checkedAt: string = '';
</script>

<section class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl overflow-hidden shadow-xs">
	<!-- Card Header -->
	<div class="px-6 py-4 border-b border-(--border-hairline) flex items-center justify-between text-xs text-(--text-muted)">
		<div>
			<span class="font-medium text-(--text-main)">Service Health</span>
			<span class="hidden sm:inline text-(--text-muted) ml-2">· Live response time & gateway probes</span>
		</div>
		<span>Target &lt;250ms</span>
	</div>

	<!-- Services List -->
	<div class="divide-y divide-(--border-hairline)">
		{#each services as service (service.id)}
			<ServiceStatusRow
				name={service.name}
				role={service.id === 'portal' ? 'Frontend UI & Dashboard Client' : service.id === 'api' ? 'Core HTTP REST API' : 'Community Hosted Projects'}
				status={service.status}
				latencyMs={service.response_time_ms}
				lastCheckedAt={service.last_checked_at}
			/>
		{/each}
	</div>

	<!-- Card Footer -->
	<div class="px-6 py-3 bg-(--bg-muted)/30 border-t border-(--border-hairline) flex items-center justify-between text-xs text-(--text-muted)">
		<span>Probe frequency: continuous</span>
		<span>Last verified: {formatTime(checkedAt)}</span>
	</div>
</section>
