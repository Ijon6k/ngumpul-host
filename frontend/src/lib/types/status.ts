export type AvailabilityStatus = 'no_data' | 'operational' | 'partial' | 'incident';

export interface AvailabilityBlock {
	index: number;
	label: string;
	date: string;
	status: AvailabilityStatus;
	uptimePercent: number | null;
	details: string;
}

export interface IncidentRecord {
	id: number;
	started_at: string;
	ended_at: string;
	duration_seconds: number;
	duration_formatted: string;
	cause: string;
	details: string;
}

export interface AvailabilityData {
	current_status: 'operational' | 'degraded' | 'down';
	availability_percent: number | null;
	availability_formatted: string;
	recorded_period: string;
	recorded_days: number;
	has_sufficient_data: boolean;
	uptime_seconds: number;
	uptime_formatted: string;
	last_heartbeat_at: string;
	last_heartbeat_text: string;
	first_monitored_at: string;
	range_filter: '1d' | '7d' | '30d';
	total_incidents: number;
	total_downtime_seconds: number;
	total_downtime_formatted: string;
	blocks: AvailabilityBlock[];
	incidents: IncidentRecord[];
}

export interface ServiceHealth {
	id: string;
	name: string;
	status: 'OPERATIONAL' | 'DEGRADED' | 'DOWN';
	last_checked_at: string;
	response_time_ms: number;
}

export interface NetworkHealth {
	link_capacity: string;
	latency_ms: number;
	source: string;
}

export interface PublicStatusResponse {
	overall_status: 'OPERATIONAL' | 'DEGRADED' | 'DOWN';
	services: ServiceHealth[];
	availability: AvailabilityData | null;
	network: NetworkHealth;
	checked_at: string;
}
