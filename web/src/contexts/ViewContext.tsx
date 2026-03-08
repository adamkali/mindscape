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
import { useAuth } from './AuthContext';

export type ActiveView = 'widgets' | 'agenda';

export type FilterType =
	| { kind: 'all' }
	| { kind: 'queue'; char: string }
	| { kind: 'status'; char: string };

export interface ViewContextValue {
	activeView: () => ActiveView;
	setActiveView: (view: ActiveView) => void;
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
	const tasksApi = new DefaultApi();

	const [activeView, setActiveView] = createSignal<ActiveView>('widgets');
	const [tasks, setTasks] = createSignal<ServicesTaskDTO[]>([]);
	const [tasksLoading, setTasksLoading] = createSignal(false);
	const [tasksLoaded, setTasksLoaded] = createSignal(false);
	const [activeFilter, setActiveFilter] = createSignal<FilterType>({
		kind: 'all',
	});

	createEffect(() => {
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
					authorization: `Bearer ${auth.token()}`,
					queueType: filter.char,
				});
			} else if (filter.kind === 'status') {
				response = await tasksApi.getTasksByTaskType({
					authorization: `Bearer ${auth.token()}`,
					taskType: filter.char,
				});
			} else {
				response = await tasksApi.readTasks({
					authorization: `Bearer ${auth.token()}`,
				});
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
				authorization: `Bearer ${auth.token()}`,
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
				authorization: `Bearer ${auth.token()}`,
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
				authorization: `Bearer ${auth.token()}`,
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
				authorization: `Bearer ${auth.token()}`,
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
