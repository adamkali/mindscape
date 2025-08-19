import { A } from '@solidjs/router';
import { createEffect, createSignal, onCleanup, onMount } from 'solid-js';
import { UsersApi } from '@/api';
import { useAuth } from '@/contexts/AuthContext';

interface TreeNode {
	id: string;
	name: string;
	type: 'folder' | 'note' | 'bookmark';
	children?: TreeNode[];
	expanded?: boolean;
	url?: string;
	content?: string;
	tags?: string[];
}

interface Mode {
	type: 'tree' | 'notes' | 'search';
	focusedNodeId?: string;
}

const Home = () => {
	const auth = useAuth();
	const [profilePicture, setProfilePicture] = createSignal<string>('');
	const [isLoadingPicture, setIsLoadingPicture] = createSignal(false);
	const [darkMode, setDarkMode] = createSignal(false);
	const [mode, setMode] = createSignal<Mode>({ type: 'tree' });
	const [searchQuery, setSearchQuery] = createSignal('');
	const [treeData, setTreeData] = createSignal<TreeNode[]>([
		{
			id: '1',
			name: 'Work',
			type: 'folder',
			expanded: true,
			children: [
				{
					id: '2',
					name: 'Project Ideas',
					type: 'note',
					content: 'Some project ideas...',
				},
				{
					id: '3',
					name: 'Documentation',
					type: 'bookmark',
					url: 'https://docs.example.com',
				},
			],
		},
		{
			id: '4',
			name: 'Personal',
			type: 'folder',
			expanded: false,
			children: [
				{
					id: '5',
					name: 'Reading List',
					type: 'note',
					content: 'Books to read...',
				},
				{ id: '6', name: 'Recipes', type: 'folder', children: [] },
			],
		},
	]);
	const [focusedNodeId, setFocusedNodeId] = createSignal<string>('1');

	const api = new UsersApi();
	const user = auth.user();

	onMount(() => {
		const savedDarkMode = localStorage.getItem('darkMode');
		if (savedDarkMode === 'true') {
			setDarkMode(true);
			document.documentElement.classList.add('dark');
		}

		document.addEventListener('keydown', handleKeyPress);
	});

	onCleanup(() => {
		document.removeEventListener('keydown', handleKeyPress);
	});

	createEffect(() => {
		if (auth.isAuthenticated() && auth.token()) {
			fetchProfilePicture();
		}
	});

	const fetchProfilePicture = async () => {
		if (!auth.token()) return;
		setIsLoadingPicture(true);
		try {
			const response = await api.getProfilePicture({
				authorization: `Bearer ${auth.token()}`,
			});
			if (response.data) {
				setProfilePicture(response.data);
			}
		} catch (error) {
			console.error('Failed to fetch profile picture:', error);
		} finally {
			setIsLoadingPicture(false);
		}
	};

	const toggleDarkMode = () => {
		const newDarkMode = !darkMode();
		setDarkMode(newDarkMode);
		localStorage.setItem('darkMode', newDarkMode.toString());
		if (newDarkMode) {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	};

	const handleLogout = () => {
		auth.logout();
	};

	const findNodeById = (nodes: TreeNode[], id: string): TreeNode | null => {
		for (const node of nodes) {
			if (node.id === id) return node;
			if (node.children) {
				const found = findNodeById(node.children, id);
				if (found) return found;
			}
		}
		return null;
	};

	const getAllNodeIds = (nodes: TreeNode[]): string[] => {
		const ids: string[] = [];
		const traverse = (nodes: TreeNode[]) => {
			for (const node of nodes) {
				ids.push(node.id);
				if (node.children && node.expanded) {
					traverse(node.children);
				}
			}
		};
		traverse(nodes);
		return ids;
	};

	const handleKeyPress = (e: KeyboardEvent) => {
		if (e.target instanceof HTMLInputElement) return;

		const currentMode = mode();

		if (e.ctrlKey && e.altKey && e.key === 'n') {
			e.preventDefault();
			setMode({ type: 'notes' });
			return;
		}

		if (e.key === 'Escape') {
			e.preventDefault();
			setMode({ type: 'tree' });
			return;
		}

		if (e.key === 'q' && currentMode.type === 'tree') {
			e.preventDefault();
			setMode({ type: 'search' });
			return;
		}

		if (currentMode.type === 'tree') {
			handleTreeNavigation(e);
		}
	};

	const handleTreeNavigation = (e: KeyboardEvent) => {
		const allIds = getAllNodeIds(treeData());
		const currentIndex = allIds.indexOf(focusedNodeId());

		switch (e.key) {
			case 'j':
			case 's':
				e.preventDefault();
				if (currentIndex < allIds.length - 1) {
					setFocusedNodeId(allIds[currentIndex + 1]);
				}
				break;
			case 'k':
			case 'w':
				e.preventDefault();
				if (currentIndex > 0) {
					setFocusedNodeId(allIds[currentIndex - 1]);
				}
				break;
			case 'h':
			case 'd': {
				e.preventDefault();
				const nodeToClose = findNodeById(treeData(), focusedNodeId());
				if (
					nodeToClose &&
					nodeToClose.type === 'folder' &&
					nodeToClose.expanded
				) {
					toggleNode(focusedNodeId());
				}
				break;
			}
			case 'l':
			case 't': {
				e.preventDefault();
				const nodeToOpen = findNodeById(treeData(), focusedNodeId());
				if (
					nodeToOpen &&
					nodeToOpen.type === 'folder' &&
					!nodeToOpen.expanded
				) {
					toggleNode(focusedNodeId());
				}
				break;
			}
			case 'Enter':
			case 'i': {
				e.preventDefault();
				const selectedNode = findNodeById(treeData(), focusedNodeId());
				if (selectedNode) {
					if (selectedNode.type === 'folder') {
						toggleNode(focusedNodeId());
					} else {
						setMode({ type: 'notes', focusedNodeId: selectedNode.id });
					}
				}
				break;
			}
		}
	};

	const toggleNode = (nodeId: string) => {
		setTreeData((prev) => {
			const updateNode = (nodes: TreeNode[]): TreeNode[] => {
				return nodes.map((node) => {
					if (node.id === nodeId) {
						return { ...node, expanded: !node.expanded };
					}
					if (node.children) {
						return { ...node, children: updateNode(node.children) };
					}
					return node;
				});
			};
			return updateNode(prev);
		});
	};

	const renderTreeNode = (node: TreeNode, depth = 0) => {
		const isFolder = node.type === 'folder';
		const isFocused = focusedNodeId() === node.id;
		const indent = depth * 20;

		return (
			<div>
				<div
					class={`flex items-center py-1 px-2 cursor-pointer transition-colors ${
						isFocused ? 'bg-primary text-primary-foreground' : 'hover:bg-card'
					}`}
					style={{ 'margin-left': `${indent}px` }}
					onClick={() => {
						setFocusedNodeId(node.id);
						if (isFolder) {
							toggleNode(node.id);
						} else {
							setMode({ type: 'notes', focusedNodeId: node.id });
						}
					}}
				>
					{isFolder && (
						<span class="mr-2 text-sm">{node.expanded ? '▼' : '▶'}</span>
					)}
					<span class="mr-2">
						{node.type === 'folder' ? '📁' : node.type === 'note' ? '📝' : '🔗'}
					</span>
					<span class="text-sm">{node.name}</span>
				</div>
				{isFolder && node.expanded && node.children && (
					<div>
						{node.children.map((child) => renderTreeNode(child, depth + 1))}
					</div>
				)}
			</div>
		);
	};

	const renderNotesMode = () => {
		const focusedNode = mode().focusedNodeId
			? findNodeById(treeData(), mode().focusedNodeId!)
			: null;

		return (
			<div class="flex-1 p-6">
				<div class="mb-4">
					<h2 class="text-xl font-bold text-foreground mb-2">
						{focusedNode ? focusedNode.name : 'Notes'}
					</h2>
					<p class="text-sm text-foreground/60">
						Press Escape to return to tree view
					</p>
				</div>

				{focusedNode && (
					<div class="bg-card rounded-lg p-4 min-h-96">
						{focusedNode.type === 'note' ? (
							<textarea
								class="w-full h-80 p-4 bg-background text-foreground border border-card-foreground/20 rounded resize-none focus:outline-none focus:ring-2 focus:ring-primary"
								placeholder="Start writing your note..."
								value={focusedNode.content || ''}
							/>
						) : focusedNode.type === 'bookmark' ? (
							<div>
								<div class="mb-4">
									<label class="block text-sm font-medium text-foreground mb-1">
										URL
									</label>
									<input
										type="url"
										class="w-full p-2 bg-background text-foreground border border-card-foreground/20 rounded focus:outline-none focus:ring-2 focus:ring-primary"
										value={focusedNode.url || ''}
										placeholder="https://..."
									/>
								</div>
								<div>
									<label class="block text-sm font-medium text-foreground mb-1">
										Description
									</label>
									<textarea
										class="w-full h-40 p-4 bg-background text-foreground border border-card-foreground/20 rounded resize-none focus:outline-none focus:ring-2 focus:ring-primary"
										placeholder="Bookmark description..."
										value={focusedNode.content || ''}
									/>
								</div>
							</div>
						) : null}
					</div>
				)}
			</div>
		);
	};

	const renderSearchMode = () => {
		return (
			<div class="flex-1 p-6">
				<div class="mb-4">
					<h2 class="text-xl font-bold text-foreground mb-2">Search</h2>
					<input
						type="text"
						class="w-full p-3 bg-background text-foreground border border-card-foreground/20 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
						placeholder="Search notes and bookmarks..."
						value={searchQuery()}
						onInput={(e) => setSearchQuery(e.currentTarget.value)}
						autofocus
					/>
				</div>

				<div class="space-y-2">
					{/* Placeholder search results */}
					<div class="p-3 bg-card rounded border-l-4 border-primary">
						<div class="font-medium text-foreground">Project Ideas</div>
						<div class="text-sm text-foreground/60">Note</div>
					</div>
				</div>
			</div>
		);
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
			{/* Header */}
			<div class="border-b border-card-foreground/20 bg-card">
				<div class="flex items-center justify-between p-4">
					<h1 class="text-2xl font-bold text-foreground">Mindscape</h1>

					<div class="flex items-center space-x-4">
						{/* Dark mode toggle */}
						<button
							onClick={toggleDarkMode}
							class="p-2 rounded-lg bg-background hover:bg-background/80 text-foreground transition-colors"
							title="Toggle dark mode"
						>
							{darkMode() ? '☀️' : '🌙'}
						</button>

						{/* Profile section */}
						<div class="flex items-center space-x-3">
							<div class="w-8 h-8 rounded-full overflow-hidden bg-card-foreground/20 flex items-center justify-center">
								{isLoadingPicture() ? (
									<div class="text-xs text-foreground/60">...</div>
								) : profilePicture() ? (
									<img
										src={profilePicture()}
										alt={`${user.username}'s profile`}
										class="w-full h-full object-cover"
									/>
								) : (
									<div class="text-sm text-foreground/60">
										{user.username?.charAt(0).toUpperCase()}
									</div>
								)}
							</div>

							<span class="text-sm text-foreground">{user.username}</span>

							<A
								href="/edit-profile"
								class="text-xs px-2 py-1 bg-secondary text-secondary-foreground rounded hover:bg-secondary/80 transition-colors"
							>
								Edit
							</A>

							<button
								onClick={handleLogout}
								class="text-xs px-2 py-1 bg-background text-foreground rounded border border-card-foreground/20 hover:bg-background/80 transition-colors"
							>
								Logout
							</button>
						</div>
					</div>
				</div>
			</div>

			{/* Main content */}
			<div class="flex h-screen">
				{/* Sidebar with tree */}
				<div class="w-80 border-r border-card-foreground/20 bg-card overflow-y-auto">
					<div class="p-4">
						<div class="flex items-center justify-between mb-4">
							<h2 class="font-semibold text-foreground">
								{mode().type === 'tree'
									? 'Tree View'
									: mode().type === 'notes'
										? 'Notes Mode'
										: 'Search Mode'}
							</h2>
							<div class="text-xs text-foreground/60">
								{mode().type === 'tree' &&
									'(j/k: navigate, l/h: expand/collapse, i: select)'}
							</div>
						</div>

						{mode().type === 'tree' && (
							<div class="space-y-1">
								{treeData().map((node) => renderTreeNode(node))}
							</div>
						)}
					</div>
				</div>

				{/* Main content area */}
				{mode().type === 'notes' && renderNotesMode()}
				{mode().type === 'search' && renderSearchMode()}
				{mode().type === 'tree' && (
					<div class="flex-1 flex items-center justify-center">
						<div class="text-center text-foreground/60">
							<div class="text-4xl mb-4">📁</div>
							<p class="text-lg">Select a note or bookmark to view</p>
							<p class="text-sm mt-2">
								Press 'q' to search, Ctrl+Alt+N for notes mode
							</p>
						</div>
					</div>
				)}
			</div>
		</div>
	);
};

export default Home;
