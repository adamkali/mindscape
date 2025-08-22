import { A } from '@solidjs/router';
import { createEffect, createSignal } from 'solid-js';
import { FoldersApi, type ResponsesFolderData } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { Header } from '@/components/Header';
import CreateFolderComponent from '@/components/CreateFolderComponent';
import { EmptyGuid } from '@/utils';
import FolderComponent from '@/components/FolderComponent';


interface Mode {
	type: 'tree' | 'notes' | 'search';
	focusedNodeId?: string;
}

const Home = () => {
	const auth = useAuth();
	const [mode, _] = createSignal<Mode>({ type: 'tree' });
	const [folders, setFolders] = createSignal<ResponsesFolderData[]>([]);
	const [focusedNodeId, setFocusedNodeId] = createSignal<string>('');
	const [isLoadingFolders, setIsLoadingFolders] = createSignal(false);
	const [showCreateFolder, setShowCreateFolder] = createSignal(false);

	const foldersApi = new FoldersApi();
	const user = auth.user();


	createEffect(() => {
		fetchRootFolders();
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

	return (
		<div class="min-h-screen bg-background">
			<Header />

			{/* Main content */}
			<div class="flex h-screen">
				{/* Sidebar with tree */}
				<div class="w-80 border-r border-card-foreground/20 bg-card overflow-y-auto">
					<div class="p-4">
						<div class="flex items-center justify-between mb-4">
							{mode().type === 'tree' && (
								<button
									onClick={() => setShowCreateFolder(true)}
									class="text-xs px-2 py-1 bg-primary text-primary-foreground rounded hover:bg-primary/80 transition-colors"
								>
									+ New Folder
								</button>
							)}
						</div>

						{showCreateFolder() && (
							<CreateFolderComponent
								userId={user.id || EmptyGuid}
								parentId={undefined}
								auth={auth}
								setShowCreateFolder={setShowCreateFolder}
								folderAPIRef={foldersApi}
							/>
						)}

						{mode().type === 'tree' && (
							<div class="space-y-1">
								{isLoadingFolders() ? (
									<div class="text-center py-4 text-foreground/60">
										Loading folders...
									</div>
								) : folders().length > 0 ? (
									folders().map((folder) => (
										<FolderComponent
											folder={folder}
											selectedFolder={setFocusedNodeId}
											deleteFolder={deleteFolder}
										/>
									))
									
								) : (
									<div class="text-center py-4 text-foreground/60">
										No folders yet. Create your first folder!
									</div>
								)}
							</div>
						)}
					</div>
				</div>

				{/* Main content area */}
				{mode().type === 'tree' && (
					<div class="flex-1 flex items-center justify-center">
						<div class="text-center text-foreground/60">
							<div class="text-4xl mb-4">📁</div>
							<p class="text-lg">Select a note or bookmark to view</p>
						</div>
					</div>
				)}
			</div>
		</div>
	);
};

export default Home;
