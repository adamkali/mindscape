import {
	type Component,
	createEffect,
	createSignal,
	onMount,
	Show,
} from 'solid-js';
import { useAuth } from '../contexts/AuthContext';
import { A } from '@solidjs/router';
import { UsersApi } from '@/api';

interface SearchEngine {
	name: string;
	placeholder: string;
	searchUrl: (query: string) => string;
}


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
	const [isDropdownOpen, setIsDropdownOpen] = createSignal(false);
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

	const toggleDropdown = () => {
		setIsDropdownOpen(!isDropdownOpen());
	};

	const closeDropdown = () => {
		setIsDropdownOpen(false);
	};

	// Close dropdown when clicking outside
	createEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			const target = event.target as HTMLElement;
			if (!target.closest('.profile-dropdown')) {
				setIsDropdownOpen(false);
			}
		};

		if (isDropdownOpen()) {
			document.addEventListener('click', handleClickOutside);
		}

		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});
	return (
		<div class="border-b border-white/20 bg-white/10 backdrop-blur-lg shadow-lg shadow-slate-900/20">
			<div class="flex items-center justify-between p-1">
				{/* Logo */}
				<a href="/">
					<img width={175} src={'banner.svg'} alt="Mindscape" />
				</a>


				<div class="flex items-center space-x-4">
					{/* Profile section with dropdown */}
					<div class="relative profile-dropdown">
						<div 
							class="flex items-center space-x-3 cursor-pointer"
							onClick={toggleDropdown}
						>
							<div class="w-8 h-8 rounded-full overflow-hidden bg-white/20 backdrop-blur-md border border-white/30 flex items-center justify-center shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105">
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
							<div class={`text-white text-xs transition-transform duration-200 ${isDropdownOpen() ? 'rotate-180' : ''}`}>
								▼
							</div>
						</div>

						{/* Dropdown Menu */}
						<Show when={isDropdownOpen()}>
							<div class="absolute right-0 top-full mt-2 w-48 bg-white/20 backdrop-blur-md border border-white/30 rounded-xl shadow-lg z-50">
								<div class="py-2">
									<A
										href="/edit-profile"
										class="flex items-center px-4 py-2 text-sm text-white hover:bg-white/20 transition-all duration-200"
										onClick={closeDropdown}
									>
										<div class="w-4 h-4 mr-3 text-center">👤</div>
										Edit Profile
									</A>
									<button
										onClick={() => {
											toggleDarkMode();
											closeDropdown();
										}}
										class="flex items-center w-full px-4 py-2 text-sm text-white hover:bg-white/20 transition-all duration-200 text-left"
									>
										<div class="w-4 h-4 mr-3 text-center">{darkMode() ? '☀️' : '🌙'}</div>
										{darkMode() ? 'Light Mode' : 'Dark Mode'}
									</button>
									<button
										onClick={() => {
											handleLogout();
											closeDropdown();
										}}
										class="flex items-center w-full px-4 py-2 text-sm text-white hover:bg-white/20 transition-all duration-200 text-left"
									>
										<div class="w-4 h-4 mr-3 text-center">🚪</div>
										Logout
									</button>
								</div>
							</div>
						</Show>
					</div>
				</div>
			</div>
		</div>
	);
};
