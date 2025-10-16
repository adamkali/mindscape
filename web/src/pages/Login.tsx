import { A, useNavigate } from '@solidjs/router';
import { createSignal } from 'solid-js';
import { useAuth } from '@/contexts/AuthContext';
import { useBackgroundStyle } from '@/hooks/useBackground';
import { usePublicApi } from '@/utils/useApi';

const Login = () => {
	const [emailOrUsername, setEmailOrUsername] = createSignal('');
	const [password, setPassword] = createSignal('');
	const [isLoading, setIsLoading] = createSignal(false);
	const [error, setError] = createSignal('');

	const auth = useAuth();
	const navigate = useNavigate();
	const api = usePublicApi();
	const backgroundStyle = useBackgroundStyle();

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

			const response = await api.users.login({
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
		<div
			class="min-h-screen flex items-center justify-center px-4"
			style={backgroundStyle()}
		>
			<div class="max-w-md w-full space-y-8 bg-white/20 backdrop-blur-md border border-white/30 rounded-xl p-8 shadow-lg">
				<div class="text-center">
					<h2 class="text-3xl font-bold text-white">Sign in to your account</h2>
					<p class="mt-2 text-sm text-white/80">
						Don't have an account?{' '}
						<A
							href="/signup"
							class="font-medium text-blue-300 hover:text-blue-200"
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
								class="block text-sm font-medium text-white/90"
							>
								Email or Username
							</label>
							<input
								id="emailOrUsername"
								name="emailOrUsername"
								type="text"
								required
								class="mt-1 appearance-none relative block w-full px-3 py-2 bg-white/10 backdrop-blur-sm border border-white/30 placeholder-white/60 text-white rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent sm:text-sm"
								placeholder="Enter your email or username"
								value={emailOrUsername()}
								onInput={(e) => setEmailOrUsername(e.currentTarget.value)}
							/>
						</div>

						<div>
							<label
								for="password"
								class="block text-sm font-medium text-white/90"
							>
								Password
							</label>
							<input
								id="password"
								name="password"
								type="password"
								required
								class="mt-1 appearance-none relative block w-full px-3 py-2 bg-white/10 backdrop-blur-sm border border-white/30 placeholder-white/60 text-white rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent sm:text-sm"
								placeholder="Enter your password"
								value={password()}
								onInput={(e) => setPassword(e.currentTarget.value)}
							/>
						</div>
					</div>

					{error() && (
						<div class="text-red-300 text-sm text-center bg-red-500/20 backdrop-blur-sm border border-red-400/30 rounded-md p-2">
							{error()}
						</div>
					)}

					<div>
						<button
							type="submit"
							disabled={isLoading()}
							class="group relative w-full flex justify-center py-2 px-4 bg-white/20 backdrop-blur-sm border border-white/30 text-white hover:bg-white/30 text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
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
