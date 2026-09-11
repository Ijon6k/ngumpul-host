<script lang="ts">
	import { onMount } from 'svelte';
	import { Moon, Sun } from 'phosphor-svelte';

	export let isTransparent: boolean = false;

	let theme = 'light';

	onMount(() => {
		const current = document.documentElement.getAttribute('data-theme') || (document.documentElement.classList.contains('dark') ? 'dark' : 'light');
		theme = current;
		if (theme === 'dark') {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	});

	function toggleTheme() {
		theme = theme === 'light' ? 'dark' : 'light';
		document.documentElement.setAttribute('data-theme', theme);
		document.documentElement.classList.toggle('dark', theme === 'dark');
		localStorage.setItem('ngumpul_theme', theme);
	}
</script>

<!-- Borderless theme toggle button with Phosphor Icons -->
<button
	type="button"
	class="inline-flex items-center justify-center w-8 h-8 rounded-[4px] border-0 transition-colors duration-200 cursor-pointer outline-none {isTransparent
		? 'text-white/80 hover:text-white hover:bg-white/10'
		: 'text-(--text-muted) hover:text-(--text-main) hover:bg-(--bg-hover)'}"
	on:click={toggleTheme}
	title={theme === 'light' ? 'Switch to Night mode' : 'Switch to Day mode'}
	aria-label="Toggle theme"
>
	{#if theme === 'light'}
		<Moon size={17} weight="regular" />
	{:else}
		<Sun size={17} weight="regular" />
	{/if}
</button>
