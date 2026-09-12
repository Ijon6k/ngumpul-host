<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';
	import {
		CaretLeft,
		Trash,
		ArrowUpRight,
		Folder,
		ShieldWarning
	} from 'phosphor-svelte';

	let projectId = $page.params.id;
	let project: any = null;
	let comments: any[] = [];
	let loading = true;
	let error: string | null = null;
	let success: string | null = null;

	// Delete Modal State
	let deleteModalOpen = false;
	let commentToDelete: any = null;
	let deleting = false;

	async function loadData() {
		loading = true;
		error = null;
		try {
			const projRes = await api.get('/admin/projects');
			const allProjects = projRes.data?.projects || [];
			project = allProjects.find((p: any) => p.id === projectId || p.slug === projectId);

			if (project) {
				const commRes = await api.get(`/projects/${project.slug}/comments`);
				comments = commRes.data?.comments || [];
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			loading = false;
		}
	}

	onMount(loadData);

	async function updateProjectStatus(newStatus: string) {
		if (!project) return;
		try {
			await api.patch(`/admin/projects/${project.id}`, { status: newStatus });
			project.status = newStatus;
			success = `Status updated to ${newStatus}.`;
		} catch (err) {
			error = extractError(err);
		}
	}

	async function toggleVisibility() {
		if (!project) return;
		const newVis = project.visibility === 'PUBLIC' ? 'UNPUBLISHED' : 'PUBLIC';
		try {
			await api.patch(`/admin/projects/${project.id}`, { visibility: newVis });
			project.visibility = newVis;
			success = `Visibility set to ${newVis}.`;
		} catch (err) {
			error = extractError(err);
		}
	}

	function openDeleteCommentModal(c: any) {
		commentToDelete = c;
		deleteModalOpen = true;
	}

	async function handleDeleteComment() {
		if (!commentToDelete) return;
		deleting = true;
		try {
			await api.delete(`/admin/comments/${commentToDelete.id}`);
			comments = comments.map((c) =>
				c.id === commentToDelete.id
					? { ...c, is_deleted: true, content: 'This comment was removed.' }
					: c
			);
			deleteModalOpen = false;
			commentToDelete = null;
			success = 'Comment removed successfully.';
		} catch (err) {
			error = extractError(err);
		} finally {
			deleting = false;
		}
	}

	function formatDate(iso?: string) {
		if (!iso) return '—';
		try {
			const d = new Date(iso);
			return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
		} catch {
			return iso;
		}
	}
</script>

<svelte:head>
	<title>{project?.name ? `${project.name} · Admin Moderation` : 'Admin Project Detail'}</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<!-- Breadcrumb -->
	<nav aria-label="Breadcrumb" class="flex items-center gap-2 text-xs text-(--text-muted)">
		<a href="/admin/projects" class="hover:text-(--text-main) transition-colors flex items-center gap-1">
			<CaretLeft size={12} weight="bold" />
			<span>Project Management</span>
		</a>
		<span class="opacity-40">/</span>
		<span class="text-(--text-main) font-medium">{project?.name || 'Details'}</span>
	</nav>

	{#if loading}
		<div class="py-20 text-center text-xs text-(--text-muted) flex flex-col items-center justify-center gap-3">
			<div class="w-5 h-5 border-2 border-(--text-muted) border-t-transparent rounded-full animate-spin"></div>
			<span>Loading project details...</span>
		</div>
	{:else if error && !project}
		<div class="p-4 rounded-sm bg-rose-500/10 text-rose-700 dark:text-rose-400 text-xs">
			{error}
		</div>
	{:else if project}
		{#if success}
			<div class="p-3 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 text-xs rounded-sm flex items-center justify-between">
				<span>{success}</span>
				<button type="button" on:click={() => (success = null)} class="font-bold cursor-pointer">×</button>
			</div>
		{/if}

		<!-- Project Overview Header -->
		<header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-(--border-hairline)">
			<div class="flex flex-col gap-1">
				<div class="flex items-center gap-3">
					<h1 class="font-sans font-semibold text-xl text-(--text-main)">
						{project.name}
					</h1>
					<StatusDot status={project.status} />
					<span class="text-xs font-mono text-(--text-muted)">/{project.slug}</span>
				</div>
				<p class="text-xs text-(--text-secondary) max-w-xl leading-relaxed">
					{project.description || 'No description provided.'}
				</p>
			</div>

			<div class="flex items-center gap-2 flex-wrap">
				<button
					type="button"
					class="btn btn-secondary btn-sm text-xs px-2.5 py-1.5 cursor-pointer"
					on:click={toggleVisibility}
				>
					Visibility: <span class="font-mono font-medium ml-1">{project.visibility}</span>
				</button>

				<select
					class="px-2.5 py-1.5 text-xs font-mono bg-(--bg-surface) border border-(--border-hairline) rounded-sm text-(--text-main) outline-none focus:border-(--accent-sky)"
					value={project.status}
					on:change={(e) => updateProjectStatus(e.currentTarget.value)}
				>
					<option value="ONLINE">ONLINE</option>
					<option value="SETUP">SETUP</option>
					<option value="OFFLINE">OFFLINE</option>
					<option value="ARCHIVED">ARCHIVED</option>
				</select>

				<a
					href="/projects/{project.slug}"
					target="_blank"
					rel="noreferrer"
					class="btn btn-secondary btn-sm text-xs px-2.5 py-1.5 inline-flex items-center gap-1 text-(--text-muted) hover:text-(--text-main)"
					title="Public Showcase"
				>
					<span>Showcase</span>
					<ArrowUpRight size={12} />
				</a>
			</div>
		</header>

		<!-- Contextual Comments Moderation Section -->
		<section class="flex flex-col gap-4">
			<div class="flex items-baseline justify-between gap-4">
				<h2 class="font-sans font-medium text-base text-(--text-main)">
					Comments ({comments.length})
				</h2>
				<span class="text-xs text-(--text-muted)">Contextual moderation for this project</span>
			</div>

			{#if comments.length === 0}
				<div class="py-8 text-center text-xs text-(--text-muted) border border-dashed border-(--border-hairline) rounded-sm">
					No comments posted on this project yet.
				</div>
			{:else}
				<div class="overflow-x-auto bg-(--bg-surface) border border-(--border-hairline) rounded-sm shadow-2xs">
					<table class="w-full text-left text-xs border-collapse font-sans">
						<thead>
							<tr class="border-b border-(--border-hairline) bg-(--bg-muted)/30 text-(--text-secondary) font-medium">
								<th class="py-2.5 px-4 font-semibold">Comment</th>
								<th class="py-2.5 px-3 font-semibold">Author</th>
								<th class="py-2.5 px-3 font-semibold">Created</th>
								<th class="py-2.5 px-3 font-semibold">Reports</th>
								<th class="py-2.5 px-4 font-semibold text-right">Action</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-(--border-hairline)">
							{#each comments as c (c.id)}
								<tr class="hover:bg-(--bg-muted)/20 transition-colors">
									<!-- Comment text -->
									<td class="py-3 px-4 max-w-sm">
										{#if c.is_deleted}
											<span class="text-neutral-400 italic text-[11px]">[Comment removed]</span>
										{:else}
											<p class="line-clamp-2 text-(--text-main) leading-relaxed">
												{c.content}
											</p>
										{/if}
									</td>

									<!-- Author -->
									<td class="py-3 px-3 whitespace-nowrap">
										<span class="font-medium text-(--text-main)">{c.author?.display_name || 'Member'}</span>
										{#if c.author?.username}
											<span class="block text-[10px] text-(--text-muted) font-mono">@{c.author.username}</span>
										{/if}
									</td>

									<!-- Created -->
									<td class="py-3 px-3 font-mono text-[11px] text-(--text-muted) whitespace-nowrap">
										{formatDate(c.created_at)}
									</td>

									<!-- Reports -->
									<td class="py-3 px-3 whitespace-nowrap">
										{#if c.report_count > 0}
											<span class="inline-flex items-center gap-1 text-[10px] font-medium text-amber-600 dark:text-amber-400">
												<ShieldWarning size={12} weight="bold" />
												<span>{c.report_count}</span>
											</span>
										{:else}
											<span class="text-neutral-400 text-[11px]">0</span>
										{/if}
									</td>

									<!-- Delete Action -->
									<td class="py-3 px-4 text-right whitespace-nowrap">
										{#if !c.is_deleted}
											<button
												type="button"
												class="p-1 rounded-sm border border-(--border-hairline) text-rose-600 dark:text-rose-400 hover:bg-rose-500/10 transition-colors cursor-pointer"
												title="Delete Comment"
												on:click={() => openDeleteCommentModal(c)}
											>
												<Trash size={13} />
											</button>
										{:else}
											<span class="text-neutral-400 text-[11px]">—</span>
										{/if}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</section>
	{/if}
</div>

<!-- Confirm Delete Comment Modal -->
<ConfirmModal
	bind:open={deleteModalOpen}
	confirmLabel={deleting ? 'Deleting...' : 'Delete comment'}
	cancelLabel="Cancel"
	destructive={true}
	loading={deleting}
	on:confirm={handleDeleteComment}
>
	<svelte:fragment slot="title">Delete Comment?</svelte:fragment>
	<p>
		The comment will be soft-deleted from this project's discussion. The text will be replaced with <em>"This comment was removed."</em>
	</p>
</ConfirmModal>
