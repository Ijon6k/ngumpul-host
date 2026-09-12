<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { reportsApi, extractError } from '$lib/api';
	import { X, ShieldWarning } from 'phosphor-svelte';

	export let open: boolean = false;
	export let targetType: 'PROJECT' | 'COMMENT' = 'PROJECT';
	export let targetId: string = ''; // project slug or comment UUID
	export let targetName: string = ''; // e.g. "Atlas" or "comment by John"

	let reason = 'SPAM';
	let details = '';
	let loading = false;
	let error: string | null = null;
	let successMessage: string | null = null;

	const dispatch = createEventDispatcher<{ success: void; close: void }>();

	const reasons = [
		{ value: 'SPAM', label: 'Spam or automated bot content' },
		{ value: 'ABUSE_HARASSMENT', label: 'Harassment, hostility, or personal attack' },
		{ value: 'INAPPROPRIATE', label: 'Inappropriate or harmful content' },
		{ value: 'MALICIOUS_SUSPICIOUS', label: 'Malicious link, phishing, or malware' },
		{ value: 'OTHER', label: 'Other violation' }
	];

	function close() {
		if (loading) return;
		open = false;
		error = null;
		successMessage = null;
		details = '';
		reason = 'SPAM';
		dispatch('close');
	}

	async function handleSubmit() {
		loading = true;
		error = null;
		successMessage = null;

		try {
			if (targetType === 'PROJECT') {
				await reportsApi.reportProject(targetId, { reason, details });
			} else {
				await reportsApi.reportComment(targetId, { reason, details });
			}

			successMessage = 'Thanks. Your report has been submitted.';
			setTimeout(() => {
				close();
				dispatch('success');
			}, 1500);
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && open && !loading) close();
	}
</script>

<svelte:window on:keydown={handleKeydown} />

{#if open}
	<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs"
		on:click={(e) => {
			if ((e.target as HTMLElement).dataset.backdrop === 'true') close();
		}}
		data-backdrop="true"
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div class="w-full max-w-md bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-6 flex flex-col gap-4 shadow-xl">
			<!-- Header -->
			<div class="flex items-start justify-between gap-3">
				<div class="flex items-center gap-2.5 min-w-0">
					<ShieldWarning size={22} weight="bold" class="text-amber-500 shrink-0" />
					<h3 class="font-sans font-semibold text-base text-(--text-main)">
						Report {targetType === 'PROJECT' ? 'Project' : 'Comment'}
					</h3>
				</div>
				<button
					type="button"
					class="shrink-0 w-7 h-7 rounded-[4px] flex items-center justify-center text-(--text-muted) hover:text-(--text-main) hover:bg-(--bg-hover) transition-colors"
					on:click={close}
					disabled={loading}
					aria-label="Close dialog"
				>
					<X size={16} weight="bold" />
				</button>
			</div>

			{#if targetName}
				<p class="text-xs text-(--text-muted)">
					Reporting: <span class="font-medium text-(--text-secondary)">{targetName}</span>
				</p>
			{/if}

			{#if successMessage}
				<div class="p-4 rounded-md bg-emerald-500/10 border border-emerald-500/30 text-emerald-700 dark:text-emerald-400 text-xs sm:text-sm">
					{successMessage}
				</div>
			{:else}
				<form on:submit|preventDefault={handleSubmit} class="flex flex-col gap-4">
					{#if error}
						<div class="p-3 rounded-md bg-rose-500/10 border border-rose-500/30 text-rose-700 dark:text-rose-400 text-xs">
							{error}
						</div>
					{/if}

					<!-- Reasons -->
					<div class="flex flex-col gap-2">
						<label for="report-reason" class="text-xs font-medium text-(--text-secondary)">
							Why are you reporting this?
						</label>
						<div class="flex flex-col gap-1.5" id="report-reason">
							{#each reasons as r}
								<label class="flex items-start gap-2 text-xs text-(--text-secondary) cursor-pointer p-2 rounded-md hover:bg-(--bg-muted)/40 transition-colors border {reason === r.value ? 'border-(--accent-sky) bg-(--accent-sky)/5 text-(--text-main)' : 'border-transparent'}">
									<input
										type="radio"
										bind:group={reason}
										value={r.value}
										class="mt-0.5 accent-(--accent-sky)"
									/>
									<span>{r.label}</span>
								</label>
							{/each}
						</div>
					</div>

					<!-- Optional details -->
					<div class="flex flex-col gap-1.5">
						<label for="report-details" class="text-xs font-medium text-(--text-secondary)">
							Additional details <span class="text-(--text-muted) font-normal">(optional)</span>
						</label>
						<textarea
							id="report-details"
							bind:value={details}
							maxlength="500"
							rows="3"
							placeholder="Provide context for our moderators..."
							class="w-full text-xs p-2.5 rounded-md border border-(--border-hairline) bg-(--bg-canvas) text-(--text-main) placeholder:text-(--text-muted) focus:outline-hidden focus:border-(--accent-sky) resize-none"
						></textarea>
						<span class="text-xs text-(--text-muted) text-right">{details.length}/500</span>
					</div>

					<!-- Actions -->
					<div class="flex items-center justify-end gap-2.5 pt-2 border-t border-(--border-hairline)">
						<button
							type="button"
							class="btn btn-secondary btn-sm text-xs px-3.5 py-1.5"
							on:click={close}
							disabled={loading}
						>
							Cancel
						</button>
						<button
							type="submit"
							class="btn btn-primary btn-sm text-xs px-3.5 py-1.5 inline-flex items-center gap-1.5"
							disabled={loading}
						>
							{#if loading}
								<span class="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin"></span>
							{/if}
							<span>{loading ? 'Submitting...' : 'Submit Report'}</span>
						</button>
					</div>
				</form>
			{/if}
		</div>
	</div>
{/if}
