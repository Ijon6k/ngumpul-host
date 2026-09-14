<script lang="ts">
	import { onMount } from 'svelte';
	import { hostingRequestsApi, projectsApi, setupApi, extractError } from '$lib/api';
	import { domainSuffix, refreshNodeDomain } from '$lib/stores/node';
	import type { HostingRequest } from '$lib/types/hosting';
	import { detectLinkInfo } from '$lib/utils/linkDetector';
	import { formatDate } from '$lib/utils/format';
	import { Alert, Textarea, ImageUpload } from '$lib/components/ui';
	import { toast } from 'svelte-sonner';
	import {
		UploadSimple,
		CheckCircle,
		XCircle,
		Clock,
		Gear,
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
	let coverImageUrl = '';
	let repositoryUrl = '';
	let documentationUrl = '';
	let deploymentNotes = '';
	let techStackRaw = '';

	// Subdomain availability check state
	let checkingSubdomain = false;
	let subdomainAvailable: boolean | null = null;
	let subdomainMessage = '';
	let subdomainDebounceTimer: any = null;

	async function loadRequests() {
		try {
			const [res] = await Promise.all([
				hostingRequestsApi.getMyHostingRequests(),
				refreshNodeDomain()
			]);
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
				cover_image_url: coverImageUrl,
				repository_url: repositoryUrl.trim(),
				documentation_url: documentationUrl.trim(),
				deployment_notes: deploymentNotes.trim(),
				technology_stack: techStack
			});

			successMsg = 'Hosting request submitted successfully. The host administrator will review your project.';
			toast.success('Hosting request submitted successfully.');
			projectName = '';
			subdomain = '';
			subdomainAvailable = null;
			subdomainMessage = '';
			description = '';
			readme = '';
			coverImageUrl = '';
			repositoryUrl = '';
			documentationUrl = '';
			deploymentNotes = '';
			techStackRaw = '';
			await loadRequests();
		} catch (err) {
			error = extractError(err);
			toast.error(error);
		} finally {
			submitting = false;
		}
	}

	function getStatusLabel(status: string): string {
		switch (status) {
			case 'COMPLETED':
				return 'Published';
			case 'PENDING':
				return 'Reviewing';
			case 'SETUP':
			case 'APPROVED':
				return 'Setup';
			case 'REJECTED':
				return 'Rejected';
			default:
				return status;
		}
	}

	function getStatusChipClass(status: string): string {
		switch (status) {
			case 'COMPLETED':
				return 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20';
			case 'PENDING':
				return 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20';
			case 'SETUP':
			case 'APPROVED':
				return 'bg-sky-500/10 text-sky-600 dark:text-sky-400 border border-sky-500/20';
			case 'REJECTED':
				return 'bg-rose-500/10 text-rose-600 dark:text-rose-400 border border-rose-500/20';
			default:
				return 'bg-(--bg-muted) text-(--text-secondary) border border-(--border-hairline)';
		}
	}

	$: repoInfo = repositoryUrl ? detectLinkInfo(repositoryUrl, 'Source code') : null;
	$: docInfo = documentationUrl ? detectLinkInfo(documentationUrl, 'Documentation') : null;
</script>

<svelte:head>
	<title>Request Hosting · Ngumpul Host</title>
</svelte:head>

<div class="max-w-3xl w-full flex flex-col gap-10">
	<!-- Header -->
	<header class="flex flex-col gap-1.5 pb-6 border-b border-(--border-hairline)">
		<h1 class="font-sans font-normal text-2xl sm:text-3xl text-(--text-main) tracking-tight">Request hosting</h1>
		<p class="text-xs text-(--text-secondary) leading-relaxed">
			Submit a side project or experiment for hosting on the community server.
		</p>
	</header>

	{#if error}
		<Alert variant="danger">{error}</Alert>
	{/if}
	{#if successMsg}
		<Alert variant="success">{successMsg}</Alert>
	{/if}

	<!-- Progressive 3-Section Form -->
	<form on:submit|preventDefault={handleSubmit} class="flex flex-col gap-8">
		<!-- Section 1 — Project -->
		<section class="flex flex-col gap-4">
			<div>
				<h2 class="font-sans font-medium text-base text-(--text-main)">1. Project</h2>
				<p class="text-xs text-(--text-secondary)">Basic identification and summary for the project catalog.</p>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="p-name" class="text-xs font-medium text-(--text-main)">Project name <span class="text-(--accent-orange)">*</span></label>
				<input
					id="p-name"
					type="text"
					class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-sm text-xs text-(--text-main) outline-none focus:border-(--accent-sky) transition-colors"
					required
					placeholder="e.g. Acme Dashboard"
					bind:value={projectName}
					on:input={handleProjectNameInput}
				/>
			</div>

			<Textarea
				id="p-desc"
				label="Short description *"
				rows={2}
				maxlength={280}
				required
				placeholder="A concise description of the project (shown on catalog cards)..."
				bind:value={description}
			>
				<svelte:fragment slot="label-extra">
					<span class="text-xs text-(--text-muted) font-mono">{description.length}/280</span>
				</svelte:fragment>
			</Textarea>

			<!-- Cover Photo Upload (16:9 Standard) -->
			<ImageUpload
				bind:value={coverImageUrl}
				projectName={projectName}
				label="Project Cover"
				helperText="(Optional · 16:9 standard)"
			/>

			<div class="flex flex-col gap-1.5 pt-1">
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
				<Textarea
					id="p-readme"
					rows={4}
					class="font-mono"
					placeholder="# Architecture and Setup Notes&#10;&#10;Provide details on setup, technology requirements, or usage..."
					bind:value={readme}
					helperText="Optional. Formatted markdown will be displayed on the project page."
				/>
			</div>
		</section>

		<div class="h-px bg-(--border-hairline)"></div>

		<!-- Section 2 — Source -->
		<section class="flex flex-col gap-4">
			<div>
				<h2 class="font-sans font-medium text-base text-(--text-main)">2. Source & Stack</h2>
				<p class="text-xs text-(--text-secondary)">Code repository location and architectural technologies.</p>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="p-repo" class="text-xs font-medium text-(--text-main) flex items-center gap-1.5">
					{#if repoInfo}
						<svelte:component this={repoInfo.icon} size={14} class="text-(--accent-sky)" />
					{:else}
						<Code size={14} class="text-(--text-muted)" />
					{/if}
					<span>Repository URL <span class="text-(--accent-orange)">*</span></span>
				</label>
				<input
					id="p-repo"
					type="url"
					class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-sm text-xs font-mono text-(--text-main) outline-none focus:border-(--accent-sky) transition-colors"
					required
					placeholder="https://github.com/username/project"
					bind:value={repositoryUrl}
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="p-docs" class="text-xs font-medium text-(--text-main) flex items-center gap-1.5">
					{#if docInfo}
						<svelte:component this={docInfo.icon} size={14} class="text-(--accent-sky)" />
					{:else}
						<BookOpen size={14} class="text-(--text-muted)" />
					{/if}
					<span>Documentation URL (optional)</span>
				</label>
				<input
					id="p-docs"
					type="url"
					class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-sm text-xs font-mono text-(--text-main) outline-none focus:border-(--accent-sky) transition-colors"
					placeholder="https://docs.example.com"
					bind:value={documentationUrl}
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="p-tech" class="text-xs font-medium text-(--text-main)">Technology stack (comma separated)</label>
				<input
					id="p-tech"
					type="text"
					class="w-full px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-sm text-xs text-(--text-main) outline-none focus:border-(--accent-sky) transition-colors"
					placeholder="SvelteKit, Go, PostgreSQL, Redis"
					bind:value={techStackRaw}
				/>
			</div>
		</section>

		<div class="h-px bg-(--border-hairline)"></div>

		<!-- Section 3 — Hosting Details -->
		<section class="flex flex-col gap-4">
			<div>
				<h2 class="font-sans font-medium text-base text-(--text-main)">3. Hosting Details</h2>
				<p class="text-xs text-(--text-secondary)">Domain configuration and runtime environment requirements.</p>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="p-subdomain" class="text-xs font-medium text-(--text-main) flex items-center justify-between">
					<span>Preferred subdomain <span class="text-(--accent-orange)">*</span></span>
					<span class="text-xs text-(--text-muted) font-mono">https://&lt;subdomain&gt;{$domainSuffix}</span>
				</label>
				<div class="flex items-stretch">
					<input
						id="p-subdomain"
						type="text"
						class="flex-1 px-3 py-2 bg-(--bg-muted) border border-(--border-hairline) rounded-l-sm text-xs font-mono text-(--text-main) outline-none focus:border-(--accent-sky) min-w-0 transition-colors"
						required
						placeholder="acme"
						bind:value={subdomain}
						on:input={handleSubdomainInput}
					/>
					<div class="px-3.5 py-2 bg-(--bg-muted)/60 border border-l-0 border-(--border-hairline) rounded-r-sm text-xs font-mono text-(--text-muted) flex items-center select-none">
						{$domainSuffix}
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

			<Textarea
				id="p-deploy"
				label="Deployment notes"
				rows={2}
				placeholder="Port requirements, Dockerfile location, or environment dependencies..."
				bind:value={deploymentNotes}
			/>
		</section>

		<div class="h-px bg-(--border-hairline)"></div>

		<!-- Footer -->
		<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pt-1">
			<p class="text-xs text-(--text-secondary) max-w-md leading-relaxed">
				The host administrator will review the project before it is deployed.
			</p>
			<button
				type="submit"
				class="btn btn-primary btn-sm rounded-sm text-xs self-start sm:self-auto cursor-pointer"
				disabled={submitting}
			>
				{#if submitting}
					<Spinner size={13} class="animate-spin" />
					<span>Submitting request...</span>
				{:else}
					<span>Submit hosting request</span>
				{/if}
			</button>
		</div>
	</form>

	<!-- Your requests list (compact rows) -->
	<section class="mt-4 pt-10 border-t border-(--border-hairline) flex flex-col gap-4">
		<div class="flex items-center justify-between">
			<div>
				<h2 class="font-sans font-medium text-lg text-(--text-main)">Your requests</h2>
				<p class="text-xs text-(--text-secondary)">Past hosting requests and review statuses.</p>
			</div>
			{#if requests.length > 0}
				<span class="text-xs text-(--text-muted)">
					{requests.length} {requests.length === 1 ? 'request' : 'requests'}
				</span>
			{/if}
		</div>

		{#if loading}
			<div class="py-8 text-center text-xs text-(--text-muted) flex items-center justify-center gap-2">
				<Spinner size={14} class="animate-spin" />
				<span>Loading requests...</span>
			</div>
		{:else if requests.length === 0}
			<div class="p-8 text-center border border-(--border-hairline) rounded-md text-xs text-(--text-muted) bg-(--bg-surface)/30">
				No hosting requests submitted yet.
			</div>
		{:else}
			<div class="border border-(--border-hairline) rounded-md divide-y divide-(--border-hairline) bg-(--bg-surface)/30 overflow-hidden">
				{#each requests as req (req.id)}
					<div class="p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:bg-(--bg-muted)/40 transition-colors">
						<div class="flex flex-col gap-1 min-w-0">
							<div class="flex items-center gap-2 flex-wrap">
								<span class="font-medium text-xs text-(--text-main)">{req.project_name}</span>
								{#if req.subdomain}
									<span class="text-xs font-mono text-(--text-muted) flex items-center gap-1">
										<Globe size={13} />
										<span>{req.subdomain}{$domainSuffix}</span>
									</span>
								{/if}
							</div>
							{#if req.description}
								<p class="text-xs text-(--text-secondary) line-clamp-1 truncate max-w-xl">
									{req.description}
								</p>
							{/if}
							{#if req.admin_notes || req.operator_notes}
								<div class="mt-1 text-xs text-amber-600 dark:text-amber-400 bg-amber-500/5 px-2.5 py-1 rounded-sm border border-amber-500/15 inline-block self-start">
									<span class="font-medium">Operator note:</span> {req.admin_notes || req.operator_notes}
								</div>
							{/if}
						</div>

						<div class="flex items-center gap-3 shrink-0 self-end sm:self-center">
							<span class="text-xs text-(--text-muted)">
								{formatDate(req.created_at)}
							</span>
							<span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-[4px] text-xs font-medium {getStatusChipClass(req.status)}">
								{#if req.status === 'COMPLETED'}
									<CheckCircle size={12} weight="bold" />
								{:else if req.status === 'PENDING'}
									<Clock size={12} weight="bold" />
								{:else if req.status === 'SETUP' || req.status === 'APPROVED'}
									<Gear size={12} weight="bold" />
								{:else}
									<XCircle size={12} weight="bold" />
								{/if}
								<span>{getStatusLabel(req.status)}</span>
							</span>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</section>
</div>
