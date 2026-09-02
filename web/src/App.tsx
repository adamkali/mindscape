import { Route, Router } from '@solidjs/router';
import type { JSX } from 'solid-js';
import ProtectedRoute from '@/components/ProtectedRoute';
import { FEATURES } from '@/config/features';
import { AuthProvider } from '@/contexts/AuthContext';
import { BackgroundProvider } from '@/contexts/BackgroundContext';
import { ApiKeys, EditProfile, Home, Households, Login, Signup } from '@/pages';
import AdminRoute from './components/AdminRoute';
import { Showcase } from './pages/admin';

const NotFound = () => <h1>404</h1>;

/**
 * Renders `children` only when the flag is on, otherwise the 404 page.
 *
 * The gate lives in the route's *component*, not around the `<Route>` itself:
 * Solid Router collects its route config from its children, and a `<Route>`
 * wrapped in `<Show>` still ends up registered — the path would resolve even
 * with the flag off. Gating the component is deterministic.
 */
const gated = (enabled: boolean, render: () => JSX.Element): JSX.Element =>
	enabled ? render() : <NotFound />;

const App = () => {
	return (
		<div class="h-screen overflow-hidden transition-colors">
			<AuthProvider>
				<BackgroundProvider>
					<Router>
						<Route path="/login" component={Login} />
						<Route path="/signup" component={Signup} />
						<Route
							path="/"
							component={() => (
								<ProtectedRoute>
									<Home />
								</ProtectedRoute>
							)}
						/>
						<Route
							path="/edit-profile"
							component={() => (
								<ProtectedRoute>
									<EditProfile />
								</ProtectedRoute>
							)}
						/>
						<Route
							path="/api-keys"
							component={() =>
								gated(FEATURES.apiKeys, () => (
									<ProtectedRoute>
										<ApiKeys />
									</ProtectedRoute>
								))
							}
						/>
						<Route
							path="/households"
							component={() => (
								<ProtectedRoute>
									<Households />
								</ProtectedRoute>
							)}
						/>
						<Route
							path="/admin/showcase"
							component={() =>
								gated(FEATURES.admin, () => (
									<ProtectedRoute>
										<AdminRoute>
											<Showcase />
										</AdminRoute>
									</ProtectedRoute>
								))
							}
						/>
						<Route path="*404" component={NotFound} />
					</Router>
				</BackgroundProvider>
			</AuthProvider>
		</div>
	);
};

export default App;
