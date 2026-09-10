<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import ProjectCard from '$lib/components/ProjectCard.svelte';
	import ActivityTimeline from '$lib/components/ActivityTimeline.svelte';
	import ServerNode3D from '$lib/components/ServerNode3D.svelte';
	import CosmicShowcase from '$lib/components/CosmicShowcase.svelte';
	import StatusDot from '$lib/components/StatusDot.svelte';

	let projects: any[] = [];
	let activities: any[] = [];
	let statusData: any = null;
	let loading = true;

	onMount(async () => {
		try {
			const [projRes, actRes, statRes] = await Promise.all([
				api.get('/projects'),
				api.get('/activity'),
				api.get('/status')
			]);
			projects = projRes.data?.projects || [];
			activities = actRes.data?.activities || [];
			statusData = statRes.data || null;
		} catch (e) {
			console.error('Failed to load portal data:', e);
		} finally {
			loading = false;
		}
	});

	$: counts = {
		projects: statusData?.counts?.projects || projects.length || 4,
		members: statusData?.counts?.members || 12
	};

	$: featuredProjects = projects.slice(0, 3);
</script>

<div class="flex flex-col">

	<!-- ══════════════════════════════════════════════════════════════════
	     1. HERO — Apple-clean split layout with 3D Server Node Visual
	     ══════════════════════════════════════════════════════════════════ -->
	<section class="pt-12 sm:pt-16 pb-10 border-b border-(--border-hairline)">
		<div class="container mx-auto px-6 max-w-6xl">
			<div class="grid grid-cols-1 lg:grid-cols-[1fr_420px] gap-6 lg:gap-10 items-stretch">

				<!-- Left — Editorial Content Card -->
				<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-2xl p-8 sm:p-10 flex flex-col justify-between gap-8 min-h-[340px]">
					<div class="flex flex-col gap-5">
						<!-- Eyebrow -->
						<div class="flex items-center gap-2.5 text-xs text-(--text-secondary)">
							<span class="w-2 h-2 rounded-full bg-(--color-success) animate-live-pulse inline-block shrink-0"></span>
							<span class="font-semibold text-(--text-main)">Independent Node</span>
							<span class="text-(--text-muted)" aria-hidden="true">·</span>
							<span class="font-mono text-[11px] text-(--text-muted)">Jakarta, ID</span>
						</div>

						<!-- Headline – Plus Jakarta Sans tight tracking -->
						<h1 class="font-display font-extrabold text-(--text-main) leading-[1.1]"
							style="font-size: clamp(2.1rem, 4.2vw, 3.25rem); letter-spacing: -0.038em;">
							A small, intentional corner<br class="hidden sm:block" /> of the internet.
						</h1>

						<!-- Sub -->
						<p class="text-[1.0625rem] text-(--text-secondary) leading-relaxed max-w-[44ch]">
							A community-hosted space for side projects, tools, experiments,
							and web creations worth keeping online.
							<span class="text-(--text-muted)">Maintained quietly on an independent server.</span>
						</p>
					</div>

					<!-- CTAs -->
					<div class="flex flex-col gap-6">
						<div class="flex items-center gap-3 flex-wrap">
							<a href="/projects" class="btn btn-primary inline-flex items-center gap-3 px-5 py-2.5">
								<span>Explore projects</span>
								<span class="w-5 h-5 rounded-full bg-(--accent-orange) text-white inline-flex items-center justify-center font-bold text-xs shrink-0">→</span>
							</a>
							<a href="/register" class="btn btn-secondary px-5 py-2.5">
								Request hosting
							</a>
						</div>

						<!-- Bottom info strip -->
						<div class="pt-5 border-t border-(--border-hairline) flex items-center justify-between flex-wrap gap-3 text-xs text-(--text-muted)">
							<span>Open source · Community operated</span>
							<div class="flex items-center gap-3">
								<!-- Real server telemetry chips -->
								<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full border border-(--border-hairline) bg-(--bg-muted) font-mono text-[10px] text-(--text-secondary)">
									<span class="w-1.5 h-1.5 rounded-full bg-(--color-success) inline-block"></span>
									Ingress Online
								</span>
								<span class="font-mono text-[10px] text-(--text-muted)">99.98% SLA</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Right Column – Server Telemetry Stack -->
				<div class="flex flex-col gap-4">

					<!-- Card A: 3D Server Node + Live Stats -->
					<div class="flex-1 bg-[#07101A] border border-[rgba(255,255,255,0.08)] rounded-2xl p-5 flex flex-col relative overflow-hidden">
						<!-- Top metrics strip -->
						<div class="grid grid-cols-3 gap-2 pb-4 border-b border-white/7 z-10">
							<div class="flex flex-col gap-0.5">
								<span class="text-[9px] font-mono font-medium text-neutral-600 tracking-widest uppercase">Pods</span>
								<div class="flex items-baseline gap-1.5">
									<span class="text-[1.6rem] font-bold text-white leading-none font-display">{counts.projects}</span>
									<span class="text-[10px] text-emerald-400 font-medium">active</span>
								</div>
							</div>
							<div class="flex flex-col gap-0.5">
								<span class="text-[9px] font-mono font-medium text-neutral-600 tracking-widest uppercase">Members</span>
								<div class="flex items-baseline gap-1.5">
									<span class="text-[1.6rem] font-bold text-white leading-none font-display">{counts.members}</span>
									<span class="text-[10px] text-sky-400 font-medium">online</span>
								</div>
							</div>
							<div class="flex flex-col gap-0.5 items-end">
								<span class="text-[9px] font-mono font-medium text-neutral-600 tracking-widest uppercase">Uptime</span>
								<span class="text-[1.1rem] font-bold text-white leading-none font-mono mt-0.5">99.98<span class="text-sm text-neutral-500">%</span></span>
							</div>
						</div>

						<!-- 3D Server Node Visualization -->
						<div class="flex-1 flex items-center justify-center py-1 z-10">
							<ServerNode3D
								cpuPct={42}
								memPct={61}
								activeContainers={18}
							/>
						</div>

						<!-- Bottom status row -->
						<div class="pt-3 border-t border-white/7 flex items-center justify-between text-[10px] text-neutral-500 z-10 gap-2 flex-wrap">
							<div class="flex items-center gap-1.5">
								<span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-live-pulse inline-block"></span>
								<span class="font-mono">Ingress Online · SSL Active</span>
							</div>
							<div class="flex items-center gap-3">
								<span class="font-mono text-[10px] text-neutral-600">14ms Jakarta Edge</span>
								<span class="font-mono text-neutral-400">3.2 GB/s NVMe</span>
							</div>
						</div>
					</div>

					<!-- Card B: Live Services -->
					<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-2xl p-5 flex flex-col gap-4">
						<div class="flex items-center justify-between pb-3 border-b border-(--border-hairline)">
							<span class="font-semibold text-sm text-(--text-main)">Ready Services</span>
							<span class="text-[10px] font-mono text-(--text-muted) bg-(--bg-muted) px-2 py-0.5 rounded-full border border-(--border-hairline)">HTTP/2</span>
						</div>

						<div class="flex flex-col divide-y divide-(--border-subtle)">
							{#if projects.length === 0}
								<div class="text-xs text-(--text-muted) py-2">Loading services...</div>
							{:else}
								{#each projects.slice(0, 3) as item}
									<div class="flex items-center justify-between py-2.5 text-xs">
										<div class="flex items-center gap-2 truncate max-w-[68%]">
											<span class="w-1.5 h-1.5 rounded-full bg-(--color-success) shrink-0"></span>
											<a href="/projects/{item.slug}" class="font-medium hover:text-(--accent-strong) transition-colors truncate text-(--text-main)">
												{item.name}
											</a>
										</div>
										<span class="text-[11px] text-(--color-success) font-medium shrink-0">Online</span>
									</div>
								{/each}
							{/if}
						</div>

						<div class="pt-2 border-t border-(--border-subtle) flex items-center justify-between text-xs">
							<span class="text-(--text-muted)">Need a domain or container?</span>
							<a href="/register" class="text-xs font-semibold text-(--accent-strong) hover:underline inline-flex items-center gap-1">
								Request host →
							</a>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- ══════════════════════════════════════════════════════════════════
	     2. COSMIC SHOWCASE — Dark planet banner with metric pillars
	     ══════════════════════════════════════════════════════════════════ -->
	<section class="py-6 sm:py-8">
		<div class="container mx-auto px-6 max-w-6xl">
			<CosmicShowcase />
		</div>
	</section>

	<!-- ══════════════════════════════════════════════════════════════════
	     3. ARCHITECTURE SECTION — Numbered 3-step with server spec cards
	     ══════════════════════════════════════════════════════════════════ -->
	<section class="py-14 sm:py-18 border-t border-(--border-hairline)">
		<div class="container mx-auto px-6 max-w-6xl flex flex-col gap-10">

			<!-- Section header -->
			<div class="max-w-xl">
				<h2 class="font-display font-bold text-(--text-main)" style="font-size: clamp(1.6rem, 3vw, 2.2rem); letter-spacing: -0.03em;">
					A small system, shared by a few people.
				</h2>
				<p class="text-base text-(--text-secondary) mt-3 leading-relaxed">
					Transparent by design. We document everything: the workload, the members,
					the infrastructure topology.
				</p>
			</div>

			<!-- 3-step process + spec cards -->
			<div class="grid grid-cols-1 lg:grid-cols-[1.1fr_1fr] gap-6">

				<!-- Left: How it works — numbered steps -->
				<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-2xl p-7 sm:p-8 flex flex-col gap-7">
					<div>
						<h3 class="font-display font-bold text-lg text-(--text-main)">Human-Scale Infrastructure</h3>
						<p class="text-sm text-(--text-secondary) leading-relaxed mt-2">
							We run on an independent Linux host. Projects are reviewed by a person,
							provisioned deliberately, not auto-generated.
						</p>
					</div>

					<div class="flex flex-col gap-5 pt-5 border-t border-(--border-hairline)">
						{#each [
							{ n: '01', title: 'Build your application', body: 'Any stack: Go, Svelte, Node, Rust, or Python. Docker-compatible or process-based.' },
							{ n: '02', title: 'Submit repository & notes', body: 'Include a Docker Compose config or a process specification. Notes help the review.' },
							{ n: '03', title: 'Manual provisioning', body: 'An administrator maps the vhost, reviews the codebase, and publishes the project.' },
						] as step}
							<div class="flex items-start gap-4">
								<span class="text-xs font-bold text-(--accent-strong) font-mono pt-0.5 shrink-0">{step.n}</span>
								<div>
									<strong class="block text-sm font-semibold text-(--text-main)">{step.title}</strong>
									<p class="text-xs text-(--text-secondary) mt-0.5 leading-relaxed">{step.body}</p>
								</div>
							</div>
						{/each}
					</div>
				</div>

				<!-- Right: Server Spec Cards Grid -->
				<div class="grid grid-cols-2 gap-4">
					{#each [
						{ icon: '⬡', value: '32 Cores', label: 'Dedicated CPU', sub: 'AMD EPYC · Multi-thread enabled', color: '#1A4A62' },
						{ icon: '▣', value: '64 GB', label: 'ECC RAM', sub: 'Error-correcting · Registered DIMM', color: '#1A3D55' },
						{ icon: '◈', value: '3.2 GB/s', label: 'NVMe I/O', sub: 'PCIe Gen 4 · RAID-1 mirror', color: '#1E4060' },
						{ icon: '◉', value: '10 GbE', label: 'Network', sub: 'Jakarta IXP direct · Sub-14ms latency', color: '#152E45' },
					] as spec}
						<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-5 flex flex-col gap-3 hover:shadow-sm transition-shadow">
							<div class="w-8 h-8 rounded-lg flex items-center justify-center text-base" style="background:{spec.color}20; color:{spec.color}">
								{spec.icon}
							</div>
							<div>
								<div class="font-display font-bold text-xl text-(--text-main) leading-none">{spec.value}</div>
								<div class="text-xs font-semibold text-(--text-secondary) mt-1">{spec.label}</div>
								<div class="text-xs text-(--text-muted) mt-1 leading-snug">{spec.sub}</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</section>

	<!-- ══════════════════════════════════════════════════════════════════
	     4. ACTIVITY TIMELINE SPLIT
	     ══════════════════════════════════════════════════════════════════ -->
	<section class="py-14 sm:py-18 border-t border-(--border-hairline)">
		<div class="container mx-auto px-6 max-w-6xl grid grid-cols-1 lg:grid-cols-[1.3fr_1fr] gap-12 lg:gap-14">
			<div>
				<div class="mb-7">
					<h2 class="font-display font-bold text-(--text-main)" style="font-size: clamp(1.5rem, 2.8vw, 2rem); letter-spacing: -0.03em;">
						Recent Activity
					</h2>
					<p class="text-base text-(--text-secondary) mt-2 leading-relaxed">
						A quiet chronicle of deployments, updates, and member events.
					</p>
				</div>

				<ActivityTimeline {activities} />

				<div class="mt-8">
					<a href="/activity" class="text-xs font-semibold text-(--accent-strong) hover:underline inline-flex items-center gap-1">
						View complete activity timeline →
					</a>
				</div>
			</div>

			<!-- Infrastructure topology nodes -->
			<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-2xl p-7 sm:p-8 flex flex-col gap-5">
				<div>
					<div class="flex items-center gap-2 mb-1">
						<span class="w-2 h-2 rounded-full bg-(--color-success)"></span>
						<h3 class="font-display font-bold text-lg text-(--text-main)">Infrastructure Topology</h3>
					</div>
					<p class="text-xs text-(--text-muted)">Port 1111 isolation · Single public entrypoint</p>
				</div>

				<div class="flex flex-col gap-2">
					{#each [
						{ n: '01', label: 'Edge Ingress', name: 'Nginx Gateway', port: ':1111', note: 'Proxy & SSL', ok: true },
						{ n: '02', label: 'Web Tier', name: 'SvelteKit SSR', port: ':3000', note: 'Hydration', ok: true },
						{ n: '03', label: 'Core Engine', name: 'Go Chi API', port: ':8080', note: 'REST /api/v1', ok: true },
						{ n: '04', label: 'Database', name: 'PostgreSQL 16', port: ':5432', note: 'ACID persistence', ok: true },
					] as node}
						<div class="flex items-center gap-3 p-3 rounded-xl bg-(--bg-muted) border border-(--border-hairline) hover:border-(--accent-sky)/30 transition-colors">
							<span class="text-[10px] font-mono font-semibold text-(--text-muted) w-6 shrink-0">{node.n}</span>
							<div class="flex-1 min-w-0">
								<div class="flex items-center gap-2">
									<span class="text-xs font-semibold text-(--text-main)">{node.name}</span>
									<code class="text-[10px] font-mono text-(--accent-strong) bg-(--accent-soft) px-1.5 py-0.5 rounded">{node.port}</code>
								</div>
								<span class="text-[10px] text-(--text-muted)">{node.label} · {node.note}</span>
							</div>
							<span class="w-1.5 h-1.5 rounded-full bg-(--color-success) shrink-0"></span>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</section>

	<!-- ══════════════════════════════════════════════════════════════════
	     5. HOSTED PROJECTS GRID
	     ══════════════════════════════════════════════════════════════════ -->
	<section class="py-14 sm:py-18 border-t border-(--border-hairline)">
		<div class="container mx-auto px-6 max-w-6xl flex flex-col gap-8">
			<div class="flex items-end justify-between flex-wrap gap-4">
				<div>
					<h2 class="font-display font-bold text-(--text-main)" style="font-size: clamp(1.5rem, 2.8vw, 2rem); letter-spacing: -0.03em;">
						Hosted Projects
					</h2>
					<p class="text-base text-(--text-secondary) mt-2 leading-relaxed">
						Tools, experiments, and web creations on this node.
					</p>
				</div>
				<a href="/projects" class="text-xs font-semibold text-(--accent-strong) hover:underline">
					View all ({projects.length}) →
				</a>
			</div>

			{#if projects.length === 0 && !loading}
				<div class="text-center py-16 px-4 text-(--text-secondary) border border-dashed border-(--border-hairline) rounded-2xl">
					<p class="text-sm">No projects registered yet.</p>
				</div>
			{:else}
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
					{#each projects.slice(0, 6) as project}
						<ProjectCard {project} />
					{/each}
				</div>
			{/if}
		</div>
	</section>

</div>
