import { BackgroundApi, BookmarksApi, FoldersApi, UsersApi } from '@/api';
import { getApiConfig, getAuthenticatedApiConfig } from './apiConfig';

/**
 * Hook to get authenticated API instances with proper interceptor configuration
 * Use this for components that require authentication
 */
export function useAuthenticatedApi() {
	const config = getAuthenticatedApiConfig();

	return {
		users: new UsersApi(config),
		bookmarks: new BookmarksApi(config),
		folders: new FoldersApi(config),
		background: new BackgroundApi(config),
	};
}

/**
 * Hook to get API instances for public endpoints (login, signup)
 * Falls back to no config if auth interceptor isn't initialized yet
 */
export function usePublicApi() {
	const config = getApiConfig();

	return {
		users: new UsersApi(config),
		bookmarks: new BookmarksApi(config),
		folders: new FoldersApi(config),
		background: new BackgroundApi(config),
	};
}

/**
 * Get authenticated API instances with proper configuration
 */
export function getUsersApi() {
	return new UsersApi(getAuthenticatedApiConfig());
}

export function getBookmarksApi() {
	return new BookmarksApi(getAuthenticatedApiConfig());
}

export function getFoldersApi() {
	return new FoldersApi(getAuthenticatedApiConfig());
}

export function getBackgroundApi() {
	return new BackgroundApi(getAuthenticatedApiConfig());
}
