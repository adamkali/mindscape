import { A } from '@solidjs/router';
import { createResource, createEffect, createSignal, Show } from 'solid-js';
import {
	BackgroundApi,
	FoldersApi,
	BookmarksApi,
	type ResponsesFolderData,
} from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { Header } from '@/components/Header';
import CreateFolderComponent from '@/components/CreateFolderComponent';
import { EmptyGuid } from '@/utils';
import FolderComponent from '@/components/FolderComponent';
import { Button } from '@/components/atoms';
import CreateBookmarkComponent from '@/components/CreateBookmarkComponent';
import { AddFolder } from '@/components/icons';

const Home = () => {
	const auth = useAuth();
	const [folders, setFolders] = createSignal<ResponsesFolderData[]>([]);
	const [focusedNodeId, setFocusedNodeId] = createSignal<string>('');
	const [isLoadingFolders, setIsLoadingFolders] = createSignal(false);
	const [showCreateFolder, setShowCreateFolder] = createSignal(false);
	const [showCreateBookmark, setShowCreateBookmark] = createSignal(false);
	const [isDragOverRoot, setIsDragOverRoot] = createSignal(false);
	const foldersApi = new FoldersApi();
	const user = auth.user();

	createEffect(() => {
		fetchRootFolders();
	});

	const [defaultBackground] = createResource(async () => {
		const api = new BackgroundApi();
		const response = await api.getDefaultBackground();
		if (response.success && response.data) {
			return response.data;
		} else {
			throw new Error(
				'Failed to fetch default background: ' + response.message,
			);
		}
	});

	createEffect(() => {
		if (focusedNodeId() === EmptyGuid) {
			setFocusedNodeId('');
		}
		if (showCreateFolder()) {
			setShowCreateBookmark(false);
		}
		if (showCreateBookmark()) {
			setShowCreateFolder(false);
		}
	});

	// Debug effect to track selected node ID changes
	createEffect(() => {
		const currentId = focusedNodeId();
		console.log(`id_change::selected_node: ${currentId || 'none'}`);
	});

	const fetchRootFolders = async () => {
		if (!auth.token()) return;
		setIsLoadingFolders(true);
		try {
			const response = await foldersApi.getRootFolders({
				authorization: `Bearer ${auth.token()}`,
			});
			if (response.success && response.data) {
				const folders = response.data;
				folders.sort((a, b) => a.name!.localeCompare(b.name!));

				// Removed auto-selection to ensure no node is selected on page load
				setFolders(folders);
			}
		} catch (error) {
			console.error('Failed to fetch folders:', error);
		} finally {
			setIsLoadingFolders(false);
		}
	};

	const deleteFolder = async (folderId: string) => {
		if (!auth.token() || !folderId) return;

		try {
			const response = await foldersApi.deleteFolder({
				folderId,
				authorization: `Bearer ${auth.token()}`,
			});
			if (response.success) {
				await fetchRootFolders();
				if (focusedNodeId() === folderId) {
					setFocusedNodeId('');
				}
			}
		} catch (error) {
			console.error('Failed to delete folder:', error);
		}
	};

	const deleteBookmark = async (bookmarkId: string) => {
		if (!auth.token() || !bookmarkId) return;

		try {
			const bookmarksApi = new BookmarksApi();
			const response = await bookmarksApi.deleteBookmark({
				bookmarkId: bookmarkId,
				authorization: `Bearer ${auth.token()}`,
			});

			if (response.success) {
				await fetchRootFolders(); // Refresh the data
			} else {
				console.error('Failed to delete bookmark:', response.message);
			}
		} catch (error) {
			console.error('Failed to delete bookmark:', error);
		}
	};

	if (!auth.isAuthenticated() || !user) {
		return (
			<div class="min-h-screen flex items-center justify-center">
				<div class="text-center">
					<h2 class="text-2xl font-bold text-foreground mb-4">
						Please sign in to continue
					</h2>
					<A href="/login" class="text-primary hover:text-primary/80">
						Sign in
					</A>
				</div>
			</div>
		);
	}

	// set overwritable callback for create bookmark component
	let bookmarkRefresh: () => void = () => {};
	const openCreateBookmarkComponent = (folderBookmarkRefresh?: () => void) => {
		setShowCreateBookmark(true);
		setShowCreateFolder(false);
		if (folderBookmarkRefresh) {
			bookmarkRefresh = folderBookmarkRefresh;
		}
	};

	const closeCreateBookmarkComponent = () => {
		setShowCreateBookmark(false);
		bookmarkRefresh();
		bookmarkRefresh = () => {};
	};

	const handleRootDragOver = (e: DragEvent) => {
		// Only handle if the target is the root container or empty space, not folder cards
		const target = e.target as HTMLElement;
		const isOverFolderCard = target.closest('[draggable="true"]') !== null;
		
		if (!isOverFolderCard) {
			e.preventDefault();
			e.dataTransfer!.dropEffect = 'move';
			setIsDragOverRoot(true);
		}
	};

	const handleRootDragLeave = (e: DragEvent) => {
		// Only clear the root drag state if we're actually leaving the root container
		const target = e.target as HTMLElement;
		const relatedTarget = e.relatedTarget as HTMLElement;
		
		// If we're moving to a child element, don't clear the drag state
		if (relatedTarget && target.contains(relatedTarget)) {
			return;
		}
		
		setIsDragOverRoot(false);
	};

	const handleRootDrop = async (e: DragEvent) => {
		// Only handle if the target is actually the root container or empty space
		const target = e.target as HTMLElement;
		const isOverFolderCard = target.closest('[draggable="true"]') !== null;
		
		if (isOverFolderCard) {
			// Let the folder card handle this drop
			return;
		}
		
		e.preventDefault();
		setIsDragOverRoot(false);
		try {
			const data = JSON.parse(e.dataTransfer!.getData('text/plain'));
			console.log('Drop data to root:', data);
			console.log(`id_change::dragged_in_node: ${data.id} (${data.type}) -> root`);
			
			if (data.type === 'folder' && data.id) {
				try {
					const response = await foldersApi.moveFolder({
						moveFolderRequest: {
							userId: user?.id,
							folderId: data.id,
							newParentId: undefined, // Moving to root
						},
						authorization: `Bearer ${auth.token()}`,
					});
					
					if (response.success) {
						console.log('Folder moved to root successfully:', response.data);
						// Refresh root folders to show the moved item
						await fetchRootFolders();
					} else {
						console.error('Failed to move folder to root:', response.message);
						alert(`Failed to move folder to root: ${response.message}`);
					}
				} catch (error) {
					console.error('Error moving folder to root:', error);
					alert(`Error moving folder to root: ${error}`);
				}
			}
		} catch (error) {
			console.error('Error parsing drop data:', error);
		}
	};

	return (
		<div
			class="min-h-screen bg-background"
			style={{ 
				'background-image': `url(${defaultBackground()})`,
				'background-size': 'cover',
				'background-position': 'center center',
				'background-repeat': 'no-repeat',
				'background-attachment': 'fixed'
			}}
			onClick={(e) => {
				// Deselect node when clicking on background areas
				const target = e.target as HTMLElement;
				if (
					target.classList.contains('bg-background') ||
					target.closest('.treeview-container') === null
				) {
					setFocusedNodeId('');
				}
			}}
		>
			<Header />

			{/* Main content */}
			<div class="flex h-screen">
				{/* Sidebar with tree */}
				<div 
					class={`treeview-container m-2 p-4 rounded-3xl overflow-y-auto bg-white/10 backdrop-blur-lg border border-white/20 max-h-[calc(100vh-2rem)] min-w-80 shadow-2xl shadow-slate-900/30 ${
						isDragOverRoot() ? 'ring-2 ring-blue-400 bg-blue-100/20' : ''
					}`}
					onDragOver={handleRootDragOver}
					onDragLeave={handleRootDragLeave}
					onDrop={handleRootDrop}
				>
					<div class="p-2">
						<div class="flex items-center justify-between mb-4">
							<Button
								variant="primary"
								onClick={() => {
									setShowCreateFolder(true);
								}}
							>
								<AddFolder />
							</Button>
						</div>

						<Show when={showCreateFolder()}>
							<CreateFolderComponent
								userId={user.id || EmptyGuid}
								parentId={focusedNodeId()}
								auth={auth}
								setShowCreateFolder={setShowCreateFolder}
								folderAPIRef={foldersApi}
								refresh={fetchRootFolders}
							/>
						</Show>

						<Show when={showCreateBookmark()}>
							<CreateBookmarkComponent
								userId={user.id || EmptyGuid}
								parentId={focusedNodeId()}
								auth={auth}
								close={closeCreateBookmarkComponent}
								refreshBookmarks={bookmarkRefresh}
							/>
						</Show>

						<Show
							when={!isLoadingFolders()}
							fallback={
								<div class="text-center py-4 text-foreground/60">
									Loading folders...
								</div>
							}
						>
							<Show
								when={folders().length > 0}
								fallback={
									<div class="text-center py-4 text-foreground/60">
										No folders yet. Create your first folder!
									</div>
								}
							>
								<div class="space-y-4">
									{folders().map((folder) => (
										<FolderComponent
											folder={folder}
											selectedFolder={focusedNodeId}
											setSelectedFolder={setFocusedNodeId}
											deleteFolder={deleteFolder}
											deleteBookmark={deleteBookmark}
											showCreateFolder={openCreateBookmarkComponent}
											indent={0}
										/>
									))}
								</div>
							</Show>
						</Show>
					</div>
				</div>

				{/* <div class="flex-1 flex items-center justify-center"> */}
				{/* 	<div class="text-center text-foreground/60"> */}
				{/* 		<div class="text-4xl mb-4">📁</div> */}
				{/* 		<p class="text-lg">Select a note or bookmark to view</p> */}
				{/* 	</div> */}
				{/* </div> */}
			</div>
		</div>
	);
};

export default Home;
