import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal, For, Show } from 'solid-js';
import { ApiKeysApi, type ServicesApiKeyDTO } from '@/api';
import {
	Button,
	Card,
	CardBody,
	CardFooter,
	CardHeader,
	DateTimeInput,
	Input,
} from '@/components/atoms';
import { Header } from '@/components/Header';
import { useAuth } from '@/contexts/AuthContext';
import { ViewProvider } from '@/contexts/ViewContext';
import { useBackgroundStyle } from '@/hooks/useBackground';

const ApiKeys = () => {
	return (
		<ViewProvider>
			<ApiKeysInner />
		</ViewProvider>
	);
};

const ApiKeysInner = () => {
	const auth = useAuth();
	const navigate = useNavigate();
	const api = new ApiKeysApi();
	const backgroundStyle = useBackgroundStyle();

	const user = auth.user();

	// Form state
	const [name, setName] = createSignal('');
	const [notBefore, setNotBefore] = createSignal('');
	const [expiration, setExpiration] = createSignal('');
	const [readAccess, setReadAccess] = createSignal(true);
	const [writeAccess, setWriteAccess] = createSignal(false);

	// List state
	const [apiKeys, setApiKeys] = createSignal<ServicesApiKeyDTO[]>([]);

	// UI state
	const [isLoading, setIsLoading] = createSignal(false);
	const [error, setError] = createSignal('');
	const [success, setSuccess] = createSignal('');
	const [newlyCreatedKey, setNewlyCreatedKey] = createSignal<string | null>(
		null,
	);
	const [copied, setCopied] = createSignal(false);
	const [deletingKeyId, setDeletingKeyId] = createSignal<string | null>(null);

	createEffect(() => {
		if (auth.isAuthenticated() && auth.token()) {
			fetchKeys();
		}
	});

	const fetchKeys = async () => {
		if (!auth.token()) return;
		try {
			const response = await api.listApiKeys({
				authorization: `Bearer ${auth.token()}`,
			});
			if (response.success && response.data) {
				setApiKeys(response.data);
			}
		} catch (err: unknown) {
			console.error('Failed to fetch API keys:', err);
		}
	};

	const handleCreate = async (e: Event) => {
		e.preventDefault();
		if (!auth.token()) return;

		setError('');
		setSuccess('');
		setIsLoading(true);

		try {
			const toRFC3339 = (v: string): string | undefined => {
				if (!v) return undefined;
				return new Date(v).toISOString();
			};

			const response = await api.createApiKey({
				authorization: `Bearer ${auth.token()}`,
				createApiKeyRequest: {
					name: name(),
					notBefore: toRFC3339(notBefore()),
					expiration: toRFC3339(expiration()),
					readAccess: readAccess(),
					writeAccess: writeAccess(),
				},
			});

			if (response.success && response.data) {
				setNewlyCreatedKey(response.data.rawKey ?? null);
				setSuccess('API key created successfully!');
				setName('');
				setNotBefore('');
				setExpiration('');
				setReadAccess(true);
				setWriteAccess(false);
				await fetchKeys();
			} else {
				setError(response.message ?? 'Failed to create API key');
			}
		} catch (err: unknown) {
			const message =
				err instanceof Error ? err.message : 'Failed to create API key';
			setError(message);
		} finally {
			setIsLoading(false);
		}
	};

	const handleDelete = async (keyId: string) => {
		if (!auth.token()) return;

		setError('');
		setSuccess('');

		try {
			await api.deleteApiKey({
				authorization: `Bearer ${auth.token()}`,
				keyId,
			});
			setSuccess('API key deleted.');
			setDeletingKeyId(null);
			await fetchKeys();
		} catch (err: unknown) {
			const message =
				err instanceof Error ? err.message : 'Failed to delete API key';
			setError(message);
		}
	};

	const handleCopy = async () => {
		const key = newlyCreatedKey();
		if (!key) return;
		await navigator.clipboard.writeText(key);
		setCopied(true);
		setTimeout(() => setCopied(false), 2000);
	};

	const getKeyStatus = (
		key: ServicesApiKeyDTO,
	): { label: string; color: string } => {
		const now = new Date();
		if (key.expiration) {
			const exp = new Date(key.expiration);
			if (exp < now)
				return {
					label: 'Expired',
					color: 'bg-red-500/30 text-red-300 border-red-400/40',
				};
		}
		if (key.notBefore) {
			const nb = new Date(key.notBefore);
			if (nb > now)
				return {
					label: 'Not Yet Active',
					color: 'bg-yellow-500/30 text-yellow-300 border-yellow-400/40',
				};
		}
		return {
			label: 'Active',
			color: 'bg-green-500/30 text-green-300 border-green-400/40',
		};
	};

	if (!auth.isAuthenticated() || !user) {
		navigate('/login');
		return null;
	}

	return (
		<div
			class="h-screen overflow-hidden bg-background"
			style={backgroundStyle()}
		>
			<Header />
			<div class="max-w-2xl mx-auto mt-4 space-y-4 overflow-y-auto max-h-[calc(100vh-4rem)] pb-8">
				{/* Create API Key */}
				<Card variant="glass">
					<CardHeader
						title="Create API Key"
						subtitle="Generate a new key for programmatic access"
					/>
					<form onSubmit={handleCreate}>
						<CardBody padding="lg">
							<div class="space-y-4">
								<Input
									id="keyName"
									name="keyName"
									type="text"
									required
									placeholder="Key name (e.g. My Integration)"
									value={name()}
									onInput={(e) =>
										setName((e.currentTarget as HTMLInputElement).value)
									}
									label="Key Name"
								/>
								<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
									<DateTimeInput
										id="notBefore"
										name="notBefore"
										value={notBefore()}
										onInput={(e) =>
											setNotBefore((e.currentTarget as HTMLInputElement).value)
										}
										label="Not Before"
									/>
									<DateTimeInput
										id="expiration"
										name="expiration"
										value={expiration()}
										onInput={(e) =>
											setExpiration((e.currentTarget as HTMLInputElement).value)
										}
										label="Expiration"
									/>
								</div>
								<div class="flex items-center gap-3">
									<button
										type="button"
										onClick={() => setReadAccess(!readAccess())}
										class={`px-3 py-1.5 text-sm rounded-lg border transition-all duration-200 ${
											readAccess()
												? 'bg-green-500/30 border-green-400/40 text-foreground'
												: 'bg-glass-bg border-glass-border text-foreground/50'
										}`}
									>
										Read Access
									</button>
									<button
										type="button"
										onClick={() => setWriteAccess(!writeAccess())}
										class={`px-3 py-1.5 text-sm rounded-lg border transition-all duration-200 ${
											writeAccess()
												? 'bg-green-500/30 border-green-400/40 text-foreground'
												: 'bg-glass-bg border-glass-border text-foreground/50'
										}`}
									>
										Write Access
									</button>
								</div>
							</div>
						</CardBody>
						<CardFooter>
							<Button
								type="submit"
								disabled={isLoading() || !name()}
								variant="secondary"
							>
								{isLoading() ? 'Creating...' : 'Create Key'}
							</Button>
							<Button type="button" onClick={() => navigate('/')}>
								Cancel
							</Button>
						</CardFooter>
					</form>
				</Card>

				{/* Newly Created Key */}
				<Show when={newlyCreatedKey()}>
					<Card variant="glass">
						<CardBody padding="lg">
							<div class="space-y-3">
								<div class="text-yellow-300 font-medium text-sm">
									Copy this key now. You won't see it again.
								</div>
								<div class="font-mono text-sm bg-black/30 rounded-lg p-3 break-all border border-yellow-400/30 text-foreground">
									{newlyCreatedKey()}
								</div>
								<div class="flex gap-2">
									<Button
										type="button"
										variant="secondary"
										onClick={handleCopy}
									>
										{copied() ? 'Copied!' : 'Copy to Clipboard'}
									</Button>
									<Button
										type="button"
										onClick={() => setNewlyCreatedKey(null)}
									>
										Done
									</Button>
								</div>
							</div>
						</CardBody>
					</Card>
				</Show>

				{/* Key List */}
				<Card variant="glass">
					<CardHeader
						title="Your API Keys"
						subtitle="Manage your existing keys"
					/>
					<CardBody padding="lg">
						<Show
							when={apiKeys().length > 0}
							fallback={
								<div class="text-foreground/50 text-sm">No API keys yet</div>
							}
						>
							<div class="space-y-3">
								<For each={apiKeys()}>
									{(key) => {
										const status = getKeyStatus(key);
										return (
											<div class="flex items-center justify-between p-3 rounded-lg bg-glass-bg border border-glass-border">
												<div class="flex-1 min-w-0 space-y-1">
													<div class="flex items-center gap-2 flex-wrap">
														<span class="font-medium text-sm text-foreground">
															{key.name}
														</span>
														<span
															class={`text-xs px-2 py-0.5 rounded-full border ${status.color}`}
														>
															{status.label}
														</span>
														<Show when={key.readAccess}>
															<span class="text-xs px-2 py-0.5 rounded-full bg-green-500/20 text-green-300 border border-green-400/30">
																Read
															</span>
														</Show>
														<Show when={key.writeAccess}>
															<span class="text-xs px-2 py-0.5 rounded-full bg-blue-500/20 text-blue-300 border border-blue-400/30">
																Write
															</span>
														</Show>
													</div>
													<div class="text-xs text-foreground/40 font-mono truncate">
														{key.keyId}
													</div>
													<Show when={key.createdAt}>
														<div class="text-xs text-foreground/40">
															Created{' '}
															{new Date(
																key.createdAt as string,
															).toLocaleDateString()}
														</div>
													</Show>
												</div>
												<div class="ml-3 flex-shrink-0">
													<Show
														when={deletingKeyId() === key.keyId}
														fallback={
															<Button
																type="button"
																variant="danger"
																onClick={() =>
																	setDeletingKeyId(key.keyId ?? null)
																}
															>
																Delete
															</Button>
														}
													>
														<div class="flex gap-2">
															<Button
																type="button"
																variant="danger"
																onClick={() =>
																	handleDelete(key.keyId as string)
																}
															>
																Confirm
															</Button>
															<Button
																type="button"
																onClick={() => setDeletingKeyId(null)}
															>
																Cancel
															</Button>
														</div>
													</Show>
												</div>
											</div>
										);
									}}
								</For>
							</div>
						</Show>
					</CardBody>
				</Card>

				{/* Status Messages */}
				<Show when={error() || success()}>
					<Card variant="glass">
						<CardBody padding="md">
							<Show when={error()}>
								<div class="text-red-400 text-sm">{error()}</div>
							</Show>
							<Show when={success() && !error()}>
								<div class="text-green-400 text-sm">{success()}</div>
							</Show>
						</CardBody>
					</Card>
				</Show>
			</div>
		</div>
	);
};

export default ApiKeys;
