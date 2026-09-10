<script lang="ts">
	import { user } from '$lib/stores/auth';
	import ThemeToggle from './ThemeToggle.svelte';
</script>

<header class="sticky top-0 z-40 bg-(--bg-canvas)/85 backdrop-blur-xl border-b border-(--border-hairline)">
	<div class="container mx-auto px-6 max-w-6xl flex items-center justify-between h-[60px] gap-6">

		<!-- Brand -->
		<div class="flex items-center gap-8">
			<a href="/" class="flex items-center gap-2.5 hover:opacity-80 transition-opacity">
				<!-- 3-bar logo mark -->
				<span class="flex items-end gap-[3px] h-[15px]">
					<span class="w-[4px] h-[13px] rounded-full bg-(--accent-sky)"></span>
					<span class="w-[4px] h-[15px] rounded-full bg-(--accent-sky)"></span>
					<span class="w-[4px] h-[9px] rounded-full bg-(--accent-sky)"></span>
				</span>
				<span class="font-display font-bold text-sm tracking-tight text-(--text-main)">
					ngumpul<span class="font-normal text-(--text-secondary)">host</span><span class="text-(--accent-orange)">.</span>
				</span>
			</a>

			<nav class="hidden md:flex items-center gap-6 text-[13px] font-medium" aria-label="Main Navigation">
				<a href="/projects" class="text-(--text-secondary) hover:text-(--text-main) transition-colors">Projects</a>
				<a href="/people"   class="text-(--text-secondary) hover:text-(--text-main) transition-colors">People</a>
				<a href="/activity" class="text-(--text-secondary) hover:text-(--text-main) transition-colors">Activity</a>
				<a href="/status"   class="text-(--text-secondary) hover:text-(--text-main) transition-colors">Status</a>
			</nav>
		</div>

		<!-- Right actions -->
		<div class="flex items-center gap-2">
			<ThemeToggle />

			{#if $user}
				{#if $user.role === 'ADMIN'}
					<a href="/admin" class="btn btn-secondary btn-sm text-[13px]">
						Console <span class="text-(--text-muted) text-xs">↗</span>
					</a>
				{/if}
				<a href="/me" class="btn btn-secondary btn-sm flex items-center gap-2 text-[13px]">
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
				<a href="/login" class="text-[13px] font-medium text-(--text-secondary) hover:text-(--text-main) px-2 py-1 transition-colors rounded-md hover:bg-(--bg-muted)">Sign in</a>
				<a href="/register" class="btn btn-primary btn-sm text-[13px]">Request hosting</a>
			{/if}
		</div>
	</div>
</header>
