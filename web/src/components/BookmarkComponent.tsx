import type { RepositoryBookmark } from '@/api';
import { A } from '@solidjs/router';
import type { ComponentProps } from 'solid-js';
import { Button } from './atoms';
import { DeleteIcon } from './icons';

interface BookmarkComponentProps extends ComponentProps<'div'> {
	bookmark: RepositoryBookmark;
	indent: number;
	selected: () => string;
	setSelected: (id: string) => void;
	deleteBookmark?: (bookmarkId: string, parentFolderId: string) => void;
}
export default function BookmarkComponent(props: BookmarkComponentProps) {
	return (
		<div class="flex items-center justify-between py-1 px-2 bg-white/20 backdrop-blur-sm border border-white/30 text-white hover:bg-white/30 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out w-64">
			<A
				class="flex flex-1 items-center cursor-pointer"
				href={props.bookmark.link || ''}
				rel="noopener noreferrer"
				target="_blank"
			>
				<span class="mr-2 text-base font-bold">{props.bookmark.name}</span>
			</A>
			{props.deleteBookmark && (
				<Button
					variant="danger"
					class="p-1 ml-2 text-xs"
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
						if (props.deleteBookmark && props.bookmark.id && props.bookmark.folderId) {
							props.deleteBookmark(props.bookmark.id, props.bookmark.folderId);
						}
					}}
				>
					<DeleteIcon />
				</Button>
			)}
		</div>
	);
}
