<script lang="ts">
	export let title: string = 'Nothing here yet';
	export let description: string = '';
	let className: string = '';
	export { className as class };

	type $$Slots = {
		default: {};
		icon: {};
		action: {};
	};
</script>

<div
	class="w-full p-8 sm:p-12 text-center border border-dashed border-(--border-hairline) rounded-[8px] bg-(--bg-surface) flex flex-col items-center justify-center gap-3 {className}"
>
	<slot name="icon">
		<div class="w-10 h-10 rounded-[6px] border border-(--border-hairline) bg-(--bg-muted) flex items-center justify-center text-(--text-muted) mb-1">
			<span class="w-2 h-2 rounded-full bg-(--text-muted)/50"></span>
		</div>
	</slot>

	<div class="max-w-md flex flex-col gap-1">
		<h3 class="font-sans font-medium text-sm sm:text-base text-(--text-main)">
			{title}
		</h3>
		{#if description}
			<p class="text-xs sm:text-sm text-(--text-secondary) leading-relaxed">
				{description}
			</p>
		{/if}
	</div>

	<slot />
	<slot name="action" />
</div>
