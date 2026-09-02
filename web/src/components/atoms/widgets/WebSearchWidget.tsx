import { createSignal } from 'solid-js';

interface WebSearchWidgetProps {
	url: string;
	engine?: string;
}

export default function WebSearchWidget(props: WebSearchWidgetProps) {
	const { url } = props;
	const [searchQuery, setSearchQuery] = createSignal('');

	const sanitizeUrl = (urlTemplate: string, query: string) => {
		const encodedQuery = encodeURIComponent(query);
		return urlTemplate.replace('%s', encodedQuery);
	};

	const handleSubmit = (e: Event) => {
		e.preventDefault();
		const query = searchQuery().trim();
		if (query) {
			const searchUrl = sanitizeUrl(url, query);
			window.open(searchUrl, '_blank');
			setSearchQuery('');
		}
	};

	return (
		<form
			class="h-full w-full flex items-center justify-center"
			onSubmit={handleSubmit}
		>
			<label class="sr-only" for="web-search">
				Search
			</label>
			<input
				id="web-search"
				type="text"
				placeholder={`🔍 Search with ${props.engine ?? 'Google'}... `}
				class="h-full w-full px-4 py-2 rounded-lg bg-glass-bg text-foreground focus:outline-none"
				value={searchQuery()}
				onInput={(e) => setSearchQuery(e.currentTarget.value)}
			/>
		</form>
	);
}
