<script lang="ts">
	import CpuPackageVisual from './CpuPackageVisual.svelte';

	export let cpu: string = '13th Gen Intel® Core™ i5-1334U';
	export let cores: number = 10;
	export let threads: number = 12;

	$: cleanCpuName = cpu.replace(/Intel\(R\)\s+/g, '').replace(/Core\(TM\)\s+/g, 'Core ');
</script>

<div class="h-full w-full bg-(--bg-surface) rounded-3xl p-5 sm:p-6 border border-(--border-hairline) flex flex-col justify-between shadow-xs min-h-[240px] transition-colors duration-300">
	<!-- Header: Clean typography matching screenshot 2 -->
	<div>
		<div class="text-base font-normal text-(--text-main)">Core</div>
		<div class="text-xs text-(--text-secondary) font-normal mt-0.5 truncate" title={cleanCpuName}>
			{cleanCpuName}
		</div>
		<div class="w-6 h-[1px] bg-(--border-hairline) my-2.5"></div>
	</div>

	<!-- Enhanced Visual Silicon Processor -->
	<div class="my-auto py-1">
		<CpuPackageVisual model={cpu} {cores} {threads} />
	</div>

	<!-- Cores / Threads Stats -->
	<div class="grid grid-cols-2 gap-3 pt-3 border-t border-(--border-hairline)">
		<div>
			<div class="text-xl sm:text-2xl font-normal text-(--text-main) tracking-tight">{cores}</div>
			<div class="text-xs text-(--text-muted) font-normal">Cores</div>
		</div>
		<div>
			<div class="text-xl sm:text-2xl font-normal text-(--text-main) tracking-tight">{threads}</div>
			<div class="text-xs text-(--text-muted) font-normal">Threads</div>
		</div>
	</div>
</div>
