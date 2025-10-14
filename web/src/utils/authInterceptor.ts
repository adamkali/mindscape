import type { Middleware } from '@/api/runtime';
import { ResponseError } from '@/api/runtime';

/**
 * Authentication interceptor that automatically handles token expiration
 * and authentication failures by triggering logout
 */
export function createAuthInterceptor(onLogout: () => void): Middleware {
	return {
		async post(context) {
			const { response } = context;

			// Check for authentication errors (401, 403)
			if (response.status === 401 || response.status === 403) {
				console.warn(
					'Authentication failed or token expired. Logging out user.',
				);

				// Trigger logout callback
				onLogout();

				// Let the error propagate normally so the calling component can handle it
				throw new ResponseError(
					response,
					'Authentication failed or token expired',
				);
			}

			return response;
		},

		async onError(context) {
			const { error, response } = context;

			// Handle network errors that might indicate auth issues
			if (response && (response.status === 401 || response.status === 403)) {
				console.warn(
					'Authentication error detected in network error handler. Logging out user.',
				);
				onLogout();
			}

			// Return undefined to let the error propagate normally
			return undefined;
		},
	};
}
