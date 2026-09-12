<script lang="ts">
	import { marked } from 'marked';
	import { FileText } from 'phosphor-svelte';

	export let content: string = '';
	export let title: string = '';
	export let bordered: boolean = false;
	let className: string = '';
	export { className as class };

	// Configure marked for GitHub Flavored Markdown
	marked.setOptions({
		gfm: true,
		breaks: true
	});

	$: html = content ? (marked.parse(content) as string) : '';
</script>

{#if html}
	{#if bordered || title}
		<div class="border border-(--border-hairline) rounded-md overflow-hidden bg-(--bg-surface)">
			{#if title}
				<div class="flex items-center gap-2 px-4 py-2.5 bg-(--bg-muted)/60 border-b border-(--border-hairline) text-xs font-mono text-(--text-secondary)">
					<FileText size={15} class="text-(--text-muted) shrink-0" />
					<span class="font-medium tracking-tight">{title}</span>
				</div>
			{/if}
			<div class="p-6 sm:p-8 markdown-content text-sm text-(--text-secondary) leading-relaxed {className}">
				{@html html}
			</div>
		</div>
	{:else}
		<div class="markdown-content text-sm text-(--text-secondary) leading-relaxed {className}">
			{@html html}
		</div>
	{/if}
{/if}

<style>
	:global(.markdown-content) {
		line-height: 1.75;
		word-break: break-word;
	}

	:global(.markdown-content > *:first-child) {
		margin-top: 0 !important;
	}

	:global(.markdown-content > *:last-child) {
		margin-bottom: 0 !important;
	}

	:global(.markdown-content h1) {
		font-size: 1.5rem;
		font-weight: 600;
		color: var(--text-main);
		border-bottom: 1px solid var(--border-hairline);
		padding-bottom: 0.5rem;
		margin-top: 2rem;
		margin-bottom: 1rem;
		letter-spacing: -0.02em;
	}

	:global(.markdown-content h2) {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
		border-bottom: 1px solid var(--border-hairline);
		padding-bottom: 0.375rem;
		margin-top: 1.75rem;
		margin-bottom: 0.75rem;
		letter-spacing: -0.01em;
	}

	:global(.markdown-content h3) {
		font-size: 1.125rem;
		font-weight: 600;
		color: var(--text-main);
		margin-top: 1.5rem;
		margin-bottom: 0.5rem;
	}

	:global(.markdown-content h4, .markdown-content h5, .markdown-content h6) {
		font-size: 1rem;
		font-weight: 600;
		color: var(--text-main);
		margin-top: 1.25rem;
		margin-bottom: 0.5rem;
	}

	:global(.markdown-content p) {
		margin-top: 0.75rem;
		margin-bottom: 0.75rem;
		color: var(--text-secondary);
	}

	:global(.markdown-content strong) {
		font-weight: 600;
		color: var(--text-main);
	}

	:global(.markdown-content em) {
		font-style: italic;
	}

	:global(.markdown-content a) {
		color: var(--accent-sky);
		text-decoration: underline;
		text-underline-offset: 2px;
		transition: opacity 0.15s ease;
	}
	:global(.markdown-content a:hover) {
		opacity: 0.85;
	}

	:global(.markdown-content code) {
		font-family: var(--font-mono, monospace);
		font-size: 0.8125rem;
		background-color: var(--bg-muted);
		color: var(--text-main);
		padding: 0.15rem 0.375rem;
		border-radius: 4px;
		border: 1px solid var(--border-hairline);
	}

	:global(.markdown-content pre) {
		background-color: var(--bg-muted);
		border: 1px solid var(--border-hairline);
		border-radius: 6px;
		padding: 1rem;
		margin: 1rem 0;
		overflow-x: auto;
		line-height: 1.5;
	}

	:global(.markdown-content pre code) {
		background-color: transparent;
		border: none;
		padding: 0;
		font-size: 0.8125rem;
		color: var(--text-main);
	}

	:global(.markdown-content ul) {
		list-style-type: disc;
		padding-left: 1.5rem;
		margin: 0.75rem 0;
	}

	:global(.markdown-content ol) {
		list-style-type: decimal;
		padding-left: 1.5rem;
		margin: 0.75rem 0;
	}

	:global(.markdown-content li) {
		margin: 0.25rem 0;
		color: var(--text-secondary);
	}

	:global(.markdown-content blockquote) {
		border-left: 3px solid var(--border-subtle);
		padding-left: 1rem;
		margin: 1rem 0;
		color: var(--text-muted);
		font-style: italic;
	}

	:global(.markdown-content hr) {
		border: 0;
		border-top: 1px solid var(--border-hairline);
		margin: 1.75rem 0;
	}

	:global(.markdown-content table) {
		width: 100%;
		border-collapse: collapse;
		margin: 1.25rem 0;
		font-size: 0.875rem;
		overflow-x: auto;
		display: block;
	}

	:global(.markdown-content th) {
		background-color: var(--bg-muted);
		color: var(--text-main);
		font-weight: 600;
		text-align: left;
		padding: 0.5rem 0.875rem;
		border: 1px solid var(--border-hairline);
	}

	:global(.markdown-content td) {
		padding: 0.5rem 0.875rem;
		border: 1px solid var(--border-hairline);
		color: var(--text-secondary);
	}

	:global(.markdown-content img) {
		max-width: 100%;
		border-radius: 6px;
		border: 1px solid var(--border-hairline);
		margin: 1rem 0;
	}
</style>
