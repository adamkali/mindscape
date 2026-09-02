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
	{
		char: 'a',
		label: 'Ambiguous',
		id: 'e56fd149-24de-4835-9dad-ae861a7c3155',
		color: '#f59e0b',
	},
	{
		char: 'c',
		label: 'Cancelled',
		id: '07bae843-7049-449c-a23e-ab78a571d7ca',
		color: '#6b7280',
	},
	{
		char: 'd',
		label: 'Done',
		id: '546d40a2-aebd-4c3e-b1b3-3fd835211c74',
		color: '#22c55e',
	},
	{
		char: 'h',
		label: 'Hold',
		id: '56fabcc6-9703-43b5-96fc-2876646a26b9',
		color: '#eab308',
	},
	{
		char: 'p',
		label: 'Pending',
		id: '11360cdc-f811-425f-b565-8b014c45ec25',
		color: '#3b82f6',
	},
	{
		char: 'r',
		label: 'Recurring',
		id: 'f1559502-fa64-419b-9b55-57842e1af279',
		color: '#a855f7',
	},
	{
		char: 'u',
		label: 'Undone',
		id: '99dee5b2-7ac9-4b02-a3e5-a1c917d90009',
		color: '#e2e8f0',
	},
	{
		char: 'i',
		label: 'Urgent',
		id: '106e703a-4dd4-4737-b38b-e4a0000ff158',
		color: '#ef4444',
	},
];

export default function TaskModal(props: TaskModalProps) {
	const view = useView();
	const auth = useAuth();
	const user = auth.user();

	const [name, setName] = createSignal('');
	const [description, setDescription] = createSignal('');
	const [taskTypeId, setTaskTypeId] = createSignal(TASK_STATUSES[4].id); // default: Pending
	const [dueAt, setDueAt] = createSignal('');
	const [showStatusDropdown, setShowStatusDropdown] = createSignal(false);
	const [showDeleteConfirm, setShowDeleteConfirm] = createSignal(false);

	const statusColor = () => {
		const id = props.task?.taskTypeId || taskTypeId();
		return TASK_STATUSES.find((s) => s.id === id)?.color ?? '#3b82f6';
	};

	const resetForm = () => {
		setName('');
		setDescription('');
		setTaskTypeId(TASK_STATUSES[4].id);
		setDueAt('');
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
			taskTypeId: taskTypeId(),
			dueAt: dueAt() ? new Date(dueAt()).toISOString() : undefined,
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
			dueAt: dueAt() ? new Date(dueAt()).toISOString() : undefined,
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
		const match = TASK_STATUSES.find((s) => s.id === props.task?.taskTypeId);
		setTaskTypeId(match?.id || TASK_STATUSES[4].id);
		setDueAt(props.task?.dueAt ? new Date(props.task.dueAt).toISOString().slice(0, 16) : '');
		props.onEdit();
	};

	return (
		<Show when={props.mode !== 'closed'}>
			<div
				class="absolute inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50 rounded-lg"
				role="dialog"
				aria-modal="true"
				onClick={handleBackdropClick}
				onKeyDown={(e) => e.key === 'Escape' && props.onClose()}
			>
				<div
					class="bg-glass-bg-strong backdrop-blur-lg border-2 rounded-2xl shadow-2xl w-full max-w-md mx-4 p-6"
					style={{ 'border-color': statusColor() }}
					role="document"
				>
					{/* Header */}
					<div class="flex justify-between items-center mb-6">
						<h2 class="text-2xl font-bold" style={{ color: statusColor() }}>
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
								<div class="flex rounded-lg border-2 border-glass-border bg-slate-300/30 overflow-hidden focus-within:border-glass-border-hover transition-all">
									<select
										value={taskTypeId()}
										onChange={(e) => setTaskTypeId(e.currentTarget.value)}
										class="px-3 py-3 bg-glass-bg text-foreground text-sm font-medium border-r border-glass-border cursor-pointer focus:outline-none appearance-none"
										aria-label="Task status type"
									>
										{TASK_STATUSES.map((s) => (
											<option value={s.id}>{s.label}</option>
										))}
									</select>
									<input
										id="task-name"
										type="text"
										class="flex-1 px-4 py-3 bg-transparent text-foreground focus:outline-none"
										value={name()}
										onInput={(e) => setName(e.currentTarget.value)}
										placeholder="Task name..."
									/>
								</div>
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
							<div>
								<label
									for="task-due-at"
									class="block text-sm font-semibold text-card-foreground mb-2"
								>
									Due Date
								</label>
								<input
									id="task-due-at"
									type="datetime-local"
									class="w-full px-4 py-3 rounded-lg border-2 bg-slate-300/30 text-foreground focus:outline-none transition-all"
									value={dueAt()}
									onInput={(e) => setDueAt(e.currentTarget.value)}
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
								<div class="flex rounded-lg border-2 border-glass-border bg-slate-300/30 overflow-hidden focus-within:border-glass-border-hover transition-all">
									<select
										value={taskTypeId()}
										onChange={(e) => setTaskTypeId(e.currentTarget.value)}
										class="px-3 py-3 bg-glass-bg text-foreground text-sm font-medium border-r border-glass-border cursor-pointer focus:outline-none appearance-none"
										aria-label="Task status type"
									>
										{TASK_STATUSES.map((s) => (
											<option value={s.id}>{s.label}</option>
										))}
									</select>
									<input
										id="task-name-edit"
										type="text"
										class="flex-1 px-4 py-3 bg-transparent text-foreground focus:outline-none"
										value={name()}
										onInput={(e) => setName(e.currentTarget.value)}
										placeholder="Task name..."
									/>
								</div>
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
							<div>
								<label
									for="task-due-at-edit"
									class="block text-sm font-semibold text-card-foreground mb-2"
								>
									Due Date
								</label>
								<input
									id="task-due-at-edit"
									type="datetime-local"
									class="w-full px-4 py-3 rounded-lg border-2 bg-slate-300/30 text-foreground focus:outline-none transition-all"
									value={dueAt()}
									onInput={(e) => setDueAt(e.currentTarget.value)}
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
								<div
									class="flex items-center gap-2 text-sm font-medium"
									style={{ color: statusColor() }}
								>
									<span
										class="inline-block w-2.5 h-2.5 rounded-full"
										style={{ 'background-color': statusColor() }}
									/>
									{(props.task?.taskType?.name || 'Pending').replace('TaskStatus', '')}
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
							<div class="flex gap-2 pt-4 border-t border-glass-border">
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
														class="flex items-center gap-2 w-full text-left px-4 py-2 text-sm text-foreground hover:bg-glass-bg-hover transition-colors"
													>
														<span
															class="inline-block w-2.5 h-2.5 rounded-full flex-shrink-0"
															style={{ 'background-color': s.color }}
														/>
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
