<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import ProjectAvailability from '$lib/components/ProjectAvailability.svelte';
	import { Tabs, Timeline, TimelineItem } from '$lib/components/ui';
	import {
		ArrowUpRight,
		GithubLogo,
		BookOpen,
		UploadSimple,
		Check,
		CaretLeft,
		Trash
	} from 'phosphor-svelte';

	let projectId = $page.params.id;
	let project: any = null;
	let availability: any = null;
	let activities: any[] = [];
	let visitsData: {
		total_page_views: number;
		total_visits: number;
		daily: Array<{ date: string; page_views: number; visits: number }>;
	} | null = null;

	let loading = true;
	let error: string | null = null;
	let activeTab = 'overview';

	// Settings Form State
	let formName = '';
	let formDesc = '';
	let formCover = '';
	let formRepo = '';
	let formDocs = '';
	let formDemo = '';
	let formTech = '';
	let saving = false;
	let saveSuccess: string | null = null;
	let saveError: string | null = null;
	let uploadingCover = false;
	let fileInput: HTMLInputElement;

	const tabsList = [
		{ id: 'overview', label: 'Overview' },
		{ id: 'visits', label: 'Visits' },
		{ id: 'activity', label: 'Activity' },
		{ id: 'settings', label: 'Settings' }
	];

	onMount(async () => {
		await loadProjectData();
	});

	async function loadProjectData() {
		loading = true;
		error = null;
		try {
			const [projRes, visitsRes] = await Promise.all([
				api.get(`/me/projects/${projectId}`),
				api.get(`/me/projects/${projectId}/visits`)
			]);

			project = projRes.data?.project;
			availability = projRes.data?.availability;
			activities = projRes.data?.activities || [];
			visitsData = visitsRes.data;

			if (project) {
				formName = project.name || '';
				formDesc = project.description || '';
				formCover = project.cover_image_url || '';
				formRepo = project.repository_url || '';
				formDocs = project.documentation_url || '';
				formDemo = project.demo_url || '';
				formTech = (project.technology_stack || []).join(', ');
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	async function handleCoverUpload(event: Event) {
		const target = event.target as HTMLInputElement;
		if (!target.files || target.files.length === 0) return;
		const file = target.files[0];

		if (file.size > 10 * 1024 * 1024) {
			saveError = 'File size exceeds maximum limit of 10 MB.';
			return;
		}

		uploadingCover = true;
		saveError = null;
		try {
			const formData = new FormData();
			formData.append('file', file);
			formData.append('purpose', 'cover');

			const res = await api.post('/upload', formData, {
				headers: { 'Content-Type': 'multipart/form-data' }
			});
			if (res.data?.url) {
				formCover = res.data.url;
			}
		} catch (err) {
			saveError = extractError(err);
		} finally {
			uploadingCover = false;
			if (fileInput) fileInput.value = '';
		}
	}

	async function handleSaveSettings() {
		saving = true;
		saveError = null;
		saveSuccess = null;

		const techStack = formTech
			.split(',')
			.map((t) => t.trim())
			.filter(Boolean);

		try {
			await api.patch(`/me/projects/${projectId}`, {
				name: formName,
				description: formDesc,
				cover_image_url: formCover,
				repository_url: formRepo,
				documentation_url: formDocs,
				demo_url: formDemo,
				technology_stack: techStack
			});

			saveSuccess = 'Project settings saved successfully.';
			if (project) {
				project.name = formName;
				project.description = formDesc;
				project.cover_image_url = formCover;
				project.repository_url = formRepo;
				project.documentation_url = formDocs;
				project.demo_url = formDemo;
				project.technology_stack = techStack;
			}
		} catch (err) {
			saveError = extractError(err);
		} finally {
			saving = false;
		}
	}

	function formatDate(dateStr?: string): string {
		if (!dateStr) return '—';
		try {
			const d = new Date(dateStr);
			return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
		} catch {
			return dateStr;
		}
	}

	function formatRelativeTime(dateStr: string): string {
		try {
			const d = new Date(dateStr);
			const diffSec = Math.floor((Date.now() - d.getTime()) / 1000);
			if (diffSec < 60) return 'Just now';
			const diffMin = Math.floor(diffSec / 60);
			if (diffMin < 60) return `${diffMin}m ago`;
			const diffHours = Math.floor(diffMin / 60);
			if (diffHours < 24) return `${diffHours}h ago`;
			const diffDays = Math.floor(diffHours / 24);
			if (diffDays < 30) return `${diffDays}d ago`;
			return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
		} catch {
			return dateStr;
		}
	}

	function getActivityDescription(a: any): string {
		if (a.metadata?.title) return a.metadata.title;
		if (a.metadata?.message) return a.metadata.message;
		switch (a.type) {
			case 'PROJECT_PUBLISHED':
				return 'Project became publicly available';
			case 'PROJECT_UPDATED':
				return 'Updated project details and links';
			case 'HEALTH_CHECK_RECOVERED':
			case 'PROJECT_ONLINE':
				return 'Health check recovered normally';
			default:
				return 'Project record updated';
		}
	}

	$: maxDailyVal = visitsData?.daily
		? Math.max(1, ...visitsData.daily.map((d) => Math.max(d.page_views, d.visits)))
		: 1;
</script>

<svelte:head>
	<title>{project?.name ? `${project.name} · Workspace` : 'Project Details · Ngumpul Host'}</title>
</svelte:head>

<div class="w-full flex flex-col">
	{#if loading}
		<div class="py-24 text-center text-xs text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-5 h-5 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading project workspace...</span>
		</div>
	{:else if error || !project}
		<div class="py-16 px-6 text-center border border-dashed border-(--border-hairline) rounded-md flex flex-col items-center gap-3 max-w-md mx-auto w-full">
			<h2 class="font-medium text-base text-(--text-main)">Project Not Located</h2>
			<p class="text-xs text-(--text-muted)">{error || 'The requested project could not be found or you do not have permission to view it.'}</p>
			<a href="/me/projects" class="btn btn-secondary btn-sm text-xs px-4 py-1.5 mt-2">
				← Back to Projects
			</a>
		</div>
	{:else}
		<!-- Breadcrumb -->
		<nav aria-label="Breadcrumb" class="flex items-center gap-2 text-xs text-(--text-muted) mb-4">
			<a href="/me/projects" class="hover:text-(--text-main) transition-colors flex items-center gap-1">
				<CaretLeft size={12} weight="bold" />
				<span>Projects</span>
			</a>
			<span class="opacity-40">/</span>
			<span class="text-(--text-main) font-medium">{project.name}</span>
		</nav>

		<!-- Project Header (Restrained Editorial Header) -->
		<header class="flex flex-col sm:flex-row sm:items-start justify-between gap-4 pb-6">
			<div class="flex flex-col gap-1.5 max-w-2xl">
				<div class="flex items-center gap-3">
					<h1 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-tight">
						{project.name}
					</h1>
					<StatusDot status={project.status} />
				</div>

				<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed font-normal">
					{project.description || 'A small collaborative workspace for teams and friends.'}
				</p>
			</div>

			<!-- Action Controls: Visit Project & Edit Project -->
			<div class="flex items-center gap-2 shrink-0 self-start">
				{#if project.public_url}
					<a
						href="/go/{project.slug}"
						target="_blank"
						rel="noreferrer"
						class="btn btn-primary btn-sm text-xs px-3 py-1.5 inline-flex items-center gap-1.5"
					>
						<span>Visit project</span>
						<ArrowUpRight size={13} weight="bold" />
					</a>
				{/if}

				<button
					type="button"
					class="btn btn-secondary btn-sm text-xs px-3 py-1.5 cursor-pointer"
					on:click={() => (activeTab = 'settings')}
				>
					Edit project
				</button>
			</div>
		</header>

		<!-- Minimal Typographic Tabs -->
		<div class="mb-8">
			<Tabs tabs={tabsList} bind:active={activeTab} />
		</div>

		<!-- TAB 1: OVERVIEW (Editorial 2-Column or Stacked Whitespace, Not 6-10 Cards) -->
		{#if activeTab === 'overview'}
			<div class="flex flex-col divide-y divide-(--border-hairline) gap-8">
				<!-- Health & Availability -->
				<section class="flex flex-col gap-3 pt-0">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Project health</h2>
					<div class="flex items-center gap-4 text-xs text-(--text-muted) flex-wrap">
						<StatusDot status={project.status} />
						{#if availability?.uptime_percent != null}
							<span>·</span>
							<span class="font-mono">{availability.uptime_percent.toFixed(1)}% uptime</span>
						{/if}
						{#if availability?.latest_response_time_ms}
							<span>·</span>
							<span class="font-mono">{availability.latest_response_time_ms} ms response time</span>
						{/if}
						{#if availability?.last_checked_at}
							<span>·</span>
							<span>Last checked {formatRelativeTime(availability.last_checked_at)}</span>
						{/if}
					</div>

					<!-- 30-day quiet availability summary -->
					<div class="max-w-md pt-2">
						<ProjectAvailability {availability} status={project.status} />
					</div>
				</section>

				<!-- Project Information -->
				<section class="flex flex-col gap-4 pt-8">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Project information</h2>
					<dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-3 text-xs max-w-xl">
						<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
							<dt class="text-(--text-muted)">Published</dt>
							<dd class="text-(--text-main) font-medium">{formatDate(project.published_at || project.created_at)}</dd>
						</div>
						<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
							<dt class="text-(--text-muted)">Hosted since</dt>
							<dd class="text-(--text-main) font-medium">{formatDate(project.created_at)}</dd>
						</div>
						<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
							<dt class="text-(--text-muted)">Last updated</dt>
							<dd class="text-(--text-main) font-medium">{formatDate(project.updated_at || project.created_at)}</dd>
						</div>
						<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
							<dt class="text-(--text-muted)">Public URL</dt>
							<dd class="text-(--text-main) font-mono truncate max-w-[200px]">
								{#if project.public_url}
									<a href="/go/{project.slug}" target="_blank" class="text-(--accent-sky) hover:underline">{project.public_url}</a>
								{:else}
									<span class="text-(--text-muted)">—</span>
								{/if}
							</dd>
						</div>
						{#if project.repository_url}
							<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
								<dt class="text-(--text-muted)">Source</dt>
								<dd class="text-(--text-main) font-mono truncate max-w-[200px]">
									<a href={project.repository_url} target="_blank" rel="noreferrer" class="hover:underline flex items-center gap-1">
										<GithubLogo size={13} />
										<span class="truncate">{project.repository_url}</span>
									</a>
								</dd>
							</div>
						{/if}
						{#if project.documentation_url}
							<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
								<dt class="text-(--text-muted)">Documentation</dt>
								<dd class="text-(--text-main) font-mono truncate max-w-[200px]">
									<a href={project.documentation_url} target="_blank" rel="noreferrer" class="hover:underline flex items-center gap-1">
										<BookOpen size={13} />
										<span class="truncate">{project.documentation_url}</span>
									</a>
								</dd>
							</div>
						{/if}
					</dl>
				</section>

				<!-- Built With -->
				{#if project.technology_stack && project.technology_stack.length > 0}
					<section class="flex flex-col gap-3 pt-8">
						<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Built with</h2>
						<p class="text-xs text-(--text-secondary) font-mono leading-relaxed">
							{project.technology_stack.join(' · ')}
						</p>
					</section>
				{/if}

				<!-- Recent Activity Snippet -->
				<section class="flex flex-col gap-4 pt-8">
					<div class="flex items-baseline justify-between gap-4">
						<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Recent activity</h2>
						<button
							type="button"
							class="text-xs text-(--text-secondary) hover:text-(--text-main) transition-colors cursor-pointer"
							on:click={() => (activeTab = 'activity')}
						>
							View all history →
						</button>
					</div>

					{#if activities.length === 0}
						<p class="text-xs text-(--text-muted)">No recent activity recorded for this project.</p>
					{:else}
						<Timeline density="compact">
							{#each activities.slice(0, 3) as act (act.id)}
								<TimelineItem
									title={getActivityDescription(act)}
									timestamp={formatRelativeTime(act.created_at)}
									density="compact"
								/>
							{/each}
						</Timeline>
					{/if}
				</section>
			</div>

		<!-- TAB 2: VISITS (Simple, No Google Analytics Clutter) -->
		{:else if activeTab === 'visits'}
			<div class="flex flex-col gap-8">
				<div class="flex flex-col gap-1">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Visits</h2>
					<p class="text-xs text-(--text-secondary)">Last 30 days metrics from community showcase views and outbound redirects.</p>
				</div>

				<!-- Two Metrics (Pure Typography) -->
				<div class="flex items-baseline gap-12 sm:gap-16 py-2">
					<div class="flex flex-col">
						<span class="text-3xl sm:text-4xl font-mono font-medium text-(--text-main)">
							{visitsData?.total_page_views || 0}
						</span>
						<span class="text-xs text-(--text-muted) mt-1">Project page views</span>
					</div>

					<div class="flex flex-col">
						<span class="text-3xl sm:text-4xl font-mono font-medium text-(--text-main)">
							{visitsData?.total_visits || 0}
						</span>
						<span class="text-xs text-(--text-muted) mt-1">Visits from Ngumpul</span>
					</div>
				</div>

				<!-- One Simple Daily Visualization (Each bar represents ONE day) -->
				<div class="flex flex-col gap-3 pt-4">
					<div class="flex items-center justify-between text-xs text-(--text-muted)">
						<span>Daily activity</span>
						<div class="flex items-center gap-4 text-[11px]">
							<span class="flex items-center gap-1.5">
								<span class="w-2 h-2 rounded-xs bg-neutral-300 dark:bg-neutral-700"></span>
								<span>Page views</span>
							</span>
							<span class="flex items-center gap-1.5">
								<span class="w-2 h-2 rounded-xs bg-(--accent-sky)"></span>
								<span>Visits</span>
							</span>
						</div>
					</div>

					{#if !visitsData?.daily || visitsData.daily.length === 0}
						<div class="py-12 text-center text-xs text-(--text-muted)">
							No traffic recorded in the last 30 days.
						</div>
					{:else}
						<div class="overflow-x-auto pb-2">
							<div class="min-w-[500px] h-36 flex items-end gap-1.5 pt-4">
								{#each visitsData.daily as day}
									{@const viewHeight = Math.max(3, Math.round((day.page_views / maxDailyVal) * 110))}
									{@const visitHeight = Math.max(3, Math.round((day.visits / maxDailyVal) * 110))}
									<div class="flex-1 flex flex-col items-center gap-1 group relative">
										<div class="w-full flex items-end justify-center gap-0.5 h-28">
											<div
												class="w-full max-w-[8px] bg-neutral-300 dark:bg-neutral-700 rounded-t-xs"
												style="height: {day.page_views > 0 ? viewHeight : 2}px;"
											></div>
											<div
												class="w-full max-w-[8px] bg-(--accent-sky) rounded-t-xs"
												style="height: {day.visits > 0 ? visitHeight : 2}px;"
											></div>
										</div>
										<span class="text-[9px] font-mono text-(--text-muted)">{day.date.slice(8)}</span>

										<!-- Tooltip -->
										<div class="absolute bottom-full mb-1 hidden group-hover:flex flex-col gap-0.5 z-10 bg-neutral-900 text-white text-[10px] px-2 py-1 rounded-sm shadow-md pointer-events-none whitespace-nowrap">
											<span>{day.date}</span>
											<span class="text-neutral-300">Views: {day.page_views}</span>
											<span class="text-(--accent-sky)">Visits: {day.visits}</span>
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			</div>

		<!-- TAB 3: ACTIVITY (Vertical Dot Timeline, Not Table/Cards) -->
		{:else if activeTab === 'activity'}
			<div class="flex flex-col gap-6">
				<div class="flex flex-col gap-1 pb-2">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Project activity</h2>
					<p class="text-xs text-(--text-secondary)">Complete chronological history of changes, deployments, and checks.</p>
				</div>

				{#if activities.length === 0}
					<p class="text-xs text-(--text-muted) py-8">No recorded activity for this project yet.</p>
				{:else}
					<Timeline density="comfortable">
						{#each activities as act (act.id)}
							<TimelineItem
								title={getActivityDescription(act)}
								timestamp={formatRelativeTime(act.created_at)}
								description={act.metadata?.message || ''}
								density="comfortable"
							/>
						{/each}
					</Timeline>
				{/if}
			</div>

		<!-- TAB 4: SETTINGS (Editorial Sections, No Boxed Form Cards) -->
		{:else if activeTab === 'settings'}
			<form on:submit|preventDefault={handleSaveSettings} class="flex flex-col gap-8 max-w-xl">
				{#if saveSuccess}
					<div class="p-3 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 text-xs rounded-sm flex items-center gap-2">
						<Check size={14} weight="bold" />
						<span>{saveSuccess}</span>
					</div>
				{/if}

				{#if saveError}
					<div class="p-3 bg-rose-500/10 text-rose-700 dark:text-rose-400 text-xs rounded-sm">
						{saveError}
					</div>
				{/if}

				<!-- Section: Project Details -->
				<section class="flex flex-col gap-4">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Project details</h2>

					<div class="flex flex-col gap-1.5">
						<label for="p-name" class="text-xs font-medium text-(--text-secondary)">Project name</label>
						<input
							id="p-name"
							type="text"
							bind:value={formName}
							required
							class="text-xs px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none"
						/>
					</div>

					<div class="flex flex-col gap-1.5">
						<label for="p-desc" class="text-xs font-medium text-(--text-secondary)">Description</label>
						<textarea
							id="p-desc"
							bind:value={formDesc}
							rows="3"
							class="text-xs p-3 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none resize-none leading-relaxed"
						></textarea>
					</div>
				</section>

				<div class="h-px bg-(--border-hairline)/60"></div>

				<!-- Section: Links -->
				<section class="flex flex-col gap-4">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Links</h2>

					<div class="flex flex-col gap-1.5">
						<label for="p-demo" class="text-xs font-medium text-(--text-secondary)">Public URL</label>
						<input
							id="p-demo"
							type="url"
							bind:value={formDemo}
							placeholder="https://..."
							class="text-xs px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none font-mono"
						/>
					</div>

					<div class="flex flex-col gap-1.5">
						<label for="p-repo" class="text-xs font-medium text-(--text-secondary)">Source code repository</label>
						<input
							id="p-repo"
							type="url"
							bind:value={formRepo}
							placeholder="https://github.com/..."
							class="text-xs px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none font-mono"
						/>
					</div>

					<div class="flex flex-col gap-1.5">
						<label for="p-docs" class="text-xs font-medium text-(--text-secondary)">Documentation</label>
						<input
							id="p-docs"
							type="url"
							bind:value={formDocs}
							placeholder="https://docs..."
							class="text-xs px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none font-mono"
						/>
					</div>
				</section>

				<div class="h-px bg-(--border-hairline)/60"></div>

				<!-- Section: Technologies -->
				<section class="flex flex-col gap-4">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Technologies</h2>
					<div class="flex flex-col gap-1.5">
						<label for="p-tech" class="text-xs text-(--text-muted)">Comma-separated list (e.g. SvelteKit, Go, PostgreSQL, Docker)</label>
						<input
							id="p-tech"
							type="text"
							bind:value={formTech}
							class="text-xs px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none font-mono"
						/>
					</div>
				</section>

				<div class="h-px bg-(--border-hairline)/60"></div>

				<!-- Section: Cover Artwork -->
				<section class="flex flex-col gap-4">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Cover</h2>
					<div class="flex flex-col sm:flex-row sm:items-center gap-4">
						<div class="w-28 h-18 rounded-xs overflow-hidden border border-(--border-hairline) bg-(--bg-muted) shrink-0">
							<ProjectCover src={formCover} alt="Cover preview" name={project.name} aspectRatio="16/9" />
						</div>

						<div class="flex flex-col gap-2">
							<input
								type="file"
								accept="image/*"
								bind:this={fileInput}
								on:change={handleCoverUpload}
								class="hidden"
							/>
							<button
								type="button"
								class="btn btn-secondary btn-sm text-xs px-3 py-1.5 inline-flex items-center gap-1.5 self-start cursor-pointer"
								on:click={() => fileInput.click()}
								disabled={uploadingCover}
							>
								{#if uploadingCover}
									<span class="w-3 h-3 border-2 border-current border-t-transparent rounded-full animate-spin"></span>
								{:else}
									<UploadSimple size={14} />
								{/if}
								<span>Change cover</span>
							</button>
							<span class="text-[11px] text-(--text-muted)">JPG, PNG, WebP up to 10MB.</span>
						</div>
					</div>
				</section>

				<div class="pt-4">
					<button
						type="submit"
						class="btn btn-primary btn-sm text-xs px-5 py-2 inline-flex items-center gap-1.5 cursor-pointer"
						disabled={saving}
					>
						{#if saving}
							<span class="w-3 h-3 border-2 border-current border-t-transparent rounded-full animate-spin"></span>
						{/if}
						<span>Save changes</span>
					</button>
				</div>
			</form>
		{/if}
	{/if}
</div>
