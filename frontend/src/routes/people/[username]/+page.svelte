<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import ProjectCard from '$lib/components/ProjectCard.svelte';

	let member: any = null;
	let projects: any[] = [];
	let loading = true;
	let error: string | null = null;

	onMount(async () => {
		const username = $page.params.username;
		try {
			const res = await api.get(`/users/${username}`);
			member = res.data?.user;
			projects = res.data?.projects || [];
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	});

	function formatJoinDate(iso: string) {
		try {
			const d = new Date(iso);
			return d.toLocaleDateString([], { month: 'long', year: 'numeric' });
		} catch (_) {
			return '';
		}
	}
</script>

<svelte:head>
	<title>{member?.display_name ? `${member.display_name} (@${member.username}) · Ngumpul Host` : 'Member Profile · Ngumpul Host'}</title>
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-12 pb-24 flex flex-col gap-8 sm:gap-10">
	{#if loading}
		<div class="py-24 text-center text-(--text-muted) font-mono text-xs">
			<p>Resolving member profile...</p>
		</div>
	{:else if error || !member}
		<div class="py-16 px-6 text-center border border-dashed border-(--border-hairline) rounded-md max-w-md mx-auto w-full">
			<h2 class="font-display text-xl font-bold text-(--text-main) mb-1.5">Member Not Located</h2>
			<p class="text-xs text-(--text-secondary)">{error || 'This user profile could not be found.'}</p>
			<a href="/people" class="btn btn-secondary btn-sm mt-4">← Back to members</a>
		</div>
	{:else}
		<!-- Breadcrumb -->
		<div class="flex items-center gap-2 text-xs text-(--text-muted)">
			<a href="/people" class="hover:text-(--text-main) transition-colors">People</a>
			<span>/</span>
			<span class="text-(--text-main) font-medium">@{member.username}</span>
		</div>

		<!-- Profile Editorial Header -->
		<header class="flex items-start gap-6 pb-8 border-b border-(--border-hairline) flex-wrap sm:flex-nowrap">
			{#if member.avatar_url}
				<img src={member.avatar_url} alt={member.display_name} class="w-16 h-16 sm:w-20 sm:h-20 rounded-full object-cover shrink-0" />
			{:else}
				<div class="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-(--accent-soft) text-(--accent-strong) font-display text-2xl sm:text-3xl font-extrabold flex items-center justify-center shrink-0">
					{member.display_name.charAt(0)}
				</div>
			{/if}

			<div class="flex flex-col gap-1.5 grow min-w-0">
				<div class="flex items-center gap-3 flex-wrap">
					<h1 class="font-display font-extrabold text-2xl sm:text-3xl text-(--text-main)">{member.display_name}</h1>
					{#if member.role === 'ADMIN'}
						<span class="text-xs text-(--accent-strong) font-medium px-2 py-0.5 rounded-md bg-(--accent-soft)">Host Maintainer</span>
					{/if}
				</div>

				<span class="text-xs text-(--text-muted)">@{member.username} · Member since {formatJoinDate(member.created_at)}</span>

				{#if member.bio}
					<p class="text-sm text-(--text-secondary) leading-relaxed max-w-2xl mt-2">{member.bio}</p>
				{/if}
			</div>
		</header>

		<!-- Member's Hosted Projects Section -->
		<section class="flex flex-col gap-6">
			<div class="flex items-baseline justify-between">
				<h2 class="font-display font-bold text-xl text-(--text-main)">Hosted Work</h2>
				<span class="text-xs text-(--text-muted)">{projects.length} {projects.length === 1 ? 'project' : 'projects'}</span>
			</div>

			{#if projects.length === 0}
				<div class="py-14 px-4 text-center border border-dashed border-(--border-hairline) rounded-md text-(--text-muted)">
					<p class="text-xs">No active projects registered under this account.</p>
				</div>
			{:else}
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each projects as project (project.id)}
						<ProjectCard {project} />
					{/each}
				</div>
			{/if}
		</section>
	{/if}
</div>
