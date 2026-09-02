import {
	type ComponentProps,
	createSignal,
	For,
	onCleanup,
	onMount,
	Show,
} from 'solid-js';
import type { RepositoryBookmark, ResponsesFolderData } from '@/api';
import { useKeyboardNav } from '@/contexts/KeyboardNavContext';
import { useTree } from '@/contexts/TreeContext';
import type { DragNodePayload } from '@/utils/dragNode';
import { useAuthenticatedApi } from '@/utils/useApi';
import BookmarkComponent from './BookmarkComponent';
import FolderCard from './FolderCard';

interface FolderComponentProps extends ComponentProps<'div'> {
	folder: ResponsesFolderData;
	editBookmark?: (bookmark: RepositoryBookmark) => void;
	openCreateBookmark: () => void;
	onFolderSelected?: (folderName: string) => void;
	indent: number;
	parentFolderId?: string | null;
}

export default function FolderComponent(props: FolderComponentProps) {
	const { folder, indent } = props;

	const api = useAuthenticatedApi();
	const nav = useKeyboardNav();
	const tree = useTree();

	const [isFolderOpen, setIsFolderOpen] = createSignal(false);
	const indentNowCN = () => `${indent * 2}rem`;
	const indentNext = () => indent + 1;

	const [children, setChildren] = createSignal<ResponsesFolderData[]>(
		folder.children || [],
	);
	const [bookmarks, setBookmarks] = createSignal<RepositoryBookmark[]>(
		folder.bookmarks || [],
	);

	const fetchContent = async () => {
		const response = await api.folders.getFolders({
			folderId: folder.id || '',
		});
		if (response.success && response.data) {
			setChildren(response.data.children || []);
			setBookmarks(response.data.bookmarks || []);
			return true;
		}
		console.error('Failed to get folder content:', response.message);
		return false;
	};

	const openFolder = async () => {
		if (isFolderOpen()) {
			setIsFolderOpen(false);
			return;
		}
		if (await fetchContent()) setIsFolderOpen(true);
	};

	const expandFolder = () => {
		if (!isFolderOpen()) openFolder();
	};
	const collapseFolder = () => setIsFolderOpen(false);

	// Registered as this folder's refresh hook: re-fetch and reveal contents so
	// items created/moved/deleted under it show without a full-tree refetch.
	const refreshFolder = async () => {
		if (await fetchContent()) setIsFolderOpen(true);
	};

	let rowRef: HTMLDivElement | undefined;

	onMount(() => {
		if (!folder.id || !rowRef) return;
		const unregister = nav.registerNode({
			id: folder.id,
			type: 'folder',
			ref: rowRef,
			parentId: props.parentFolderId ?? null,
			name: folder.name ?? '',
			isOpen: isFolderOpen,
			expand: expandFolder,
			collapse: collapseFolder,
			refresh: refreshFolder,
		});
		onCleanup(unregister);
	});

	const handleDrop = (payload: DragNodePayload) => {
		// Dropping a node onto its own card is a no-op.
		if (!folder.id || payload.id === folder.id) return;
		tree.move(payload.id, folder.id);
	};

	return (
		<>
			<div
				ref={(r) => {
					rowRef = r;
				}}
				tabIndex={-1}
				style={{ 'margin-left': indentNowCN() }}
				class={
					folder.id === tree.focusedId()
						? 'ring-2 ring-white/50 dark:bg-white/20 bg-black/10 rounded-md'
						: ''
				}
			>
				<FolderCard
					folder={folder}
					isSelected={folder.id === tree.focusedId()}
					onSelect={(folderId) => {
						tree.setFocus(folderId);
						openFolder();
						props.onFolderSelected?.(folder.name || '');
					}}
					onDelete={(folderId) => tree.deleteFolder(folderId)}
					onCreateBookmark={(folderId) => {
						tree.setFocus(folderId);
						props.openCreateBookmark();
					}}
					onDrop={handleDrop}
				/>
			</div>
			<Show when={isFolderOpen()}>
				<div class="space-y-4">
					<For each={children()}>
						{(child) => (
							<FolderComponent
								folder={child}
								editBookmark={props.editBookmark}
								indent={indentNext()}
								openCreateBookmark={props.openCreateBookmark}
								onFolderSelected={props.onFolderSelected}
								parentFolderId={folder.id}
							/>
						)}
					</For>
					<For each={bookmarks()}>
						{(bookmark) => (
							<BookmarkComponent
								bookmark={bookmark}
								indent={indentNext()}
								editBookmark={props.editBookmark}
							/>
						)}
					</For>
				</div>
			</Show>
		</>
	);
}
