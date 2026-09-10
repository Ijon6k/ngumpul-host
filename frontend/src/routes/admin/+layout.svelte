<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { user, logout } from '$lib/stores/auth';
	import { api } from '$lib/api';

	let pendingCount = 0;
	let theme = 'light';
	let mobileMenuOpen = false;

	onMount(async () => {
		const storedTheme = localStorage.getItem('theme');
		if (storedTheme) {
			theme = storedTheme;
			document.documentElement.setAttribute('data-theme', theme);
		} else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
			theme = 'dark';
			document.documentElement.setAttribute('data-theme', 'dark');
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

	$: currentPath = $page.url.pathname;

	const navItems = [
		{ label: 'Overview', href: '/admin', exact: true },
		{ label: 'Hosting Requests', href: '/admin/requests', badge: pendingCount },
		{ label: 'Projects', href: '/admin/projects' },
		{ label: 'Members', href: '/admin/users' },
	];

	const infraItems = [
		{ label: 'System Health', href: '/admin/system' },
		{ label: 'Audit Trail', href: '/admin/audit' },
	];

	function isActive(href: string, exact = false) {
		if (exact) return currentPath === href;
		return currentPath.startsWith(href);
	}

	function getPageTitle(path: string) {
		if (path === '/admin') return 'Instance Overview';
		if (path.startsWith('/admin/requests')) return 'Hosting Requests';
		if (path.startsWith('/admin/projects')) return 'Project Management';
		if (path.startsWith('/admin/users')) return 'Member Moderation';
		if (path.startsWith('/admin/system')) return 'System Diagnostics';
		if (path.startsWith('/admin/audit')) return 'Immutable Audit Log';
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
		     Linear-style Left Sidebar
		     ══════════════════════════════════════════════ -->
		<aside
			class="fixed lg:static top-0 bottom-0 left-0 z-50 w-[220px] border-r border-(--border-hairline) bg-(--bg-surface) flex flex-col transition-transform duration-200 {mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}"
		>
			<!-- Brand Header -->
			<div class="h-12 px-4 border-b border-(--border-hairline) flex items-center justify-between shrink-0">
				<a href="/admin" class="flex items-center gap-2 text-sm font-semibold tracking-tight hover:opacity-80 transition-opacity">
					<span class="flex items-end gap-[2px] h-[13px]">
						<span class="w-[3px] h-[10px] rounded-full bg-(--accent-sky)"></span>
						<span class="w-[3px] h-[13px] rounded-full bg-(--accent-sky)"></span>
						<span class="w-[3px] h-[8px] rounded-full bg-(--accent-sky)"></span>
					</span>
					<span class="font-display font-bold text-sm text-(--text-main)">
						ngumpul<span class="font-normal text-(--text-secondary)">ops</span>
					</span>
				</a>
				<!-- Live node chip -->
				<span class="inline-flex items-center gap-1 text-[10px] font-mono px-2 py-0.5 rounded-full bg-(--bg-muted) text-(--text-secondary) border border-(--border-hairline)">
					<span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-live-pulse inline-block"></span>
					NODE
				</span>
			</div>

			<!-- Navigation -->
			<div class="flex-1 overflow-y-auto px-3 py-4 flex flex-col gap-5">

				<!-- Management section -->
				<div>
					<div class="px-2 mb-1 text-[10px] font-semibold text-(--text-muted) uppercase tracking-widest">Management</div>
					<nav class="flex flex-col gap-0.5">
						{#each navItems as item}
							<a
								href={item.href}
								class="flex items-center justify-between px-2.5 py-1.5 rounded-lg text-[13px] transition-colors
									{isActive(item.href, item.exact)
										? 'bg-(--bg-muted) text-(--text-main) font-medium'
										: 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)'}"
								on:click={() => (mobileMenuOpen = false)}
							>
								<span>{item.label}</span>
								{#if item.badge && item.badge > 0}
									<span class="text-[10px] px-1.5 py-0.5 rounded-full bg-(--accent-orange) text-white font-semibold leading-none">
										{item.badge}
									</span>
								{/if}
							</a>
						{/each}
					</nav>
				</div>

				<!-- Diagnostics section -->
				<div>
					<div class="px-2 mb-1 text-[10px] font-semibold text-(--text-muted) uppercase tracking-widest">Diagnostics</div>
					<nav class="flex flex-col gap-0.5">
						{#each infraItems as item}
							<a
								href={item.href}
								class="flex items-center justify-between px-2.5 py-1.5 rounded-lg text-[13px] transition-colors
									{isActive(item.href)
										? 'bg-(--bg-muted) text-(--text-main) font-medium'
										: 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)'}"
								on:click={() => (mobileMenuOpen = false)}
							>
								<span>{item.label}</span>
							</a>
						{/each}
					</nav>
				</div>

				<!-- Shortcuts -->
				<div>
					<div class="px-2 mb-1 text-[10px] font-semibold text-(--text-muted) uppercase tracking-widest">Shortcuts</div>
					<nav class="flex flex-col gap-0.5">
						<a
							href="/"
							target="_blank"
							rel="noreferrer"
							class="flex items-center justify-between px-2.5 py-1.5 rounded-lg text-[13px] text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors"
						>
							<span>Public Showcase</span>
							<span class="text-xs text-(--text-muted)">↗</span>
						</a>
						<a
							href="/status"
							target="_blank"
							rel="noreferrer"
							class="flex items-center justify-between px-2.5 py-1.5 rounded-lg text-[13px] text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) transition-colors"
						>
							<span>System Status</span>
							<span class="text-xs text-(--text-muted)">↗</span>
						</a>
					</nav>
				</div>
			</div>

			<!-- User footer -->
			<div class="h-[56px] px-3 border-t border-(--border-hairline) flex items-center justify-between shrink-0">
				<div class="flex items-center gap-2.5 overflow-hidden min-w-0">
					<div class="w-7 h-7 rounded-lg bg-(--accent-soft) text-(--accent-strong) flex items-center justify-center text-[10px] font-bold shrink-0">
						{$user.username ? $user.username.slice(0, 2).toUpperCase() : 'AD'}
					</div>
					<div class="overflow-hidden min-w-0">
						<div class="text-[13px] font-medium text-(--text-main) truncate">{$user.display_name}</div>
						<div class="text-[10px] font-mono text-(--text-muted) truncate">@{$user.username}</div>
					</div>
				</div>
				<div class="flex items-center gap-1 shrink-0">
					<button
						type="button"
						class="text-[10px] font-mono px-2 py-1 rounded-md text-(--text-muted) hover:bg-(--bg-muted) hover:text-(--text-main) transition-colors cursor-pointer border border-transparent hover:border-(--border-hairline)"
						title="Toggle theme"
						on:click={toggleTheme}
					>
						{theme === 'light' ? '◑' : '○'}
					</button>
					<button
						type="button"
						class="text-[10px] font-mono px-2 py-1 rounded-md text-(--text-muted) hover:bg-(--bg-muted) hover:text-(--color-danger) transition-colors cursor-pointer"
						title="Sign out"
						on:click={logout}
					>
						Exit
					</button>
				</div>
			</div>
		</aside>

		<!-- ══════════════════════════════════════════════
		     Main Console Content
		     ══════════════════════════════════════════════ -->
		<div class="flex-1 flex flex-col min-w-0 min-h-screen">
			<!-- Top bar — Linear breadcrumb style -->
			<header class="h-12 border-b border-(--border-hairline) bg-(--bg-surface)/90 backdrop-blur-md px-6 flex items-center justify-between shrink-0 sticky top-0 z-30">
				<div class="flex items-center gap-3">
					<button
						type="button"
						class="lg:hidden p-1.5 -ml-1.5 text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted) rounded-md transition-colors cursor-pointer"
						on:click={() => (mobileMenuOpen = !mobileMenuOpen)}
						aria-label="Toggle navigation"
					>
						<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
						</svg>
					</button>

					<div class="flex items-center gap-2 text-[13px]">
						<span class="text-(--text-secondary)">Admin</span>
						<span class="text-(--text-muted)">/</span>
						<span class="text-(--text-main) font-medium">{getPageTitle(currentPath)}</span>
					</div>
				</div>

				<div class="hidden sm:flex items-center gap-3 text-[12px]">
					<div class="flex items-center gap-1.5 text-(--text-secondary)">
						<span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-live-pulse inline-block"></span>
						<span>Jakarta Node 01 · Healthy</span>
					</div>
				</div>
			</header>

			<!-- Page content -->
			<main class="flex-1 p-6 md:p-8 max-w-[1100px] w-full">
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
			<div class="w-2 h-2 rounded-full bg-(--accent-sky) animate-live-pulse mx-auto"></div>
			<h1 class="font-display font-bold text-lg text-(--text-main)">Authenticating...</h1>
			<p class="text-xs text-(--text-muted)">Checking session credentials...</p>
			<a href="/login" class="btn btn-secondary btn-sm mt-2">Sign in to console</a>
		</div>
	</div>
{/if}
