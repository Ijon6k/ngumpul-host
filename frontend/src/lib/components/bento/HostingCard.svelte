<script lang="ts">
	// 3-step request flow
	const steps = [
		{ num: '1', label: 'Dockerize or build' },
		{ num: '2', label: 'Submit request' },
		{ num: '3', label: 'Get port & domain' }
	];
</script>

<div class="h-full w-full bg-[#406882] text-white rounded-3xl p-5 sm:p-6 flex flex-col justify-between shadow-xs min-h-[140px] relative overflow-hidden">
	<!-- Background subtle architectural circles -->
	<div class="absolute -right-6 -bottom-6 w-32 h-32 rounded-full border border-white/10 pointer-events-none"></div>

	<div class="relative z-10 flex items-center justify-between">
		<div>
			<div class="text-xs uppercase tracking-widest text-white/70 font-mono">Hosting</div>
			<div class="text-sm font-normal text-white mt-0.5">3-step request flow</div>
		</div>

		<a
			href="/register"
			class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg bg-white text-neutral-900 text-xs font-medium hover:bg-neutral-100 transition-all active:scale-[0.98] shadow-xs shrink-0"
		>
			Request a spot →
		</a>
	</div>

	<!-- Steps Row -->
	<div class="relative z-10 grid grid-cols-3 gap-2 pt-3 border-t border-white/15">
		{#each steps as step}
			<div>
				<div class="text-[11px] font-mono text-white/60">0{step.num}</div>
				<div class="text-xs text-white/90 font-normal mt-0.5 leading-snug">{step.label}</div>
			</div>
		{/each}
	</div>
</div>
