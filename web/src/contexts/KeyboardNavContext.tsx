import {
	createContext,
	createSignal,
	onCleanup,
	onMount,
	type ParentComponent,
	useContext,
} from 'solid-js';

export type NavMode = 'mouse' | 'nav' | 'edit' | 'create' | 'delete' | 'move';

// Returns true if the handler consumed the event (stops dispatch).
type KeyHandler = (e: KeyboardEvent) => boolean;

export interface TreeNodeEntry {
	id: string;
	type: 'folder' | 'bookmark';
	ref: HTMLElement;
	parentId: string | null;
	name?: string;
	onEdit?: () => void;
	// Folder-only
	isOpen?: () => boolean;
	expand?: () => void;
	collapse?: () => void;
	refresh?: () => Promise<void>;
	// Bookmark-only
	link?: string;
}

export type SearchMode = 'web' | 'tree';

export interface KeyboardNavContextValue {
	mode: () => NavMode;
	setMode: (m: NavMode) => void;
	enter: () => void;
	exit: () => void;
	register: (handler: KeyHandler) => () => void;
	registerNode: (entry: TreeNodeEntry) => () => void;
	getNodes: () => TreeNodeEntry[];
	moveSourceId: () => string | null;
	startMove: (nodeId: string) => void;
	cancelMove: () => void;
	searchMode: () => SearchMode | null;
	openSearch: (mode: SearchMode) => void;
	closeSearch: () => void;
}

const KeyboardNavContext = createContext<KeyboardNavContextValue>();

const isInputFocused = (): boolean => {
	const el = document.activeElement as HTMLElement | null;
	if (!el) return false;
	const tag = el.tagName.toLowerCase();
	return tag === 'input' || tag === 'textarea' || el.isContentEditable;
};

const isModalOpen = (): boolean =>
	document.querySelector('[role="dialog"]') !== null;

export const KeyboardNavProvider: ParentComponent = (props) => {
	const [mode, setMode] = createSignal<NavMode>('mouse');
	const [moveSourceId, setMoveSourceId] = createSignal<string | null>(null);
	const [searchMode, setSearchMode] = createSignal<SearchMode | null>(null);
	const handlers: KeyHandler[] = [];
	const nodeRegistry = new Map<string, TreeNodeEntry>();

	const enter = () => setMode('nav');
	const exit = () => setMode('mouse');
	const startMove = (nodeId: string) => {
		setMoveSourceId(nodeId);
		setMode('move');
	};
	const cancelMove = () => {
		setMoveSourceId(null);
		setMode('nav');
	};
	const openSearch = (m: SearchMode) => setSearchMode(m);
	const closeSearch = () => setSearchMode(null);

	const register = (handler: KeyHandler): (() => void) => {
		handlers.push(handler);
		return () => {
			const idx = handlers.indexOf(handler);
			if (idx !== -1) handlers.splice(idx, 1);
		};
	};

	const registerNode = (entry: TreeNodeEntry): (() => void) => {
		nodeRegistry.set(entry.id, entry);
		return () => {
			nodeRegistry.delete(entry.id);
		};
	};

	const getNodes = (): TreeNodeEntry[] =>
		[...nodeRegistry.values()].sort((a, b) => {
			const pos = a.ref.compareDocumentPosition(b.ref);
			if (pos & Node.DOCUMENT_POSITION_FOLLOWING) return -1;
			if (pos & Node.DOCUMENT_POSITION_PRECEDING) return 1;
			return 0;
		});

	const handleKeyDown = (e: KeyboardEvent) => {
		if (isInputFocused() || isModalOpen()) return;

		// Search shortcuts work in any nav mode; skip if overlay is already open.
		if (e.key === '/' && searchMode() === null) {
			e.preventDefault();
			openSearch('web');
			return;
		}
		if (e.key === '?' && searchMode() === null) {
			e.preventDefault();
			openSearch('tree');
			return;
		}

		const current = mode();

		if (e.key === 'm' && current === 'mouse') {
			e.preventDefault();
			setMode('nav');
			return;
		}

		if (e.key === 'q' && current !== 'mouse') {
			e.preventDefault();
			setMoveSourceId(null);
			setMode('mouse');
			return;
		}

		if (e.key === 'Escape' && current !== 'mouse' && current !== 'nav') {
			e.preventDefault();
			setMoveSourceId(null);
			setMode('nav');
			return;
		}

		if (current !== 'mouse') {
			for (const h of handlers) {
				if (h(e)) return;
			}
		}
	};

	onMount(() => {
		document.addEventListener('keydown', handleKeyDown);
	});

	onCleanup(() => {
		document.removeEventListener('keydown', handleKeyDown);
	});

	const value: KeyboardNavContextValue = {
		mode,
		setMode,
		enter,
		exit,
		register,
		registerNode,
		getNodes,
		moveSourceId,
		startMove,
		cancelMove,
		searchMode,
		openSearch,
		closeSearch,
	};

	return (
		<KeyboardNavContext.Provider value={value}>
			{props.children}
		</KeyboardNavContext.Provider>
	);
};

export const useKeyboardNav = (): KeyboardNavContextValue => {
	const context = useContext(KeyboardNavContext);
	if (!context) {
		throw new Error('useKeyboardNav must be used within a KeyboardNavProvider');
	}
	return context;
};
