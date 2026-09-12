<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import { unreadNotificationsCount, refreshUnreadNotifications } from '$lib/stores/notifications';
	import { adminPendingCount, adminOpenReportsCount, refreshAdminStats } from '$lib/stores/adminStats';
	import ThemeToggle from './ThemeToggle.svelte';
	import { ArrowSquareOut } from 'phosphor-svelte';

	let scrollY = 0;

	$: isHome = $page.url.pathname === '/';
	$: isTransparent = isHome && scrollY < 60;

	onMount(() => {
		if ($user) {
			refreshUnreadNotifications();
			if ($user.role === 'ADMIN') {
				refreshAdminStats();
			}
		}
	});

	$: if ($user) {
		refreshUnreadNotifications();
		if ($user.role === 'ADMIN') {
			refreshAdminStats();
		}
	}
</script>

<svelte:window bind:scrollY />

<header
	class="fixed top-0 left-0 right-0 z-40 transition-all duration-200 {isTransparent
		? 'bg-transparent border-b border-transparent text-white'
		: 'bg-(--bg-canvas)/90 backdrop-blur-xl border-b border-(--border-hairline) text-(--text-main)'}"
>
	<div class="container mx-auto px-6 max-w-6xl flex items-center justify-between h-[60px] gap-6">

		<!-- Brand -->
		<div class="flex items-center gap-8">
			<a href="/" class="flex items-center gap-2 hover:opacity-85 transition-opacity" title="Ngumpul Host">
				<img src="/logo.webp" alt="" class="h-6 w-6 object-contain" />
				<span class="font-display font-bold text-sm tracking-tight {isTransparent ? 'text-white' : 'text-(--text-main)'}">
					ngumpul<span class="font-normal {isTransparent ? 'text-neutral-300' : 'text-(--text-secondary)'}">host</span><span class="text-(--accent-orange)">.</span>
				</span>
			</a>

			<nav class="hidden md:flex items-center gap-6 text-xs font-medium" aria-label="Main Navigation">
				<a
					href="/projects"
					class="transition-colors {isTransparent ? 'text-neutral-300 hover:text-white' : 'text-(--text-secondary) hover:text-(--text-main)'}"
				>Projects</a>
				<a
					href="/people"
					class="transition-colors {isTransparent ? 'text-neutral-300 hover:text-white' : 'text-(--text-secondary) hover:text-(--text-main)'}"
				>People</a>
				<a
					href="/activity"
					class="transition-colors {isTransparent ? 'text-neutral-300 hover:text-white' : 'text-(--text-secondary) hover:text-(--text-main)'}"
				>Activity</a>
				<a
					href="/status"
					class="transition-colors {isTransparent ? 'text-neutral-300 hover:text-white' : 'text-(--text-secondary) hover:text-(--text-main)'}"
				>Status</a>
			</nav>
		</div>

		<!-- Right actions -->
		<div class="flex items-center gap-2">
			<ThemeToggle {isTransparent} />

			{#if $user}
				{#if $user.role === 'ADMIN'}
					<a
						href="/admin"
						class="btn btn-sm text-xs inline-flex items-center gap-1.5 relative {isTransparent ? 'bg-white/10 text-white hover:bg-white/20 border-white/20' : 'btn-secondary'}"
					>
						<span>Console</span>
						{#if ($adminPendingCount + $adminOpenReportsCount) > 0}
							<span class="w-2 h-2 rounded-full bg-rose-500 shrink-0 animate-pulse" title="Pending items awaiting review"></span>
						{/if}
						<ArrowSquareOut size={14} weight="regular" class="opacity-70" />
					</a>
				{/if}
				<a
					href="/me"
					class="btn btn-sm flex items-center gap-2 text-xs relative {isTransparent ? 'bg-white text-neutral-900 hover:bg-neutral-100' : 'btn-secondary'}"
				>
					{#if $user.avatar_url}
						<img src={$user.avatar_url} alt={$user.display_name} class="w-4 h-4 rounded-full object-cover" />
					{:else}
						<span class="w-4 h-4 rounded-full bg-(--accent-soft) text-(--accent-strong) text-xs flex items-center justify-center font-bold">
							{($user.display_name || 'U').slice(0, 1)}
						</span>
					{/if}
					<span>Workspace</span>
					{#if $unreadNotificationsCount > 0}
						<span class="w-2 h-2 rounded-full bg-rose-500 shrink-0 animate-pulse" title="You have unread notifications"></span>
					{/if}
				</a>
			{:else}
				<a
					href="/login"
					class="text-xs font-medium px-2.5 py-1.5 transition-colors rounded-md {isTransparent ? 'text-neutral-200 hover:text-white hover:bg-white/10' : 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)'}"
				>Sign in</a>
				<a
					href="/register"
					class="btn btn-sm text-xs {isTransparent ? 'bg-white text-neutral-950 hover:bg-neutral-100' : 'btn-primary'}"
				>Request hosting</a>
			{/if}
		</div>
	</div>
</header>
