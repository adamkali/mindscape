import { type Component, createEffect, createSignal, onMount } from 'solid-js';
import { useAuth } from '../contexts/AuthContext';
import { A } from '@solidjs/router';
import { UsersApi } from '@/api';

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
	return (
		<div class="border-b border-card-foreground/20 bg-card">
			<div class="flex items-center justify-between p-4">
				<h1 class="text-2xl font-bold text-foreground">Mindscape</h1>

				<div class="flex items-center space-x-4">
					{/* Dark mode toggle */}
					<button
						onClick={toggleDarkMode}
						class="p-2 rounded-lg bg-background hover:bg-background/80 text-foreground transition-colors"
						title="Toggle dark mode"
					>
						{darkMode() ? '☀' : '🌙'}
					</button>

					{/* Profile section */}
					<div class="flex items-center space-x-3">
						<div class="w-8 h-8 rounded-full overflow-hidden bg-card-foreground/20 flex items-center justify-center">
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
							class="text-xs px-2 py-1 bg-secondary text-secondary-foreground rounded hover:bg-secondary/80 transition-colors"
						>
							Dev Showcase	
						</A>

						<A
							href="/edit-profile"
							class="text-xs px-2 py-1 bg-secondary text-secondary-foreground rounded hover:bg-secondary/80 transition-colors"
						>
							Edit
						</A>

						<button
							onClick={handleLogout}
							class="text-xs px-2 py-1 bg-background text-foreground rounded border border-card-foreground/20 hover:bg-background/80 transition-colors"
						>
							Logout
						</button>
					</div>
				</div>
			</div>
		</div>
	);
};
