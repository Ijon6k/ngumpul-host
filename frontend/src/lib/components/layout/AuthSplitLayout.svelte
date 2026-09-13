<script lang="ts">
	export let title: string;
	export let description: string = '';
</script>

<div class="min-h-screen w-full flex flex-col lg:grid lg:grid-cols-12 bg-(--bg-canvas)">
	<!-- Left 2/3 Editorial Visual (Desktop) - Clean photo with zero text overlay -->
	<div class="relative hidden lg:block lg:col-span-8 overflow-hidden bg-neutral-950 border-r border-(--border-hairline) select-none">
		<img
			src="/assets/loginserver.webp"
			alt="Ngumpul Host server hardware"
			class="w-full h-full object-cover object-center"
		/>
	</div>

	<!-- Right 1/3 Form Surface -->
	<div class="w-full lg:col-span-4 flex flex-col justify-center min-h-screen px-6 sm:px-10 lg:px-10 py-12 bg-(--bg-canvas) overflow-y-auto">
		<div class="w-full max-w-md lg:max-w-sm mx-auto flex flex-col gap-6">
			<!-- Header with Navbar Logo Component (slightly larger) -->
			<header class="flex flex-col gap-4">
				<a href="/" class="flex items-center gap-2.5 hover:opacity-85 transition-opacity w-fit" title="Ngumpul Host">
					<img src="/logo.webp" alt="Ngumpul Host" class="h-8 w-8 object-contain" />
					<span class="font-display font-bold text-lg tracking-tight text-(--text-main)">
						ngumpul<span class="font-normal text-(--text-secondary)">host</span><span class="text-(--accent-orange)">.</span>
					</span>
				</a>

				<div class="flex flex-col gap-1">
					<h1 class="font-sans font-semibold text-xl sm:text-2xl text-(--text-main) tracking-tight">
						{title}
					</h1>
					{#if description}
						<p class="text-sm text-(--text-secondary) leading-relaxed">
							{description}
						</p>
					{/if}
				</div>
			</header>

			<!-- Form Content -->
			<div class="flex flex-col gap-4">
				<slot />
			</div>

			<!-- Footer Slot -->
			<slot name="footer" />
		</div>
	</div>
</div>
