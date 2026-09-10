<script lang="ts">
	export let text = '04';
	export let dotSize = 3; // radius of dot
	export let gap = 3; // gap between dots
	export let color = 'currentColor';
	export let showInactive = false;
	export let inactiveOpacity = 0.08;
	export let charSpacing = 5;

	// 5x7 dot-matrix bitmap definitions
	const BITMAPS: Record<string, number[]> = {
		'0': [0b01110, 0b10001, 0b10011, 0b10101, 0b11001, 0b10001, 0b01110],
		'1': [0b00100, 0b01100, 0b00100, 0b00100, 0b00100, 0b00100, 0b01110],
		'2': [0b01110, 0b10001, 0b00001, 0b00010, 0b00100, 0b01000, 0b11111],
		'3': [0b11110, 0b00001, 0b00001, 0b01110, 0b00001, 0b00001, 0b11110],
		'4': [0b00010, 0b00110, 0b01010, 0b10010, 0b11111, 0b00010, 0b00010],
		'5': [0b11111, 0b10000, 0b11110, 0b00001, 0b00001, 0b10001, 0b01110],
		'6': [0b00110, 0b01000, 0b10000, 0b11110, 0b10001, 0b10001, 0b01110],
		'7': [0b11111, 0b00001, 0b00010, 0b00100, 0b01000, 0b01000, 0b01000],
		'8': [0b01110, 0b10001, 0b10001, 0b01110, 0b10001, 0b10001, 0b01110],
		'9': [0b01110, 0b10001, 0b10001, 0b01111, 0b00001, 0b00010, 0b01100],
		':': [0b00000, 0b00100, 0b00000, 0b00000, 0b00100, 0b00000, 0b00000],
		'/': [0b00001, 0b00010, 0b00100, 0b01000, 0b10000, 0b00000, 0b00000],
		'.': [0b00000, 0b00000, 0b00000, 0b00000, 0b00000, 0b00100, 0b00100],
		'A': [0b01110, 0b10001, 0b10001, 0b11111, 0b10001, 0b10001, 0b10001],
		'B': [0b11110, 0b10001, 0b10001, 0b11110, 0b10001, 0b10001, 0b11110],
		'C': [0b01110, 0b10001, 0b10000, 0b10000, 0b10000, 0b10001, 0b01110],
		'D': [0b11110, 0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b11110],
		'E': [0b11111, 0b10000, 0b10000, 0b11110, 0b10000, 0b10000, 0b11111],
		'F': [0b11111, 0b10000, 0b10000, 0b11110, 0b10000, 0b10000, 0b10000],
		'H': [0b10001, 0b10001, 0b10001, 0b11111, 0b10001, 0b10001, 0b10001],
		'I': [0b01110, 0b00100, 0b00100, 0b00100, 0b00100, 0b00100, 0b01110],
		'L': [0b10000, 0b10000, 0b10000, 0b10000, 0b10000, 0b10000, 0b11111],
		'N': [0b10001, 0b11001, 0b10101, 0b10011, 0b10001, 0b10001, 0b10001],
		'O': [0b01110, 0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b01110],
		'P': [0b11110, 0b10001, 0b10001, 0b11110, 0b10000, 0b10000, 0b10000],
		'R': [0b11110, 0b10001, 0b10001, 0b11110, 0b10100, 0b10010, 0b10001],
		'S': [0b01111, 0b10000, 0b10000, 0b01110, 0b00001, 0b00001, 0b11110],
		'T': [0b11111, 0b00100, 0b00100, 0b00100, 0b00100, 0b00100, 0b00100],
		'U': [0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b01110],
		'V': [0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b01010, 0b00100],
		'W': [0b10001, 0b10001, 0b10001, 0b10101, 0b10101, 0b11011, 0b01010],
		' ': [0b00000, 0b00000, 0b00000, 0b00000, 0b00000, 0b00000, 0b00000]
	};

	const COLS = 5;
	const ROWS = 7;
	const PITCH = dotSize * 2 + gap;

	$: chars = text.toUpperCase().split('');
	$: totalWidth = chars.length > 0 ? chars.length * (COLS * PITCH) + (chars.length - 1) * charSpacing : 0;
	$: totalHeight = ROWS * PITCH;
</script>

<svg
	width={totalWidth}
	height={totalHeight}
	viewBox="0 0 {totalWidth} {totalHeight}"
	class="inline-block shrink-0 select-none"
	aria-label={text}
	role="img"
>
	{#each chars as char, charIdx}
		{@const bitmap = BITMAPS[char] || BITMAPS[' ']}
		{@const charOffsetX = charIdx * (COLS * PITCH + charSpacing)}

		{#each Array(ROWS) as _, rowIdx}
			{@const rowBits = bitmap[rowIdx] || 0}
			{#each Array(COLS) as _, colIdx}
				{@const bitMask = 1 << (COLS - 1 - colIdx)}
				{@const isLit = (rowBits & bitMask) !== 0}
				{@const cx = charOffsetX + colIdx * PITCH + dotSize}
				{@const cy = rowIdx * PITCH + dotSize}

				{#if isLit}
					<circle {cx} {cy} r={dotSize} fill={color} />
				{:else if showInactive}
					<circle {cx} {cy} r={dotSize * 0.75} fill={color} opacity={inactiveOpacity} />
				{/if}
			{/each}
		{/each}
	{/each}
</svg>
