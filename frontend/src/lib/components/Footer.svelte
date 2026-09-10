<script lang="ts">
	import StatusDot from './StatusDot.svelte';
</script>

<footer class="mt-auto border-t border-(--border-hairline) py-12 px-6">
	<div class="container mx-auto max-w-6xl flex flex-col gap-6">
		<div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
			<div>
				<span class="font-display font-bold text-sm text-(--text-main)">Ngumpul Host</span>
				<p class="text-xs text-(--text-secondary) mt-0.5">
					A quiet community server for friends, experiments, and independent software.
				</p>
			</div>

			<div class="flex items-center gap-4 text-xs">
				<a href="/status" class="inline-flex items-center gap-2 text-(--text-secondary) hover:text-(--text-main) transition-colors">
					<StatusDot status="OPERATIONAL" showLabel={false} />
					<span>All systems operational</span>
				</a>
			</div>
		</div>

		<div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 pt-6 border-t border-(--border-subtle) text-xs text-(--text-muted)">
			<div class="flex items-center gap-4">
				<a href="/projects" class="hover:text-(--text-main) transition-colors">Projects</a>
				<span>/</span>
				<a href="/people" class="hover:text-(--text-main) transition-colors">Members</a>
				<span>/</span>
				<a href="/activity" class="hover:text-(--text-main) transition-colors">Activity</a>
				<span>/</span>
				<a href="/status" class="hover:text-(--text-main) transition-colors">Status</a>
			</div>

			<div>
				<span>Open source (MIT) · Self-hosted modular monolith</span>
			</div>
		</div>
	</div>
</footer>
