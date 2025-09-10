import { A } from '@solidjs/router';
import { createResource, createEffect, createSignal, Show } from 'solid-js';
import { BackgroundApi, FoldersApi, type ResponsesFolderData } from '@/api';
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
			throw new Error('Failed to fetch default background: ' + response.message);
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

				if (folders.length > 0 && !focusedNodeId()) {
					setFocusedNodeId(folders[0].id!);
				}
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
	let bookmarkRefresh: () => void = () => { };
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
		bookmarkRefresh = () => { };
	};


	return (
		<div
			class="min-h-screen bg-background"
			style={{ "background-image": `url(${defaultBackground()})` }}
		>
			<Header />

			{/* Main content */}
			<div class="flex h-screen">
				{/* Sidebar with tree */}
				<div class="m-2 p-1 rounded-2xl overflow-y-auto">
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
