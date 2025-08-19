import { Route, Router } from '@solidjs/router';
import ProtectedRoute from '@/components/ProtectedRoute';
import { AuthProvider } from '@/contexts/AuthContext';
import EditProfile from '@/pages/EditProfile';
import Home from '@/pages/Home';
import Login from '@/pages/Login';
import Signup from '@/pages/Signup';

const App = () => {
	return (
		<div class="min-h-screen transition-colors">
			<AuthProvider>
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
				</Router>
			</AuthProvider>
		</div>
	);
};

export default App;
