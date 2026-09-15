<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { initAuth } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import { Header, Footer } from '$lib/components/layout';
	import { nodeDomain } from '$lib/stores/node';
	import { Toaster } from 'svelte-sonner';
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';

	interface Props {
		data: LayoutData;
		children?: Snippet;
	}

	let { data, children }: Props = $props();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 1000 * 60 * 2, // 2 minutes
				retry: 1
			}
		}
	});

	$effect(() => {
		if (data?.setupDomain) {
			const clean = data.setupDomain
				.replace(/^https?:\/\//, '')
				.replace(/:\d+$/, '')
				.trim();
			if (clean) {
				nodeDomain.set(clean);
			}
		}
	});

	onMount(() => {
		initAuth();
	});

	let isHome = $derived($page.url.pathname === '/');
	let isAdminRoute = $derived($page.url.pathname.startsWith('/admin'));
	let isWorkspaceRoute = $derived($page.url.pathname.startsWith('/me'));
	let isAuthOrSetupRoute = $derived(
		$page.url.pathname === '/login' ||
		$page.url.pathname === '/setup' ||
		$page.url.pathname === '/register'
	);
</script>

<QueryClientProvider client={queryClient}>
	{#if isAdminRoute || isWorkspaceRoute || isAuthOrSetupRoute}
		{@render children?.()}
	{:else}
		<div class="min-h-screen flex flex-col bg-(--bg-canvas)">
			<Header />
			<!-- Page content offset:
			     Home ('/') has 0 top padding so hero image flows seamlessly behind transparent navbar.
			     All other pages have pt-[60px] to start cleanly below the fixed 60px navbar without overlap. -->
			<main class="flex-1 shrink-0 {isHome ? '' : 'pt-[60px]'}">
				{@render children?.()}
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
