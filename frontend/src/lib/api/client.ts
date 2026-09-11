import axios from 'axios';

/**
 * Base Axios client with credentials and JSON headers.
 */
export const apiClient = axios.create({
	baseURL: '/api',
	withCredentials: true,
	headers: {
		'Content-Type': 'application/json'
	}
});

/**
 * Extracts a user-friendly error message from Axios errors or exceptions.
 */
export function extractError(err: any): string {
	if (err.response?.data?.error) {
		return err.response.data.error;
	}
	if (err.response?.data?.message) {
		return err.response.data.message;
	}
	if (err.message) {
		return err.message;
	}
	return 'An unexpected error occurred. Please try again.';
}
