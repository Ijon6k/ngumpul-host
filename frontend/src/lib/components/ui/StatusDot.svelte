<script lang="ts">
	interface Props {
		status?: string;
		showLabel?: boolean;
		variant?: 'inline' | 'dot' | 'badge';
	}

	let {
		status = 'ONLINE',
		showLabel = true,
		variant = 'inline'
	}: Props = $props();

	let normalized = $derived((status || 'ONLINE').toUpperCase());

	let label = $derived.by(() => {
		switch (normalized) {
			case 'ONLINE':
			case 'REACHABLE':
			case 'OPERATIONAL':
				return 'Online';
			case 'UNREACHABLE':
			case 'OFFLINE':
			case 'DEGRADED':
				return 'Unreachable';
			case 'SUSPENDED':
				return 'Suspended';
			case 'SETUP':
				return 'Setting up';
			case 'PENDING':
				return 'Pending';
			case 'CHECKING':
			case 'UNKNOWN':
				return 'Checking';
			case 'ARCHIVED':
				return 'Archived';
			default:
				return normalized;
		}
	});

	let colorClasses = $derived.by(() => {
		switch (normalized) {
			case 'ONLINE':
			case 'REACHABLE':
			case 'OPERATIONAL':
				return {
					dot: 'bg-emerald-500',
					badge: 'bg-(--bg-surface)/85 dark:bg-neutral-900/85 text-emerald-700 dark:text-emerald-300 border-emerald-500/30'
				};
			case 'UNREACHABLE':
			case 'OFFLINE':
			case 'DEGRADED':
				return {
					dot: 'bg-rose-500',
					badge: 'bg-(--bg-surface)/85 dark:bg-neutral-900/85 text-rose-700 dark:text-rose-300 border-rose-500/30'
				};
			case 'SUSPENDED':
				return {
					dot: 'bg-amber-500',
					badge: 'bg-(--bg-surface)/85 dark:bg-neutral-900/85 text-amber-700 dark:text-amber-300 border-amber-500/30'
				};
			case 'SETUP':
				return {
					dot: 'bg-(--accent-sky)',
					badge: 'bg-(--bg-surface)/85 dark:bg-neutral-900/85 text-sky-700 dark:text-sky-300 border-(--accent-sky)/30'
				};
			case 'CHECKING':
			case 'PENDING':
			case 'UNKNOWN':
				return {
					dot: 'bg-amber-400',
					badge: 'bg-(--bg-surface)/85 dark:bg-neutral-900/85 text-amber-600 dark:text-amber-300 border-amber-400/30'
				};
			case 'ARCHIVED':
			default:
				return {
					dot: 'bg-neutral-400 dark:bg-neutral-500',
					badge: 'bg-(--bg-surface)/85 dark:bg-neutral-900/85 text-(--text-muted) border-(--border-hairline)'
				};
		}
	});

	let activeVariant = $derived(!showLabel ? 'dot' : variant);
</script>

{#if activeVariant === 'dot'}
	<span class="inline-block w-1.5 h-1.5 rounded-full {colorClasses.dot} shrink-0" title={label}></span>
{:else if activeVariant === 'badge'}
	<span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border backdrop-blur-xs shadow-2xs {colorClasses.badge}" title={label}>
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
