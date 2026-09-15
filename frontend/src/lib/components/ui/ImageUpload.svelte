<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Image, UploadSimple, Spinner, Trash } from 'phosphor-svelte';
	import { compressImageToWebP } from '$lib/utils/imageCompressor';
	import { projectsApi, extractError } from '$lib/api';
	import ProjectCover from '$lib/components/ui/ProjectCover.svelte';
	import { cn } from '$lib/utils';

	const dispatch = createEventDispatcher<{
		change: string;
		error: string;
	}>();

	export let value: string = '';
	export let label: string = 'Project Preview';
	export let helperText: string = '(Optional · 16:9 standard)';
	export let projectName: string = '';
	export let disabled: boolean = false;
	export let aspectRatio: string = '16/9';

	let className: string = '';
	export { className as class };

	let uploading = false;
	let errorMessage: string | null = null;
	let fileInput: HTMLInputElement;

	async function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		if (!target.files || target.files.length === 0) return;
		const file = target.files[0];

		if (file.size > 10 * 1024 * 1024) {
			errorMessage = 'Original file size exceeds maximum limit of 10 MB.';
			dispatch('error', errorMessage);
			return;
		}

		errorMessage = null;
		uploading = true;

		try {
			const result = await compressImageToWebP(file, {
				maxWidth: 1600,
				maxHeight: 1200,
				quality: 0.82
			});

			const url = await projectsApi.uploadCover(result.file);
			if (url) {
				value = url;
				dispatch('change', url);
			}
		} catch (err) {
			errorMessage = extractError(err);
			dispatch('error', errorMessage);
		} finally {
			uploading = false;
			if (fileInput) fileInput.value = '';
		}
	}

	function handleRemove() {
		value = '';
		errorMessage = null;
		dispatch('change', '');
	}
</script>

<div class={cn('flex flex-col gap-2 w-full', className)}>
	<div class="flex items-center justify-between">
		<span class="text-xs font-medium text-(--text-secondary) flex items-center gap-1.5">
			<span>{label}</span>
			{#if helperText}
				<span class="text-xs text-(--text-muted) font-normal">{helperText}</span>
			{/if}
		</span>
		{#if value && !disabled}
			<button
				type="button"
				on:click={handleRemove}
				class="text-xs text-rose-600 dark:text-rose-400 hover:underline flex items-center gap-1 cursor-pointer"
			>
				<Trash size={13} />
				<span>Remove preview</span>
			</button>
		{/if}
	</div>

	{#if value}
		<div class="relative w-full aspect-[16/9] rounded-md overflow-hidden border border-(--border-hairline) bg-(--bg-muted) group">
			<ProjectCover src={value} alt={projectName || 'Project preview'} name={projectName || 'Project'} {aspectRatio} />
			{#if !disabled}
				<div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
					<label class="btn btn-secondary btn-sm text-xs px-3 py-1.5 cursor-pointer inline-flex items-center gap-1.5 shadow-sm">
						<UploadSimple size={13} weight="bold" />
						<span>Change preview</span>
						<input
							type="file"
							accept="image/jpeg,image/png,image/webp"
							class="hidden"
							on:change={handleFileSelect}
							disabled={uploading}
						/>
					</label>
				</div>
			{/if}
		</div>
	{:else}
		<label
			class={cn(
				'border border-dashed border-(--border-hairline) hover:border-(--accent-sky) bg-(--bg-muted)/30 hover:bg-(--bg-muted)/60 rounded-md p-6 flex flex-col items-center justify-center gap-2.5 transition-colors text-center',
				disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
			)}
		>
			<input
				bind:this={fileInput}
				type="file"
				accept="image/jpeg,image/png,image/webp"
				class="hidden"
				on:change={handleFileSelect}
				disabled={uploading || disabled}
			/>
			{#if uploading}
				<div class="flex items-center gap-2 text-xs text-(--accent-sky)">
					<Spinner size={16} class="animate-spin" />
					<span>Optimizing and uploading image...</span>
				</div>
			{:else}
				<div class="w-10 h-10 rounded-full bg-(--bg-muted) flex items-center justify-center text-(--text-muted)">
					<Image size={20} />
				</div>
				<div>
					<p class="text-xs sm:text-sm font-medium text-(--text-main)">Click to upload project preview</p>
					<p class="text-xs text-(--text-muted) mt-0.5">JPEG, PNG, or WebP (auto-compressed to WebP, max 10 MB)</p>
				</div>
			{/if}
		</label>
	{/if}

	{#if errorMessage}
		<p class="text-xs text-(--color-danger)">{errorMessage}</p>
	{/if}
</div>
