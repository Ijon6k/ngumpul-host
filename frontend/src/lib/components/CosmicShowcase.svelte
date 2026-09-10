<script lang="ts">
	// The cosmic showcase section – modeled on the reference "Every project is a fusion..." card.
	// Metrics here represent real server infrastructure facts.
	const metrics = [
		{
			value: '99.98%',
			label: 'Availability SLA',
			sub: 'Continuous monitoring · Jakarta IXP direct peering',
			arrow: true,
		},
		{
			value: '14ms',
			label: 'Edge Latency',
			sub: 'Average round-trip to Jakarta gateway node',
			arrow: true,
		},
		{
			value: '3.2 GB/s',
			label: 'Storage I/O',
			sub: 'NVMe PCIe Gen 4 throughput · No throttling',
			arrow: true,
		},
	];
</script>

<section class="relative overflow-hidden rounded-2xl" style="background:#060B10;">
	<!-- Cosmic background image -->
	<div class="absolute inset-0">
		<img
			src="/assets/cosmic_horizon.jpg"
			alt=""
			aria-hidden="true"
			class="w-full h-full object-cover object-bottom opacity-55"
			loading="lazy"
		/>
		<!-- Overlay gradient for text legibility -->
		<div class="absolute inset-0" style="background: linear-gradient(to bottom, rgba(6,11,16,0.15) 0%, rgba(6,11,16,0.45) 40%, rgba(6,11,16,0.88) 75%, #060B10 100%);"></div>
		<!-- Left vignette -->
		<div class="absolute inset-0" style="background: linear-gradient(to right, rgba(6,11,16,0.75) 0%, transparent 55%);"></div>
	</div>

	<div class="relative z-10 container mx-auto px-8 sm:px-10 max-w-6xl">
		<!-- Top pad + eyebrow -->
		<div class="pt-14 sm:pt-16">
			<span class="inline-flex items-center gap-2 text-[11px] font-mono tracking-widest uppercase text-white/40">
				<span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-live-pulse inline-block"></span>
				ngumpul.host · Independent Node · Jakarta, ID
			</span>
		</div>

		<!-- Headline -->
		<div class="mt-5 max-w-2xl">
			<h2 class="text-white font-display font-bold tracking-tight" style="font-size: clamp(2rem, 4.5vw, 3.2rem); line-height: 1.1; letter-spacing: -0.035em;">
				Every project is a fusion of
				<em class="italic" style="color: rgba(180,220,240,0.85);">craft</em>,
				resilience, and curiosity.
			</h2>
		</div>

		<!-- Subtitle -->
		<p class="mt-4 max-w-xl text-base text-white/45 leading-relaxed">
			A quiet, independent Linux host. Maintained by humans, reviewed by humans,
			run on real hardware — not an abstracted cloud scaffold.
		</p>

		<!-- Metrics row with vertical pill arrows -->
		<div class="mt-10 sm:mt-12 pb-12 sm:pb-14 flex flex-col sm:flex-row sm:items-end gap-0">
			{#each metrics as metric, i}
				<div
					class="flex flex-col gap-2 py-6 sm:py-0 sm:pr-10 border-b border-white/8 sm:border-b-0 sm:border-l sm:pl-10"
					class:sm:border-l-0={i === 0}
					class:sm:pl-0={i === 0}
				>
					<!-- Large number -->
					<div class="flex items-end gap-4">
						<span
							class="font-display font-bold text-white"
							style="font-size: clamp(2.2rem, 5vw, 3.5rem); line-height: 1; letter-spacing: -0.04em;"
						>{metric.value}</span>
						{#if metric.arrow}
							<a
								href="/status"
								class="btn-arrow mb-1 shrink-0"
								aria-label="View status details"
								tabindex="-1"
							>
								<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
									<path d="M5 12h14M12 5l7 7-7 7" />
								</svg>
							</a>
						{/if}
					</div>

					<!-- Label -->
					<span class="text-sm font-semibold text-white/70">{metric.label}</span>

					<!-- Sub -->
					<span class="text-xs text-white/35 leading-relaxed max-w-[14rem]">{metric.sub}</span>
				</div>
			{/each}
		</div>
	</div>
</section>
