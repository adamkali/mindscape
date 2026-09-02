import { createSignal, For, onCleanup, onMount, Show } from 'solid-js';
import { WidgetsApi } from '@/api';
import type { GithubPRData, GithubPRsWidgetData } from '@/api/models';
import { getAuthenticatedApiConfig } from '@/utils/apiConfig';

interface GitHubPRsWidgetProps {
	widgetId: string;
	/** Refresh interval in seconds; defaults to 300 (5 minutes). */
	refreshInterval?: number;
}

/** Returns a human-readable relative time string, e.g. "3d ago". */
function relativeAge(isoDate: string | undefined): string {
	if (!isoDate) return '';
	const ms = Date.now() - new Date(isoDate).getTime();
	const minutes = Math.floor(ms / 60_000);
	if (minutes < 60) return `${minutes}m ago`;
	const hours = Math.floor(minutes / 60);
	if (hours < 24) return `${hours}h ago`;
	const days = Math.floor(hours / 24);
	return `${days}d ago`;
}

interface PRSectionProps {
	title: string;
	items: GithubPRData[] | undefined;
}

function PRSection(props: PRSectionProps) {
	const items = () => props.items ?? [];
	return (
		<Show when={items().length > 0}>
			<div class="flex flex-col gap-1">
				<h4 class="text-xs font-semibold uppercase tracking-wide text-glass-text-muted px-1">
					{props.title}
					<span class="ml-1 text-primary">{items().length}</span>
				</h4>
				<For each={items()}>
					{(pr) => (
						<a
							href={pr.htmlUrl}
							target="_blank"
							rel="noopener noreferrer"
							class="flex items-start gap-2 px-2 py-1.5 rounded-lg bg-black/20 hover:bg-black/30 transition-colors no-underline"
						>
							<span class="text-xs text-glass-text-muted shrink-0 pt-0.5 w-14 truncate">
								{pr.repo?.split('/')[1] ?? pr.repo}
							</span>
							<span class="text-xs text-primary shrink-0 pt-0.5">
								#{pr.number}
							</span>
							<span class="text-xs text-foreground flex-1 line-clamp-1">
								{pr.title}
							</span>
							<span class="text-xs text-glass-text-muted shrink-0 pt-0.5">
								{pr.author}
							</span>
							<span class="text-xs text-glass-text-muted shrink-0 pt-0.5 w-12 text-right">
								{relativeAge(pr.createdAt)}
							</span>
						</a>
					)}
				</For>
			</div>
		</Show>
	);
}

export default function GitHubPRsWidget(props: GitHubPRsWidgetProps) {
	const [data, setData] = createSignal<GithubPRsWidgetData | null>(null);
	const [loading, setLoading] = createSignal(true);
	const [error, setError] = createSignal<string | null>(null);

	const fetchPRs = () => {
		const api = new WidgetsApi(getAuthenticatedApiConfig());

		api
			.getGithubPRsWidgetData({
				userWidgetId: props.widgetId,
			})
			.then((response) => {
				if (response.success && response.data) {
					setData(response.data);
					setError(null);
				} else {
					setError(response.message ?? 'Failed to load PRs');
				}
			})
			.catch((err: Error) => setError(err.message))
			.finally(() => setLoading(false));
	};

	onMount(() => {
		fetchPRs();

		const intervalSeconds = props.refreshInterval ?? 300;
		const id = setInterval(fetchPRs, intervalSeconds * 1000);
		onCleanup(() => clearInterval(id));
	});

	const hasAnyPRs = () => {
		const d = data();
		if (!d) return false;
		return (
			(d.reviewRequested?.length ?? 0) > 0 ||
			(d.authored?.length ?? 0) > 0 ||
			(d.mentioned?.length ?? 0) > 0 ||
			(d.assigned?.length ?? 0) > 0
		);
	};

	return (
		<Show
			when={!loading()}
			fallback={
				<div class="w-full flex items-center justify-center bg-gray-800/50 rounded-lg text-white p-6">
					<span class="text-sm">Loading pull requests...</span>
				</div>
			}
		>
			<Show
				when={!error()}
				fallback={
					<div class="w-full flex flex-col items-center justify-center bg-red-900/20 rounded-lg text-white p-6">
						<div class="text-3xl mb-2">&#9888;</div>
						<div class="text-sm text-red-300">Error: {error()}</div>
					</div>
				}
			>
				<div class="w-full flex flex-col gap-3 p-4 bg-glass-bg rounded-lg text-foreground overflow-y-auto max-h-full">
					<Show
						when={hasAnyPRs()}
						fallback={
							<p class="text-sm text-glass-text-muted">
								No open pull requests.
							</p>
						}
					>
						<PRSection
							title="Review Requested"
							items={data()?.reviewRequested}
						/>
						<PRSection title="Authored" items={data()?.authored} />
						<PRSection title="Mentioned" items={data()?.mentioned} />
						<PRSection title="Assigned" items={data()?.assigned} />
					</Show>
				</div>
			</Show>
		</Show>
	);
}
