<script lang="ts">
	export let maxWidth: '4xl' | '5xl' | '6xl' | 'full' = '6xl';
	export let paddingY: string = 'py-8 sm:py-12 pb-24';

	let className: string = '';
	export { className as class };

	const maxWidthStyles = {
		'4xl': 'max-w-4xl',
		'5xl': 'max-w-5xl',
		'6xl': 'max-w-6xl',
		full: 'max-w-full'
	};
</script>

<div class="container mx-auto px-4 sm:px-6 {maxWidthStyles[maxWidth]} {paddingY} flex flex-col gap-8 sm:gap-10 {className}">
	<slot />
</div>
