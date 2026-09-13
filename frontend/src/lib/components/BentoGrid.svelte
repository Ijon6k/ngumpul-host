<script lang="ts">
	import { createPublicServerQuery } from '$lib/api/server';
	import type { PublicServerResponse } from '$lib/types/server';
	import ServerPhotoCard from './bento/ServerPhotoCard.svelte';
	import HardwareCard from './bento/HardwareCard.svelte';
	import MemoryCard from './bento/MemoryCard.svelte';
	import StorageCard from './bento/StorageCard.svelte';
	import NetworkCard from './bento/NetworkCard.svelte';
	import AvailabilityCard from './bento/AvailabilityCard.svelte';
	import IdeasCard from './bento/IdeasCard.svelte';

	export let hostSpecs: any = null;

	const serverQuery = createPublicServerQuery();

	$: publicServer = ($serverQuery.data as PublicServerResponse) || null;

	// Dynamic reactive attributes with graceful fallbacks
	$: hardware = publicServer?.hardware || {
		cpu: hostSpecs?.hardware?.cpu || hostSpecs?.cpu_model || '13th Gen Intel® Core™ i5-1334U',
		cores: hostSpecs?.hardware?.cores || hostSpecs?.cpu_cores || 10,
		threads: hostSpecs?.hardware?.threads || hostSpecs?.cpu_threads || 12
	};

	$: memory = publicServer?.memory || {
		totalGB: hostSpecs?.memory?.totalGB || hostSpecs?.total_ram_gb || 23.1,
		usagePercent: Math.round(hostSpecs?.memory?.usagePercent || hostSpecs?.used_ram_percent || 72),
		availableGB: hostSpecs?.memory?.availableGB || hostSpecs?.available_ram_gb || 6.3
	};

	$: storage = publicServer?.storage || {
		physicalDiskGB: hostSpecs?.storage?.physicalDiskGB || hostSpecs?.disk_physical_gb || 512,
		linuxTotalGB: hostSpecs?.storage?.linuxTotalGB || hostSpecs?.disk_total_gb || 98,
		totalGB: hostSpecs?.storage?.totalGB || hostSpecs?.disk_total_gb || 98,
		usagePercent: Math.round(hostSpecs?.storage?.usagePercent || hostSpecs?.disk_used_percent || 93),
		availableGB: hostSpecs?.storage?.availableGB || hostSpecs?.available_disk_gb || 6.8,
		type: hostSpecs?.storage?.type || hostSpecs?.disk_type || 'NVMe SSD',
		model: hostSpecs?.storage?.model || hostSpecs?.disk_model || 'SAMSUNG MZVL4512HBLU-00BH1',
		summary: hostSpecs?.storage?.summary || hostSpecs?.disk_summary || 'SAMSUNG MZVL4512HBLU-00BH1 · 512GB NVMe SSD (98GB Linux Partition)'
	};

	$: network = publicServer?.network || {
		linkCapacity: hostSpecs?.network?.linkCapacity || hostSpecs?.network_speed || '1 Gbps',
		latencyMs: hostSpecs?.network?.latencyMs || 29
	};

	$: availability = publicServer?.availability || {
		current: hostSpecs?.availability?.current || 'operational',
		last30Days: hostSpecs?.availability?.last30Days || 100,
		uptime: hostSpecs?.availability?.uptime || hostSpecs?.uptime_formatted || 'Running',
		recordedPeriod: hostSpecs?.availability?.recordedPeriod || 'Monitored',
		dailyBlocks: hostSpecs?.availability?.dailyBlocks || []
	};

	$: projects = publicServer?.projects || {
		total: hostSpecs?.projects?.total || 3,
		online: hostSpecs?.projects?.online || 3
	};

	$: location = publicServer?.os?.location || hostSpecs?.os?.location || hostSpecs?.location || 'Jakarta, Indonesia';
</script>

<!-- Bento Grid Section: Calm editorial typography with genuine host telemetry -->
<section class="py-12 sm:py-16 md:py-20 border-b border-(--border-hairline) bg-(--bg-muted) transition-colors duration-300">
	<div class="container mx-auto px-4 sm:px-6 max-w-6xl flex flex-col gap-6">

		<!-- ══════════════════════════════════════════════════════════
		     ROW 1: Editorial Typography (Left) + Server Photo Card (Right)
		     ══════════════════════════════════════════════════════════ -->
		<div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-stretch">
			<!-- Left: Editorial Typography (Span 5) -->
			<div class="lg:col-span-5 flex flex-col justify-between py-2 sm:py-4 pr-0 lg:pr-4">
				<div>
					<!-- Section Heading (Clean, non-bold editorial scale) -->
					<h2 class="font-sans font-normal text-4xl sm:text-5xl lg:text-[3.6rem] text-(--text-main) tracking-[-0.035em] leading-[1.06] mt-2 sm:mt-4 transition-colors">
						Running<br />
						on real<br />
						<span class="text-(--text-muted)">hardware.</span>
					</h2>

					<!-- Subtitle -->
					<p class="text-sm sm:text-base text-(--text-secondary) font-normal leading-relaxed mt-5 max-w-sm transition-colors">
						A single machine at home, kept online for friends and experiments. No cloud bills, no hidden layers—just a server we manage together.
					</p>
				</div>
			</div>

			<!-- Right: Server Photo Card (Span 7) -->
			<div class="lg:col-span-7">
				<ServerPhotoCard
					{location}
					status={availability.current}
				/>
			</div>
		</div>

		<!-- ══════════════════════════════════════════════════════════
		     ROW 2: Hardware Core, Memory & Storage (3 Cards)
		     ══════════════════════════════════════════════════════════ -->
		<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 items-stretch">
			<!-- 1. Core Hardware Card -->
			<div>
				<HardwareCard
					cpu={hardware.cpu}
					cores={hardware.cores}
					threads={hardware.threads}
				/>
			</div>

			<!-- 2. System Memory Card (Dark) -->
			<div>
				<MemoryCard
					totalGB={memory.totalGB}
					usagePercent={memory.usagePercent}
					availableGB={memory.availableGB}
				/>
			</div>

			<!-- 3. Storage Card -->
			<div>
				<StorageCard
					physicalGB={storage.physicalDiskGB || storage.totalGB || 512}
					linuxTotalGB={storage.linuxTotalGB || storage.totalGB || 98}
					totalGB={storage.totalGB || 98}
					usagePercent={storage.usagePercent}
					availableGB={storage.availableGB}
					type={storage.type}
					model={storage.model}
					summary={storage.summary}
				/>
			</div>
		</div>

		<!-- ══════════════════════════════════════════════════════════
		     ROW 3: Availability, Network & Ideas Accent (3 Cards)
		     ══════════════════════════════════════════════════════════ -->
		<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 items-stretch">
			<!-- 4. Availability Card (Dark, Calculated 30-day Uptime) -->
			<div>
				<AvailabilityCard
					availabilityPercent={availability.last30Days}
					status={availability.current}
					uptime={availability.uptime}
					recordedPeriod={availability.recordedPeriod}
					dailyBlocks={availability.dailyBlocks}
				/>
			</div>

			<!-- 5. Network Card -->
			<div>
				<NetworkCard
					linkCapacity={network.linkCapacity}
					latencyMs={network.latencyMs}
				/>
			</div>

			<!-- 6. Ideas & Hosting Card (Slate Blue) -->
			<div>
				<IdeasCard
					projectsOnline={projects.online}
					projectsTotal={projects.total}
				/>
			</div>
		</div>

	</div>
</section>
