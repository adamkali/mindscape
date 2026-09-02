import { useNavigate, useSearchParams } from '@solidjs/router';
import {
	type Component,
	createEffect,
	createResource,
	createSignal,
	For,
	Show,
} from 'solid-js';
import type {
	RepositoryGetMembersByHouseholdIDRow,
	RepositoryHousehold,
	RepositorySearchUsersByUsernameRow,
} from '@/api';
import { Button, Card, CardBody, CardHeader, Input } from '@/components/atoms';
import { Header } from '@/components/Header';
import { QRCode } from '@/components/QRCode';
import { useAuth } from '@/contexts/AuthContext';
import { HouseholdProvider, useHousehold } from '@/contexts/HouseholdContext';
import { ViewProvider } from '@/contexts/ViewContext';
import { useBackgroundStyle } from '@/hooks/useBackground';

const Households = () => {
	return (
		<ViewProvider>
			<HouseholdProvider>
				<HouseholdsInner />
			</HouseholdProvider>
		</ViewProvider>
	);
};

const HouseholdsInner = () => {
	const auth = useAuth();
	const navigate = useNavigate();
	const household = useHousehold();
	const backgroundStyle = useBackgroundStyle();
	const [searchParams, setSearchParams] = useSearchParams();

	const user = auth.user();

	const [name, setName] = createSignal('');
	const [creating, setCreating] = createSignal(false);
	const [error, setError] = createSignal('');
	const [success, setSuccess] = createSignal('');

	// Auto-accept an invite arriving via a scanned QR deep link (?invite=<code>).
	createEffect(() => {
		const code = searchParams.invite;
		if (typeof code === 'string' && code) {
			setSearchParams({ invite: undefined });
			household
				.acceptInvite(code)
				.then(() => setSuccess('Invite accepted — welcome to the household!'))
				.catch(() => setError('That invite is no longer valid.'));
		}
	});

	const handleCreate = async (e: Event) => {
		e.preventDefault();
		if (!name().trim()) return;
		setError('');
		setSuccess('');
		setCreating(true);
		try {
			await household.createHousehold(name().trim());
			setSuccess('Household created.');
			setName('');
		} catch (err: unknown) {
			setError(
				err instanceof Error ? err.message : 'Failed to create household',
			);
		} finally {
			setCreating(false);
		}
	};

	const handleAccept = async (code: string) => {
		setError('');
		setSuccess('');
		try {
			await household.acceptInvite(code);
			setSuccess('Invite accepted.');
		} catch {
			setError('Failed to accept invite.');
		}
	};

	const handleReject = async (code: string) => {
		setError('');
		setSuccess('');
		try {
			await household.rejectInvite(code);
		} catch {
			setError('Failed to reject invite.');
		}
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
				{/* Pending Invites */}
				<Show when={(household.invites() ?? []).length > 0}>
					<Card variant="glass">
						<CardHeader
							title="Pending Invites"
							subtitle="Households that have invited you"
						/>
						<CardBody padding="lg">
							<div class="space-y-3">
								<For each={household.invites()}>
									{(invite) => (
										<div class="flex items-center justify-between p-3 rounded-lg bg-glass-bg border border-glass-border">
											<div class="min-w-0">
												<div class="font-medium text-sm text-foreground">
													{invite.householdName}
												</div>
												<Show when={invite.expiresAt}>
													<div class="text-xs text-foreground/40">
														Expires{' '}
														{new Date(
															invite.expiresAt as string,
														).toLocaleDateString()}
													</div>
												</Show>
											</div>
											<div class="flex gap-2 flex-shrink-0">
												<Button
													type="button"
													variant="secondary"
													onClick={() => handleAccept(invite.code as string)}
												>
													Accept
												</Button>
												<Button
													type="button"
													variant="danger"
													onClick={() => handleReject(invite.code as string)}
												>
													Reject
												</Button>
											</div>
										</div>
									)}
								</For>
							</div>
						</CardBody>
					</Card>
				</Show>

				{/* Create Household */}
				<Card variant="glass">
					<CardHeader
						title="Create Household"
						subtitle="Start a household and invite others to join"
					/>
					<form onSubmit={handleCreate}>
						<CardBody padding="lg">
							<div class="flex gap-2 items-end">
								<div class="flex-1">
									<Input
										id="householdName"
										name="householdName"
										type="text"
										required
										placeholder="Household name"
										value={name()}
										onInput={(e) =>
											setName((e.currentTarget as HTMLInputElement).value)
										}
										label="Name"
									/>
								</div>
								<Button
									type="submit"
									variant="secondary"
									disabled={creating() || !name().trim()}
								>
									{creating() ? 'Creating…' : 'Create'}
								</Button>
							</div>
						</CardBody>
					</form>
				</Card>

				{/* My Households */}
				<Card variant="glass">
					<CardHeader
						title="My Households"
						subtitle="Households you belong to"
					/>
					<CardBody padding="lg">
						<Show
							when={(household.households() ?? []).length > 0}
							fallback={
								<div class="text-foreground/50 text-sm">
									You're not part of any household yet.
								</div>
							}
						>
							<div class="space-y-3">
								<For each={household.households()}>
									{(h) => (
										<HouseholdRow
											household={h}
											isOwner={h.ownerId === user.id}
											onError={setError}
											onSuccess={setSuccess}
										/>
									)}
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

interface HouseholdRowProps {
	household: RepositoryHousehold;
	isOwner: boolean;
	onError: (message: string) => void;
	onSuccess: (message: string) => void;
}

const HouseholdRow: Component<HouseholdRowProps> = (props) => {
	const ctx = useHousehold();
	const [expanded, setExpanded] = createSignal(false);
	const [inviteOpen, setInviteOpen] = createSignal(false);

	const [members, { refetch: refetchMembers }] = createResource(
		expanded,
		async (isExpanded) => {
			if (!isExpanded) return [] as RepositoryGetMembersByHouseholdIDRow[];
			return ctx.getMembers(props.household.id as string);
		},
	);

	const handleLeave = async () => {
		try {
			await ctx.leaveHousehold(props.household.id as string);
			props.onSuccess('You left the household.');
		} catch {
			props.onError('Failed to leave household.');
		}
	};

	const handleDelete = async () => {
		try {
			await ctx.deleteHousehold(props.household.id as string);
			props.onSuccess('Household deleted.');
		} catch {
			props.onError('Failed to delete household.');
		}
	};

	const handleRemoveMember = async (userId: string) => {
		try {
			await ctx.removeMember(props.household.id as string, userId);
			refetchMembers();
		} catch {
			props.onError('Failed to remove member.');
		}
	};

	return (
		<div class="p-3 rounded-lg bg-glass-bg border border-glass-border space-y-3">
			<div class="flex items-center justify-between">
				<div class="flex items-center gap-2">
					<span class="font-medium text-sm text-foreground">
						{props.household.name}
					</span>
					<Show when={props.isOwner}>
						<span class="text-xs px-2 py-0.5 rounded-full bg-yellow-500/20 text-yellow-300 border border-yellow-400/30">
							Owner
						</span>
					</Show>
				</div>
				<div class="flex gap-2 flex-shrink-0">
					<Button type="button" onClick={() => setExpanded(!expanded())}>
						{expanded() ? 'Hide' : 'Members'}
					</Button>
					<Show when={props.isOwner}>
						<Button
							type="button"
							variant="secondary"
							onClick={() => setInviteOpen(!inviteOpen())}
						>
							Invite
						</Button>
						<Button type="button" variant="danger" onClick={handleDelete}>
							Delete
						</Button>
					</Show>
					<Show when={!props.isOwner}>
						<Button type="button" variant="danger" onClick={handleLeave}>
							Leave
						</Button>
					</Show>
				</div>
			</div>

			<Show when={inviteOpen() && props.isOwner}>
				<InvitePanel
					householdId={props.household.id as string}
					householdName={props.household.name as string}
					onError={props.onError}
				/>
			</Show>

			<Show when={expanded()}>
				<div class="space-y-2 pt-2 border-t border-glass-border">
					<Show
						when={(members() ?? []).length > 0}
						fallback={
							<div class="text-xs text-foreground/50">Loading members…</div>
						}
					>
						<For each={members()}>
							{(member) => (
								<div class="flex items-center justify-between text-sm">
									<div class="flex items-center gap-2">
										<span class="text-foreground">{member.username}</span>
										<span class="text-xs px-2 py-0.5 rounded-full bg-glass-bg border border-glass-border text-foreground/60">
											{member.role}
										</span>
									</div>
									<Show
										when={
											props.isOwner && member.userId !== props.household.ownerId
										}
									>
										<Button
											type="button"
											variant="danger"
											onClick={() =>
												handleRemoveMember(member.userId as string)
											}
										>
											Remove
										</Button>
									</Show>
								</div>
							)}
						</For>
					</Show>
				</div>
			</Show>
		</div>
	);
};

interface InvitePanelProps {
	householdId: string;
	householdName: string;
	onError: (message: string) => void;
}

const InvitePanel: Component<InvitePanelProps> = (props) => {
	const ctx = useHousehold();
	const [query, setQuery] = createSignal('');
	const [results, setResults] = createSignal<
		RepositorySearchUsersByUsernameRow[]
	>([]);
	const [searching, setSearching] = createSignal(false);
	const [inviteCode, setInviteCode] = createSignal('');
	const [invitedName, setInvitedName] = createSignal('');
	const [copied, setCopied] = createSignal(false);

	const joinUrl = () =>
		`${window.location.origin}/households?invite=${inviteCode()}`;

	const handleSearch = async () => {
		if (!query().trim()) return;
		setSearching(true);
		try {
			setResults(await ctx.searchUsers(query()));
		} catch {
			props.onError('Failed to search users.');
		} finally {
			setSearching(false);
		}
	};

	const handleInvite = async (invitee: RepositorySearchUsersByUsernameRow) => {
		try {
			const code = await ctx.createInvite(
				props.householdId,
				invitee.id as string,
			);
			setInviteCode(code);
			setInvitedName(invitee.username ?? '');
			setResults([]);
			setQuery('');
		} catch {
			props.onError('Failed to create invite.');
		}
	};

	const handleCopy = async () => {
		await navigator.clipboard.writeText(joinUrl());
		setCopied(true);
		setTimeout(() => setCopied(false), 2000);
	};

	return (
		<div class="space-y-3 pt-2 border-t border-glass-border">
			<Show when={!inviteCode()}>
				<div class="flex gap-2 items-end">
					<div class="flex-1">
						<Input
							id={`invite-search-${props.householdId}`}
							name="inviteSearch"
							type="text"
							placeholder="Search users by username"
							value={query()}
							onInput={(e) =>
								setQuery((e.currentTarget as HTMLInputElement).value)
							}
							label="Invite a user"
						/>
					</div>
					<Button
						type="button"
						onClick={handleSearch}
						disabled={searching() || !query().trim()}
					>
						{searching() ? 'Searching…' : 'Search'}
					</Button>
				</div>
				<Show when={results().length > 0}>
					<div class="space-y-2">
						<For each={results()}>
							{(u) => (
								<div class="flex items-center justify-between text-sm">
									<span class="text-foreground">{u.username}</span>
									<Button
										type="button"
										variant="secondary"
										onClick={() => handleInvite(u)}
									>
										Invite
									</Button>
								</div>
							)}
						</For>
					</div>
				</Show>
			</Show>

			<Show when={inviteCode()}>
				<div class="flex flex-col items-center text-center gap-3">
					<div class="text-sm text-foreground">
						Invite sent to <span class="font-medium">{invitedName()}</span>.
						They can scan this code to join{' '}
						<span class="font-medium">{props.householdName}</span>.
					</div>
					<QRCode
						value={joinUrl()}
						size={180}
						label="Household invite QR code"
					/>
					<div class="font-mono text-xs bg-black/30 rounded-lg p-2 break-all border border-glass-border text-foreground w-full">
						{inviteCode()}
					</div>
					<div class="flex gap-2">
						<Button type="button" variant="secondary" onClick={handleCopy}>
							{copied() ? 'Copied!' : 'Copy Link'}
						</Button>
						<Button type="button" onClick={() => setInviteCode('')}>
							Invite Another
						</Button>
					</div>
				</div>
			</Show>
		</div>
	);
};

export default Households;
