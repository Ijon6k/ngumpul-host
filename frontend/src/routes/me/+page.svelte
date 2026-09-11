<script lang="ts">
	import { onMount } from 'svelte';
	import { user, logout } from '$lib/stores/auth';
	import { api } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import { Plus, Gear, SignOut, ArrowRight, ArrowSquareOut } from 'phosphor-svelte';

	let myProjects: any[] = [];
	let myRequests: any[] = [];
	let notifications: any[] = [];
	let loading = true;

	onMount(async () => {
		try {
			const [projRes, reqRes, notifRes] = await Promise.all([
				api.get('/me/projects'),
				api.get('/me/hosting-requests'),
				api.get('/me/notifications')
			]);
			myProjects = projRes.data?.projects || [];
			myRequests = reqRes.data?.requests || [];
			notifications = notifRes.data?.notifications || [];
		} catch (e) {
			console.error('Failed to load personal space data:', e);
		} finally {
			loading = false;
		}
	});

	function getGreeting() {
		const hour = new Date().getHours();
		if (hour < 12) return 'Good morning';
		if (hour < 18) return 'Good afternoon';
		return 'Good evening';
	}
</script>

<svelte:head>
	<title>My Workspace · Ngumpul Host</title>
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-12 pb-24 flex flex-col gap-8 sm:gap-9">
	{#if $user}
		<header class="flex flex-col sm:flex-row sm:items-center justify-between gap-6 pb-6 border-b border-(--border-hairline)">
			<div class="flex items-center gap-4">
				{#if $user.avatar_url}
					<img src={$user.avatar_url} alt={$user.display_name} class="w-12 h-12 rounded-full object-cover border border-(--border-hairline)" />
				{:else}
					<div class="w-12 h-12 rounded-full bg-(--accent-soft) text-(--accent-strong) font-bold flex items-center justify-center text-base border border-(--border-hairline)">
						{($user.display_name || 'U').slice(0, 2).toUpperCase()}
					</div>
				{/if}
				<div>
					<h1 class="font-display font-bold text-2xl sm:text-3xl text-(--text-main) tracking-tight">
						{getGreeting()}, {$user.display_name}.
					</h1>
					<p class="text-xs text-(--text-muted) mt-0.5">
						@{$user.username} · Member since {$user.created_at ? new Date($user.created_at).getFullYear() : '2026'}
					</p>
				</div>
			</div>

			<div class="flex items-center gap-2.5 flex-wrap">
				<a href="/me/requests" class="btn btn-primary btn-sm inline-flex items-center gap-1.5">
					<Plus size={14} weight="bold" />
					<span>Request Project</span>
				</a>
				<a href="/me/settings" class="btn btn-secondary btn-sm inline-flex items-center gap-1.5">
					<Gear size={14} weight="regular" />
					<span>Settings</span>
				</a>
				<button type="button" class="btn btn-secondary btn-sm cursor-pointer inline-flex items-center gap-1.5" on:click={logout}>
					<SignOut size={14} weight="regular" />
					<span>Sign out</span>
				</button>
			</div>
		</header>

		<div class="grid grid-cols-1 lg:grid-cols-[2fr_1fr] gap-8">
			<!-- Main Column: Projects & Requests -->
			<div class="flex flex-col gap-8">
				<!-- Your Projects -->
				<section class="flex flex-col gap-4">
					<div class="flex items-baseline justify-between">
						<h2 class="font-display text-lg font-bold text-(--text-main)">Hosted Applications</h2>
						<a href="/me/projects" class="text-xs text-(--accent-strong) hover:underline font-medium inline-flex items-center gap-1">
							<span>Manage details</span>
							<ArrowRight size={12} weight="bold" />
						</a>
					</div>

					{#if myProjects.length === 0}
						<div class="p-8 text-center bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col items-center gap-2">
							<p class="text-xs text-(--text-secondary)">You do not have any published projects on this node.</p>
							<a href="/me/requests" class="btn btn-primary btn-sm mt-2">
								Submit a hosting request
							</a>
						</div>
					{:else}
						<div class="flex flex-col gap-3">
							{#each myProjects as proj}
								<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-5 flex flex-col gap-3 hover:border-(--border-hairline) transition-colors">
									<div class="flex items-start justify-between gap-4">
										<div>
											<h3 class="font-display text-base font-bold text-(--text-main)">
												<a href="/projects/{proj.slug}" class="hover:text-(--accent-strong) transition-colors">{proj.name}</a>
											</h3>
											{#if proj.public_url}
												<a href={proj.public_url} target="_blank" rel="noreferrer" class="text-xs text-(--accent-strong) hover:underline inline-flex items-center gap-1 mt-0.5">
													<span>{proj.public_url}</span>
													<ArrowSquareOut size={12} weight="regular" />
												</a>
											{/if}
										</div>
										<StatusDot status={proj.status} />
									</div>

									{#if proj.description}
										<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">{proj.description}</p>
									{/if}

									{#if proj.technology_stack && proj.technology_stack.length > 0}
										<div class="text-[11px] text-(--text-muted) pt-2 border-t border-(--border-hairline)">
											{proj.technology_stack.join(' · ')}
										</div>
									{/if}
								</div>
							{/each}
						</div>
					{/if}
				</section>

				<!-- Pending Requests -->
				{#if myRequests.length > 0}
					<section class="flex flex-col gap-4">
						<div class="flex items-baseline justify-between">
							<h2 class="font-display text-lg font-bold text-(--text-main)">Submitted Requests</h2>
							<a href="/me/requests" class="text-xs text-(--accent-strong) hover:underline font-medium inline-flex items-center gap-1">
								<span>New request</span>
								<ArrowRight size={12} weight="bold" />
							</a>
						</div>

						<div class="flex flex-col gap-3">
							{#each myRequests as req}
								<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-5 flex flex-col gap-2">
									<div class="flex items-center justify-between gap-3">
										<strong class="text-sm font-semibold text-(--text-main)">{req.project_name}</strong>
										<span class="inline-flex items-center gap-1.5 text-xs px-2 py-0.5 rounded-md border {req.status === 'PENDING' ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20' : req.status === 'COMPLETED' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20' : 'bg-(--bg-muted) text-(--text-secondary) border-(--border-hairline)'}">
											{req.status}
										</span>
									</div>
									<p class="text-xs text-(--text-muted) truncate">{req.repository_url}</p>
									{#if req.admin_notes}
										<div class="mt-2 p-3 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs">
											<span class="font-semibold text-(--text-main) block text-xs">Administrator note:</span>
											<p class="text-(--text-secondary) mt-0.5">{req.admin_notes}</p>
										</div>
									{/if}
								</div>
							{/each}
						</div>
					</section>
				{/if}
			</div>

			<!-- Sidebar Column: Activity Notifications -->
			<aside class="flex flex-col gap-4">
				<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-5 flex flex-col gap-4">
					<div class="flex items-center justify-between pb-3 border-b border-(--border-hairline)">
						<h3 class="font-semibold text-sm text-(--text-main)">Notifications</h3>
						<span class="text-xs text-(--text-muted)">{notifications.length}</span>
					</div>

					{#if notifications.length === 0}
						<p class="text-xs text-(--text-muted) py-4 text-center">Inbox clear.</p>
					{:else}
						<ul class="flex flex-col gap-3 list-none">
							{#each notifications.slice(0, 5) as n}
								<li class="pb-3 border-b border-(--border-subtle) last:border-none last:pb-0">
									<strong class="block text-xs font-medium text-(--text-main)">{n.title}</strong>
									<p class="text-xs text-(--text-secondary) mt-0.5 leading-relaxed">{n.body}</p>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			</aside>
		</div>
	{:else if !loading}
		<div class="py-16 px-4 text-center bg-(--bg-surface) border border-(--border-hairline) rounded-md max-w-md mx-auto w-full flex flex-col gap-3">
			<h2 class="font-display text-xl font-bold text-(--text-main)">Sign In Required</h2>
			<p class="text-xs text-(--text-secondary)">Please authenticate to access your personal space.</p>
			<a href="/login" class="btn btn-primary btn-sm mx-auto mt-2">Sign in</a>
		</div>
	{/if}
</div>
