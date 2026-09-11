<script lang="ts">
	let className: string = '';
	export { className as class };
</script>

<div class="w-full overflow-x-auto rounded-[8px] border border-(--border-hairline) bg-(--bg-surface) shadow-xs transition-colors {className}">
	<table class="w-full text-left text-xs border-collapse">
		<slot />
	</table>
</div>
