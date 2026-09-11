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

	$: styleClasses = (() => {
		switch (normalized) {
			case 'ONLINE':
			case 'OPERATIONAL':
				return 'bg-(--cf-green-pastel) text-(--cf-green-text)';
			case 'SETUP':
			case 'PENDING':
				return 'bg-amber-500/10 text-amber-700 dark:text-amber-400';
			case 'OFFLINE':
			case 'DEGRADED':
				return 'bg-rose-500/10 text-rose-700 dark:text-rose-400';
			default:
				return 'bg-(--bg-muted) text-(--text-muted)';
		}
	})();
</script>

<span class="inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-medium {styleClasses}" title={label}>
	{#if showLabel}
		<span>{label}</span>
	{/if}
</span>
