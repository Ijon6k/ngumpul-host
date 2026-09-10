<script lang="ts">
	import { onMount } from 'svelte';

	let theme = 'light';

	onMount(() => {
		const current = document.documentElement.getAttribute('data-theme') || 'light';
		theme = current;
	});

	function toggleTheme() {
		theme = theme === 'light' ? 'dark' : 'light';
		document.documentElement.setAttribute('data-theme', theme);
		localStorage.setItem('ngumpul_theme', theme);
	}
</script>

<button
	type="button"
	class="inline-flex items-center justify-center w-9 h-9 rounded-md border border-(--border) bg-(--surface) text-(--ink-secondary) hover:text-(--ink) hover:bg-(--surface-muted) hover:border-(--ink-muted) transition-all cursor-pointer"
	on:click={toggleTheme}
	title={theme === 'light' ? 'Switch to Night mode' : 'Switch to Day mode'}
	aria-label="Toggle theme"
>
	{#if theme === 'light'}
		<!-- Moon icon (Night) -->
		<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z" />
		</svg>
	{:else}
		<!-- Sun icon (Day) -->
		<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="12" cy="12" r="4" />
			<path d="M12 2v2" />
			<path d="M12 20v2" />
			<path d="m4.93 4.93 1.41 1.41" />
			<path d="m17.66 17.66 1.41 1.41" />
			<path d="M2 12h2" />
			<path d="M20 12h2" />
			<path d="m6.34 17.66-1.41 1.41" />
			<path d="m19.07 4.93-1.41 1.41" />
		</svg>
	{/if}
</button>
