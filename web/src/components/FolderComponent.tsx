import type { ResponsesFolderData } from "@/api";
import { createEffect, type ComponentProps } from "solid-js";

interface FolderComponentProps extends ComponentProps<'div'> {
	folder: ResponsesFolderData
	selectedFolder: (id: string) => void
	deleteFolder: (id: string) => void
}

export default function FolderComponent(props: FolderComponentProps) {
	const { folder, selectedFolder, deleteFolder } = props
	createEffect(() => {
		console.log(folder)
		
	})

	return <div
		class='flex items-center py-1 px-2 cursor-pointer transition-colors bg-primary text-primary-foreground hover:bg-primary/80'
		onClick={() => selectedFolder(folder.id || '')}
	>
		<span class="mr-2 text-sm">📁</span>
		<span class="mr-2 text-sm">{folder.name}</span>
		<button
			onClick={(e) => {
				e.stopPropagation();
				deleteFolder(folder.id || '');
			}}
			class="ml-2 text-red-500 hover:text-red-700 text-xs px-1"
			title="Delete folder"
		>
			🗑
		</button>
	</div>
}
