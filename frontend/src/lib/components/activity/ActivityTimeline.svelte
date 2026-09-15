<script lang="ts">
	import { Timeline, TimelineItem } from '$lib/components/ui';

	interface Props {
		activities?: Array<{
			id: string;
			type: string;
			metadata: any;
			created_at: string;
			actor_name?: string;
			actor_username?: string;
			project_name?: string;
			project_slug?: string;
		}>;
	}

	let { activities = [] }: Props = $props();

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
		<p class="text-xs text-(--text-muted) py-6 text-center">No activity events recorded yet.</p>
	{:else}
		<Timeline density="comfortable">
			{#each activities as item (item.id)}
				<TimelineItem
					title={item.metadata?.title || 'Community event'}
					timestamp="{formatDate(item.created_at)} · {formatTime(item.created_at)}"
					description={item.metadata?.message || ''}
					href={item.project_slug ? `/projects/${item.project_slug}` : ''}
					linkText={item.project_name || item.project_slug || ''}
					density="comfortable"
				/>
			{/each}
		</Timeline>
	{/if}
</div>
