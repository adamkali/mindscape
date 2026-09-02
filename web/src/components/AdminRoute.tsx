import { useNavigate } from '@solidjs/router';
import { createEffect, type ParentComponent } from 'solid-js';
import { FEATURES } from '@/config/features';
import { useAuth } from '../contexts/AuthContext';

const AdminRoute: ParentComponent = (props) => {
	const auth = useAuth();
	const navigate = useNavigate();

	// Admin pages need BOTH an admin account and an admin-enabled build. A
	// production build without VITE_ADMIN can never reach them, even for an
	// admin user.
	const allowed = () => FEATURES.admin && auth.isAdmin();

	createEffect(() => {
		if (!auth.isInitializing() && !allowed()) {
			navigate('/', { replace: true });
		}
	});

	// Show loading while initializing
	if (auth.isInitializing()) {
		return (
			<div class="min-h-screen flex items-center justify-center bg-background">
				<div class="text-center">
					<div class="text-lg text-foreground">Loading...</div>
				</div>
			</div>
		);
	}

	// we do not trust if the user is still here
	//
	if (!allowed()) {
		navigate('/', { replace: true });
		return null;
	}
	return <>{props.children}</>;
};

export default AdminRoute;
