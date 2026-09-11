<script lang="ts">
	import BentoCard from './BentoCard.svelte';

	const steps = [
		{
			n: '1',
			title: 'Share your repository',
			desc: 'Drop your project link with a simple Dockerfile or compose setup.'
		},
		{
			n: '2',
			title: 'Quick review & setup',
			desc: 'We check the container, allocate fair resources, and set up your subdomain & SSL.'
		},
		{
			n: '3',
			title: 'Live on the web',
			desc: 'Your experiment or bot stays online 24/7 for friends and community to visit.'
		}
	];
</script>

<BentoCard title="How to Host With Us" badge="Simple Workflow" badgeClass="text-(--text-muted)">
	<div class="flex flex-col gap-4">
		<p class="text-xs text-(--text-secondary) leading-relaxed">
			We keep hosting simple and personal. No complex cloud consoles or surprise billing—just a friendly homelab space for things you've built.
		</p>

		<div class="flex flex-col gap-3">
			{#each steps as step}
				<div class="flex items-start gap-3">
					<span
						class="w-6 h-6 rounded-full bg-(--cf-pastel-bg) text-(--cf-pastel-text) border border-(--cf-pastel-border) flex items-center justify-center font-bold text-xs shrink-0 mt-0.5"
					>
						{step.n}
					</span>
					<div>
						<h4 class="font-semibold text-xs text-(--text-main)">{step.title}</h4>
						<p class="text-[11px] text-(--text-secondary) leading-snug mt-0.5">{step.desc}</p>
					</div>
				</div>
			{/each}
		</div>
	</div>

	<svelte:fragment slot="footer">
		<div class="flex items-center justify-between text-xs">
			<span class="text-(--text-muted)">Have an open-source project?</span>
			<a href="/register" class="font-semibold text-(--cf-blue) hover:underline flex items-center gap-1">
				<span>Request a spot</span>
				<span>→</span>
			</a>
		</div>
	</svelte:fragment>
</BentoCard>
