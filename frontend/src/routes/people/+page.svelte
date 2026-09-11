<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	let members: any[] = [];
	let loading = true;

	onMount(async () => {
		try {
			const res = await api.get('/users');
			members = res.data?.members || [];
		} catch (err) {
			console.error('Failed to load members:', err);
		} finally {
			loading = false;
		}
	});

	function formatJoinDate(iso: string) {
		try {
			const d = new Date(iso);
			return d.toLocaleDateString([], { month: 'short', year: 'numeric' });
		} catch (_) {
			return '';
		}
	}
</script>

<svelte:head>
	<title>People · Ngumpul Host</title>
	<meta name="description" content="The developers, creators, and tinkerers hosting their work on Ngumpul Host." />
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-12 pb-24 flex flex-col gap-8 sm:gap-10">
	<header class="max-w-2xl">
		<h1 class="font-sans font-normal text-3xl sm:text-4xl text-(--text-main) tracking-[-0.03em] mb-2">People</h1>
		<p class="text-sm sm:text-base text-(--text-secondary) leading-relaxed font-normal">
			The developers, creators, and tinkerers hosting their work on this system.
		</p>
	</header>

	{#if loading}
		<div class="py-20 text-center text-(--text-muted) text-xs">
			<p>Loading member directory...</p>
		</div>
	{:else if members.length === 0}
		<div class="py-16 px-4 text-center border border-dashed border-(--border-hairline) rounded-md text-(--text-muted)">
			<p class="text-sm">No members registered yet.</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			{#each members as member (member.id)}
				<a
					href="/people/{member.username}"
					class="group p-6 bg-(--bg-surface) border border-(--border-hairline) hover:border-(--text-muted) rounded-md flex flex-col justify-between gap-4 transition-colors"
				>
					<div class="flex items-start gap-4">
						{#if member.avatar_url}
							<img src={member.avatar_url} alt={member.display_name} class="w-12 h-12 rounded-full object-cover shrink-0" />
						{:else}
							<div class="w-12 h-12 rounded-full bg-(--accent-soft) text-(--accent-strong) font-bold flex items-center justify-center text-lg shrink-0">
								{member.display_name.charAt(0)}
							</div>
						{/if}
						<div class="min-w-0">
							<h3 class="font-display font-bold text-base text-(--text-main) group-hover:text-(--accent-strong) transition-colors truncate">
								{member.display_name}
							</h3>
							<span class="text-xs text-(--text-muted)">@{member.username}</span>
							{#if member.bio}
								<p class="text-xs text-(--text-secondary) line-clamp-2 leading-relaxed mt-1.5">{member.bio}</p>
							{/if}
						</div>
					</div>

					<div class="flex items-center justify-between pt-3 border-t border-(--border-hairline) text-xs text-(--text-muted)">
						<span class="text-(--text-main) font-medium">{member.projects_count || 0} {member.projects_count === 1 ? 'project' : 'projects'}</span>
						<span>Joined {formatJoinDate(member.created_at)}</span>
					</div>
				</a>
			{/each}
		</div>
	{/if}
</div>
