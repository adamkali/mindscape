import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal, createResource, For, Show } from 'solid-js';
import { UsersApi, BackgroundApi, UserApi, type UpdateCredentialsRequest } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
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

const EditProfile = () => {
	const auth = useAuth();
	const navigate = useNavigate();
	const api = new UsersApi();

	const user = auth.user();

	const [defaultBackground] = createResource(async () => {
		const backgroundApi = new BackgroundApi();
		const response = await backgroundApi.getDefaultBackground();
		if (response.success && response.data) {
			return response.data;
		} else {
			throw new Error(
				'Failed to fetch default background: ' + response.message,
			);
		}
	});

	const [backgroundChoices] = createResource(async () => {
		const backgroundApi = new BackgroundApi();
		const userApi = new UserApi();

		// Fetch global background choices
		const globalResponse = await backgroundApi.getBackgroundChoices();
		let globalChoices: string[] = [];

		if (globalResponse.success && globalResponse.data) {
			try {
				if (typeof globalResponse.data === 'string') {
					globalChoices = JSON.parse(globalResponse.data);
				} else {
					globalChoices = globalResponse.data;
				}
			} catch (error) {
				console.warn('Failed to parse global background choices as JSON, treating as string:', error);
				if (typeof globalResponse.data === 'string') {
					globalChoices = globalResponse.data.split(/[,\n]/).map(url => url.trim()).filter(url => url);
				} else {
					globalChoices = [globalResponse.data];
				}
			}
		}

		// Fetch user-specific background choices if authenticated
		let userChoices: string[] = [];
		if (auth.token()) {
			try {
				const userResponse = await userApi.getUserBackgroundChoices({
					authorization: `Bearer ${auth.token()}`,
				});

				if (userResponse.success && userResponse.data) {
					try {
						if (typeof userResponse.data === 'string') {
							userChoices = JSON.parse(userResponse.data);
						} else {
							userChoices = userResponse.data;
						}
					} catch (error) {
						console.warn('Failed to parse user background choices as JSON, treating as string:', error);
						if (typeof userResponse.data === 'string') {
							userChoices = userResponse.data.split(/[,\n]/).map(url => url.trim()).filter(url => url);
						} else {
							userChoices = [userResponse.data];
						}
					}
				}
			} catch (error) {
				console.warn('Failed to fetch user background choices:', error);
				// Continue with just global choices if user choices fail
			}
		}

		// Combine both global and user-specific choices, removing duplicates
		const allChoices = [...globalChoices, ...userChoices];
		return [...new Set(allChoices)]; // Remove duplicates using Set
	});

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
	const [selectedBackground, setSelectedBackground] = createSignal<string>('');
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

	createEffect(() => {
		// Set the initial selected background to the current default background
		if (defaultBackground()) {
			setSelectedBackground(defaultBackground()!);
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

	const handleBackgroundSelect = (backgroundUrl: string) => {
		setSelectedBackground(backgroundUrl);
		setSuccess('Background updated! Changes are applied immediately.');
		setError('');
	};

	const handleCustomBackgroundFileSelect = (e: Event) => {
		const target = e.currentTarget as HTMLInputElement;
		const file = target.files?.[0];
		if (file) {
			setCustomBackgroundFile(file);
			// Create a local URL for immediate preview
			const previewUrl = URL.createObjectURL(file);
			setSelectedBackground(previewUrl);
			setSuccess('Custom background selected! Upload to save permanently.');
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
			style={{ 
				'background-image': `url(${selectedBackground() || defaultBackground()})`,
				'background-size': 'cover',
				'background-position': 'center center',
				'background-repeat': 'no-repeat',
				'background-attachment': 'fixed'
			}}
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
					<CardHeader title="Update Credentials Data" subtitle="Update your credentials information. You do not need to change your password, just leave it blank." />
					<form class="space-y-6" onSubmit={handleCredentialsSubmit}>
						<CardBody padding="lg">
							<div class="space-x-4 flex flex-row">
								<div class="flex-1">
									<Input
										id="username"
										name="username"
										type="text"
										required
										placeholder="Enter your username"
										value={updateCredentialsRequest().username}
										onInput={(e) => handleUsernameChange(e)}
										label="Username"
										variant="primary"
									/>
								</div>

								<div class="flex-1">
									<Input
										id="email"
										name="email"
										type="email"
										required
										placeholder="Enter your email"
										value={updateCredentialsRequest().email}
										onInput={(e) => handleEmailChange(e)}
										label="Email"
										variant="primary"
									/>
								</div>
							</div>

							<Input
								id="password"
								name="password"
								type="password"
								placeholder="Enter your Current Password"
								value={updateCredentialsRequest().oldPassword}
								onInput={(e) => handleOldPasswordChange(e)}
								label="Current Password"
								variant="primary"
							/>
							<Input
								id="newPassword"
								name="newPassword"
								type="password"
								placeholder="Enter your New Password"
								value={updateCredentialsRequest().password}
								onInput={(e) => handlePasswordChange(e)}
								variant="primary"
								label="New Password"
							/>
							<Input
								id="confirmPassword"
								name="confirmPassword"
								type="password"
								placeholder="Confirm your New Password"
								value={confirmPassword()}
								onInput={(e) => handleConfirmPasswordChange(e)}
								variant="primary"
								label="Confirm Password"
							/>
							{!passwordMatch() && (
								<div class="text-red-400 text-sm">Passwords do not match</div>
							)}
						</CardBody>
						<CardFooter>
							<Button type="submit" disabled={isLoading()} variant="secondary">
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
							<Show when={backgroundChoices.loading}>
								<div class="text-white/70">Loading backgrounds...</div>
							</Show>
							<Show when={backgroundChoices.error}>
								<div class="text-red-400 text-sm">Failed to load background choices</div>
							</Show>
							<Show when={backgroundChoices()}>
								<div class="grid grid-cols-2 md:grid-cols-3 gap-4">
									<For each={backgroundChoices()}>
										{(backgroundUrl) => (
											<div 
												class={`relative aspect-video rounded-lg overflow-hidden cursor-pointer border-2 transition-all duration-300 hover:scale-105 ${
													selectedBackground() === backgroundUrl 
														? 'border-white/70 ring-2 ring-white/50' 
														: 'border-white/20 hover:border-white/40'
												}`}
												onClick={() => handleBackgroundSelect(backgroundUrl)}
											>
												<img 
													src={backgroundUrl} 
													alt="Background option"
													class="w-full h-full object-cover"
													onError={(e) => {
														(e.target as HTMLImageElement).style.display = 'none';
													}}
												/>
												<Show when={selectedBackground() === backgroundUrl}>
													<div class="absolute inset-0 bg-white/20 backdrop-blur-sm flex items-center justify-center">
														<div class="text-white font-semibold">Selected</div>
													</div>
												</Show>
											</div>
										)}
									</For>
								</div>
							</Show>
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
