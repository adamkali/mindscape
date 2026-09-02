import { createSignal, For, onCleanup, onMount, Show } from 'solid-js';
import { WidgetsApi } from '@/api';
import type { PlexMediaItemData } from '@/api/models';
import { getAuthenticatedApiConfig } from '@/utils/apiConfig';

interface PlexWidgetProps {
	widgetId: string;
	authToken: string;
	/** Refresh interval in seconds; defaults to 3600 (1 hour). */
	refreshInterval?: number;
}

function formatAddedAt(addedAt: number | undefined): string {
	if (!addedAt) return '';
	return new Date(addedAt * 1000).toLocaleDateString(undefined, {
		month: 'short',
		day: 'numeric',
		year: 'numeric',
	});
}

export default function PlexWidget(props: PlexWidgetProps) {
	const [items, setItems] = createSignal<PlexMediaItemData[]>([]);
	const [loading, setLoading] = createSignal(true);
	const [error, setError] = createSignal<string | null>(null);

	const fetchItems = () => {
		const api = new WidgetsApi(getAuthenticatedApiConfig());

		api
			.getPlexRecentlyAdded({
				userWidgetId: props.widgetId,
			})
			.then((response) => {
				if (response.success && response.data) {
					setItems(response.data);
					setError(null);
				} else {
					setError(response.message ?? 'Failed to load Plex data');
				}
			})
			.catch((err: Error) => setError(err.message))
			.finally(() => setLoading(false));
	};

	onMount(() => {
		fetchItems();

		const intervalSeconds = props.refreshInterval ?? 3600;
		const intervalId = setInterval(fetchItems, intervalSeconds * 1000);

		onCleanup(() => clearInterval(intervalId));
	});

	return (
		<Show
			when={!loading()}
			fallback={
				<div class="w-full flex items-center justify-center bg-gray-800/50 rounded-lg text-white p-8">
					<div class="text-sm">Loading Plex recently added...</div>
				</div>
			}
		>
			<Show
				when={!error()}
				fallback={
					<div class="w-full flex flex-col items-center justify-center bg-red-900/20 rounded-lg text-white p-8">
						<div class="text-3xl mb-2">&#9888;</div>
						<div class="text-sm text-red-300">Error: {error()}</div>
					</div>
				}
			>
				<div class="w-full p-4 bg-gradient-to-br bg-glass-bg rounded-lg text-foreground">
					<h3 class="text-sm font-semibold text-glass-text-muted mb-3 uppercase tracking-wide">
						Recently Added
					</h3>
					<Show
						when={items().length > 0}
						fallback={
							<p class="text-sm text-glass-text-muted">No items found.</p>
						}
					>
						<div class="grid grid-cols-3 gap-3 overflow-y-auto max-h-96">
							<For each={items()}>
								{(item) => (
									<div class="flex flex-col gap-1">
										<Show
											when={item.thumb}
											fallback={
												<div class="w-full aspect-[2/3] bg-gray-700 rounded flex items-center justify-center">
													<span class="text-xs text-gray-400">No art</span>
												</div>
											}
										>
											<img
												src={item.thumb}
												alt={item.title ?? 'Plex media item'}
												class="w-full aspect-[2/3] object-cover rounded"
												loading="lazy"
											/>
										</Show>
										<p
											class="text-xs font-medium text-foreground truncate"
											title={item.title}
										>
											{item.title}
										</p>
										<Show when={item.type === 'episode' && item.parentTitle}>
											<p class="text-xs text-glass-text-muted truncate">
												{item.parentTitle}
											</p>
										</Show>
										<p class="text-xs text-glass-text-muted">
											{formatAddedAt(item.addedAt)}
										</p>
									</div>
								)}
							</For>
						</div>
					</Show>
				</div>
			</Show>
		</Show>
	);
}
