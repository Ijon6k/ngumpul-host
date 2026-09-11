<script lang="ts">
	export let label: string = '';
	export let value: string = '';
	export let subtext: string = '';
	export let dark: boolean = false;
</script>

<div class="flex flex-col p-4 sm:p-5 rounded-lg border {dark ? 'bg-[#0E1013] border-neutral-800 text-white' : 'bg-(--bg-surface) border-(--border-hairline)'}">
	<span class="text-xs {dark ? 'text-neutral-400' : 'text-(--text-muted)'} font-normal">
		{label}
	</span>

	<span class="font-sans font-normal text-xl sm:text-2xl mt-1 tracking-tight {dark ? 'text-white' : 'text-(--text-main)'}">
		{value}
	</span>

	{#if subtext}
		<span class="text-xs {dark ? 'text-neutral-500' : 'text-(--text-secondary)'} font-normal mt-1 leading-normal">
			{subtext}
		</span>
	{/if}
</div>
