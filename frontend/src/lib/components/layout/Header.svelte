<script lang="ts">
	import { user } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import ThemeToggle from './ThemeToggle.svelte';
	import { ArrowSquareOut } from 'phosphor-svelte';

	let scrollY = 0;

	$: isHome = $page.url.pathname === '/';
	$: isTransparent = isHome && scrollY < 60;
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
			<a href="/" class="flex items-center gap-2.5 hover:opacity-85 transition-opacity">
				<!-- 3-bar logo mark -->
				<span class="flex items-end gap-[3px] h-[15px]">
					<span class="w-[4px] h-[13px] rounded-full {isTransparent ? 'bg-sky-400' : 'bg-(--accent-sky)'}"></span>
					<span class="w-[4px] h-[15px] rounded-full {isTransparent ? 'bg-sky-400' : 'bg-(--accent-sky)'}"></span>
					<span class="w-[4px] h-[9px] rounded-full {isTransparent ? 'bg-sky-400' : 'bg-(--accent-sky)'}"></span>
				</span>
				<span class="font-display font-bold text-sm tracking-tight {isTransparent ? 'text-white' : 'text-(--text-main)'}">
					ngumpul<span class="font-normal {isTransparent ? 'text-neutral-300' : 'text-(--text-secondary)'}">host</span><span class="text-(--accent-orange)">.</span>
				</span>
			</a>

			<nav class="hidden md:flex items-center gap-6 text-[13px] font-medium" aria-label="Main Navigation">
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
						class="btn btn-sm text-[13px] inline-flex items-center gap-1.5 {isTransparent ? 'bg-white/10 text-white hover:bg-white/20 border-white/20' : 'btn-secondary'}"
					>
						<span>Console</span>
						<ArrowSquareOut size={13} weight="regular" class="opacity-70" />
					</a>
				{/if}
				<a
					href="/me"
					class="btn btn-sm flex items-center gap-2 text-[13px] {isTransparent ? 'bg-white text-neutral-900 hover:bg-neutral-100' : 'btn-secondary'}"
				>
					{#if $user.avatar_url}
						<img src={$user.avatar_url} alt={$user.display_name} class="w-4 h-4 rounded-full object-cover" />
					{:else}
						<span class="w-4 h-4 rounded-full bg-(--accent-soft) text-(--accent-strong) text-[9px] flex items-center justify-center font-bold">
							{($user.display_name || 'U').slice(0, 1)}
						</span>
					{/if}
					Workspace
				</a>
			{:else}
				<a
					href="/login"
					class="text-[13px] font-medium px-2 py-1 transition-colors rounded-md {isTransparent ? 'text-neutral-200 hover:text-white hover:bg-white/10' : 'text-(--text-secondary) hover:text-(--text-main) hover:bg-(--bg-muted)'}"
				>Sign in</a>
				<a
					href="/register"
					class="btn btn-sm text-[13px] {isTransparent ? 'bg-white text-neutral-950 hover:bg-neutral-100' : 'btn-primary'}"
				>Request hosting</a>
			{/if}
		</div>
	</div>
</header>
