import type { AvailabilityBlock } from './status';

export interface ServerHardware {
	cpu: string;
	cores: number;
	threads: number;
	mhz?: number;
}

export interface ServerMemory {
	totalGB: number;
	usagePercent: number;
	availableGB: number;
	totalBytes?: number;
	usedBytes?: number;
	availableBytes?: number;
}

export interface ServerStorage {
	physicalDiskGB: number;
	linuxTotalGB: number;
	totalGB: number;
	usagePercent: number;
	availableGB: number;
	type: string;
	model: string;
	summary: string;
	physicalModel?: string;
	totalBytes?: number;
	usedBytes?: number;
	availableBytes?: number;
}

export interface ServerNetwork {
	linkCapacity: string;
	latencyMs: number;
	linkMbps?: number;
	source?: string;
}

export interface ServerAvailability {
	current: 'operational' | 'degraded' | 'down';
	last30Days: number;
	last30DaysFormatted: string;
	recordedPeriod: string;
	uptime: string;
	uptimeSec: number;
	dailyBlocks?: AvailabilityBlock[];
}

export interface ServerProjects {
	total: number;
	online: number;
}

export interface ServerOS {
	arch: string;
	distro: string;
	kernel: string;
	location: string;
}

export interface PublicServerResponse {
	hardware: ServerHardware;
	memory: ServerMemory;
	storage: ServerStorage;
	network: ServerNetwork;
	availability: ServerAvailability;
	projects: ServerProjects;
	os: ServerOS;
}
