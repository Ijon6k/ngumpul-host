<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi, extractError } from '$lib/api';
	import type { Report } from '$lib/types/report';
	import { refreshAdminStats } from '$lib/stores/adminStats';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';
	import {
		ShieldWarning,
		ArrowUpRight,
		Check,
		X,
		Trash,
		Funnel,
		ChatDots,
		Folder
	} from 'phosphor-svelte';

	let reports: Report[] = [];
	let loading = true;
	let error: string | null = null;
	let statusFilter = 'ALL';

	// Action states
	let resolvingReport: any = null;
	let resolutionStatus = 'RESOLVED';
	let resolutionNotes = '';
	let updating = false;

	// Delete comment modal state
	let deleteCommentModal = false;
	let commentIdToDelete: string | null = null;
	let deletingComment = false;

	async function loadReports() {
		loading = true;
		error = null;
		try {
			const res = await adminApi.reports.listReports(statusFilter);
			reports = res.reports || [];
			await refreshAdminStats();
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadReports);

	function setFilter(f: string) {
		statusFilter = f;
		loadReports();
	}

	function openResolutionDialog(report: any, newStatus: string) {
		resolvingReport = report;
		resolutionStatus = newStatus;
		resolutionNotes = '';
	}

	async function submitResolution() {
		if (!resolvingReport) return;
		updating = true;
		try {
			await adminApi.reports.patchReport(resolvingReport.id, {
				status: resolutionStatus,
				resolution_notes: resolutionNotes
			});
			resolvingReport = null;
			await loadReports();
			await refreshAdminStats();
		} catch (err) {
			alert('Failed to update report: ' + extractError(err));
		} finally {
			updating = false;
		}
	}

	function confirmDeleteComment(commentId: string) {
		commentIdToDelete = commentId;
		deleteCommentModal = true;
	}

	async function handleDeleteComment() {
		if (!commentIdToDelete) return;
		deletingComment = true;
		try {
			await adminApi.comments.deleteComment(commentIdToDelete);
			deleteCommentModal = false;
			commentIdToDelete = null;
			await loadReports();
		} catch (err) {
			alert('Failed to delete comment: ' + extractError(err));
		} finally {
			deletingComment = false;
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
</script>

<svelte:head>
	<title>Moderation Reports · Operator Console</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<!-- Top Bar -->
	<header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-(--border-hairline)">
		<div class="flex flex-col gap-1">
			<div class="flex items-center gap-2.5">
				<ShieldWarning size={22} weight="bold" class="text-amber-500" />
				<h1 class="font-sans font-semibold text-xl text-(--text-main)">Moderation Reports</h1>
			</div>
			<p class="text-xs text-(--text-secondary) leading-relaxed">
				Review and act on community-submitted spam, violation, or abuse reports.
			</p>
		</div>

		<!-- Status Filters -->
		<div class="flex items-center gap-1.5 p-1 rounded-md bg-(--bg-muted) border border-(--border-hairline) self-start sm:self-auto overflow-x-auto">
			{#each ['ALL', 'OPEN', 'REVIEWED', 'RESOLVED', 'DISMISSED'] as f}
				<button
					type="button"
					class="px-2.5 py-1 text-xs font-medium rounded-sm transition-colors cursor-pointer {statusFilter === f ? 'bg-(--bg-surface) text-(--text-main) shadow-2xs font-semibold' : 'text-(--text-secondary) hover:text-(--text-main)'}"
					on:click={() => setFilter(f)}
				>
					{f}
				</button>
			{/each}
		</div>
	</header>

	{#if loading}
		<div class="py-24 text-center text-xs text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-6 h-6 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading reports queue...</span>
		</div>
	{:else if error}
		<div class="p-4 rounded-md bg-rose-500/10 border border-rose-500/30 text-rose-700 dark:text-rose-400 text-xs">
			{error}
		</div>
	{:else if reports.length === 0}
		<div class="p-12 text-center bg-(--bg-surface) border border-dashed border-(--border-hairline) rounded-md flex flex-col items-center gap-2">
			<span class="text-xs font-medium text-(--text-main)">No reports in this view</span>
			<p class="text-xs text-(--text-muted)">The queue is clean. Community reports will appear here when submitted.</p>
		</div>
	{:else}
		<div class="flex flex-col bg-(--bg-surface) border border-(--border-hairline) rounded-md divide-y divide-(--border-hairline) overflow-hidden shadow-xs">
			{#each reports as r (r.id)}
				<div class="p-5 flex flex-col gap-3">
					<div class="flex items-start justify-between flex-wrap gap-3">
						<div class="flex items-center gap-2.5 flex-wrap">
							<!-- Target Badge -->
							<span class="inline-flex items-center gap-1 text-[11px] font-mono px-2 py-0.5 rounded-sm bg-(--bg-muted) text-(--text-secondary) border border-(--border-hairline)">
								{#if r.target_type === 'PROJECT'}
									<Folder size={12} />
									<span>PROJECT</span>
								{:else}
									<ChatDots size={12} />
									<span>COMMENT</span>
								{/if}
							</span>

							<!-- Reason Badge -->
							<span class="text-[11px] font-mono px-2 py-0.5 rounded-sm bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20 font-medium">
								{r.reason}
							</span>

							<!-- Status Badge -->
							<span class="text-[11px] font-mono px-2 py-0.5 rounded-sm border {r.status === 'OPEN' ? 'bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/30 font-semibold' : r.status === 'RESOLVED' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30' : r.status === 'REVIEWED' ? 'bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/30' : 'bg-(--bg-muted) text-(--text-muted) border-(--border-hairline)'}">
								{r.status}
							</span>
						</div>

						<span class="text-[11px] text-(--text-muted) font-mono">
							{formatDate(r.created_at)}
						</span>
					</div>

					<!-- Target Details -->
					<div class="p-3 bg-(--bg-muted)/50 border border-(--border-hairline) rounded-md text-xs flex flex-col gap-1.5">
						{#if r.target_type === 'PROJECT'}
							<div class="flex items-center justify-between">
								<span class="font-medium text-(--text-main)">Project: {r.target_name || r.target_id}</span>
								{#if r.target_slug}
									<a href="/projects/{r.target_slug}" target="_blank" rel="noreferrer" class="text-(--accent-sky) hover:underline inline-flex items-center gap-0.5 font-mono">
										<span>/{r.target_slug}</span>
										<ArrowUpRight size={11} />
									</a>
								{/if}
							</div>
						{:else}
							<div class="flex flex-col gap-1">
								<span class="text-[11px] text-(--text-muted)">Comment Content:</span>
								<p class="text-(--text-main) italic leading-relaxed">"{r.comment_content || 'Content not available or soft-deleted'}"</p>
							</div>
						{/if}
					</div>

					<!-- Report Details from Reporter -->
					{#if r.details}
						<div class="text-xs text-(--text-secondary) leading-relaxed">
							<span class="font-medium text-(--text-main)">Reporter Notes:</span> {r.details}
						</div>
					{/if}

					<!-- Reporter info & Resolution info -->
					<div class="flex items-center justify-between flex-wrap gap-3 pt-2 border-t border-(--border-hairline) text-[11px] text-(--text-muted)">
						<span>Reported by: <span class="text-(--text-secondary)">{r.reporter_name || 'Anonymous'}</span></span>

						{#if r.resolved_by_name}
							<span>Resolved by: <span class="text-(--text-secondary)">{r.resolved_by_name}</span> ({r.status})</span>
						{/if}

						<!-- Operator Actions -->
						<div class="flex items-center gap-1.5 ml-auto">
							{#if r.target_type === 'COMMENT' && r.comment_content}
								<button
									type="button"
									class="btn btn-secondary btn-sm text-[11px] px-2 py-1 text-rose-600 dark:text-rose-400 hover:bg-rose-500/10 inline-flex items-center gap-1"
									on:click={() => confirmDeleteComment(r.target_id)}
								>
									<Trash size={12} />
									<span>Delete Comment</span>
								</button>
							{/if}

							{#if r.status === 'OPEN'}
								<button
									type="button"
									class="btn btn-secondary btn-sm text-[11px] px-2 py-1"
									on:click={() => openResolutionDialog(r, 'REVIEWED')}
								>
									Mark Reviewed
								</button>
								<button
									type="button"
									class="btn btn-secondary btn-sm text-[11px] px-2 py-1"
									on:click={() => openResolutionDialog(r, 'DISMISSED')}
								>
									Dismiss
								</button>
								<button
									type="button"
									class="btn btn-primary btn-sm text-[11px] px-2.5 py-1 inline-flex items-center gap-1"
									on:click={() => openResolutionDialog(r, 'RESOLVED')}
								>
									<Check size={12} weight="bold" />
									<span>Resolve</span>
								</button>
							{:else if r.status === 'REVIEWED'}
								<button
									type="button"
									class="btn btn-secondary btn-sm text-[11px] px-2 py-1"
									on:click={() => openResolutionDialog(r, 'DISMISSED')}
								>
									Dismiss
								</button>
								<button
									type="button"
									class="btn btn-primary btn-sm text-[11px] px-2.5 py-1 inline-flex items-center gap-1"
									on:click={() => openResolutionDialog(r, 'RESOLVED')}
								>
									<Check size={12} weight="bold" />
									<span>Resolve</span>
								</button>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Resolution Modal -->
{#if resolvingReport}
	<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs"
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div class="w-full max-w-md bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-6 flex flex-col gap-4 shadow-xl">
			<div class="flex items-center justify-between">
				<h3 class="font-sans font-semibold text-base text-(--text-main)">
					{resolutionStatus === 'RESOLVED' ? 'Resolve Report' : resolutionStatus === 'DISMISSED' ? 'Dismiss Report' : 'Mark Report as Reviewed'}
				</h3>
				<button
					type="button"
					class="p-1 rounded-sm text-(--text-muted) hover:text-(--text-main)"
					on:click={() => (resolvingReport = null)}
				>
					<X size={16} />
				</button>
			</div>

			<p class="text-xs text-(--text-secondary)">
				Add resolution notes for the audit trail regarding this report.
			</p>

			<textarea
				bind:value={resolutionNotes}
				rows="3"
				placeholder="Resolution notes / rationale (optional)..."
				class="w-full text-xs p-2.5 rounded-md border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) placeholder:text-(--text-muted) focus:outline-hidden focus:border-(--accent-sky)"
			></textarea>

			<div class="flex items-center justify-end gap-2 pt-2 border-t border-(--border-hairline)">
				<button
					type="button"
					class="btn btn-secondary btn-sm text-xs px-3.5 py-1.5"
					on:click={() => (resolvingReport = null)}
					disabled={updating}
				>
					Cancel
				</button>
				<button
					type="button"
					class="btn btn-primary btn-sm text-xs px-4 py-1.5"
					on:click={submitResolution}
					disabled={updating}
				>
					{updating ? 'Saving...' : 'Confirm'}
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Confirm Delete Comment Modal -->
<ConfirmModal
	bind:open={deleteCommentModal}
	confirmLabel={deletingComment ? 'Deleting...' : 'Delete comment'}
	cancelLabel="Cancel"
	destructive={true}
	loading={deletingComment}
	on:confirm={handleDeleteComment}
>
	<svelte:fragment slot="title">Delete Reported Comment?</svelte:fragment>
	<p>
		This will permanently soft-delete the comment from the discussion. The text will be replaced with <em>"This comment was removed"</em>.
	</p>
</ConfirmModal>
