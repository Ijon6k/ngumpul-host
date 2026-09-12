<script lang="ts">
	export let value: string | number;
	export let label: string;
	export let subtext: string = '';
	let className: string = '';
	export { className as class };
</script>

<div class="flex flex-col gap-1 {className}">
	<span class="font-sans font-normal text-xl sm:text-2xl text-white tracking-tight">
		{value}
	</span>
	<span class="text-xs font-normal text-neutral-400">{label}</span>
	{#if subtext}
		<span class="text-xs font-normal text-neutral-500 truncate">
			{subtext}
		</span>
	{/if}
</div>
