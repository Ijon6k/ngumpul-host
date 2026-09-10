<script lang="ts">
	export let activities: Array<{
		id: string;
		type: string;
		metadata: any;
		created_at: string;
		actor_name?: string;
		actor_username?: string;
		project_name?: string;
		project_slug?: string;
	}> = [];

	function formatTime(iso: string) {
		try {
			const d = new Date(iso);
			return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
		} catch (_) {
			return '';
		}
	}

	function formatDate(iso: string) {
		try {
			const d = new Date(iso);
			return d.toLocaleDateString([], { month: 'short', day: 'numeric' });
		} catch (_) {
			return '';
		}
	}
</script>

<div class="w-full">
	{#if activities.length === 0}
		<div class="p-8 text-center text-(--text-muted) border border-dashed border-(--border-hairline) rounded-md">
			<p class="text-sm">No activity events recorded yet.</p>
		</div>
	{:else}
		<div class="relative pl-5 border-l border-(--border-hairline) space-y-6">
			{#each activities as item}
				<div class="relative pb-1">
					<!-- Timeline node indicator -->
					<div class="absolute -left-[26px] top-1.5 w-2.5 h-2.5 rounded-full bg-(--accent-sky) ring-4 ring-(--bg-canvas)"></div>

					<div class="flex flex-col gap-1">
						<div class="flex items-baseline justify-between gap-4 flex-wrap">
							<span class="font-medium text-sm text-(--text-main)">
								{item.metadata?.title || 'Community event'}
							</span>
							<span class="text-xs text-(--text-muted) shrink-0">
								{formatDate(item.created_at)} · {formatTime(item.created_at)}
							</span>
						</div>

						{#if item.metadata?.message}
							<p class="text-xs text-(--text-secondary) leading-relaxed">{item.metadata.message}</p>
						{/if}

						{#if item.project_slug}
							<div class="mt-1">
								<a href="/projects/{item.project_slug}" class="text-xs font-medium text-(--accent-strong) hover:underline inline-flex items-center gap-1">
									<span>{item.project_name || item.project_slug}</span>
									<span>↗</span>
								</a>
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
