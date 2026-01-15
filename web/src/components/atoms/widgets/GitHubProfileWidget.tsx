import { createSignal, onMount, Show } from 'solid-js';

interface GitHubProfileWidgetProps {
    username?: string;
    showStats?: boolean;
    personalAccessToken?: string;
    variant?: 'sm' | 'lg' | 'wide';
}

interface GitHubProfile {
    login: string;
    name: string;
    avatar_url: string;
    bio: string;
    followers: number;
    following: number;
    public_repos: number;
    public_gists: number;
}

export default function GitHubProfileWidget(props: GitHubProfileWidgetProps) {
    const [profile, setProfile] = createSignal<GitHubProfile | null>(null);
    const [loading, setLoading] = createSignal(true);
    const [error, setError] = createSignal<string | null>(null);

    const username = () => props.username || 'github';
    const variant = () => props.variant || 'sm';
    const showStats = () => props.showStats ?? true;

    onMount(async () => {
        try {
            // GitHub API: https://api.github.com/users/{username}
            const headers: HeadersInit = {
                'Accept': 'application/vnd.github.v3+json',
            };

            // Add PAT to headers if provided for higher rate limits
            if (props.personalAccessToken) {
                headers['Authorization'] = `Bearer ${props.personalAccessToken}`;
            }

            const response = await fetch(`https://api.github.com/users/${username()}`, {
                headers,
            });

            if (!response.ok) {
                throw new Error(`GitHub API error: ${response.status}`);
            }

            const data = await response.json();
            setProfile(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch profile');
        } finally {
            setLoading(false);
        }
    });

    const handleClick = () => {
        window.open(`https://github.com/${profile()?.login}`, '_blank')

    }

    // Small variant: Avatar + Name only (minimal)
    const SmallVariant = () => (
        <div onClick={handleClick} class="h-full w-full flex flex-col items-center justify-center p-2 bg-gradient-to-br from-gray-800 to-gray-900 rounded-lg text-white">
            <img
                src={profile()?.avatar_url}
                alt={profile()?.name}
                class="w-16 h-16 rounded-full border-2 border-primary mb-2"
            />
            <h3 class="text-sm font-bold text-center truncate w-full">
                {profile()?.name || profile()?.login}
            </h3>
            <p class="text-xs text-gray-400">@{profile()?.login}</p>
        </div>
    );

    // Large variant: Full stats and details (vertical)
    const LargeVariant = () => (
        <div onClick={handleClick} class="h-full w-full flex flex-col items-center p-4 bg-gradient-to-br from-gray-800 to-gray-900 rounded-lg text-white">
            <img
                src={profile()?.avatar_url}
                alt={profile()?.name}
                class="w-24 h-24 rounded-full border-2 border-primary mb-3"
            />
            <h3 class="text-lg font-bold text-center">{profile()?.name || profile()?.login}</h3>
            <p class="text-sm text-gray-400 mb-2">@{profile()?.login}</p>

            <Show when={profile()?.bio}>
                <p class="text-xs text-gray-300 text-center mb-3 line-clamp-2">{profile()?.bio}</p>
            </Show>

            <Show when={showStats()}>
                <div class="w-full grid grid-cols-2 gap-2 mt-auto">
                    <div class="bg-black/20 rounded p-2 text-center">
                        <div class="text-lg font-bold text-primary">{profile()?.followers}</div>
                        <div class="text-xs text-gray-400">Followers</div>
                    </div>
                    <div class="bg-black/20 rounded p-2 text-center">
                        <div class="text-lg font-bold text-secondary">{profile()?.following}</div>
                        <div class="text-xs text-gray-400">Following</div>
                    </div>
                    <div class="bg-black/20 rounded p-2 text-center">
                        <div class="text-lg font-bold text-green-400">{profile()?.public_repos}</div>
                        <div class="text-xs text-gray-400">Repos</div>
                    </div>
                    <div class="bg-black/20 rounded p-2 text-center">
                        <div class="text-lg font-bold text-yellow-400">{profile()?.public_gists}</div>
                        <div class="text-xs text-gray-400">Gists</div>
                    </div>
                </div>
            </Show>
        </div>
    );

    // Wide variant: Horizontal layout
    const WideVariant = () => (
        <div
            class="h-full w-full flex flex-row items-center p-3 bg-gradient-to-r from-gray-800 to-gray-900 rounded-lg text-white gap-4"
            onclick={handleClick}
        >
            <img
                src={profile()?.avatar_url}
                alt={profile()?.name}
                class="w-20 h-20 rounded-full border-2 border-primary flex-shrink-0"
            />

            <div class="flex-1 min-w-0">
                <h3 class="text-base font-bold truncate">{profile()?.name || profile()?.login}</h3>
                <p class="text-xs text-gray-400 mb-1">@{profile()?.login}</p>
                <Show when={profile()?.bio}>
                    <p class="text-xs text-gray-300 line-clamp-1">{profile()?.bio}</p>
                </Show>
            </div>

            <Show when={showStats()}>
                <div class="flex gap-3 flex-shrink-0">
                    <div class="text-center">
                        <div class="text-sm font-bold text-primary">{profile()?.followers}</div>
                        <div class="text-xs text-gray-400">Followers</div>
                    </div>
                    <div class="text-center">
                        <div class="text-sm font-bold text-secondary">{profile()?.following}</div>
                        <div class="text-xs text-gray-400">Following</div>
                    </div>
                    <div class="text-center">
                        <div class="text-sm font-bold text-green-400">{profile()?.public_repos}</div>
                        <div class="text-xs text-gray-400">Repos</div>
                    </div>
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
                when={!error() && profile()}
                fallback={
                    <div class="h-full w-full flex flex-col items-center justify-center bg-red-900/20 rounded-lg text-white p-4">
                        <div class="text-3xl mb-2">⚠️</div>
                        <div class="text-sm text-red-300">Error: {error()}</div>
                    </div>
                }
            >
                <Show when={variant() === 'sm'}>
                    <SmallVariant />
                </Show>
                <Show when={variant() === 'lg'}>
                    <LargeVariant />
                </Show>
                <Show when={variant() === 'wide'}>
                    <WideVariant />
                </Show>
            </Show>
        </Show>
    );
}
