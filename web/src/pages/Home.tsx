import { A } from '@solidjs/router';
import {
	createEffect,
	createSignal,
	For,
	onCleanup,
	onMount,
	Show,
} from 'solid-js';
import type { RepositoryBookmark } from '@/api';
import Components from '@/components';
import AgendaContainer from '@/components/AgendaContainer';
import { Button, Input } from '@/components/atoms';
import CreateBookmarkComponent from '@/components/CreateBookmarkComponent';
import CreateFolderComponent from '@/components/CreateFolderComponent';
import EditBookmarkModal from '@/components/EditBookmarkModal';
import FolderComponent from '@/components/FolderComponent';
import { Header } from '@/components/Header';
import { AddFolder, EditIcon, SaveIcon } from '@/components/icons';
import ModePill from '@/components/ModePill';
import SearchOverlay from '@/components/SearchOverlay';
import { FEATURES } from '@/config/features';
import { useAuth } from '@/contexts/AuthContext';
import {
	KeyboardNavProvider,
	useKeyboardNav,
} from '@/contexts/KeyboardNavContext';
import { TreeProvider, useTree } from '@/contexts/TreeContext';
import { useView, ViewProvider } from '@/contexts/ViewContext';
import { WidgetProvider } from '@/contexts/WidgetContext';
import { useBackgroundStyle } from '@/hooks/useBackground';
import { EmptyGuid } from '@/utils';
import { parseNode } from '@/utils/dragNode';

const Home = () => {
	return (
		<KeyboardNavProvider>
			<TreeProvider>
				<ViewProvider>
					<HomeInner />
				</ViewProvider>
			</TreeProvider>
		</KeyboardNavProvider>
	);
};

const HomeInner = () => {
	const auth = useAuth();
	const { activeView, homepageMode, setHomepageMode } = useView();
	const nav = useKeyboardNav();
	const tree = useTree();
	const backgroundStyle = useBackgroundStyle();
	const [isLoadingFolders, setIsLoadingFolders] = createSignal(false);
	const [showCreateFolder, setShowCreateFolder] = createSignal(false);
	const [showCreateBookmark, setShowCreateBookmark] = createSignal(false);
	const [isDragOverRoot, setIsDragOverRoot] = createSignal(false);
	const [focusedFolderName, setFocusedFolderName] = createSignal('');
	const [isEditingFolderName, setIsEditingFolderName] = createSignal(false);
	const [editFolderName, setEditFolderName] = createSignal('');
	const [editingBookmark, setEditingBookmark] =
		createSignal<RepositoryBookmark | null>(null);
	const [showEditBookmark, setShowEditBookmark] = createSignal(false);
	const user = auth.user();

	// Movement key handler — registered into the keyboard nav context.
	// Runs only when mode === 'nav'; delegates to sub-mode handlers in TASK-011.
	onMount(() => {
		const unregister = nav.register((e: KeyboardEvent) => {
			if (nav.mode() !== 'nav' && nav.mode() !== 'move') return false;

			const nodes = nav.getNodes();
			const currentId = tree.focusedId();
			const currentIdx = currentId
				? nodes.findIndex((n) => n.id === currentId)
				: -1;
			const current = currentIdx >= 0 ? nodes[currentIdx] : undefined;

			if (e.key === 'j' || e.key === 'ArrowDown') {
				e.preventDefault();
				if (nodes.length === 0) return true;
				if (currentIdx < 0) {
					tree.setFocus(nodes[0].id);
					nodes[0].ref.scrollIntoView({ block: 'nearest' });
				} else {
					const next = nodes[Math.min(currentIdx + 1, nodes.length - 1)];
					tree.setFocus(next.id);
					next.ref.scrollIntoView({ block: 'nearest' });
				}
				return true;
			}

			if (e.key === 'k' || e.key === 'ArrowUp') {
				e.preventDefault();
				if (nodes.length === 0) return true;
				if (currentIdx <= 0) {
					tree.setFocus(nodes[0].id);
					nodes[0].ref.scrollIntoView({ block: 'nearest' });
				} else {
					const prev = nodes[currentIdx - 1];
					tree.setFocus(prev.id);
					prev.ref.scrollIntoView({ block: 'nearest' });
				}
				return true;
			}

			if (e.key === 'l' || e.key === 'ArrowRight') {
				e.preventDefault();
				if (!current) return true;
				if (current.type === 'folder') {
					if (current.isOpen?.()) {
						// Folder already open — move focus into first visible child.
						const firstChild = nodes.find((n) => n.parentId === current.id);
						if (firstChild) {
							tree.setFocus(firstChild.id);
							firstChild.ref.scrollIntoView({ block: 'nearest' });
						}
					} else {
						current.expand?.();
					}
				} else if (current.type === 'bookmark' && current.link) {
					window.open(current.link, '_blank');
				}
				return true;
			}

			if (e.key === 'h' || e.key === 'ArrowLeft') {
				e.preventDefault();
				if (!current) return true;
				if (current.type === 'folder' && current.isOpen?.()) {
					current.collapse?.();
				} else if (current.parentId) {
					const parent = nodes.find((n) => n.id === current.parentId);
					if (parent) {
						tree.setFocus(parent.id);
						parent.ref.scrollIntoView({ block: 'nearest' });
					}
				}
				return true;
			}

			if (e.key === 'g') {
				e.preventDefault();
				if (nodes.length > 0) {
					tree.setFocus(nodes[0].id);
					nodes[0].ref.scrollIntoView({ block: 'nearest' });
				}
				return true;
			}

			if (e.key === 'G') {
				e.preventDefault();
				if (nodes.length > 0) {
					const last = nodes[nodes.length - 1];
					tree.setFocus(last.id);
					last.ref.scrollIntoView({ block: 'nearest' });
				}
				return true;
			}

			return false;
		});
		onCleanup(unregister);
	});

	createEffect(() => {
		if (!auth.token()) return;
		setIsLoadingFolders(true);
		tree.refreshRoot().finally(() => setIsLoadingFolders(false));
	});

	createEffect(() => {
		if (tree.focusedId() === EmptyGuid) {
			tree.setFocus('');
		}
		if (showCreateFolder()) {
			setShowCreateBookmark(false);
		}
		if (showCreateBookmark()) {
			setShowCreateFolder(false);
		}
	});

	// Keep focusedFolderName in sync when keyboard nav moves focus to a folder.
	createEffect(() => {
		const id = tree.focusedId();
		if (!id) return;
		const node = nav.getNodes().find((n) => n.id === id);
		if (node?.type === 'folder' && node.name) {
			setFocusedFolderName(node.name);
		}
	});

	const handleEditBookmark = (bookmark: RepositoryBookmark) => {
		setEditingBookmark(bookmark);
		setShowEditBookmark(true);
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

	const openCreateBookmark = () => {
		setShowCreateBookmark(true);
		setShowCreateFolder(false);
	};

	const openCreateFolder = () => {
		setShowCreateFolder(true);
		setShowCreateBookmark(false);
	};

	const handleFolderSelected = (folderName: string) => {
		setFocusedFolderName(folderName);
		setIsEditingFolderName(false);
	};

	const handleSaveFolderName = async () => {
		const folderId = tree.focusedId();
		const newName = editFolderName().trim();
		if (!folderId || !newName) return;

		const ok = await tree.updateFolder(folderId, { name: newName });
		if (ok) {
			setFocusedFolderName(newName);
			setIsEditingFolderName(false);
		}
	};

	const handleRootDragOver = (e: DragEvent) => {
		// Only handle if the target is the root container or empty space, not folder cards
		const target = e.target as HTMLElement;
		const isOverFolderCard = target.closest('[draggable="true"]') !== null;

		if (!isOverFolderCard) {
			e.preventDefault();
			if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
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

	const handleRootDrop = (e: DragEvent) => {
		// Only handle if the target is actually the root container or empty space
		const target = e.target as HTMLElement;
		const isOverFolderCard = target.closest('[draggable="true"]') !== null;

		if (isOverFolderCard) {
			// Let the folder card handle this drop
			return;
		}

		e.preventDefault();
		setIsDragOverRoot(false);
		const payload = parseNode(e);
		if (payload) tree.move(payload.id, null);
	};

	// Action key handler — i/d/c/s and their sub-modes (delete/create/move).
	// Registered here so openCreateBookmarkComponent / openCreateFolderComponent are in scope.
	onMount(() => {
		const unregister = nav.register((e: KeyboardEvent) => {
			const currentMode = nav.mode();

			// ── delete sub-mode ──────────────────────────────────────────────
			if (currentMode === 'delete') {
				if (e.key === 'd') {
					e.preventDefault();
					const id = tree.focusedId();
					if (id) {
						const node = nav.getNodes().find((n) => n.id === id);
						if (node?.type === 'folder') {
							tree.deleteFolder(id);
						} else if (node?.type === 'bookmark') {
							tree.deleteBookmark(id);
						}
					}
					nav.setMode('nav');
					return true;
				}
				if (e.key === 'c') {
					e.preventDefault();
					nav.setMode('nav');
					return true;
				}
				// Swallow all other keys so stray presses don't leak out.
				e.preventDefault();
				return true;
			}

			// ── create sub-mode ──────────────────────────────────────────────
			if (currentMode === 'create') {
				if (e.key === 'b') {
					e.preventDefault();
					openCreateBookmark();
					nav.setMode('nav');
					return true;
				}
				if (e.key === 'f') {
					e.preventDefault();
					openCreateFolder();
					nav.setMode('nav');
					return true;
				}
				// Swallow every other key so a stray press can't leak into nav mode.
				e.preventDefault();
				return true;
			}

			// ── move sub-mode ('p' to place; movement keys handled by TASK-010) ──
			if (currentMode === 'move') {
				if (e.key === 'p') {
					e.preventDefault();
					const sourceId = nav.moveSourceId();
					if (sourceId) tree.move(sourceId, tree.focusedId() || null);
					nav.cancelMove();
					return true;
				}
				// Let movement keys fall through to the TASK-010 handler.
				return false;
			}

			// ── nav mode actions ─────────────────────────────────────────────
			if (currentMode === 'nav') {
				if (e.key === 'i') {
					e.preventDefault();
					const id = tree.focusedId();
					if (id) {
						const node = nav.getNodes().find((n) => n.id === id);
						if (node?.type === 'folder') {
							setEditFolderName(focusedFolderName());
							setIsEditingFolderName(true);
						} else if (node?.type === 'bookmark') {
							node.onEdit?.();
						}
					}
					return true;
				}
				if (e.key === 'd') {
					e.preventDefault();
					if (tree.focusedId()) nav.setMode('delete');
					return true;
				}
				if (e.key === 'c') {
					e.preventDefault();
					nav.setMode('create');
					return true;
				}
				if (e.key === 's') {
					e.preventDefault();
					const id = tree.focusedId();
					if (id) nav.startMove(id);
					return true;
				}
				// Layout switching. 'n' is deliberately unbound — it was Notes mode,
				// which shipped cut; 'c n' (create note) is gone with it.
				if (e.key === 't') {
					e.preventDefault();
					setHomepageMode('tree');
					return true;
				}
				if (e.key === 'r') {
					e.preventDefault();
					setHomepageMode('normal');
					return true;
				}
			}

			return false;
		});
		onCleanup(unregister);
	});

	return (
		<main
			class="h-screen overflow-hidden bg-background flex flex-col"
			style={backgroundStyle()}
			onClick={(e) => {
				// Deselect node when clicking on background areas
				const target = e.target as HTMLElement;
				if (
					target.classList.contains('bg-background') ||
					target.closest('.treeview-container') === null
				) {
					tree.setFocus('');
				}
			}}
			onKeyDown={(e) => {
				if (e.key === 'Escape') tree.setFocus('');
			}}
		>
			<Header />

			<div class="flex flex-1 min-h-0 flex-row">
				<section
					aria-label="Folder tree"
					class={`treeview-container m-2 p-4 rounded-lg flex flex-col overflow-hidden bg-background backdrop-blur-lg border border-white/20 max-h-[calc(100vh-2rem)] shadow-2xl shadow-slate-900/30 dark:border-slate-700/50 dark:shadow-black/30 ${
						homepageMode() === 'tree' ? 'flex-1' : 'min-w-80 max-w-80'
					} ${isDragOverRoot() ? 'ring-2 ring-blue-400 bg-blue-100/20' : ''}`}
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
									openCreateFolder();
								}}
							>
								<AddFolder />
							</Button>

							<Show when={tree.focusedId()}>
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
							<CreateFolderComponent close={() => setShowCreateFolder(false)} />
						</Show>

						<Show when={showCreateBookmark()}>
							<CreateBookmarkComponent
								close={() => setShowCreateBookmark(false)}
							/>
						</Show>
					</div>

					<div class="treeview-scroll-wrapper relative flex-1 min-h-0">
						<div class="treeview-scroll p-2 pt-0 overflow-y-auto h-full">
							<Show
								when={!isLoadingFolders()}
								fallback={
									<div class="text-center py-4 text-foreground/60">
										Loading folders...
									</div>
								}
							>
								<Show
									when={tree.rootFolders().length > 0}
									fallback={
										<div class="text-center py-4 text-foreground/60">
											No folders yet. Create your first folder!
										</div>
									}
								>
									<div class="space-y-4">
										<For each={tree.rootFolders()}>
											{(folder) => (
												<FolderComponent
													folder={folder}
													editBookmark={handleEditBookmark}
													openCreateBookmark={openCreateBookmark}
													onFolderSelected={handleFolderSelected}
													indent={0}
												/>
											)}
										</For>
									</div>
								</Show>
							</Show>
						</div>
					</div>
				</section>

				<div
					class={
						homepageMode() === 'normal' && activeView() === 'widgets'
							? 'w-full flex'
							: 'hidden'
					}
				>
					<WidgetProvider>
						<Components.WidgetContainer />
					</WidgetProvider>
				</div>
				<Show when={FEATURES.tasks}>
					<div
						class={
							homepageMode() === 'normal' && activeView() === 'agenda'
								? 'w-full flex'
								: 'hidden'
						}
					>
						<AgendaContainer />
					</div>
				</Show>
			</div>

			<EditBookmarkModal
				isOpen={showEditBookmark()}
				onClose={() => setShowEditBookmark(false)}
				bookmark={editingBookmark()}
			/>

			<ModePill />
			<SearchOverlay />
		</main>
	);
};

export default Home;
