<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { user, authLoading, logout } from '$lib/stores/auth';
	import { unreadNotificationsCount, refreshUnreadNotifications } from '$lib/stores/notifications';
	import {
		SquaresFour,
		Folder,
		Tray,
		ClockCounterClockwise,
		Bell,
		User,
		Gear,
		ShieldCheck,
		Sun,
		Moon,
		SignOut,
		List,
		X
	} from 'phosphor-svelte';

	let theme = 'light';
	let mobileMenuOpen = false;

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

		if ($user) {
			await refreshUnreadNotifications();
		}
	});

	function toggleTheme() {
		theme = theme === 'light' ? 'dark' : 'light';
		localStorage.setItem('ngumpul_theme', theme);
		localStorage.setItem('theme', theme);
		document.documentElement.setAttribute('data-theme', theme);
		document.documentElement.classList.toggle('dark', theme === 'dark');
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}

	$: mainNav = [
		{
			label: 'Overview',
			href: '/me',
			exact: true,
			icon: SquaresFour
		},
		{
			label: 'Projects',
			href: '/me/projects',
			exact: false,
			icon: Folder
		},
		{
			label: 'Hosting Requests',
			href: '/me/requests',
			exact: false,
			icon: Tray
		},
		{
			label: 'Activity',
			href: '/me/activity',
			exact: false,
			icon: ClockCounterClockwise
		},
		{
			label: 'Notifications',
			href: '/me/notifications',
			exact: false,
			icon: Bell,
			badge: $unreadNotificationsCount
		}
	];

	$: secondaryNav = [
		{
			label: 'Profile',
			href: '/me/profile',
			exact: false,
			icon: User
		},
		{
			label: 'Settings',
			href: '/me/settings',
			exact: false,
			icon: Gear
		}
	];

	function isItemActive(pathname: string, href: string, exact = false): boolean {
		const path = (pathname || '').replace(/\/$/, '') || '/';
		const target = (href || '').replace(/\/$/, '') || '/';
		if (exact) {
			return path === target;
		}
		return path === target || path.startsWith(target + '/');
	}
</script>

{#if $authLoading}
	<div class="min-h-screen bg-(--bg-canvas) flex flex-col items-center justify-center gap-3">
		<div class="w-5 h-5 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
		<span class="text-xs text-(--text-muted)">Loading workspace...</span>
	</div>
{:else if !$user}
	<div class="min-h-screen bg-(--bg-canvas) flex items-center justify-center p-4">
		<div class="max-w-md w-full p-8 rounded-md bg-(--bg-surface) border border-(--border-hairline) text-center flex flex-col items-center gap-4 shadow-sm">
			<div class="w-12 h-12 rounded-full bg-(--bg-muted) flex items-center justify-center text-(--text-muted)">
				<User size={22} weight="regular" />
			</div>
			<div class="flex flex-col gap-1">
				<h1 class="font-sans font-semibold text-lg text-(--text-main)">Authentication Required</h1>
				<p class="text-xs text-(--text-secondary) leading-relaxed">
					Please sign in to access your personal workspace, manage hosted projects, and track requests.
				</p>
			</div>
			<a href="/login" class="btn btn-primary text-xs px-5 py-2 mt-2">
				Sign In →
			</a>
		</div>
	</div>
{:else}
	<div class="min-h-screen bg-(--bg-canvas) text-(--text-main) flex flex-col">
		<!-- Mobile Header Bar (< lg) -->
		<header class="lg:hidden sticky top-0 z-40 bg-(--bg-surface)/95 backdrop-blur-md border-b border-(--border-hairline) px-4 py-3 flex items-center justify-between">
			<div class="flex items-center gap-3">
				<button
					type="button"
					class="p-1.5 rounded-sm border border-(--border-hairline) text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors cursor-pointer"
					on:click={() => (mobileMenuOpen = !mobileMenuOpen)}
					aria-label="Toggle navigation menu"
				>
					{#if mobileMenuOpen}
						<X size={18} weight="bold" />
					{:else}
						<List size={18} weight="bold" />
					{/if}
				</button>
				<a href="/" class="font-sans font-semibold text-sm tracking-tight text-(--text-main) flex items-center gap-1.5">
					<span class="w-2 h-2 rounded-full bg-(--accent-sky)"></span>
					<span>ngumpul<span class="text-(--text-muted)">host</span><span class="text-(--accent-sky)">.</span></span>
				</a>
			</div>

			<div class="flex items-center gap-2">
				<button
					type="button"
					on:click={toggleTheme}
					class="p-1.5 rounded-sm border border-(--border-hairline) text-(--text-muted) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors cursor-pointer"
					aria-label="Toggle theme"
				>
					{#if theme === 'dark'}
						<Sun size={15} weight="regular" />
					{:else}
						<Moon size={15} weight="regular" />
					{/if}
				</button>
				<a href="/me/profile" class="shrink-0">
					{#if $user.avatar_url}
						<img src={$user.avatar_url} alt={$user.display_name} class="w-7 h-7 rounded-full object-cover border border-(--border-hairline)" />
					{:else}
						<div class="w-7 h-7 rounded-full bg-(--bg-muted) text-(--text-secondary) flex items-center justify-center text-xs font-semibold border border-(--border-hairline)">
							{($user.display_name || 'U').charAt(0)}
						</div>
					{/if}
				</a>
			</div>
		</header>

		<!-- Mobile Navigation Drawer Overlay (< lg) -->
		{#if mobileMenuOpen}
			<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
			<div
				class="lg:hidden fixed inset-0 z-50 bg-black/60 backdrop-blur-xs flex"
				on:click={(e) => {
					if ((e.target as HTMLElement).dataset.backdrop === 'true') closeMobileMenu();
				}}
				data-backdrop="true"
				role="dialog"
				aria-modal="true"
				tabindex="-1"
			>
				<div class="w-72 bg-(--bg-surface) border-r border-(--border-hairline) h-full flex flex-col p-5 gap-6 shadow-2xl overflow-y-auto">
					<div class="flex items-center justify-between pb-3 border-b border-(--border-hairline)">
						<a href="/" class="font-sans font-semibold text-sm tracking-tight text-(--text-main) flex items-center gap-1.5" on:click={closeMobileMenu}>
							<span class="w-2 h-2 rounded-full bg-(--accent-sky)"></span>
							<span>ngumpul<span class="text-(--text-muted)">host</span><span class="text-(--accent-sky)">.</span></span>
						</a>
						<button
							type="button"
							on:click={closeMobileMenu}
							class="p-1 rounded-sm text-(--text-muted) hover:text-(--text-main)"
							aria-label="Close menu"
						>
							<X size={18} weight="bold" />
						</button>
					</div>

					<nav class="flex flex-col gap-1">
						<span class="px-3 py-1 text-xs font-medium text-(--text-muted)">Workspace</span>
						{#each mainNav as item}
							{@const active = isItemActive($page.url.pathname, item.href, item.exact)}
							<a
								href={item.href}
								on:click={closeMobileMenu}
								class="flex items-center justify-between px-3 py-2 rounded-sm text-sm transition-colors {active ? 'text-(--text-main) font-medium bg-(--bg-muted)' : 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)/50'}"
							>
								<span class="flex items-center gap-2.5">
									<svelte:component this={item.icon} size={18} weight={active ? 'bold' : 'regular'} />
									<span>{item.label}</span>
								</span>
								{#if item.badge && item.badge > 0}
									<span class="px-2 py-0.5 rounded-full text-xs font-mono font-medium bg-(--accent-sky) text-white">
										{item.badge}
									</span>
								{/if}
							</a>
						{/each}
					</nav>

					<div class="h-px bg-(--border-hairline)"></div>

					<nav class="flex flex-col gap-1">
						{#each secondaryNav as item}
							{@const active = isItemActive($page.url.pathname, item.href, item.exact)}
							<a
								href={item.href}
								on:click={closeMobileMenu}
								class="flex items-center gap-2.5 px-3 py-2 rounded-sm text-sm transition-colors {active ? 'text-(--text-main) font-medium bg-(--bg-muted)' : 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)/50'}"
							>
								<svelte:component this={item.icon} size={18} weight={active ? 'bold' : 'regular'} />
								<span>{item.label}</span>
							</a>
						{/each}

						{#if $user.role === 'ADMIN'}
							<a
								href="/admin"
								on:click={closeMobileMenu}
								class="flex items-center gap-2.5 px-3 py-2 rounded-sm text-sm font-medium text-amber-600 dark:text-amber-400 hover:bg-amber-500/10 transition-colors mt-1"
							>
								<ShieldCheck size={18} weight="regular" />
								<span>Admin Console</span>
							</a>
						{/if}
					</nav>

					<div class="mt-auto pt-4 border-t border-(--border-hairline) flex items-center justify-between">
						<button
							type="button"
							on:click={() => { closeMobileMenu(); logout(); }}
							class="flex items-center gap-2 text-xs text-rose-600 dark:text-rose-400 hover:opacity-80 transition-opacity cursor-pointer px-1 py-1"
						>
							<SignOut size={16} />
							<span>Sign out</span>
						</button>
					</div>
				</div>
			</div>
		{/if}

		<!-- Centered Desktop Workspace Canvas (lg+) -->
		<div class="max-w-6xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 lg:py-12 flex-1 flex items-start gap-10">
			<!-- Floating Navigation Rail (Desktop) -->
			<aside class="hidden lg:flex flex-col w-56 shrink-0 self-start sticky top-10 select-none gap-6">
				<!-- Top Brand Logo (Clicking returns to public site) -->
				<div class="flex items-center justify-between">
					<a
						href="/"
						class="font-sans font-semibold text-base tracking-tight text-(--text-main) flex items-center gap-2 group py-1"
						title="Return to Ngumpul Host Public Site"
					>
						<span class="w-2.5 h-2.5 rounded-full bg-(--accent-sky) group-hover:scale-110 transition-transform"></span>
						<span>ngumpul<span class="text-(--text-muted)">host</span><span class="text-(--accent-sky)">.</span></span>
					</a>

					<button
						type="button"
						on:click={toggleTheme}
						class="p-1.5 rounded-sm text-(--text-muted) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors cursor-pointer"
						aria-label="Toggle theme"
						title="Toggle Theme"
					>
						{#if theme === 'dark'}
							<Sun size={15} weight="regular" />
						{:else}
							<Moon size={15} weight="regular" />
						{/if}
					</button>
				</div>

				<!-- Workspace Primary Navigation -->
				<nav class="flex flex-col gap-1">
					<span class="px-3 py-1 text-xs font-medium text-(--text-muted)">
						Workspace
					</span>
					{#each mainNav as item}
						{@const active = isItemActive($page.url.pathname, item.href, item.exact)}
						<a
							href={item.href}
							class="flex items-center justify-between px-3 py-2 rounded-sm text-sm transition-colors {active ? 'text-(--text-main) font-medium bg-(--bg-muted)' : 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)/50'}"
						>
							<span class="flex items-center gap-2.5">
								<svelte:component this={item.icon} size={18} weight={active ? 'bold' : 'regular'} />
								<span>{item.label}</span>
							</span>
							{#if item.badge && item.badge > 0}
								<span class="px-2 py-0.5 rounded-full text-xs font-mono font-medium bg-(--accent-sky) text-white">
									{item.badge}
								</span>
							{/if}
						</a>
					{/each}
				</nav>

				<!-- Subtle Divider -->
				<div class="h-px bg-(--border-hairline)/60 mx-1"></div>

				<!-- Secondary Profile / Account Navigation -->
				<nav class="flex flex-col gap-1">
					{#each secondaryNav as item}
						{@const active = isItemActive($page.url.pathname, item.href, item.exact)}
						<a
							href={item.href}
							class="flex items-center gap-2.5 px-3 py-2 rounded-sm text-sm transition-colors {active ? 'text-(--text-main) font-medium bg-(--bg-muted)' : 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)/50'}"
						>
							<svelte:component this={item.icon} size={18} weight={active ? 'bold' : 'regular'} />
							<span>{item.label}</span>
						</a>
					{/each}

					{#if $user.role === 'ADMIN'}
						<a
							href="/admin"
							class="flex items-center gap-2.5 px-3 py-2 rounded-sm text-sm font-medium text-amber-600 dark:text-amber-400 hover:bg-amber-500/10 transition-colors mt-2"
						>
							<ShieldCheck size={18} weight="regular" />
							<span>Admin Console</span>
						</a>
					{/if}
				</nav>

				<!-- Rail Footer: User note & Sign Out -->
				<div class="mt-4 pt-4 border-t border-(--border-hairline)/60 flex items-center justify-between text-xs text-(--text-muted) px-2">
					<span class="truncate max-w-[120px] text-xs">@{$user.username}</span>
					<button
						type="button"
						on:click={logout}
						class="text-(--text-muted) hover:text-rose-500 p-1.5 rounded-sm transition-colors cursor-pointer"
						title="Sign Out"
						aria-label="Sign Out"
					>
						<SignOut size={16} />
					</button>
				</div>
			</aside>

			<!-- Centered Workspace Content -->
			<main class="flex-1 min-w-0">
				<slot />
			</main>
		</div>
	</div>
{/if}
