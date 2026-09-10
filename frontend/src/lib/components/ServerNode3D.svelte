<script lang="ts">
	import { onMount } from 'svelte';

	export let cpuPct = 42;
	export let memPct = 61;
	export let activeContainers = 18;
	export let networkIn = 340; // MB/s
	export let networkOut = 128; // MB/s

	let ledOn = true;
	let packet1Phase = 0;
	let packet2Phase = 0.33;
	let packet3Phase = 0.66;

	onMount(() => {
		const ledInterval = setInterval(() => {
			ledOn = !ledOn;
		}, 900);

		// Animate packets
		const packetInterval = setInterval(() => {
			packet1Phase = (packet1Phase + 0.022) % 1;
			packet2Phase = (packet2Phase + 0.018) % 1;
			packet3Phase = (packet3Phase + 0.025) % 1;
		}, 30);

		// Animate CPU/MEM
		const metricsInterval = setInterval(() => {
			cpuPct = Math.max(22, Math.min(78, cpuPct + (Math.random() - 0.5) * 8));
			memPct = Math.max(44, Math.min(72, memPct + (Math.random() - 0.5) * 4));
			networkIn = Math.max(180, Math.min(580, networkIn + (Math.random() - 0.5) * 60));
		}, 1200);

		return () => {
			clearInterval(ledInterval);
			clearInterval(packetInterval);
			clearInterval(metricsInterval);
		};
	});

	// Bezier path for packet animation
	function getPacketPos(t: number, x1: number, y1: number, x2: number, y2: number, cx: number, cy: number) {
		const mt = 1 - t;
		return {
			x: mt * mt * x1 + 2 * mt * t * cx + t * t * x2,
			y: mt * mt * y1 + 2 * mt * t * cy + t * t * y2,
		};
	}

	// Node grid positions for the 3D rack face
	const nodeGrid = Array.from({ length: 12 }, (_, i) => ({
		row: Math.floor(i / 4),
		col: i % 4,
		active: i < 9,
		primary: i < 3,
	}));

	// Reactive packet positions — computed here because {@const} cannot be
	// a direct child of an SVG root in Svelte (must be inside a block directive).
	$: p1 = getPacketPos(packet1Phase, 36, 118, 0, 100, 18, 118);
	$: p2 = getPacketPos(packet2Phase, 36, 125, 0, 108, 18, 125);
	$: p3 = getPacketPos(packet3Phase, 36, 112, 0, 94, 18, 112);
</script>

<div class="relative w-full flex items-center justify-center py-3 pointer-events-none select-none">
	<div class="node-float" style="width: 300px; height: 180px;">
		<svg
			viewBox="0 0 300 180"
			xmlns="http://www.w3.org/2000/svg"
			style="width:100%;height:100%;display:block"
			aria-hidden="true"
		>
			<defs>
				<!-- Isometric top face gradient -->
				<linearGradient id="sg-top" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#1D3545" />
					<stop offset="100%" stop-color="#0D1E2A" />
				</linearGradient>
				<!-- Isometric left face gradient -->
				<linearGradient id="sg-left" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" stop-color="#0A1520" />
					<stop offset="100%" stop-color="#081018" />
				</linearGradient>
				<!-- Isometric right face gradient -->
				<linearGradient id="sg-right" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" stop-color="#10222F" />
					<stop offset="100%" stop-color="#0A1820" />
				</linearGradient>
				<!-- Front display gradient -->
				<linearGradient id="sg-front" x1="0%" y1="0%" x2="0%" y2="100%">
					<stop offset="0%" stop-color="#0E2233" />
					<stop offset="100%" stop-color="#07111B" />
				</linearGradient>
				<!-- Cyan glow filter -->
				<filter id="glow-cyan" x="-50%" y="-50%" width="200%" height="200%">
					<feGaussianBlur in="SourceGraphic" stdDeviation="2.5" result="blur" />
					<feMerge>
						<feMergeNode in="blur" />
						<feMergeNode in="SourceGraphic" />
					</feMerge>
				</filter>
				<filter id="glow-orange" x="-50%" y="-50%" width="200%" height="200%">
					<feGaussianBlur in="SourceGraphic" stdDeviation="2" result="blur" />
					<feMerge>
						<feMergeNode in="blur" />
						<feMergeNode in="SourceGraphic" />
					</feMerge>
				</filter>
				<!-- Horizontal ambient glow at bottom -->
				<radialGradient id="ambient-glow" cx="50%" cy="100%" r="60%">
					<stop offset="0%" stop-color="#1A4A62" stop-opacity="0.35" />
					<stop offset="100%" stop-color="#000000" stop-opacity="0" />
				</radialGradient>
			</defs>

			<!-- Ambient glow under the rack -->
			<ellipse cx="148" cy="172" rx="110" ry="12" fill="url(#ambient-glow)" />

			<!-- ── Isometric Server Rack Body ── -->
			<!-- Top face -->
			<polygon
				points="148,20 260,62 260,80 148,38"
				fill="url(#sg-top)"
				stroke="#1E3A50"
				stroke-width="0.7"
			/>
			<!-- Left face -->
			<polygon
				points="36,62 148,20 148,38 36,80"
				fill="url(#sg-left)"
				stroke="#0D2030"
				stroke-width="0.7"
			/>
			<!-- Front face (main display panel) -->
			<polygon
				points="36,80 260,80 260,155 36,155"
				fill="url(#sg-front)"
				stroke="#1A3245"
				stroke-width="0.7"
			/>
			<!-- Front face inner bezel -->
			<polygon
				points="44,86 252,86 252,149 44,149"
				fill="none"
				stroke="#1E3850"
				stroke-width="0.5"
			/>

			<!-- ── Top Face Vent Lines ── -->
			{#each [0, 1, 2, 3, 4] as i}
				<line
					x1={88 + i * 22}
					y1={25 + i * 3}
					x2={172 + i * 18}
					y2={55 + i * 3}
					stroke="#1A3040"
					stroke-width="0.7"
					stroke-dasharray="3 4"
				/>
			{/each}

			<!-- ── Left Side Ridge Lines ── -->
			{#each [0, 1, 2] as i}
				<line
					x1={36}
					y1={95 + i * 20}
					x2={148}
					y2={95 + i * 20 - 14}
					stroke="#0F2535"
					stroke-width="0.5"
				/>
			{/each}

			<!-- ── Status Bar (top strip of front face) ── -->
			<rect x="44" y="86" width="208" height="14" rx="1" fill="#0A1C2C" />

			<!-- Power LED -->
			<circle
				cx="55"
				cy="93"
				r="3.5"
				fill={ledOn ? '#22C55E' : '#0A2215'}
				filter="url(#glow-cyan)"
			/>

			<!-- Node label text -->
			<text
				x="66"
				y="96.5"
				fill="#2A4A5E"
				font-size="5.5"
				font-family="JetBrains Mono, ui-monospace, monospace"
				font-weight="500"
			>NODE-01 · ngumpul.host</text>

			<!-- Network activity LED (right side) -->
			<circle
				cx="242"
				cy="93"
				r="2.5"
				fill="#FF5500"
				filter="url(#glow-orange)"
				class="inner-glow"
			/>
			<text
				x="228"
				y="96.5"
				fill="#3A5568"
				font-size="5"
				font-family="JetBrains Mono, ui-monospace, monospace"
				text-anchor="end"
			>NET ●</text>

			<!-- ── 12-Node Grid (3×4) ── -->
			{#each nodeGrid as node}
				{@const nx = 52 + node.col * 28}
				{@const ny = 105 + node.row * 14}
				<!-- Node slot background -->
				<rect
					x={nx}
					y={ny}
					width="22"
					height="10"
					rx="2"
					fill={node.active ? '#0C2235' : '#080F18'}
					stroke={node.primary ? '#1E5C7A' : (node.active ? '#133045' : '#0D1820')}
					stroke-width="0.7"
				/>
				<!-- Node LED indicator -->
				<circle
					cx={nx + 4}
					cy={ny + 5}
					r="2.2"
					fill={node.primary ? '#22C55E' : (node.active ? '#1A8C4E' : '#0C1C28')}
					filter={node.primary ? 'url(#glow-cyan)' : undefined}
				/>
				<!-- Node label -->
				<text
					x={nx + 9}
					y={ny + 7.5}
					fill={node.active ? '#2A5068' : '#0F1E2A'}
					font-size="4.5"
					font-family="JetBrains Mono, ui-monospace, monospace"
				>C{String(node.row * 4 + node.col + 1).padStart(2, '0')}</text>
			{/each}

			<!-- ── CPU / MEM Bars (right column) ── -->
			<!-- CPU bar -->
			<text x="172" y="109" fill="#2A4A5E" font-size="5" font-family="JetBrains Mono, ui-monospace, monospace">CPU</text>
			<rect x="186" y="104" width="60" height="6" rx="2" fill="#080F18" />
			<rect x="186" y="104" width={Math.round(cpuPct * 0.6)} height="6" rx="2" fill="#1A5C7A" style="transition: width 0.8s ease" />
			<text x="250" y="109" fill="#2A5A70" font-size="4.5" font-family="JetBrains Mono, ui-monospace, monospace">{Math.round(cpuPct)}%</text>

			<!-- MEM bar -->
			<text x="172" y="121" fill="#2A4A5E" font-size="5" font-family="JetBrains Mono, ui-monospace, monospace">MEM</text>
			<rect x="186" y="116" width="60" height="6" rx="2" fill="#080F18" />
			<rect x="186" y="116" width={Math.round(memPct * 0.6)} height="6" rx="2" fill="#1A4A60" style="transition: width 0.8s ease" />
			<text x="250" y="121" fill="#2A5A70" font-size="4.5" font-family="JetBrains Mono, ui-monospace, monospace">{Math.round(memPct)}%</text>

			<!-- Containers count -->
			<text x="172" y="133" fill="#2A4A5E" font-size="5" font-family="JetBrains Mono, ui-monospace, monospace">PODS</text>
			<text x="186" y="133" fill="#22C55E" font-size="5.5" font-family="JetBrains Mono, ui-monospace, monospace" font-weight="600">{activeContainers}</text>

			<!-- ── Bottom Port Row ── -->
			<line x1="44" y1="143" x2="252" y2="143" stroke="#0D2030" stroke-width="0.6" />
			{#each Array(8) as _, i}
				{@const px = 52 + i * 24}
				<rect x={px} y="144" width="16" height="9" rx="1.5" fill={i < 6 ? '#0A1E2E' : '#060D15'} stroke={i < 6 ? '#133045' : '#0D1820'} stroke-width="0.6" />
				<circle cx={px + 8} cy={148.5} r="1.8" fill={i < 6 ? '#22C55E' : '#0D1C28'} />
			{/each}
			<text x="248" y="151" fill="#1A3A50" font-size="4.5" font-family="JetBrains Mono, ui-monospace, monospace" text-anchor="end">10GbE</text>

			<!-- ── Animated Data Stream Packets (positions computed in script) ── -->
			<circle cx={p1.x} cy={p1.y} r="2.5" fill="#4A8FA8" opacity={0.7} filter="url(#glow-cyan)" />
			<circle cx={p2.x} cy={p2.y} r="2" fill="#22C55E" opacity={0.6} filter="url(#glow-cyan)" />
			<circle cx={p3.x} cy={p3.y} r="1.8" fill="#FF5500" opacity={0.65} filter="url(#glow-orange)" />

			<!-- Packet trail lines -->
			<path
				d="M 36 118 Q 18 118 0 100"
				fill="none"
				stroke="#4A8FA8"
				stroke-width="1"
				stroke-dasharray="3 9"
				opacity="0.25"
			/>
			<path
				d="M 36 125 Q 18 125 0 108"
				fill="none"
				stroke="#22C55E"
				stroke-width="0.8"
				stroke-dasharray="3 9"
				opacity="0.18"
			/>
			<path
				d="M 36 112 Q 18 112 0 94"
				fill="none"
				stroke="#FF5500"
				stroke-width="0.8"
				stroke-dasharray="3 9"
				opacity="0.2"
			/>

			<!-- ── Edge Connector / Port Badge ── -->
			<rect x="0" y="90" width="28" height="20" rx="4" fill="#0A1C2C" stroke="#1A3A50" stroke-width="0.8" />
			<text x="14" y="99" fill="#2A5A70" font-size="5" font-family="JetBrains Mono, ui-monospace, monospace" text-anchor="middle">ETH</text>
			<text x="14" y="106" fill="#1E4A62" font-size="4.5" font-family="JetBrains Mono, ui-monospace, monospace" text-anchor="middle">:1111</text>
		</svg>
	</div>
</div>
