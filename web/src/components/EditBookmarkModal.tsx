import { createEffect, createSignal, Show } from 'solid-js';
import { BookmarksApi, type RepositoryBookmark } from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Input } from './atoms';

interface EditBookmarkModalProps {
	isOpen: boolean;
	onClose: () => void;
	bookmark: RepositoryBookmark | null;
	onSaved: () => void;
}

export default function EditBookmarkModal(props: EditBookmarkModalProps) {
	const [name, setName] = createSignal('');
	const [link, setLink] = createSignal('');
	const [saving, setSaving] = createSignal(false);
	const [error, setError] = createSignal('');
	const auth = useAuth();

	createEffect(() => {
		if (props.bookmark) {
			setName(props.bookmark.name || '');
			setLink(props.bookmark.link || '');
			setError('');
		}
	});

	const handleSave = async (e: Event) => {
		e.preventDefault();
		if (!props.bookmark?.id || !auth.token()) return;

		setError('');
		setSaving(true);
		try {
			const bookmarksApi = new BookmarksApi();
			const response = await bookmarksApi.updateBookmark({
				bookmarkId: props.bookmark.id,
				updateBookmarkRequest: {
					userId: auth.user()?.id,
					bookmarkId: props.bookmark.id,
					name: name(),
					link: link(),
				},
				authorization: `Bearer ${auth.token()}`,
			});

			if (response.success) {
				props.onSaved();
				props.onClose();
			} else {
				setError(response.message || 'Failed to update bookmark');
			}
		} catch (err) {
			console.error('Error updating bookmark:', err);
			setError('Failed to update bookmark');
		} finally {
			setSaving(false);
		}
	};

	const handleBackdropClick = (e: MouseEvent) => {
		if (e.target === e.currentTarget) {
			props.onClose();
		}
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Escape') {
			props.onClose();
		}
	};

	return (
		<Show when={props.isOpen}>
			<div
				role="dialog"
				class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50"
				onClick={handleBackdropClick}
				onKeyDown={handleKeyDown}
			>
				<div
					role="document"
					class="bg-gradient-to-br from-card to-card/80 backdrop-blur-lg border-2 border-slate-700/20 rounded-2xl shadow-2xl w-full max-w-md p-6"
					onClick={(e) => e.stopPropagation()}
					onKeyDown={(e) => e.stopPropagation()}
				>
					<div class="flex justify-between items-center mb-6">
						<h2 class="text-2xl font-bold text-card-foreground">
							Edit Bookmark
						</h2>
						<button
							type="button"
							onClick={props.onClose}
							class="text-card-foreground/60 hover:text-card-foreground transition-colors"
						>
							<svg
								class="w-6 h-6"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
								role="img"
								aria-label="Close"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
						</button>
					</div>

					<form onSubmit={handleSave} class="space-y-4">
						<Input
							label="Name"
							type="text"
							value={name()}
							onInput={(e) => setName(e.currentTarget.value)}
						/>
						<Input
							label="Link"
							type="url"
							value={link()}
							onInput={(e) => setLink(e.currentTarget.value)}
						/>

						<Show when={error()}>
							<p class="text-sm text-red-400">{error()}</p>
						</Show>

						<div class="flex gap-3 pt-4">
							<Button
								type="button"
								variant="tertiary"
								onClick={props.onClose}
								class="flex-1"
							>
								Cancel
							</Button>
							<Button
								type="submit"
								variant="primary"
								class="flex-1"
								disabled={saving() || !name() || !link()}
							>
								{saving() ? 'Saving...' : 'Save'}
							</Button>
						</div>
					</form>
				</div>
			</div>
		</Show>
	);
}
