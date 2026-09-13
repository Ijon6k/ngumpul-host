export { apiClient, apiClient as api, http, extractError } from './client';
export { projectsApi } from './projects';
export { commentsApi } from './comments';
export { reportsApi } from './reports';
export { notificationsApi } from './notifications';
export { hostingRequestsApi } from './hostingRequests';
export { activityApi } from './activity';
export { systemApi } from './system';
export { authApi } from './auth';
export { setupApi } from './setup';
export { usersApi } from './users';
export * as adminApi from './admin';

// Retain backwards compatibility for legacy helpers
export { fetchPublicStatus, createStatusQuery } from './status';
export { fetchPublicServer, createPublicServerQuery } from './server';
