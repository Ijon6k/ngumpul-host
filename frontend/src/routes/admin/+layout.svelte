<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { user, logout, authLoading } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import {
		SquaresFour,
		Tray,
		Folder,
		Users,
		Pulse,
		Scroll,
		ArrowSquareOut,
		Lightning,
		SidebarSimple,
		SignOut,
		List,
		Sun,
		Moon,
		Gear,
		ShieldWarning,
		ChatDots
	} from 'phosphor-svelte';
	import {
		adminPendingCount,
		adminOpenReportsCount,
		refreshAdminStats,
		adminBreadcrumb
	} from '$lib/stores';
	import { Breadcrumb } from '$lib/components/ui';

	let theme = 'light';
	let mobileMenuOpen = false;
	let sidebarCollapsed = false;

	onMount(async () => {
		const storedTheme = localStorage.getItem('ngumpul_theme') || localStorage.getItem('theme');
		if (storedTheme) {
			theme = storedTheme;
			document.documentElement.setAttribute('data-theme', theme);
			document.documentElement.classList.toggle('dark', theme === 'dark');
		} else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
			theme = 'dark';
			document.documentElement.setAttribute('data-theme', 'dark');
			document.documentElement.classList.add('dark');
		}

		const storedCollapsed = localStorage.getItem('admin_sidebar_collapsed');
		if (storedCollapsed !== null) {
			sidebarCollapsed = storedCollapsed === 'true';
		}

		await refreshAdminStats();
	});

	let prevPath = '';
	$: if ($page.url.pathname !== prevPath) {
		prevPath = $page.url.pathname;
		adminBreadcrumb.set(null);
	}

	$: if (browser && !$authLoading && !$user) {
		goto('/login');
	}

	function toggleTheme() {
		theme = theme === 'light' ? 'dark' : 'light';
		localStorage.setItem('ngumpul_theme', theme);
		localStorage.setItem('theme', theme);
		document.documentElement.setAttribute('data-theme', theme);
		document.documentElement.classList.toggle('dark', theme === 'dark');
	}

	function toggleSidebar() {
		sidebarCollapsed = !sidebarCollapsed;
		localStorage.setItem('admin_sidebar_collapsed', String(sidebarCollapsed));
	}

	$: navItems = [
		{
			label: 'Overview',
			href: '/admin',
			exact: true,
			icon: SquaresFour
		},
		{
			label: 'Hosting Requests',
			href: '/admin/requests',
			exact: false,
			badge: $adminPendingCount,
			icon: Tray
		},
		{
			label: 'Reports',
			href: '/admin/reports',
			exact: false,
			badge: $adminOpenReportsCount,
			icon: ShieldWarning
		},
		{
			label: 'Comments',
			href: '/admin/comments',
			exact: false,
			icon: ChatDots
		},
		{
			label: 'Projects',
			href: '/admin/projects',
			exact: false,
			icon: Folder
		},
		{
			label: 'Members',
			href: '/admin/users',
			exact: false,
			icon: Users
		}
	];

	const infraItems = [
		{
			label: 'System Health',
			href: '/admin/system',
			exact: false,
			icon: Pulse
		},
		{
			label: 'Audit Trail',
			href: '/admin/audit',
			exact: false,
			icon: Scroll
		},
		{
			label: 'Settings',
			href: '/admin/settings',
			exact: false,
			icon: Gear
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
		if (p.startsWith('/admin/reports')) return 'Moderation Reports';
		if (p.startsWith('/admin/comments')) return 'Comments Moderation';
		if (p.startsWith('/admin/projects')) return 'Projects';
		if (p.startsWith('/admin/users')) return 'Members';
		if (p.startsWith('/admin/system')) return 'System Health';
		if (p.startsWith('/admin/audit')) return 'Audit Trail';
		if (p.startsWith('/admin/settings')) return 'Instance Settings';
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
					<a href="/admin" class="flex items-center gap-2 hover:opacity-85 transition-opacity min-w-0" title="ngumpul-host console">
						<img src="/logo.webp" alt="" class="h-5 w-5 object-contain shrink-0" />
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
						<SidebarSimple size={16} weight="regular" />
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
						<!-- Brand logo mark (Default) -->
						<span class="flex items-center justify-center transition-all duration-200 group-hover:opacity-0 group-hover:scale-75">
							<img src="/logo.webp" alt="Ngumpul Host" class="w-6 h-6 object-contain" />
						</span>
						<!-- Sidebar icon (Hover) -->
						<span class="absolute inset-0 flex items-center justify-center text-(--text-main) opacity-0 scale-75 group-hover:opacity-100 group-hover:scale-100 transition-all duration-200">
							<SidebarSimple size={18} weight="regular" />
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
									<svelte:component this={item.icon} size={16} weight="regular" class="shrink-0" />
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
									<svelte:component this={item.icon} size={16} weight="regular" class="shrink-0" />
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
								<ArrowSquareOut size={16} weight="regular" class="shrink-0" />
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
								<Lightning size={16} weight="regular" class="shrink-0" />
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
							class="p-1.5 rounded-md text-(--text-muted) hover:bg-(--bg-muted) hover:text-(--text-main) transition-colors cursor-pointer flex items-center justify-center"
							title="Toggle theme"
							on:click={toggleTheme}
						>
							{#if theme === 'light'}
								<Moon size={15} weight="regular" />
							{:else}
								<Sun size={15} weight="regular" />
							{/if}
						</button>
						<button
							type="button"
							class="px-2 py-1 text-xs rounded-md text-(--text-muted) hover:bg-(--bg-muted) hover:text-(--color-danger) transition-colors cursor-pointer font-medium flex items-center gap-1"
							title="Sign out"
							on:click={logout}
						>
							<SignOut size={13} weight="regular" />
							<span>Exit</span>
						</button>
					</div>
				{:else}
					<button
						type="button"
						class="w-full flex items-center justify-center p-1.5 rounded-md text-(--text-muted) hover:bg-(--bg-muted) hover:text-(--color-danger) transition-colors cursor-pointer"
						title="Sign out"
						on:click={logout}
					>
						<SignOut size={16} weight="regular" />
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
						<List size={20} weight="regular" />
					</button>

					<!-- Breadcrumb Navigation -->
					<Breadcrumb
						items={$adminBreadcrumb && $adminBreadcrumb.length > 0
							? [{ label: 'Console', href: '/admin' }, ...$adminBreadcrumb]
							: [
									{ label: 'Console', href: '/admin' },
									{ label: getPageTitle($page.url.pathname) }
							  ]}
					/>
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
		<div class="max-w-md w-full p-8 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col gap-4 text-center shadow-xs">
			<div class="w-10 h-10 rounded-full bg-(--color-danger)/10 text-(--color-danger) flex items-center justify-center mx-auto text-base font-bold">!</div>
			<h1 class="font-sans font-semibold text-xl text-(--text-main)">Access Restricted</h1>
			<p class="text-sm text-(--text-secondary) leading-relaxed">
				Your account <strong class="text-(--text-main)">@{$user.username}</strong> does not have administrator privileges to access the operational console.
			</p>
			<div class="flex justify-center gap-3 pt-2">
				<a href="/" class="btn btn-secondary btn-sm">Return Home</a>
				<a href="/me" class="btn btn-primary btn-sm">Your Workspace</a>
			</div>
		</div>
	</div>

{:else if $authLoading}
	<div class="min-h-screen bg-(--bg-canvas) flex items-center justify-center">
		<div class="w-5 h-5 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
	</div>
{/if}
