import { For, Show } from 'solid-js';
import type { ServicesTaskDTO } from '@/api';
import AgendaRow from './AgendaRow';

interface AgendaListProps {
	tasks: ServicesTaskDTO[];
	loading: boolean;
	onTaskClick: (task: ServicesTaskDTO) => void;
}

export default function AgendaList(props: AgendaListProps) {
	return (
		<div class="flex-1 overflow-y-auto treeview-container">
			<Show
				when={!props.loading}
				fallback={
					<div class="text-foreground/60 text-center py-8">
						Loading tasks...
					</div>
				}
			>
				<Show
					when={props.tasks.length > 0}
					fallback={
						<div class="bg-yellow-500/20 border-2 border-yellow-500/50 rounded-xl p-8 text-foreground text-center">
							<div class="text-4xl mb-4">📋</div>
							<h3 class="text-xl font-bold mb-2">No Tasks Found</h3>
							<p class="text-sm text-foreground/70">
								Create a new task to get started.
							</p>
						</div>
					}
				>
					<div class="flex flex-col gap-2">
						<For each={props.tasks}>
							{(task) => (
								<AgendaRow
									task={task}
									onClick={() => props.onTaskClick(task)}
								/>
							)}
						</For>
					</div>
				</Show>
			</Show>
		</div>
	);
}
