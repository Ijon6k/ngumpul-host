<script lang="ts">
	export let hostSpecs: any = null;
</script>

<div class="relative w-full h-full min-h-[340px] sm:min-h-[400px] rounded-2xl overflow-hidden border border-(--border-hairline) bg-(--bg-surface) group">
	<!-- Minimalist Pinterest-aesthetic 3D Server Unit Product Image -->
	<img
		src="/assets/minimalist_server_rack.jpg"
		alt="Minimalist homelab host device"
		class="w-full h-full object-cover object-center group-hover:scale-[1.01] transition-transform duration-500 ease-out"
	/>

	<!-- Subtle gradient scrim at bottom and top for readable tags -->
	<div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-black/30 pointer-events-none"></div>

	<!-- Top Badging: Host Identity & Status -->
	<div class="absolute top-4 left-4 right-4 flex items-center justify-between pointer-events-none">
		<div class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-black/40 backdrop-blur-md border border-white/10 text-white text-xs font-mono">
			<span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
			<span class="font-medium tracking-tight">{hostSpecs?.hostname ?? 'ngumpul-host'}</span>
		</div>

		<span class="px-2.5 py-1 rounded-md bg-white/20 backdrop-blur-md border border-white/10 text-xs font-mono text-white/90">
			{hostSpecs?.arch ?? 'amd64'} · {hostSpecs?.disk_type ?? 'NVMe SSD'}
		</span>
	</div>

	<!-- Bottom Badging: Quiet Hardware Profile Caption -->
	<div class="absolute bottom-4 left-4 right-4 flex items-end justify-between pointer-events-none text-white">
		<div class="flex flex-col">
			<span class="text-xs font-semibold text-white/90">Homelab node</span>
			<span class="font-display font-medium text-sm sm:text-base text-white drop-shadow-xs">
				Bare-metal Linux host for community projects
			</span>
		</div>

		<div class="text-right font-mono text-xs text-white/80">
			<span>{hostSpecs?.uptime_formatted ?? '3 hours'} uptime</span>
		</div>
	</div>
</div>
