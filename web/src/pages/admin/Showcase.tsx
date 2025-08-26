import { FoldersApi } from '@/api';
import {
	Button,
	Card,
	CardBody,
	CardFooter,
	CardHeader,
} from '@/components/atoms';
import Input from '@/components/atoms/Input';
import CreateFolderComponent from '@/components/CreateFolderComponent';
import FolderComponent from '@/components/FolderComponent';
import { Header } from '@/components/Header';
import { useAuth } from '@/contexts/AuthContext';
import { EmptyGuid } from '@/utils';
import { createSignal, Show } from 'solid-js';

export default function Showcase() {
	const [formSuccessful, setFormSuccessful] = createSignal<boolean | undefined>(
		undefined,
	);
	const folderAPIRef = new FoldersApi();
	const auth = useAuth();
	const [formMessage, setFormMessage] =
		createSignal<string>('This is a message');

	const submit = async (e: Event) => {
		e.preventDefault();
		setFormMessage('This is a successful message');
		setFormSuccessful(true);
	};

	return (
		<div class="min-h-screen bg-background space-y-4">
			<Header />
			<Card class="flex flex-col mx-auto w-1/2 ">
				<CardHeader>
					<div class="text-lg text-foreground">Mindscape UI Buttons</div>
				</CardHeader>
				<CardBody padding='lg' spac>
					<Button variant="primary">Primary</Button>
					<Button variant="secondary">Secondary</Button>
					<Button variant="tertiary">Tertiary</Button>
					<Button variant="danger">Danger</Button>
				</CardBody>
			</Card>
			<Card class="flex flex-col mx-auto w-1/2 ">
				<CardHeader>
					<div class="text-lg text-foreground">Mindscape UI Forms</div>
				</CardHeader>
				<CardBody>
					<form onSubmit={submit} class="flex flex-row text-center" id="form">
						<Input
							variant="primary"
							title="Search"
							value="search"
							onValueChange={(value) => console.log(value)}
							placeholder="Enter text..."
							error={false}
							required
							errorMessage="This field is required"
						/>
						<Input
							variant="secondary"
							title="Search"
							value="search"
							onValueChange={(value) => console.log(value)}
							placeholder="Enter text..."
							error={false}
							errorMessage="This field is required"
						/>
						<Input
							variant="tertiary"
							title="Search"
							value="search"
							onValueChange={(value) => console.log(value)}
							placeholder="Enter text..."
							error={false}
							errorMessage="This field is required"
						/>
						<Input
							variant="danger"
							title="Search"
							value="search"
							onValueChange={(value) => console.log(value)}
							placeholder="Enter text..."
							error={false}
							errorMessage="This field is required"
						/>
					</form>
				</CardBody>
				<CardFooter>
					<Button
						variant="primary"
						on:click={() => {
							const form = document.getElementById('form') as HTMLFormElement;
							// get it to submit the form data
							// we do not connect to a server in this example so just
							// log the form data to the console
							form.submit();
						}}
					>
						Submit
					</Button>
					<Show
						when={formSuccessful() === true}
						fallback={<span class="text-error">{formMessage()}</span>}
					>
						<span class="text-green-500">{formMessage()}</span>
					</Show>
				</CardFooter>
			</Card>
			<Card class="flex flex-col mx-auto w-1/2 ">
				<CardHeader>
					<div class="text-lg text-foreground">Mindscape UI Folders</div>
				</CardHeader>
				<CardBody>
					<div class="flex flex-col space-y-4">
						<FolderComponent
							folder={{
								id: '09870d37-a8e1-4d3c-b9fa-c74475653315',
								userId: '3fa3ebb5-7f5d-49d7-a36a-cf3878d49aea',
								name: 'Folder 1',
								createdDatetime: '2023-01-01T00:00:00.000Z',
								updatedDatetime: '2025-01-01T00:00:00.000Z',
								bookmarks: [],
								notes: [],
								children: [],
							}}
							selectedFolder={() => { }}
							deleteFolder={() => { }}
						/>
						<FolderComponent
							folder={{
								id: '50b15788-4553-4840-8118-3bd15250fbf9',
								userId: '7efc610e-5cb0-4dbb-95df-3507dd919202',
								name: 'Folder 2',
								createdDatetime: '2023-02-02T00:00:00.000Z',
								updatedDatetime: '2025-02-02T00:00:00.000Z',
								bookmarks: [],
								notes: [],
								children: [],
							}}
							selectedFolder={() => { }}
							deleteFolder={() => { }}
						/>
						<CreateFolderComponent
							userId={auth.user()?.id ?? EmptyGuid}
							parentId={undefined}
							auth={auth}
							setShowCreateFolder={() => true}
							folderAPIRef={folderAPIRef}
						/>
					</div>
				</CardBody>
			</Card>
		</div>
	);
}
