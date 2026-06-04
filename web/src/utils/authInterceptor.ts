import type { Middleware } from '@/api/runtime';
import { ResponseError } from '@/api/runtime';
import { refreshAccessToken } from './refreshAuth';

/**
 * Replace the Authorization header (any casing) in a fetch init.
 */
function withAuthorization(init: RequestInit, token: string): RequestInit {
	const headers: Record<string, string> = {
		...(init.headers as Record<string, string>),
	};
	for (const key of Object.keys(headers)) {
		if (key.toLowerCase() === 'authorization') {
			delete headers[key];
		}
	}
	headers.Authorization = `Bearer ${token}`;
	return { ...init, headers };
}

/**
 * Authentication interceptor.
 *
 * 401 (expired access token): silently refresh the session once — all
 * concurrent 401s share a single refresh — then retry the original request
 * with the new token. Only logs the user out if the refresh itself fails
 * (session revoked or expired).
 *
 * 403 (real authorization failure): logout as before.
 */
export function createAuthInterceptor(onLogout: () => void): Middleware {
	return {
		async post(context) {
			const { response } = context;

			if (response.status === 401) {
				const newToken = await refreshAccessToken();
				if (newToken) {
					// Retry the original request once with the fresh token.
					// Plain fetch (not context.fetch) so the retry does not
					// re-enter this middleware — no infinite loop possible.
					return fetch(context.url, withAuthorization(context.init, newToken));
				}

				console.warn('Session refresh failed. Logging out user.');
				onLogout();
				throw new ResponseError(
					response,
					'Authentication failed or token expired',
				);
			}

			if (response.status === 403) {
				console.warn('Authorization failed. Logging out user.');
				onLogout();
				throw new ResponseError(response, 'Authorization failed');
			}

			return response;
		},
	};
}
