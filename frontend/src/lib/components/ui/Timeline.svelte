<script lang="ts">
	export let density: 'comfortable' | 'compact' = 'comfortable';
</script>

<div class="relative w-full {density === 'compact' ? 'space-y-4' : 'space-y-6'}">
	<slot />
</div>
