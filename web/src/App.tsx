import { Route, Router } from '@solidjs/router';
import { Home, EditProfile, Login, Signup } from '@/pages';
import { Showcase } from './pages/admin';

import ProtectedRoute from '@/components/ProtectedRoute';
import { AuthProvider } from '@/contexts/AuthContext';
import { BackgroundProvider } from '@/contexts/BackgroundContext';
import AdminRoute from './components/AdminRoute';

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
						path="/admin/showcase"
						component={() => (
							<ProtectedRoute>
								<AdminRoute>
									<Showcase />
								</AdminRoute>
							</ProtectedRoute>
						)}
					/>
					<Route path="*404" component={() => <h1>404</h1>} />
					</Router>
				</BackgroundProvider>
			</AuthProvider>
		</div>
	);
};

export default App;
