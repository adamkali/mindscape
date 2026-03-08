import type { ComponentProps } from 'solid-js';
import type { RepositoryBookmark } from '@/api';
import BookmarkCard from './BookmarkCard';

interface BookmarkComponentProps extends ComponentProps<'div'> {
	bookmark: RepositoryBookmark;
	indent: number;
	selected: () => string;
	setSelected: (id: string) => void;
	deleteBookmark?: (bookmarkId: string, parentFolderId: string) => void;
	editBookmark?: (bookmark: RepositoryBookmark) => void;
}

export default function BookmarkComponent(props: BookmarkComponentProps) {
	return (
		<div style={{ 'margin-left': `${props.indent * 2}rem` }}>
			<BookmarkCard
				bookmark={props.bookmark}
				onDelete={props.deleteBookmark}
				onEdit={props.editBookmark}
				draggable={true}
			/>
		</div>
	);
}
