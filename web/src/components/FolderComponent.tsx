import {
	FoldersApi,
	type ResponsesFolderData,
	type RepositoryBookmark,
	BookmarksApi,
} from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { cn } from '@/utils/cn';
import {
	createMemo,
	createSignal,
	For,
	Show,
	type ComponentProps,
} from 'solid-js';
import { Button } from './atoms';
import BookmarkComponent from './BookmarkComponent';
import { DeleteIcon } from './icons';

interface FolderComponentProps extends ComponentProps<'div'> {
	folder: ResponsesFolderData;
	selectedFolder: () => string;
	setSelectedFolder: (id: string) => void;
	deleteFolder: (id: string) => void;
	deleteBookmark?: (bookmarkId: string, parentFolderId: string) => void;
	showCreateFolder: (folderBookmarkRefresh?: () => void) => void;
	indent: number;
}

export default function FolderComponent(props: FolderComponentProps) {
	const { folder, selectedFolder, setSelectedFolder, deleteFolder, indent } =
		props;

	const foldersApi = new FoldersApi();
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
			const response = await foldersApi.getFolders({
				folderId: folder.id || '',
				authorization: `Bearer ${auth.token()}`,
			});
			if (response.success && response.data?.children) {
				setChildren(response.data.children || []);
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

	const refreshBookmarks = async () => {
		const response = await new BookmarksApi().align({
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
								indent={indentNext()}
								showCreateFolder={props.showCreateFolder}
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
							/>
						)}
					</For>
					<Show when={selectedFolder() === folder.id}>
						<Button
							style={{ 'margin-left': indentNextCN() }}
							variant="secondary"
							onClick={() => {
								props.showCreateFolder(async () => {
									await refreshBookmarks();
								});
							}}
						>
							Create Bookmark
						</Button>
					</Show>
				</div>
			);
		}, [children]);

	const childClassName = (): string => {
		let isFolder = '';
		if (folder.id === selectedFolder()) {
			isFolder = 'bg-white/40 backdrop-blur-md border-white/50';
		}
		return cn(
			`flex flex-row items-center py-1 px-2 cursor-pointer bg-white/20 backdrop-blur-sm border border-white/30 text-white hover:bg-white/30 rounded-lg
		shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out space-x-1 justify-between w-64`,
			isFolder,
		);
	};

	return (
		<>
			<div
				class={childClassName()}
				style={{ 'margin-left': indentNowCN() }}
				onClick={() => {
					setSelectedFolder(folder.id || '');
					openFolder();
				}}
			>
				<div>
					<span class="mr-2 text-base font-bold">{folder.name}</span>
				</div>
				<Button
					variant="danger"
					class="p-1 text-xs"
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
						if (folder.id) {
							deleteFolder(folder.id);
						}
					}}
				>
					<DeleteIcon />
				</Button>
			</div>
			{isFolderOpen() && renderChildren()}
		</>
	);
}
