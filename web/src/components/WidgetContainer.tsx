import { createSignal, For, type JSX, Suspense } from 'solid-js';
import { useWidgets } from '@/contexts/WidgetContext';
import AddWidgetModal from './AddWidgetModal';
import { Button, Input } from './atoms';
import RenderWidget from './RenderWidget';

interface WidgetContainerProps extends JSX.HTMLAttributes<HTMLDivElement> {}

export default function WidgetContainer(_props: WidgetContainerProps) {
	const { widgets } = useWidgets();
	const [isAddWidgetModalOpen, setIsAddWidgetModalOpen] = createSignal(false);
	const [search, setSearch] = createSignal('');

	const widgetsFiltered = () => {
		const searchInput = search().toLowerCase();
		return widgets().filter((widget) =>
			widget.schemaTitle?.toLowerCase().includes(searchInput),
		);
	};

	return (
		<div
			class="bg-background backdrop-blur-lg border border-white/20 w-full flex flex-col flex-1 m-2 mr-2 rounded-lg backdrop-blur-sm p-4 dark:border-slate-700/50 dark:shadow-black/30"
			id="widget-container"
		>
			<div class="mb-4 text-foreground flex justify-between items-center gap-4">
				<Input
					label="Search widgets"
					variant="glass"
					value={search()}
					onInput={(e) => setSearch(e.currentTarget.value)}
					placeholder="Search widgets..."
					class="py-2 flex-1"
					id="widget-search"
					title="Search widgets"
				/>

				<Button
					variant="primary"
					onClick={() => setIsAddWidgetModalOpen(true)}
					class="flex items-center gap-2"
				>
					<svg
						class="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4v16m8-8H4"
						/>
					</svg>
					Add Widget
				</Button>
			</div>

			<div class="flex-1 overflow-y-auto overflow-x-hidden treeview-container">
				<Suspense
					fallback={
						<div class="text-foreground text-center py-8">
							Loading widgets...
						</div>
					}
				>
					{widgetsFiltered().length > 0 ? (
						<div class="flex flex-col gap-4 w-full">
							<For each={widgetsFiltered()}>
								{(widget) => (
									<div class="w-full">
										<RenderWidget spacing={2} widget={widget} />
									</div>
								)}
							</For>
						</div>
					) : (
						<div class="bg-yellow-500/20 border-2 border-yellow-500/50 rounded-xl p-8 text-foreground text-center">
							<div class="text-4xl mb-4">📦</div>
							<h3 class="text-xl font-bold mb-2">No Widgets Found</h3>
							<p class="text-sm text-foreground/70">
								No widgets are currently assigned to your account.
							</p>
						</div>
					)}
				</Suspense>
			</div>

			<AddWidgetModal
				isOpen={isAddWidgetModalOpen()}
				onClose={() => setIsAddWidgetModalOpen(false)}
				widgetLength={widgets().length}
			/>
		</div>
	);
}
