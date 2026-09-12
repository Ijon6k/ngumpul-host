import { http } from './client';

export const systemApi = {
	/**
	 * Gets public host telemetry status.
	 */
	getPublicStatus: () => http.get('/status'),

	/**
	 * Gets public node device information.
	 */
	getPublicServer: () => http.get('/public/server'),

	/**
	 * Checks backend service health.
	 */
	getHealth: () => http.get('/health')
};
