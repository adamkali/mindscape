import { A } from '@solidjs/router';
import { createEffect, createSignal, Show } from 'solid-js';
import {
    BookmarksApi,
    FoldersApi,
    type ResponseError,
    type ResponsesFolderData,
} from '@/api';
import Components from '@/components';
import { Button, Input } from '@/components/atoms';
import CreateBookmarkComponent from '@/components/CreateBookmarkComponent';
import CreateFolderComponent from '@/components/CreateFolderComponent';
import FolderComponent from '@/components/FolderComponent';
import { Header } from '@/components/Header';
import { AddFolder, EditIcon, SaveIcon } from '@/components/icons';
import { useAuth } from '@/contexts/AuthContext';
import { useBackgroundStyle } from '@/hooks/useBackground';
import { EmptyGuid } from '@/utils';

const Home = () => {
    const auth = useAuth();
    const backgroundStyle = useBackgroundStyle();
    const [folders, setFolders] = createSignal<ResponsesFolderData[]>([]);
    const [focusedNodeId, setFocusedNodeId] = createSignal<string>('');
    const [isLoadingFolders, setIsLoadingFolders] = createSignal(false);
    const [showCreateFolder, setShowCreateFolder] = createSignal(false);
    const [showCreateBookmark, setShowCreateBookmark] = createSignal(false);
    const [isDragOverRoot, setIsDragOverRoot] = createSignal(false);
    const [focusedFolderName, setFocusedFolderName] = createSignal('');
    const [isEditingFolderName, setIsEditingFolderName] = createSignal(false);
    const [editFolderName, setEditFolderName] = createSignal('');
    const foldersApi = new FoldersApi();
    const user = auth.user();

    createEffect(() => {
        fetchRootFolders();
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
            (error as ResponseError).response.status === 401 && auth.logout();
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
            <div class="h-screen flex items-center justify-center overflow-hidden">
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

    // set overwritable callback for create folder component (parent folder refresh)
    let folderRefresh: () => void = () => { };
    const openCreateFolderComponent = () => {
        setShowCreateFolder(true);
        setShowCreateBookmark(false);
    };

    const handleFolderSelected = (
        refreshFn: () => Promise<void>,
        folderName: string,
    ) => {
        folderRefresh = refreshFn;
        setFocusedFolderName(folderName);
        setIsEditingFolderName(false);
    };

    const handleSaveFolderName = async () => {
        const folderId = focusedNodeId();
        const newName = editFolderName().trim();
        if (!folderId || !newName || !auth.token()) return;

        try {
            const response = await foldersApi.updateFolder({
                folderId,
                updateFolderRequest: {
                    userId: user?.id,
                    folderId,
                    name: newName,
                },
                authorization: `Bearer ${auth.token()}`,
            });
            if (response.success) {
                setFocusedFolderName(newName);
                setIsEditingFolderName(false);
                await fetchRootFolders();
                folderRefresh();
            }
        } catch (error) {
            console.error('Failed to rename folder:', error);
        }
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
            console.log(
                `id_change::dragged_in_node: ${data.id} (${data.type}) -> root`,
            );

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
            class="h-screen overflow-hidden bg-background"
            style={backgroundStyle()}
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

            <div class="flex h-screen flex-row">
                <div
                    class={`treeview-container m-2 p-4 rounded-lg flex flex-col overflow-hidden bg-background backdrop-blur-lg border border-white/20 max-h-[calc(100vh-2rem)] min-w-80 shadow-2xl shadow-slate-900/30 dark:border-slate-700/50 dark:shadow-black/30 ${isDragOverRoot() ? 'ring-2 ring-blue-400 bg-blue-100/20' : ''
                        }`}
                    onDragOver={handleRootDragOver}
                    onDragLeave={handleRootDragLeave}
                    onDrop={handleRootDrop}
                >
                    <div class="pb-4 p-2 bg-glass-bg rounded-md flex-shrink-0">
                        <div class="flex items-center justify-between mb-2">
                            <Button
                                class="p-1 text-xs flex-shrink-0"
                                variant="secondary"
                                onClick={() => {
                                    openCreateFolderComponent();
                                }}
                            >
                                <AddFolder />
                            </Button>

                            <Show when={focusedNodeId()}>
                                <Show
                                    when={isEditingFolderName()}
                                    fallback={
                                        <span class="text-sm font-semibold text-foreground truncate">
                                            {focusedFolderName()}
                                        </span>
                                    }
                                >
                                    <Input
                                        label="Folder name"
                                        type="text"
                                        value={editFolderName()}
                                        onInput={(e) => setEditFolderName(e.currentTarget.value)}
                                        onKeyDown={(e) => {
                                            if (e.key === 'Enter') handleSaveFolderName();
                                            if (e.key === 'Escape') setIsEditingFolderName(false);
                                        }}
                                        class="text-sm h-7 min-w-0 flex-1"
                                    />
                                </Show>
                                <Button
                                    variant="secondary"
                                    class="p-1 text-xs flex-shrink-0"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        if (isEditingFolderName()) {
                                            handleSaveFolderName();
                                        } else {
                                            setEditFolderName(focusedFolderName());
                                            setIsEditingFolderName(true);
                                        }
                                    }}
                                >
                                    <Show when={isEditingFolderName()} fallback={<EditIcon />}>
                                        <SaveIcon />
                                    </Show>
                                </Button>
                            </Show>
                        </div>

                        <Show when={showCreateFolder()}>
                            <CreateFolderComponent
                                userId={user.id || EmptyGuid}
                                parentId={focusedNodeId()}
                                auth={auth}
                                setShowCreateFolder={setShowCreateFolder}
                                folderAPIRef={foldersApi}
                                refresh={async () => {
                                    await fetchRootFolders();
                                    folderRefresh();
                                    folderRefresh = () => { };
                                }}
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
                    </div>

                    <div class="treeview-scroll p-2 pt-0 overflow-y-auto flex-1 min-h-0">
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
                                            onFolderSelected={handleFolderSelected}
                                            indent={0}
                                        />
                                    ))}
                                </div>
                            </Show>
                        </Show>
                    </div>
                </div>

                <Components.WidgetContainer />
            </div>
        </div>
    );
};

export default Home;
