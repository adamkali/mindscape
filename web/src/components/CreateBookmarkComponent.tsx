import type { AuthContextValue } from '@/contexts/AuthContext';
import { createSignal, type ComponentProps } from 'solid-js';
import {
	BookmarksApi,
	type CreateBookmarkRequest,
	type RepositoryCreateBookmarkParams,
} from '@/api';
import { Button, Card, CardBody, CardFooter, CardHeader } from './atoms';

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
		<Card variant="elevated">
			<CardHeader>Create Bookmark</CardHeader>
			<CardBody>
				<form onSubmit={create} class="space-y-6 flex flex-col">
					<div class="mb-4">
						<label
							for="linkUrl"
							class="block text-sm font-medium text-gray-700"
						>
							Link URL
						</label>
						<input
							required
							type="text"
							id="linkUrl"
							name="linkUrl"
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
							onChange={(e) => setLinkUrl(e.currentTarget.value)}
						/>
					</div>

					<div class="mb-4">
						<label
							for="linkName"
							class="block text-sm font-medium text-gray-700"
						>
							Give a Name
						</label>
						<input
							type="text"
							id="linkName"
							name="linkName"
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
							onChange={(e) => setLinkName(e.currentTarget.value)}
						/>
					</div>
				</form>
			</CardBody>
			<CardFooter>
				<Button variant="primary" onClick={create} type="submit">
					Create
				</Button>
			</CardFooter>
		</Card>
	);
}
