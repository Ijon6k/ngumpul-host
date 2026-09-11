<script lang="ts">
	export let src: string | null | undefined = undefined;
	export let name: string = 'Project';
	export let alt: string | undefined = undefined;
	export let aspectRatio: string = 'aspect-[4/3]';
	let className: string = '';
	export { className as class };

	let imageError = false;

	// Deterministic palette generator for consistent, intentional fallback styling
	function getDeterministicPalette(str: string) {
		let hash = 0;
		for (let i = 0; i < str.length; i++) {
			hash = str.charCodeAt(i) + ((hash << 5) - hash);
		}
		const palettes = [
			{ bg: 'bg-slate-800 dark:bg-slate-900', text: 'text-slate-100', subtext: 'text-slate-400', border: 'border-slate-700/60' },
			{ bg: 'bg-zinc-800 dark:bg-zinc-900', text: 'text-zinc-100', subtext: 'text-zinc-400', border: 'border-zinc-700/60' },
			{ bg: 'bg-stone-800 dark:bg-stone-900', text: 'text-stone-100', subtext: 'text-stone-400', border: 'border-stone-700/60' },
			{ bg: 'bg-cyan-950 dark:bg-cyan-950', text: 'text-cyan-100', subtext: 'text-cyan-400', border: 'border-cyan-800/60' },
			{ bg: 'bg-teal-950 dark:bg-teal-950', text: 'text-teal-100', subtext: 'text-teal-400', border: 'border-teal-800/60' },
			{ bg: 'bg-emerald-950 dark:bg-emerald-950', text: 'text-emerald-100', subtext: 'text-emerald-400', border: 'border-emerald-800/60' },
			{ bg: 'bg-indigo-950 dark:bg-indigo-950', text: 'text-indigo-100', subtext: 'text-indigo-400', border: 'border-indigo-800/60' }
		];
		const index = Math.abs(hash) % palettes.length;
		return palettes[index];
	}

	$: palette = getDeterministicPalette(name || 'Project');
	$: initials = (name || 'P')
		.split(/\s+/)
		.map((part) => part[0])
		.slice(0, 2)
		.join('')
		.toUpperCase();
</script>

<div
	class="relative w-full {aspectRatio} overflow-hidden rounded-[6px] border border-(--border-hairline) bg-(--bg-muted) select-none {className}"
>
	{#if src && !imageError}
		<img
			{src}
			alt={alt || name}
			loading="lazy"
			class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-[1.01]"
			on:error={() => (imageError = true)}
		/>
	{:else}
		<!-- Deterministic Editorial Fallback Artwork -->
		<div class="w-full h-full flex flex-col items-center justify-center p-6 text-center {palette.bg} transition-colors">
			<div
				class="w-12 h-12 mb-3 rounded-[4px] border {palette.border} bg-white/5 flex items-center justify-center {palette.text} font-mono text-base font-semibold tracking-wider shadow-xs"
			>
				{initials}
			</div>
			<p class="font-sans font-medium text-sm sm:text-base {palette.text} line-clamp-1 max-w-[85%]">
				{name}
			</p>
			<span class="text-[11px] {palette.subtext} font-mono mt-0.5">
				ngumpul.host
			</span>
		</div>
	{/if}

	<slot />
</div>
