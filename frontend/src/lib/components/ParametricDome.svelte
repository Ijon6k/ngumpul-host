<script lang="ts">
	import { onMount } from 'svelte';

	export let width = 310;
	export let height = 160;
	export const dotColor = 'rgba(255,255,255,0.75)';

	// === Drive activity state ===
	// 16 drives: 2 rows × 8 cols
	const DRIVES = 16;
	let driveActive = Array.from({ length: DRIVES }, (_, i) => i < 11);

	// === Network port state ===
	const PORTS = 10;
	let portActive = Array.from({ length: PORTS }, (_, i) => i < 6);

	// === CPU / MEM usage bars (0-100) ===
	let cpuPct = 42;
	let memPct = 61;

	// === LED blink ===
	let ledOn = true;
	let netLedOn = true;

	// === Activity LED (random drive write) ===
	let writeDrive = -1;

	onMount(() => {
		// Blink status LED
		const ledInterval = setInterval(() => {
			ledOn = !ledOn;
		}, 1400);

		// Blink net LED faster
		const netLedInterval = setInterval(() => {
			netLedOn = !netLedOn;
		}, 300);

		// Random drive activity
		const driveInterval = setInterval(() => {
			// Toggle a random drive
			const idx = Math.floor(Math.random() * DRIVES);
			driveActive[idx] = !driveActive[idx];
			driveActive = [...driveActive];
			writeDrive = idx;
		}, 600);

		// Fluctuate CPU usage
		const cpuInterval = setInterval(() => {
			cpuPct = Math.max(20, Math.min(85, cpuPct + (Math.random() - 0.5) * 12));
			memPct = Math.max(40, Math.min(78, memPct + (Math.random() - 0.5) * 6));
		}, 900);

		// Random port toggle
		const portInterval = setInterval(() => {
			const idx = Math.floor(Math.random() * PORTS);
			portActive[idx] = Math.random() > 0.3;
			portActive = [...portActive];
		}, 1200);

		return () => {
			clearInterval(ledInterval);
			clearInterval(netLedInterval);
			clearInterval(driveInterval);
			clearInterval(cpuInterval);
			clearInterval(portInterval);
		};
	});

	// Layout constants
	const PAD = 10;
	const W = width;
	const H = height;

	// Top status bar
	const BAR_H = 24;

	// Drive section: 2 rows × 8 cols
	const DRIVE_COLS = 8;
	const DRIVE_ROWS = 2;
	const DRIVE_W = 28;
	const DRIVE_H = 16;
	const DRIVE_GAP_X = 4;
	const DRIVE_GAP_Y = 4;
	const DRIVE_AREA_X = PAD + 2;
	const DRIVE_AREA_Y = BAR_H + 16;

	// Network ports row
	const PORT_W = 14;
	const PORT_H = 10;
	const PORT_GAP = 4;
	const PORT_Y = DRIVE_AREA_Y + DRIVE_ROWS * (DRIVE_H + DRIVE_GAP_Y) + 10;
</script>

<svg
	{width}
	{height}
	viewBox="0 0 {W} {H}"
	xmlns="http://www.w3.org/2000/svg"
	style="display:block; max-width:100%"
	aria-hidden="true"
>
	<!-- Chassis border -->
	<rect x="2" y="2" width={W - 4} height={H - 4} rx="4"
		fill="none" stroke="#1E2A35" stroke-width="1"/>

	<!-- ── Top Status Bar ── -->
	<rect x="2" y="2" width={W - 4} height={BAR_H} rx="4" fill="#0C1015"/>
	<!-- square bottom corners of top bar -->
	<rect x="2" y={BAR_H - 4} width={W - 4} height="6" fill="#0C1015"/>

	<!-- Power LED -->
	<circle cx="16" cy={BAR_H / 2 + 1} r="3.5"
		fill={ledOn ? '#22C55E' : '#0F2A1A'}/>

	<!-- Node label -->
	<text x="26" y="16"
		fill="#4A6070"
		font-size="8"
		font-family="ui-monospace, JetBrains Mono, monospace">
		NODE-01 · ngumpul.host · :1111
	</text>

	<!-- CPU usage bar -->
	<text x={W - 88} y="11"
		fill="#334D5C"
		font-size="6.5"
		font-family="ui-monospace, monospace">CPU</text>
	<rect x={W - 72} y="5" width="64" height="9" rx="1.5" fill="#0A1520"/>
	<rect x={W - 71} y="6" width={Math.round(cpuPct * 0.62)} height="7" rx="1"
		fill="#1E4A60"
		style="transition: width 0.7s ease"/>

	<!-- MEM usage bar -->
	<text x={W - 88} y="22"
		fill="#334D5C"
		font-size="6.5"
		font-family="ui-monospace, monospace">MEM</text>
	<rect x={W - 72} y="16" width="64" height="9" rx="1.5" fill="#0A1520"/>
	<rect x={W - 71} y="17" width={Math.round(memPct * 0.62)} height="7" rx="1"
		fill="#1A4060"
		style="transition: width 0.7s ease"/>

	<!-- ── Drive Bay Section ── -->
	<text x={DRIVE_AREA_X} y={DRIVE_AREA_Y - 5}
		fill="#2A3D4A"
		font-size="6.5"
		font-family="ui-monospace, monospace">SSD/NVMe ARRAY</text>

	{#each Array(DRIVE_ROWS) as _, row}
		{#each Array(DRIVE_COLS) as _, col}
			{@const idx = row * DRIVE_COLS + col}
			{@const dx = DRIVE_AREA_X + col * (DRIVE_W + DRIVE_GAP_X)}
			{@const dy = DRIVE_AREA_Y + row * (DRIVE_H + DRIVE_GAP_Y)}
			{@const active = driveActive[idx]}
			{@const isWriting = writeDrive === idx}

			<!-- Drive bay slot -->
			<rect x={dx} y={dy} width={DRIVE_W} height={DRIVE_H} rx="2"
				fill={active ? '#0A1F18' : '#080E14'}
				stroke={active ? (isWriting ? '#FF5500' : '#1A4232') : '#131E28'}
				stroke-width="0.8"/>

			<!-- Drive activity LED -->
			<circle cx={dx + 4} cy={dy + 8} r="2"
				fill={isWriting ? '#FF5500' : (active ? '#22C55E' : '#111D26')}/>

			<!-- Drive face lines -->
			<line x1={dx + 9} y1={dy + 5} x2={dx + DRIVE_W - 3} y2={dy + 5}
				stroke={active ? '#152E22' : '#0D1820'} stroke-width="0.8"/>
			<line x1={dx + 9} y1={dy + 8} x2={dx + DRIVE_W - 3} y2={dy + 8}
				stroke={active ? '#152E22' : '#0D1820'} stroke-width="0.8"/>
			<line x1={dx + 9} y1={dy + 11} x2={dx + DRIVE_W - 3} y2={dy + 11}
				stroke={active ? '#152E22' : '#0D1820'} stroke-width="0.8"/>
		{/each}
	{/each}

	<!-- ── Network Port Section ── -->
	<line x1={PAD} y1={PORT_Y - 5} x2={W - PAD} y2={PORT_Y - 5}
		stroke="#131E28" stroke-width="0.7"/>

	<text x={PAD + 2} y={PORT_Y + PORT_H - 1}
		fill="#2A3D4A"
		font-size="6.5"
		font-family="ui-monospace, monospace">NET</text>

	{#each Array(PORTS) as _, i}
		{@const px = PAD + 34 + i * (PORT_W + PORT_GAP)}
		{@const isUp = portActive[i]}

		<!-- Port housing -->
		<rect x={px} y={PORT_Y} width={PORT_W} height={PORT_H} rx="1.5"
			fill={isUp ? '#0A1F18' : '#080D14'}
			stroke={isUp ? '#1A4232' : '#131E28'}
			stroke-width="0.7"/>

		<!-- Port LED -->
		<circle cx={px + PORT_W / 2} cy={PORT_Y + PORT_H / 2} r="2"
			fill={isUp ? (netLedOn ? '#22C55E' : '#114A28') : '#0F1820'}/>
	{/each}

	<!-- Uptime label far right of port row -->
	<text x={W - PAD - 2} y={PORT_Y + PORT_H - 1}
		fill="#2A5060"
		font-size="6.5"
		font-family="ui-monospace, monospace"
		text-anchor="end">99.85% UP</text>
</svg>
