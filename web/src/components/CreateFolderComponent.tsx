import type { FoldersApi } from '@/api';
import type { AuthContextValue } from '@/contexts/AuthContext';
import { createSignal, type ComponentProps } from 'solid-js';

interface CreateFolderComponentProps extends ComponentProps<'div'> {
	userId: string;
	parentId: string | undefined;
	auth: AuthContextValue | undefined;
	setShowCreateFolder: (show: boolean) => void;
	folderAPIRef: FoldersApi;
}

export default function CreateFolderComponent(
	props: CreateFolderComponentProps,
) {
	const { userId, parentId, auth , setShowCreateFolder} = props;
	const api = props.folderAPIRef;
	const [folderName, setFolderName] = createSignal('');
	const [folderDescription, setFolderDescription] = createSignal('');

	const create = async (event: Event) => {
		event.preventDefault();
		if (!folderName()) return;
		console.log('create folder', folderName());
		const response = await api.createFolder({
			createFolderRequest: {
				userId,
				parentId,
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
		<div class="mb-4 p-3 bg-background rounded border">
			<input
				type="text"
				placeholder="Enter folder name"
				value={folderName()}
				onInput={(e) => setFolderName(e.currentTarget.value)}
				class="w-full p-2 text-sm border border-card-foreground/20 rounded mb-2 focus:outline-none focus:ring-2 focus:ring-primary"
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
					class="text-xs px-2 py-1 bg-primary text-primary-foreground rounded hover:bg-primary/80"
				>
					Create
				</button>
				<button
					onClick={() => {
						setShowCreateFolder(false);
						setFolderName('');
					}}
					class="text-xs px-2 py-1 bg-background border border-card-foreground/20 rounded hover:bg-background/80"
				>
					Cancel
				</button>
			</div>
		</div>
	);
}
