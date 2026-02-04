import { Route, Router } from '@solidjs/router';
import ProtectedRoute from '@/components/ProtectedRoute';
import { AuthProvider } from '@/contexts/AuthContext';
import { BackgroundProvider } from '@/contexts/BackgroundContext';
import { EditProfile, Home, Login, Signup } from '@/pages';
import AdminRoute from './components/AdminRoute';
import { Showcase } from './pages/admin';

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
