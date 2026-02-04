import { createSignal } from 'solid-js';

interface SearchbarWidgetProps {
	url: string;
	engine?: string;
}

export default function SearchbarWidget(props: SearchbarWidgetProps) {
	const { url } = props;
	const [searchQuery, setSearchQuery] = createSignal('');

	const sanitizeUrl = (urlTemplate: string, query: string) => {
		// Encode the search query to make it URL-safe
		const encodedQuery = encodeURIComponent(query);
		return urlTemplate.replace('%s', encodedQuery);
	};

	const handleSubmit = (e: Event) => {
		e.preventDefault(); // Prevent default form submission
		const query = searchQuery().trim();
		if (query) {
			const searchUrl = sanitizeUrl(url, query);
			window.open(searchUrl, '_blank');
			setSearchQuery(''); // Clear search after opening
		}
	};

	return (
		<form
			class="h-full w-full flex items-center justify-center"
			onSubmit={handleSubmit}
		>
			<label class="sr-only" for="search">
				Search
			</label>
			<input
				id="search"
				type="text"
				placeholder={`🔍 Search with ${props.engine}... `}
				class="h-full w-full px-4 py-2 rounded-lg bg-glass-bg text-foreground focus:outline-none"
				value={searchQuery()}
				onInput={(e) => setSearchQuery(e.currentTarget.value)}
			/>
		</form>
	);
}
