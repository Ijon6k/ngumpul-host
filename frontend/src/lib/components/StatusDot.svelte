<script lang="ts">
	export let status: string = 'ONLINE';
	export let showLabel: boolean = true;
	export let variant: 'inline' | 'dot' | 'badge' = 'inline';

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

	$: colorClasses = (() => {
		switch (normalized) {
			case 'ONLINE':
			case 'OPERATIONAL':
				return {
					dot: 'bg-emerald-500',
					badge: 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border-emerald-500/20'
				};
			case 'SETUP':
			case 'PENDING':
				return {
					dot: 'bg-amber-500',
					badge: 'bg-amber-500/10 text-amber-700 dark:text-amber-400 border-amber-500/20'
				};
			case 'OFFLINE':
			case 'DEGRADED':
				return {
					dot: 'bg-rose-500',
					badge: 'bg-rose-500/10 text-rose-700 dark:text-rose-400 border-rose-500/20'
				};
			case 'ARCHIVED':
			default:
				return {
					dot: 'bg-neutral-400 dark:bg-neutral-500',
					badge: 'bg-(--bg-muted) text-(--text-muted) border-(--border-hairline)'
				};
		}
	})();

	$: activeVariant = !showLabel ? 'dot' : variant;
</script>

{#if activeVariant === 'dot'}
	<span class="inline-block w-1.5 h-1.5 rounded-full {colorClasses.dot} shrink-0" title={label}></span>
{:else if activeVariant === 'badge'}
	<span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[11px] font-medium border {colorClasses.badge}" title={label}>
		<span class="w-1.5 h-1.5 rounded-full {colorClasses.dot} shrink-0"></span>
		<span>{label}</span>
	</span>
{:else}
	<!-- Default 'inline' calm indicator: small dot + text, no boxed pill -->
	<span class="inline-flex items-center gap-1.5 text-xs text-(--text-secondary) font-normal" title={label}>
		<span class="w-1.5 h-1.5 rounded-full {colorClasses.dot} shrink-0"></span>
		<span class="text-(--text-main)">{label}</span>
	</span>
{/if}
