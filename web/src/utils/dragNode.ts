/**
 * The drag-and-drop contract between tree cards and the tree mutation layer.
 * Cards serialize this payload on drag start; drop targets parse it and hand
 * the id off to `TreeContext.move`.
 */
export type DragNodeType = 'folder' | 'bookmark';

export interface DragNodePayload {
	type: DragNodeType;
	id: string;
	name?: string;
	link?: string;
}

const DRAG_MIME = 'text/plain';

export const serializeNode = (e: DragEvent, payload: DragNodePayload): void => {
	if (!e.dataTransfer) return;
	e.dataTransfer.setData(DRAG_MIME, JSON.stringify(payload));
	e.dataTransfer.effectAllowed = 'move';
};

export const parseNode = (e: DragEvent): DragNodePayload | null => {
	try {
		const raw = e.dataTransfer?.getData(DRAG_MIME);
		if (!raw) return null;
		const data = JSON.parse(raw);
		if (data && typeof data.id === 'string' && typeof data.type === 'string') {
			return data as DragNodePayload;
		}
	} catch (error) {
		console.error('Error parsing drag payload:', error);
	}
	return null;
};
