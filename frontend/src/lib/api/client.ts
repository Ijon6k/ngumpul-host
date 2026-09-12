import axios, { type AxiosRequestConfig } from 'axios';

/**
 * Base Axios client with credentials and JSON headers.
 * Centralizes all communication targeting the `/api` backend prefix.
 */
export const apiClient = axios.create({
	baseURL: '/api',
	withCredentials: true,
	headers: {
		'Content-Type': 'application/json'
	}
});

/**
 * Convenience typed HTTP wrapper around apiClient.
 */
export const http = {
	get: <T = any>(url: string, config?: AxiosRequestConfig) =>
		apiClient.get<T>(url, config).then((r) => r.data),

	post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
		apiClient.post<T>(url, data, config).then((r) => r.data),

	patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
		apiClient.patch<T>(url, data, config).then((r) => r.data),

	delete: <T = any>(url: string, config?: AxiosRequestConfig) =>
		apiClient.delete<T>(url, config).then((r) => r.data)
};

/**
 * Extracts a user-friendly error message from Axios errors or exceptions.
 */
export function extractError(err: any): string {
	if (err?.response?.data?.error) {
		return err.response.data.error;
	}
	if (err?.response?.data?.message) {
		return err.response.data.message;
	}
	if (err?.message) {
		return err.message;
	}
	return 'An unexpected error occurred. Please try again.';
}
