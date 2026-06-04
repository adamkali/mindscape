/**
 * Silent session refresh.
 *
 * The refresh token lives in an httpOnly cookie scoped to
 * /api/users/refresh — JavaScript never sees it; the browser attaches it
 * automatically when we call the endpoint with credentials: 'include'.
 *
 * refreshAccessToken() is single-flight: if several requests hit a 401 at
 * the same time, they all await the same in-flight refresh instead of
 * racing each other (the backend rotates the token on every refresh, so a
 * second concurrent refresh would 401).
 */

const REFRESH_URL = '/api/users/refresh';

let refreshPromise: Promise<string | null> | null = null;
let onTokenRefreshed: ((token: string) => void) | null = null;

/**
 * Register a callback invoked whenever a refresh succeeds, so app state
 * (e.g. the AuthContext token signal) can pick up the new access token.
 */
export function setTokenRefreshedHandler(handler: (token: string) => void) {
	onTokenRefreshed = handler;
}

/**
 * Exchange the refresh cookie for a new access token.
 * Returns the new token, or null if the session is gone (re-login needed).
 */
export function refreshAccessToken(): Promise<string | null> {
	if (!refreshPromise) {
		refreshPromise = doRefresh().finally(() => {
			refreshPromise = null;
		});
	}
	return refreshPromise;
}

async function doRefresh(): Promise<string | null> {
	try {
		const res = await fetch(REFRESH_URL, {
			method: 'POST',
			credentials: 'include',
		});
		if (!res.ok) {
			return null;
		}
		const body = await res.json();
		if (body?.success && body?.jwt) {
			localStorage.setItem('jwt', body.jwt);
			onTokenRefreshed?.(body.jwt);
			return body.jwt as string;
		}
		return null;
	} catch {
		return null;
	}
}

/**
 * Revoke this device's session on the server (DELETE on the refresh path so
 * the path-scoped cookie is sent) and let the server clear the cookie.
 * Other devices' sessions are untouched.
 */
export async function revokeSession(): Promise<void> {
	try {
		await fetch(REFRESH_URL, {
			method: 'DELETE',
			credentials: 'include',
		});
	} catch {
		// best-effort: local logout proceeds regardless
	}
}
