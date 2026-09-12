<script lang="ts">
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import {
		Bell,
		Check,
		CheckCircle,
		Checks
	} from 'phosphor-svelte';

	let notifications: any[] = [];
	let loading = true;
	let error: string | null = null;
	let markingAll = false;

	async function loadNotifications() {
		loading = true;
		error = null;
		try {
			const res = await api.get('/me/notifications');
			notifications = res.data?.notifications || [];
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadNotifications);

	async function markAsRead(id: string) {
		try {
			await api.patch(`/me/notifications/${id}/read`, {});
			notifications = notifications.map((n) =>
				n.id === id ? { ...n, is_read: true } : n
			);
		} catch (err) {
			console.error('Failed to mark notification as read:', err);
		}
	}

	async function markAllAsRead() {
		markingAll = true;
		try {
			await api.post('/me/notifications/read-all', {});
			notifications = notifications.map((n) => ({ ...n, is_read: true }));
		} catch (err) {
			console.error('Failed to mark all as read:', err);
		} finally {
			markingAll = false;
		}
	}

	function formatRelativeTime(dateStr: string): string {
		try {
			const d = new Date(dateStr);
			const diffSec = Math.floor((Date.now() - d.getTime()) / 1000);
			if (diffSec < 60) return 'Just now';
			const diffMin = Math.floor(diffSec / 60);
			if (diffMin < 60) return `${diffMin}m ago`;
			const diffHours = Math.floor(diffMin / 60);
			if (diffHours < 24) return `${diffHours}h ago`;
			const diffDays = Math.floor(diffHours / 24);
			if (diffDays < 30) return `${diffDays}d ago`;
			return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
		} catch {
			return dateStr;
		}
	}

	$: unreadCount = notifications.filter((n) => !n.is_read).length;
</script>

<svelte:head>
	<title>Notifications · Ngumpul Host</title>
</svelte:head>

<div class="p-6 sm:p-8 lg:p-10 max-w-4xl w-full flex flex-col gap-8">
	<!-- Page Header -->
	<header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-(--border-hairline)">
		<div class="flex flex-col gap-1">
			<div class="flex items-center gap-2.5">
				<h1 class="font-sans font-semibold text-2xl sm:text-3xl text-(--text-main) tracking-tight">
					Notifications
				</h1>
				{#if unreadCount > 0}
					<span class="px-2 py-0.5 rounded-full text-[11px] font-medium bg-(--accent-sky)/10 text-(--accent-sky) border border-(--accent-sky)/30">
						{unreadCount} unread
					</span>
				{/if}
			</div>
			<p class="text-xs text-(--text-secondary) leading-relaxed">
				Updates on hosting requests, project reviews, and system alerts.
			</p>
		</div>

		{#if unreadCount > 0}
			<button
				type="button"
				class="btn btn-secondary btn-sm text-xs px-3.5 py-1.5 inline-flex items-center gap-1.5 self-start sm:self-auto cursor-pointer"
				on:click={markAllAsRead}
				disabled={markingAll}
			>
				<Checks size={14} weight="bold" />
				<span>{markingAll ? 'Marking...' : 'Mark all as read'}</span>
			</button>
		{/if}
	</header>

	<!-- Notifications List -->
	{#if loading}
		<div class="py-24 text-center text-xs text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-6 h-6 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading notifications...</span>
		</div>
	{:else if error}
		<div class="p-4 rounded-md bg-rose-500/10 border border-rose-500/30 text-rose-700 dark:text-rose-400 text-xs">
			{error}
		</div>
	{:else if notifications.length === 0}
		<div class="p-12 text-center bg-(--bg-surface) border border-dashed border-(--border-hairline) rounded-md flex flex-col items-center gap-3 max-w-md mx-auto w-full">
			<div class="w-12 h-12 rounded-full bg-(--bg-muted) flex items-center justify-center text-(--text-muted)">
				<Bell size={24} />
			</div>
			<div class="flex flex-col gap-1">
				<h3 class="font-sans font-semibold text-sm text-(--text-main)">Inbox clear</h3>
				<p class="text-xs text-(--text-secondary) leading-relaxed">
					You don't have any notifications at this moment.
				</p>
			</div>
		</div>
	{:else}
		<div class="flex flex-col bg-(--bg-surface) border border-(--border-hairline) rounded-md divide-y divide-(--border-hairline) overflow-hidden shadow-xs">
			{#each notifications as notif (notif.id)}
				<div
					class="p-4 sm:p-5 flex items-start justify-between gap-4 transition-colors {notif.is_read ? 'opacity-85 hover:bg-(--bg-muted)/30' : 'bg-(--accent-sky)/5 hover:bg-(--accent-sky)/10 border-l-2 border-l-(--accent-sky)'}"
				>
					<div class="flex items-start gap-3.5 min-w-0">
						<div class="mt-0.5 shrink-0">
							{#if notif.is_read}
								<CheckCircle size={18} class="text-(--text-muted)" />
							{:else}
								<span class="w-2.5 h-2.5 rounded-full bg-(--accent-sky) block mt-1"></span>
							{/if}
						</div>

						<div class="flex flex-col gap-1 min-w-0">
							<h3 class="font-sans font-semibold text-xs sm:text-sm text-(--text-main) leading-snug">
								{notif.title}
							</h3>
							<p class="text-xs text-(--text-secondary) leading-relaxed">
								{notif.body}
							</p>
							<span class="text-[11px] text-(--text-muted) font-mono pt-1">
								{formatRelativeTime(notif.created_at)}
							</span>
						</div>
					</div>

					{#if !notif.is_read}
						<button
							type="button"
							on:click={() => markAsRead(notif.id)}
							class="text-xs text-(--text-muted) hover:text-(--text-main) p-1.5 rounded-sm hover:bg-(--bg-muted) transition-colors cursor-pointer shrink-0"
							title="Mark as read"
							aria-label="Mark as read"
						>
							<Check size={14} weight="bold" />
						</button>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
