<script lang="ts">
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
</script>

<!-- Mirrors the exact visual structure of ProjectCard.svelte during loading. -->
<article class="flex flex-col bg-(--bg-surface) border border-(--border-hairline) rounded-md overflow-hidden" aria-hidden="true">
	<!-- Cover area (16:9) -->
	<div class="w-full aspect-[16/9] bg-(--bg-muted)">
		<Skeleton variant="block" />
	</div>

	<!-- Card body -->
	<div class="flex flex-col grow p-5 gap-3">
		<div class="flex flex-col gap-2">
			<Skeleton variant="line" class="w-2/3 sm:w-1/2" />
			<Skeleton variant="line" class="w-full" />
			<Skeleton variant="line" class="w-11/12" />
		</div>

		<div class="flex items-center gap-2 mt-1">
			<Skeleton variant="line" class="w-20" />
			<Skeleton variant="line" class="w-28" />
		</div>

		<div class="flex items-center justify-between mt-auto pt-3 border-t border-(--border-hairline)">
			<Skeleton variant="line" class="w-28" />
			<Skeleton variant="line" class="w-12" />
		</div>
	</div>
</article>