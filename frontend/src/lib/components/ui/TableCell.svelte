<script lang="ts">
	export let align: 'left' | 'center' | 'right' = 'left';
	export let alignTop: boolean = false;
	export let mono: boolean = false;
	export let colSpan: number | undefined = undefined;
	let className: string = '';
	export { className as class };

	$: alignClass = align === 'right' ? 'text-right' : align === 'center' ? 'text-center' : 'text-left';
	$: verticalAlign = alignTop ? 'align-top' : 'align-middle';
	$: fontClass = mono ? 'font-mono text-xs sm:text-sm' : 'font-sans text-sm';
</script>

<td colspan={colSpan} class="px-5 py-3.5 {verticalAlign} {alignClass} {fontClass} text-(--text-main) {className}">
	<slot />
</td>
