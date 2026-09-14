<script lang="ts">
	import { cn } from '$lib/utils';

	type Variant = 'block' | 'line' | 'avatar' | 'badge' | 'button';

	export let variant: Variant = 'block';
	export let animate: boolean = true;
	export let ariaLabel: string = '';

	let className: string = '';
	export { className as class };

	const variantStyles: Record<Variant, string> = {
		block: 'h-full w-full rounded-[6px]',
		line: 'h-3.5 w-full rounded-[4px]',
		avatar: 'h-10 w-10 rounded-full',
		badge: 'h-5 w-14 rounded-full',
		button: 'h-9 w-24 rounded-[4px]'
	};
</script>

<div
	class={cn(
		'bg-(--bg-muted) shrink-0 select-none pointer-events-none',
		animate && 'animate-pulse',
		variantStyles[variant],
		className
	)}
	role={ariaLabel ? 'status' : undefined}
	aria-hidden={ariaLabel ? undefined : 'true'}
	aria-label={ariaLabel || undefined}
	{...$$restProps}
></div>