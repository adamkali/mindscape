import { A } from '@solidjs/router';
import {
	type Component,
	createEffect,
	createSignal,
	onMount,
	Show,
} from 'solid-js';
import { UsersApi } from '@/api';
import { FEATURES } from '@/config/features';
import { getAuthenticatedApiConfig } from '@/utils/apiConfig';
import { useAuth } from '../contexts/AuthContext';
import { useView } from '../contexts/ViewContext';
import { AgendaIcon, WidgetIcon } from './icons';

export const Header: Component = () => {
	const usersApi = new UsersApi(getAuthenticatedApiConfig());
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
			const response = await usersApi.getProfilePicture({});
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
	const view = useView();
	const user = auth.user();

	const [profilePicture, setProfilePicture] = createSignal<string>('');
	const [isLoadingPicture, setIsLoadingPicture] = createSignal(false);
	const [darkMode, setDarkMode] = createSignal(false);
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
		<div class="relative z-50 border-b border-white/20 bg-glass-bg backdrop-blur-lg shadow-lg shadow-slate-900/20 dark:border-slate-700/50 dark:shadow-black/30">
			<div class="flex items-center justify-between p-1">
				{/* Logo */}
				<a href="/">
					<img width={175} src={'banner.svg'} alt="Mindscape" />
				</a>

				{/* Homepage layout switcher — mouse equivalent of nav-mode t / r */}
				<div class="flex items-center rounded-lg border border-glass-border bg-glass-bg/30 backdrop-blur-md overflow-hidden">
					<button
						type="button"
						onClick={() => view.setHomepageMode('normal')}
						disabled={view.homepageMode() === 'normal'}
						title="Normal layout (r)"
						class={`flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium transition-all duration-200 ${
							view.homepageMode() === 'normal'
								? 'bg-glass-bg-hover text-foreground'
								: 'text-foreground/50 hover:text-foreground/80'
						}`}
					>
						Normal
					</button>
					<button
						type="button"
						onClick={() => view.setHomepageMode('tree')}
						disabled={view.homepageMode() === 'tree'}
						title="Tree layout (t)"
						class={`flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium transition-all duration-200 ${
							view.homepageMode() === 'tree'
								? 'bg-glass-bg-hover text-foreground'
								: 'text-foreground/50 hover:text-foreground/80'
						}`}
					>
						Tree
					</button>
				</div>

				{/* View Switcher — only meaningful when Tasks/Agenda is enabled,
				    and only applies inside the Normal layout */}
				<Show when={FEATURES.tasks && view.homepageMode() === 'normal'}>
					<div class="flex items-center rounded-lg border border-glass-border bg-glass-bg/30 backdrop-blur-md overflow-hidden">
						<button
							type="button"
							onClick={() => view.setActiveView('widgets')}
							disabled={view.activeView() === 'widgets'}
							class={`flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium transition-all duration-200 ${
								view.activeView() === 'widgets'
									? 'bg-glass-bg-hover text-foreground'
									: 'text-foreground/50 hover:text-foreground/80'
							}`}
						>
							<WidgetIcon class="text-base" />
							Widgets
						</button>
						<button
							type="button"
							onClick={() => view.setActiveView('agenda')}
							disabled={view.activeView() === 'agenda'}
							class={`flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium transition-all duration-200 ${
								view.activeView() === 'agenda'
									? 'bg-glass-bg-hover text-foreground'
									: 'text-foreground/50 hover:text-foreground/80'
							}`}
						>
							<AgendaIcon class="text-base" />
							Agenda
						</button>
					</div>
				</Show>

				<div class="flex items-center space-x-4">
					{/* Profile section with dropdown */}
					<div class="relative profile-dropdown">
						<button
							type="button"
							class="flex items-center space-x-3 cursor-pointer"
							onClick={toggleDropdown}
							aria-haspopup="true"
							aria-expanded={isDropdownOpen()}
						>
							<div class="w-8 h-8 rounded-full overflow-hidden bg-white/20 backdrop-blur-md border border-white/30 flex items-center justify-center shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105 dark:bg-slate-900/40 dark:border-slate-700/50 dark:shadow-black/30">
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
							<div
								class={`text-foreground text-xs transition-transform duration-200 ${isDropdownOpen() ? 'rotate-180' : ''}`}
							>
								▼
							</div>
						</button>

						{/* Dropdown Menu */}
						<Show when={isDropdownOpen()}>
							<div class="absolute right-0 top-full mt-2 w-48 bg-glass-bg-strong backdrop-blur-md border border-white/30 rounded-xl shadow-lg z-50 dark:shadow-black/30">
								<div class="py-2">
									<A
										href="/edit-profile"
										class="flex items-center px-4 py-2 text-sm text-foreground transition-all duration-200"
										onClick={closeDropdown}
									>
										<div class="w-4 h-4 mr-3 text-center">👤</div>
										Edit Profile
									</A>
									<A
										href="/households"
										class="flex items-center px-4 py-2 text-sm text-foreground transition-all duration-200"
										onClick={closeDropdown}
									>
										<div class="w-4 h-4 mr-3 text-center">🏠</div>
										Households
									</A>
									<Show when={FEATURES.apiKeys}>
										<A
											href="/api-keys"
											class="flex items-center px-4 py-2 text-sm text-foreground transition-all duration-200"
											onClick={closeDropdown}
										>
											<div class="w-4 h-4 mr-3 text-center">🔑</div>
											API Keys
										</A>
									</Show>
									<button
										type="button"
										onClick={() => {
											toggleDarkMode();
											closeDropdown();
										}}
										class="flex items-center w-full px-4 py-2 text-sm text-foreground transition-all duration-200 text-left"
									>
										<div class="w-4 h-4 mr-3 text-center">
											{darkMode() ? '☀️' : '🌙'}
										</div>
										{darkMode() ? 'Light Mode' : 'Dark Mode'}
									</button>
									<button
										type="button"
										onClick={() => {
											handleLogout();
											closeDropdown();
										}}
										class="flex items-center w-full px-4 py-2 text-sm text-foreground transition-all duration-200 text-left"
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
