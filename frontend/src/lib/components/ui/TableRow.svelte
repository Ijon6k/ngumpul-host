<script lang="ts">
	let className: string = '';
	export { className as class };
</script>

<tr class="hover:bg-(--bg-muted)/25 transition-colors {className}">
	<slot />
</tr>
