import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal } from 'solid-js';
import { UsersApi, type UpdateCredentialsRequest } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { EmptyGuid } from '@/utils';
import { Header } from '@/components/Header';
import { Button, Input } from '@/components/atoms';

const EditProfile = () => {
	const auth = useAuth();
	const navigate = useNavigate();
	const api = new UsersApi();

	const user = auth.user();

	const [updateCredentialsRequest, setUpdateCredentialsRequest] = createSignal({
		username: '',
		email: '',
		password: '',
		oldPassword: '',
		id: EmptyGuid, // Default to an empty string
	} as UpdateCredentialsRequest);
	const [confirmPassword, setConfirmPassword] = createSignal('');
	const [passwordMatch, setPasswordMatch] = createSignal(true);

	const [profilePicture, setProfilePicture] = createSignal<string>('');
	const [selectedFile, setSelectedFile] = createSignal<File | null>(null);
	const [isLoading, setIsLoading] = createSignal(false);
	const [isLoadingPicture, setIsLoadingPicture] = createSignal(false);
	const [error, setError] = createSignal('');
	const [success, setSuccess] = createSignal('');

	createEffect(() => {
		if (user) {
			setUpdateCredentialsRequest({
				username: user.username,
				email: user.email,
				password: '',
				oldPassword: '',
				id: user.id,
			});
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

	const handleCredentialsSubmit = async (e: Event) => {
		e.preventDefault();
		if (!auth.token() || !user) return;

		setError('');
		setSuccess('');
		setIsLoading(true);

		if (passwordMatch()) {
			try {
				const response = await api.updateUser({
					authorization: `Bearer ${auth.token()}`,
					updateCredentialsRequest: updateCredentialsRequest(),
				});
				if (!response.success) {
					setError(response.message || 'Failed to update credentials');
				}
				if (response.success && response.jwt && response.data) {
					auth.update(response.data, response.jwt);
				} else {
					setError(response.message || 'Failed to update credentials');
				}
				setSuccess('Credentials updated successfully!');
			} catch (error: any) {
				setError(error.message || 'Failed to update credentials');
			} finally {
				setIsLoading(false);
			}
		}
	};

	const handleUsernameChange = (e: Event) => {
		const target = e.currentTarget as HTMLInputElement;
		setUpdateCredentialsRequest({
			...updateCredentialsRequest(),
			username: target.value,
		});
	};

	const handleEmailChange = (e: Event) => {
		const target = e.currentTarget as HTMLInputElement;
		setUpdateCredentialsRequest({
			...updateCredentialsRequest(),
			email: target.value,
		});
	};

	const handlePasswordChange = (e: Event) => {
		const target = e.currentTarget as HTMLInputElement;
		setUpdateCredentialsRequest({
			...updateCredentialsRequest(),
			password: target.value,
		});
	};

	const handleOldPasswordChange = (e: Event) => {
		const target = e.currentTarget as HTMLInputElement;
		setUpdateCredentialsRequest({
			...updateCredentialsRequest(),
			oldPassword: target.value,
		});
	};

	const handleConfirmPasswordChange = (e: Event) => {
		const target = e.currentTarget as HTMLInputElement;
		setConfirmPassword(target.value);
		if (target.value === updateCredentialsRequest().password) {
			setPasswordMatch(true);
		} else {
			setPasswordMatch(false);
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
			} else {
				setError('No changes to save');
				setIsLoading(false);
				return;
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
		<div class="min-h-screen bg-background">
			<Header />
			<div class="max-w-2xl mx-auto mt-4">
				<div class="bg-card text-card-foreground rounded-lg shadow-lg p-8 space-y-2">
					<div class="mb-8">
						<h1 class="text-3xl font-bold text-gray-900 dark:text-white">
							Edit Profile
						</h1>
						<p class="text-gray-600 dark:text-gray-400 mt-2">
							Update your profile information
						</p>
					</div>

					<form class="space-y-6" onSubmit={handleSubmit}>
						<div class="flex flex-row items-center justify-evenly space-y-4">
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

							<div class="flex flex-col space-y-2">
								<div >
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
									class="block w-full text-sm text-secondary-foreground file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-secondary file:text-secondary-foreground hover:file:bg-secondary/80"
									onChange={handleFileSelect}
								/>
								</div>
								<button
									type="submit"
									class="bg-secondary hover:bg-secondary/80 text-secondary-foreground font-bold py-2 px-4 rounded"
									disabled={isLoading()}
								>
									Save
								</button>
							</div>
						</div>
					</form>
					<form class="space-y-6" onSubmit={handleCredentialsSubmit}>
						<div class="space-x-4 flex flex-row">
							<div class="flex-1">
								<label
									for="username"
									class="block text-sm font-medium bg-primary text-foreground dark:text-slate-100 rounded-t-md pl-4"
								>
									Username
								</label>
								<input
									id="username"
									name="username"
									type="text"
									required
									class="appearance-none relative block w-full px-3 py-2 border border-primary text-primary background-card rounded-b-md sm:text-sm"
									placeholder="Enter your username"
									value={updateCredentialsRequest().username}
									onInput={(e) => handleUsernameChange(e)}
								/>
							</div>

							<div class="flex-1">
								<label
									for="email"
									class="block text-sm font-medium bg-primary text-foreground dark:text-slate-100 rounded-t-md pl-4"
								>
									Email
								</label>
								<input
									id="email"
									name="email"
									type="email"
									required
									class="appearance-none relative block w-full px-3 py-2 border border-primary text-primary background-card rounded-b-md sm:text-sm"
									placeholder="Enter your email"
									value={updateCredentialsRequest().email}
									onInput={(e) => handleEmailChange(e)}
								/>
							</div>
						</div>

						<div>
							<label
								for="password"
								class="block text-sm font-medium bg-primary text-slate-100"
							>
								Current Password
							</label>
							<input
								id="password"
								name="password"
								type="password"
								class="mt-1 appearance-none relative block w-full px-3 py-2 border border-primary text-primary background-card rounded-md sm:text-sm"
								placeholder="Enter your Current Password"
								value={updateCredentialsRequest().oldPassword}
								onInput={(e) => handleOldPasswordChange(e)}
							/>
						</div>
						<Input
							id="newPassword"
							name="newPassword"
							type="password"
							placeholder="Enter your New Password"
							value={updateCredentialsRequest().password}
							onInput={(e) => handlePasswordChange(e)}
							variant="primary"
							label={
								<span class="block text-sm font-medium bg-primary text-slate-100">
									New Password
								</span>
							}
						/>
						<Input 
							id="confirmPassword"
							name="confirmPassword"
							type="password"
							placeholder="Confirm your New Password"
							value={confirmPassword()}
							onInput={(e) => handleConfirmPasswordChange(e)}
							variant="primary"
							label={
								<span class="block text-sm font-medium bg-primary text-slate-100">
									Confirm Password
								</span>
							}
						/>
						<div>
							<Button 
								type="submit"
								disabled={isLoading()}
								variant="secondary"
							>
								{isLoading() ? 'Updating...' : 'Update Password'}
							</Button>

						</div>
					</form>
					<div class="mt-8 flex justify-between">
						{error() && <div class="text-red-600 text-sm">{error()}</div>}

						{success() && <div class="text-green-600 text-sm">{success()}</div>}

						<div class="flex space-x-4">
							<button
								type="button"
								onClick={handleCancel}
								class="flex-1 flex justify-center py-2 px-4 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white dark:border-gray-600 dark:hover:bg-gray-600"
							>
								Cancel
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default EditProfile;
