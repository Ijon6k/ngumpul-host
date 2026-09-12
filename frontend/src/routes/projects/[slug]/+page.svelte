<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import Comments from '$lib/components/Comments.svelte';
	import ReportModal from '$lib/components/ReportModal.svelte';
	import ProjectAvailability from '$lib/components/ProjectAvailability.svelte';
	import { Tabs, Timeline, TimelineItem } from '$lib/components/ui';
	import {
		ArrowUpRight,
		GithubLogo,
		BookOpen,
		Globe,
		User,
		Calendar,
		Flag,
		ShareNetwork,
		Check
	} from 'phosphor-svelte';

	let project: any = null;
	let activities: any[] = [];
	let availability: any = null;
	let loading = true;
	let error: string | null = null;
	let copied = false;
	let reportModalOpen = false;
	let activeTab = 'overview';

	const tabsList = [
		{ id: 'overview', label: 'Overview' },
		{ id: 'activity', label: 'Activity' }
	];

	function copyUrl() {
		if (!navigator?.clipboard) return;
		navigator.clipboard.writeText(window.location.href);
		copied = true;
		setTimeout(() => {
			copied = false;
		}, 2000);
	}

	onMount(async () => {
		const slug = $page.params.slug;
		try {
			const res = await api.get(`/projects/${slug}`);
			project = res.data?.project;
			activities = res.data?.activities || [];
			availability = res.data?.availability;
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	});

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

<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-12 pb-24 flex flex-col">
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
		<nav aria-label="Breadcrumb" class="flex items-center gap-2 text-xs text-(--text-muted) mb-6">
			<a href="/projects" class="hover:text-(--text-main) transition-colors">Projects</a>
			<span class="opacity-40">/</span>
			<span class="text-(--text-main) font-medium">{project.name}</span>
		</nav>

		<!-- Public Project Header (Editorial Showcase) -->
		<header class="flex flex-col gap-3">
			<!-- Status & Node Info -->
			<div class="flex items-center gap-2.5 text-xs text-(--text-muted)">
				<StatusDot status={project.status} />
				<span>·</span>
				<span>
					{#if project.hosting_type === 'HOSTED_HERE'}
						Hosted on community server
					{:else}
						External registered service
					{/if}
				</span>
			</div>

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

			<!-- Creator and Hosted Date -->
			<div class="flex items-center flex-wrap gap-4 pt-1 text-xs text-(--text-muted)">
				{#if project.owner}
					<a href="/people/{project.owner.username}" class="flex items-center gap-1.5 font-medium text-(--text-main) hover:underline">
						<User size={14} />
						<span>by {project.owner.display_name}</span>
					</a>
					<span class="opacity-40">·</span>
				{/if}

				<span class="flex items-center gap-1.5">
					<Calendar size={14} />
					<span>Hosted since {formatDate(project.published_at || project.created_at)}</span>
				</span>
			</div>

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
					<a
						href={project.repository_url}
						target="_blank"
						rel="noopener noreferrer"
						class="btn btn-secondary btn-sm text-xs px-3 py-2 inline-flex items-center gap-1.5"
					>
						<GithubLogo size={14} />
						<span>Source</span>
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

		<!-- Large Visual Cover Area -->
		<div class="w-full aspect-video sm:aspect-[21/9] max-h-[400px] rounded-md overflow-hidden border border-(--border-hairline) bg-(--bg-muted) shadow-xs my-8">
			<ProjectCover
				src={project.cover_image_url}
				alt={project.name}
				name={project.name}
				aspectRatio="16/9"
			/>
		</div>

		<!-- Minimal Typographic Tabs: Overview | Activity -->
		<div class="mb-8">
			<Tabs tabs={tabsList} bind:active={activeTab} />
		</div>

		<!-- TAB 1: OVERVIEW (About, Built with, Availability, Links, Comments) -->
		{#if activeTab === 'overview'}
			<div class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
				<!-- Main Column (Left, 8 cols) -->
				<div class="lg:col-span-8 flex flex-col gap-10">
					<!-- Project Description -->
					<section class="flex flex-col gap-3">
						<h2 class="font-sans font-medium text-base text-(--text-main)">
							About
						</h2>
						<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">
							{project.description || 'No detailed description provided for this project.'}
						</p>
					</section>

					<!-- Built With -->
					{#if project.technology_stack && project.technology_stack.length > 0}
						<section class="flex flex-col gap-3">
							<h2 class="font-sans font-medium text-base text-(--text-main)">
								Built with
							</h2>
							<p class="text-xs text-(--text-secondary) font-mono leading-relaxed">
								{project.technology_stack.join(' · ')}
							</p>
						</section>
					{/if}

					<!-- Community Comments (Stay on Overview) -->
					<Comments projectSlug={project.slug} projectOwnerId={project.owner_id} />
				</div>

				<!-- Sidebar Column (Right, 4 cols) -->
				<aside class="lg:col-span-4 flex flex-col gap-8">
					<!-- Availability Probe -->
					<ProjectAvailability {availability} status={project.status} />

					<!-- Connected Links -->
					<div class="flex flex-col gap-3 text-xs">
						<h3 class="font-medium text-xs text-(--text-main)">Links</h3>
						<div class="flex flex-col gap-2">
							{#if project.public_url}
								<a
									href="/go/{project.slug}"
									target="_blank"
									rel="noopener noreferrer"
									class="text-(--accent-sky) hover:underline inline-flex items-center justify-between py-1"
								>
									<span class="truncate font-mono">{project.public_url}</span>
									<ArrowUpRight size={12} />
								</a>
							{/if}

							{#if project.repository_url}
								<a
									href={project.repository_url}
									target="_blank"
									rel="noopener noreferrer"
									class="text-(--text-secondary) hover:text-(--text-main) inline-flex items-center justify-between py-1 transition-colors"
								>
									<span class="flex items-center gap-2">
										<GithubLogo size={13} />
										<span>Source code</span>
									</span>
									<ArrowUpRight size={12} />
								</a>
							{/if}

							{#if project.documentation_url}
								<a
									href={project.documentation_url}
									target="_blank"
									rel="noopener noreferrer"
									class="text-(--text-secondary) hover:text-(--text-main) inline-flex items-center justify-between py-1 transition-colors"
								>
									<span class="flex items-center gap-2">
										<BookOpen size={13} />
										<span>Documentation</span>
									</span>
									<ArrowUpRight size={12} />
								</a>
							{/if}

							{#if project.demo_url && project.demo_url !== project.public_url}
								<a
									href={project.demo_url}
									target="_blank"
									rel="noopener noreferrer"
									class="text-(--text-secondary) hover:text-(--text-main) inline-flex items-center justify-between py-1 transition-colors"
								>
									<span class="flex items-center gap-2">
										<Globe size={13} />
										<span>Live demo</span>
									</span>
									<ArrowUpRight size={12} />
								</a>
							{/if}
						</div>
					</div>

					<!-- Creator Snippet -->
					{#if project.owner}
						<div class="flex flex-col gap-3 pt-6 border-t border-(--border-hairline) text-xs">
							<h3 class="font-medium text-xs text-(--text-main)">Creator</h3>
							<div class="flex items-center gap-3">
								{#if project.owner.avatar_url}
									<img src={project.owner.avatar_url} alt={project.owner.display_name} class="w-9 h-9 rounded-full object-cover border border-(--border-hairline)" />
								{:else}
									<div class="w-9 h-9 rounded-full bg-(--bg-muted) border border-(--border-hairline) flex items-center justify-center font-medium text-xs text-(--text-main)">
										{project.owner.display_name?.charAt(0) || 'U'}
									</div>
								{/if}
								<div class="flex flex-col">
									<a href="/people/{project.owner.username}" class="font-medium text-(--text-main) hover:underline">
										{project.owner.display_name}
									</a>
									<span class="text-[11px] text-(--text-muted)">@{project.owner.username}</span>
								</div>
							</div>
							{#if project.owner.bio}
								<p class="text-xs text-(--text-secondary) leading-relaxed">
									{project.owner.bio}
								</p>
							{/if}
						</div>
					{/if}
				</aside>
			</div>

		<!-- TAB 2: ACTIVITY (Timeline, Not Table) -->
		{:else if activeTab === 'activity'}
			<div class="flex flex-col gap-6 max-w-2xl">
				<div class="flex flex-col gap-1 pb-2">
					<h2 class="font-sans font-medium text-base text-(--text-main)">Project activity</h2>
					<p class="text-xs text-(--text-secondary)">Recorded change log and operational events.</p>
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
		{/if}

		<!-- Report Modal -->
		<ReportModal
			bind:open={reportModalOpen}
			targetType="PROJECT"
			targetId={project.slug}
			targetName={project.name}
		/>
	{/if}
</div>
