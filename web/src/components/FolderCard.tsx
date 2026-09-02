import type { ComponentProps } from 'solid-js';
import type { ResponsesFolderData } from '@/api';
import { useDraggableNode, useDropTarget } from '@/hooks/useDragNode';
import type { DragNodePayload } from '@/utils/dragNode';
import { Button, Card } from './atoms';
import { AddBookmarkIcon, DeleteIcon } from './icons';

interface FolderCardProps
	extends Omit<ComponentProps<'div'>, 'onDrop' | 'onSelect'> {
	folder: ResponsesFolderData;
	isSelected?: boolean;
	onSelect?: (folderId: string) => void;
	onDelete?: (folderId: string) => void;
	onCreateBookmark?: (folderId: string) => void;
	onDrop?: (payload: DragNodePayload) => void;
}

export default function FolderCard(props: FolderCardProps) {
	const { isDragging, draggableProps } = useDraggableNode(() => ({
		type: 'folder',
		id: props.folder.id ?? '',
		name: props.folder.name ?? undefined,
	}));
	const { isDragOver, dropProps } = useDropTarget((payload) =>
		props.onDrop?.(payload),
	);

	const cardClasses = () => {
		let classes =
			'w-64 hover:scale-105 active:scale-95 cursor-pointer transition-all duration-300';

		if (isDragOver()) {
			classes += ' ring-2 ring-blue-400 bg-blue-100/20';
		}

		if (isDragging()) {
			classes += ' opacity-50';
		}

		return classes;
	};

	return (
		<Card
			variant="glass"
			class={cardClasses()}
			{...draggableProps}
			{...dropProps}
			onClick={(e) => {
				e.stopPropagation();
				if (props.onSelect && props.folder.id) {
					props.onSelect(props.folder.id);
				}
			}}
		>
			<div class="flex items-center justify-between px-4">
				<div class="flex items-center flex-1">
					<span class="text-base font-bold truncate">{props.folder.name}</span>
				</div>

				<div class="flex items-center space-x-1 flex-shrink-0">
					{props.onCreateBookmark && (
						<div class="relative group">
							<Button
								variant="secondary"
								class="p-1 text-xs"
								onClick={(e) => {
									e.preventDefault();
									e.stopPropagation();
									if (props.onCreateBookmark && props.folder.id) {
										props.onCreateBookmark(props.folder.id);
									}
								}}
							>
								<AddBookmarkIcon />
							</Button>
							<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 text-xs text-white bg-black/80 backdrop-blur-sm rounded-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-50">
								Create Bookmark
							</div>
						</div>
					)}
					{props.onDelete && (
						<div class="relative group">
							<Button
								variant="danger"
								class="p-1 text-xs"
								onClick={(e) => {
									e.preventDefault();
									e.stopPropagation();
									if (props.onDelete && props.folder.id) {
										props.onDelete(props.folder.id);
									}
								}}
							>
								<DeleteIcon />
							</Button>
							<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 text-xs text-white bg-black/80 backdrop-blur-sm rounded-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-50">
								Delete Folder
							</div>
						</div>
					)}
				</div>
			</div>
		</Card>
	);
}
