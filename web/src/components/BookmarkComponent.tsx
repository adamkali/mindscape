import type { ComponentProps } from 'solid-js';
import type { RepositoryBookmark } from '@/api';
import BookmarkCard from './BookmarkCard';

interface BookmarkComponentProps extends ComponentProps<'div'> {
	bookmark: RepositoryBookmark;
	indent: number;
	selected: () => string;
	setSelected: (id: string) => void;
	deleteBookmark?: (bookmarkId: string, parentFolderId: string) => void;
}

export default function BookmarkComponent(props: BookmarkComponentProps) {
	const bookmarkClassName = (): string => {
		let isSelected = '';
		if (props.bookmark.id === props.selected()) {
			isSelected =
				'bg-gray-300/40 backdrop-blur-sm text-white hover:bg-gray-300/50';
		}
		return cn(
			'flex flex-1 items-center py-3 px-4 cursor-pointer bg-gray-500/30 backdrop-blur-sm text-white hover:bg-gray-500/40 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out w-80',
			isSelected,
		);
	};

	return (
		<div style={{ 'margin-left': `${props.indent * 2}rem` }}>
			<BookmarkCard
				bookmark={props.bookmark}
				onDelete={props.deleteBookmark}
				draggable={true}
			/>
		</div>
	);
}
