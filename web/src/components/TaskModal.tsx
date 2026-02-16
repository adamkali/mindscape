import { createSignal, Show } from 'solid-js';
import type { ServicesTaskDTO } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { useView } from '@/contexts/ViewContext';
import type { ModalMode } from './AgendaContainer';
import { Button } from './atoms';

interface TaskModalProps {
	mode: ModalMode;
	task?: ServicesTaskDTO;
	onClose: () => void;
	onEdit: () => void;
	onView: () => void;
}

const TASK_STATUSES = [
	{ char: 'a', label: 'Ambiguous' },
	{ char: 'c', label: 'Cancelled' },
	{ char: 'd', label: 'Done' },
	{ char: 'h', label: 'Hold' },
	{ char: 'p', label: 'Pending' },
	{ char: 'r', label: 'Recurring' },
	{ char: 'u', label: 'Undone' },
	{ char: 'i', label: 'Urgent' },
];

export default function TaskModal(props: TaskModalProps) {
	const view = useView();
	const auth = useAuth();
	const user = auth.user();

	const [name, setName] = createSignal('');
	const [description, setDescription] = createSignal('');
	const [showStatusDropdown, setShowStatusDropdown] = createSignal(false);
	const [showDeleteConfirm, setShowDeleteConfirm] = createSignal(false);

	const resetForm = () => {
		setName('');
		setDescription('');
		setShowStatusDropdown(false);
		setShowDeleteConfirm(false);
	};

	const handleBackdropClick = (e: MouseEvent) => {
		if (e.target === e.currentTarget) {
			resetForm();
			props.onClose();
		}
	};

	const handleCreate = async (e: Event) => {
		e.preventDefault();
		if (!name().trim()) return;
		await view.createTask({
			name: name().trim(),
			description: description().trim(),
			userId: user?.id,
		});
		resetForm();
		props.onClose();
	};

	const handleEdit = async (e: Event) => {
		e.preventDefault();
		if (!props.task?.id || !name().trim()) return;
		await view.updateTaskContent({
			id: props.task.id,
			name: name().trim(),
			description: description().trim(),
		});
		resetForm();
		props.onClose();
	};

	const handleStatusChange = async (statusChar: string) => {
		if (!props.task?.id) return;
		await view.updateTaskStatus(props.task.id, statusChar);
		setShowStatusDropdown(false);
		props.onClose();
	};

	const handleDelete = async () => {
		if (!props.task?.id) return;
		await view.deleteTask(props.task.id);
		setShowDeleteConfirm(false);
		resetForm();
		props.onClose();
	};

	const startEdit = () => {
		setName(props.task?.name || '');
		setDescription(props.task?.description || '');
		props.onEdit();
	};

	return (
		<Show when={props.mode !== 'closed'}>
			<div
				class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50"
				role="dialog"
				aria-modal="true"
				onClick={handleBackdropClick}
				onKeyDown={(e) => e.key === 'Escape' && props.onClose()}
			>
				<div
					class="bg-gradient-to-br from-card to-card/80 backdrop-blur-lg border-2 border-slate-700/20 rounded-2xl shadow-2xl w-full max-w-md p-6"
					role="document"
				>
					{/* Header */}
					<div class="flex justify-between items-center mb-6">
						<h2 class="text-2xl font-bold text-card-foreground">
							{props.mode === 'create'
								? 'New Task'
								: props.mode === 'edit'
									? 'Edit Task'
									: props.task?.name || 'Task'}
						</h2>
						<button
							type="button"
							onClick={() => {
								resetForm();
								props.onClose();
							}}
							class="text-card-foreground/60 hover:text-card-foreground transition-colors"
						>
							<svg
								class="w-6 h-6"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
								role="img"
								aria-label="Close"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
						</button>
					</div>

					{/* Create Mode */}
					<Show when={props.mode === 'create'}>
						<form onSubmit={handleCreate} class="space-y-4">
							<div>
								<label
									for="task-name"
									class="block text-sm font-semibold text-card-foreground mb-2"
								>
									Name
								</label>
								<input
									id="task-name"
									type="text"
									class="w-full px-4 py-3 rounded-lg border-2 bg-slate-300/30 text-foreground focus:outline-none transition-all"
									value={name()}
									onInput={(e) => setName(e.currentTarget.value)}
									placeholder="Task name..."
								/>
							</div>
							<div>
								<label
									for="task-description"
									class="block text-sm font-semibold text-card-foreground mb-2"
								>
									Description
								</label>
								<textarea
									id="task-description"
									class="w-full px-4 py-3 rounded-lg border-2 bg-slate-300/30 text-foreground focus:outline-none transition-all min-h-24 resize-y"
									value={description()}
									onInput={(e) => setDescription(e.currentTarget.value)}
									placeholder="Task description..."
								/>
							</div>
							<div class="flex gap-3 pt-4">
								<Button
									type="button"
									variant="tertiary"
									onClick={() => {
										resetForm();
										props.onClose();
									}}
									class="flex-1"
								>
									Cancel
								</Button>
								<Button
									type="submit"
									variant="primary"
									class="flex-1"
									disabled={!name().trim()}
								>
									Create Task
								</Button>
							</div>
						</form>
					</Show>

					{/* Edit Mode */}
					<Show when={props.mode === 'edit'}>
						<form onSubmit={handleEdit} class="space-y-4">
							<div>
								<label
									for="task-name-edit"
									class="block text-sm font-semibold text-card-foreground mb-2"
								>
									Name
								</label>
								<input
									id="task-name-edit"
									type="text"
									class="w-full px-4 py-3 rounded-lg border-2 bg-slate-300/30 text-foreground focus:outline-none transition-all"
									value={name()}
									onInput={(e) => setName(e.currentTarget.value)}
									placeholder="Task name..."
								/>
							</div>
							<div>
								<label
									for="task-description-edit"
									class="block text-sm font-semibold text-card-foreground mb-2"
								>
									Description
								</label>
								<textarea
									id="task-description-edit"
									class="w-full px-4 py-3 rounded-lg border-2 bg-slate-300/30 text-foreground focus:outline-none transition-all min-h-24 resize-y"
									value={description()}
									onInput={(e) => setDescription(e.currentTarget.value)}
									placeholder="Task description..."
								/>
							</div>
							<div class="flex gap-3 pt-4">
								<Button
									type="button"
									variant="tertiary"
									onClick={() => props.onView()}
									class="flex-1"
								>
									Cancel
								</Button>
								<Button
									type="submit"
									variant="primary"
									class="flex-1"
									disabled={!name().trim()}
								>
									Save
								</Button>
							</div>
						</form>
					</Show>

					{/* View Mode */}
					<Show when={props.mode === 'view' && props.task}>
						<div class="space-y-4">
							<div>
								<div class="text-xs text-card-foreground/50 uppercase tracking-wider mb-1">
									Status
								</div>
								<div class="text-sm text-card-foreground">
									{props.task?.taskType?.name || 'Pending'}
								</div>
							</div>

							<Show when={props.task?.description}>
								<div>
									<div class="text-xs text-card-foreground/50 uppercase tracking-wider mb-1">
										Description
									</div>
									<div class="text-sm text-card-foreground whitespace-pre-wrap">
										{props.task?.description}
									</div>
								</div>
							</Show>

							<Show when={props.task?.dueAt}>
								<div>
									<div class="text-xs text-card-foreground/50 uppercase tracking-wider mb-1">
										Due
									</div>
									<div class="text-sm text-card-foreground">
										{props.task?.dueAt &&
											new Date(props.task.dueAt).toLocaleString()}
									</div>
								</div>
							</Show>

							<div class="flex gap-2">
								<Show when={props.task?.createdAt}>
									<div class="flex-1">
										<div class="text-xs text-card-foreground/50 uppercase tracking-wider mb-1">
											Created
										</div>
										<div class="text-xs text-card-foreground/70">
											{props.task?.createdAt &&
												new Date(props.task.createdAt).toLocaleString()}
										</div>
									</div>
								</Show>
								<Show when={props.task?.updatedAt}>
									<div class="flex-1">
										<div class="text-xs text-card-foreground/50 uppercase tracking-wider mb-1">
											Updated
										</div>
										<div class="text-xs text-card-foreground/70">
											{props.task?.updatedAt &&
												new Date(props.task.updatedAt).toLocaleString()}
										</div>
									</div>
								</Show>
							</div>

							{/* Footer Actions */}
							<div class="flex gap-2 pt-4 border-t border-white/10">
								<Button variant="primary" onClick={startEdit} class="flex-1">
									Edit
								</Button>

								{/* Status Dropdown */}
								<div class="relative flex-1">
									<Button
										variant="secondary"
										onClick={() => setShowStatusDropdown(!showStatusDropdown())}
										class="w-full"
									>
										Change Status
									</Button>
									<Show when={showStatusDropdown()}>
										<div class="absolute bottom-full mb-1 left-0 w-full bg-glass-bg-strong backdrop-blur-md border border-white/30 rounded-xl shadow-lg z-10 dark:shadow-black/30">
											<div class="py-1">
												{TASK_STATUSES.map((s) => (
													<button
														type="button"
														onClick={() => handleStatusChange(s.char)}
														class="block w-full text-left px-4 py-2 text-sm text-foreground hover:bg-glass-bg-hover transition-colors"
													>
														{s.label}
													</button>
												))}
											</div>
										</div>
									</Show>
								</div>

								{/* Delete */}
								<Show
									when={showDeleteConfirm()}
									fallback={
										<Button
											variant="danger"
											onClick={() => setShowDeleteConfirm(true)}
											class="flex-shrink-0"
										>
											Delete
										</Button>
									}
								>
									<div class="flex gap-1 flex-shrink-0">
										<Button variant="danger" onClick={handleDelete}>
											Confirm
										</Button>
										<Button
											variant="tertiary"
											onClick={() => setShowDeleteConfirm(false)}
										>
											No
										</Button>
									</div>
								</Show>
							</div>
						</div>
					</Show>
				</div>
			</div>
		</Show>
	);
}
