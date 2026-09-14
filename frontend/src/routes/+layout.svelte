<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { initAuth } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import { Header, Footer } from '$lib/components/layout';
	import { setupApi } from '$lib/api/setup';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';

	import { nodeDomain } from '$lib/stores/node';
	import { Toaster } from 'svelte-sonner';

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 1000 * 60 * 2, // 2 minutes
				retry: 1
			}
		}
	});

	let setupChecked = false;
	let isInitialized = true;

	async function checkSetupState() {
		if (!browser) return;
		try {
			const res = await setupApi.getStatus();
			isInitialized = res.initialized;
			if (res.domain) {
				const clean = res.domain
					.replace(/^https?:\/\//, '')
					.replace(/:\d+$/, '')
					.trim();
				if (clean) {
					nodeDomain.set(clean);
				}
			}
			if (!res.initialized && $page.url.pathname !== '/setup') {
				goto('/setup');
			} else if (res.initialized && $page.url.pathname === '/setup') {
				goto('/');
			}
		} catch (_) {
			// In case of unexpected connection error, do not block
		} finally {
			setupChecked = true;
		}
	}

	onMount(() => {
		initAuth();
		checkSetupState();
	});

	$: if (browser && setupChecked) {
		if (!isInitialized && $page.url.pathname !== '/setup') {
			goto('/setup');
		} else if (isInitialized && $page.url.pathname === '/setup') {
			goto('/');
		}
	}

	$: isHome = $page.url.pathname === '/';
	$: isAdminRoute = $page.url.pathname.startsWith('/admin');
	$: isWorkspaceRoute = $page.url.pathname.startsWith('/me');
	$: isAuthOrSetupRoute = $page.url.pathname === '/login' || $page.url.pathname === '/setup' || $page.url.pathname === '/register';
</script>

<QueryClientProvider client={queryClient}>
	{#if !setupChecked && !isInitialized && $page.url.pathname !== '/setup'}
		<div class="min-h-screen flex items-center justify-center bg-(--bg-canvas)">
			<div class="w-5 h-5 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
		</div>
	{:else if isAdminRoute || isWorkspaceRoute || isAuthOrSetupRoute}
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

	<Toaster
		position="bottom-right"
		toastOptions={{
			classes: {
				toast: 'bg-(--bg-surface) text-(--text-main) border border-(--border-hairline) shadow-md rounded-md text-sm',
				description: 'text-(--text-muted) text-xs',
				actionButton: 'bg-(--accent-sky) text-white rounded-sm text-xs font-medium',
				cancelButton: 'bg-(--bg-muted) text-(--text-main) rounded-sm text-xs'
			}
		}}
	/>
</QueryClientProvider>
