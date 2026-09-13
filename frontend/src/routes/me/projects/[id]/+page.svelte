<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { projectsApi, extractError } from '$lib/api';
	import { hostingRequestsApi } from '$lib/api/hostingRequests';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import ProjectAvailability from '$lib/components/ProjectAvailability.svelte';
	import { Tabs, Timeline, TimelineItem, Breadcrumb, BreadcrumbItem, BreadcrumbDropdown, MarkdownView, UptimeHistory } from '$lib/components/ui';
	import {
		ArrowUpRight,
		UploadSimple,
		Check,
		CaretLeft,
		Trash,
		X
	} from 'phosphor-svelte';
	import { detectLinkInfo } from '$lib/utils/linkDetector';
	import { getPingColorClass } from '$lib/utils/format';
	import { resolveProjectStatus } from '$lib/utils/projectStatus';
	import { compressImageToWebP } from '$lib/utils/imageCompressor';
	import { domainSuffix } from '$lib/stores/node';

	let projectId = '';
	let project: any = null;
	$: projectStatus = resolveProjectStatus(project);
	let availability: any = null;
	let activities: any[] = [];
	let allProjects: Array<{ id: string; label: string }> = [];
	let visitsData: {
		total_page_views: number;
		total_visits: number;
		daily: Array<{ date: string; page_views: number; visits: number }>;
	} | null = null;

	let loading = true;
	let error: string | null = null;
	let activeTab = 'overview';

	// Reactive route parameter watcher to support breadcrumb project switching
	$: if ($page.params.id && $page.params.id !== projectId) {
		projectId = $page.params.id;
		loadProjectData();
	}

	// Settings Form State
	let formName = '';
	let formDesc = '';
	let formReadme = '';
	let formCover = '';
	let formRepo = '';
	let formDocs = '';
	let formTech = '';
	let saving = false;
	let saveSuccess: string | null = null;
	let saveError: string | null = null;
	let uploadingCover = false;
	let fileInput: HTMLInputElement;
	let readmeInput: HTMLInputElement;

	// Subdomain Change Modal State
	let showSubdomainModal = false;
	let newSubdomain = '';
	let subdomainReason = '';
	let checkingSubdomain = false;
	let subdomainAvailable: boolean | null = null;
	let subdomainMessage = '';
	let submittingSubdomainChange = false;
	let subdomainChangeSuccess = '';
	let subdomainChangeError = '';
	let subdomainTimer: any = null;


	const tabsList = [
		{ id: 'overview', label: 'Overview' },
		{ id: 'visits', label: 'Visits' },
		{ id: 'activity', label: 'Activity' },
		{ id: 'settings', label: 'Settings' }
	];

	onMount(async () => {
		// Load all user projects for the breadcrumb switcher
		try {
			const res = await projectsApi.getMyProjects();
			allProjects = (res.projects || []).map((p: any) => ({ id: p.id, label: p.name }));
		} catch {
			// non-critical — switcher just won't populate
		}
	});

	async function loadProjectData() {
		loading = true;
		error = null;
		try {
			const [projRes, visitsRes] = await Promise.all([
				projectsApi.getMyProject(projectId),
				projectsApi.getMyProjectVisits(projectId)
			]);

			project = projRes.project;
			availability = projRes.availability;
			activities = projRes.activities || [];
			visitsData = visitsRes;

			if (project) {
				formName = project.name || '';
				formDesc = project.description || '';
				formReadme = project.readme || '';
				formCover = project.cover_image_url || '';
				formRepo = project.repository_url || '';
				formDocs = project.documentation_url || '';
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
			// Client-side WebP compression (1600x1200 max, 0.82 quality)
			const result = await compressImageToWebP(file, {
				maxWidth: 1600,
				maxHeight: 1200,
				quality: 0.82
			});
			const url = await projectsApi.uploadCover(result.file);
			if (url) {
				formCover = url;
			}
		} catch (err) {
			saveError = extractError(err);
		} finally {
			uploadingCover = false;
			if (fileInput) fileInput.value = '';
		}
	}

	async function handleReadmeUpload(event: Event) {
		const target = event.target as HTMLInputElement;
		if (!target.files || target.files.length === 0) return;
		const file = target.files[0];

		if (file.size > 2 * 1024 * 1024) {
			saveError = 'README file exceeds maximum limit of 2 MB.';
			return;
		}

		const reader = new FileReader();
		reader.onload = (e) => {
			if (typeof e.target?.result === 'string') {
				formReadme = e.target.result;
			}
		};
		reader.readAsText(file);
		if (readmeInput) readmeInput.value = '';
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
			await projectsApi.updateMyProject(projectId, {
				name: formName,
				description: formDesc,
				readme: formReadme,
				cover_image_url: formCover,
				repository_url: formRepo,
				documentation_url: formDocs,
				technology_stack: techStack
			});

			saveSuccess = 'Project settings saved successfully.';
			if (project) {
				project.name = formName;
				project.description = formDesc;
				project.readme = formReadme;
				project.cover_image_url = formCover;
				project.repository_url = formRepo;
				project.documentation_url = formDocs;
				project.technology_stack = techStack;
			}
		} catch (err) {
			saveError = extractError(err);
		} finally {
			saving = false;
		}
	}

	function handleSubdomainInput() {
		subdomainAvailable = null;
		subdomainMessage = '';
		clearTimeout(subdomainTimer);

		const clean = newSubdomain.trim().toLowerCase();
		if (!clean) return;

		if (clean === project?.slug) {
			subdomainAvailable = false;
			subdomainMessage = 'Current subdomain';
			return;
		}

		checkingSubdomain = true;
		subdomainTimer = setTimeout(async () => {
			try {
				const res = await hostingRequestsApi.checkSubdomain(clean);
				subdomainAvailable = res.available;
				subdomainMessage = res.message;
			} catch (err) {
				subdomainAvailable = false;
				subdomainMessage = extractError(err);
			} finally {
				checkingSubdomain = false;
			}
		}, 400);
	}

	async function handleSubmitSubdomainChange() {
		if (!newSubdomain || !subdomainAvailable) return;
		submittingSubdomainChange = true;
		subdomainChangeError = '';
		subdomainChangeSuccess = '';

		try {
			const res = await hostingRequestsApi.requestSubdomainChange(projectId, {
				new_subdomain: newSubdomain.trim().toLowerCase(),
				reason: subdomainReason.trim()
			});
			subdomainChangeSuccess = res.message || 'Subdomain change request submitted to admin.';
			setTimeout(() => {
				showSubdomainModal = false;
				newSubdomain = '';
				subdomainReason = '';
				subdomainChangeSuccess = '';
			}, 2500);
		} catch (err) {
			subdomainChangeError = extractError(err);
		} finally {
			submittingSubdomainChange = false;
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
		<!-- Breadcrumb with project switcher dropdown -->
		<Breadcrumb class="mb-4">
			<BreadcrumbItem href="/me/projects" icon={CaretLeft} separator={false}>
				Projects
			</BreadcrumbItem>
			<BreadcrumbItem current>
				<BreadcrumbDropdown
					label={project.name}
					activeId={projectId}
					options={allProjects}
					listboxLabel="Switch project"
					dropdownHeader="Switch project"
					on:select={(e) => goto(`/me/projects/${e.detail.id}`)}
				/>
			</BreadcrumbItem>
		</Breadcrumb>

		<!-- Project Header (Restrained Editorial Header) -->
		<header class="flex flex-col sm:flex-row sm:items-start justify-between gap-4 pb-6">
			<div class="flex flex-col gap-1.5 max-w-2xl">
				<div class="flex items-center gap-3">
					<h1 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-tight">
						{project.name}
					</h1>
					<StatusDot status={projectStatus.dotStatus} />
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

		{#if projectStatus.key === 'suspended'}
			<div class="mb-6 p-4 rounded-md bg-amber-500/10 border border-amber-500/30 text-amber-800 dark:text-amber-200 text-xs flex items-center justify-between">
				<div>
					<strong>Administrative Notice:</strong> This project has been placed under administrative suspension by the host operator.
				</div>
			</div>
		{/if}

		<!-- Minimal Typographic Tabs -->
		<div class="mb-8">
			<Tabs tabs={tabsList} bind:active={activeTab} />
		</div>

		<!-- TAB 1: OVERVIEW (Editorial 2-Column or Stacked Whitespace, Not 6-10 Cards) -->
		{#if activeTab === 'overview'}
			<div class="flex flex-col divide-y divide-(--border-hairline) gap-8">
				<!-- Cover Artwork (16:9 Ratio) -->
				{#if project.cover_image_url}
					<div class="w-full aspect-[16/9] max-h-[420px] overflow-hidden rounded-md border border-(--border-hairline) bg-(--bg-muted)">
						<ProjectCover src={project.cover_image_url} alt={project.name} name={project.name} aspectRatio="16/9" />
					</div>
				{/if}

				<!-- Health & Availability -->
				<section class="flex flex-col gap-3 pt-0">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Project health</h2>
					<div class="flex items-center gap-3 text-xs text-(--text-muted) flex-wrap">
						<StatusDot status={projectStatus.dotStatus} showLabel={false} />
						<span class="font-medium {projectStatus.colorClass}">
							{projectStatus.label}
						</span>
						{#if project.availability_reason && (projectStatus.key === 'unreachable' || project.availability === 'UNREACHABLE')}
							<span class="px-2 py-0.5 rounded-sm bg-rose-500/10 text-rose-600 dark:text-rose-400 font-mono text-xs">
								{project.availability_reason}
							</span>
						{/if}
						{#if availability?.uptime_percent != null}
							<span>·</span>
							<span class="font-mono">{availability.uptime_percent.toFixed(1)}% uptime</span>
						{/if}
						{#if availability?.latest_response_time_ms}
							<span>·</span>
							<span class="font-mono font-medium {getPingColorClass(availability.latest_response_time_ms)}">
								{availability.latest_response_time_ms} ms
							</span>
							<span>latency</span>
						{/if}
						{#if availability?.last_checked_at}
							<span>·</span>
							<span>Last checked {formatRelativeTime(availability.last_checked_at)}</span>
						{/if}
					</div>

					<!-- Multi-range quiet availability bar history -->
					<div class="max-w-xl pt-2">
						<UptimeHistory {availability} status={project.status} variant="detailed" showRangeSelector={true} />
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
							{@const repoInfo = detectLinkInfo(project.repository_url, 'Source')}
							<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
								<dt class="text-(--text-muted)">{repoInfo.label}</dt>
								<dd class="text-(--text-main) font-mono truncate max-w-[200px]">
									<a href={project.repository_url} target="_blank" rel="noreferrer" class="hover:underline flex items-center gap-1.5">
										<svelte:component this={repoInfo.icon} size={13} />
										<span class="truncate">{project.repository_url}</span>
									</a>
								</dd>
							</div>
						{/if}
						{#if project.documentation_url}
							{@const docInfo = detectLinkInfo(project.documentation_url, 'Documentation')}
							<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
								<dt class="text-(--text-muted)">{docInfo.label}</dt>
								<dd class="text-(--text-main) font-mono truncate max-w-[200px]">
									<a href={project.documentation_url} target="_blank" rel="noreferrer" class="hover:underline flex items-center gap-1.5">
										<svelte:component this={docInfo.icon} size={13} />
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

				<!-- Project README / Markdown View (Framed with marker border & README.md header) -->
				{#if project.readme}
					<section class="flex flex-col gap-3 pt-8">
						<MarkdownView content={project.readme} title="README.md" bordered={true} />
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
						<div class="flex items-center gap-4 text-xs">
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
										<span class="text-xs font-mono text-(--text-muted)">{day.date.slice(8)}</span>

										<!-- Tooltip -->
										<div class="absolute bottom-full mb-1 hidden group-hover:flex flex-col gap-0.5 z-10 bg-neutral-900 text-white text-xs px-2 py-1 rounded-sm shadow-md pointer-events-none whitespace-nowrap">
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
						<div class="flex items-center justify-between">
							<label for="p-desc" class="text-xs font-medium text-(--text-secondary)">Short description</label>
							<span class="text-xs font-mono text-(--text-muted)">{formDesc.length}/280</span>
						</div>
						<textarea
							id="p-desc"
							bind:value={formDesc}
							maxlength="280"
							rows="3"
							class="text-xs p-3 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none resize-none leading-relaxed"
							placeholder="Brief description of the project (max 280 characters)..."
						></textarea>
					</div>

					<div class="flex flex-col gap-1.5">
						<div class="flex items-center justify-between">
							<label for="p-readme" class="text-xs font-medium text-(--text-secondary)">Documentation / README (Markdown)</label>
							<div class="flex items-center gap-2">
								<input
									type="file"
									accept=".md,.markdown,text/markdown,text/plain"
									bind:this={readmeInput}
									on:change={handleReadmeUpload}
									class="hidden"
								/>
								<button
									type="button"
									class="btn btn-secondary btn-sm text-xs px-2.5 py-1 inline-flex items-center gap-1 cursor-pointer"
									on:click={() => readmeInput.click()}
								>
									<UploadSimple size={13} />
									<span>Upload .md file</span>
								</button>
							</div>
						</div>
						<textarea
							id="p-readme"
							bind:value={formReadme}
							rows="8"
							placeholder="# Project Documentation..."
							class="text-xs p-3 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) font-mono focus:border-(--accent-sky) outline-none leading-relaxed"
						></textarea>
						<span class="text-xs text-(--text-muted)">Supports GitHub-style markdown syntax. Rendered on the project overview page.</span>
					</div>
				</section>

				<div class="h-px bg-(--border-hairline)/60"></div>

				<!-- Section: Links -->
				<section class="flex flex-col gap-4">
					<div>
						<h2 class="text-sm font-medium text-(--text-main) tracking-tight">External links</h2>
						<p class="text-xs text-(--text-secondary) mt-0.5">Two primary reference links for your project.</p>
					</div>

					{#if project.public_url}
						<div class="flex flex-col gap-2 pb-2">
							<span class="text-xs font-medium text-(--text-secondary)">Live Application Domain</span>
							<div class="flex items-center justify-between gap-4 p-3 bg-(--bg-muted)/40 rounded-sm border border-(--border-hairline)">
								<a href="/go/{project.slug}" target="_blank" class="text-xs font-mono text-(--accent-sky) hover:underline flex items-center gap-1.5">
									<span>{project.public_url}</span>
									<ArrowUpRight size={13} />
								</a>
								<button
									type="button"
									class="btn btn-secondary btn-sm text-xs px-3 py-1 cursor-pointer"
									on:click={() => {
										showSubdomainModal = true;
										newSubdomain = '';
										subdomainReason = '';
										subdomainAvailable = null;
										subdomainMessage = '';
										subdomainChangeSuccess = '';
										subdomainChangeError = '';
									}}
								>
									Change subdomain
								</button>
							</div>
						</div>
					{/if}

					<div class="flex flex-col gap-1.5">
						<div class="flex items-center justify-between">
							<label for="p-repo" class="text-xs font-medium text-(--text-secondary)">Source code repository</label>
							{#if formRepo}
								{@const info = detectLinkInfo(formRepo)}
								<span class="text-xs text-(--text-muted) flex items-center gap-1 font-mono">
									<svelte:component this={info.icon} size={13} />
									<span>{info.label}</span>
								</span>
							{/if}
						</div>
						<input
							id="p-repo"
							type="url"
							bind:value={formRepo}
							placeholder="https://github.com/... or https://gitlab.com/..."
							class="text-xs px-3 py-2 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none font-mono"
						/>
					</div>

					<div class="flex flex-col gap-1.5">
						<div class="flex items-center justify-between">
							<label for="p-docs" class="text-xs font-medium text-(--text-secondary)">Documentation</label>
							{#if formDocs}
								{@const info = detectLinkInfo(formDocs)}
								<span class="text-xs text-(--text-muted) flex items-center gap-1 font-mono">
									<svelte:component this={info.icon} size={13} />
									<span>{info.label}</span>
								</span>
							{/if}
						</div>
						<input
							id="p-docs"
							type="url"
							bind:value={formDocs}
							placeholder="https://docs.google.com/... or Notion / GitBook / Drive"
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

				<!-- Section: Cover Artwork (16:9 Ratio) -->
				<section class="flex flex-col gap-4">
					<h2 class="text-sm font-medium text-(--text-main) tracking-tight">Cover artwork (16:9)</h2>
					<div class="flex flex-col sm:flex-row sm:items-center gap-4">
						<div class="w-56 aspect-[16/9] rounded-xs overflow-hidden border border-(--border-hairline) bg-(--bg-muted) shrink-0">
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
							<span class="text-xs text-(--text-muted)">JPG, PNG, WebP up to 10MB (16:9 standard ratio).</span>
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

	<!-- Subdomain Change Modal -->
	{#if showSubdomainModal}
		<div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs">
			<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-6 max-w-lg w-full shadow-lg flex flex-col gap-5">
				<div class="flex items-start justify-between">
					<div>
						<h3 class="font-sans font-semibold text-base text-(--text-main)">Request Subdomain Change</h3>
						<p class="text-xs text-(--text-secondary) mt-1 leading-relaxed">
							Requests are reviewed by administrators before proxy and DNS routes are updated.
						</p>
					</div>
					<button
						type="button"
						class="text-(--text-muted) hover:text-(--text-main) cursor-pointer p-1"
						on:click={() => (showSubdomainModal = false)}
					>
						<X size={16} />
					</button>
				</div>

				{#if subdomainChangeSuccess}
					<div class="p-3 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 text-xs rounded-sm flex items-center gap-2">
						<Check size={14} weight="bold" />
						<span>{subdomainChangeSuccess}</span>
					</div>
				{:else}
					{#if subdomainChangeError}
						<div class="p-3 bg-rose-500/10 text-rose-700 dark:text-rose-400 text-xs rounded-sm">
							{subdomainChangeError}
						</div>
					{/if}

					<div class="flex flex-col gap-2">
						<label for="new-subdomain-input" class="text-xs font-medium text-(--text-secondary)">New subdomain</label>
						<!-- Cloudflare-style Split Input -->
						<div class="flex items-center border border-(--border-hairline) rounded-sm bg-(--bg-canvas) overflow-hidden focus-within:border-(--accent-sky)">
							<input
								id="new-subdomain-input"
								type="text"
								placeholder="new-subdomain"
								bind:value={newSubdomain}
								on:input={handleSubdomainInput}
								class="grow text-xs px-3 py-2 bg-transparent text-(--text-main) font-mono outline-none"
							/>
							<span class="text-xs text-(--text-muted) font-mono px-3 py-2 bg-(--bg-muted)/60 border-l border-(--border-hairline) select-none">
								{$domainSuffix}
							</span>
						</div>

						<!-- Availability Feedback -->
						<div class="text-xs flex items-center gap-1.5 min-h-[1.25rem]">
							{#if checkingSubdomain}
								<span class="text-(--text-muted)">Checking availability...</span>
							{:else if subdomainAvailable === true}
								<span class="text-emerald-600 dark:text-emerald-400 font-medium">✓ Available</span>
							{:else if subdomainAvailable === false}
								<span class="text-rose-600 dark:text-rose-400 font-medium">✕ {subdomainMessage}</span>
							{/if}
						</div>
					</div>

					<div class="flex flex-col gap-1.5">
						<label for="subdomain-reason" class="text-xs font-medium text-(--text-secondary)">Reason for change (optional)</label>
						<textarea
							id="subdomain-reason"
							bind:value={subdomainReason}
							rows="2"
							placeholder="Briefly state why you'd like to rename the project subdomain..."
							class="text-xs p-2.5 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) focus:border-(--accent-sky) outline-none"
						></textarea>
					</div>

					<div class="flex items-center justify-end gap-2.5 pt-2 border-t border-(--border-hairline)">
						<button
							type="button"
							class="btn btn-secondary btn-sm text-xs px-3 py-1.5 cursor-pointer"
							on:click={() => (showSubdomainModal = false)}
							disabled={submittingSubdomainChange}
						>
							Cancel
						</button>
						<button
							type="button"
							class="btn btn-primary btn-sm text-xs px-4 py-1.5 cursor-pointer font-medium"
							on:click={handleSubmitSubdomainChange}
							disabled={!newSubdomain || !subdomainAvailable || submittingSubdomainChange}
						>
							{#if submittingSubdomainChange}
								Submitting...
							{:else}
								Submit request
							{/if}
						</button>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

