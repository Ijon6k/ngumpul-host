<script lang="ts">
	export let location: string = 'Jakarta, Indonesia';
	export let status: string = 'online';
</script>

<div class="h-full w-full bg-[#0E1013] text-white rounded-3xl p-6 sm:p-7 flex flex-col justify-between relative overflow-hidden shadow-xs border border-transparent dark:border-(--border-hairline) group min-h-[240px] transition-colors duration-300">
	<!-- User-specified bentoserver.webp photograph -->
	<img
		src="/assets/bentoserver.webp"
		alt="Ngumpul Host Server Rack"
		class="absolute inset-0 w-full h-full object-cover object-center brightness-100 group-hover:scale-[1.015] transition-transform duration-700 ease-out"
	/>

	<!-- Minimal gentle gradient for text legibility without muddy darkness -->
	<div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/25 to-black/35 pointer-events-none"></div>

	<!-- Top: Clean text with location pin, zero chips, zero capslock monospace -->
	<div class="relative z-10 flex items-center justify-between">
		<div class="flex items-center gap-1.5 text-xs text-neutral-200 font-normal">
			<svg class="w-3.5 h-3.5 text-neutral-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
				<path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z" />
				<circle cx="12" cy="10" r="3" />
			</svg>
			<span>{location}</span>
			<span class="text-neutral-500">·</span>
			<span class="flex items-center gap-1.5 text-emerald-400 font-normal">
				<span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
				<span class="capitalize">{status}</span>
			</span>
		</div>

		<span class="text-xs text-neutral-400 font-normal">
			Bare metal
		</span>
	</div>

	<!-- Bottom Editorial Tagline -->
	<div class="relative z-10 pt-16">
		<div class="text-base sm:text-lg text-white font-normal leading-snug">
			Always online
		</div>
		<div class="text-xs sm:text-sm text-neutral-300 font-normal mt-0.5">
			for what we build and test next.
		</div>
	</div>
</div>
