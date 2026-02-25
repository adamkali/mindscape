import { A } from '@solidjs/router';
import { type ComponentProps, createMemo, createSignal } from 'solid-js';
import type { RepositoryBookmark } from '@/api';
import { Button, Card } from './atoms';
import { DeleteIcon, EditIcon } from './icons';

interface BookmarkCardProps extends ComponentProps<'div'> {
	bookmark: RepositoryBookmark;
	onDelete?: (bookmarkId: string, parentFolderId: string) => void;
	onEdit?: (bookmark: RepositoryBookmark) => void;
	draggable?: boolean;
}

const faviconeUrl = 'https://favicone.com/';
const getFaviconUrl = (url: string): string => {
	try {
		const urlObj = new URL(url);
		return `${faviconeUrl}${urlObj.hostname}`;
	} catch {
		return '';
	}
};

export default function BookmarkCard(props: BookmarkCardProps) {
	const faviconUrl = createMemo(() => getFaviconUrl(props.bookmark.link || ''));
	const [isDragging, setIsDragging] = createSignal(false);

	const handleDragStart = (e: DragEvent) => {
		if (!props.draggable) return;

		setIsDragging(true);
		e.dataTransfer!.setData(
			'text/plain',
			JSON.stringify({
				type: 'bookmark',
				id: props.bookmark.id,
				name: props.bookmark.name,
				link: props.bookmark.link,
			}),
		);
		e.dataTransfer!.effectAllowed = 'move';
	};

	const handleDragEnd = () => {
		setIsDragging(false);
	};

	return (
		<Card
			variant="glass"
			class={`w-64 hover:scale-105 active:scale-95 cursor-pointer ${isDragging() ? 'opacity-50' : ''}`}
			draggable={props.draggable}
			onDragStart={handleDragStart}
			onDragEnd={handleDragEnd}
			onClick={(e) => e.stopPropagation()}
		>
			<div class="flex items-center justify-between px-4">
				<A
					class="flex flex-1 items-center"
					href={props.bookmark.link || ''}
					rel="noopener noreferrer"
					target="_blank"
				>
					{faviconUrl() && (
						<img
							src={faviconUrl()}
							alt="favicon"
							class="w-4 h-4 mr-3 rounded-sm"
							onError={(e) => {
								(e.target as HTMLImageElement).style.display = 'none';
							}}
						/>
					)}
					<span class="text-base font-bold truncate">
						{props.bookmark.name}
					</span>
				</A>

				{props.onEdit && (
					<Button
						variant="secondary"
						class="p-1 ml-2 text-xs flex-shrink-0"
						onClick={(e) => {
							e.preventDefault();
							e.stopPropagation();
							props.onEdit?.(props.bookmark);
						}}
					>
						<EditIcon />
					</Button>
				)}

				{props.onDelete && (
					<Button
						variant="danger"
						class="p-1 ml-2 text-xs flex-shrink-0"
						onClick={(e) => {
							e.preventDefault();
							e.stopPropagation();
							if (
								props.onDelete &&
								props.bookmark.id &&
								props.bookmark.folderId
							) {
								props.onDelete(props.bookmark.id, props.bookmark.folderId);
							}
						}}
					>
						<DeleteIcon />
					</Button>
				)}
			</div>
		</Card>
	);
}
