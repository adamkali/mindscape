import type { AuthContextValue } from '@/contexts/AuthContext';
import { createSignal, type ComponentProps } from 'solid-js';
import {
	BookmarksApi,
	type CreateBookmarkRequest,
	type RepositoryCreateBookmarkParams,
} from '@/api';

interface CreateBookmarkComponentProps extends ComponentProps<'div'> {
	userId: string;
	parentId: string | undefined;
	auth: AuthContextValue | undefined;
	close: () => void;
	refreshBookmarks?: () => void;
}

export default function CreateBookmarkComponent(
	props: CreateBookmarkComponentProps,
) {
	const { userId, parentId, auth, close, refreshBookmarks } = props;
	const [linkName, setLinkName] = createSignal('');
	const [linkUrl, setLinkUrl] = createSignal('');
	const api = new BookmarksApi();

	const create = async (event: Event) => {
		event.preventDefault();
		const response = await api.createBookmark({
			authorization: `Bearer ${auth?.token()}`,
			createBookmarkRequest: {
				userId,
				folderId: parentId,
				link: linkUrl(),
				name: linkName(),
			} as RepositoryCreateBookmarkParams,
		} as CreateBookmarkRequest);

		if (response.success && response.data) {
			close();
			refreshBookmarks?.();
		} else {
			console.error('Failed to create bookmark:', response.message);
		}
		return;
	};

	return (
		<div class="mb-4 p-4 bg-white/10 backdrop-blur-md border border-white/20 rounded-xl shadow-lg">
			<div class="mb-3">
				<h3 class="text-sm font-medium text-white/90 mb-3">Create Bookmark</h3>
			</div>
			
			<input
				type="url"
				placeholder="Enter bookmark URL"
				value={linkUrl()}
				onInput={(e) => setLinkUrl(e.currentTarget.value)}
				class="w-full p-3 text-sm bg-white/20 backdrop-blur-md border border-white/30 rounded-lg mb-3 focus:outline-none focus:border-white/50 focus:bg-white/25
				placeholder:text-white/60 text-white transition-all duration-200 shadow-sm hover:shadow-md"
				autofocus
				required
				onKeyDown={(e) => {
					if (e.key === 'Enter' && linkUrl() && linkName()) {
						create(e);
					} else if (e.key === 'Escape') {
						close();
						setLinkUrl('');
						setLinkName('');
					}
				}}
			/>
			
			<input
				type="text"
				placeholder="Enter bookmark name"
				value={linkName()}
				onInput={(e) => setLinkName(e.currentTarget.value)}
				class="w-full p-3 text-sm bg-white/20 backdrop-blur-md border border-white/30 rounded-lg mb-3 focus:outline-none focus:border-white/50 focus:bg-white/25
				placeholder:text-white/60 text-white transition-all duration-200 shadow-sm hover:shadow-md"
				onKeyDown={(e) => {
					if (e.key === 'Enter' && linkUrl() && linkName()) {
						create(e);
					} else if (e.key === 'Escape') {
						close();
						setLinkUrl('');
						setLinkName('');
					}
				}}
			/>
			
			<div class="flex space-x-2">
				<button
					onClick={create}
					disabled={!linkUrl() || !linkName()}
					class="text-xs px-3 py-1.5 bg-white/20 backdrop-blur-md border border-white/30 text-white rounded-lg hover:bg-white/30 transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 disabled:hover:bg-white/20"
				>
					Create
				</button>
				<button
					onClick={() => {
						close();
						setLinkUrl('');
						setLinkName('');
					}}
					class="text-xs px-3 py-1.5 bg-white/15 backdrop-blur-md border border-white/25 text-white rounded-lg hover:bg-white/25 transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95"
				>
					Cancel
				</button>
			</div>
		</div>
	);
}
