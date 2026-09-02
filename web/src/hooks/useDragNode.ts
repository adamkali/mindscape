import { createSignal } from 'solid-js';
import {
	type DragNodePayload,
	parseNode,
	serializeNode,
} from '@/utils/dragNode';

/**
 * Makes an element a draggable tree node. Spread `draggableProps` onto the
 * element/card; read `isDragging()` for styling. The payload is read lazily on
 * drag start so it always reflects current values.
 */
export function useDraggableNode(payload: () => DragNodePayload) {
	const [isDragging, setIsDragging] = createSignal(false);

	const draggableProps = {
		draggable: true,
		onDragStart: (e: DragEvent) => {
			serializeNode(e, payload());
			setIsDragging(true);
		},
		onDragEnd: () => setIsDragging(false),
	};

	return { isDragging, draggableProps };
}

/**
 * Makes an element a drop target for tree nodes. Spread `dropProps` onto the
 * element; read `isDragOver()` for styling. `onDrop` receives the parsed
 * payload (invalid drops are ignored).
 */
export function useDropTarget(onDrop: (payload: DragNodePayload) => void) {
	const [isDragOver, setIsDragOver] = createSignal(false);

	const dropProps = {
		onDragOver: (e: DragEvent) => {
			e.preventDefault();
			e.stopPropagation();
			if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
			setIsDragOver(true);
		},
		onDragLeave: () => setIsDragOver(false),
		onDrop: (e: DragEvent) => {
			e.preventDefault();
			e.stopPropagation();
			setIsDragOver(false);
			const payload = parseNode(e);
			if (payload) onDrop(payload);
		},
	};

	return { isDragOver, dropProps };
}
