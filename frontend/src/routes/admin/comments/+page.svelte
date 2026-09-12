<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi, extractError } from '$lib/api';
	import type { Comment } from '$lib/types/comment';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';
	import {
		ChatDots,
		MagnifyingGlass,
		Trash,
		ArrowUpRight,
		ShieldWarning,
		Funnel
	} from 'phosphor-svelte';

	let comments: Comment[] = [];
	let loading = true;
	let error: string | null = null;
	let searchTerm = '';
	let searchDebounce: any = null;
	let filterFlaggedOnly = false;

	// Delete Modal State
	let deleteModalOpen = false;
	let commentToDelete: any = null;
	let deleting = false;

	async function loadComments() {
		loading = true;
		error = null;
		try {
			const query = searchTerm ? `search=${encodeURIComponent(searchTerm)}` : '';
			const res = await adminApi.comments.listComments(query);
			comments = res.comments || [];
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	function handleSearchInput() {
		clearTimeout(searchDebounce);
		searchDebounce = setTimeout(() => {
			loadComments();
		}, 300);
	}

	onMount(loadComments);

	function openDeleteModal(c: any) {
		commentToDelete = c;
		deleteModalOpen = true;
	}

	async function handleDelete() {
		if (!commentToDelete) return;
		deleting = true;
		try {
			await adminApi.comments.deleteComment(commentToDelete.id);
			deleteModalOpen = false;
			commentToDelete = null;
			await loadComments();
		} catch (err) {
			alert('Failed to delete comment: ' + extractError(err));
		} finally {
			deleting = false;
		}
	}

	function formatDate(iso?: string) {
		if (!iso) return '—';
		try {
			const d = new Date(iso);
			return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
		} catch {
			return iso;
		}
	}

	$: filteredComments = filterFlaggedOnly
		? comments.filter((c) => (c.report_count ?? 0) > 0)
		: comments;
</script>

<svelte:head>
	<title>Comments Moderation · Operator Console</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<!-- Page Header -->
	<header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-(--border-hairline)">
		<div class="flex flex-col gap-1">
			<div class="flex items-center gap-2.5">
				<ChatDots size={22} weight="bold" class="text-(--accent-sky)" />
				<h1 class="font-sans font-semibold text-xl text-(--text-main)">Comments Moderation</h1>
			</div>
			<p class="text-xs text-(--text-secondary) leading-relaxed">
				Directory of public community comments across all hosted project showcases.
			</p>
		</div>

		<!-- Search & Filter Controls -->
		<div class="flex items-center gap-2.5 flex-wrap">
			<div class="relative">
				<MagnifyingGlass size={14} class="absolute left-3 top-1/2 -translate-y-1/2 text-(--text-muted)" />
				<input
					type="search"
					bind:value={searchTerm}
					on:input={handleSearchInput}
					placeholder="Search content or author..."
					class="text-xs pl-8 pr-3 py-1.5 rounded-md border border-(--border-hairline) bg-(--bg-surface) text-(--text-main) placeholder:text-(--text-muted) focus:outline-hidden focus:border-(--accent-sky) w-56"
				/>
			</div>

			<button
				type="button"
				class="btn btn-secondary btn-sm text-xs px-3 py-1.5 inline-flex items-center gap-1.5 cursor-pointer {filterFlaggedOnly ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/30 font-medium' : ''}"
				on:click={() => (filterFlaggedOnly = !filterFlaggedOnly)}
			>
				<ShieldWarning size={14} />
				<span>Flagged Only</span>
			</button>
		</div>
	</header>

	{#if loading}
		<div class="py-24 text-center text-xs text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-6 h-6 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading comments...</span>
		</div>
	{:else if error}
		<div class="p-4 rounded-md bg-rose-500/10 border border-rose-500/30 text-rose-700 dark:text-rose-400 text-xs">
			{error}
		</div>
	{:else if filteredComments.length === 0}
		<div class="p-12 text-center bg-(--bg-surface) border border-dashed border-(--border-hairline) rounded-md flex flex-col items-center gap-2">
			<span class="text-xs font-medium text-(--text-main)">No comments located</span>
			<p class="text-xs text-(--text-muted)">Try adjusting your search filter.</p>
		</div>
	{:else}
		<div class="overflow-x-auto bg-(--bg-surface) border border-(--border-hairline) rounded-md shadow-xs">
			<table class="w-full text-left text-xs border-collapse font-sans">
				<thead>
					<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/40 text-(--text-secondary) font-medium">
						<th class="py-3 px-5 font-semibold">Author</th>
						<th class="py-3 px-4 font-semibold">Project</th>
						<th class="py-3 px-4 font-semibold">Content</th>
						<th class="py-3 px-3 font-semibold">Status</th>
						<th class="py-3 px-4 font-semibold">Posted</th>
						<th class="py-3 px-5 font-semibold text-right">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-(--border-hairline)">
					{#each filteredComments as c (c.id)}
						<tr class="hover:bg-(--bg-muted)/20 transition-colors">
							<!-- Author -->
							<td class="py-3.5 px-5">
								<div class="flex items-center gap-2.5">
									{#if c.author?.avatar_url}
										<img src={c.author.avatar_url} alt={c.author.display_name} class="w-7 h-7 rounded-full object-cover border border-(--border-hairline)" />
									{:else}
										<div class="w-7 h-7 rounded-full bg-(--bg-muted) text-(--text-secondary) flex items-center justify-center text-xs font-semibold border border-(--border-hairline)">
											{(c.author?.display_name || 'U').charAt(0)}
										</div>
									{/if}
									<div class="flex flex-col min-w-0">
										<span class="font-medium text-(--text-main) truncate">{c.author?.display_name || 'Member'}</span>
										<span class="text-xs text-(--text-muted) font-mono">@{c.author?.username}</span>
									</div>
								</div>
							</td>

							<!-- Project -->
							<td class="py-3.5 px-4 font-medium text-(--text-main)">
								<a
									href="/projects/{c.project_slug}"
									target="_blank"
									rel="noreferrer"
									class="hover:text-(--accent-sky) inline-flex items-center gap-1 font-mono"
								>
									<span>/{c.project_slug}</span>
									<ArrowUpRight size={11} />
								</a>
							</td>

							<!-- Content -->
							<td class="py-3.5 px-4 max-w-sm">
								{#if c.deleted_at}
									<span class="text-neutral-400 italic text-xs">[Comment deleted]</span>
								{:else}
									<p class="line-clamp-2 text-(--text-secondary) leading-relaxed">
										{c.content}
									</p>
								{/if}
							</td>

							<!-- Status / Reports -->
							<td class="py-3.5 px-3 whitespace-nowrap">
								{#if (c.report_count ?? 0) > 0}
									<a
										href="/admin/reports"
										class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/30"
									>
										<ShieldWarning size={12} weight="bold" />
										<span>{c.report_count} reported</span>
									</a>
								{:else if c.deleted_at}
									<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-neutral-200 dark:bg-neutral-800 text-neutral-500">
										Deleted
									</span>
								{:else}
									<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20">
										Active
									</span>
								{/if}
							</td>

							<!-- Posted Date -->
							<td class="py-3.5 px-4 font-mono text-xs text-(--text-muted) whitespace-nowrap">
								{formatDate(c.created_at)}
							</td>

							<!-- Actions -->
							<td class="py-3.5 px-5 text-right whitespace-nowrap">
								{#if !c.deleted_at}
									<button
										type="button"
										class="p-1.5 rounded-sm border border-(--border-hairline) text-rose-600 dark:text-rose-400 hover:bg-rose-500/10 transition-colors cursor-pointer"
										title="Delete Comment"
										on:click={() => openDeleteModal(c)}
									>
										<Trash size={14} />
									</button>
								{:else}
									<span class="text-neutral-400 text-xs">—</span>
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<!-- Confirm Delete Modal -->
<ConfirmModal
	bind:open={deleteModalOpen}
	confirmLabel={deleting ? 'Deleting...' : 'Delete comment'}
	cancelLabel="Cancel"
	destructive={true}
	loading={deleting}
	on:confirm={handleDelete}
>
	<svelte:fragment slot="title">Delete Community Comment?</svelte:fragment>
	<p>
		This comment will be soft-deleted. The text will be replaced with <em>"This comment was removed"</em> and will no longer be visible to catalog visitors.
	</p>
</ConfirmModal>
