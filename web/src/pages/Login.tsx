import { A, useNavigate } from '@solidjs/router';
import { createSignal } from 'solid-js';
import { UsersApi } from '@/api';
import { useAuth } from '@/contexts/AuthContext';

const Login = () => {
	const [emailOrUsername, setEmailOrUsername] = createSignal('');
	const [password, setPassword] = createSignal('');
	const [isLoading, setIsLoading] = createSignal(false);
	const [error, setError] = createSignal('');

	const auth = useAuth();
	const navigate = useNavigate();
	const api = new UsersApi();

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		setError('');

		if (!emailOrUsername() || !password()) {
			setError('All fields are required');
			return;
		}

		setIsLoading(true);

		try {
			const loginRequest = emailOrUsername().includes('@')
				? { email: emailOrUsername(), password: password() }
				: { username: emailOrUsername(), password: password() };

			const response = await api.login({
				loginRequest,
			});

			if (response.success && response.jwt && response.data) {
				auth.login(response.jwt, response.data);
				navigate('/');
			} else {
				setError(response.message || 'Login failed');
			}
		} catch (error: any) {
			setError(error.message || 'An error occurred during login');
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div class="min-h-screen flex items-center justify-center px-4">
			<div class="max-w-md w-full space-y-8">
				<div class="text-center">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white">
						Sign in to your account
					</h2>
					<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
						Don't have an account?{' '}
						<A
							href="/signup"
							class="font-medium text-blue-600 hover:text-blue-500"
						>
							Sign up
						</A>
					</p>
				</div>

				<form class="mt-8 space-y-6" onSubmit={handleSubmit}>
					<div class="space-y-4">
						<div>
							<label
								for="emailOrUsername"
								class="block text-sm font-medium text-gray-700 dark:text-gray-300"
							>
								Email or Username
							</label>
							<input
								id="emailOrUsername"
								name="emailOrUsername"
								type="text"
								required
								class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white bg-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Enter your email or username"
								value={emailOrUsername()}
								onInput={(e) => setEmailOrUsername(e.currentTarget.value)}
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
							{isLoading() ? 'Signing in...' : 'Sign in'}
						</button>
					</div>
				</form>
			</div>
		</div>
	);
};

export default Login;
