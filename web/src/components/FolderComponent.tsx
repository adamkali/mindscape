import { type ComponentProps, createMemo, createSignal, For } from 'solid-js';
import type { RepositoryBookmark, ResponsesFolderData } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { useAuthenticatedApi } from '@/utils/useApi';
import BookmarkComponent from './BookmarkComponent';
import FolderCard from './FolderCard';

interface FolderComponentProps extends ComponentProps<'div'> {
	folder: ResponsesFolderData;
	selectedFolder: () => string;
	setSelectedFolder: (id: string) => void;
	deleteFolder: (id: string) => void;
	deleteBookmark?: (bookmarkId: string) => void;
	editBookmark?: (bookmark: RepositoryBookmark) => void;
	showCreateFolder: (folderBookmarkRefresh?: () => void) => void;
	onFolderSelected?: (
		refreshFn: () => Promise<void>,
		folderName: string,
	) => void;
	indent: number;
}

export default function FolderComponent(props: FolderComponentProps) {
	const { folder, selectedFolder, setSelectedFolder, deleteFolder, indent } =
		props;

	const api = useAuthenticatedApi();
	const auth = useAuth();

	const [isFolderOpen, setIsFolderOpen] = createSignal(false);
	const [indentNow, _setIndentNow] = createSignal(indent);
	const indentNowCN = () => {
		return indent * 2 + 'rem';
	};
	const indentNext = () => {
		return indent + 1;
	};
	const indentNextCN = () => indentNext() * 2 + 'rem';

	console.log({
		Prev: indent,
		Indent: indentNow(),
		IndentClass: indentNowCN(),
		IndentNext: indentNext(),
		IndentNextClass: indentNextCN(),
	});

	const [children, setChildren] = createSignal<ResponsesFolderData[]>(
		folder.children || [],
	);
	const [bookmarks, setBookmarks] = createSignal<RepositoryBookmark[]>(
		folder.bookmarks || [],
	);

	const openFolder = async () => {
		if (isFolderOpen()) {
			setIsFolderOpen(false);
		} else {
			setIsFolderOpen(true);
			// try to get the folder content from the server
			const response = await api.folders.getFolders({
				folderId: folder.id || '',
				authorization: `Bearer ${auth.token()}`,
			});
			if (response.success && response.data) {
				setChildren(response.data.children || []);
				setBookmarks(response.data.bookmarks || []);
				setIsFolderOpen(true);
			} else {
				console.error(
					'Failed to get folder content:',
					response.message,
					response.data,
					response.success,
				);
			}
		}
	};

	const refreshFolder = async () => {
		const response = await api.folders.getFolders({
			folderId: folder.id || '',
			authorization: `Bearer ${auth.token()}`,
		});
		if (response.success && response.data) {
			setChildren(response.data.children || []);
			setBookmarks(response.data.bookmarks || []);
		} else {
			console.error(
				'Failed to refresh folder:',
				response.message,
				response.data,
				response.success,
			);
		}
	};

	const refreshBookmarks = async () => {
		const response = await api.bookmarks.getBookmarks({
			parentId: folder.id || '',
			authorization: `Bearer ${auth.token()}`,
		});
		if (response.success && response.data) {
			setBookmarks(response.data || []);
			setIsFolderOpen(true);
		} else {
			console.error(
				'Failed to get folder content:',
				response.message,
				response.data,
				response.success,
			);
		}
	};

	const renderChildren = () =>
		createMemo(() => {
			console.log({
				'Rendering children for folder:': folder.id,
				'Children:': children(),
				'Bookmarks:': bookmarks(),
			});
			return (
				<div class="space-y-4">
					<For each={children()}>
						{(child) => (
							<FolderComponent
								folder={child}
								selectedFolder={selectedFolder}
								setSelectedFolder={setSelectedFolder}
								deleteFolder={deleteFolder}
								deleteBookmark={props.deleteBookmark}
								editBookmark={props.editBookmark}
								indent={indentNext()}
								showCreateFolder={props.showCreateFolder}
								onFolderSelected={props.onFolderSelected}
							/>
						)}
					</For>
					<For each={bookmarks()}>
						{(bookmark) => (
							<BookmarkComponent
								bookmark={bookmark}
								selected={selectedFolder}
								setSelected={setSelectedFolder}
								indent={indentNext()}
								deleteBookmark={props.deleteBookmark}
								editBookmark={props.editBookmark}
							/>
						)}
					</For>
				</div>
			);
		}, [children]);

	const handleDrop = async (data: any) => {
		console.log('Drop data:', data, 'into folder:', folder.id);
		console.log(
			`id_change::dragged_in_node: ${data.id} (${data.type}) -> ${folder.id}`,
		);

		if (data.type === 'folder' && data.id && folder.id) {
			// Prevent dropping folder into itself
			if (data.id === folder.id) {
				alert('Cannot move folder into itself');
				return;
			}

			try {
				const response = await api.folders.moveFolder({
					moveFolderRequest: {
						userId: auth.user()?.id,
						folderId: data.id,
						newParentId: folder.id,
					},
					authorization: `Bearer ${auth.token()}`,
				});

				if (response.success) {
					console.log('Folder moved successfully:', response.data);
					// Refresh the folder to show the moved item
					await openFolder();
				} else {
					console.error('Failed to move folder:', response.message);
					alert(`Failed to move folder: ${response.message}`);
				}
			} catch (error) {
				console.error('Error moving folder:', error);
				alert(`Error moving folder: ${error}`);
			}
		} else if (data.type === 'bookmark' && data.id && folder.id) {
			try {
				const response = await api.bookmarks.moveBookmark({
					moveBookmarkRequest: {
						userId: auth.user()?.id,
						bookmarkId: data.id,
						newParentId: folder.id,
					},
					authorization: `Bearer ${auth.token()}`,
				});

				if (response.success) {
					console.log('Bookmark moved successfully:', response.data);
					// Refresh the folder to show the moved item
					await openFolder();
					// Also refresh bookmarks to update the display
					await refreshBookmarks();
				} else {
					console.error('Failed to move bookmark:', response.message);
					alert(`Failed to move bookmark: ${response.message}`);
				}
			} catch (error) {
				console.error('Error moving bookmark:', error);
				alert(`Error moving bookmark: ${error}`);
			}
		} else {
			console.warn('Invalid drop data or missing IDs:', data);
		}
	};

	return (
		<>
			<div style={{ 'margin-left': indentNowCN() }}>
				<FolderCard
					folder={folder}
					isSelected={folder.id === selectedFolder()}
					onSelect={(folderId) => {
						setSelectedFolder(folderId);
						openFolder();
						props.onFolderSelected?.(refreshFolder, folder.name || '');
					}}
					onDelete={deleteFolder}
					onCreateBookmark={(folderId) => {
						setSelectedFolder(folderId);
						props.showCreateFolder(async () => {
							await refreshFolder();
						});
					}}
					onDrop={handleDrop}
					draggable={true}
				/>
			</div>
			{isFolderOpen() && renderChildren()}
		</>
	);
}
