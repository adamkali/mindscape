import {
	createContext,
	createResource,
	type ParentComponent,
	useContext,
} from 'solid-js';
import type {
	AddUserWidgetRequest,
	ResponsesUserWidgetData,
	ResponsesWidgetData,
} from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { useAuthenticatedApi } from '@/utils/useApi';

export interface WidgetContextValue {
	widgets: () => ResponsesUserWidgetData[];
	widgetsLoading: () => boolean;
	schemas: () => ResponsesWidgetData[];
	schemasLoading: () => boolean;
	refreshWidgets: () => void;
	addWidget: (body: AddUserWidgetRequest) => Promise<boolean>;
}

const WidgetContext = createContext<WidgetContextValue>();

export const WidgetProvider: ParentComponent = (props) => {
	const auth = useAuth();
	const api = useAuthenticatedApi();

	// User's widgets, re-fetched whenever the auth token changes.
	const [widgets, { refetch }] = createResource(
		() => auth.token(),
		async (token) => {
			if (!token) return [];
			const response = await api.widgets.getUserWidgets({});
			if (response.success && response.data) return response.data;
			console.error('Failed to fetch widgets:', response.message);
			return [];
		},
	);

	// Available widget schemas — public, fetched once.
	const [schemas] = createResource(async () => {
		const response = await api.widgets.getWidgetSchemas();
		if (response.success && response.data) return response.data;
		console.error('Failed to fetch widget schemas:', response.message);
		return [];
	});

	const addWidget = async (body: AddUserWidgetRequest) => {
		const token = auth.token();
		if (!token) return false;
		try {
			const response = await api.widgets.addUserWidget({
				addUserWidgetRequest: body,
			});
			if (response.success) {
				refetch();
				return true;
			}
			console.error('Failed to add widget:', response.message);
		} catch (error) {
			console.error('Error adding widget:', error);
		}
		return false;
	};

	const value: WidgetContextValue = {
		widgets: () => widgets() ?? [],
		widgetsLoading: () => widgets.loading,
		schemas: () => schemas() ?? [],
		schemasLoading: () => schemas.loading,
		refreshWidgets: refetch,
		addWidget,
	};

	return (
		<WidgetContext.Provider value={value}>
			{props.children}
		</WidgetContext.Provider>
	);
};

export const useWidgets = (): WidgetContextValue => {
	const context = useContext(WidgetContext);
	if (!context) {
		throw new Error('useWidgets must be used within a WidgetProvider');
	}
	return context;
};
