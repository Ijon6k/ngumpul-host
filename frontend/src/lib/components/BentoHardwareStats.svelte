<script lang="ts">
	import BentoCard from './BentoCard.svelte';
	import ResourceBar from './ResourceBar.svelte';

	export let hostSpecs: any = null;
</script>

<BentoCard title="Detected Host Hardware" badge="Live Telemetry" badgeClass="text-(--cf-blue)">
	<div class="flex flex-col gap-4 text-xs">
		<!-- Processor -->
		<div class="flex flex-col gap-1">
			<span class="text-(--text-secondary) font-medium">Processor</span>
			<span class="font-semibold text-sm text-(--text-main) leading-snug">
				{hostSpecs?.cpu_model ?? '13th Gen Intel Core i5-1334U'}
			</span>
			<span class="text-neutral-400 font-mono text-[11px]">{hostSpecs?.cpu_cores ?? 12} physical threads</span>
		</div>

		<!-- Memory Gauge -->
		<div class="pt-2 border-t border-(--border-hairline)">
			<ResourceBar
				label="Memory (RAM)"
				valueText="{hostSpecs?.total_ram_gb ?? 23.1} GB ({hostSpecs?.used_ram_percent ?? 38.4}%)"
				percent={hostSpecs?.used_ram_percent ?? 38.4}
				colorClass="bg-(--cf-blue)"
				subText="{hostSpecs?.available_ram_gb ?? 14.2} GB available for container workloads"
			/>
		</div>

		<!-- Storage Gauge -->
		<div class="pt-2 border-t border-(--border-hairline)">
			<ResourceBar
				label="NVMe Solid State Storage"
				valueText="{hostSpecs?.disk_total_gb ?? 98} GB ({hostSpecs?.disk_used_percent ?? 86.7}%)"
				percent={hostSpecs?.disk_used_percent ?? 86.7}
				colorClass="bg-(--accent-sky)"
				subText="Direct PCIe local flash storage"
			/>
		</div>

		<!-- Architecture & OS -->
		<div class="pt-2 border-t border-(--border-hairline) flex items-center justify-between text-[11px] text-(--text-muted)">
			<span>Architecture</span>
			<span class="font-mono">{hostSpecs?.arch ?? 'amd64'} · {hostSpecs?.os ?? 'linux'}</span>
		</div>
	</div>

	<svelte:fragment slot="footer">
		<a href="/status" class="font-semibold text-(--cf-blue) hover:underline flex items-center justify-between">
			<span>Detailed health diagnostics</span>
			<span>→</span>
		</a>
	</svelte:fragment>
</BentoCard>
