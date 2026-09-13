<script lang="ts">
	import MetricPillar from './MetricPillar.svelte';
	import { ArrowRight } from 'phosphor-svelte';

	export let hostSpecs: any = null;
</script>

<section class="relative min-h-screen flex flex-col justify-between overflow-hidden bg-[#0A0D14] text-white">
	<!-- Ambient Background with heroserver.webp in natural lighting -->
	<div class="absolute inset-0 z-0 overflow-hidden">
		<!-- User-specified hero server photography: bright and crisp, no excessive dimming -->
		<img
			src="/assets/heroserver.webp"
			alt="Homelab Server Infrastructure"
			class="w-full h-full object-cover object-center brightness-100 opacity-95"
		/>

		<!-- Minimalist gentle directional vignette for comfortable text legibility only on left -->
		<div class="absolute inset-0 bg-gradient-to-r from-[#070A0E]/75 via-[#070A0E]/35 to-transparent pointer-events-none"></div>
		<div class="absolute inset-x-0 bottom-0 h-24 bg-gradient-to-t from-[#070A0E] to-transparent pointer-events-none"></div>
	</div>

	<!-- Top spacer for fixed header clearance -->
	<div class="h-20 shrink-0"></div>

	<!-- Center / Hero Content -->
	<div class="relative z-10 container mx-auto px-6 max-w-6xl py-12 md:py-16 my-auto">
		<div class="max-w-3xl flex flex-col gap-6">

			<!-- Headline: Minimalist Apple/Samsung editorial typography, non-bold, brand sky blue highlight -->
			<h1
				class="font-sans font-normal text-white tracking-[-0.025em] leading-[1.22] text-2xl sm:text-3xl md:text-4xl lg:text-[2.65rem] max-w-2xl"
			>
				A quiet corner of the internet for
				<span class="text-[#79AFC4]">side projects, experiments,</span>
				and friends.
			</h1>

			<!-- Sub-headline: Honest, casual, non-corporate, clear and grounded -->
			<p class="text-base sm:text-lg text-neutral-300/90 leading-relaxed max-w-xl font-normal">
				A community-hosted server running quietly on real hardware. Kept online for side projects, hobby bots, and creative ideas—without cloud subscription traps or corporate noise.
			</p>

			<!-- Action CTAs: Refined, un-bolded -->
			<div class="flex items-center gap-3.5 pt-1 flex-wrap">
				<a
					href="/projects"
					class="px-5 py-2.5 text-sm font-medium rounded-lg bg-white text-neutral-950 hover:bg-neutral-100 transition-all shadow-sm active:scale-[0.98] inline-flex items-center gap-1.5"
				>
					<span>Explore projects</span>
					<ArrowRight size={14} weight="bold" />
				</a>
				<a
					href="/register"
					class="px-5 py-2.5 text-sm font-medium rounded-lg text-neutral-200 hover:text-white bg-white/[0.06] hover:bg-white/[0.1] border border-white/15 backdrop-blur-sm transition-all active:scale-[0.98]"
				>
					Request a spot
				</a>
			</div>
		</div>
	</div>

	<!-- Bottom Dynamic Metric Pillars (Direct Host Data) -->
	<div class="relative z-10 border-t border-white/10 bg-[#070A0E]/80 backdrop-blur-md">
		<div class="container mx-auto px-6 max-w-6xl py-6 grid grid-cols-2 md:grid-cols-4 gap-6 divide-y md:divide-y-0 md:divide-x divide-white/10">
			<!-- Pillar 1: CPU Cores -->
			<div class="pt-4 md:pt-0 md:pr-6">
				<MetricPillar
					value="{hostSpecs?.hardware?.cores ?? hostSpecs?.cpu_cores ?? 10} Cores"
					label="Dedicated CPU"
					subtext={hostSpecs?.hardware?.cpu ?? hostSpecs?.cpu_model ?? '13th Gen Intel Core i5-1334U'}
					class="text-white"
				/>
			</div>

			<!-- Pillar 2: System Memory -->
			<div class="pt-4 md:pt-0 md:px-6">
				<MetricPillar
					value="{hostSpecs?.memory?.totalGB ?? hostSpecs?.total_ram_gb ?? 23.1} GB"
					label="System Memory"
					subtext="{hostSpecs?.memory?.usagePercent ?? hostSpecs?.used_ram_percent ?? 68}% allocated to containers"
					class="text-white"
				/>
			</div>

			<!-- Pillar 3: Storage -->
			<div class="pt-4 md:pt-0 md:px-6">
				<MetricPillar
					value="{hostSpecs?.storage?.totalGB ?? hostSpecs?.disk_total_gb ?? 98} GB"
					label="Linux Storage"
					subtext="{hostSpecs?.storage?.availableGB ?? 7.8} GB free · {hostSpecs?.storage?.type ?? 'NVMe SSD'}"
					class="text-white"
				/>
			</div>

			<!-- Pillar 4: Uptime -->
			<div class="pt-4 md:pt-0 md:pl-6">
				<MetricPillar
					value="{hostSpecs?.availability?.uptime ?? hostSpecs?.uptime_formatted ?? 'Running'}"
					label="System Uptime"
					subtext="Continuous host session without reboots"
					class="text-white"
				/>
			</div>
		</div>
	</div>
</section>
