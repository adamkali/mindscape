import {
	createContext,
	createResource,
	createSignal,
	type ParentComponent,
	type Resource,
	useContext,
} from 'solid-js';
import type {
	RepositoryGetInvitesByUserIDRow,
	RepositoryGetMembersByHouseholdIDRow,
	RepositoryHousehold,
	RepositorySearchUsersByUsernameRow,
} from '@/api';
import { HouseholdsApi, UsersApi } from '@/api';
import { getAuthenticatedApiConfig } from '@/utils/apiConfig';

export interface HouseholdContextValue {
	households: Resource<RepositoryHousehold[]>;
	invites: Resource<RepositoryGetInvitesByUserIDRow[]>;
	refetchHouseholds: () => void;
	refetchInvites: () => void;
	createHousehold: (name: string) => Promise<RepositoryHousehold | null>;
	deleteHousehold: (householdId: string) => Promise<void>;
	leaveHousehold: (householdId: string) => Promise<void>;
	getMembers: (
		householdId: string,
	) => Promise<RepositoryGetMembersByHouseholdIDRow[]>;
	removeMember: (householdId: string, userId: string) => Promise<void>;
	searchUsers: (q: string) => Promise<RepositorySearchUsersByUsernameRow[]>;
	createInvite: (householdId: string, invitedUserId: string) => Promise<string>;
	acceptInvite: (code: string) => Promise<void>;
	rejectInvite: (code: string) => Promise<void>;
}

const HouseholdContext = createContext<HouseholdContextValue>();

export const HouseholdProvider: ParentComponent = (props) => {
	const api = new HouseholdsApi(getAuthenticatedApiConfig());
	const usersApi = new UsersApi(getAuthenticatedApiConfig());

	const [householdsKey, setHouseholdsKey] = createSignal(0);
	const [invitesKey, setInvitesKey] = createSignal(0);

	const [households] = createResource(householdsKey, async () => {
		const response = await api.getMyHouseholds();
		return response.data ?? [];
	});

	const [invites] = createResource(invitesKey, async () => {
		const response = await api.getMyInvites();
		return response.data ?? [];
	});

	const refetchHouseholds = () => setHouseholdsKey((k) => k + 1);
	const refetchInvites = () => setInvitesKey((k) => k + 1);

	const createHousehold = async (name: string) => {
		const response = await api.createHousehold({
			createHouseholdRequest: { name },
		});
		refetchHouseholds();
		return response.data ?? null;
	};

	const deleteHousehold = async (householdId: string) => {
		await api.deleteHousehold({ householdId });
		refetchHouseholds();
	};

	const leaveHousehold = async (householdId: string) => {
		await api.leaveHousehold({ householdId });
		refetchHouseholds();
	};

	const getMembers = async (householdId: string) => {
		const response = await api.getHouseholdMembers({ householdId });
		return response.data ?? [];
	};

	const removeMember = async (householdId: string, userId: string) => {
		await api.removeMember({ householdId, userId });
	};

	const searchUsers = async (q: string) => {
		if (!q.trim()) return [];
		const response = await usersApi.searchUsers({ q });
		return response.data ?? [];
	};

	const createInvite = async (householdId: string, invitedUserId: string) => {
		const response = await api.createInvite({
			householdId,
			createInviteRequest: { invitedUserId },
		});
		return response.data ?? '';
	};

	const acceptInvite = async (code: string) => {
		await api.acceptInvite({ code });
		refetchInvites();
		refetchHouseholds();
	};

	const rejectInvite = async (code: string) => {
		await api.rejectInvite({ code });
		refetchInvites();
	};

	const value: HouseholdContextValue = {
		households,
		invites,
		refetchHouseholds,
		refetchInvites,
		createHousehold,
		deleteHousehold,
		leaveHousehold,
		getMembers,
		removeMember,
		searchUsers,
		createInvite,
		acceptInvite,
		rejectInvite,
	};

	return (
		<HouseholdContext.Provider value={value}>
			{props.children}
		</HouseholdContext.Provider>
	);
};

export const useHousehold = () => {
	const context = useContext(HouseholdContext);
	if (!context) {
		throw new Error('useHousehold must be used within a HouseholdProvider');
	}
	return context;
};
