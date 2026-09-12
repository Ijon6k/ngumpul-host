<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { X, WarningCircle } from 'phosphor-svelte';

	/**
	 * Reusable destructive-action confirmation dialog.
	 *
	 * Design invariants:
	 *   - Backdrop: bg-black/60 backdrop-blur-xs (same as existing admin modals)
	 *   - Container: radius-lg (rounded-xl), max-w-md
	 *   - Confirm button: danger variant when destructive=true, primary otherwise
	 *   - Cancel: always secondary, keyboard Esc closes immediately
	 *
	 * Usage:
	 *   <ConfirmModal open={bool} on:confirm={handler} on:cancel={handler}>
	 *     <svelte:fragment slot="title">Revoke invitation?</svelte:fragment>
	 *     <p>This token will immediately be invalidated.</p>
	 *   </ConfirmModal>
	 */

	export let open: boolean = false;

	/** Label for the confirm action button. */
	export let confirmLabel: string = 'Confirm';

	/** Label for the cancel button. */
	export let cancelLabel: string = 'Cancel';

	/**
	 * When true: confirm button renders in danger styling (rose).
	 * When false: confirm button renders in primary styling (sky).
	 */
	export let destructive: boolean = true;

	/** While true, both buttons are disabled and confirm shows a spinner. */
	export let loading: boolean = false;

	/**
	 * When true, the confirm button is additionally disabled (e.g. for
	 * phrase-validation gates like PROMOTE). Cancel is still available.
	 */
	export let disabled: boolean = false;

	/** Optional inline error message shown above action row. */
	export let error: string = '';

	const dispatch = createEventDispatcher<{ confirm: void; cancel: void }>();

	function confirm() {
		if (!loading) dispatch('confirm');
	}

	function cancel() {
		if (!loading) {
			open = false;
			dispatch('cancel');
		}
	}

	function handleBackdrop(e: MouseEvent) {
		if ((e.target as HTMLElement).dataset.backdrop === 'true') cancel();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && open && !loading) cancel();
	}

	$: confirmClass = destructive
		? 'inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-[4px] border border-rose-300 dark:border-rose-900/50 bg-rose-600/10 text-rose-700 dark:text-rose-400 hover:bg-rose-600/20 active:scale-[0.99] font-medium text-sm transition-all disabled:opacity-50 disabled:pointer-events-none'
		: 'inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-[4px] border border-transparent bg-(--accent-sky) text-white hover:opacity-90 active:scale-[0.99] font-medium text-sm transition-all disabled:opacity-50 disabled:pointer-events-none';
</script>

<svelte:window on:keydown={handleKeydown} />

{#if open}
	<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs"
		data-backdrop="true"
		on:click={handleBackdrop}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div
			class="w-full max-w-md bg-(--bg-surface) border border-(--border-hairline) rounded-xl p-6 flex flex-col gap-4 shadow-xl"
		>
			<!-- Header row: title slot + close button -->
			<div class="flex items-start justify-between gap-3">
				<div class="flex items-center gap-2.5 min-w-0">
					{#if destructive}
						<WarningCircle
							size={22}
							weight="bold"
							class="text-rose-500 shrink-0 mt-px"
						/>
					{/if}
					<h3 class="font-semibold text-base text-(--text-main) leading-snug">
						<slot name="title">Confirm action</slot>
					</h3>
				</div>
				<button
					type="button"
					class="shrink-0 w-7 h-7 rounded-[4px] flex items-center justify-center text-(--text-muted) hover:text-(--text-main) hover:bg-(--bg-hover) transition-colors"
					on:click={cancel}
					disabled={loading}
					aria-label="Close dialog"
				>
					<X size={16} weight="bold" />
				</button>
			</div>

			<!-- Body: custom content goes here -->
			<div class="text-sm text-(--text-secondary) leading-relaxed">
				<slot />
			</div>

			<!-- Optional extra slot (e.g. input confirmation phrase) -->
			{#if $$slots.extra}
				<div>
					<slot name="extra" />
				</div>
			{/if}

			<!-- Error banner -->
			{#if error}
				<div
					class="px-3.5 py-2.5 rounded-[6px] bg-(--color-danger)/10 border border-(--color-danger)/30 text-(--color-danger) text-sm"
				>
					{error}
				</div>
			{/if}

			<!-- Action row -->
			<div class="flex items-center justify-end gap-2.5 pt-2 border-t border-(--border-hairline)">
				<button
					type="button"
					class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-[4px] border border-(--border-hairline) bg-(--bg-surface) text-(--text-main) hover:bg-(--bg-hover) active:scale-[0.99] font-normal text-sm transition-all disabled:opacity-50 disabled:pointer-events-none"
					on:click={cancel}
					disabled={loading}
				>
					{cancelLabel}
				</button>

				<button
					type="button"
					class={confirmClass}
					on:click={confirm}
					disabled={loading || disabled}
				>
					{#if loading}
						<span class="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin"></span>
					{/if}
					{loading ? 'Processing...' : confirmLabel}
				</button>
			</div>
		</div>
	</div>
{/if}
