import type { ResponsesFolderData } from '@/api';
import { createSignal, type ComponentProps } from 'solid-js';
import { Button, Card } from './atoms';
import { DeleteIcon } from './icons';

interface FolderCardProps extends ComponentProps<'div'> {
	folder: ResponsesFolderData;
	isSelected?: boolean;
	onSelect?: (folderId: string) => void;
	onDelete?: (folderId: string) => void;
	draggable?: boolean;
	onDrop?: (dragData: any) => void;
}

export default function FolderCard(props: FolderCardProps) {
	const [isDragging, setIsDragging] = createSignal(false);
	const [isDragOver, setIsDragOver] = createSignal(false);

	const handleDragStart = (e: DragEvent) => {
		if (!props.draggable) return;
		
		setIsDragging(true);
		e.dataTransfer!.setData(
			'text/plain',
			JSON.stringify({
				type: 'folder',
				id: props.folder.id,
				name: props.folder.name,
			}),
		);
		e.dataTransfer!.effectAllowed = 'move';
	};

	const handleDragEnd = () => {
		setIsDragging(false);
	};

	const handleDragOver = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation(); // Prevent root container from handling this event
		e.dataTransfer!.dropEffect = 'move';
		setIsDragOver(true);
	};

	const handleDragLeave = () => {
		setIsDragOver(false);
	};

	const handleDrop = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation(); // Prevent event from bubbling up to root container
		setIsDragOver(false);

		try {
			const data = JSON.parse(e.dataTransfer!.getData('text/plain'));
			if (props.onDrop) {
				props.onDrop(data);
			}
		} catch (error) {
			console.error('Error parsing drop data:', error);
		}
	};

	const cardClasses = () => {
		let classes = 'w-64 hover:scale-105 active:scale-95 cursor-pointer transition-all duration-300';
		
		if (props.isSelected) {
			classes += ' ring-2 ring-white/50 bg-white/40';
		}
		
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
			draggable={props.draggable}
			onDragStart={handleDragStart}
			onDragEnd={handleDragEnd}
			onDragOver={handleDragOver}
			onDragLeave={handleDragLeave}
			onDrop={handleDrop}
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
				
				{props.onDelete && (
					<Button
						variant="danger"
						class="p-1 ml-2 text-xs flex-shrink-0"
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
				)}
			</div>
		</Card>
	);
}
