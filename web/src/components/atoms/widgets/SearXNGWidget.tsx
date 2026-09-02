import { createSignal, For, Show } from 'solid-js';

// SearXNG widget — direct client-side fetch to the configured SearXNG instance.
//
// Design decision — client-side direct:
//   The SearXNG server is user-owned and typically on the local network or a
//   trusted domain. No sensitive credentials are involved, so a server-side
//   proxy is unnecessary overhead. CORS and JSON format must be configured on
//   the SearXNG instance (see requirements below).
//
// SearXNG instance requirements:
//   1. In settings.yml add:  search: { formats: [html, json] }
//   2. Set CORS to allow the Mindscape origin, or set it to '*' for internal use.
//      Without this the browser will block the request.

interface SearXNGWidgetProps {
	serverUrl: string;
	defaultEngines?: string[];
	categories?: string[];
	safeSearch?: boolean;
}

interface SearXNGResult {
	title: string;
	url: string;
	content?: string;
	engine?: string;
}

interface SearXNGResponse {
	results: SearXNGResult[];
	query: string;
	number_of_results?: number;
}

export default function SearXNGWidget(props: SearXNGWidgetProps) {
	const [query, setQuery] = createSignal('');
	const [results, setResults] = createSignal<SearXNGResult[]>([]);
	const [loading, setLoading] = createSignal(false);
	const [error, setError] = createSignal<string | null>(null);
	const [searched, setSearched] = createSignal(false);

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		const q = query().trim();
		if (!q || !props.serverUrl) return;

		const params = new URLSearchParams({ q, format: 'json' });

		const engines = props.defaultEngines ?? [];
		if (engines.length > 0) {
			params.set('engines', engines.join(','));
		}

		const cats = props.categories ?? [];
		if (cats.length > 0) {
			params.set('categories', cats.join(','));
		}

		params.set('safesearch', props.safeSearch ? '1' : '0');

		const url = `${props.serverUrl}/search?${params.toString()}`;

		setLoading(true);
		setError(null);
		setSearched(true);

		try {
			const resp = await fetch(url, {
				headers: { Accept: 'application/json' },
			});
			if (!resp.ok) {
				setError(
					`SearXNG returned ${resp.status}. Check that the server URL is correct.`,
				);
				return;
			}
			const data = (await resp.json()) as SearXNGResponse;
			setResults(data.results ?? []);
		} catch (err) {
			const message = err instanceof Error ? err.message : String(err);
			setError(
				`Request failed: ${message}. Ensure your SearXNG instance has format=json enabled and allows this origin in CORS settings.`,
			);
		} finally {
			setLoading(false);
		}
	};

	return (
		<div class="w-full flex flex-col gap-3 p-4 bg-glass-bg rounded-lg text-foreground">
			<form class="w-full flex gap-2" onSubmit={handleSubmit}>
				<label class="sr-only" for="searxng-input">
					Search SearXNG
				</label>
				<input
					id="searxng-input"
					type="text"
					placeholder="🔍 Search your SearXNG..."
					class="flex-1 px-4 py-2 rounded-lg bg-black/20 text-foreground focus:outline-none focus:ring-1 focus:ring-primary"
					value={query()}
					onInput={(e) => setQuery(e.currentTarget.value)}
				/>
				<button
					type="submit"
					class="px-4 py-2 rounded-lg bg-primary text-white font-medium hover:opacity-90 transition-opacity"
					disabled={loading()}
				>
					{loading() ? '...' : 'Search'}
				</button>
			</form>

			<Show when={error()}>
				<div class="text-sm text-red-300 bg-red-900/20 rounded p-3">
					{error()}
				</div>
			</Show>

			<Show when={searched() && !loading() && !error()}>
				<Show
					when={results().length > 0}
					fallback={
						<p class="text-sm text-glass-text-muted">No results found.</p>
					}
				>
					<div class="flex flex-col gap-2 overflow-y-auto max-h-80">
						<For each={results()}>
							{(result) => (
								<a
									href={result.url}
									target="_blank"
									rel="noopener noreferrer"
									class="flex flex-col gap-0.5 p-3 rounded-lg bg-black/20 hover:bg-black/30 transition-colors no-underline"
								>
									<span class="text-sm font-medium text-primary truncate">
										{result.title}
									</span>
									<span class="text-xs text-glass-text-muted truncate">
										{result.url}
									</span>
									<Show when={result.content}>
										<span class="text-xs text-glass-text line-clamp-2">
											{result.content}
										</span>
									</Show>
								</a>
							)}
						</For>
					</div>
				</Show>
			</Show>
		</div>
	);
}
