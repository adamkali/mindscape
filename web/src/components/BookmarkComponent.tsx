import { type ComponentProps, onCleanup, onMount } from 'solid-js';
import type { RepositoryBookmark } from '@/api';
import { useKeyboardNav } from '@/contexts/KeyboardNavContext';
import { useTree } from '@/contexts/TreeContext';
import BookmarkCard from './BookmarkCard';

interface BookmarkComponentProps extends ComponentProps<'div'> {
	bookmark: RepositoryBookmark;
	indent: number;
	editBookmark?: (bookmark: RepositoryBookmark) => void;
}

export default function BookmarkComponent(props: BookmarkComponentProps) {
	const nav = useKeyboardNav();
	const tree = useTree();
	let rowRef: HTMLDivElement | undefined;
	const bookmark = props.bookmark;
	const editFn = props.editBookmark;

	onMount(() => {
		if (!bookmark.id || !rowRef) return;
		const unregister = nav.registerNode({
			id: bookmark.id,
			type: 'bookmark',
			ref: rowRef,
			parentId: bookmark.folderId ?? null,
			name: bookmark.name ?? '',
			link: bookmark.link ?? undefined,
			onEdit: editFn ? () => editFn(bookmark) : undefined,
		});
		onCleanup(unregister);
	});

	const isFocused = () => tree.focusedId() === props.bookmark.id;

	return (
		<div
			ref={(r) => {
				rowRef = r;
			}}
			tabIndex={-1}
			style={{ 'margin-left': `${props.indent * 2}rem` }}
			class={
				isFocused()
					? 'rounded-lg ring-2 ring-white/50 dark:bg-white/20 bg-black/10'
					: ''
			}
		>
			<BookmarkCard
				bookmark={props.bookmark}
				onDelete={(bookmarkId) => tree.deleteBookmark(bookmarkId)}
				onEdit={props.editBookmark}
			/>
		</div>
	);
}
