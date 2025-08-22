import {
	createContext,
	createEffect,
	createSignal,
	type ParentComponent,
	useContext,
} from 'solid-js';
import type { ResponsesUserData } from '@/api';
import { UsersApi } from '@/api';

export interface AuthContextValue {
	user: () => ResponsesUserData | null;
	token: () => string | null;
	login: (token: string, userData: ResponsesUserData) => void;
	logout: () => void;
	isAuthenticated: () => boolean;
	isInitializing: () => boolean;
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
	const [isInitializing, setIsInitializing] = createSignal(!!storedToken && !parsedUser);

	const api = new UsersApi();

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

	const logout = () => {
		setToken(null);
		setUser(null);
		setIsInitializing(false);
	};

	const isAuthenticated = () => {
		return !!token() && !!user();
	};

	const update = (userData: ResponsesUserData, token: string) => {
		setUser(userData);
		setToken(token);
	};

	const value: AuthContextValue = {
		user,
		token,
		login,
		logout,
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
