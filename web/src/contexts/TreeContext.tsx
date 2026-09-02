import {
	createContext,
	createMemo,
	createSignal,
	type ParentComponent,
	useContext,
} from 'solid-js';
import type { ResponsesFolderData } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { useKeyboardNav } from '@/contexts/KeyboardNavContext';
import { useAuthenticatedApi } from '@/utils/useApi';

export interface TreeContextValue {
	// Selection / focus — the single source of truth for "which node is active".
	focusedId: () => string;
	setFocus: (id: string) => void;
	// The folder new items are created into: the focused folder itself, or the
	// containing folder of a focused bookmark, or null for root.
	currentParentFolder: () => string | null;
	// Root folder list (parent === null refresh target).
	rootFolders: () => ResponsesFolderData[];
	refreshRoot: () => Promise<void>;
	// Refresh a specific folder by id, or the root list when id is null.
	refreshFolderById: (id: string | null) => Promise<void>;
	// Mutations — the only place tree-changing API calls live.
	createBookmark: (params: { name: string; link: string }) => Promise<boolean>;
	createFolder: (params: {
		name: string;
		description?: string;
	}) => Promise<boolean>;
	deleteBookmark: (id: string) => Promise<void>;
	deleteFolder: (id: string) => Promise<void>;
	updateBookmark: (
		id: string,
		params: { name: string; link: string },
	) => Promise<boolean>;
	updateFolder: (
		id: string,
		params: { name?: string; description?: string },
	) => Promise<boolean>;
	move: (sourceId: string, destId: string | null) => Promise<void>;
}

const TreeContext = createContext<TreeContextValue>();

export const TreeProvider: ParentComponent = (props) => {
	const auth = useAuth();
	const api = useAuthenticatedApi();
	const nav = useKeyboardNav();

	const [focusedId, setFocusedId] = createSignal<string>('');
	const [rootFolders, setRootFolders] = createSignal<ResponsesFolderData[]>([]);

	const setFocus = (id: string) => setFocusedId(id);

	const currentParentFolder = createMemo<string | null>(() => {
		const id = focusedId();
		if (!id) return null;
		const node = nav.getNodes().find((n) => n.id === id);
		if (!node) return null;
		// A focused folder is itself the parent for new children; a focused
		// bookmark resolves up to its containing folder.
		return node.type === 'folder' ? node.id : (node.parentId ?? null);
	});

	const refreshRoot = async () => {
		if (!auth.token()) return;
		const response = await api.folders.getRootFolders({});
		if (response.success && response.data) {
			const data = [...response.data].sort((a, b) =>
				(a.name ?? '').localeCompare(b.name ?? ''),
			);
			setRootFolders(data);
		} else {
			console.error('Failed to fetch root folders:', response.message);
		}
	};

	// Refresh whichever folder owns the mutated node. null === root list.
	const refreshFolderById = async (id: string | null) => {
		if (!id) {
			await refreshRoot();
			return;
		}
		const node = nav.getNodes().find((n) => n.id === id);
		await node?.refresh?.();
	};

	// The containing folder of a node, or null when it lives at the root.
	const parentOf = (id: string): string | null =>
		nav.getNodes().find((n) => n.id === id)?.parentId ?? null;

	const createBookmark = async (params: { name: string; link: string }) => {
		if (!auth.token()) return false;
		const parent = currentParentFolder();
		const response = await api.bookmarks.createBookmark({
			createBookmarkRequest: {
				userId: auth.user()?.id,
				folderId: parent ?? undefined,
				link: params.link,
				name: params.name,
			},
		});
		if (response.success && response.data) {
			await refreshFolderById(parent);
			return true;
		}
		console.error('Failed to create bookmark:', response.message);
		return false;
	};

	const createFolder = async (params: {
		name: string;
		description?: string;
	}) => {
		if (!auth.token()) return false;
		const parent = currentParentFolder();
		const response = await api.folders.createFolder({
			createFolderRequest: {
				userId: auth.user()?.id,
				parentId: parent ?? undefined,
				name: params.name,
				description: params.description,
			},
		});
		if (response.success) {
			await refreshFolderById(parent);
			return true;
		}
		console.error('Failed to create folder:', response.message);
		return false;
	};

	const deleteBookmark = async (id: string) => {
		if (!auth.token() || !id) return;
		const parentId = parentOf(id);
		try {
			const response = await api.bookmarks.deleteBookmark({
				bookmarkId: id,
			});
			if (response.success) {
				if (focusedId() === id) setFocus('');
				await refreshFolderById(parentId);
			} else {
				console.error('Failed to delete bookmark:', response.message);
			}
		} catch (error) {
			console.error('Failed to delete bookmark:', error);
		}
	};

	const deleteFolder = async (id: string) => {
		if (!auth.token() || !id) return;
		const parentId = parentOf(id);
		try {
			const response = await api.folders.deleteFolder({
				folderId: id,
			});
			if (response.success) {
				if (focusedId() === id) setFocus('');
				await refreshFolderById(parentId);
			} else {
				console.error('Failed to delete folder:', response.message);
			}
		} catch (error) {
			console.error('Failed to delete folder:', error);
		}
	};

	const updateBookmark = async (
		id: string,
		params: { name: string; link: string },
	) => {
		if (!auth.token() || !id) return false;
		try {
			const response = await api.bookmarks.updateBookmark({
				bookmarkId: id,
				updateBookmarkRequest: {
					userId: auth.user()?.id,
					bookmarkId: id,
					name: params.name,
					link: params.link,
				},
			});
			if (response.success) {
				await refreshFolderById(parentOf(id));
				return true;
			}
			console.error('Failed to update bookmark:', response.message);
		} catch (error) {
			console.error('Failed to update bookmark:', error);
		}
		return false;
	};

	const updateFolder = async (
		id: string,
		params: { name?: string; description?: string },
	) => {
		if (!auth.token() || !id) return false;
		try {
			const response = await api.folders.updateFolder({
				folderId: id,
				updateFolderRequest: {
					userId: auth.user()?.id,
					folderId: id,
					name: params.name,
					description: params.description,
				},
			});
			if (response.success) {
				// A folder's card is rendered by its parent, so refresh the parent
				// (or the root list for a top-level folder) to reflect the change.
				await refreshFolderById(parentOf(id));
				return true;
			}
			console.error('Failed to update folder:', response.message);
		} catch (error) {
			console.error('Failed to update folder:', error);
		}
		return false;
	};

	const move = async (sourceId: string, destId: string | null) => {
		if (!auth.token() || !sourceId || sourceId === destId) return;
		const nodes = nav.getNodes();
		const src = nodes.find((n) => n.id === sourceId);
		if (!src) return;
		// Destination must be a folder (or root when null).
		if (destId) {
			const dst = nodes.find((n) => n.id === destId);
			if (dst?.type !== 'folder') return;
		}
		const oldParent = src.parentId ?? null;

		try {
			let ok = false;
			if (src.type === 'folder') {
				const response = await api.folders.moveFolder({
					moveFolderRequest: {
						userId: auth.user()?.id,
						folderId: sourceId,
						newParentId: destId ?? undefined,
					},
				});
				ok = response.success ?? false;
				if (!ok) console.error('Failed to move folder:', response.message);
			} else if (src.type === 'bookmark') {
				const response = await api.bookmarks.moveBookmark({
					moveBookmarkRequest: {
						userId: auth.user()?.id,
						bookmarkId: sourceId,
						newParentId: destId ?? undefined,
					},
				});
				ok = response.success ?? false;
				if (!ok) console.error('Failed to move bookmark:', response.message);
			}
			if (ok) {
				await refreshFolderById(oldParent);
				if (destId !== oldParent) await refreshFolderById(destId);
			}
		} catch (error) {
			console.error('Failed to move node:', error);
		}
	};

	const value: TreeContextValue = {
		focusedId,
		setFocus,
		currentParentFolder,
		rootFolders,
		refreshRoot,
		refreshFolderById,
		createBookmark,
		createFolder,
		deleteBookmark,
		deleteFolder,
		updateBookmark,
		updateFolder,
		move,
	};

	return (
		<TreeContext.Provider value={value}>{props.children}</TreeContext.Provider>
	);
};

export const useTree = (): TreeContextValue => {
	const context = useContext(TreeContext);
	if (!context) {
		throw new Error('useTree must be used within a TreeProvider');
	}
	return context;
};
