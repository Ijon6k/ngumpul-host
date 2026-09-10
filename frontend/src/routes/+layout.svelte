<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { initAuth } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 1000 * 60 * 2, // 2 minutes
				retry: 1
			}
		}
	});

	onMount(() => {
		initAuth();
	});

	$: isAdminRoute = $page.url.pathname.startsWith('/admin');
</script>

<QueryClientProvider client={queryClient}>
	{#if isAdminRoute}
		<slot />
	{:else}
		<div class="min-h-screen flex flex-col bg-(--bg-canvas)">
			<Header />
			<main class="flex-1 shrink-0">
				<slot />
			</main>
			<Footer />
		</div>
	{/if}
</QueryClientProvider>
