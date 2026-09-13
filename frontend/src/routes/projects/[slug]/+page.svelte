<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, commentsApi, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import Comments from '$lib/components/Comments.svelte';
	import ReportModal from '$lib/components/ReportModal.svelte';
	import { Tabs, Timeline, TimelineItem, Breadcrumb, MarkdownView, UptimeHistory } from '$lib/components/ui';
	import {
		ArrowUpRight,
		Flag,
		ShareNetwork,
		Check
	} from 'phosphor-svelte';
	import { detectLinkInfo } from '$lib/utils/linkDetector';
	import { getPingColorClass } from '$lib/utils/format';
	import { resolveProjectStatus } from '$lib/utils/projectStatus';

	import type { PageData } from './$types';

	export let data: PageData;

	let copied = false;
	let reportModalOpen = false;
	let activeTab = 'overview';
	let commentCount = 0;

	$: project = data.project;
	$: projectStatus = resolveProjectStatus(project);
	$: availability = data.availability;
	$: activities = data.activities || [];
	let loading = false;
	let error: string | null = null;

	$: tabsList = [
		{ id: 'overview', label: 'Overview' },
		{ id: 'comments', label: commentCount > 0 ? `Comments (${commentCount})` : 'Comments' },
		{ id: 'activity', label: 'Activity' }
	];

	onMount(async () => {
		if (project?.slug) {
			try {
				const res = await commentsApi.getProjectComments(project.slug);
				commentCount = res?.comments?.length || 0;
			} catch {
				// non-critical
			}
		}
	});

	function copyUrl() {
		if (!navigator?.clipboard) return;
		navigator.clipboard.writeText(window.location.href);
		copied = true;
		setTimeout(() => {
			copied = false;
		}, 2000);
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
				return 'Project details and links updated';
			case 'HEALTH_CHECK_RECOVERED':
			case 'PROJECT_ONLINE':
				return 'Project came online and operational';
			case 'PROJECT_OFFLINE':
				return 'Project became temporarily unavailable';
			default:
				return 'Project event recorded';
		}
	}
</script>

<svelte:head>
	<title>{project?.name ? `${project.name} · Ngumpul Host` : 'Project Showcase · Ngumpul Host'}</title>
	<meta name="description" content={project?.description || 'Independent software project hosted on Ngumpul Host.'} />
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-screen-xl py-8 sm:py-12 pb-24 flex flex-col">
	{#if loading}
		<div class="py-32 text-center text-(--text-muted) text-sm flex flex-col items-center justify-center gap-3">
			<div class="w-5 h-5 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
			<p class="text-xs">Loading showcase...</p>
		</div>
	{:else if error || !project}
		<div class="py-20 px-6 text-center border border-dashed border-(--border-hairline) rounded-md max-w-lg mx-auto w-full bg-(--bg-surface) flex flex-col items-center gap-4">
			<h2 class="font-normal text-xl text-(--text-main)">Project Not Found</h2>
			<p class="text-xs text-(--text-secondary) leading-relaxed">
				{error || 'The requested project might not be published yet, or the link has expired.'}
			</p>
			<a href="/projects" class="btn btn-secondary text-xs px-4 py-2">
				← Back to Catalog
			</a>
		</div>
	{:else}
		<!-- Breadcrumb -->
		<Breadcrumb
			class="mb-6"
			items={[
				{ label: 'Projects', href: '/projects' },
				{ label: project.name }
			]}
		/>

		<!-- Public Project Header (Editorial Showcase) -->
		<header class="flex flex-col gap-3">
			<!-- Project Title -->
			<h1 class="font-sans font-normal text-3xl sm:text-4xl lg:text-5xl text-(--text-main) tracking-[-0.02em] leading-tight">
				{project.name}
			</h1>

			<!-- Short Description -->
			{#if project.description}
				<p class="text-sm sm:text-base text-(--text-secondary) max-w-2xl leading-relaxed font-normal">
					{project.description}
				</p>
			{/if}

			<!-- Header Actions (Restrained) -->
			<div class="flex items-center gap-2.5 pt-3 flex-wrap">
				{#if project.public_url}
					<a
						href="/go/{project.slug}"
						target="_blank"
						rel="noopener noreferrer"
						class="btn btn-primary btn-sm text-xs px-4 py-2 inline-flex items-center gap-1.5 font-medium"
					>
						<span>Visit project</span>
						<ArrowUpRight size={13} weight="bold" />
					</a>
				{/if}

				{#if project.repository_url}
					{@const sourceInfo = detectLinkInfo(project.repository_url, 'Source')}
					<a
						href={project.repository_url}
						target="_blank"
						rel="noopener noreferrer"
						class="btn btn-secondary btn-sm text-xs px-3 py-2 inline-flex items-center gap-1.5"
					>
						<svelte:component this={sourceInfo.icon} size={14} />
						<span>{sourceInfo.label === 'GitHub' || sourceInfo.label === 'GitLab' ? sourceInfo.label : 'Source'}</span>
						<ArrowUpRight size={11} />
					</a>
				{/if}

				{#if project.documentation_url}
					{@const docInfo = detectLinkInfo(project.documentation_url, 'Documentation')}
					<a
						href={project.documentation_url}
						target="_blank"
						rel="noopener noreferrer"
						class="btn btn-secondary btn-sm text-xs px-3 py-2 inline-flex items-center gap-1.5"
					>
						<svelte:component this={docInfo.icon} size={14} />
						<span>Docs</span>
						<ArrowUpRight size={11} />
					</a>
				{/if}

				<button
					type="button"
					on:click={copyUrl}
					class="btn btn-secondary btn-sm text-xs px-3 py-2 inline-flex items-center gap-1.5 cursor-pointer text-(--text-secondary) hover:text-(--text-main)"
					title="Share link"
				>
					{#if copied}
						<Check size={13} weight="bold" class="text-emerald-500" />
						<span class="text-emerald-500 font-medium">Copied</span>
					{:else}
						<ShareNetwork size={14} />
						<span>Share</span>
					{/if}
				</button>

				<button
					type="button"
					on:click={() => (reportModalOpen = true)}
					class="btn btn-secondary btn-sm text-xs px-3 py-2 inline-flex items-center gap-1.5 cursor-pointer text-(--text-muted) hover:text-(--text-main)"
					title="Report this project"
				>
					<Flag size={13} />
					<span>Report</span>
				</button>
			</div>
		</header>

		<!-- Desktop Side-by-Side Composition (Cover & Tabs on Left, Sticky Context Column on Right) -->
		<div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-start mt-8">
			<!-- Left Column: 16:9 Cover + Main Tabs (Overview | Comments | Activity) -->
			<div class="lg:col-span-8 order-2 lg:order-1 flex flex-col gap-8">
				<!-- Project Cover Artwork (16:9 Ratio) -->
				<div class="w-full aspect-[16/9] max-h-[460px] rounded-md overflow-hidden border border-(--border-hairline) bg-(--bg-muted) shadow-xs flex">
					<ProjectCover
						src={project.cover_image_url}
						alt={project.name}
						name={project.name}
						aspectRatio="16/9"
						class="w-full h-full !rounded-none !border-none"
					/>
				</div>

				<!-- Main Tabs Container -->
				<div class="flex flex-col gap-6">
					<Tabs tabs={tabsList} bind:active={activeTab} />

					<!-- TAB 1: OVERVIEW -->
					{#if activeTab === 'overview'}
						<div class="flex flex-col gap-8">
							<!-- 1. About -->
							{#if project.description}
								<section class="flex flex-col gap-2.5">
									<h2 class="font-sans font-medium text-sm text-(--text-main)">
										About
									</h2>
									<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">
										{project.description}
									</p>
								</section>
							{/if}

							<!-- 2. README / Documentation -->
							{#if project.readme}
								<section class="flex flex-col gap-2.5">
									<MarkdownView content={project.readme} title="README.md" bordered={true} />
								</section>
							{/if}

							<!-- 3. Built With -->
							{#if project.technology_stack && project.technology_stack.length > 0}
								<section class="flex flex-col gap-2.5">
									<h2 class="font-sans font-medium text-sm text-(--text-main)">
										Built with
									</h2>
									<p class="text-xs text-(--text-secondary) font-mono leading-relaxed">
										{project.technology_stack.join(' · ')}
									</p>
								</section>
							{/if}

							<!-- 4. Relevant Project Information -->
							<section class="flex flex-col gap-3 pt-4 border-t border-(--border-hairline)">
								<h2 class="font-sans font-medium text-sm text-(--text-main)">
									Project information
								</h2>
								<dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-2.5 text-xs">
									<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
										<dt class="text-(--text-muted)">Published</dt>
										<dd class="text-(--text-main) font-medium">{formatDate(project.published_at || project.created_at)}</dd>
									</div>
									<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
										<dt class="text-(--text-muted)">Hosted since</dt>
										<dd class="text-(--text-main) font-medium">{formatDate(project.created_at)}</dd>
									</div>
									<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
										<dt class="text-(--text-muted)">Public URL</dt>
										<dd class="text-(--text-main) font-mono truncate max-w-[180px]">
											{#if project.public_url}
												<a href="/go/{project.slug}" target="_blank" class="text-(--accent-sky) hover:underline">{project.public_url}</a>
											{:else}
												<span class="text-(--text-muted)">—</span>
											{/if}
										</dd>
									</div>
									{#if project.repository_url}
										<div class="flex items-baseline justify-between gap-4 py-1 border-b border-(--border-hairline)/50">
											<dt class="text-(--text-muted)">Repository</dt>
											<dd class="text-(--text-main) font-mono truncate max-w-[180px]">
												<a href={project.repository_url} target="_blank" rel="noopener noreferrer" class="hover:underline text-(--accent-sky)">
													{project.repository_url.replace(/^https?:\/\//, '')}
												</a>
											</dd>
										</div>
									{/if}
								</dl>
							</section>
						</div>

					<!-- TAB 2: COMMENTS -->
					{:else if activeTab === 'comments'}
						<div class="flex flex-col gap-6">
							<Comments projectSlug={project.slug} projectOwnerId={project.owner_id} bind:commentCount />
						</div>

					<!-- TAB 3: ACTIVITY -->
					{:else if activeTab === 'activity'}
						<div class="flex flex-col gap-6 max-w-2xl">
							{#if activities.length === 0}
								<p class="text-xs text-(--text-muted) py-6">No recorded activity for this project yet.</p>
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
					{/if}
				</div>
			</div>

			<!-- Right Column: Sticky Project Context Column (shows above tabs on mobile, sticky on desktop) -->
			<aside class="lg:col-span-4 order-1 lg:order-2 lg:sticky lg:top-24 self-start flex flex-col gap-6 pt-2 lg:pt-0 border border-(--border-hairline) lg:border-0 rounded-md lg:rounded-none p-4 lg:p-0 bg-(--bg-surface) lg:bg-transparent">
				<!-- Availability & Health -->
				<div class="flex flex-col gap-3">
					<div class="flex items-center justify-between">
						<span class="text-xs font-semibold text-(--text-main)">Availability</span>
						<div class="flex items-center gap-1.5 text-xs">
							<StatusDot status={projectStatus.dotStatus} showLabel={false} />
							<span class="font-medium {projectStatus.colorClass}">
								{projectStatus.label}
							</span>
						</div>
					</div>

					{#if availability?.uptime_percent != null || availability?.latest_response_time_ms || availability?.last_checked_at}
						<div class="flex items-center gap-2 text-xs text-(--text-muted) flex-wrap">
							{#if availability?.uptime_percent != null}
								<span class="font-mono">{availability.uptime_percent.toFixed(1)}% uptime</span>
							{/if}
							{#if availability?.latest_response_time_ms}
								<span class="opacity-40">·</span>
								<span class="font-mono font-medium {getPingColorClass(availability.latest_response_time_ms)}">
									{availability.latest_response_time_ms} ms
								</span>
							{/if}
							{#if availability?.last_checked_at}
								<span class="opacity-40">·</span>
								<span>Checked {formatRelativeTime(availability.last_checked_at)}</span>
							{/if}
						</div>
					{/if}

					<!-- Multi-Range Bar Sequence -->
					<div class="pt-1">
						<UptimeHistory {availability} status={project.status} variant="compact" showRangeSelector={true} />
					</div>
				</div>

				<!-- Connected Links (External/Companion resources only) -->
				{#if project.repository_url || project.documentation_url || project.demo_url}
					<div class="h-px bg-(--border-hairline)"></div>

					<div class="flex flex-col gap-3">
						<h3 class="font-semibold text-xs text-(--text-main)">Links</h3>
						<div class="flex flex-col gap-2 text-xs">
							{#if project.repository_url}
								{@const repoInfo = detectLinkInfo(project.repository_url, 'Source code')}
								<a
									href={project.repository_url}
									target="_blank"
									rel="noopener noreferrer"
									class="text-(--text-secondary) hover:text-(--text-main) inline-flex items-center justify-between py-1 transition-colors"
								>
									<span class="flex items-center gap-2">
										<svelte:component this={repoInfo.icon} size={14} />
										<span>{repoInfo.label}</span>
									</span>
									<ArrowUpRight size={13} />
								</a>
							{/if}

							{#if project.documentation_url}
								{@const docInfo = detectLinkInfo(project.documentation_url, 'Documentation')}
								<a
									href={project.documentation_url}
									target="_blank"
									rel="noopener noreferrer"
									class="text-(--text-secondary) hover:text-(--text-main) inline-flex items-center justify-between py-1 transition-colors"
								>
									<span class="flex items-center gap-2">
										<svelte:component this={docInfo.icon} size={14} />
										<span>{docInfo.label}</span>
									</span>
									<ArrowUpRight size={13} />
								</a>
							{/if}

							{#if project.demo_url && project.demo_url !== project.public_url}
								{@const demoInfo = detectLinkInfo(project.demo_url, 'Live demo')}
								<a
									href={project.demo_url}
									target="_blank"
									rel="noopener noreferrer"
									class="text-(--text-secondary) hover:text-(--text-main) inline-flex items-center justify-between py-1 transition-colors"
								>
									<span class="flex items-center gap-2">
										<svelte:component this={demoInfo.icon} size={14} />
										<span>{demoInfo.label}</span>
									</span>
									<ArrowUpRight size={13} />
								</a>
							{/if}
						</div>
					</div>
				{/if}

				{#if project.owner}
					<div class="h-px bg-(--border-hairline)"></div>

					<!-- Creator Section -->
					<div class="flex flex-col gap-3">
						<h3 class="font-semibold text-xs text-(--text-main)">Creator</h3>
						<div class="flex items-center gap-3">
							{#if project.owner.avatar_url}
								<img src={project.owner.avatar_url} alt={project.owner.display_name} class="w-9 h-9 rounded-full object-cover border border-(--border-hairline)" />
							{:else}
								<div class="w-9 h-9 rounded-full bg-(--bg-muted) border border-(--border-hairline) flex items-center justify-center font-medium text-xs text-(--text-main)">
									{project.owner.display_name?.charAt(0) || 'U'}
								</div>
							{/if}
							<div class="flex flex-col text-xs">
								<a href="/people/{project.owner.username}" class="font-medium text-(--text-main) hover:underline">
									{project.owner.display_name}
								</a>
								<span class="text-(--text-muted)">@{project.owner.username}</span>
							</div>
						</div>
						{#if project.owner.bio}
							<p class="text-xs text-(--text-secondary) leading-relaxed">
								{project.owner.bio}
							</p>
						{/if}
					</div>
				{/if}

				<div class="h-px bg-(--border-hairline)"></div>

				<!-- Hosted Since -->
				<div class="flex items-center justify-between text-xs text-(--text-muted)">
					<span>Hosted since</span>
					<span class="text-(--text-main) font-mono">{formatDate(project.published_at || project.created_at)}</span>
				</div>
			</aside>
		</div>

		<!-- Report Modal -->
		<ReportModal
			bind:open={reportModalOpen}
			targetType="PROJECT"
			targetId={project.slug}
			targetName={project.name}
		/>
	{/if}
</div>
