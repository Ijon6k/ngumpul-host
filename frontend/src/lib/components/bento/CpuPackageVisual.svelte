<script lang="ts">
	export let model: string = '13th Gen Intel® Core™ i5-1334U';
	export let cores: number = 10;
	export let threads: number = 12;

	// Extract clean short name, e.g. "i5-1334U"
	$: shortModel = model.includes('i5-') ? 'i5-1334U' : model.replace(/^(13th Gen )?(Intel\(R\)\s+)?(Core\(TM\)\s+)?/, '').slice(0, 16);
</script>

<div class="relative w-full flex items-center justify-center select-none py-1">
	<svg
		viewBox="0 0 170 140"
		fill="none"
		xmlns="http://www.w3.org/2000/svg"
		class="w-full max-w-[175px] h-auto drop-shadow-[0_4px_12px_rgba(0,0,0,0.06)]"
		aria-label="Intel Processor Package"
	>
		<defs>
			<!-- Metallic Heat Spreader Gradient -->
			<linearGradient id="ihs-metal" x1="28" y1="18" x2="142" y2="122" gradientUnits="userSpaceOnUse">
				<stop offset="0%" stop-color="#F1F5F9" />
				<stop offset="35%" stop-color="#E2E8F0" />
				<stop offset="70%" stop-color="#CBD5E1" />
				<stop offset="100%" stop-color="#94A3B8" />
			</linearGradient>

			<!-- Substrate Fiber Green/Charcoal Gradient -->
			<linearGradient id="substrate-grad" x1="10" y1="10" x2="160" y2="130" gradientUnits="userSpaceOnUse">
				<stop offset="0%" stop-color="#1E293B" />
				<stop offset="100%" stop-color="#0F172A" />
			</linearGradient>

			<!-- Inner Cavity Bevel -->
			<linearGradient id="bevel-grad" x1="32" y1="22" x2="138" y2="118" gradientUnits="userSpaceOnUse">
				<stop offset="0%" stop-color="#FFFFFF" stop-opacity="0.6" />
				<stop offset="100%" stop-color="#64748B" stop-opacity="0.3" />
			</linearGradient>
		</defs>

		<!-- 1. PCB Substrate (Outer Carrier with Socket Keying Notches) -->
		<!-- Outer dimensions: 154 x 124, corners rounded r=5, with alignment notches on left/right -->
		<path
			d="M 14 10 
			   L 156 10 
			   A 4 4 0 0 1 160 14 
			   L 160 62 
			   A 3 3 0 0 0 157 65 
			   A 3 3 0 0 0 160 68 
			   L 160 126 
			   A 4 4 0 0 1 156 130 
			   L 14 130 
			   A 4 4 0 0 1 10 126 
			   L 10 68 
			   A 3 3 0 0 0 13 65 
			   A 3 3 0 0 0 10 62 
			   L 10 14 
			   A 4 4 0 0 1 14 10 Z"
			fill="url(#substrate-grad)"
			stroke="#334155"
			stroke-width="1.2"
		/>

		<!-- Pin 1 Alignment Gold Triangle (Top-Left) -->
		<polygon points="14,14 24,14 14,24" fill="#F59E0B" />

		<!-- Subtle Substrate Trace Routing Lines -->
		<g stroke="#F59E0B" stroke-opacity="0.25" stroke-width="0.75">
			<line x1="16" y1="36" x2="28" y2="36" />
			<line x1="16" y1="42" x2="26" y2="42" />
			<line x1="16" y1="96" x2="26" y2="96" />
			<line x1="16" y1="102" x2="28" y2="102" />
			<line x1="142" y1="36" x2="154" y2="36" />
			<line x1="144" y1="42" x2="154" y2="42" />
			<line x1="144" y1="96" x2="154" y2="96" />
			<line x1="142" y1="102" x2="154" y2="102" />
		</g>

		<!-- Miniature SMD Decoupling Capacitors around perimeter -->
		<g fill="#D97706" stroke="#92400E" stroke-width="0.4">
			<rect x="20" y="52" width="5" height="2.5" rx="0.5" />
			<rect x="20" y="78" width="5" height="2.5" rx="0.5" />
			<rect x="145" y="52" width="5" height="2.5" rx="0.5" />
			<rect x="145" y="78" width="5" height="2.5" rx="0.5" />
		</g>

		<!-- 2. Integrated Heat Spreader (IHS) Base Step -->
		<rect
			x="29"
			y="19"
			width="112"
			height="102"
			rx="6"
			fill="#475569"
			stroke="#1E293B"
			stroke-width="0.8"
		/>

		<!-- 3. Integrated Heat Spreader Main Metal Plate -->
		<rect
			x="32"
			y="22"
			width="106"
			height="96"
			rx="5"
			fill="url(#ihs-metal)"
			stroke="url(#bevel-grad)"
			stroke-width="1.2"
		/>

		<!-- Inner Rim Chamfer Accent -->
		<rect
			x="35"
			y="25"
			width="100"
			height="90"
			rx="3.5"
			fill="none"
			stroke="#FFFFFF"
			stroke-opacity="0.4"
			stroke-width="0.8"
		/>

		<!-- 4. Laser Engraved Typography (Minimalist Industrial Precision) -->
		<!-- Intel Logo & Brand -->
		<text
			x="42"
			y="44"
			fill="#334155"
			font-family="system-ui, -apple-system, sans-serif"
			font-size="8.5"
			font-weight="600"
			letter-spacing="-0.02em"
		>intel<tspan font-size="6" dy="-2">®</tspan></text>

		<!-- Processor Model -->
		<text
			x="42"
			y="61"
			fill="#0F172A"
			font-family="system-ui, -apple-system, sans-serif"
			font-size="11"
			font-weight="600"
			letter-spacing="-0.03em"
		>CORE™ <tspan fill="#0284C7" font-weight="600">{shortModel}</tspan></text>

		<!-- Spec / Stepping Code -->
		<text
			x="42"
			y="77"
			fill="#64748B"
			font-family="ui-monospace, monospace"
			font-size="6.5"
			letter-spacing="0.08em"
		>SRMD0 · 1.30GHZ</text>

		<!-- Cores & Threads Batch String -->
		<text
			x="42"
			y="89"
			fill="#64748B"
			font-family="ui-monospace, monospace"
			font-size="6"
			letter-spacing="0.06em"
		>{cores}C / {threads}T · HYBRID DIE</text>

		<!-- 2D Matrix / QR Code Serializer Placeholder (Bottom-Right) -->
		<g transform="translate(109, 87)">
			<rect x="0" y="0" width="18" height="18" fill="#1E293B" rx="1.5" />
			<rect x="2" y="2" width="4" height="4" fill="#E2E8F0" />
			<rect x="8" y="2" width="2" height="2" fill="#E2E8F0" />
			<rect x="12" y="2" width="4" height="4" fill="#E2E8F0" />
			<rect x="2" y="8" width="2" height="2" fill="#E2E8F0" />
			<rect x="6" y="6" width="4" height="4" fill="#0284C7" />
			<rect x="12" y="8" width="2" height="2" fill="#E2E8F0" />
			<rect x="2" y="12" width="4" height="4" fill="#E2E8F0" />
			<rect x="8" y="14" width="2" height="2" fill="#E2E8F0" />
			<rect x="12" y="12" width="4" height="4" fill="#E2E8F0" />
		</g>

		<!-- Tiny subtle country of origin text -->
		<text
			x="42"
			y="103"
			fill="#94A3B8"
			font-family="ui-monospace, monospace"
			font-size="5"
			letter-spacing="0.1em"
		>MADE IN MALAYSIA</text>
	</svg>
</div>
