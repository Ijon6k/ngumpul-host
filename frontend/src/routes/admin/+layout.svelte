<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { user, logout } from '$lib/stores/auth';
	import { api } from '$lib/api';

	let pendingCount = 0;
	let theme = 'light';
	let mobileMenuOpen = false;
	let sidebarCollapsed = false;

	onMount(async () => {
		const storedTheme = localStorage.getItem('theme');
		if (storedTheme) {
			theme = storedTheme;
			document.documentElement.setAttribute('data-theme', theme);
		} else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
			theme = 'dark';
			document.documentElement.setAttribute('data-theme', 'dark');
		}

		const storedCollapsed = localStorage.getItem('admin_sidebar_collapsed');
		if (storedCollapsed !== null) {
			sidebarCollapsed = storedCollapsed === 'true';
		}

		try {
			const res = await api.get('/admin/stats');
			if (res.data?.pending_requests != null) {
				pendingCount = res.data.pending_requests;
			}
		} catch (e) {
			// Silently fail if not loaded
		}
	});

	function toggleTheme() {
		theme = theme === 'light' ? 'dark' : 'light';
		localStorage.setItem('theme', theme);
		document.documentElement.setAttribute('data-theme', theme);
	}

	function toggleSidebar() {
		sidebarCollapsed = !sidebarCollapsed;
		localStorage.setItem('admin_sidebar_collapsed', String(sidebarCollapsed));
	}

	const navItems = [
		{
			label: 'Overview',
			href: '/admin',
			exact: true,
			icon: `<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg>`
		},
		{
			label: 'Hosting Requests',
			href: '/admin/requests',
			exact: false,
			badge: pendingCount,
			icon: `<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>`
		},
		{
			label: 'Projects',
			href: '/admin/projects',
			exact: false,
			icon: `<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" /></svg>`
		},
		{
			label: 'Members',
			href: '/admin/users',
			exact: false,
			icon: `<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /></svg>`
		}
	];

	const infraItems = [
		{
			label: 'System Health',
			href: '/admin/system',
			exact: false,
			icon: `<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" /></svg>`
		},
		{
			label: 'Audit Trail',
			href: '/admin/audit',
			exact: false,
			icon: `<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" /></svg>`
		}
	];

	function checkActive(currentPathname: string, href: string, exact = false) {
		const path = (currentPathname || '').replace(/\/$/, '') || '/';
		const target = (href || '').replace(/\/$/, '') || '/';
		if (exact || target === '/admin') {
			return path === target;
		}
		return path.startsWith(target);
	}

	function getPageTitle(path: string) {
		const p = (path || '').replace(/\/$/, '') || '/';
		if (p === '/admin') return 'Overview';
		if (p.startsWith('/admin/requests')) return 'Hosting Requests';
		if (p.startsWith('/admin/projects')) return 'Projects';
		if (p.startsWith('/admin/users')) return 'Members';
		if (p.startsWith('/admin/system')) return 'System Health';
		if (p.startsWith('/admin/audit')) return 'Audit Trail';
		return 'Console';
	}
</script>

{#if $user && $user.role === 'ADMIN'}
	<div class="min-h-screen flex bg-(--bg-canvas) text-(--text-main) antialiased">
		<!-- Mobile Backdrop -->
		{#if mobileMenuOpen}
			<button
				type="button"
				class="fixed inset-0 z-40 bg-black/40 lg:hidden cursor-default"
				aria-label="Close menu"
				on:click={() => (mobileMenuOpen = false)}
			></button>
		{/if}

		<!-- ══════════════════════════════════════════════
		     Left Sidebar (Collapsible with Landing Page Logo)
		     ══════════════════════════════════════════════ -->
		<aside
			class="fixed lg:sticky top-0 bottom-0 left-0 z-50 h-screen border-r border-(--border-hairline) bg-(--bg-surface) flex flex-col transition-all duration-200 shrink-0
				{mobileMenuOpen ? 'translate-x-0 w-60' : '-translate-x-full lg:translate-x-0'}
				{sidebarCollapsed ? 'lg:w-[68px]' : 'lg:w-60'}"
		>
			<!-- Brand Header: Logo on Left, Sidebar Icon on Right (when expanded). When collapsed, logo morphs to sidebar icon on hover -->
			{#if !sidebarCollapsed}
				<div class="h-16 px-4 border-b border-(--border-hairline) flex items-center justify-between shrink-0">
					<!-- Brand Logo on Left -->
					<a href="/admin" class="flex items-center gap-2.5 hover:opacity-85 transition-opacity min-w-0" title="ngumpul-host console">
						<span class="flex items-end gap-[3px] h-[15px] shrink-0">
							<span class="w-[3.5px] h-[13px] rounded-full bg-(--accent-sky)"></span>
							<span class="w-[3.5px] h-[15px] rounded-full bg-(--accent-sky)"></span>
							<span class="w-[3.5px] h-[9px] rounded-full bg-(--accent-sky)"></span>
						</span>
						<span class="font-display font-medium text-sm tracking-tight text-(--text-main) truncate">
							ngumpul<span class="font-normal text-(--text-secondary)">host</span><span class="text-(--accent-orange)">.</span>
						</span>
					</a>

					<!-- Sidebar Toggle Icon on Right -->
					<button
						type="button"
						class="w-7 h-7 rounded-md flex items-center justify-center text-(--text-muted) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors cursor-pointer shrink-0"
						title="Minimize sidebar"
						aria-label="Minimize sidebar"
						on:click={toggleSidebar}
					>
						<svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<rect x="3" y="3" width="18" height="18" rx="2.5" />
							<path d="M9 3v18" />
							<path d="m16 15-3-3 3-3" />
						</svg>
					</button>
				</div>
			{:else}
				<!-- Collapsed State: Logo morphs to sidebar icon on hover -->
				<div class="h-16 px-3 border-b border-(--border-hairline) flex items-center justify-center shrink-0">
					<button
						type="button"
						class="group relative w-10 h-10 rounded-lg flex items-center justify-center hover:bg-(--bg-muted) transition-colors cursor-pointer"
						title="Expand sidebar"
						aria-label="Expand sidebar"
						on:click={toggleSidebar}
					>
						<!-- Brand 3-bar logo mark (Default) -->
						<span class="flex items-end gap-[3px] h-[15px] transition-all duration-200 group-hover:opacity-0 group-hover:scale-75">
							<span class="w-[3.5px] h-[13px] rounded-full bg-(--accent-sky)"></span>
							<span class="w-[3.5px] h-[15px] rounded-full bg-(--accent-sky)"></span>
							<span class="w-[3.5px] h-[9px] rounded-full bg-(--accent-sky)"></span>
						</span>
						<!-- Sidebar icon (Hover) -->
						<span class="absolute inset-0 flex items-center justify-center text-(--text-main) opacity-0 scale-75 group-hover:opacity-100 group-hover:scale-100 transition-all duration-200">
							<svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
								<rect x="3" y="3" width="18" height="18" rx="2.5" />
								<path d="M9 3v18" />
								<path d="m14 9 3 3-3 3" />
							</svg>
						</span>
					</button>
				</div>
			{/if}

			<!-- Navigation Items -->
			<div class="flex-1 overflow-y-auto px-2.5 py-4 flex flex-col gap-6">
				<!-- Management Section -->
				<div>
					{#if !sidebarCollapsed}
						<div class="px-2.5 mb-1.5 text-xs font-semibold text-(--text-muted)">Management</div>
					{/if}
					<nav class="flex flex-col gap-1">
						{#each navItems as item}
							{@const active = checkActive($page.url.pathname, item.href, item.exact)}
							<a
								href={item.href}
								class="flex items-center rounded-lg text-sm transition-colors relative group
									{sidebarCollapsed ? 'justify-center p-2.5' : 'justify-between px-3 py-2.5'}
									{active
										? 'bg-(--cf-pastel-bg) text-(--cf-pastel-text) font-semibold border border-(--cf-pastel-border)/50'
										: 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)'}"
								title={sidebarCollapsed ? item.label : undefined}
								on:click={() => (mobileMenuOpen = false)}
							>
								<div class="flex items-center gap-3">
									<span class="shrink-0">{@html item.icon}</span>
									{#if !sidebarCollapsed}
										<span>{item.label}</span>
									{/if}
								</div>

								{#if item.badge && item.badge > 0}
									{#if !sidebarCollapsed}
										<span class="text-xs px-2 py-0.5 rounded-full bg-(--accent-orange) text-white font-semibold leading-none">
											{item.badge}
										</span>
									{:else}
										<span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-(--accent-orange)"></span>
									{/if}
								{/if}
							</a>
						{/each}
					</nav>
				</div>

				<!-- Diagnostics Section -->
				<div>
					{#if !sidebarCollapsed}
						<div class="px-2.5 mb-1.5 text-xs font-semibold text-(--text-muted)">Diagnostics</div>
					{/if}
					<nav class="flex flex-col gap-1">
						{#each infraItems as item}
							{@const active = checkActive($page.url.pathname, item.href, item.exact)}
							<a
								href={item.href}
								class="flex items-center rounded-lg text-sm transition-colors relative group
									{sidebarCollapsed ? 'justify-center p-2.5' : 'justify-between px-3 py-2.5'}
									{active
										? 'bg-(--cf-pastel-bg) text-(--cf-pastel-text) font-semibold border border-(--cf-pastel-border)/50'
										: 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)'}"
								title={sidebarCollapsed ? item.label : undefined}
								on:click={() => (mobileMenuOpen = false)}
							>
								<div class="flex items-center gap-3">
									<span class="shrink-0">{@html item.icon}</span>
									{#if !sidebarCollapsed}
										<span>{item.label}</span>
									{/if}
								</div>
							</a>
						{/each}
					</nav>
				</div>

				<!-- External Links -->
				<div>
					{#if !sidebarCollapsed}
						<div class="px-2.5 mb-1.5 text-xs font-semibold text-(--text-muted)">External</div>
					{/if}
					<nav class="flex flex-col gap-1">
						<a
							href="/"
							target="_blank"
							rel="noreferrer"
							class="flex items-center rounded-lg text-sm text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors
								{sidebarCollapsed ? 'justify-center p-2.5' : 'justify-between px-3 py-2'}"
							title="Public Showcase"
						>
							<div class="flex items-center gap-3">
								<svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
									<path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
								</svg>
								{#if !sidebarCollapsed}
									<span>Public Site</span>
								{/if}
							</div>
							{#if !sidebarCollapsed}
								<span class="text-xs text-(--text-muted)">↗</span>
							{/if}
						</a>
						<a
							href="/status"
							target="_blank"
							rel="noreferrer"
							class="flex items-center rounded-lg text-sm text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors
								{sidebarCollapsed ? 'justify-center p-2.5' : 'justify-between px-3 py-2'}"
							title="Public Status Page"
						>
							<div class="flex items-center gap-3">
								<svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
									<path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
								</svg>
								{#if !sidebarCollapsed}
									<span>Status Page</span>
								{/if}
							</div>
							{#if !sidebarCollapsed}
								<span class="text-xs text-(--text-muted)">↗</span>
							{/if}
						</a>
					</nav>
				</div>
			</div>

			<!-- User Footer -->
			<div class="p-3 border-t border-(--border-hairline) flex items-center justify-between shrink-0">
				{#if !sidebarCollapsed}
					<div class="flex items-center gap-2.5 overflow-hidden min-w-0">
						<div class="w-8 h-8 rounded-lg bg-(--cf-pastel-bg) text-(--cf-pastel-text) flex items-center justify-center text-xs font-bold shrink-0">
							{$user.username ? $user.username.slice(0, 2).toUpperCase() : 'AD'}
						</div>
						<div class="overflow-hidden min-w-0">
							<div class="text-sm font-medium text-(--text-main) truncate">{$user.display_name}</div>
							<div class="text-xs text-(--text-muted) truncate">@{$user.username}</div>
						</div>
					</div>
					<div class="flex items-center gap-1 shrink-0">
						<button
							type="button"
							class="p-1.5 rounded-md text-(--text-muted) hover:bg-(--bg-muted) hover:text-(--text-main) transition-colors cursor-pointer"
							title="Toggle theme"
							on:click={toggleTheme}
						>
							{theme === 'light' ? '◑' : '○'}
						</button>
						<button
							type="button"
							class="px-2 py-1 text-xs rounded-md text-(--text-muted) hover:bg-(--bg-muted) hover:text-(--color-danger) transition-colors cursor-pointer font-medium"
							title="Sign out"
							on:click={logout}
						>
							Exit
						</button>
					</div>
				{:else}
					<button
						type="button"
						class="w-full flex items-center justify-center p-1.5 rounded-md text-(--text-muted) hover:bg-(--bg-muted) hover:text-(--text-main) transition-colors cursor-pointer"
						title="Sign out"
						on:click={logout}
					>
						<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
						</svg>
					</button>
				{/if}
			</div>
		</aside>

		<!-- ══════════════════════════════════════════════
		     Main Console Content (Full-Screen Layout)
		     ══════════════════════════════════════════════ -->
		<div class="flex-1 flex flex-col min-w-0 min-h-screen">
			<!-- Top Bar: Header with Logo, Sidebar Toggle Button & Breadcrumbs -->
			<header class="h-16 border-b border-(--border-hairline) bg-(--bg-surface) px-6 md:px-10 flex items-center justify-between shrink-0 sticky top-0 z-30">
				<div class="flex items-center gap-3">
					<!-- Mobile menu toggle (only on mobile when sidebar is off-canvas) -->
					<button
						type="button"
						class="lg:hidden p-1.5 -ml-1.5 rounded-md text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors cursor-pointer"
						title="Open navigation menu"
						aria-label="Open navigation menu"
						on:click={() => (mobileMenuOpen = !mobileMenuOpen)}
					>
						<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h7" />
						</svg>
					</button>

					<!-- Breadcrumb Navigation: Clean, no logo -->
					<nav class="flex items-center gap-2 text-sm" aria-label="Breadcrumb">
						<a href="/admin" class="text-(--text-muted) hover:text-(--text-main) transition-colors">Console</a>
						<span class="text-(--text-muted)">/</span>
						<span class="text-(--text-main) font-medium">{getPageTitle($page.url.pathname)}</span>
					</nav>
				</div>

				<div class="flex items-center gap-4 text-xs sm:text-sm">
					<div class="hidden sm:flex items-center gap-2 px-3 py-1 rounded-full bg-(--cf-green-pastel) text-(--cf-green-text) font-medium">
						<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
						<span>All services healthy</span>
					</div>
					<a href="/status" target="_blank" rel="noreferrer" class="text-(--text-secondary) hover:text-(--text-main) font-medium flex items-center gap-1">
						<span>Diagnostics</span>
						<span>↗</span>
					</a>
				</div>
			</header>

			<!-- Page Content: Expansive Full-Screen Canvas -->
			<main class="flex-1 p-6 md:p-10 w-full max-w-[1700px] mx-auto">
				<slot />
			</main>
		</div>
	</div>

{:else if $user && $user.role !== 'ADMIN'}
	<div class="min-h-screen flex items-center justify-center bg-(--bg-canvas) px-4">
		<div class="max-w-md w-full p-8 bg-(--bg-surface) border border-(--border-hairline) rounded-2xl flex flex-col gap-4 text-center">
			<div class="w-10 h-10 rounded-full bg-(--color-danger)/10 text-(--color-danger) flex items-center justify-center mx-auto text-base font-bold">!</div>
			<h1 class="font-display font-bold text-xl text-(--text-main)">Access Restricted</h1>
			<p class="text-sm text-(--text-secondary) leading-relaxed">
				Your account <strong class="text-(--text-main)">@{$user.username}</strong> does not have administrator privileges to access the operational console.
			</p>
			<div class="flex justify-center gap-3 pt-2">
				<a href="/" class="btn btn-secondary btn-sm">Return Home</a>
				<a href="/me" class="btn btn-primary btn-sm">Your Workspace</a>
			</div>
		</div>
	</div>

{:else}
	<div class="min-h-screen flex items-center justify-center bg-(--bg-canvas) px-4">
		<div class="max-w-md w-full p-8 bg-(--bg-surface) border border-(--border-hairline) rounded-2xl flex flex-col gap-4 text-center">
			<div class="w-2.5 h-2.5 rounded-full bg-(--accent-sky) animate-ping mx-auto"></div>
			<h1 class="font-display font-bold text-lg text-(--text-main)">Authenticating...</h1>
			<p class="text-xs text-(--text-muted)">Checking session credentials...</p>
			<a href="/login" class="btn btn-secondary btn-sm mt-2">Sign in to console</a>
		</div>
	</div>
{/if}
