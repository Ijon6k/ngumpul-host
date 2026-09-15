<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import { unreadNotificationsCount, refreshUnreadNotifications } from '$lib/stores/notifications';
	import { adminPendingCount, adminOpenReportsCount, refreshAdminStats } from '$lib/stores/adminStats';
	import ThemeToggle from './ThemeToggle.svelte';
	import { ArrowSquareOut, List } from 'phosphor-svelte';

	let scrollY = 0;
	let mobileNavOpen = false;

	$: isHome = $page.url.pathname === '/';
	$: isTransparent = isHome && scrollY < 60;

	// Close mobile nav on route change
	$: $page.url.pathname, (mobileNavOpen = false);

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
	<div class="container mx-auto px-4 sm:px-6 max-w-6xl flex items-center justify-between h-[60px] gap-4">

		<!-- Brand -->
		<div class="flex items-center gap-6 sm:gap-8">
			<a href="/" class="flex items-center gap-2 hover:opacity-85 transition-opacity shrink-0" title="Ngumpul Host">
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

		<!-- Desktop Right actions (hidden on mobile) -->
		<div class="hidden md:flex items-center gap-2">
			<ThemeToggle {isTransparent} />

			{#if $user}
				{#if $user.role === 'ADMIN'}
					<a
						href="/admin"
						class="btn btn-sm text-xs items-center gap-1.5 relative {isTransparent ? 'bg-white/10 text-white hover:bg-white/20 border-white/20' : 'btn-secondary'}"
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
					class="btn btn-sm text-xs {isTransparent ? 'bg-white text-neutral-950 hover:bg-neutral-100' : 'btn-primary'}"
				>Sign in</a>
			{/if}
		</div>

		<!-- Mobile Hamburger Button (visible < md) -->
		<div class="md:hidden flex items-center">
			<button
				type="button"
				class="p-1.5 rounded-sm border {isTransparent ? 'border-white/20 text-white hover:bg-white/10' : 'border-(--border-hairline) text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)'} transition-colors cursor-pointer"
				aria-label={mobileNavOpen ? 'Close navigation' : 'Open navigation'}
				aria-expanded={mobileNavOpen}
				on:click={() => (mobileNavOpen = !mobileNavOpen)}
			>
				<List size={18} weight="bold" />
			</button>
		</div>
	</div>

	<!-- Mobile Nav Drawer (slide-down, visible < md) -->
	{#if mobileNavOpen}
		<div class="md:hidden bg-(--bg-canvas)/98 backdrop-blur-xl border-t border-(--border-hairline) px-4 py-3 flex flex-col gap-0.5 shadow-lg">
			<a href="/projects" on:click={() => (mobileNavOpen = false)} class="px-3 py-2.5 text-sm text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) rounded-sm transition-colors">Projects</a>
			<a href="/people" on:click={() => (mobileNavOpen = false)} class="px-3 py-2.5 text-sm text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) rounded-sm transition-colors">People</a>
			<a href="/activity" on:click={() => (mobileNavOpen = false)} class="px-3 py-2.5 text-sm text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) rounded-sm transition-colors">Activity</a>
			<a href="/status" on:click={() => (mobileNavOpen = false)} class="px-3 py-2.5 text-sm text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) rounded-sm transition-colors">Status</a>

			<div class="h-px bg-(--border-hairline) my-1.5"></div>

			{#if $user}
				<a
					href="/me"
					on:click={() => (mobileNavOpen = false)}
					class="px-3 py-2.5 text-sm flex items-center justify-between text-(--text-main) hover:bg-(--bg-muted) rounded-sm transition-colors"
				>
					<div class="flex items-center gap-2.5">
						{#if $user.avatar_url}
							<img src={$user.avatar_url} alt={$user.display_name} class="w-5 h-5 rounded-full object-cover" />
						{:else}
							<span class="w-5 h-5 rounded-full bg-(--accent-soft) text-(--accent-strong) text-xs flex items-center justify-center font-bold">
								{($user.display_name || 'U').slice(0, 1)}
							</span>
						{/if}
						<span>Workspace</span>
					</div>
					{#if $unreadNotificationsCount > 0}
						<span class="px-2 py-0.5 rounded-full text-xs font-semibold bg-(--accent-orange) text-white">
							{$unreadNotificationsCount}
						</span>
					{/if}
				</a>

				{#if $user.role === 'ADMIN'}
					<a
						href="/admin"
						on:click={() => (mobileNavOpen = false)}
						class="px-3 py-2.5 text-sm flex items-center justify-between text-(--text-main) hover:bg-(--bg-muted) rounded-sm transition-colors"
					>
						<span>Operator Console</span>
						{#if ($adminPendingCount + $adminOpenReportsCount) > 0}
							<span class="w-2 h-2 rounded-full bg-rose-500 shrink-0"></span>
						{/if}
					</a>
				{/if}
			{:else}
				<a
					href="/login"
					on:click={() => (mobileNavOpen = false)}
					class="px-3 py-2.5 text-sm font-medium text-(--text-main) hover:bg-(--bg-muted) rounded-sm transition-colors flex items-center justify-between"
				>
					<span>Sign in</span>
					<span>→</span>
				</a>
			{/if}

			<div class="h-px bg-(--border-hairline) my-1.5"></div>

			<!-- Mobile Theme Switcher Row -->
			<div class="px-3 py-2 flex items-center justify-between text-xs text-(--text-secondary)">
				<span>Appearance</span>
				<ThemeToggle isTransparent={false} />
			</div>
		</div>
	{/if}
</header>
