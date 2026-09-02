import {
	createContext,
	createEffect,
	createSignal,
	type ParentComponent,
	useContext,
} from 'solid-js';
import {
	DefaultApi,
	type RepositoryInsertNewTaskParams,
	type RepositoryUpdateTaskContentParams,
	type ResponsesTasksResponse,
	type ServicesTaskDTO,
} from '@/api';
import { FEATURES } from '@/config/features';
import { getAuthenticatedApiConfig } from '@/utils/apiConfig';
import { useAuth } from './AuthContext';

export type ActiveView = 'widgets' | 'agenda';

/**
 * Outer page layout, orthogonal to `activeView`:
 *  - 'normal' — tree in a fixed-width left rail, widgets/agenda fill the rest
 *  - 'tree'   — tree takes the full width, right pane hidden
 *
 * `activeView` only matters inside 'normal'; in 'tree' the right pane is hidden
 * but never unmounted, so widget state survives a round trip.
 */
export type HomepageMode = 'tree' | 'normal';

const HOMEPAGE_MODE_KEY = 'homepageMode';

const isHomepageMode = (v: string | null): v is HomepageMode =>
	v === 'tree' || v === 'normal';

export type FilterType =
	| { kind: 'all' }
	| { kind: 'queue'; char: string }
	| { kind: 'status'; char: string };

export interface ViewContextValue {
	activeView: () => ActiveView;
	setActiveView: (view: ActiveView) => void;
	homepageMode: () => HomepageMode;
	setHomepageMode: (mode: HomepageMode) => void;
	tasks: () => ServicesTaskDTO[];
	tasksLoading: () => boolean;
	activeFilter: () => FilterType;
	setActiveFilter: (filter: FilterType) => void;
	refreshTasks: () => Promise<void>;
	createTask: (params: RepositoryInsertNewTaskParams) => Promise<void>;
	updateTaskContent: (
		params: RepositoryUpdateTaskContentParams,
	) => Promise<void>;
	updateTaskStatus: (
		taskId: string,
		status: string,
		dueDate?: string,
	) => Promise<void>;
	deleteTask: (taskId: string) => Promise<void>;
}

const ViewContext = createContext<ViewContextValue>();

export const ViewProvider: ParentComponent = (props) => {
	const auth = useAuth();
	const tasksApi = new DefaultApi(getAuthenticatedApiConfig());

	const [activeView, setActiveViewSignal] = createSignal<ActiveView>('widgets');

	// Guard the setter rather than the callers: with tasks off, 'agenda' is not
	// a reachable state no matter who asks for it.
	const setActiveView = (view: ActiveView) => {
		if (view === 'agenda' && !FEATURES.tasks) return;
		setActiveViewSignal(view);
	};

	// Restored from localStorage so a reload keeps the last layout; 'normal'
	// is the default for a first visit or an unreadable/legacy value.
	const [homepageMode, setHomepageModeSignal] = createSignal<HomepageMode>(
		(() => {
			try {
				const stored = localStorage.getItem(HOMEPAGE_MODE_KEY);
				return isHomepageMode(stored) ? stored : 'normal';
			} catch {
				return 'normal';
			}
		})(),
	);

	const setHomepageMode = (mode: HomepageMode) => {
		setHomepageModeSignal(mode);
		try {
			localStorage.setItem(HOMEPAGE_MODE_KEY, mode);
		} catch {
			// A blocked/full storage must not break layout switching.
		}
	};
	const [tasks, setTasks] = createSignal<ServicesTaskDTO[]>([]);
	const [tasksLoading, setTasksLoading] = createSignal(false);
	const [tasksLoaded, setTasksLoaded] = createSignal(false);
	const [activeFilter, setActiveFilter] = createSignal<FilterType>({
		kind: 'all',
	});

	createEffect(() => {
		// With tasks flagged off the agenda pane is never rendered, so the fetch
		// must not fire either — its endpoints aren't mounted.
		if (!FEATURES.tasks) return;
		if (activeView() === 'agenda' && !tasksLoaded()) {
			fetchTasks();
		}
	});

	const fetchTasks = async () => {
		if (!auth.token()) return;
		setTasksLoading(true);
		try {
			const filter = activeFilter();
			let response: ResponsesTasksResponse;
			if (filter.kind === 'queue') {
				response = await tasksApi.getTasksByQueueType({
					queueType: filter.char,
				});
			} else if (filter.kind === 'status') {
				response = await tasksApi.getTasksByTaskType({
					taskType: filter.char,
				});
			} else {
				response = await tasksApi.readTasks({});
			}
			if (response.success && response.data) {
				setTasks(response.data);
			}
			setTasksLoaded(true);
		} catch (error) {
			console.error('Failed to fetch tasks:', error);
		} finally {
			setTasksLoading(false);
		}
	};

	const refreshTasks = async () => {
		setTasksLoaded(false);
		await fetchTasks();
	};

	const applyFilter = async (filter: FilterType) => {
		setActiveFilter(filter);
		setTasksLoaded(false);
		await fetchTasks();
	};

	const createTask = async (params: RepositoryInsertNewTaskParams) => {
		if (!auth.token()) return;
		try {
			const response = await tasksApi.createTask({
				createTaskRequest: params,
			});
			if (response.success && response.data) {
				const newTask = response.data;
				setTasks((prev) => [...prev, newTask]);
			}
		} catch (error) {
			console.error('Failed to create task:', error);
		}
	};

	const updateTaskContent = async (
		params: RepositoryUpdateTaskContentParams,
	) => {
		if (!auth.token()) return;
		try {
			const response = await tasksApi.updateTask({
				updateTaskRequest: params,
			});
			if (response.success && response.data) {
				const updated = response.data;
				setTasks((prev) =>
					prev.map((t) => (t.id === updated.id ? updated : t)),
				);
			}
		} catch (error) {
			console.error('Failed to update task:', error);
		}
	};

	const updateTaskStatus = async (
		taskId: string,
		status: string,
		dueDate?: string,
	) => {
		if (!auth.token()) return;
		try {
			const response = await tasksApi.updateTaskStatus({
				taskId,
				status,
				dueDate,
			});
			if (response.success && response.data) {
				const updated = response.data;
				setTasks((prev) =>
					prev.map((t) => (t.id === updated.id ? updated : t)),
				);
			}
		} catch (error) {
			console.error('Failed to update task status:', error);
		}
	};

	const deleteTask = async (taskId: string) => {
		if (!auth.token()) return;
		try {
			const response = await tasksApi.deleteTask({
				taskId,
			});
			if (response.success) {
				setTasks((prev) => prev.filter((t) => t.id !== taskId));
			}
		} catch (error) {
			console.error('Failed to delete task:', error);
		}
	};

	const value: ViewContextValue = {
		activeView,
		setActiveView,
		homepageMode,
		setHomepageMode,
		tasks,
		tasksLoading,
		activeFilter,
		setActiveFilter: applyFilter,
		refreshTasks,
		createTask,
		updateTaskContent,
		updateTaskStatus,
		deleteTask,
	};

	return (
		<ViewContext.Provider value={value}>{props.children}</ViewContext.Provider>
	);
};

export const useView = () => {
	const context = useContext(ViewContext);
	if (!context) {
		throw new Error('useView must be used within a ViewProvider');
	}
	return context;
};
