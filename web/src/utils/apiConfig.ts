import { Configuration } from '@/api/runtime';
import { createAuthInterceptor } from './authInterceptor';

let globalApiConfig: Configuration | null = null;

/**
 * Initialize the global API configuration with authentication interceptor
 */
export function initializeApiConfig(logoutHandler: () => void): Configuration {
	globalApiConfig = new Configuration({
		// Inject the bearer token (BearerAuth security scheme) on every request,
		// read at request time so it always reflects the current session.
		apiKey: () => {
			const token = localStorage.getItem('jwt');
			return token ? `Bearer ${token}` : '';
		},
		middleware: [createAuthInterceptor(logoutHandler)],
	});

	return globalApiConfig;
}

/**
 * Get the global API configuration
 * Returns undefined if not initialized (for use in login/signup)
 */
export function getApiConfig(): Configuration | undefined {
	return globalApiConfig ?? undefined;
}

/**
 * Get the global API configuration with auth interceptor
 * Throws an error if not initialized - use for authenticated routes
 */
export function getAuthenticatedApiConfig(): Configuration {
	if (!globalApiConfig) {
		throw new Error(
			'API configuration not initialized. Call initializeApiConfig first.',
		);
	}
	return globalApiConfig;
}

/**
 * Check if API configuration is initialized
 */
export function isApiConfigInitialized(): boolean {
	return globalApiConfig !== null;
}
