import type { ServicesTaskDTO } from '@/api';

interface AgendaRowProps {
	task: ServicesTaskDTO;
	onClick: () => void;
}

const STATUS_ICONS: Record<string, { icon: string; color: string }> = {
	Done: { icon: '\u2713', color: 'text-green-400' },
	Urgent: { icon: '!', color: 'text-red-400' },
	Pending: { icon: '\u25CB', color: 'text-blue-400' },
	Hold: { icon: '\u2759\u2759', color: 'text-yellow-400' },
	Cancelled: { icon: '\u2717', color: 'text-red-400' },
	Recurring: { icon: '\u21BB', color: 'text-purple-400' },
	Ambiguous: { icon: '?', color: 'text-gray-400' },
	Undone: { icon: '\u25CB', color: 'text-foreground/40' },
};

function formatDate(dateStr: string): string {
	const date = new Date(dateStr);
	const now = new Date();
	const diffMs = date.getTime() - now.getTime();
	const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24));

	if (diffDays === 0) return 'Today';
	if (diffDays === 1) return 'Tomorrow';
	if (diffDays === -1) return 'Yesterday';
	if (diffDays > 0 && diffDays <= 7) return `In ${diffDays}d`;
	if (diffDays < 0 && diffDays >= -7) return `${Math.abs(diffDays)}d ago`;

	return date.toLocaleDateString('en-US', {
		month: 'short',
		day: 'numeric',
	});
}

export default function AgendaRow(props: AgendaRowProps) {
	const statusName = () =>
		(props.task.taskType?.name || 'Pending').replace('TaskStatus', '');
	const statusInfo = () => STATUS_ICONS[statusName()] || STATUS_ICONS.Pending;

	const dateDisplay = () => {
		if (props.task.dueAt) {
			return { date: formatDate(props.task.dueAt), icon: 'due' };
		}
		if (props.task.createdAt) {
			return { date: formatDate(props.task.createdAt), icon: 'created' };
		}
		return null;
	};

	return (
		<button
			type="button"
			onClick={props.onClick}
			class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg bg-glass-bg/30 border border-glass-border/40 hover:bg-glass-bg-hover hover:border-glass-border-hover transition-all duration-200 cursor-pointer text-left"
		>
			{/* Status Icon */}
			<span
				class={`text-lg font-bold w-6 text-center flex-shrink-0 ${statusInfo().color}`}
			>
				{statusInfo().icon}
			</span>

			{/* Date */}
			<span class="text-xs text-foreground/50 w-20 flex-shrink-0 flex items-center gap-1">
				{dateDisplay() && (
					<>
						{dateDisplay()?.icon === 'due' ? (
							<svg
								class="w-3 h-3"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								role="img"
								aria-label="Due date"
							>
								<circle cx="12" cy="12" r="10" />
								<path d="M12 6v6l4 2" />
							</svg>
						) : (
							<svg
								class="w-3 h-3"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								role="img"
								aria-label="Created date"
							>
								<rect x="3" y="4" width="18" height="18" rx="2" />
								<path d="M16 2v4M8 2v4M3 10h18" />
							</svg>
						)}
						{dateDisplay()?.date}
					</>
				)}
			</span>

			{/* Title */}
			<span class="text-sm text-foreground truncate flex-1">
				{props.task.name || 'Untitled Task'}
			</span>
		</button>
	);
}
