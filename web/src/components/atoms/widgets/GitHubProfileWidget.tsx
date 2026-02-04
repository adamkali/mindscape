import { createSignal, For, onMount, Show } from 'solid-js';
import { Configuration, WidgetsApi } from '@/api';
import type {
	ResponsesGithubWidgetCommitsDayData,
	ResponsesGithubWidgetCommitsWeekData,
	ResponsesGithubWidgetCommitsData,
	ResponsesGithubWidgetProfileData,
} from '@/api/models';
import { useAuth } from '@/contexts/AuthContext';

interface GitHubProfileWidgetProps {
	widgetId: string;
	authToken: string;
}

export default function GitHubProfileWidget(props: GitHubProfileWidgetProps) {
	const [profile, setProfile] = createSignal<ResponsesGithubWidgetProfileData | null>(null);
	const [commits, setCommits] = createSignal<ResponsesGithubWidgetCommitsData | null>(null);
	const [profileLoading, setProfileLoading] = createSignal(true);
	const [commitsLoading, setCommitsLoading] = createSignal(true);
	const [profileError, setProfileError] = createSignal<string | null>(null);
	const [commitsError, setCommitsError] = createSignal<string | null>(null);
	const auth = useAuth();

	onMount(async () => {
		const config = new Configuration({
			basePath: '/api',
		});
		const api = new WidgetsApi(config);

		// Fetch profile (fast)
		api.getGithubProfileWidgetData({
			authorization: `Bearer ${auth.token()}`,
			userWidgetId: props.widgetId,
		}).then(response => {
			if (response.success && response.data) {
				setProfile(response.data);
			} else {
				setProfileError(response.message || 'Failed to load profile');
			}
		}).catch(err => setProfileError(err.message))
		  .finally(() => setProfileLoading(false));

		// Fetch commits (slow, loads in background)
		api.getGithubCommitsWidgetData({
			authorization: `Bearer ${auth.token()}`,
			userWidgetId: props.widgetId,
		}).then(response => {
			if (response.success && response.data) {
				setCommits(response.data);
			} else {
				setCommitsError(response.message || 'Failed to load commits');
			}
		}).catch(err => setCommitsError(err.message))
		  .finally(() => setCommitsLoading(false));
	});

	const handleClick = () => {
		const profileData = profile();
		if (profileData?.htmlUrl) {
			window.open(profileData.htmlUrl, '_blank');
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

	// Commit Graph Component - uses CSS for responsiveness
	const CommitGraph = () => {
		const allWeeks = () => commits()?.weeks || [];
		const total = () => commits()?.total || 0;

		// Show all weeks, let CSS handle overflow
		const visibleWeeks = () => allWeeks();

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

		return (
			<div class="w-full mt-3">
				<Show
					when={!commitsLoading()}
					fallback={
						<div class="flex flex-col gap-2">
							<div class="flex items-center justify-between mb-2">
								<span class="text-sm text-gray-400">Loading commits...</span>
							</div>
							<div class="flex gap-0.5 overflow-x-auto pb-2">
								<For each={Array(12).fill(null)}>
									{() => (
										<div class="flex flex-col gap-0.5">
											<For each={Array(7).fill(null)}>
												{() => (
													<div
														class="w-3 h-3 rounded-sm bg-gray-700 animate-pulse"
													/>
												)}
											</For>
										</div>
									)}
								</For>
							</div>
						</div>
					}
				>
					<Show
						when={!commitsError()}
						fallback={
							<div class="text-sm text-red-300">
								Failed to load commits: {commitsError()}
							</div>
						}
					>
						<div class="flex items-center justify-between mb-2">
							<span class="text-sm text-gray-400">{total()} commits this year</span>
							<div class="flex items-center gap-1 text-xs text-gray-500">
								<span>Less</span>
								<For each={legendColors()}>
									{(color) => (
										<div
											class="w-3 h-3 rounded-sm"
											style={{ 'background-color': color }}
										/>
									)}
								</For>
								<span>More</span>
							</div>
						</div>
						<div class="flex gap-0.5 overflow-x-auto pb-2">
							<For each={visibleWeeks()}>
								{(week) => (
									<div class="flex flex-col gap-0.5">
										<For each={getDaysFromWeek(week)}>
											{(day) => (
												<div
													class="w-3 h-3 rounded-sm cursor-pointer transition-transform hover:scale-110"
													style={{ 'background-color': day.color || '#161b22' }}
													title={`${day.date}: ${day.count || 0} commits`}
												/>
											)}
										</For>
									</div>
								)}
							</For>
						</div>
					</Show>
				</Show>
			</div>
		);
	};

	return (
		<Show
			when={!profileLoading()}
			fallback={
				<div class="w-full flex items-center justify-center bg-gray-800/50 rounded-lg text-white p-8">
					<div class="text-sm">Loading profile...</div>
				</div>
			}
		>
			<Show
				when={!profileError() && profile()}
				fallback={
					<div class="w-full flex flex-col items-center justify-center bg-red-900/20 rounded-lg text-white p-8">
						<div class="text-3xl mb-2">&#9888;</div>
						<div class="text-sm text-red-300">Error: {profileError()}</div>
					</div>
				}
			>
				<div class="w-full p-4 bg-gradient-to-br bg-glass-bg rounded-lg text-foreground">
					{/* Profile Header */}
					<div class="flex items-center gap-4 mb-4">
						<img
							src={profile()?.avatarUrl}
							alt={profile()?.name}
							class="w-16 h-16 rounded-full border-2 border-primary cursor-pointer flex-shrink-0"
							onClick={handleClick}
						/>
						<div class="flex-1 min-w-0">
							<h3
								class="text-lg font-bold truncate cursor-pointer hover:text-primary"
								onClick={handleClick}
							>
								{profile()?.name}
							</h3>
							<Show when={profile()?.company}>
								<p class="text-sm text-glass-text-muted truncate">
									{profile()?.company}
								</p>
							</Show>
							<Show when={profile()?.bio}>
								<p class="text-sm text-glass-text line-clamp-2 mt-1">
									{profile()?.bio}
								</p>
							</Show>
						</div>
					</div>

					{/* Stats Row */}
					<div class="grid grid-cols-4 gap-2 mb-4">
						<div class="bg-black/20 rounded p-2 text-center">
							<div class="text-lg font-bold text-primary">
								{profile()?.followers}
							</div>
							<div class="text-xs text-gray-400">Followers</div>
						</div>
						<div class="bg-black/20 rounded p-2 text-center">
							<div class="text-lg font-bold text-blue-400">
								{profile()?.following}
							</div>
							<div class="text-xs text-gray-400">Following</div>
						</div>
						<div class="bg-black/20 rounded p-2 text-center">
							<div class="text-lg font-bold text-green-400">
								{profile()?.publicRepos}
							</div>
							<div class="text-xs text-gray-400">Repos</div>
						</div>
						<div class="bg-black/20 rounded p-2 text-center">
							<div class="text-lg font-bold text-yellow-400">
								{profile()?.publicGists}
							</div>
							<div class="text-xs text-gray-400">Gists</div>
						</div>
					</div>

					{/* Commit Graph */}
					<CommitGraph />
				</div>
			</Show>
		</Show>
	);
}
