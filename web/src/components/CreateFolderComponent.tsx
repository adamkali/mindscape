import { type ComponentProps, createSignal } from 'solid-js';
import { useTree } from '@/contexts/TreeContext';

interface CreateFolderComponentProps extends ComponentProps<'div'> {
	close: () => void;
}

export default function CreateFolderComponent(
	props: CreateFolderComponentProps,
) {
	const tree = useTree();
	const [folderName, setFolderName] = createSignal('');
	const [folderDescription, setFolderDescription] = createSignal('');

	const reset = () => {
		setFolderName('');
		setFolderDescription('');
	};

	const create = async (event: Event) => {
		event.preventDefault();
		if (!folderName()) return;
		const ok = await tree.createFolder({
			name: folderName(),
			description: folderDescription(),
		});
		if (ok) {
			reset();
			props.close();
		}
	};

	const cancel = () => {
		reset();
		props.close();
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
						cancel();
					}
				}}
			/>
			<div class="flex space-x-2">
				<button
					type="button"
					onClick={create}
					class="text-xs px-3 py-1.5 bg-glass-bg backdrop-blur-md border border-glass-border text-foreground rounded-lg hover:bg-glass-bg-hover transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
				>
					Create
				</button>
				<button
					type="button"
					onClick={cancel}
					class="text-xs px-3 py-1.5 bg-glass-bg/75 backdrop-blur-md border border-glass-border/85 text-foreground rounded-lg hover:bg-glass-bg-hover transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
				>
					Cancel
				</button>
			</div>
		</div>
	);
}
