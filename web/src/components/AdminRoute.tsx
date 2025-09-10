import { useNavigate } from "@solidjs/router";
import { createEffect, type ParentComponent } from "solid-js";
import { useAuth } from "../contexts/AuthContext";

const AdminRoute: ParentComponent = (props) => {
	const auth = useAuth();
	const navigate = useNavigate();

	createEffect(() => {
		if (!auth.isInitializing() && !auth.isAdmin()) {
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
	if (!auth.isAdmin()) {
		navigate('/', { replace: true });
	}
	return <>{props.children}</>;
}

export default AdminRoute
