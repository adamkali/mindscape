import {
	createEffect,
	createSignal,
	For,
	type JSX,
	onMount,
	Suspense,
} from 'solid-js';
import { type ResponsesUserWidgetData, WidgetsApi } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import AddWidgetModal from './AddWidgetModal';
import { Button } from './atoms';
import RenderWidget from './RenderWidget';
import type { ElementSize } from './types';

interface WidgetContainerProps extends JSX.HTMLAttributes<HTMLDivElement> {}

const DEFAULT_COLUMN_COUNT = 24;
const REM = 2;

function calculateBestColumnSize(elementSize: ElementSize) {
	// get the best possible size of the column based on the width
	// of the container and the number of columns
	const columnCount = DEFAULT_COLUMN_COUNT;
	const squareSize = Math.floor(elementSize.width / columnCount);
	return squareSize;
}

function calculateSquareElementSize(elementSize: ElementSize) {
	const width = elementSize.width;
	const columnCount = DEFAULT_COLUMN_COUNT;
	// Get current rem size
	const remSize = getComputedStyle(document.documentElement).fontSize;

	// calculate width minus 2rem padding
	const newWidthMinusPadding = width - 2 * REM * parseFloat(remSize);
	const squareSize = newWidthMinusPadding / columnCount;
	const squareElementSize = {
		width: squareSize,
		height: squareSize,
	};
	return squareElementSize;
}

function roundDown(size: ElementSize) {
	return {
		width: Math.floor(size.width),
		height: Math.floor(size.height),
	};
}

export default function WidgetContainer(props: WidgetContainerProps) {
	const [elementSize, setElementSize] = createSignal<ElementSize>({
		width: 0,
		height: 0,
	});
	const [widgets, setWidgets] = createSignal<ResponsesUserWidgetData[]>([]);
	const [firstWidget, setFirstWidget] = createSignal<ResponsesUserWidgetData>();
	const [isAddWidgetModalOpen, setIsAddWidgetModalOpen] = createSignal(false);
	const auth = useAuth();
	const token = auth.token();

	const getCurrentSizeOfContainer = () => {
		const container = window.document.querySelector('#widget-container');

		if (container) {
			const rect = container.getBoundingClientRect();
			return {
				width: rect.width,
				height: rect.height,
			};
		} else {
			return undefined;
		}
	};

	const getColumnCount = () => {
		const currentSize = elementSize();
		if (currentSize) {
			return calculateBestColumnSize(currentSize);
		}
	};

	const handleResize = () => {
		const newSize = getCurrentSizeOfContainer();
		if (newSize !== undefined) {
			if (
				newSize.width === elementSize().width &&
				newSize.height === elementSize().height
			) {
				return;
			}
			setElementSize(newSize);
		}
	};

	window.addEventListener('resize', handleResize);

	// set initial size
	createEffect(() => {
		const newSize = getCurrentSizeOfContainer();
		if (newSize !== undefined) {
			setElementSize(newSize);
		}
	});

	onMount(() => {
		handleResize();
		getWidget().then((widgets) => {
			setWidgets(widgets);
		});
	});

	createEffect(() => {
		setFirstWidget(widgets()[0]);
		console.log('first widget', firstWidget());
	});

	const squareSize = (): ElementSize =>
		roundDown(calculateSquareElementSize(elementSize()));

	const dumbList = () => {
		return [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12];
	};

	const getWidget = async () => {
		const api = new WidgetsApi();
		const widgets = await api.getUserWidgets({
			authorization: 'Bearer ' + token,
		});
		if (widgets.data && widgets.success) {
			return widgets.data;
		} else {
			console.error({ error: widgets.message });
		}
		return [];
	};

	const handleWidgetAdded = async () => {
		// Refresh widgets list after adding a new widget
		const updatedWidgets = await getWidget();
		setWidgets(updatedWidgets);
	};

	return (
		<div
			class="bg-slate-700/50 h-full w-full flex flex-col m-2 rounded-3xl backdrop-blur-sm p-4"
			id="widget-container"
		>
			<div class="mb-4 text-white flex justify-between items-start">
				<div>
					<h2 class="text-2xl font-bold mb-2">Widget Debug Panel</h2>
					<div class="text-sm text-white/70">
						Total widgets loaded:{' '}
						<span class="font-bold text-green-300">{widgets().length}</span>
					</div>
					<div class="text-xs text-white/50 mt-1">
						Grid size: {squareSize().width.toFixed(0)}px ×{' '}
						{squareSize().height.toFixed(0)}px
					</div>
				</div>

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

			<div class="flex-1 overflow-y-auto overflow-x-hidden">
				<Suspense
					fallback={
						<div class="text-white text-center py-8">Loading widgets...</div>
					}
				>
					{widgets().length > 0 ? (
						<div
							class="grid w-full"
							style={{
								'grid-template-columns': `repeat(${DEFAULT_COLUMN_COUNT}, 1fr)`,
								'grid-auto-rows': `${squareSize().height}px`,
								gap: '0px',
							}}
						>
							<For each={widgets()}>
								{(widget) => (
									<div
										style={{
											'grid-column': `${(widget.positionX ?? 0) + 1} / span ${widget.width ?? 1}`,
											'grid-row': `${(widget.positionY ?? 0) + 1} / span ${widget.height ?? 1}`,
											'z-index': widget.zIndex ?? 1,
										}}
									>
										<RenderWidget
											elementSize={squareSize()}
											spacing={2}
											widget={widget}
										/>
									</div>
								)}
							</For>
						</div>
					) : (
						<div class="bg-yellow-500/20 border-2 border-yellow-500/50 rounded-xl p-8 text-white text-center">
							<div class="text-4xl mb-4">📦</div>
							<h3 class="text-xl font-bold mb-2">No Widgets Found</h3>
							<p class="text-sm text-white/70">
								No widgets are currently assigned to your account.
							</p>
						</div>
					)}
				</Suspense>
			</div>

			<AddWidgetModal
				isOpen={isAddWidgetModalOpen()}
				onClose={() => setIsAddWidgetModalOpen(false)}
				onWidgetAdded={handleWidgetAdded}
			/>
		</div>
	);
}
