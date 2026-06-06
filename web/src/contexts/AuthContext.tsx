import {
	createContext,
	createEffect,
	createSignal,
	type ParentComponent,
	useContext,
} from 'solid-js';
import type { ResponsesUserData } from '@/api';
import { UsersApi } from '@/api';
import { getApiConfig, initializeApiConfig } from '@/utils/apiConfig';
import { revokeSession, setTokenRefreshedHandler } from '@/utils/refreshAuth';

export interface AuthContextValue {
	user: () => ResponsesUserData | null;
	token: () => string | null;
	login: (token: string, userData: ResponsesUserData) => void;
	logout: () => void;
	isAuthenticated: () => boolean;
	isInitializing: () => boolean;
	isAdmin: () => boolean;
	update: (userData: ResponsesUserData, token: string) => void;
}

const AuthContext = createContext<AuthContextValue>();

export const AuthProvider: ParentComponent = (props) => {
	const storedToken = localStorage.getItem('jwt');
	const storedUser = localStorage.getItem('user');
	let parsedUser = null;
	try {
		parsedUser = storedUser ? JSON.parse(storedUser) : null;
	} catch {
		localStorage.removeItem('user');
	}

	const [token, setToken] = createSignal<string | null>(storedToken);
	const [user, setUser] = createSignal<ResponsesUserData | null>(parsedUser);
	const [isInitializing, setIsInitializing] = createSignal(
		!!storedToken && !parsedUser,
	);

	// Create logout function that will be used by the auth interceptor
	const logout = () => {
		console.log('Logging out user');
		// Best-effort: revoke this device's session server-side (other
		// browsers/devices stay logged in), then clear local state.
		void revokeSession();
		setToken(null);
		setUser(null);
		setIsInitializing(false);
		// Redirect to login page using window.location to avoid router dependency
		window.location.href = '/login';
	};

	// When the interceptor silently refreshes the access token, push the new
	// token into the auth signal so localStorage and headers stay current.
	setTokenRefreshedHandler((newToken) => {
		setToken(newToken);
	});

	// Initialize global API configuration with auth interceptor
	const apiConfig = initializeApiConfig(logout);
	const api = new UsersApi(apiConfig);

	createEffect(() => {
		const currentToken = token();
		const currentUser = user();

		if (currentToken) {
			localStorage.setItem('jwt', currentToken);
			if (currentUser) {
				localStorage.setItem('user', JSON.stringify(currentUser));
				// Make sure we're not initializing if we have both token and user
				setIsInitializing(false);
			} else {
				// Only fetch user if we don't have user data
				fetchCurrentUser(currentToken);
			}
		} else {
			localStorage.removeItem('jwt');
			localStorage.removeItem('user');
			setUser(null);
			setIsInitializing(false);
		}
	});

	const fetchCurrentUser = async (authToken: string) => {
		setIsInitializing(true);
		try {
			const response = await api.getCurrentLoggedInUser({
				authorization: `Bearer ${authToken}`,
			});
			if (response.success && response.data) {
				setUser(response.data);
			} else {
				logout();
			}
		} catch (error) {
			console.error('Failed to fetch current user:', error);
			logout();
		} finally {
			setIsInitializing(false);
		}
	};

	const login = (newToken: string, userData: ResponsesUserData) => {
		setToken(newToken);
		setUser(userData);
		setIsInitializing(false);
	};

	const isAuthenticated = () => {
		return !!token() && !!user();
	};

	const update = (userData: ResponsesUserData, token: string) => {
		setUser(userData);
		setToken(token);
	};

	const isAdmin = () => {
		return user()?.admin || false;
	};

	const value: AuthContextValue = {
		user,
		token,
		login,
		logout,
		isAdmin,
		isAuthenticated,
		isInitializing,
		update,
	};

	return (
		<AuthContext.Provider value={value}>{props.children}</AuthContext.Provider>
	);
};

export const useAuth = () => {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error('useAuth must be used within an AuthProvider');
	}
	return context;
};
