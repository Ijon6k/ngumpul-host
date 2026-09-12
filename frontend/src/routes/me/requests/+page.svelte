<script lang="ts">
	import { onMount } from 'svelte';
	import { api, extractError } from '$lib/api';

	let requests: any[] = [];
	let loading = true;
	let submitting = false;
	let error: string | null = null;
	let successMsg: string | null = null;

	let projectName = '';
	let description = '';
	let repositoryUrl = '';
	let documentationUrl = '';
	let deploymentNotes = '';
	let techStackRaw = '';

	async function loadRequests() {
		try {
			const res = await api.get('/me/hosting-requests');
			requests = res.data?.requests || [];
		} catch (err) {
			console.error('Failed to load requests:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadRequests);

	async function handleSubmit() {
		error = null;
		successMsg = null;
		submitting = true;

		const techStack = techStackRaw
			.split(',')
			.map((t) => t.trim())
			.filter(Boolean);

		try {
			const res = await api.post('/hosting-requests', {
				project_name: projectName,
				description,
				repository_url: repositoryUrl,
				documentation_url: documentationUrl,
				deployment_notes: deploymentNotes,
				technology_stack: techStack
			});

			successMsg = res.data?.message || 'Hosting request submitted successfully.';
			projectName = '';
			description = '';
			repositoryUrl = '';
			documentationUrl = '';
			deploymentNotes = '';
			techStackRaw = '';
			await loadRequests();
		} catch (err) {
			error = extractError(err);
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:head>
	<title>Request Hosting · Ngumpul Host</title>
</svelte:head>

<div class="p-6 sm:p-8 lg:p-10 max-w-5xl w-full flex flex-col gap-8">
	<header class="flex flex-col gap-1 pb-6 border-b border-(--border-hairline)">
		<h1 class="font-sans font-semibold text-2xl sm:text-3xl text-(--text-main) tracking-tight">Request Project Hosting</h1>
		<p class="text-xs text-(--text-secondary) leading-relaxed">
			Submit your application repository for server allocation. Operators review specifications and provision reverse-proxy routing.
		</p>
	</header>

	<div class="grid grid-cols-1 lg:grid-cols-[1.3fr_1fr] gap-10">
		<!-- Form Column -->
		<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-6 sm:p-7 flex flex-col gap-5">
			<h2 class="font-display font-bold text-base text-(--text-main)">Application Specifications</h2>

			{#if error}
				<div class="p-3 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-xs">{error}</div>
			{/if}
			{#if successMsg}
				<div class="p-3 bg-(--accent-soft) text-(--accent-strong) border border-(--accent-sky)/30 rounded-md text-xs">{successMsg}</div>
			{/if}

			<form on:submit|preventDefault={handleSubmit} class="flex flex-col gap-4">
				<div class="flex flex-col gap-1.5">
					<label for="p-name" class="text-xs font-medium text-(--text-main)">Project Name *</label>
					<input
						id="p-name"
						type="text"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)"
						required
						placeholder="e.g. Nexus Dashboard"
						bind:value={projectName}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="p-desc" class="text-xs font-medium text-(--text-main)">Brief Purpose</label>
					<textarea
						id="p-desc"
						rows="2"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)"
						placeholder="What problem does this project solve for the community?"
						bind:value={description}
					></textarea>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="p-repo" class="text-xs font-medium text-(--text-main)">Git Repository URL *</label>
					<input
						id="p-repo"
						type="url"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)"
						required
						placeholder="https://github.com/username/project"
						bind:value={repositoryUrl}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="p-docs" class="text-xs font-medium text-(--text-main)">Documentation Link (Optional)</label>
					<input
						id="p-docs"
						type="url"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)"
						placeholder="https://..."
						bind:value={documentationUrl}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="p-tech" class="text-xs font-medium text-(--text-main)">Tech Stack (comma separated)</label>
					<input
						id="p-tech"
						type="text"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)"
						placeholder="SvelteKit, Go, PostgreSQL, Redis"
						bind:value={techStackRaw}
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label for="p-deploy" class="text-xs font-medium text-(--text-main)">Deployment Architecture Notes</label>
					<textarea
						id="p-deploy"
						rows="3"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--text-main)"
						placeholder="Port requirements, Dockerfile location, or environment dependencies..."
						bind:value={deploymentNotes}
					></textarea>
				</div>

				<button type="submit" class="btn btn-primary btn-sm mt-2 self-start" disabled={submitting}>
					{submitting ? 'Submitting to queue...' : 'Submit hosting request'}
				</button>
			</form>
		</div>

		<!-- History Column -->
		<div class="flex flex-col gap-4">
			<h2 class="font-display font-bold text-base text-(--text-main)">Submission History</h2>

			{#if loading}
				<p class="text-xs text-(--text-muted)">Loading requests...</p>
			{:else if requests.length === 0}
				<div class="p-8 text-center bg-(--bg-surface) border border-(--border-hairline) rounded-md text-xs text-(--text-secondary)">
					No hosting requests submitted yet.
				</div>
			{:else}
				<div class="flex flex-col gap-3">
					{#each requests as req (req.id)}
						<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-5 flex flex-col gap-2">
							<div class="flex items-center justify-between gap-2">
								<strong class="font-semibold text-sm text-(--text-main)">{req.project_name}</strong>
								<span class="inline-flex items-center gap-1.5 text-xs px-2 py-0.5 rounded-md border {req.status === 'PENDING' ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20' : req.status === 'COMPLETED' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20' : 'bg-(--bg-muted) text-(--text-secondary) border-(--border-hairline)'}">
									{req.status}
								</span>
							</div>
							<p class="text-xs text-(--text-muted) truncate">{req.repository_url}</p>
							{#if req.description}
								<p class="text-xs text-(--text-secondary) leading-relaxed">{req.description}</p>
							{/if}
							{#if req.admin_notes}
								<div class="mt-1 p-2.5 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs">
									<span class="font-semibold text-(--text-main) text-xs block">Operator response:</span>
									<p class="text-(--text-secondary) mt-0.5">{req.admin_notes}</p>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
