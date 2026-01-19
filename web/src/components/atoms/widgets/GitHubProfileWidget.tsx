import { createSignal, For, onMount, Show } from 'solid-js';
import { Configuration, WidgetsApi } from '@/api';
import type {
	ResponsesGithubWidgetCommitsDayData,
	ResponsesGithubWidgetCommitsWeekData,
	ResponsesGithubWidgetData,
} from '@/api/models';
import { useAuth } from '@/contexts/AuthContext';

interface GitHubProfileWidgetProps {
	widgetId: string;
	authToken: string;
	containerWidth: number;
	containerHeight: number;
}

// Responsive breakpoints for layout detection
const BREAKPOINTS = {
	// Width thresholds
	COMPACT_WIDTH: 150,      // Below this: avatar + name only
	MEDIUM_WIDTH: 250,       // Below this: hide bio, minimal stats
	WIDE_WIDTH: 350,         // Above this: horizontal layout if aspect ratio allows
	// Height thresholds
	COMPACT_HEIGHT: 150,     // Below this: hide commit graph and stats
	MEDIUM_HEIGHT: 250,      // Below this: hide bio
};

export default function GitHubProfileWidget(props: GitHubProfileWidgetProps) {
	const [data, setData] = createSignal<ResponsesGithubWidgetData | null>(null);
	const [loading, setLoading] = createSignal(true);
	const [error, setError] = createSignal<string | null>(null);
	const auth = useAuth();

	// Responsive layout calculations
	const width = () => props.containerWidth;
	const height = () => props.containerHeight;
	const aspectRatio = () => width() / height();

	// Layout mode: 'compact' | 'vertical' | 'wide'
	const layoutMode = () => {
		if (width() < BREAKPOINTS.COMPACT_WIDTH || height() < BREAKPOINTS.COMPACT_HEIGHT) {
			return 'compact';
		}
		if (width() >= BREAKPOINTS.WIDE_WIDTH && aspectRatio() > 1.5) {
			return 'wide';
		}
		return 'vertical';
	};

	// What to show based on available space
	const showStats = () => width() >= BREAKPOINTS.MEDIUM_WIDTH && height() >= BREAKPOINTS.COMPACT_HEIGHT;
	const showBio = () => width() >= BREAKPOINTS.MEDIUM_WIDTH && height() >= BREAKPOINTS.MEDIUM_HEIGHT;
	const showCommitGraph = () => height() >= BREAKPOINTS.COMPACT_HEIGHT && width() >= BREAKPOINTS.COMPACT_WIDTH;

	// Calculate commit graph sizing
	const commitSquareSize = () => {
		// Base size on available width, min 8px, max 12px
		const availableWidth = width() - 16; // padding
		const maxWeeks = 52;
		const idealSize = Math.floor(availableWidth / maxWeeks);
		return Math.max(8, Math.min(12, idealSize));
	};

	// Limit visible weeks based on available width
	const maxVisibleWeeks = () => {
		const squareSize = commitSquareSize();
		const gap = 2;
		const availableWidth = width() - 16; // padding
		return Math.floor(availableWidth / (squareSize + gap));
	};

	// Avatar size based on container
	const avatarSize = () => {
		if (layoutMode() === 'compact') return 'w-12 h-12';
		if (layoutMode() === 'wide') return 'w-14 h-14';
		return 'w-16 h-16';
	};

	onMount(async () => {
		try {
			const config = new Configuration({
				basePath: '/api',
			});
			const api = new WidgetsApi(config);

			console.log({'Fetching GitHub widget data...': {
				authToken: auth.token(),
				userWidgetId: props.widgetId
			}});
			const response = await api.getGithubWidgetData({
				authorization: `Bearer ${auth.token()}`,
				userWidgetId: props.widgetId,
			});

			if (response.success && response.data) {
				setData(response.data);
			} else {
				throw new Error(
					response.message || 'Failed to fetch GitHub widget data',
				);
			}
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Failed to fetch data');
		} finally {
			setLoading(false);
		}
	});

	const handleClick = () => {
		const profile = data()?.profile;
		if (profile?.htmlUrl) {
			window.open(profile.htmlUrl, '_blank');
		}
	};

	// Helper to get day data from a week
	const getDaysFromWeek = (
		week: ResponsesGithubWidgetCommitsWeekData,
	): ResponsesGithubWidgetCommitsDayData[] => {
		return [
			week.monday,
			week.tuesday,
			week.wednesday,
			week.thursday,
			week.friday,
			week.saturday,
			week.sunday,
		].filter(
			(day): day is ResponsesGithubWidgetCommitsDayData => day !== undefined,
		);
	};

	// Commit Graph Component - responsive to container size
	const CommitGraph = () => {
		const allWeeks = () => data()?.commits?.weeks || [];
		const total = () => data()?.commits?.total || 0;

		// Get only the most recent weeks that fit in available space
		const visibleWeeks = () => {
			const max = maxVisibleWeeks();
			const weeks = allWeeks();
			// Show most recent weeks (end of array)
			return weeks.slice(-Math.min(max, weeks.length));
		};

		// Responsive square size
		const squareSize = () => `${commitSquareSize()}px`;
		const gapSize = () => width() < 200 ? '1px' : '2px';

		// Extract legend colors from actual data
		const legendColors = () => {
			const allDays: ResponsesGithubWidgetCommitsDayData[] = [];
			for (const week of allWeeks()) {
				allDays.push(...getDaysFromWeek(week));
			}

			if (allDays.length === 0) return [];

			const sorted = [...allDays]
				.filter((d) => d.color && d.percent !== undefined)
				.sort((a, b) => (a.percent || 0) - (b.percent || 0));

			if (sorted.length === 0) return [];

			const indices = [
				0,
				Math.floor(sorted.length * 0.25),
				Math.floor(sorted.length * 0.5),
				Math.floor(sorted.length * 0.75),
				sorted.length - 1,
			];

			const colors: string[] = [];
			for (const idx of indices) {
				const color = sorted[Math.min(idx, sorted.length - 1)]?.color;
				if (color && !colors.includes(color)) {
					colors.push(color);
				}
			}
			return colors;
		};

		// Only show legend if there's enough width
		const showLegend = () => width() >= 200;

		return (
			<div class="w-full mt-2 flex-shrink-0">
				<div class="flex items-center justify-between mb-1">
					<span class="text-xs text-gray-400 truncate">{total()} commits</span>
					<Show when={showLegend()}>
						<div class="flex items-center gap-1 text-xs text-gray-500 flex-shrink-0">
							<span>Less</span>
							<For each={legendColors()}>
								{(color) => (
									<div
										class="rounded-sm"
										style={{
											'background-color': color,
											width: squareSize(),
											height: squareSize(),
										}}
									/>
								)}
							</For>
							<span>More</span>
						</div>
					</Show>
				</div>
				<div class="flex overflow-hidden" style={{ gap: gapSize() }}>
					<For each={visibleWeeks()}>
						{(week) => (
							<div class="flex flex-col" style={{ gap: gapSize() }}>
								<For each={getDaysFromWeek(week)}>
									{(day) => (
										<div
											class="rounded-sm cursor-pointer transition-transform hover:scale-110"
											style={{
												'background-color': day.color || '#161b22',
												width: squareSize(),
												height: squareSize(),
											}}
											title={`${day.date}: ${day.count || 0} commits`}
										/>
									)}
								</For>
							</div>
						)}
					</For>
				</div>
			</div>
		);
	};

	// Compact layout: Just avatar and name (for very small containers)
	const CompactLayout = () => (
		<div
			onClick={handleClick}
			class="h-full w-full flex flex-col items-center justify-center p-2 bg-gradient-to-br from-gray-800 to-gray-900 rounded-lg text-white cursor-pointer overflow-hidden"
		>
			<img
				src={data()?.profile?.avatarUrl}
				alt={data()?.profile?.name}
				class={`${avatarSize()} rounded-full border-2 border-primary mb-1 flex-shrink-0`}
			/>
			<h3 class="text-xs font-bold text-center truncate w-full">
				{data()?.profile?.name}
			</h3>
		</div>
	);

	// Vertical layout: Stacked content (for taller containers)
	const VerticalLayout = () => (
		<div class="h-full w-full flex flex-col items-center p-3 bg-gradient-to-br from-gray-800 to-gray-900 rounded-lg text-white overflow-hidden">
			<div
				onClick={handleClick}
				class="flex flex-col items-center cursor-pointer flex-shrink-0"
			>
				<img
					src={data()?.profile?.avatarUrl}
					alt={data()?.profile?.name}
					class={`${avatarSize()} rounded-full border-2 border-primary mb-2`}
				/>
				<h3 class="text-sm font-bold text-center truncate max-w-full">
					{data()?.profile?.name}
				</h3>
				<Show when={data()?.profile?.company}>
					<p class="text-xs text-gray-400 truncate max-w-full">
						{data()?.profile?.company}
					</p>
				</Show>
			</div>

			<Show when={showBio() && data()?.profile?.bio}>
				<p class="text-xs text-gray-300 text-center mt-2 line-clamp-2 flex-shrink-0">
					{data()?.profile?.bio}
				</p>
			</Show>

			<Show when={showStats()}>
				<div class="w-full grid grid-cols-2 gap-1 mt-2 flex-shrink-0">
					<div class="bg-black/20 rounded p-1.5 text-center">
						<div class="text-sm font-bold text-primary">
							{data()?.profile?.followers}
						</div>
						<div class="text-xs text-gray-400">Followers</div>
					</div>
					<div class="bg-black/20 rounded p-1.5 text-center">
						<div class="text-sm font-bold text-green-400">
							{data()?.profile?.publicRepos}
						</div>
						<div class="text-xs text-gray-400">Repos</div>
					</div>
				</div>
			</Show>

			<Show when={showCommitGraph()}>
				<div class="flex-1 min-h-0 w-full overflow-hidden">
					<CommitGraph />
				</div>
			</Show>
		</div>
	);

	// Wide layout: Horizontal content (for wider containers with appropriate aspect ratio)
	const WideLayout = () => (
		<div class="h-full w-full flex flex-col p-3 bg-gradient-to-r from-gray-800 to-gray-900 rounded-lg text-white overflow-hidden">
			<div class="flex flex-row items-center gap-3 flex-shrink-0">
				<img
					src={data()?.profile?.avatarUrl}
					alt={data()?.profile?.name}
					class={`${avatarSize()} rounded-full border-2 border-primary flex-shrink-0 cursor-pointer`}
					onClick={handleClick}
				/>

				<div class="flex-1 min-w-0 overflow-hidden">
					<h3
						class="text-sm font-bold truncate cursor-pointer hover:text-primary"
						onClick={handleClick}
					>
						{data()?.profile?.name}
					</h3>
					<Show when={data()?.profile?.company}>
						<p class="text-xs text-gray-400 truncate">
							{data()?.profile?.company}
						</p>
					</Show>
					<Show when={showBio() && data()?.profile?.bio}>
						<p class="text-xs text-gray-300 line-clamp-1 mt-0.5">
							{data()?.profile?.bio}
						</p>
					</Show>
				</div>

				<Show when={showStats()}>
					<div class="flex gap-2 flex-shrink-0">
						<div class="text-center">
							<div class="text-sm font-bold text-primary">
								{data()?.profile?.followers}
							</div>
							<div class="text-xs text-gray-400">Followers</div>
						</div>
						<div class="text-center">
							<div class="text-sm font-bold text-green-400">
								{data()?.profile?.publicRepos}
							</div>
							<div class="text-xs text-gray-400">Repos</div>
						</div>
					</div>
				</Show>
			</div>

			<Show when={showCommitGraph()}>
				<div class="flex-1 min-h-0 overflow-hidden">
					<CommitGraph />
				</div>
			</Show>
		</div>
	);

	return (
		<Show
			when={!loading()}
			fallback={
				<div class="h-full w-full flex items-center justify-center bg-gray-800/50 rounded-lg text-white">
					<div class="text-sm">Loading...</div>
				</div>
			}
		>
			<Show
				when={!error() && data()}
				fallback={
					<div class="h-full w-full flex flex-col items-center justify-center bg-red-900/20 rounded-lg text-white p-4">
						<div class="text-3xl mb-2">&#9888;</div>
						<div class="text-sm text-red-300">Error: {error()}</div>
					</div>
				}
			>
				<Show when={layoutMode() === 'compact'}>
					<CompactLayout />
				</Show>
				<Show when={layoutMode() === 'vertical'}>
					<VerticalLayout />
				</Show>
				<Show when={layoutMode() === 'wide'}>
					<WideLayout />
				</Show>
			</Show>
		</Show>
	);
}
