<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/auth';
	import { api } from '$lib/api';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import {
		User,
		Gear,
		ArrowUpRight,
		Calendar,
		Folder,
		EnvelopeSimple
	} from 'phosphor-svelte';

	let myProjects: any[] = [];
	let loading = true;

	onMount(async () => {
		try {
			const res = await api.get('/me/projects');
			myProjects = res.data?.projects || [];
		} catch {
			// Silently fail
		} finally {
			loading = false;
		}
	});

	function formatJoinDate(iso?: string) {
		if (!iso) return '2026';
		try {
			return new Date(iso).toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
		} catch {
			return iso;
		}
	}
</script>

<svelte:head>
	<title>{$user?.display_name || 'My Profile'} · Ngumpul Host</title>
</svelte:head>

<div class="p-6 sm:p-8 lg:p-10 max-w-4xl w-full flex flex-col gap-8">
	<!-- Page Header -->
	<header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-(--border-hairline)">
		<div class="flex flex-col gap-1">
			<h1 class="font-sans font-semibold text-2xl sm:text-3xl text-(--text-main) tracking-tight">
				Member Profile
			</h1>
			<p class="text-xs text-(--text-secondary) leading-relaxed">
				How your developer identity appears across the community roster and project showcases.
			</p>
		</div>

		<div class="flex items-center gap-2">
			{#if $user?.username}
				<a
					href="/people/{$user.username}"
					target="_blank"
					rel="noreferrer"
					class="btn btn-secondary btn-sm text-xs px-3 py-1.5 inline-flex items-center gap-1.5"
				>
					<span>Public View</span>
					<ArrowUpRight size={13} />
				</a>
			{/if}
			<a
				href="/me/settings"
				class="btn btn-primary btn-sm text-xs px-3.5 py-1.5 inline-flex items-center gap-1.5"
			>
				<Gear size={13} weight="bold" />
				<span>Edit Profile</span>
			</a>
		</div>
	</header>

	{#if $user}
		<!-- Profile Card Preview -->
		<div class="p-6 sm:p-8 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col sm:flex-row items-start gap-6 shadow-xs">
			{#if $user.avatar_url}
				<img
					src={$user.avatar_url}
					alt={$user.display_name}
					class="w-20 h-20 rounded-full object-cover border border-(--border-hairline) shrink-0"
				/>
			{:else}
				<div class="w-20 h-20 rounded-full bg-(--bg-muted) text-(--text-secondary) flex items-center justify-center font-bold text-2xl border border-(--border-hairline) shrink-0">
					{($user.display_name || 'U').charAt(0)}
				</div>
			{/if}

			<div class="flex flex-col gap-3 flex-1 min-w-0">
				<div class="flex flex-col gap-0.5">
					<div class="flex items-center gap-2 flex-wrap">
						<h2 class="font-sans font-semibold text-xl text-(--text-main)">
							{$user.display_name}
						</h2>
						{#if $user.role === 'ADMIN'}
							<span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/30">
								Operator
							</span>
						{/if}
					</div>
					<span class="text-xs font-mono text-(--text-muted)">@{$user.username}</span>
				</div>

				{#if $user.bio}
					<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">
						{$user.bio}
					</p>
				{:else}
					<p class="text-xs text-(--text-muted) italic">
						No bio added yet. Click 'Edit Profile' to add a short note about what you build.
					</p>
				{/if}

				<div class="flex items-center gap-5 pt-3 border-t border-(--border-hairline) text-xs text-(--text-muted) flex-wrap">
					{#if $user.email}
						<div class="flex items-center gap-1.5">
							<EnvelopeSimple size={14} />
							<span>{$user.email}</span>
						</div>
					{/if}

					<div class="flex items-center gap-1.5">
						<Calendar size={14} />
						<span>Joined {formatJoinDate($user.created_at)}</span>
					</div>

					<div class="flex items-center gap-1.5">
						<Folder size={14} />
						<span>{myProjects.length} hosted {myProjects.length === 1 ? 'project' : 'projects'}</span>
					</div>
				</div>
			</div>
		</div>

		<!-- My Published Projects List -->
		<section class="flex flex-col gap-4">
			<div class="flex items-center justify-between pb-2 border-b border-(--border-hairline)">
				<h3 class="font-sans font-semibold text-base text-(--text-main)">
					Authored Projects
				</h3>
				<a href="/me/projects" class="text-xs text-(--text-secondary) hover:text-(--text-main) transition-colors">
					Manage all →
				</a>
			</div>

			{#if loading}
				<div class="p-8 text-center text-xs text-(--text-muted) bg-(--bg-surface) border border-(--border-hairline) rounded-md">
					Loading projects...
				</div>
			{:else if myProjects.length === 0}
				<div class="p-8 text-center bg-(--bg-surface) border border-dashed border-(--border-hairline) rounded-md text-xs text-(--text-muted)">
					No projects hosted yet under this profile.
				</div>
			{:else}
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					{#each myProjects as proj (proj.id)}
						<a
							href="/me/projects/{proj.id}"
							class="p-4 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex items-center justify-between gap-3 hover:border-(--border-subtle) transition-colors group"
						>
							<div class="flex items-center gap-3 min-w-0">
								<div class="w-12 h-9 rounded-xs overflow-hidden shrink-0 border border-(--border-hairline) bg-(--bg-muted)">
									<ProjectCover
										src={proj.cover_image_url}
										alt={proj.name}
										name={proj.name}
										aspectRatio="4/3"
									/>
								</div>
								<div class="flex flex-col min-w-0">
									<span class="font-medium text-xs text-(--text-main) group-hover:text-(--accent-sky) transition-colors truncate">
										{proj.name}
									</span>
									<span class="text-[11px] text-(--text-muted) font-mono">/{proj.slug}</span>
								</div>
							</div>
							<StatusDot status={proj.status} />
						</a>
					{/each}
				</div>
			{/if}
		</section>
	{/if}
</div>
