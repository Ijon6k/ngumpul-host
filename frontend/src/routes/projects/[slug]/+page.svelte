<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import {
		Globe,
		GithubLogo,
		ArrowUpRight,
		ClockCounterClockwise,
		CaretLeft,
		Calendar,
		User,
		BookOpen,
		Copy,
		Check
	} from 'phosphor-svelte';

	let project: any = null;
	let activities: any[] = [];
	let loading = true;
	let error: string | null = null;
	let copied = false;

	function copyUrl(url: string) {
		if (!navigator?.clipboard) return;
		navigator.clipboard.writeText(url);
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
			return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
		} catch {
			return dateStr;
		}
	}

	function getActivityMessage(act: any, defaultName: string): string {
		if (act.metadata?.title) return act.metadata.title;
		if (act.metadata?.message) return act.metadata.message;
		switch (act.type) {
			case 'PROJECT_PUBLISHED':
				return `${defaultName} was published to the catalog`;
			case 'PROJECT_UPDATED':
				return `Project metadata for ${defaultName} was updated`;
			case 'PROJECT_ONLINE':
				return `${defaultName} is online and operational`;
			case 'PROJECT_OFFLINE':
				return `${defaultName} is currently unavailable`;
			default:
				return `Update recorded for ${defaultName}`;
		}
	}
</script>

<svelte:head>
	<title>{project?.name ? `${project.name} — Ngumpul Host` : 'Project Showcase — Ngumpul Host'}</title>
	<meta name="description" content={project?.description || 'Independent software project hosted on Ngumpul Host.'} />
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-12 pb-24 flex flex-col gap-8">
	{#if loading}
		<div class="py-32 text-center text-(--text-muted) text-sm flex flex-col items-center justify-center gap-3">
			<div class="w-6 h-6 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
			<p class="text-xs">Loading project showcase...</p>
		</div>
	{:else if error || !project}
		<div class="py-20 px-6 text-center border border-dashed border-(--border-hairline) rounded-md max-w-lg mx-auto w-full bg-(--bg-surface) flex flex-col items-center gap-4">
			<h2 class="font-sans font-semibold text-xl text-(--text-main)">Project Not Found</h2>
			<p class="text-xs text-(--text-secondary) leading-relaxed">
				{error || 'The requested project might not be published yet, or the link has expired.'}
			</p>
			<a href="/projects" class="btn btn-secondary text-xs px-4 py-2">
				← Back to Catalog
			</a>
		</div>
	{:else}
		<!-- Breadcrumb & Top Navigation Controls -->
		<div class="flex items-center justify-between flex-wrap gap-4 text-xs">
			<nav aria-label="Breadcrumb" class="flex items-center gap-2 text-(--text-muted)">
				<a href="/projects" class="hover:text-(--text-main) transition-colors font-medium">Projects</a>
				<span class="opacity-40">/</span>
				<span class="text-(--text-main) font-medium">{project.name}</span>
			</nav>

			<div class="flex items-center gap-2.5">
				<a
					href="/projects"
					class="text-(--text-secondary) hover:text-(--text-main) text-xs flex items-center gap-1.5 transition-colors px-3 py-1.5 rounded-md border border-(--border-hairline) bg-(--bg-surface)"
				>
					<CaretLeft size={12} weight="bold" />
					<span>Catalog</span>
				</a>

				<button
					on:click={() => copyUrl(project.public_url || window.location.href)}
					class="text-(--text-secondary) hover:text-(--text-main) text-xs flex items-center gap-1.5 transition-colors px-3 py-1.5 rounded-md border border-(--border-hairline) bg-(--bg-surface) cursor-pointer"
					title="Copy Project Link"
				>
					{#if copied}
						<Check size={12} weight="bold" class="text-emerald-500" />
						<span class="text-emerald-500 font-medium">Copied</span>
					{:else}
						<Copy size={12} weight="regular" />
						<span>Share</span>
					{/if}
				</button>
			</div>
		</div>

		<!-- Project Header: Editorial Showcase -->
		<header class="flex flex-col gap-3">
			<div class="flex items-center gap-3">
				<StatusDot status={project.status} />
				<span class="text-xs text-(--text-muted)">
					{#if project.hosting_type === 'HOSTED_HERE'}
						Hosted on community server
					{:else}
						External registered service
					{/if}
				</span>
			</div>

			<h1 class="font-sans font-semibold text-3xl sm:text-4xl text-(--text-main) tracking-[-0.03em] leading-tight">
				{project.name}
			</h1>

			{#if project.description}
				<p class="text-sm sm:text-base text-(--text-secondary) max-w-3xl leading-relaxed font-normal">
					{project.description}
				</p>
			{/if}

			<div class="flex items-center flex-wrap gap-4 pt-2 text-xs text-(--text-muted)">
				{#if project.owner}
					<a href="/people/{project.owner.username}" class="flex items-center gap-1.5 font-medium text-(--text-main) hover:underline">
						<User size={14} />
						<span>{project.owner.display_name}</span>
					</a>
				{/if}

				<span class="opacity-40">·</span>

				<span class="flex items-center gap-1.5">
					<Calendar size={14} />
					<span>{project.published_at ? `Hosted since ${formatDate(project.published_at)}` : `Added ${formatDate(project.created_at)}`}</span>
				</span>
			</div>
		</header>

		<!-- Large Project Artwork (Consistent 16:9 Showcase Cover with fallback) -->
		<div class="relative w-full aspect-video sm:aspect-[21/9] max-h-[380px] rounded-md overflow-hidden border border-(--border-hairline) bg-(--bg-muted) shadow-xs">
			<ProjectCover
				src={project.cover_image_url}
				alt={project.name}
				name={project.name}
				aspectRatio="16/9"
			/>
		</div>

		<!-- Main Asymmetric 2-Column Content Layout -->
		<div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-start pt-2">
			<!-- Main Column (Left): About, Technology, Project Activity -->
			<div class="lg:col-span-8 flex flex-col gap-10">
				<!-- Section 1: About the project -->
				<section class="flex flex-col gap-3">
					<h2 class="font-sans font-semibold text-base sm:text-lg text-(--text-main) tracking-[-0.01em]">
						About the Project
					</h2>
					<div class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed flex flex-col gap-3">
						<p>{project.description || 'No additional detailed description provided for this project.'}</p>
					</div>
				</section>

				<!-- Section 2: Technology Stack -->
				{#if project.technology_stack && project.technology_stack.length > 0}
					<section class="flex flex-col gap-3">
						<h2 class="font-sans font-semibold text-base sm:text-lg text-(--text-main) tracking-[-0.01em]">
							Built With
						</h2>
						<div class="flex items-center flex-wrap gap-2">
							{#each project.technology_stack as tech}
								<span class="text-xs font-mono bg-(--bg-surface) text-(--text-main) px-2.5 py-1 rounded-sm border border-(--border-hairline)">
									{tech}
								</span>
							{/each}
						</div>
					</section>
				{/if}

				<!-- Section 3: Honest Project Activity History -->
				<section class="flex flex-col gap-4">
					<div class="flex items-center justify-between pb-2 border-b border-(--border-hairline)">
						<h2 class="font-sans font-semibold text-base sm:text-lg text-(--text-main) tracking-[-0.01em] flex items-center gap-2">
							<ClockCounterClockwise size={18} />
							<span>Project History</span>
						</h2>
						<span class="text-[11px] text-(--text-muted)">Recorded activities</span>
					</div>

					{#if activities.length === 0}
						<div class="p-6 bg-(--bg-surface) border border-dashed border-(--border-hairline) rounded-md text-center text-xs text-(--text-muted)">
							No change history has been recorded for this project yet.
						</div>
					{:else}
						<div class="flex flex-col border border-(--border-hairline) rounded-md bg-(--bg-surface) divide-y divide-(--border-hairline)">
							{#each activities as act (act.id)}
								<div class="p-4 flex items-start justify-between gap-4">
									<div class="flex flex-col gap-1">
										<p class="text-xs font-medium text-(--text-main)">
											{getActivityMessage(act, project.name)}
										</p>
										{#if act.actor_name}
											<p class="text-[11px] text-(--text-muted)">
												by {act.actor_name}
											</p>
										{/if}
									</div>
									<span class="text-[11px] text-(--text-muted) shrink-0 font-mono">
										{formatRelativeTime(act.created_at)}
									</span>
								</div>
							{/each}
						</div>
					{/if}
				</section>
			</div>

			<!-- Sidebar Column (Right): Access Links & Creator Info -->
			<aside class="lg:col-span-4 flex flex-col gap-6">
				<!-- Action Box: Launch / Visit -->
				<div class="p-5 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-4">
					<div class="flex flex-col gap-1">
						<span class="text-xs font-semibold text-(--text-main)">Access Link</span>
						<span class="text-[11px] text-(--text-secondary)">Open the application directly in your browser.</span>
					</div>

					{#if project.public_url}
						<a
							href={project.public_url}
							target="_blank"
							rel="noopener noreferrer"
							class="btn btn-primary w-full py-2.5 text-xs text-center flex items-center justify-center gap-1.5 font-medium"
						>
							<span>Open Application</span>
							<ArrowUpRight size={14} weight="bold" />
						</a>
					{:else}
						<div class="p-3 bg-(--bg-muted) rounded-md text-xs text-(--text-muted) text-center">
							Public URL has not been configured yet.
						</div>
					{/if}

					<div class="flex flex-col gap-2 pt-3 border-t border-(--border-hairline) text-xs">
						{#if project.repository_url}
							<a
								href={project.repository_url}
								target="_blank"
								rel="noopener noreferrer"
								class="text-(--text-secondary) hover:text-(--text-main) flex items-center justify-between py-1 transition-colors"
							>
								<span class="flex items-center gap-2">
									<GithubLogo size={14} />
									<span>Source Code</span>
								</span>
								<ArrowUpRight size={12} />
							</a>
						{/if}

						{#if project.documentation_url}
							<a
								href={project.documentation_url}
								target="_blank"
								rel="noopener noreferrer"
								class="text-(--text-secondary) hover:text-(--text-main) flex items-center justify-between py-1 transition-colors"
							>
								<span class="flex items-center gap-2">
									<BookOpen size={14} />
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
								class="text-(--text-secondary) hover:text-(--text-main) flex items-center justify-between py-1 transition-colors"
							>
								<span class="flex items-center gap-2">
									<Globe size={14} />
									<span>Live Demo</span>
								</span>
								<ArrowUpRight size={12} />
							</a>
						{/if}
					</div>
				</div>

				<!-- Creator Card -->
				{#if project.owner}
					<div class="p-5 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-3">
						<span class="text-xs font-semibold text-(--text-main)">Project Creator</span>
						<div class="flex items-center gap-3">
							{#if project.owner.avatar_url}
								<img src={project.owner.avatar_url} alt={project.owner.display_name} class="w-10 h-10 rounded-full object-cover border border-(--border-hairline)" />
							{:else}
								<div class="w-10 h-10 rounded-full bg-(--bg-muted) border border-(--border-hairline) flex items-center justify-center font-medium text-xs text-(--text-main)">
									{project.owner.display_name?.charAt(0) || 'U'}
								</div>
							{/if}
							<div class="flex flex-col">
								<span class="text-xs font-semibold text-(--text-main)">{project.owner.display_name}</span>
								<span class="text-[11px] text-(--text-muted)">@{project.owner.username}</span>
							</div>
						</div>

						{#if project.owner.bio}
							<p class="text-xs text-(--text-secondary) leading-relaxed pt-1">
								{project.owner.bio}
							</p>
						{/if}

						<a
							href="/people/{project.owner.username}"
							class="btn btn-secondary w-full py-2 mt-1 text-xs text-center"
						>
							View Creator Profile
						</a>
					</div>
				{/if}

				<!-- Calm Node Info (Privacy compliant, no cgroups or internal telemetry leaks) -->
				<div class="p-4 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-[11px] text-(--text-muted) flex flex-col gap-1.5">
					<span class="font-medium text-(--text-main)">Community Hosting Node</span>
					<p class="leading-relaxed">
						This project runs on Ngumpul Host shared community infrastructure on node Jakarta (cgk01).
					</p>
				</div>
			</aside>
		</div>
	{/if}
</div>
