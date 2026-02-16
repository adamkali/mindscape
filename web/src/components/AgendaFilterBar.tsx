import { type FilterType, useView } from '@/contexts/ViewContext';
import { Button } from './atoms';

const QUEUE_FILTERS: { char: string; label: string }[] = [
	{ char: 'a', label: 'Available' },
	{ char: 'c', label: 'Completed' },
	{ char: 'x', label: 'Cancelled' },
	{ char: 's', label: 'Scheduled' },
];

const STATUS_FILTERS: { char: string; label: string }[] = [
	{ char: 'a', label: 'Ambiguous' },
	{ char: 'c', label: 'Cancelled' },
	{ char: 'd', label: 'Done' },
	{ char: 'h', label: 'Hold' },
	{ char: 'p', label: 'Pending' },
	{ char: 'r', label: 'Recurring' },
	{ char: 'u', label: 'Undone' },
	{ char: 'i', label: 'Urgent' },
];

function isActive(current: FilterType, target: FilterType): boolean {
	if (current.kind !== target.kind) return false;
	if (current.kind === 'all') return true;
	if (current.kind === 'queue' && target.kind === 'queue')
		return current.char === target.char;
	if (current.kind === 'status' && target.kind === 'status')
		return current.char === target.char;
	return false;
}

export default function AgendaFilterBar() {
	const view = useView();

	return (
		<div class="flex items-center gap-2 flex-wrap flex-1 min-w-0">
			<Button
				variant={view.activeFilter().kind === 'all' ? 'primary' : 'tertiary'}
				onClick={() => view.setActiveFilter({ kind: 'all' })}
				class="text-xs !py-1 !px-2"
			>
				All
			</Button>

			<span class="text-foreground/30 text-xs">|</span>

			{QUEUE_FILTERS.map((q) => (
				<Button
					variant={
						isActive(view.activeFilter(), { kind: 'queue', char: q.char })
							? 'primary'
							: 'tertiary'
					}
					onClick={() => view.setActiveFilter({ kind: 'queue', char: q.char })}
					class="text-xs !py-1 !px-2"
				>
					{q.label}
				</Button>
			))}

			<span class="text-foreground/30 text-xs">|</span>

			{STATUS_FILTERS.map((s) => (
				<Button
					variant={
						isActive(view.activeFilter(), { kind: 'status', char: s.char })
							? 'primary'
							: 'tertiary'
					}
					onClick={() => view.setActiveFilter({ kind: 'status', char: s.char })}
					class="text-xs !py-1 !px-2"
				>
					{s.label}
				</Button>
			))}
		</div>
	);
}
