import { createSignal, For, onCleanup, onMount, Show } from 'solid-js';
import { Configuration, WidgetsApi } from '@/api';
import type {
	CoolifyWidgetApplication,
	CoolifyWidgetService,
} from '@/api/models';
import { useAuth } from '@/contexts/AuthContext';

interface CoolifyWidgetProps {
	widgetId: string;
	authToken: string;
	foldInitially?: boolean;
}

// Status color mapping
function getStatusColor(status: string | undefined): string {
	if (!status) return 'bg-gray-500';
	const s = status.toLowerCase();
	if (s === 'running:healthy') return 'bg-green-500';
	if (s === 'running:unknown') return 'bg-yellow-500';
	if (s === 'exited:unhealthy') return 'bg-red-500';
	return 'bg-gray-500';
}

// Check if application is running
function isRunning(status: string | undefined): boolean {
	if (!status) return false;
	return status.toLowerCase().startsWith('running');
}

export default function CoolifyWidget(props: CoolifyWidgetProps) {
	const [applications, setApplications] = createSignal<
		CoolifyWidgetApplication[]
	>([]);
	const [services, setServices] = createSignal<CoolifyWidgetService[]>([]);
	const [appsLoading, setAppsLoading] = createSignal(true);
	const [servicesLoading, setServicesLoading] = createSignal(true);
	const [appsError, setAppsError] = createSignal<string | null>(null);
	const [servicesError, setServicesError] = createSignal<string | null>(null);
	const [appsExpanded, setAppsExpanded] = createSignal(!props.foldInitially);
	const [servicesExpanded, setServicesExpanded] = createSignal(
		!props.foldInitially,
	);
	const [actionLoading, setActionLoading] = createSignal<
		Record<string, string | null>
	>({});
	const auth = useAuth();

	const fetchData = async () => {
		const config = new Configuration({
			basePath: '/api',
		});
		const api = new WidgetsApi(config);

		// Fetch applications
		api
			.getUserCoolifyApplications({
				authorization: `Bearer ${auth.token()}`,
				userWidgetId: props.widgetId,
			})
			.then((response) => {
				if (response.success && response.data) {
					setApplications(response.data);
					setAppsError(null);
				} else {
					setAppsError(response.message || 'Failed to load applications');
				}
			})
			.catch((err) => setAppsError(err.message))
			.finally(() => setAppsLoading(false));

		// Fetch services
		api
			.getUserCoolifyServices({
				authorization: `Bearer ${auth.token()}`,
				userWidgetId: props.widgetId,
			})
			.then((response) => {
				if (response.success && response.data) {
					setServices(response.data);
					setServicesError(null);
				} else {
					setServicesError(response.message || 'Failed to load services');
				}
			})
			.catch((err) => setServicesError(err.message))
			.finally(() => setServicesLoading(false));
	};

	const handleAction = async (
		appUuid: string,
		action: 'start' | 'stop' | 'restart',
	) => {
		setActionLoading((prev) => ({ ...prev, [appUuid]: action }));
		try {
			const config = new Configuration({
				basePath: '/api',
			});
			const api = new WidgetsApi(config);
			const requestParams = {
				authorization: `Bearer ${auth.token()}`,
				userWidgetId: props.widgetId,
				appUuid: appUuid,
			};

			if (action === 'start') {
				await api.startCoolifyApplication(requestParams);
			} else if (action === 'stop') {
				await api.stopCoolifyApplication(requestParams);
			} else {
				await api.restartCoolifyApplication(requestParams);
			}
			// Refresh data after action
			await fetchData();
		} catch (err) {
			console.error(`Failed to ${action} application:`, err);
		} finally {
			setActionLoading((prev) => ({ ...prev, [appUuid]: null }));
		}
	};

	onMount(() => {
		fetchData();
		// Auto-refresh every 30 seconds
		const interval = setInterval(fetchData, 30000);
		onCleanup(() => clearInterval(interval));
	});

	const truncateNameMaxLength = 20;
	const truncateName = (name: string) => {
		if (name.length > truncateNameMaxLength) {
			return name.substring(0, truncateNameMaxLength - 3) + '...';
		}
		return name;
	};

	const ApplicationItem = (app: CoolifyWidgetApplication) => {
		const appRunning = isRunning(app.status);
		const loadingAction = () => actionLoading()[app.uuid || ''];

		return (
			<div class="flex items-center justify-between p-2 bg-black/20 rounded mb-1">
				<div class="flex items-center gap-2 min-w-0 flex-1">
					<div
						class={`w-2 h-2 rounded-full flex-shrink-0 ${getStatusColor(app.status)}`}
					/>
					<div class="min-w-0 flex-1">
						<div class="font-medium text-sm truncate">
							{truncateName(app.name || 'Unnamed')}
						</div>
						<Show when={app.fqdn}>
							<a
								href={app.fqdn}
								target="_blank"
								rel="noopener noreferrer"
								class="text-xs text-primary hover:text-primary-hover -blue-300 truncate block"
							>
								{app.fqdn}
							</a>
						</Show>
					</div>
				</div>
				<div class="flex gap-1 flex-shrink-0 ml-2">
					<Show when={appRunning}>
						<button
							type="button"
							onClick={() => handleAction(app.uuid!, 'restart')}
							disabled={!!loadingAction()}
							class="px-2 py-1 text-xs bg-yellow-600 hover:bg-yellow-500 rounded disabled:opacity-50 transition-colors"
							title="Restart"
						>
							{loadingAction() === 'restart' ? '...' : '\u21BB'}
						</button>
						<button
							type="button"
							onClick={() => handleAction(app.uuid!, 'stop')}
							disabled={!!loadingAction()}
							class="px-2 py-1 text-xs bg-red-600 hover:bg-red-500 rounded disabled:opacity-50 transition-colors"
							title="Stop"
						>
							{loadingAction() === 'stop' ? '...' : '\u25A0'}
						</button>
					</Show>
					<Show when={!appRunning}>
						<button
							type="button"
							onClick={() => handleAction(app.uuid!, 'start')}
							disabled={!!loadingAction()}
							class="px-2 py-1 text-xs bg-green-600 hover:bg-green-500 rounded disabled:opacity-50 transition-colors"
							title="Start"
						>
							{loadingAction() === 'start' ? '...' : '\u25B6'}
						</button>
					</Show>
				</div>
			</div>
		);
	};

	const ServiceItem = (service: CoolifyWidgetService) => (
		<div class="flex items-center justify-between p-2 bg-black/20 rounded mb-1">
			<div class="flex items-center gap-2 min-w-0 flex-1">
				<div class="w-2 h-2 rounded-full flex-shrink-0 bg-blue-500" />
				<div class="min-w-0 flex-1">
					<div class="font-medium text-sm truncate">
						{service.name || 'Unnamed'}
					</div>
				</div>
			</div>
		</div>
	);

	const LoadingSkeleton = () => (
		<div class="space-y-1">
			<For each={[1, 2, 3]}>
				{() => (
					<div class="flex items-center justify-between p-2 bg-black/20 rounded animate-pulse">
						<div class="flex items-center gap-2">
							<div class="w-2 h-2 rounded-full bg-gray-600" />
							<div class="w-24 h-4 bg-gray-600 rounded" />
						</div>
						<div class="w-12 h-3 bg-gray-600 rounded" />
					</div>
				)}
			</For>
		</div>
	);

	return (
		<div class="w-full p-4 bg-glass-bg rounded-lg glass-text">
			{/* Header */}
			<div class="flex items-center justify-between mb-4">
				<h3 class="text-lg font-bold">Coolify Dashboard</h3>
				<span class="text-xs text-gray-400">Auto-refresh: 30s</span>
			</div>

			{/* Applications Section */}
			<div class="mb-4">
				<button
					type="button"
					class="flex items-center gap-2 w-full text-left mb-2 hover:text-blue-400 transition-colors"
					onClick={() => setAppsExpanded(!appsExpanded())}
				>
					<span class="text-sm">{appsExpanded() ? '\u25BC' : '\u25B6'}</span>
					<span class="font-semibold">
						Applications ({applications().length})
					</span>
				</button>
				<Show when={appsExpanded()}>
					<Show when={!appsLoading()} fallback={<LoadingSkeleton />}>
						<Show
							when={!appsError()}
							fallback={
								<div class="text-sm text-red-300 p-2 bg-red-900/20 rounded">
									{appsError()}
								</div>
							}
						>
							<Show
								when={applications().length > 0}
								fallback={
									<div class="text-sm text-gray-400 p-2">
										No applications configured
									</div>
								}
							>
								<For each={applications()}>{(app) => ApplicationItem(app)}</For>
							</Show>
						</Show>
					</Show>
				</Show>
			</div>

			{/* Services Section */}
			<div>
				<button
					type="button"
					class="flex items-center gap-2 w-full text-left mb-2 hover:text-blue-400 transition-colors"
					onClick={() => setServicesExpanded(!servicesExpanded())}
				>
					<span class="text-sm">
						{servicesExpanded() ? '\u25BC' : '\u25B6'}
					</span>
					<span class="font-semibold">Services ({services().length})</span>
				</button>
				<Show when={servicesExpanded()}>
					<Show when={!servicesLoading()} fallback={<LoadingSkeleton />}>
						<Show
							when={!servicesError()}
							fallback={
								<div class="text-sm text-red-300 p-2 bg-red-900/20 rounded">
									{servicesError()}
								</div>
							}
						>
							<Show
								when={services().length > 0}
								fallback={
									<div class="text-sm text-gray-400 p-2">
										No services configured
									</div>
								}
							>
								<For each={services()}>{(service) => ServiceItem(service)}</For>
							</Show>
						</Show>
					</Show>
				</Show>
			</div>
		</div>
	);
}
