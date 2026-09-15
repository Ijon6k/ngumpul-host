<script lang="ts">
	import { onMount } from 'svelte';
	import { commentsApi } from '$lib/api';
	import StatusDot from '$lib/components/ui/StatusDot.svelte';
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

	let { data }: { data: PageData } = $props();

	let copied = $state(false);
	let reportModalOpen = $state(false);
	let activeTab = $state('overview');
	let commentCount = $state(0);
	let loading = $state(false);
	let error: string | null = $state(null);

	let project = $derived(data.project);
	let projectStatus = $derived(resolveProjectStatus(project));
	let availability = $derived(data.availability);
	let activities = $derived(data.activities || []);

	let tabsList = $derived([
		{ id: 'overview', label: 'Overview' },
		{ id: 'comments', label: commentCount > 0 ? `Comments (${commentCount})` : 'Comments' },
		{ id: 'activity', label: 'Activity' }
	]);

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

<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-12 pb-24 flex flex-col gap-6 sm:gap-8">
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
			class="mb-1"
			items={[
				{ label: 'Projects', href: '/projects' },
				{ label: project.name }
			]}
		/>

		<!-- 1. Title & Description (Urutan 1: Judul, Deskripsi) -->
		<header class="flex flex-col gap-3">
			<div class="flex flex-col sm:flex-row sm:items-baseline justify-between gap-3">
				<h1 class="font-sans font-normal text-2xl sm:text-3xl lg:text-4xl text-(--text-main) tracking-[-0.03em] leading-tight text-balance break-words">
					{project.name}
				</h1>
				<div class="self-start sm:self-auto shrink-0">
					<StatusDot status={projectStatus.dotStatus} variant="badge" />
				</div>
			</div>

			{#if project.description}
				<p class="text-sm sm:text-base text-(--text-secondary) max-w-3xl leading-relaxed font-normal break-words">
					{project.description}
				</p>
			{/if}
		</header>

		<!-- 2. Foto Cover (Urutan 2: Foto 16:9 di tengah selebar container) -->
		<div class="w-full aspect-[16/9] rounded-md overflow-hidden border border-(--border-hairline) bg-(--bg-muted) shadow-xs flex">
			<ProjectCover
				src={project.cover_image_url}
				alt={project.name}
				name={project.name}
				aspectRatio="16/9"
				class="w-full h-full !rounded-none !border-none"
			/>
		</div>

		<!-- 3. Tombol-tombol Aksi (Urutan 3: Tombol setelah foto) -->
		<div class="flex items-center gap-2 sm:gap-2.5 flex-wrap">
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
				{@const SourceIcon = sourceInfo.icon}
				<a
					href={project.repository_url}
					target="_blank"
					rel="noopener noreferrer"
					class="btn btn-secondary btn-sm text-xs px-3 py-2 inline-flex items-center gap-1.5"
				>
					<SourceIcon size={14} />
					<span>{sourceInfo.label === 'GitHub' || sourceInfo.label === 'GitLab' ? sourceInfo.label : 'Source'}</span>
					<ArrowUpRight size={11} />
				</a>
			{/if}

			{#if project.documentation_url}
				{@const docInfo = detectLinkInfo(project.documentation_url, 'Documentation')}
				{@const DocIcon = docInfo.icon}
				<a
					href={project.documentation_url}
					target="_blank"
					rel="noopener noreferrer"
					class="btn btn-secondary btn-sm text-xs px-3 py-2 inline-flex items-center gap-1.5"
				>
					<DocIcon size={14} />
					<span>Docs</span>
					<ArrowUpRight size={11} />
				</a>
			{/if}

			<button
				type="button"
				onclick={copyUrl}
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
				onclick={() => (reportModalOpen = true)}
				class="btn btn-secondary btn-sm text-xs px-3 py-2 inline-flex items-center gap-1.5 cursor-pointer text-(--text-muted) hover:text-(--text-main)"
				title="Report this project"
			>
				<Flag size={13} />
				<span>Report</span>
			</button>
		</div>

		<!-- 4. Main Tabs Container (Selebar Foto) -->
		<div class="flex flex-col gap-6 pt-2">
			<Tabs tabs={tabsList} bind:active={activeTab} />

			<!-- TAB 1: OVERVIEW -->
			{#if activeTab === 'overview'}
				<div class="flex flex-col gap-8">
					<!-- 1. Metadata Utama & Creator (At the top, clean minimalist layout, NO card borders) -->
					<section class="flex flex-col gap-5 pt-1">
						<!-- Creator Profile -->
						{#if project.owner}
							<div class="flex items-start gap-4">
								{#if project.owner.avatar_url}
									<img
										src={project.owner.avatar_url}
										alt={project.owner.display_name}
										class="w-11 h-11 rounded-full object-cover border border-(--border-hairline) shrink-0"
									/>
								{:else}
									<div class="w-11 h-11 rounded-full bg-(--accent-soft) text-(--accent-strong) flex items-center justify-center font-bold text-sm shrink-0">
										{project.owner.display_name?.charAt(0) || 'U'}
									</div>
								{/if}
								<div class="flex flex-col min-w-0 grow">
									<div class="flex items-center gap-2 flex-wrap">
										<a href="/people/{project.owner.username}" class="font-sans font-semibold text-sm sm:text-base text-(--text-main) hover:underline">
											{project.owner.display_name}
										</a>
										<span class="text-xs text-(--text-muted)">@{project.owner.username}</span>
									</div>
									{#if project.owner.bio}
										<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed mt-1 max-w-2xl break-words">
											{project.owner.bio}
										</p>
									{/if}
								</div>
							</div>
						{/if}

						<!-- Project Specs / Details Grid (Clean definition list with hairline dividers, no heavy card boxes) -->
						<dl class="grid grid-cols-2 sm:grid-cols-4 gap-4 sm:gap-6 pt-4 border-t border-(--border-hairline) text-xs">
							<div class="flex flex-col gap-1">
								<dt class="text-(--text-muted)">Published</dt>
								<dd class="text-(--text-main) font-medium">{formatDate(project.published_at || project.created_at)}</dd>
							</div>

							<div class="flex flex-col gap-1">
								<dt class="text-(--text-muted)">Hosted since</dt>
								<dd class="text-(--text-main) font-medium">{formatDate(project.created_at)}</dd>
							</div>

							{#if project.public_url}
								<div class="flex flex-col gap-1">
									<dt class="text-(--text-muted)">Public URL</dt>
									<dd class="text-(--text-main) font-mono truncate">
										<a href="/go/{project.slug}" target="_blank" class="text-(--accent-sky) hover:underline break-all">
											{project.public_url.replace(/^https?:\/\//, '')}
										</a>
									</dd>
								</div>
							{/if}

							{#if project.technology_stack && project.technology_stack.length > 0}
								<div class="flex flex-col gap-1">
									<dt class="text-(--text-muted)">Technology</dt>
									<dd class="text-(--text-secondary) font-mono truncate" title={project.technology_stack.join(' · ')}>
										{project.technology_stack.join(' · ')}
									</dd>
								</div>
							{/if}
						</dl>
					</section>

					<!-- 2. Status & Availability (Placed below metadata, minimalist, no boxed card slop) -->
					<section class="flex flex-col gap-3.5 pt-6 border-t border-(--border-hairline)">
						<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
							<div class="flex items-center gap-2">
								<h2 class="font-sans font-medium text-sm text-(--text-main)">Availability & Telemetry</h2>
								<StatusDot status={projectStatus.dotStatus} showLabel={false} />
								<span class="text-xs font-medium {projectStatus.colorClass}">
									{projectStatus.label}
								</span>
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
						</div>

						<div class="pt-1">
							<UptimeHistory {availability} status={project.status} variant="detailed" showRangeSelector={true} />
						</div>
					</section>

					<!-- 3. README / Dokumentasi (At the very bottom, full width) -->
					{#if project.readme}
						<section class="flex flex-col gap-2 pt-4 border-t border-(--border-hairline)">
							<MarkdownView content={project.readme} title="README.md" bordered={true} />
						</section>
					{/if}
				</div>

			<!-- TAB 2: COMMENTS -->
			{:else if activeTab === 'comments'}
				<div class="flex flex-col gap-6">
					<Comments projectSlug={project.slug} projectOwnerId={project.owner_id} bind:commentCount />
				</div>

			<!-- TAB 3: ACTIVITY -->
			{:else if activeTab === 'activity'}
				<div class="flex flex-col gap-6">
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

		<!-- Report Modal -->
		<ReportModal
			bind:open={reportModalOpen}
			targetType="PROJECT"
			targetId={project.slug}
			targetName={project.name}
		/>
	{/if}
</div>
