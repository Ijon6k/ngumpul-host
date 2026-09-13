import { writable } from 'svelte/store';

export interface AdminBreadcrumbItem {
	label: string;
	href?: string;
}

export const adminBreadcrumb = writable<AdminBreadcrumbItem[] | null>(null);

export function setAdminBreadcrumb(items: AdminBreadcrumbItem[] | null) {
	adminBreadcrumb.set(items);
}
