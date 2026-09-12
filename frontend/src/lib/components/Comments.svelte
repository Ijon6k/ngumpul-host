<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/auth';
	import { api, extractError } from '$lib/api';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';
	import ReportModal from '$lib/components/ReportModal.svelte';
	import { DotsThree, Trash, Flag, PaperPlaneTilt } from 'phosphor-svelte';

	export let projectSlug: string;
	export let projectOwnerId: string;

	let comments: any[] = [];
	let loading = true;
	let error: string | null = null;
	let submitError: string | null = null;
	let newCommentText = '';
	let submitting = false;

	// Delete confirmation modal state
	let deleteModalOpen = false;
	let commentToDelete: any = null;
	let deleting = false;

	// Report modal state
	let reportModalOpen = false;
	let commentToReport: any = null;

	// Dropdown menu state
	let activeMenuId: string | null = null;

	function toggleMenu(id: string) {
		activeMenuId = activeMenuId === id ? null : id;
	}

	function closeMenu() {
		activeMenuId = null;
	}

	async function loadComments() {
		loading = true;
		error = null;
		try {
			const res = await api.get(`/projects/${projectSlug}/comments`);
			comments = res.data?.comments || [];
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadComments);

	async function handlePostComment() {
		if (!newCommentText.trim()) return;
		submitting = true;
		submitError = null;

		try {
			const res = await api.post(`/projects/${projectSlug}/comments`, {
				content: newCommentText.trim()
			});
			if (res.data?.comment) {
				comments = [...comments, res.data.comment];
				newCommentText = '';
			}
		} catch (err) {
			submitError = extractError(err);
		} finally {
			submitting = false;
		}
	}

	function openDeleteDialog(comment: any) {
		commentToDelete = comment;
		deleteModalOpen = true;
		closeMenu();
	}

	async function confirmDelete() {
		if (!commentToDelete) return;
		deleting = true;
		try {
			await api.delete(`/comments/${commentToDelete.id}`);
			comments = comments.map((c) => {
				if (c.id === commentToDelete.id) {
					return { ...c, is_deleted: true, content: 'This comment was removed.' };
				}
				return c;
			});
			deleteModalOpen = false;
			commentToDelete = null;
		} catch (err) {
			submitError = extractError(err);
		} finally {
			deleting = false;
		}
	}

	function openReportDialog(comment: any) {
		commentToReport = comment;
		reportModalOpen = true;
		closeMenu();
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
			return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
		} catch {
			return dateStr;
		}
	}

	$: canDelete = (authorId: string) => {
		if (!$user) return false;
		return $user.role === 'ADMIN' || $user.id === projectOwnerId;
	};
</script>

<svelte:window on:click={closeMenu} />

<section class="flex flex-col gap-6 pt-6 border-t border-(--border-hairline)">
	<!-- Section Title -->
	<div class="flex items-center justify-between">
		<h2 class="font-sans font-medium text-base text-(--text-main)">
			Comments
			{#if comments.length > 0}
				<span class="text-xs font-normal text-(--text-muted) ml-1">({comments.length})</span>
			{/if}
		</h2>
	</div>

	<!-- Comments List (Whitespace-first, no heavy bordered outer container) -->
	{#if loading}
		<div class="py-6 text-xs text-(--text-muted)">
			Loading comments...
		</div>
	{:else if error}
		<div class="p-3 rounded-sm bg-rose-500/10 text-rose-700 dark:text-rose-400 text-xs">
			{error}
		</div>
	{:else if comments.length === 0}
		<p class="text-xs text-(--text-muted) py-2">
			No comments yet. Leave a note below.
		</p>
	{:else}
		<div class="flex flex-col space-y-6">
			{#each comments as c (c.id)}
				<div class="flex items-start gap-3.5 group relative">
					<!-- Author Avatar -->
					{#if c.author?.avatar_url}
						<img
							src={c.author.avatar_url}
							alt={c.author.display_name}
							class="w-7 h-7 rounded-full object-cover shrink-0 border border-(--border-hairline) mt-0.5"
						/>
					{:else}
						<div class="w-7 h-7 rounded-full bg-(--bg-muted) border border-(--border-hairline) flex items-center justify-center text-[11px] font-medium text-(--text-main) shrink-0 mt-0.5">
							{c.author?.display_name?.charAt(0) || 'U'}
						</div>
					{/if}

					<!-- Main Body -->
					<div class="flex flex-col gap-1 flex-1 min-w-0">
						<div class="flex items-center justify-between gap-2">
							<div class="flex items-center gap-2">
								<span class="text-xs font-medium text-(--text-main)">
									{c.author?.display_name || 'Anonymous'}
								</span>
								{#if c.author_id === projectOwnerId}
									<span class="text-[10px] text-(--text-muted) font-mono">
										(author)
									</span>
								{/if}
							</div>

							<!-- Overflow Action Menu (Report / Delete) -->
							{#if !c.is_deleted && $user}
								<div class="relative shrink-0">
									<button
										type="button"
										class="p-1 rounded-sm text-(--text-muted) hover:text-(--text-main) transition-colors cursor-pointer"
										on:click|stopPropagation={() => toggleMenu(c.id)}
										aria-label="Comment actions"
									>
										<DotsThree size={16} weight="bold" />
									</button>

									{#if activeMenuId === c.id}
										<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
										<div
											class="absolute right-0 top-6 z-20 w-32 bg-(--bg-surface) border border-(--border-hairline) rounded-sm shadow-md py-1 flex flex-col text-xs"
											role="menu"
											tabindex="-1"
											on:click|stopPropagation
										>
											<button
												type="button"
												class="w-full px-3 py-1.5 text-left hover:bg-(--bg-muted) flex items-center gap-2 text-(--text-secondary) hover:text-(--text-main) transition-colors cursor-pointer"
												on:click={() => openReportDialog(c)}
											>
												<Flag size={13} />
												<span>Report</span>
											</button>

											{#if canDelete(c.author_id)}
												<button
													type="button"
													class="w-full px-3 py-1.5 text-left hover:bg-rose-500/10 flex items-center gap-2 text-rose-600 dark:text-rose-400 transition-colors border-t border-(--border-hairline) cursor-pointer"
													on:click={() => openDeleteDialog(c)}
												>
													<Trash size={13} />
													<span>Delete</span>
												</button>
											{/if}
										</div>
									{/if}
								</div>
							{/if}
						</div>

						{#if c.is_deleted}
							<p class="text-xs italic text-(--text-muted)">
								{c.content}
							</p>
						{:else}
							<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed break-words whitespace-pre-wrap">
								{c.content}
							</p>
						{/if}

						<span class="text-[11px] text-(--text-muted) pt-0.5">
							{formatRelativeTime(c.created_at)}
						</span>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	<!-- Comment Composer -->
	<div class="pt-4 mt-2">
		{#if $user}
			<form on:submit|preventDefault={handlePostComment} class="flex flex-col gap-3">
				<label for="new-comment" class="text-xs font-medium text-(--text-main)">
					Leave a comment
				</label>

				{#if submitError}
					<div class="p-2.5 rounded-sm bg-rose-500/10 text-rose-700 dark:text-rose-400 text-xs">
						{submitError}
					</div>
				{/if}

				<textarea
					id="new-comment"
					bind:value={newCommentText}
					rows="3"
					maxlength="1000"
					placeholder="Write a comment..."
					class="w-full text-xs sm:text-sm p-3 rounded-sm border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) placeholder:text-(--text-muted) focus:outline-hidden focus:border-(--accent-sky) resize-none leading-relaxed"
				></textarea>

				<div class="flex items-center justify-between">
					<span class="text-[11px] text-(--text-muted)">
						Posting as {$user.display_name}
					</span>

					<button
						type="submit"
						class="btn btn-primary btn-sm text-xs px-4 py-1.5 inline-flex items-center gap-1.5 cursor-pointer"
						disabled={submitting || !newCommentText.trim()}
					>
						{#if submitting}
							<span class="w-3 h-3 border-2 border-current border-t-transparent rounded-full animate-spin"></span>
						{:else}
							<PaperPlaneTilt size={13} weight="bold" />
						{/if}
						<span>Post comment</span>
					</button>
				</div>
			</form>
		{:else}
			<div class="py-4 text-xs text-(--text-muted) flex items-center gap-2">
				<span>Sign in with your community account to join the conversation.</span>
				<a href="/login" class="text-(--accent-sky) font-medium hover:underline">Sign in →</a>
			</div>
		{/if}
	</div>
</section>

<!-- Delete Confirmation Dialog -->
<ConfirmModal
	bind:open={deleteModalOpen}
	confirmLabel={deleting ? 'Deleting...' : 'Delete comment'}
	cancelLabel="Cancel"
	destructive={true}
	loading={deleting}
	on:confirm={confirmDelete}
>
	<svelte:fragment slot="title">Remove comment?</svelte:fragment>
	<p>
		Are you sure you want to remove this comment? The text will be replaced with <em>"This comment was removed."</em>
	</p>
</ConfirmModal>

<!-- Report Modal -->
<ReportModal
	bind:open={reportModalOpen}
	targetType="COMMENT"
	targetId={commentToReport?.id || ''}
	targetName={commentToReport?.author ? `Comment by ${commentToReport.author.display_name}` : 'Comment'}
/>
