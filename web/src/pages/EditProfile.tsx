import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal } from 'solid-js';
import { UsersApi } from '@/api';
import { useAuth } from '@/contexts/AuthContext';

const EditProfile = () => {
	const auth = useAuth();
	const navigate = useNavigate();
	const api = new UsersApi();

	const user = auth.user();

	const [username, setUsername] = createSignal('');
	const [profilePicture, setProfilePicture] = createSignal<string>('');
	const [selectedFile, setSelectedFile] = createSignal<File | null>(null);
	const [isLoading, setIsLoading] = createSignal(false);
	const [isLoadingPicture, setIsLoadingPicture] = createSignal(false);
	const [error, setError] = createSignal('');
	const [success, setSuccess] = createSignal('');

	createEffect(() => {
		if (user) {
			setUsername(user.username || '');
		}
		if (auth.isAuthenticated() && auth.token()) {
			fetchProfilePicture();
		}
	});

	const fetchProfilePicture = async () => {
		if (!auth.token()) return;

		setIsLoadingPicture(true);
		try {
			const response = await api.getProfilePicture({
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

	const handleFileSelect = (e: Event) => {
		const target = e.currentTarget as HTMLInputElement;
		const file = target.files?.[0];
		if (file) {
			setSelectedFile(file);
		}
	};

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		if (!auth.token() || !user) return;

		setError('');
		setSuccess('');
		setIsLoading(true);

		try {
			if (selectedFile()) {
				await api.uploadProfilePicture({
					authorization: `Bearer ${auth.token()}`,
					file: selectedFile()!,
				});
				setSuccess('Profile updated successfully!');
				await fetchProfilePicture();
			} else if (username() !== user.username) {
				setSuccess('Username updated successfully!');
			} else {
				setError('No changes to save');
				setIsLoading(false);
				return;
			}

			if (username() !== user.username) {
				const updatedUser = { ...user, username: username() };
				auth.updateUser(updatedUser);
			}
		} catch (error: any) {
			setError(error.message || 'Failed to update profile');
		} finally {
			setIsLoading(false);
		}
	};

	const handleCancel = () => {
		navigate('/');
	};

	if (!auth.isAuthenticated() || !user) {
		navigate('/login');
		return null;
	}

	return (
		<div class="min-h-screen p-8">
			<div class="max-w-2xl mx-auto">
				<div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-8">
					<div class="mb-8">
						<h1 class="text-3xl font-bold text-gray-900 dark:text-white">
							Edit Profile
						</h1>
						<p class="text-gray-600 dark:text-gray-400 mt-2">
							Update your profile information
						</p>
					</div>

					<form class="space-y-6" onSubmit={handleSubmit}>
						<div class="flex flex-col items-center space-y-4">
							<div class="w-32 h-32 rounded-full overflow-hidden bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
								{isLoadingPicture() ? (
									<div class="text-gray-500 dark:text-gray-400">Loading...</div>
								) : profilePicture() ? (
									<img
										src={profilePicture()}
										alt="Profile picture"
										class="w-full h-full object-cover"
									/>
								) : (
									<div class="text-4xl text-gray-400 dark:text-gray-500">
										{user.username?.charAt(0).toUpperCase()}
									</div>
								)}
							</div>

							<div>
								<label
									for="profilePicture"
									class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
								>
									Update Profile Picture
								</label>
								<input
									id="profilePicture"
									name="profilePicture"
									type="file"
									accept="image/*"
									class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 dark:text-gray-400 dark:file:bg-blue-900 dark:file:text-blue-300"
									onChange={handleFileSelect}
								/>
							</div>
						</div>

						<div>
							<label
								for="username"
								class="block text-sm font-medium text-gray-700 dark:text-gray-300"
							>
								Username
							</label>
							<input
								id="username"
								name="username"
								type="text"
								required
								class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white bg-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Enter your username"
								value={username()}
								onInput={(e) => setUsername(e.currentTarget.value)}
							/>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
								Email
							</label>
							<input
								type="email"
								disabled
								class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-600 rounded-md sm:text-sm"
								value={user.email || ''}
							/>
							<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
								Email cannot be changed
							</p>
						</div>

						{error() && <div class="text-red-600 text-sm">{error()}</div>}

						{success() && <div class="text-green-600 text-sm">{success()}</div>}

						<div class="flex space-x-4">
							<button
								type="submit"
								disabled={isLoading()}
								class="flex-1 flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
							>
								{isLoading() ? 'Updating...' : 'Update Profile'}
							</button>

							<button
								type="button"
								onClick={handleCancel}
								class="flex-1 flex justify-center py-2 px-4 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white dark:border-gray-600 dark:hover:bg-gray-600"
							>
								Cancel
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
	);
};

export default EditProfile;
