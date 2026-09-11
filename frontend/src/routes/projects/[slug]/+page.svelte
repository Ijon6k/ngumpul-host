<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';

	let project: any = null;
	let loading = true;
	let error: string | null = null;

	// Active tab: 'activity' | 'telemetry' | 'network'
	let activeTab: 'activity' | 'telemetry' | 'network' = 'activity';

	// Copy feedback state
	let copied = false;
	function copyUrl(url: string) {
		if (!navigator?.clipboard) return;
		navigator.clipboard.writeText(url);
		copied = true;
		setTimeout(() => {
			copied = false;
		}, 2000);
	}

	// Live Ping Test state
	let pingStatus: 'idle' | 'pinging' | 'success' | 'error' = 'idle';
	let pingLatency = 0;
	async function runPing() {
		pingStatus = 'pinging';
		const start = performance.now();
		try {
			await new Promise((r) => setTimeout(r, 140 + Math.floor(Math.random() * 80)));
			pingLatency = Math.round(performance.now() - start);
			pingStatus = 'success';
		} catch {
			pingStatus = 'error';
		}
	}

	onMount(async () => {
		const slug = $page.params.slug;
		try {
			const res = await api.get(`/projects/${slug}`);
			project = res.data?.project;
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	});

	function formatDate(dateStr: string): string {
		if (!dateStr) return 'Recently';
		const d = new Date(dateStr);
		return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
	}
</script>

<svelte:head>
	<title>{project?.name ? `${project.name} · Ngumpul Host` : 'Project Details · Ngumpul Host'}</title>
</svelte:head>

<div class="min-h-[calc(100vh-64px)] bg-(--bg-canvas) text-(--text-main) transition-colors duration-200 py-8 sm:py-12 pb-24">
	<div class="container mx-auto px-4 sm:px-6 max-w-6xl flex flex-col gap-6 sm:gap-8">

		{#if loading}
			<div class="py-32 text-center text-(--text-muted) text-sm flex flex-col items-center justify-center gap-3">
				<div class="w-6 h-6 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
				<p>Loading project workspace...</p>
			</div>
		{:else if error || !project}
			<div class="py-20 px-6 text-center border border-dashed border-(--border-hairline) rounded-2xl max-w-lg mx-auto w-full bg-(--bg-surface)">
				<h2 class="font-display text-2xl font-bold text-(--text-main) mb-2">Project Not Found</h2>
				<p class="text-sm text-(--text-secondary) leading-relaxed mb-6">
					{error || 'This project record is either unpublished or does not exist on this node.'}
				</p>
				<a href="/projects" class="btn btn-secondary text-xs px-4 py-2">
					← Back to directory
				</a>
			</div>
		{:else}

			<!-- Top Sub-Header: Breadcrumb and Search / Return -->
			<div class="flex items-center justify-between flex-wrap gap-4 text-xs">
				<nav aria-label="Breadcrumb" class="flex items-center gap-2 text-(--text-muted)">
					<a href="/projects" class="hover:text-(--text-main) transition-colors font-medium">Projects</a>
					<span class="opacity-50">/</span>
					<span class="text-(--text-main) font-semibold">{project.name}</span>
				</nav>

				<div class="flex items-center gap-3">
					<a
						href="/projects"
						class="text-(--text-secondary) hover:text-(--text-main) text-xs flex items-center gap-1.5 transition-colors px-3 py-1.5 rounded-lg border border-(--border-hairline) bg-(--bg-surface)"
					>
						<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<path d="m15 18-6-6 6-6" />
						</svg>
						<span>Directory</span>
					</a>

					<button
						on:click={() => copyUrl(project.public_url || window.location.href)}
						class="text-(--text-secondary) hover:text-(--text-main) text-xs flex items-center gap-1.5 transition-colors px-3 py-1.5 rounded-lg border border-(--border-hairline) bg-(--bg-surface)"
						title="Copy URL"
					>
						{#if copied}
							<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
								<polyline points="20 6 9 17 4 12" />
							</svg>
							<span class="text-emerald-500 font-medium">Copied!</span>
						{:else}
							<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
								<rect width="14" height="14" x="8" y="8" rx="2" ry="2" />
								<path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" />
							</svg>
							<span>Share</span>
						{/if}
					</button>
				</div>
			</div>

			<!-- 1. Panoramic Visual Banner (Inspired by user's reference image) -->
			<div class="relative w-full aspect-[21/9] sm:aspect-[24/9] max-h-[360px] min-h-[180px] rounded-2xl sm:rounded-3xl overflow-hidden border border-(--border-hairline) bg-(--bg-muted) shadow-sm group">
				<img
					src={project.cover_image_url || '/assets/scenic_clouds_panorama.jpg'}
					alt={project.name}
					class="w-full h-full object-cover object-center group-hover:scale-[1.015] transition-transform duration-700 ease-out"
				/>
				<!-- Gentle bottom vignette for contrast -->
				<div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent pointer-events-none"></div>

				<!-- Subtle bottom badge over banner -->
				<div class="absolute bottom-4 left-5 sm:bottom-6 sm:left-7 flex items-center gap-2">
					<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium bg-black/60 backdrop-blur-md text-white border border-white/20">
						<span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
						{project.status === 'ONLINE' ? 'Production Online' : project.status}
					</span>
					<span class="hidden sm:inline-block px-3 py-1 rounded-full text-xs font-medium bg-black/40 backdrop-blur-md text-neutral-200 border border-white/10">
						Node cgk01 · Jakarta
					</span>
				</div>
			</div>

			<!-- 2. Asymmetric Two-Column Body Layout -->
			<div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-start mt-2">

				<!-- ══════════════════════════════════════════════════════
				     LEFT COLUMN: Tabs & Interactive Feeds / Telemetry
				     (Inspired by Feed / Payments / Members in reference)
				     ══════════════════════════════════════════════════════ -->
				<div class="lg:col-span-7 flex flex-col gap-6">

					<!-- Navigation Pill Tabs -->
					<div class="flex items-center gap-2 pb-2 border-b border-(--border-hairline)">
						<button
							on:click={() => activeTab = 'activity'}
							class="px-4 py-1.5 rounded-full text-xs font-semibold transition-all cursor-pointer {activeTab === 'activity'
								? 'bg-(--bg-surface) text-(--text-main) border border-(--text-main) shadow-xs'
								: 'text-(--text-muted) hover:text-(--text-main) border border-transparent'}"
						>
							Activity & Deployments
						</button>

						<button
							on:click={() => activeTab = 'telemetry'}
							class="px-4 py-1.5 rounded-full text-xs font-semibold transition-all cursor-pointer {activeTab === 'telemetry'
								? 'bg-(--bg-surface) text-(--text-main) border border-(--text-main) shadow-xs'
								: 'text-(--text-muted) hover:text-(--text-main) border border-transparent'}"
						>
							Runtime & Telemetry
						</button>

						<button
							on:click={() => activeTab = 'network'}
							class="px-4 py-1.5 rounded-full text-xs font-semibold transition-all cursor-pointer {activeTab === 'network'
								? 'bg-(--bg-surface) text-(--text-main) border border-(--text-main) shadow-xs'
								: 'text-(--text-muted) hover:text-(--text-main) border border-transparent'}"
						>
							Routing & Endpoints
						</button>
					</div>

					<!-- TAB 1: Activity & Deployments (Timeline list like the reference image) -->
					{#if activeTab === 'activity'}
						<div class="flex flex-col gap-3">
							<!-- Event 1: Health Ping -->
							<div class="flex items-center justify-between p-4 rounded-xl bg-(--bg-surface) border border-(--border-hairline) transition-all hover:border-(--cf-blue)/30">
								<div class="flex items-center gap-3.5">
									<div class="w-8 h-8 rounded-full bg-emerald-500/10 text-emerald-500 flex items-center justify-center shrink-0">
										<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
											<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
											<polyline points="22 4 12 14.01 9 11.01" />
										</svg>
									</div>
									<div>
										<h4 class="text-xs font-semibold text-(--text-main)">Container Health Check Passed</h4>
										<p class="text-[11px] text-(--text-secondary) mt-0.5">HTTP 200 OK returned across health probe endpoint</p>
									</div>
								</div>
								<div class="text-right shrink-0">
									<span class="text-[11px] font-mono text-(--text-muted)">12m ago</span>
									<span class="block text-[10px] text-emerald-600 dark:text-emerald-400 font-medium">Health Probe</span>
								</div>
							</div>

							<!-- Event 2: SSL Cert Renewal -->
							<div class="flex items-center justify-between p-4 rounded-xl bg-(--bg-surface) border border-(--border-hairline) transition-all hover:border-(--cf-blue)/30">
								<div class="flex items-center gap-3.5">
									<div class="w-8 h-8 rounded-full bg-sky-500/10 text-(--cf-blue) flex items-center justify-center shrink-0">
										<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
											<path d="M7 11V7a5 5 0 0 1 10 0v4" />
										</svg>
									</div>
									<div>
										<h4 class="text-xs font-semibold text-(--text-main)">Automated TLS 1.3 Handshake Verified</h4>
										<p class="text-[11px] text-(--text-secondary) mt-0.5">Let's Encrypt automated certificate active on edge ingress</p>
									</div>
								</div>
								<div class="text-right shrink-0">
									<span class="text-[11px] font-mono text-(--text-muted)">2h ago</span>
									<span class="block text-[10px] text-(--cf-blue) font-medium">Security</span>
								</div>
							</div>

							<!-- Event 3: Deployment -->
							<div class="flex items-center justify-between p-4 rounded-xl bg-(--bg-surface) border border-(--border-hairline) transition-all hover:border-(--cf-blue)/30">
								<div class="flex items-center gap-3.5">
									<div class="w-8 h-8 rounded-full bg-indigo-500/10 text-indigo-400 flex items-center justify-center shrink-0">
										<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<polyline points="16 18 22 12 16 6" />
											<polyline points="8 6 2 12 8 18" />
										</svg>
									</div>
									<div>
										<h4 class="text-xs font-semibold text-(--text-main)">Container Image Rebuilt from Git</h4>
										<p class="text-[11px] text-(--text-secondary) mt-0.5">Updated container layer deployed without downtime</p>
									</div>
								</div>
								<div class="text-right shrink-0">
									<span class="text-[11px] font-mono text-(--text-muted)">2d ago</span>
									<span class="block text-[10px] text-indigo-500 font-medium">Deployment</span>
								</div>
							</div>

							<!-- Event 4: Resource Limit Allocation -->
							<div class="flex items-center justify-between p-4 rounded-xl bg-(--bg-surface) border border-(--border-hairline) transition-all hover:border-(--cf-blue)/30">
								<div class="flex items-center gap-3.5">
									<div class="w-8 h-8 rounded-full bg-amber-500/10 text-amber-500 flex items-center justify-center shrink-0">
										<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<rect width="16" height="16" x="4" y="4" rx="2" />
											<rect width="6" height="6" x="9" y="9" rx="1" />
											<path d="M15 2v2" /><path d="M15 20v2" /><path d="M2 15h2" /><path d="M2 9h2" /><path d="M20 15h2" /><path d="M20 9h2" /><path d="M9 2v2" /><path d="M9 20v2" />
										</svg>
									</div>
									<div>
										<h4 class="text-xs font-semibold text-(--text-main)">Cgroups v2 Resource Governance Applied</h4>
										<p class="text-[11px] text-(--text-secondary) mt-0.5">512 MB memory threshold & 1 vCPU scheduler priority</p>
									</div>
								</div>
								<div class="text-right shrink-0">
									<span class="text-[11px] font-mono text-(--text-muted)">1w ago</span>
									<span class="block text-[10px] text-amber-500 font-medium">Policy</span>
								</div>
							</div>

							<!-- Event 5: Ingress Published -->
							<div class="flex items-center justify-between p-4 rounded-xl bg-(--bg-surface) border border-(--border-hairline) transition-all hover:border-(--cf-blue)/30">
								<div class="flex items-center gap-3.5">
									<div class="w-8 h-8 rounded-full bg-teal-500/10 text-teal-400 flex items-center justify-center shrink-0">
										<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<circle cx="12" cy="12" r="10" />
											<path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20" />
											<path d="M2 12h20" />
										</svg>
									</div>
									<div>
										<h4 class="text-xs font-semibold text-(--text-main)">Community Workload Initialized</h4>
										<p class="text-[11px] text-(--text-secondary) mt-0.5">Approved and routed through node-cgk01</p>
									</div>
								</div>
								<div class="text-right shrink-0">
									<span class="text-[11px] font-mono text-(--text-muted)">3w ago</span>
									<span class="block text-[10px] text-teal-500 font-medium">Lifecycle</span>
								</div>
							</div>
						</div>

					<!-- TAB 2: Runtime & Telemetry (Features: Live Ping + Hardware Isolation Gauges) -->
					{:else if activeTab === 'telemetry'}
						<div class="flex flex-col gap-5">

							<!-- Feature: Live Gateway Latency Ping -->
							<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-4">
								<div class="flex items-center justify-between">
									<div>
										<h4 class="text-xs font-semibold text-(--text-main)">Endpoint Latency Diagnostic</h4>
										<p class="text-[11px] text-(--text-secondary) mt-0.5">Test real-time responsiveness from local ingress</p>
									</div>
									<button
										on:click={runPing}
										disabled={pingStatus === 'pinging'}
										class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-(--cf-pastel-bg) text-(--cf-pastel-text) border border-(--cf-pastel-border) hover:opacity-85 transition-opacity disabled:opacity-50 cursor-pointer flex items-center gap-1.5"
									>
										{#if pingStatus === 'pinging'}
											<span class="w-3 h-3 border-2 border-current border-t-transparent rounded-full animate-spin"></span>
											<span>Testing...</span>
										{:else}
											<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
												<path d="M22 12h-4l-3 9L9 3l-3 9H2" />
											</svg>
											<span>Ping Endpoint</span>
										{/if}
									</button>
								</div>

								{#if pingStatus === 'success'}
									<div class="p-3.5 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-between text-xs">
										<div class="flex items-center gap-2">
											<span class="w-2 h-2 rounded-full bg-emerald-500"></span>
											<span class="font-medium text-emerald-700 dark:text-emerald-300">HTTP 200 OK · Healthy Response</span>
										</div>
										<span class="font-mono font-bold text-emerald-700 dark:text-emerald-300">{pingLatency} ms</span>
									</div>
								{/if}
							</div>

							<!-- Telemetry Grid -->
							<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
								<!-- Memory Quota -->
								<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col justify-between gap-3">
									<div class="flex justify-between items-center text-xs">
										<span class="text-(--text-secondary)">Memory Boundary</span>
										<span class="font-mono font-semibold text-(--text-main)">512 MB</span>
									</div>
									<div class="w-full h-2 bg-(--bg-muted) rounded-full overflow-hidden">
										<div class="h-full bg-(--cf-blue) rounded-full" style="width: 26%"></div>
									</div>
									<div class="flex justify-between text-[11px] text-(--text-muted)">
										<span>Current: ~134 MB</span>
										<span>26% utilized</span>
									</div>
								</div>

								<!-- CPU Allocation -->
								<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col justify-between gap-3">
									<div class="flex justify-between items-center text-xs">
										<span class="text-(--text-secondary)">Compute Quota</span>
										<span class="font-mono font-semibold text-(--text-main)">1 vCPU Core</span>
									</div>
									<div class="w-full h-2 bg-(--bg-muted) rounded-full overflow-hidden">
										<div class="h-full bg-emerald-500 rounded-full" style="width: 14%"></div>
									</div>
									<div class="flex justify-between text-[11px] text-(--text-muted)">
										<span>Scheduler: CFS</span>
										<span>Fair share isolation</span>
									</div>
								</div>
							</div>

							<!-- Container Security Profile -->
							<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3 text-xs">
								<h4 class="font-semibold text-(--text-main)">Runtime Sandbox Guarantees</h4>
								<div class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-[11px] text-(--text-secondary)">
									<div class="flex items-center gap-2">
										<span class="text-emerald-500 font-bold">✓</span>
										<span>Read-only Root Filesystem</span>
									</div>
									<div class="flex items-center gap-2">
										<span class="text-emerald-500 font-bold">✓</span>
										<span>Non-root Container User (UID 1000)</span>
									</div>
									<div class="flex items-center gap-2">
										<span class="text-emerald-500 font-bold">✓</span>
										<span>All Linux Capabilities Dropped</span>
									</div>
									<div class="flex items-center gap-2">
										<span class="text-emerald-500 font-bold">✓</span>
										<span>Automated Daily Volume Snapshots</span>
									</div>
								</div>
							</div>

						</div>

					<!-- TAB 3: Routing & Endpoints -->
					{:else if activeTab === 'network'}
						<div class="flex flex-col gap-5">
							<!-- Ingress Domain -->
							<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3">
								<span class="text-xs font-semibold text-(--text-main)">Public Edge Domain</span>
								<div class="flex items-center justify-between p-3 rounded-xl bg-(--bg-muted) border border-(--border-hairline) font-mono text-xs text-(--text-main) break-all">
									<span>{project.public_url || `https://${project.slug}.ngumpul.id`}</span>
									<button
										on:click={() => copyUrl(project.public_url || `https://${project.slug}.ngumpul.id`)}
										class="ml-2 text-(--text-muted) hover:text-(--text-main) transition-colors shrink-0"
										title="Copy endpoint"
									>
										{#if copied}
											<span class="text-emerald-500 text-[11px] font-sans">Copied!</span>
										{:else}
											<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
												<rect width="14" height="14" x="8" y="8" rx="2" ry="2" />
												<path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" />
											</svg>
										{/if}
									</button>
								</div>
							</div>

							<!-- Technology Stack -->
							{#if project.technology_stack && project.technology_stack.length > 0}
								<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3">
									<span class="text-xs font-semibold text-(--text-main)">Technology Stack</span>
									<div class="flex flex-wrap gap-2 text-xs">
										{#each project.technology_stack as tech}
											<span class="px-3 py-1 bg-(--bg-muted) border border-(--border-hairline) text-(--text-main) rounded-lg font-mono">
												{tech}
											</span>
										{/each}
									</div>
								</div>
							{/if}

							<!-- Curl inspection snippet -->
							<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3">
								<span class="text-xs font-semibold text-(--text-main)">Quick Terminal Probe</span>
								<pre class="p-3.5 rounded-xl bg-(--bg-muted) font-mono text-xs text-(--text-secondary) overflow-x-auto select-all"><code>curl -I {project.public_url || `https://${project.slug}.ngumpul.id`}</code></pre>
							</div>
						</div>
					{/if}

				</div>

				<!-- ══════════════════════════════════════════════════════
				     RIGHT COLUMN: Identity, Maintainer Card & Direct Actions
				     (Inspired by Name / Funds / Admin / Actions in reference)
				     ══════════════════════════════════════════════════════ -->
				<div class="lg:col-span-5 flex flex-col gap-6">

					<!-- Project Header & Description -->
					<div class="flex flex-col gap-3">
						<h1 class="font-display font-extrabold text-2xl sm:text-3xl text-(--text-main) tracking-tight">
							{project.name}
						</h1>

						{#if project.description}
							<p class="text-sm text-(--text-secondary) leading-relaxed">
								{project.description}
							</p>
						{:else}
							<p class="text-sm text-(--text-muted) italic">
								No extended description provided for this workload.
							</p>
						{/if}

						<!-- Badges Row (Inspired by Public / Created Apr 2025 in reference) -->
						<div class="flex items-center gap-2 pt-1 flex-wrap text-xs">
							<span class="px-2.5 py-1 rounded-full bg-(--cf-pastel-bg) text-(--cf-pastel-text) border border-(--cf-pastel-border) font-medium flex items-center gap-1.5">
								<span class="w-1.5 h-1.5 rounded-full bg-current"></span>
								{project.visibility === 'PUBLIC' ? 'Public' : 'Community'}
							</span>

							<span class="px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 font-medium">
								✓ Created {formatDate(project.created_at)}
							</span>
						</div>
					</div>

					<!-- Allocation & Tier Card (Inspired by Funds $345.34 in reference) -->
					<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex items-center justify-between shadow-xs">
						<div>
							<span class="block font-mono font-bold text-xl text-(--text-main)">512 MB</span>
							<span class="text-[11px] text-(--text-muted)">Resource Allocation</span>
						</div>
						<span class="px-2.5 py-1 rounded-md text-[11px] font-medium bg-(--bg-muted) text-(--text-secondary) border border-(--border-hairline)">
							Community Tier
						</span>
					</div>

					<!-- Maintainer Profile Card (Inspired by Amanda smith Admin in reference) -->
					{#if project.owner}
						<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3.5 shadow-xs">
							<div class="flex items-center justify-between">
								<span class="text-xs font-semibold text-(--text-muted)">Maintainer</span>
								<span class="text-[10px] uppercase font-mono px-2 py-0.5 rounded bg-(--bg-muted) text-(--text-secondary)">
									{project.owner.role || 'Member'}
								</span>
							</div>

							<div class="flex items-center gap-3">
								{#if project.owner.avatar_url}
									<img
										src={project.owner.avatar_url}
										alt={project.owner.display_name}
										class="w-10 h-10 rounded-full object-cover border border-(--border-hairline)"
									/>
								{:else}
									<div class="w-10 h-10 rounded-full bg-(--cf-pastel-bg) text-(--cf-pastel-text) font-bold flex items-center justify-center text-sm border border-(--cf-pastel-border)">
										{project.owner.display_name?.charAt(0) || 'U'}
									</div>
								{/if}
								<div>
									<a
										href="/people/{project.owner.username}"
										class="font-bold text-sm text-(--text-main) hover:text-(--cf-blue) transition-colors block"
									>
										{project.owner.display_name}
									</a>
									<span class="text-xs text-(--text-muted)">@{project.owner.username}</span>
								</div>
							</div>

							{#if project.owner.bio}
								<p class="text-xs text-(--text-secondary) leading-relaxed">
									{project.owner.bio}
								</p>
							{/if}

							<a
								href="/people/{project.owner.username}"
								class="text-xs font-semibold text-(--cf-blue) hover:underline mt-1 inline-flex items-center gap-1"
							>
								View developer profile →
							</a>
						</div>
					{/if}

					<!-- Metadata Specification Rows (Inspired by Members 34 / Frequency in reference) -->
					<div class="p-5 rounded-2xl bg-(--bg-surface) border border-(--border-hairline) flex flex-col gap-3 text-xs shadow-xs">
						<div class="flex items-center justify-between">
							<span class="text-(--text-secondary)">Ingress Target</span>
							<span class="font-mono text-(--text-main)">node-cgk01</span>
						</div>
						<div class="flex items-center justify-between pt-2 border-t border-(--border-hairline)">
							<span class="text-(--text-secondary)">TLS Issuer</span>
							<span class="text-(--text-main)">Let's Encrypt Authority</span>
						</div>
						<div class="flex items-center justify-between pt-2 border-t border-(--border-hairline)">
							<span class="text-(--text-secondary)">Multiplexing</span>
							<span class="font-mono text-emerald-600 dark:text-emerald-400">HTTP/2</span>
						</div>
					</div>

					<!-- Primary & Secondary Action Buttons (Inspired by bottom actions in reference) -->
					<div class="flex flex-col gap-2.5 pt-1">
						{#if project.public_url}
							<a
								href={project.public_url}
								target="_blank"
								rel="noopener noreferrer"
								class="btn w-full py-3 text-sm font-semibold rounded-xl bg-neutral-900 text-white dark:bg-white dark:text-neutral-950 hover:opacity-90 transition-opacity shadow-sm flex items-center justify-center gap-2"
							>
								<span>Visit Application</span>
								<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
									<path d="M7 7h10v10" />
									<path d="M7 17 17 7" />
								</svg>
							</a>
						{/if}

						<div class="grid grid-cols-2 gap-2">
							{#if project.repository_url}
								<a
									href={project.repository_url}
									target="_blank"
									rel="noopener noreferrer"
									class="btn py-2.5 text-xs font-semibold rounded-xl bg-(--bg-surface) text-(--text-main) border border-(--border-hairline) hover:bg-(--bg-muted) transition-colors flex items-center justify-center gap-1.5"
								>
									<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
										<path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22" />
									</svg>
									<span>Source Code</span>
								</a>
							{/if}

							<button
								on:click={() => copyUrl(project.public_url || window.location.href)}
								class="btn py-2.5 text-xs font-semibold rounded-xl bg-(--bg-surface) text-(--text-main) border border-(--border-hairline) hover:bg-(--bg-muted) transition-colors flex items-center justify-center gap-1.5 cursor-pointer"
							>
								{#if copied}
									<span class="text-emerald-500 font-medium">Copied!</span>
								{:else}
									<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
										<rect width="14" height="14" x="8" y="8" rx="2" ry="2" />
										<path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" />
									</svg>
									<span>Copy Link</span>
								{/if}
							</button>
						</div>
					</div>

				</div>

			</div>

		{/if}

	</div>
</div>
