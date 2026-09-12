<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { CaretDown, Check } from 'phosphor-svelte';

	/**
	 * The label shown on the breadcrumb trigger (current item name).
	 */
	export let label: string;

	/**
	 * The id of the currently active item — used to highlight it in the list.
	 */
	export let activeId: string = '';

	/**
	 * List of options to show in the dropdown.
	 * Requires at minimum `id` and `label` fields.
	 */
	export let options: Array<{ id: string; label: string; [key: string]: unknown }> = [];

	/**
	 * Accessible label for the listbox (e.g. "Switch project").
	 */
	export let listboxLabel: string = 'Switch item';

	/**
	 * Header text shown at the top of the dropdown panel.
	 */
	export let dropdownHeader: string = listboxLabel;

	let open = false;

	const dispatch = createEventDispatcher<{ select: { id: string } }>();

	function select(id: string) {
		open = false;
		dispatch('select', { id });
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') open = false;
	}
</script>

{#if options.length > 0}
	<div class="relative">
		<!-- Trigger: always shows caret when there are options -->
		<button
			type="button"
			class="flex items-center gap-1.5 text-(--text-main) font-medium hover:text-(--accent-sky) transition-colors cursor-pointer group"
			on:click={() => (open = !open)}
			aria-haspopup="listbox"
			aria-expanded={open}
			aria-label={listboxLabel}
		>
			<span>{label}</span>
			<CaretDown
				size={13}
				weight="bold"
				class="opacity-50 group-hover:opacity-100 transition-all duration-150 {open ? 'rotate-180' : ''}"
			/>
		</button>

		{#if open}
			<!-- Backdrop: closes dropdown on outside click -->
			<!-- svelte-ignore a11y-no-static-element-interactions -->
			<div
				class="fixed inset-0 z-40"
				on:click={() => (open = false)}
				on:keydown={handleKeydown}
			></div>

			<!-- Dropdown panel -->
			<div
				class="absolute left-0 top-full mt-1.5 z-50 min-w-[200px] max-w-[280px] bg-(--bg-surface) border border-(--border-hairline) rounded-md shadow-lg overflow-hidden"
				role="listbox"
				aria-label={listboxLabel}
			>
				<div class="px-3 py-1.5 text-xs font-semibold text-(--text-muted) border-b border-(--border-hairline) tracking-wide">
					{dropdownHeader}
				</div>

				<ul class="py-1 max-h-52 overflow-y-auto">
					{#each options as opt (opt.id)}
						<li>
							<button
								type="button"
								role="option"
								aria-selected={opt.id === activeId}
								class="w-full text-left px-3 py-2 text-xs flex items-center justify-between gap-3 transition-colors cursor-pointer
									{opt.id === activeId
										? 'bg-(--cf-pastel-bg) text-(--cf-pastel-text) font-medium'
										: 'text-(--text-main) hover:bg-(--bg-muted)'}"
								on:click={() => select(opt.id)}
							>
								<span class="truncate">{opt.label}</span>
								{#if opt.id === activeId}
									<Check size={12} weight="bold" class="shrink-0" />
								{/if}
							</button>
						</li>
					{/each}
				</ul>
			</div>
		{/if}
	</div>
{:else}
	<!-- No options yet loaded — render plain text -->
	<span class="text-(--text-main) font-medium">{label}</span>
{/if}
