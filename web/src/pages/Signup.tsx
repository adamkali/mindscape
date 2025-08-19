import { A, useNavigate } from '@solidjs/router';
import { createSignal } from 'solid-js';
import { UsersApi } from '@/api';
import { useAuth } from '@/contexts/AuthContext';

const Signup = () => {
	const [email, setEmail] = createSignal('');
	const [username, setUsername] = createSignal('');
	const [password, setPassword] = createSignal('');
	const [confirmPassword, setConfirmPassword] = createSignal('');
	const [isLoading, setIsLoading] = createSignal(false);
	const [error, setError] = createSignal('');

	const auth = useAuth();
	const navigate = useNavigate();
	const api = new UsersApi();

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		setError('');

		if (password() !== confirmPassword()) {
			setError('Passwords do not match');
			return;
		}

		if (!email() || !username() || !password()) {
			setError('All fields are required');
			return;
		}

		setIsLoading(true);

		try {
			const response = await api.signup({
				signupRequest: {
					email: email(),
					username: username(),
					password: password(),
					isAdmin: false,
				},
			});

			if (response.success && response.jwt && response.data) {
				auth.login(response.jwt, response.data);
				navigate('/');
			} else {
				setError(response.message || 'Signup failed');
			}
		} catch (error: any) {
			setError(error.message || 'An error occurred during signup');
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div class="min-h-screen flex items-center justify-center px-4">
			<div class="max-w-md w-full space-y-8">
				<div class="text-center">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white">
						Create your account
					</h2>
					<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
						Already have an account?{' '}
						<A
							href="/login"
							class="font-medium text-blue-600 hover:text-blue-500"
						>
							Sign in
						</A>
					</p>
				</div>

				<form class="mt-8 space-y-6" onSubmit={handleSubmit}>
					<div class="space-y-4">
						<div>
							<label
								for="email"
								class="block text-sm font-medium text-gray-700 dark:text-gray-300"
							>
								Email address
							</label>
							<input
								id="email"
								name="email"
								type="email"
								required
								class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white bg-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Enter your email"
								value={email()}
								onInput={(e) => setEmail(e.currentTarget.value)}
							/>
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
								placeholder="Choose a username"
								value={username()}
								onInput={(e) => setUsername(e.currentTarget.value)}
							/>
						</div>

						<div>
							<label
								for="password"
								class="block text-sm font-medium text-gray-700 dark:text-gray-300"
							>
								Password
							</label>
							<input
								id="password"
								name="password"
								type="password"
								required
								class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white bg-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Enter your password"
								value={password()}
								onInput={(e) => setPassword(e.currentTarget.value)}
							/>
						</div>

						<div>
							<label
								for="confirmPassword"
								class="block text-sm font-medium text-gray-700 dark:text-gray-300"
							>
								Confirm Password
							</label>
							<input
								id="confirmPassword"
								name="confirmPassword"
								type="password"
								required
								class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white bg-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Confirm your password"
								value={confirmPassword()}
								onInput={(e) => setConfirmPassword(e.currentTarget.value)}
							/>
						</div>
					</div>

					{error() && (
						<div class="text-red-600 text-sm text-center">{error()}</div>
					)}

					<div>
						<button
							type="submit"
							disabled={isLoading()}
							class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							{isLoading() ? 'Creating Account...' : 'Sign up'}
						</button>
					</div>
				</form>
			</div>
		</div>
	);
};

export default Signup;
