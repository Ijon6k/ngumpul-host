<script lang="ts">
	export let src: string | null | undefined = undefined;
	export let name: string = 'Project';
	export let alt: string | undefined = undefined;
	export let aspectRatio: string = '16/9';
	let className: string = '';

	export { className as class };

	let imageError = false;

	$: isFull = aspectRatio === 'full' || aspectRatio === '100%' || aspectRatio === 'auto';

	$: cleanAspect = isFull
		? ''
		: aspectRatio.startsWith('aspect-')
			? aspectRatio
			: aspectRatio.includes('/')
				? `aspect-[${aspectRatio}]`
				: `aspect-${aspectRatio}`;

	$: inlineAspect = !isFull && aspectRatio && !aspectRatio.startsWith('aspect-') && aspectRatio.includes('/')
		? `aspect-ratio: ${aspectRatio};`
		: '';

	$: initials = (name || 'P')
		.split(/\s+/)
		.map((part) => part[0])
		.slice(0, 2)
		.join('')
		.toUpperCase();
</script>

<div
	class="relative w-full {isFull ? 'h-full' : cleanAspect} overflow-hidden rounded-[6px] select-none {className}"
	style={inlineAspect}
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
		<!-- Calm Pastel Blue Fallback Artwork (Guaranteed Full-Bleed) -->
		<div
			class="w-full h-full min-h-full flex flex-col items-center justify-center p-6 text-center bg-sky-100/90 dark:bg-sky-950/50 text-sky-900 dark:text-sky-100 transition-colors"
		>
			<div
				class="w-11 h-11 mb-2.5 rounded-sm bg-sky-200/80 dark:bg-sky-900/60 border border-sky-300/70 dark:border-sky-800/60 flex items-center justify-center font-mono text-sm font-semibold text-sky-800 dark:text-sky-200 shadow-xs"
			>
				{initials}
			</div>
			<p class="font-sans font-medium text-xs sm:text-sm text-sky-900 dark:text-sky-100 line-clamp-1 max-w-[80%]">
				{name}
			</p>
		</div>
	{/if}

	<slot />
</div>
