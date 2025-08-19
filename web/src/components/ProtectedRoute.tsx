import { useNavigate } from '@solidjs/router';
import { createEffect, type ParentComponent } from 'solid-js';
import { useAuth } from '@/contexts/AuthContext';

const ProtectedRoute: ParentComponent = (props) => {
	const auth = useAuth();
	const navigate = useNavigate();

	createEffect(() => {
		// Only redirect if not initializing and not authenticated
		if (!auth.isInitializing() && !auth.isAuthenticated()) {
			navigate('/login', { replace: true });
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

	// Don't render if not authenticated after initialization
	if (!auth.isAuthenticated()) {
		return null;
	}

	return <>{props.children}</>;
};

export default ProtectedRoute;
