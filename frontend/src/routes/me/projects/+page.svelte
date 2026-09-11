<script lang="ts">
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';
	import StatusDot from '$lib/components/StatusDot.svelte';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import { UploadSimple, ArrowSquareOut } from 'phosphor-svelte';

	let projects: any[] = [];
	let loading = true;
	let editingProject: any = null;
	let saving = false;
	let uploadingCover = false;
	let error: string | null = null;
	let success: string | null = null;

	let editDesc = '';
	let editCover = '';
	let editRepo = '';
	let editDocs = '';
	let editDemo = '';
	let editTechRaw = '';

	let fileInput: HTMLInputElement;

	async function handleCoverUpload(event: Event) {
		const target = event.target as HTMLInputElement;
		if (!target.files || target.files.length === 0) return;
		const file = target.files[0];

		if (file.size > 10 * 1024 * 1024) {
			error = 'File size exceeds maximum limit of 10 MB.';
			return;
		}

		uploadingCover = true;
		error = null;
		try {
			const formData = new FormData();
			formData.append('file', file);
			formData.append('purpose', 'cover');

			const res = await api.post('/upload', formData, {
				headers: { 'Content-Type': 'multipart/form-data' }
			});
			if (res.data?.url) {
				editCover = res.data.url;
			}
		} catch (err) {
			error = extractError(err);
		} finally {
			uploadingCover = false;
			if (fileInput) fileInput.value = '';
		}
	}

	async function loadProjects() {
		loading = true;
		try {
			const res = await api.get('/me/projects');
			projects = res.data?.projects || [];
		} catch (err) {
			console.error('Failed to load projects:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadProjects);

	function startEdit(proj: any) {
		editingProject = proj;
		editDesc = proj.description || '';
		editCover = proj.cover_image_url || '';
		editRepo = proj.repository_url || '';
		editDocs = proj.documentation_url || '';
		editDemo = proj.demo_url || '';
		editTechRaw = (proj.technology_stack || []).join(', ');
		error = null;
		success = null;
	}

	function cancelEdit() {
		editingProject = null;
	}

	async function saveEdit() {
		if (!editingProject) return;
		saving = true;
		error = null;
		success = null;

		const techStack = editTechRaw
			.split(',')
			.map((t) => t.trim())
			.filter(Boolean);

		try {
			await api.patch(`/me/projects/${editingProject.id}`, {
				description: editDesc,
				cover_image_url: editCover,
				repository_url: editRepo,
				documentation_url: editDocs,
				demo_url: editDemo,
				technology_stack: techStack
			});

			success = 'Project metadata updated successfully.';
			await loadProjects();
			editingProject = null;
		} catch (err) {
			error = extractError(err);
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Project Management · Ngumpul Host</title>
</svelte:head>

<div class="container mx-auto px-4 sm:px-6 max-w-4xl py-8 sm:py-12 pb-24 flex flex-col gap-6 sm:gap-8">
	<div class="flex items-center gap-2 text-xs text-(--text-muted)">
		<a href="/me" class="hover:text-(--text-main) transition-colors">Workspace</a>
		<span>/</span>
		<span class="text-(--text-main) font-medium">Manage Projects</span>
	</div>

	<header class="max-w-2xl">
		<h1 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-[-0.03em]">Project Management</h1>
		<p class="text-xs sm:text-sm text-(--text-secondary) mt-1.5 leading-relaxed font-normal">
			Update presentation-level details, repository links, documentation, and technology tags for your published workloads.
		</p>
	</header>

	{#if success}
		<div class="p-3 bg-(--accent-soft) text-(--accent-strong) border border-(--accent-sky)/30 rounded-md text-xs">{success}</div>
	{/if}
	{#if error}
		<div class="p-3 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-xs">{error}</div>
	{/if}

	{#if loading}
		<p class="py-16 text-center text-xs text-(--text-muted)">Loading your projects...</p>
	{:else if projects.length === 0}
		<div class="text-center py-14 px-6 bg-(--bg-surface) border border-(--border-hairline) rounded-md flex flex-col items-center gap-2">
			<p class="text-xs text-(--text-secondary)">You do not have any published projects to manage.</p>
			<a href="/me/requests" class="btn btn-primary btn-sm mt-2">
				Submit a hosting request
			</a>
		</div>
	{:else}
		<div class="flex flex-col gap-4">
			{#each projects as proj (proj.id)}
				<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-6 flex flex-col gap-4">
					<div class="flex items-center justify-between flex-wrap gap-4">
						<div>
							<h3 class="font-display font-bold text-base text-(--text-main)">{proj.name}</h3>
							{#if proj.public_url}
								<a href={proj.public_url} target="_blank" rel="noreferrer" class="text-xs text-(--accent-strong) hover:underline block mt-0.5">
									{proj.public_url} ↗
								</a>
							{/if}
						</div>
						<div class="flex items-center gap-3">
							<StatusDot status={proj.status} />
							<button type="button" class="btn btn-secondary btn-sm text-xs" on:click={() => startEdit(proj)}>
								Edit details
							</button>
						</div>
					</div>

					<!-- Inline Editor -->
					{#if editingProject && editingProject.id === proj.id}
						<div class="mt-2 pt-4 border-t border-(--border-hairline) flex flex-col gap-3">
							<h4 class="font-semibold text-xs text-(--accent-strong)">Editing Metadata: {proj.name}</h4>

							<div class="flex flex-col gap-1.5">
								<label for="e-desc" class="text-xs font-medium text-(--text-main)">Description</label>
								<textarea id="e-desc" class="w-full px-3 py-1.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)" rows="2" bind:value={editDesc}></textarea>
							</div>

							<div class="flex flex-col gap-2">
								<label for="e-cover" class="text-xs font-medium text-(--text-main)">Project Cover Image</label>
								<div class="flex flex-col sm:flex-row items-start gap-4">
									<div class="w-32 sm:w-40 aspect-4/3 rounded-md overflow-hidden border border-(--border-hairline) bg-(--bg-muted) shrink-0">
										<ProjectCover
											src={editCover}
											alt={proj.name}
											name={proj.name}
											aspectRatio="4/3"
										/>
									</div>
									<div class="flex-1 flex flex-col gap-2 w-full">
										<div class="flex items-center gap-2">
											<input
												type="file"
												accept="image/jpeg,image/png,image/webp"
												class="hidden"
												bind:this={fileInput}
												on:change={handleCoverUpload}
											/>
											<button
												type="button"
												class="btn btn-secondary btn-sm text-xs inline-flex items-center gap-1.5"
												on:click={() => fileInput.click()}
												disabled={uploadingCover}
											>
												<UploadSimple size={14} weight="bold" />
												<span>{uploadingCover ? 'Uploading...' : 'Upload Cover (Max 10 MB)'}</span>
											</button>
											{#if editCover}
												<button
													type="button"
													class="text-xs text-(--color-danger) hover:underline"
													on:click={() => (editCover = '')}
												>
													Remove
												</button>
											{/if}
										</div>
										<input
											id="e-cover"
											type="text"
											class="w-full px-3 py-1.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs font-mono text-(--text-main) outline-none focus:border-(--text-main)"
											placeholder="Image URL /uploads/covers/... or https://..."
											bind:value={editCover}
										/>
									</div>
								</div>
							</div>

							<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
								<div class="flex flex-col gap-1.5">
									<label for="e-repo" class="text-xs font-medium text-(--text-main)">Repository URL</label>
									<input id="e-repo" type="url" class="w-full px-3 py-1.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)" bind:value={editRepo} />
								</div>
								<div class="flex flex-col gap-1.5">
									<label for="e-docs" class="text-xs font-medium text-(--text-main)">Documentation URL</label>
									<input id="e-docs" type="url" class="w-full px-3 py-1.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)" bind:value={editDocs} />
								</div>
							</div>

							<div class="flex flex-col gap-1.5">
								<label for="e-tech" class="text-xs font-medium text-(--text-main)">Technology Stack (comma separated)</label>
								<input id="e-tech" type="text" class="w-full px-3 py-1.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)" bind:value={editTechRaw} />
							</div>

							<div class="flex items-center gap-2 mt-2">
								<button type="button" class="btn btn-primary btn-sm text-xs" on:click={saveEdit} disabled={saving}>
									{saving ? 'Saving...' : 'Save changes'}
								</button>
								<button type="button" class="btn btn-secondary btn-sm text-xs" on:click={cancelEdit}>
									Cancel
								</button>
							</div>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
