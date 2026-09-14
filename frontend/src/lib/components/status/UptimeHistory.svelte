<script lang="ts">
	import { Clock } from 'phosphor-svelte';

	export let availability: any = null;
	export let blocks: any[] = [];
	export let status: string = 'ONLINE';
	export let variant: 'compact' | 'detailed' = 'compact';
	export let timeUnitLabel: string = '';
	export let range: '1d' | '7d' | '30d' = '30d';
	export let showRangeSelector: boolean = false;
	export let onRangeChange: ((r: '1d' | '7d' | '30d') => void) | null = null;

	interface BarItem {
		date: string;
		label: string;
		status: 'operational' | 'partial' | 'incident' | 'no_data';
		uptimePercent: number | null;
		successfulChecks?: number;
		totalChecks?: number;
		details: string;
	}

	let activeHover: BarItem | null = null;
	let hoverPos = { x: 0, y: 0 };

	// Reactive selected range data
	$: currentRangeData = availability?.ranges?.[range];

	// Reactive blocks resolution: prefer passed blocks, then current range blocks, then availability blocks
	$: activeBlocks = (blocks && blocks.length > 0)
		? blocks
		: (currentRangeData?.blocks && currentRangeData.blocks.length > 0)
		? currentRangeData.blocks
		: (availability?.blocks && availability.blocks.length > 0)
		? availability.blocks
		: [];

	// Map raw blocks into BarItem format
	$: bars = activeBlocks.map((b: any) => ({
		date: b.date,
		label: b.label,
		status: (b.status || (b.uptimePercent != null && b.uptimePercent >= 99 ? 'operational' : b.uptimePercent != null && b.uptimePercent > 0 ? 'partial' : b.uptimePercent === 0 ? 'incident' : 'no_data')) as BarItem['status'],
		uptimePercent: b.uptimePercent,
		successfulChecks: b.successful_checks,
		totalChecks: b.total_checks,
		details: b.details || (b.uptimePercent != null ? `${b.uptimePercent.toFixed(1)}% availability` : 'No telemetry recorded')
	}));

	// Determine if project has any recorded availability telemetry across all ranges
	$: hasAnyTelemetry = (() => {
		if (!availability && (!blocks || blocks.length === 0)) return false;
		if ((availability?.total_checks_30d ?? 0) > 0) return true;
		if ((availability?.total_checks ?? 0) > 0) return true;
		if (availability?.last_checked_at != null) return true;
		if (availability?.ranges) {
			for (const r of Object.values(availability.ranges) as any[]) {
				if ((r?.total_checks ?? 0) > 0) return true;
			}
		}
		return bars.some((b: BarItem) => b.status !== 'no_data' && (b.totalChecks ?? 0) > 0);
	})();

	// Determine if the currently active window has recorded telemetry
	$: hasCurrentRangeData = (() => {
		if (!hasAnyTelemetry) return false;
		if (currentRangeData && (currentRangeData.total_checks ?? 0) > 0) return true;
		return bars.some((b: BarItem) => b.status !== 'no_data' && (b.totalChecks ?? 0) > 0);
	})();

	// Overall uptime percentage for active window
	$: currentUptime = currentRangeData?.uptime_percent ?? availability?.uptime_percent ?? (status === 'ONLINE' && bars.some((b: BarItem) => b.status === 'operational') ? 100 : null);

	// Dynamic label resolution
	$: defaultUnitLabel = range === '1d'
		? '48 bars · 30 mins / bar'
		: range === '7d'
		? '56 bars · 3 hrs / bar'
		: '30 bars · 1 bar / day';

	$: effectiveUnitLabel = timeUnitLabel || defaultUnitLabel;

	$: startLabel = range === '1d' ? '24 hours ago' : range === '7d' ? '7 days ago' : '30 days ago';
	$: endLabel = range === '1d' ? 'Now' : 'Today';

	function handleRangeClick(newRange: '1d' | '7d' | '30d') {
		range = newRange;
		if (onRangeChange) {
			onRangeChange(newRange);
		}
	}

	function getBarColor(bar: BarItem): string {
		if (bar.status === 'no_data' || bar.uptimePercent == null) {
			return 'bg-neutral-200 dark:bg-neutral-800 hover:opacity-80';
		}
		if (bar.uptimePercent >= 99) {
			return 'bg-[#2E8555] dark:bg-emerald-500 hover:opacity-80';
		}
		if (bar.uptimePercent >= 85) {
			// Minor degradation: vibrant lime/emerald
			return 'bg-emerald-400 dark:bg-emerald-400/80 hover:opacity-80';
		}
		if (bar.uptimePercent >= 60) {
			// Moderate degradation: yellow/amber
			return 'bg-amber-400 dark:bg-amber-400 hover:opacity-80';
		}
		if (bar.uptimePercent >= 30) {
			// Significant degradation: warm orange
			return 'bg-orange-500 dark:bg-orange-400 hover:opacity-80';
		}
		if (bar.uptimePercent > 0) {
			// Severe degradation: coral/red-orange
			return 'bg-rose-400 dark:bg-rose-400 hover:opacity-80';
		}
		// Complete outage (0%)
		return 'bg-rose-600 dark:bg-rose-500 hover:opacity-80';
	}

	function handleMouseEnter(e: MouseEvent, bar: BarItem) {
		activeHover = bar;
		const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
		hoverPos = {
			x: rect.left + rect.width / 2,
			y: rect.top - 8
		};
	}

	function handleMouseLeave() {
		activeHover = null;
	}
</script>

{#if !hasAnyTelemetry}
	<div class="p-3.5 rounded-md bg-(--bg-muted)/40 border border-(--border-hairline) flex items-start gap-2.5 text-xs text-(--text-secondary)">
		<Clock size={16} class="shrink-0 text-(--text-muted) mt-0.5" />
		<div class="flex flex-col gap-0.5">
			<span class="font-medium text-(--text-main)">Telemetry pending</span>
			<span class="text-(--text-muted) leading-relaxed">
				No availability records yet. Telemetry probes will begin recording once the endpoint is active.
			</span>
		</div>
	</div>
{:else}
	<div class="flex flex-col gap-2 w-full select-none">
		<!-- Top Row: Metric Header & Optional Range Controls -->
		<div class="flex items-center justify-between text-xs gap-3 flex-wrap">
			<div class="flex items-baseline gap-2">
				<span class="font-normal text-(--text-muted)">
					{range === '1d' ? '24-hour' : range === '7d' ? '7-day' : '30-day'} availability
				</span>
				<span class="font-mono font-medium text-(--text-main)">
					{currentUptime != null ? `${currentUptime.toFixed(1)}%` : '—'}
				</span>
			</div>

			{#if showRangeSelector || (variant === 'detailed' && availability?.ranges)}
				<div class="inline-flex items-center p-0.5 rounded-sm bg-(--bg-muted) border border-(--border-hairline) text-xs">
					<button
						type="button"
						class="px-2 py-0.5 rounded-xs font-mono transition-colors cursor-pointer {range === '1d' ? 'bg-(--bg-surface) text-(--text-main) font-medium shadow-xs' : 'text-(--text-muted) hover:text-(--text-main)'}"
						on:click={() => handleRangeClick('1d')}
					>
						24h
					</button>
					<button
						type="button"
						class="px-2 py-0.5 rounded-xs font-mono transition-colors cursor-pointer {range === '7d' ? 'bg-(--bg-surface) text-(--text-main) font-medium shadow-xs' : 'text-(--text-muted) hover:text-(--text-main)'}"
						on:click={() => handleRangeClick('7d')}
					>
						7d
					</button>
					<button
						type="button"
						class="px-2 py-0.5 rounded-xs font-mono transition-colors cursor-pointer {range === '30d' ? 'bg-(--bg-surface) text-(--text-main) font-medium shadow-xs' : 'text-(--text-muted) hover:text-(--text-main)'}"
						on:click={() => handleRangeClick('30d')}
					>
						30d
					</button>
				</div>
			{/if}
		</div>

		{#if !hasCurrentRangeData}
			<div class="p-3 rounded-md bg-(--bg-muted)/30 border border-(--border-hairline) text-xs text-(--text-muted) flex items-center justify-between">
				<span>No probes recorded in the last {range === '1d' ? '24 hours' : range === '7d' ? '7 days' : '30 days'}.</span>
				<span class="font-mono text-xs opacity-60">{startLabel} → {endLabel}</span>
			</div>
		{:else}
			<!-- Vertical Bar Sequence -->
			<div
				class="flex items-center gap-[2px] sm:gap-[3px] w-full py-0.5"
				role="region"
				aria-label="{range} availability history"
			>
				{#each bars as bar, idx (bar.date + '-' + idx)}
					<div
						class="flex-1 min-w-[2px] rounded-[1.5px] transition-all duration-150 cursor-pointer {getBarColor(bar)} {variant === 'compact' ? 'h-4 sm:h-5' : 'h-6 sm:h-7'}"
						role="button"
						tabindex="0"
						aria-label="{bar.label}: {bar.details}"
						on:mouseenter={(e) => handleMouseEnter(e, bar)}
						on:mouseleave={handleMouseLeave}
					></div>
				{/each}
			</div>

			<!-- Time Bounds Footer -->
			<div class="flex items-center justify-between text-xs text-(--text-muted) font-mono pt-0.5 gap-2 flex-wrap">
				<span>{startLabel}</span>
				<span class="opacity-75">{effectiveUnitLabel}</span>
				<span>{endLabel}</span>
			</div>
		{/if}

		<!-- Detailed Variant Legend -->
		{#if variant === 'detailed'}
			<div class="flex items-center justify-between pt-2.5 mt-1 border-t border-(--border-hairline) text-xs text-(--text-secondary) flex-wrap gap-2">
				<div class="flex items-center gap-3.5 flex-wrap">
					<div class="flex items-center gap-1.5">
						<span class="w-2.5 h-2.5 rounded-[2px] bg-[#2E8555] dark:bg-emerald-500"></span>
						<span class="text-xs">Operational</span>
					</div>
					<div class="flex items-center gap-1.5">
						<span class="w-4 h-2.5 rounded-[2px] bg-gradient-to-r from-emerald-400 via-amber-400 to-rose-400"></span>
						<span class="text-xs">Degraded (Graduated)</span>
					</div>
					<div class="flex items-center gap-1.5">
						<span class="w-2.5 h-2.5 rounded-[2px] bg-rose-600 dark:bg-rose-500"></span>
						<span class="text-xs">Downtime</span>
					</div>
					<div class="flex items-center gap-1.5">
						<span class="w-2.5 h-2.5 rounded-[2px] bg-neutral-200 dark:bg-neutral-800"></span>
						<span class="text-xs">No data</span>
					</div>
				</div>

				<span class="text-xs text-(--text-muted) font-mono">
					Periodic 5m probe
				</span>
			</div>
		{/if}

		<!-- Floating Tooltip -->
		{#if activeHover}
			<div
				class="fixed z-50 pointer-events-none transform -translate-x-1/2 -translate-y-full px-3 py-1.5 rounded-md text-xs shadow-md border backdrop-blur-md transition-opacity duration-100 bg-neutral-900/95 text-white border-neutral-800 dark:bg-neutral-950/95 dark:border-neutral-700"
				style="left: {hoverPos.x}px; top: {hoverPos.y}px;"
			>
				<div class="font-medium text-neutral-100 font-sans">{activeHover.label}</div>
				<div class="text-xs text-neutral-300 font-mono mt-0.5">
					{activeHover.details}
				</div>
			</div>
		{/if}
	</div>
{/if}
