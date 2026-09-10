import axios from 'axios';

// Browser client with credentials
export const api = axios.create({
	baseURL: '/api/v1',
	withCredentials: true,
	headers: {
		'Content-Type': 'application/json'
	}
});

// Response helper
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
