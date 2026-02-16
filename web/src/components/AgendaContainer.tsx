import { createSignal } from 'solid-js';
import type { ServicesTaskDTO } from '@/api';
import { useView } from '@/contexts/ViewContext';
import AgendaFilterBar from './AgendaFilterBar';
import AgendaList from './AgendaList';
import { Button } from './atoms';
import TaskModal from './TaskModal';

export type ModalMode = 'closed' | 'create' | 'view' | 'edit';

export default function AgendaContainer() {
	const view = useView();
	const [modalMode, setModalMode] = createSignal<ModalMode>('closed');
	const [selectedTask, setSelectedTask] = createSignal<
		ServicesTaskDTO | undefined
	>();

	const handleTaskClick = (task: ServicesTaskDTO) => {
		setSelectedTask(task);
		setModalMode('view');
	};

	const handleNewTask = () => {
		setSelectedTask(undefined);
		setModalMode('create');
	};

	const closeModal = () => {
		setModalMode('closed');
		setSelectedTask(undefined);
	};

	return (
		<div
			class="bg-background backdrop-blur-lg border border-white/20 w-full flex flex-col m-2 rounded-lg p-4 max-h-[calc(100vh-2rem)] dark:border-slate-700/50 dark:shadow-black/30"
			id="agenda-container"
		>
			<div class="mb-4 text-foreground flex justify-between items-center gap-4">
				<AgendaFilterBar />
				<Button
					variant="primary"
					onClick={handleNewTask}
					class="flex items-center gap-2 flex-shrink-0"
				>
					<svg
						class="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						role="img"
						aria-label="Add"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4v16m8-8H4"
						/>
					</svg>
					New Task
				</Button>
			</div>

			<AgendaList
				tasks={view.tasks()}
				loading={view.tasksLoading()}
				onTaskClick={handleTaskClick}
			/>

			<TaskModal
				mode={modalMode()}
				task={selectedTask()}
				onClose={closeModal}
				onEdit={() => setModalMode('edit')}
				onView={() => setModalMode('view')}
			/>
		</div>
	);
}
