import {
	type Component,
	createEffect,
	createSignal,
	onMount,
	For,
} from 'solid-js';
import { useAuth } from '../contexts/AuthContext';
import { A } from '@solidjs/router';
import { UsersApi } from '@/api';

interface SearchEngine {
	name: string;
	placeholder: string;
	searchUrl: (query: string) => string;
}

const searchEngine: SearchEngine = {
	name: 'SearXNG',
	placeholder: 'Search SearXNG...',
	searchUrl: (query: string) =>
		`https://search.kalilarosa.xyz/?q=${encodeURIComponent(query)}`,
};

export const Header: Component = () => {
	const usersApi = new UsersApi();
	createEffect(() => {
		if (auth.isAuthenticated() && auth.token()) {
			fetchProfilePicture();
		}
	});

	onMount(() => {
		const savedDarkMode = localStorage.getItem('darkMode');
		if (savedDarkMode === 'true') {
			setDarkMode(true);
			document.documentElement.classList.add('dark');
		}
	});
	const handleLogout = () => {
		auth.logout();
	};
	const fetchProfilePicture = async () => {
		if (!auth.token()) return;
		setIsLoadingPicture(true);
		try {
			const response = await usersApi.getProfilePicture({
				authorization: `Bearer ${auth.token()}`,
			});
			if (response.data) {
				setProfilePicture(response.data);
			}
		} catch (error) {
			console.error('Failed to fetch profile picture:', error);
		} finally {
			setIsLoadingPicture(false);
		}
	};

	const auth = useAuth();
	const user = auth.user();

	const [profilePicture, setProfilePicture] = createSignal<string>('');
	const [isLoadingPicture, setIsLoadingPicture] = createSignal(false);
	const [darkMode, setDarkMode] = createSignal(false);
	const [searchQueries, setSearchQueries] = createSignal<
		Record<string, string>
	>({});
	const toggleDarkMode = () => {
		const newDarkMode = !darkMode();
		setDarkMode(newDarkMode);
		localStorage.setItem('darkMode', newDarkMode.toString());
		if (newDarkMode) {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	};

	const handleSearch = (engine: SearchEngine, query: string) => {
		if (query.trim()) {
			window.open(engine.searchUrl(query.trim()), '_blank');
		}
	};

	const updateSearchQuery = (engineName: string, query: string) => {
		setSearchQueries((prev) => ({ ...prev, [engineName]: query }));
	};
	return (
		<div class="border-b border-white/20 bg-white/10 backdrop-blur-lg shadow-lg shadow-slate-900/20">
			<div class="flex items-center justify-between p-1">
				{/* Logo */}
				<a href="/">
					<img width={200} src={'banner.svg'} alt="Mindscape" />
				</a>

				{/* Search Bars */}
				<div class="flex items-center space-x-2">
					<div class="flex items-center w-full">
						<input
							type="text"
							placeholder={searchEngine.placeholder}
							value={searchQueries()[searchEngine.name] || ''}
							onInput={(e) =>
								updateSearchQuery(searchEngine.name, e.currentTarget.value)
							}
							onKeyPress={(e) => {
								if (e.key === 'Enter') {
									handleSearch(searchEngine, e.currentTarget.value);
									updateSearchQuery(searchEngine.name, '');
								}
							}}
							class="px-3 py-1.5 text-xs bg-white/20 backdrop-blur-md border border-white/30 rounded-lg text-white placeholder:text-white/60 focus:outline-none focus:border-white/50 focus:bg-white/25 w-32 transition-all duration-200 shadow-sm hover:shadow-md"
						/>
					</div>
				</div>

				<div class="flex items-center space-x-4">
					{/* Dark mode toggle */}
					<button
						onClick={toggleDarkMode}
						class="p-2 rounded-xl bg-white/20 backdrop-blur-md border border-white/30 hover:bg-white/30 text-white transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
						title="Toggle dark mode"
					>
						{darkMode() ? '☀' : '🌙'}
					</button>

					{/* Profile section */}
					<div class="flex items-center space-x-3">
						<div class="w-8 h-8 rounded-full overflow-hidden bg-white/20 backdrop-blur-md border border-white/30 flex items-center justify-center shadow-lg">
							{isLoadingPicture() ? (
								<div class="text-xs text-foreground/60">...</div>
							) : profilePicture() ? (
								<img
									src={profilePicture()}
									alt={`${user?.username}'s profile`}
									class="w-full h-full object-cover"
								/>
							) : (
								<div class="text-sm text-foreground/60">
									{user?.username?.charAt(0).toUpperCase()}
								</div>
							)}
						</div>

						<span class="text-sm text-foreground">{user?.username}</span>

						<A
							href="/admin/showcase"
							class="text-xs px-3 py-1.5 bg-white/20 backdrop-blur-md border border-white/30 text-white rounded-lg hover:bg-white/30 transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
						>
							Dev Showcase
						</A>

						<A
							href="/edit-profile"
							class="text-xs px-3 py-1.5 bg-white/20 backdrop-blur-md border border-white/30 text-white rounded-lg hover:bg-white/30 transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
						>
							Edit
						</A>

						<button
							onClick={handleLogout}
							class="text-xs px-3 py-1.5 bg-white/15 backdrop-blur-md border border-white/25 text-white rounded-lg hover:bg-white/25 transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
						>
							Logout
						</button>
					</div>
				</div>
			</div>
		</div>
	);
};
