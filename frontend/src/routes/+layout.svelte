<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { initAuth } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import { Header, Footer } from '$lib/components/layout';

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

	$: isHome = $page.url.pathname === '/';
	$: isAdminRoute = $page.url.pathname.startsWith('/admin');
	$: isWorkspaceRoute = $page.url.pathname.startsWith('/me');
</script>

<QueryClientProvider client={queryClient}>
	{#if isAdminRoute || isWorkspaceRoute}
		<slot />
	{:else}
		<div class="min-h-screen flex flex-col bg-(--bg-canvas)">
			<Header />
			<!-- Page content offset:
			     Home ('/') has 0 top padding so hero image flows seamlessly behind transparent navbar.
			     All other pages have pt-[60px] to start cleanly below the fixed 60px navbar without overlap. -->
			<main class="flex-1 shrink-0 {isHome ? '' : 'pt-[60px]'}">
				<slot />
			</main>
			<Footer />
		</div>
	{/if}
</QueryClientProvider>
