<script lang="ts">
	// Pinterest Minimalist Industrial Server Rack Blade Enclosure
	// Aesthetic: Dieter Rams / Teenage Engineering / Modern Architectural Tech
</script>

<div class="relative w-full h-full min-h-[300px] flex items-center justify-center p-4 select-none pointer-events-none">
	<!-- Ambient Studio Backlight -->
	<div class="absolute inset-0 bg-gradient-to-t from-sky-500/10 via-sky-500/5 to-transparent blur-2xl rounded-full opacity-60"></div>

	<!-- Minimalist Server Box SVG -->
	<svg
		viewBox="0 0 540 320"
		class="w-full max-w-[480px] h-auto drop-shadow-[0_20px_35px_rgba(0,0,0,0.7)]"
		fill="none"
		xmlns="http://www.w3.org/2000/svg"
		aria-hidden="true"
	>
		<defs>
			<!-- Anodized Chassis Gradient -->
			<linearGradient id="chassis-body" x1="0%" y1="0%" x2="0%" y2="100%">
				<stop offset="0%" stop-color="#18202A" />
				<stop offset="12%" stop-color="#0F151D" />
				<stop offset="85%" stop-color="#0A0E14" />
				<stop offset="100%" stop-color="#06090D" />
			</linearGradient>

			<!-- Top Beveled Rim Highlight -->
			<linearGradient id="top-highlight" x1="0%" y1="0%" x2="100%" y2="0%">
				<stop offset="0%" stop-color="#334155" stop-opacity="0.2" />
				<stop offset="30%" stop-color="#94A3B8" stop-opacity="0.8" />
				<stop offset="70%" stop-color="#38BDF8" stop-opacity="0.6" />
				<stop offset="100%" stop-color="#334155" stop-opacity="0.2" />
			</linearGradient>

			<!-- Ambient Floor Shadow -->
			<radialGradient id="floor-glow" cx="50%" cy="50%" r="50%">
				<stop offset="0%" stop-color="#38BDF8" stop-opacity="0.25" />
				<stop offset="60%" stop-color="#0EA5E9" stop-opacity="0.08" />
				<stop offset="100%" stop-color="#000000" stop-opacity="0" />
			</radialGradient>

			<!-- Minimalist Soft LED Glow -->
			<filter id="led-soft-glow" x="-30%" y="-30%" width="160%" height="160%">
				<feGaussianBlur stdDeviation="3" result="blur" />
				<feMerge>
					<feMergeNode in="blur" />
					<feMergeNode in="SourceGraphic" />
				</feMerge>
			</filter>
		</defs>

		<!-- Diffuse Floor Glow -->
		<ellipse cx="270" cy="285" rx="220" ry="25" fill="url(#floor-glow)" />

		<!-- Main Rack Chassis Monolith -->
		<rect
			x="40"
			y="50"
			width="460"
			height="210"
			rx="14"
			fill="url(#chassis-body)"
			stroke="#1E293B"
			stroke-width="1.2"
		/>

		<!-- Precision Top Edge Light Reflection -->
		<line
			x1="54"
			y1="50.6"
			x2="486"
			y2="50.6"
			stroke="url(#top-highlight)"
			stroke-width="1"
		/>

		<!-- Recessed Inner Front Facet -->
		<rect
			x="55"
			y="65"
			width="430"
			height="180"
			rx="8"
			fill="#080C10"
			stroke="#131D28"
			stroke-width="0.8"
		/>

		<!-- Architectural Top LED Strip (Razor Thin Light Accent) -->
		<rect
			x="75"
			y="78"
			width="390"
			height="2"
			rx="1"
			fill="#38BDF8"
			filter="url(#led-soft-glow)"
			opacity="0.85"
		/>

		<!-- Blade Bay Partitions (3 Symmetrical Minimalist Drive Modules) -->
		{#each [0, 1, 2] as i}
			{@const bx = 75 + i * 134}
			<g transform="translate({bx}, 95)">
				<!-- Bay Outer Panel -->
				<rect
					x="0"
					y="0"
					width="122"
					height="132"
					rx="6"
					fill="#0B1016"
					stroke="#16222F"
					stroke-width="0.8"
				/>

				<!-- Top Eject Latch / Clean Handle Groove -->
				<rect
					x="12"
					y="10"
					width="98"
					height="6"
					rx="3"
					fill="#070A0E"
					stroke="#1E2D3D"
					stroke-width="0.6"
				/>

				<!-- Precision Acoustic Micro-Vent Slits (Minimalist Grid) -->
				{#each Array(6) as _, row}
					<line
						x1="16"
						y1={32 + row * 10}
						x2="106"
						y2={32 + row * 10}
						stroke="#15212E"
						stroke-width="1"
						stroke-linecap="round"
						stroke-dasharray="2 4"
					/>
				{/each}

				<!-- Subtle Module Identifier -->
				<text
					x="16"
					y="108"
					fill="#334155"
					font-family="JetBrains Mono, monospace"
					font-size="6.5"
					font-weight="600"
					letter-spacing="0.1em"
				>NVMe B0{i+1}</text>

				<!-- Micro Status Light (Minimalist Pin-point) -->
				<circle
					cx="103"
					cy="106"
					r="1.8"
					fill={i === 0 ? '#38BDF8' : (i === 1 ? '#0284C7' : '#0369A1')}
					opacity="0.9"
				/>

				<!-- Lower Thermal Relief Notch -->
				<rect
					x="45"
					y="120"
					width="32"
					height="3"
					rx="1.5"
					fill="#06090D"
				/>
			</g>
		{/each}

		<!-- Minimalist Chassis Badging (Bottom Right) -->
		<text
			x="485"
			y="237"
			text-anchor="end"
			fill="#475569"
			font-family="JetBrains Mono, monospace"
			font-size="6"
			letter-spacing="0.16em"
		>NODE // INDEPENDENT HOST</text>

		<!-- Subdued Power Touch Contact (Bottom Left) -->
		<circle
			cx="68"
			cy="235"
			r="3"
			fill="#080C10"
			stroke="#1E293B"
			stroke-width="0.8"
		/>
		<circle
			cx="68"
			cy="235"
			r="1.2"
			fill="#38BDF8"
			opacity="0.7"
		/>
	</svg>
</div>
