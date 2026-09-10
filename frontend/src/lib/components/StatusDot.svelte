<script lang="ts">
	export let status: string = 'ONLINE';
	export let showLabel: boolean = true;

	$: normalized = (status || 'ONLINE').toUpperCase();

	$: label = (() => {
		switch (normalized) {
			case 'ONLINE':
			case 'OPERATIONAL':
				return 'Online';
			case 'SETUP':
				return 'Setup';
			case 'PENDING':
				return 'Pending';
			case 'OFFLINE':
			case 'DEGRADED':
				return 'Offline';
			case 'ARCHIVED':
				return 'Archived';
			default:
				return normalized;
		}
	})();

	$: dotColorClass = (() => {
		switch (normalized) {
			case 'ONLINE':
			case 'OPERATIONAL':
				return 'bg-(--color-success)';
			case 'SETUP':
			case 'PENDING':
				return 'bg-(--color-warning)';
			case 'OFFLINE':
			case 'DEGRADED':
				return 'bg-(--color-danger)';
			default:
				return 'bg-(--text-muted)';
		}
	})();
</script>

<span class="inline-flex items-center gap-1.5 text-xs text-(--text-secondary)" title={label}>
	<span class="w-1.5 h-1.5 rounded-full shrink-0 {dotColorClass}"></span>
	{#if showLabel}
		<span class="font-medium">{label}</span>
	{/if}
</span>
