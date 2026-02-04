import { type ComponentProps, createSignal } from 'solid-js';
import type { FoldersApi } from '@/api';
import type { AuthContextValue } from '@/contexts/AuthContext';
import { EmptyGuid } from '@/utils';

interface CreateFolderComponentProps extends ComponentProps<'div'> {
	userId: string;
	parentId: string | undefined;
	auth: AuthContextValue | undefined;
	setShowCreateFolder: (show: boolean) => void;
	folderAPIRef: FoldersApi;
	refresh?: () => void;
}

export default function CreateFolderComponent(
	props: CreateFolderComponentProps,
) {
	const { userId, parentId, auth, setShowCreateFolder, folderAPIRef, refresh } =
		props;
	const [folderName, setFolderName] = createSignal('');
	const [folderDescription, setFolderDescription] = createSignal('');

	const create = async (event: Event) => {
		event.preventDefault();
		if (!folderName()) return;
		console.log('create folder', folderName());
		console.log('parent id', parentId === '' ? 'null' : parentId);

		const response = await folderAPIRef.createFolder({
			createFolderRequest: {
				userId,
				parentId: parentId === '' ? undefined : parentId,
				name: folderName(),
				description: folderDescription(),
			},
			authorization: `Bearer ${auth?.token()}`,
		});

		if (response.success) {
			setFolderName('');
			setFolderDescription('');
		} else {
			console.error('Failed to create folder:', response.message);
			return;
		}

		return response;
	};

	return (
		<div class="mb-4 p-4 bg-white/10 backdrop-blur-md border border-white/20 rounded-xl shadow-lg">
			<input
				type="text"
				placeholder="Enter folder name"
				value={folderName()}
				onInput={(e) => setFolderName(e.currentTarget.value)}
				class="w-full p-3 text-sm bg-glass-bg backdrop-blur-md border border-glass-border rounded-lg mb-3 focus:outline-none focus:border-glass-border-hover focus:bg-glass-bg-hover
				placeholder:text-foreground/60 text-foreground transition-all duration-200 shadow-sm hover:shadow-md"
				autofocus
				onKeyDown={(e) => {
					if (e.key === 'Enter') {
						create(e);
					} else if (e.key === 'Escape') {
						setShowCreateFolder(false);
						setFolderName('');
					}
				}}
			/>
			<div class="flex space-x-2">
				<button
					onClick={create}
					class="text-xs px-3 py-1.5 bg-glass-bg backdrop-blur-md border border-glass-border text-foreground rounded-lg hover:bg-glass-bg-hover transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
				>
					Create
				</button>
				<button
					onClick={() => {
						setShowCreateFolder(false);
						setFolderName('');
					}}
					class="text-xs px-3 py-1.5 bg-glass-bg/75 backdrop-blur-md border border-glass-border/85 text-foreground rounded-lg hover:bg-glass-bg-hover transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
				>
					Cancel
				</button>
			</div>
		</div>
	);
}
