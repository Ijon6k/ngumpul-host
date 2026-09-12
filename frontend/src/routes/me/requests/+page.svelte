<script lang="ts">
	import { onMount } from 'svelte';
	import { hostingRequestsApi, extractError } from '$lib/api';
	import type { HostingRequest } from '$lib/types/hosting';
	import { detectLinkInfo } from '$lib/utils/linkDetector';
	import {
		UploadSimple,
		CheckCircle,
		XCircle,
		Spinner,
		Globe,
		BookOpen,
		Code
	} from 'phosphor-svelte';

	let requests: HostingRequest[] = [];
	let loading = true;
	let submitting = false;
	let error: string | null = null;
	let successMsg: string | null = null;

	let projectName = '';
	let subdomain = '';
	let description = '';
	let readme = '';
	let repositoryUrl = '';
	let documentationUrl = '';
	let deploymentNotes = '';
	let techStackRaw = '';

	// Subdomain availability check state
	let checkingSubdomain = false;
	let subdomainAvailable: boolean | null = null;
	let subdomainMessage = '';
	let subdomainDebounceTimer: any = null;

	const domainSuffix = '.ngumpul.local';

	async function loadRequests() {
		try {
			const res = await hostingRequestsApi.getMyHostingRequests();
			requests = res.requests || [];
		} catch (err) {
			console.error('Failed to load requests:', err);
		} finally {
			loading = false;
		}
	}

	onMount(loadRequests);

	function handleProjectNameInput() {
		if (!subdomain && projectName) {
			subdomain = projectName
				.toLowerCase()
				.replace(/[^a-z0-9]+/g, '-')
				.replace(/^-+|-+$/g, '');
			triggerSubdomainCheck();
		}
	}

	function handleSubdomainInput() {
		subdomain = subdomain
			.toLowerCase()
			.replace(/[^a-z0-9-]/g, '');
		triggerSubdomainCheck();
	}

	function triggerSubdomainCheck() {
		clearTimeout(subdomainDebounceTimer);
		const clean = subdomain.trim();
		if (!clean) {
			subdomainAvailable = null;
			subdomainMessage = '';
			checkingSubdomain = false;
			return;
		}

		if (clean.length < 2) {
			subdomainAvailable = false;
			subdomainMessage = 'Subdomain must be at least 2 characters';
			checkingSubdomain = false;
			return;
		}

		checkingSubdomain = true;
		subdomainDebounceTimer = setTimeout(async () => {
			try {
				const res = await hostingRequestsApi.checkSubdomain(clean);
				subdomainAvailable = res.available;
				subdomainMessage = res.message;
			} catch (err) {
				subdomainAvailable = false;
				subdomainMessage = 'Error checking subdomain availability';
			} finally {
				checkingSubdomain = false;
			}
		}, 300);
	}

	function handleMdUpload(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;

		const reader = new FileReader();
		reader.onload = (e) => {
			const text = e.target?.result as string;
			if (text) {
				readme = text;
			}
		};
		reader.readAsText(file);
	}

	async function handleSubmit() {
		error = null;
		successMsg = null;

		if (subdomainAvailable === false) {
			error = 'Please choose an available subdomain before submitting.';
			return;
		}

		submitting = true;

		const techStack = techStackRaw
			.split(',')
			.map((t) => t.trim())
			.filter(Boolean);

		try {
			await hostingRequestsApi.submitHostingRequest({
				project_name: projectName,
				subdomain: subdomain.trim(),
				description: description.trim(),
				readme: readme.trim(),
				repository_url: repositoryUrl.trim(),
				documentation_url: documentationUrl.trim(),
				deployment_notes: deploymentNotes.trim(),
				technology_stack: techStack
			});

			successMsg = 'Hosting request submitted successfully. Operators will review your allocation.';
			projectName = '';
			subdomain = '';
			subdomainAvailable = null;
			subdomainMessage = '';
			description = '';
			readme = '';
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

	$: repoInfo = repositoryUrl ? detectLinkInfo(repositoryUrl, 'Source code') : null;
	$: docInfo = documentationUrl ? detectLinkInfo(documentationUrl, 'Documentation') : null;
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

	<div class="grid grid-cols-1 lg:grid-cols-[1.3fr_1fr] gap-10 items-start">
		<!-- Form Column -->
		<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-6 sm:p-7 flex flex-col gap-5 shadow-xs">
			<h2 class="font-display font-bold text-base text-(--text-main)">Application Specifications</h2>

			{#if error}
				<div class="p-3 bg-red-950/10 text-(--color-danger) border border-(--color-danger)/30 rounded-md text-xs">{error}</div>
			{/if}
			{#if successMsg}
				<div class="p-3 bg-(--accent-soft) text-(--accent-strong) border border-(--accent-sky)/30 rounded-md text-xs">{successMsg}</div>
			{/if}

			<form on:submit|preventDefault={handleSubmit} class="flex flex-col gap-5">
				<!-- Project Name -->
				<div class="flex flex-col gap-1.5">
					<label for="p-name" class="text-xs font-medium text-(--text-main)">Project Name *</label>
					<input
						id="p-name"
						type="text"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
						required
						placeholder="e.g. Nexus Dashboard"
						bind:value={projectName}
						on:input={handleProjectNameInput}
					/>
				</div>

				<!-- Cloudflare-style Subdomain + Domain Picker -->
				<div class="flex flex-col gap-1.5">
					<label for="p-subdomain" class="text-xs font-medium text-(--text-main) flex items-center justify-between">
						<span>Desired Subdomain / Live Route *</span>
						<span class="text-xs text-(--text-muted) font-mono">https://&lt;subdomain&gt;{domainSuffix}</span>
					</label>
					<div class="flex items-stretch">
						<input
							id="p-subdomain"
							type="text"
							class="flex-1 px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-l-md text-xs font-mono text-(--text-main) outline-none focus:border-(--accent-sky) min-w-0"
							required
							placeholder="nexus"
							bind:value={subdomain}
							on:input={handleSubdomainInput}
						/>
						<div class="px-3.5 py-2 bg-(--bg-muted)/60 border border-l-0 border-(--border-hairline) rounded-r-md text-xs font-mono text-(--text-muted) flex items-center select-none">
							{domainSuffix}
						</div>
					</div>

					<!-- Real-time availability indicator -->
					{#if checkingSubdomain}
						<div class="flex items-center gap-1.5 text-xs text-(--text-muted) pt-0.5">
							<Spinner size={13} class="animate-spin" />
							<span>Checking availability...</span>
						</div>
					{:else if subdomainAvailable === true}
						<div class="flex items-center gap-1.5 text-xs text-emerald-600 dark:text-emerald-400 pt-0.5">
							<CheckCircle size={14} weight="bold" />
							<span>{subdomainMessage || 'Subdomain is available'}</span>
						</div>
					{:else if subdomainAvailable === false}
						<div class="flex items-center gap-1.5 text-xs text-rose-600 dark:text-rose-400 pt-0.5">
							<XCircle size={14} weight="bold" />
							<span>{subdomainMessage || 'Subdomain is unavailable'}</span>
						</div>
					{/if}
				</div>

				<!-- Short Description (Brief Purpose) with character counter -->
				<div class="flex flex-col gap-1.5">
					<div class="flex items-center justify-between">
						<label for="p-desc" class="text-xs font-medium text-(--text-main)">Short Summary (Overview Card)</label>
						<span class="text-xs text-(--text-muted) font-mono">{description.length}/280</span>
					</div>
					<textarea
						id="p-desc"
						rows="2"
						maxlength="280"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
						placeholder="A quiet and concise description of the project (shown on showcase cards)..."
						bind:value={description}
					></textarea>
				</div>

				<!-- Extended Documentation / README (.md upload or direct input) -->
				<div class="flex flex-col gap-1.5">
					<div class="flex items-center justify-between">
						<label for="p-readme" class="text-xs font-medium text-(--text-main) flex items-center gap-1.5">
							<BookOpen size={14} class="text-(--text-muted)" />
							<span>Extended README / Documentation (.md)</span>
						</label>
						<label class="btn btn-secondary btn-sm text-xs px-2.5 py-1 inline-flex items-center gap-1.5 cursor-pointer">
							<UploadSimple size={13} weight="bold" />
							<span>Upload .md file</span>
							<input
								type="file"
								accept=".md,.markdown,.txt"
								class="hidden"
								on:change={handleMdUpload}
							/>
						</label>
					</div>
					<textarea
						id="p-readme"
						rows="5"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs font-mono text-(--text-main) outline-none focus:border-(--accent-sky)"
						placeholder="# Project Overview&#10;&#10;Detailed architecture, setup notes, or user guide rendered like GitHub README..."
						bind:value={readme}
					></textarea>
					<span class="text-xs text-(--text-muted)">
						Supports markdown formatting: headings, code blocks, lists, links, and blockquotes.
					</span>
				</div>

				<!-- Simplified External Links: 1. Source Code, 2. Documentation -->
				<div class="flex flex-col gap-3 pt-2 border-t border-(--border-hairline)">
					<div class="flex items-center justify-between">
						<span class="text-xs font-semibold text-(--text-main)">External Links (Simplified)</span>
						<span class="text-xs text-(--text-muted)">Other demo/social links can go into the README</span>
					</div>

					<!-- 1. Source Code Link -->
					<div class="flex flex-col gap-1.5">
						<label for="p-repo" class="text-xs font-medium text-(--text-main) flex items-center gap-1.5">
							{#if repoInfo}
								<svelte:component this={repoInfo.icon} size={14} class="text-(--accent-sky)" />
							{:else}
								<Code size={14} class="text-(--text-muted)" />
							{/if}
							<span>Source Code Repository *</span>
						</label>
						<input
							id="p-repo"
							type="url"
							class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
							required
							placeholder="https://github.com/username/project"
							bind:value={repositoryUrl}
						/>
					</div>

					<!-- 2. Documentation Link -->
					<div class="flex flex-col gap-1.5">
						<label for="p-docs" class="text-xs font-medium text-(--text-main) flex items-center gap-1.5">
							{#if docInfo}
								<svelte:component this={docInfo.icon} size={14} class="text-(--accent-sky)" />
							{:else}
								<BookOpen size={14} class="text-(--text-muted)" />
							{/if}
							<span>Documentation Manual (Optional)</span>
						</label>
						<input
							id="p-docs"
							type="url"
							class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
							placeholder="https://docs.example.com"
							bind:value={documentationUrl}
						/>
					</div>
				</div>

				<!-- Technology Stack -->
				<div class="flex flex-col gap-1.5">
					<label for="p-tech" class="text-xs font-medium text-(--text-main)">Technology Stack (comma separated)</label>
					<input
						id="p-tech"
						type="text"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
						placeholder="SvelteKit, Go, PostgreSQL, Redis"
						bind:value={techStackRaw}
					/>
				</div>

				<!-- Deployment Architecture Notes -->
				<div class="flex flex-col gap-1.5">
					<label for="p-deploy" class="text-xs font-medium text-(--text-main)">Deployment Architecture Notes</label>
					<textarea
						id="p-deploy"
						rows="2"
						class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-md text-xs text-(--text-main) outline-none focus:border-(--accent-sky)"
						placeholder="Port requirements, Dockerfile location, or environment dependencies..."
						bind:value={deploymentNotes}
					></textarea>
				</div>

				<button type="submit" class="btn btn-primary btn-sm mt-2 self-start text-xs" disabled={submitting}>
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
						<div class="bg-(--bg-surface) border border-(--border-hairline) rounded-md p-5 flex flex-col gap-2.5 shadow-xs">
							<div class="flex items-center justify-between gap-2">
								<strong class="font-semibold text-sm text-(--text-main)">{req.project_name}</strong>
								<span class="inline-flex items-center gap-1.5 text-xs px-2.5 py-0.5 rounded-md border {req.status === 'PENDING' ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20' : req.status === 'COMPLETED' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20' : 'bg-(--bg-muted) text-(--text-secondary) border-(--border-hairline)'}">
									{req.status}
								</span>
							</div>

							{#if req.subdomain}
								<div class="text-xs font-mono text-(--accent-sky) flex items-center gap-1">
									<Globe size={13} />
									<span>https://{req.subdomain}{domainSuffix}</span>
								</div>
							{/if}

							{#if req.repository_url}
								<p class="text-xs text-(--text-muted) truncate font-mono">{req.repository_url}</p>
							{/if}

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
