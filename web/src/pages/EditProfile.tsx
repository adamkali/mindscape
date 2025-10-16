import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal, Show } from 'solid-js';
import { UsersApi,  type UpdateCredentialsRequest, ResponseError, BackgroundApi } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { useBackground, useBackgroundStyle } from '@/hooks/useBackground';
import { EmptyGuid } from '@/utils';
import { Header } from '@/components/Header';
import {
	Button,
	Card,
	CardBody,
	CardFooter,
	CardHeader,
	Input,
} from '@/components/atoms';
import BackgroundChoices from '@/components/BackgroundChoices';

const EditProfile = () => {
	const auth = useAuth();
	const navigate = useNavigate();
	const api = new UsersApi();
	const { 
		setUserBackground, 
		isLoadingChoices 
	} = useBackground();
	const backgroundStyle = useBackgroundStyle();

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
	const [customBackgroundFile, setCustomBackgroundFile] = createSignal<File | null>(null);
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

	const handleBackgroundSelect = async (backgroundUrl: string) => {
		try {
			await setUserBackground(backgroundUrl);
			setSuccess('Background updated! Changes are applied immediately.');
			setError('');
		} catch (error: any) {
			setError(error.message || 'Failed to update background');
			if ((error as ResponseError).response?.status === 401) {
				auth.logout();
			}
		}
	};

	const handleCustomBackgroundFileSelect = (e: Event) => {
		const target = e.currentTarget as HTMLInputElement;
		const file = target.files?.[0];
		if (file) {
			setCustomBackgroundFile(file);
			// Note: Custom background upload functionality would need backend support
			setSuccess('Custom background selected! Upload feature coming soon.');
			setError('');
		}
	};

	const handleCustomBackgroundUpload = async () => {
		if (!customBackgroundFile() || !auth.token()) return;

		setIsLoading(true);
		setError('');
		setSuccess('');

		try {
			const backgroundApi = new BackgroundApi();
			await backgroundApi.uploadBackground({
				authorization: `Bearer ${auth.token()}`,
				file: customBackgroundFile()!,

			});
			setSuccess('Custom background uploaded successfully!');
		} catch (error: any) {
			setError(error.message || 'Failed to upload custom background');
		} finally {
			setIsLoading(false);
		}
	};

	if (!auth.isAuthenticated() || !user) {
		navigate('/login');
		return null;
	}

	return (
		<div 
			class="min-h-screen bg-background"
			style={backgroundStyle()}
		>
			<Header />
			<div class="max-w-2xl mx-auto mt-4 space-y-4">
				<Card variant="glass">
					<CardHeader title="Edit Profile" subtitle="Update your profile information" />
					<form class="space-y-6" onSubmit={handleSubmit}>
						<CardBody padding="lg" class="flex flex-row space-x-4">
							<div class="flex flex-row items-center justify-evenly space-y-4">
								<div class="w-32 h-32 rounded-full overflow-hidden bg-white/10 backdrop-blur-md border border-white/20 flex items-center justify-center">
									{isLoadingPicture() ? (
										<div class="text-white/70">
											Loading...
										</div>
									) : profilePicture() ? (
										<img
											src={profilePicture()}
											alt="Profile picture"
											class="w-full h-full object-cover"
										/>
									) : (
										<div class="text-4xl text-white/70">
											{user.username?.charAt(0).toUpperCase()}
										</div>
									)}
								</div>

								<div class="flex flex-col space-y-2">
									<div>
										<label
											for="profilePicture"
											class="block text-sm font-medium text-white mb-2"
										>
											Update Profile Picture
										</label>
										<input
											id="profilePicture"
											name="profilePicture"
											type="file"
											accept="image/*"
											class="block w-full text-sm text-white file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-white/20 file:text-white hover:file:bg-white/30 file:backdrop-blur-md"
											onChange={handleFileSelect}
										/>
									</div>
								</div>
							</div>
						</CardBody>
						<CardFooter>
							<Button type="submit" disabled={isLoading()} variant="secondary">
								Save
							</Button>
							<Button onClick={handleCancel}>Cancel</Button>
						</CardFooter>
					</form>
				</Card>
				<Card variant="glass">
					<CardHeader title="Account Credentials" subtitle="Update your account information and change your password" />
					<form class="space-y-6" onSubmit={handleCredentialsSubmit}>
						<CardBody padding="lg">
							<div class="space-y-6">
								{/* Account Information */}
								<div>
									<h3 class="text-lg font-medium text-white mb-4">Account Information</h3>
									<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
										<Input
											id="username"
											name="username"
											type="text"
											required
											placeholder="Enter your username"
											value={updateCredentialsRequest().username}
											onInput={(e) => handleUsernameChange(e)}
											label="Username"
										/>
										<Input
											id="email"
											name="email"
											type="email"
											required
											placeholder="Enter your email"
											value={updateCredentialsRequest().email}
											onInput={(e) => handleEmailChange(e)}
											label="Email"
										/>
									</div>
								</div>

								{/* Password Section */}
								<div>
									<h3 class="text-lg font-medium text-white mb-4">Change Password</h3>
									<div class="space-y-4">
										<Input
											id="currentPassword"
											name="currentPassword"
											type="password"
											placeholder="Enter your current password"
											value={updateCredentialsRequest().oldPassword}
											onInput={(e) => handleOldPasswordChange(e)}
											label="Current Password"
										/>
										<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
											<Input
												id="newPassword"
												name="newPassword"
												type="password"
												placeholder="Enter your new password"
												value={updateCredentialsRequest().password}
												onInput={(e) => handlePasswordChange(e)}
												label="New Password"
											/>
											<Input
												id="confirmPassword"
												name="confirmPassword"
												type="password"
												placeholder="Confirm your new password"
												value={confirmPassword()}
												onInput={(e) => handleConfirmPasswordChange(e)}
												label="Confirm New Password"
											/>
										</div>
										{!passwordMatch() && (
											<div class="text-red-400 text-sm mt-2">Passwords do not match</div>
										)}
										<div class="text-sm text-white/70">
											Leave password fields blank if you don't want to change your password.
										</div>
									</div>
								</div>
							</div>
						</CardBody>
						<CardFooter>
							<Button type="submit" disabled={isLoading() || !passwordMatch()} variant="secondary">
								{isLoading() ? 'Updating...' : 'Update Credentials'}
							</Button>
							<Button type="button" onClick={handleCancel}>Cancel</Button>
						</CardFooter>
					</form>
				</Card>

				{/* Background Selection Section */}
				<Card variant="glass">
					<CardHeader title="Background Settings" subtitle="Choose from available backgrounds or upload a custom one" />
					<CardBody padding="lg">
						{/* Custom Background Upload */}
						<div class="mb-6">
							<label class="block text-sm font-medium text-white mb-2">
								Upload Custom Background
							</label>
							<div class="flex items-center space-x-4">
								<input
									type="file"
									accept="image/*"
									class="block w-full text-sm text-white file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-white/20 file:text-white hover:file:bg-white/30 file:backdrop-blur-md"
									onChange={handleCustomBackgroundFileSelect}
								/>
								<Show when={customBackgroundFile()}>
									<Button 
										onClick={handleCustomBackgroundUpload}
										disabled={isLoading()}
										variant="secondary"
									>
										{isLoading() ? 'Uploading...' : 'Upload'}
									</Button>
								</Show>
							</div>
						</div>

						{/* Available Background Choices */}
						<div>
							<label class="block text-sm font-medium text-white mb-4">
								Choose from Available Backgrounds
							</label>
							<Show when={isLoadingChoices()}>
								<div class="text-white/70">Loading backgrounds...</div>
							</Show>
							<BackgroundChoices handleBackgroundSelect={handleBackgroundSelect} />
						</div>
					</CardBody>
				</Card>

				{/* Status Messages */}
				{(error() || success()) && (
					<Card variant="glass">
						<CardBody padding="md">
							{error() && <div class="text-red-400 text-sm">{error()}</div>}
							{success() && <div class="text-green-400 text-sm">{success()}</div>}
						</CardBody>
					</Card>
				)}
			</div>
		</div>
	);
};

export default EditProfile;
