import { type ComponentProps, createSignal } from 'solid-js';
import { useTree } from '@/contexts/TreeContext';

interface CreateBookmarkComponentProps extends ComponentProps<'div'> {
	close: () => void;
}

export default function CreateBookmarkComponent(
	props: CreateBookmarkComponentProps,
) {
	const tree = useTree();
	const [linkName, setLinkName] = createSignal('');
	const [linkUrl, setLinkUrl] = createSignal('');

	const reset = () => {
		setLinkUrl('');
		setLinkName('');
	};

	const create = async (event: Event) => {
		event.preventDefault();
		const ok = await tree.createBookmark({ name: linkName(), link: linkUrl() });
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
			<div class="mb-3">
				<h3 class="text-sm font-medium text-foreground/90 mb-3">
					Create Bookmark
				</h3>
			</div>

			<input
				type="url"
				placeholder="Enter bookmark URL"
				value={linkUrl()}
				onInput={(e) => setLinkUrl(e.currentTarget.value)}
				class="w-full p-3 text-sm bg-glass-bg backdrop-blur-md border border-glass-border rounded-lg mb-3 focus:outline-none focus:border-glass-border-hover focus:bg-glass-bg-hover
				placeholder:text-foreground/60 text-foreground transition-all duration-200 shadow-sm hover:shadow-md"
				autofocus
				required
				onKeyDown={(e) => {
					if (e.key === 'Enter' && linkUrl() && linkName()) {
						create(e);
					} else if (e.key === 'Escape') {
						cancel();
					}
				}}
			/>

			<input
				type="text"
				placeholder="Enter bookmark name"
				value={linkName()}
				onInput={(e) => setLinkName(e.currentTarget.value)}
				class="w-full p-3 text-sm bg-glass-bg backdrop-blur-md border border-glass-border rounded-lg mb-3 focus:outline-none focus:border-glass-border-hover focus:bg-glass-bg-hover
				placeholder:text-foreground/60 text-foreground transition-all duration-200 shadow-sm hover:shadow-md"
				onKeyDown={(e) => {
					if (e.key === 'Enter' && linkUrl() && linkName()) {
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
					disabled={!linkUrl() || !linkName()}
					class="text-xs px-3 py-1.5 bg-glass-bg backdrop-blur-md border border-glass-border text-foreground rounded-lg hover:bg-glass-bg-hover transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 disabled:hover:bg-glass-bg"
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
