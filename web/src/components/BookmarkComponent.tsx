import type { RepositoryBookmark } from '@/api';
import { A } from '@solidjs/router';
import type { ComponentProps } from 'solid-js';

interface BookmarkComponentProps extends ComponentProps<'div'> {
	bookmark: RepositoryBookmark;
	indent: number;
	selected: () => string;
	setSelected: (id: string) => void;
}
export default function BookmarkComponent(props: BookmarkComponentProps) {
	return (
		<A
			class="flex flex-1 items-center py-1 px-2 cursor-pointer bg-primary text-primary-foreground hover:bg-primary/80 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out"
			href={props.bookmark.link || ''}
			rel="noopener noreferrer"
			target="_blank"
		>
			<span class="mr-2 text-base font-bold">{props.bookmark.name}</span>
		</A>
	);
}
