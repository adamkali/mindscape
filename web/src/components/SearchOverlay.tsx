import Fuse from 'fuse.js';
import { createEffect, createMemo, createSignal, For, Show } from 'solid-js';
import {
	type SearchMode,
	type TreeNodeEntry,
	useKeyboardNav,
} from '@/contexts/KeyboardNavContext';
import { useTree } from '@/contexts/TreeContext';

const TYPE_LABEL: Record<TreeNodeEntry['type'], string> = {
	folder: 'folder',
	bookmark: 'link',
};

export default function SearchOverlay() {
	const nav = useKeyboardNav();
	const tree = useTree();
	const [query, setQuery] = createSignal('');
	const [selectedIdx, setSelectedIdx] = createSignal(0);
	let inputRef: HTMLInputElement | undefined;
	let fuseInstance: Fuse<TreeNodeEntry> | null = null;

	// Build the Fuse index once when tree-search opens; discard when closed.
	createEffect(() => {
		const sm = nav.searchMode();
		if (sm === 'tree') {
			fuseInstance = new Fuse(nav.getNodes(), {
				keys: ['name'],
				threshold: 0.4,
				ignoreLocation: true,
			});
		} else {
			fuseInstance = null;
		}
	});

	// Autofocus and reset state when overlay opens.
	createEffect(() => {
		if (nav.searchMode() !== null) {
			setQuery('');
			setSelectedIdx(0);
			setTimeout(() => inputRef?.focus(), 0);
		}
	});

	const results = createMemo((): TreeNodeEntry[] => {
		if (nav.searchMode() !== 'tree') return [];
		const q = query();
		if (!q) return nav.getNodes().slice(0, 10);
		if (!fuseInstance) return [];
		return fuseInstance
			.search(q)
			.map((r) => r.item)
			.slice(0, 10);
	});

	// Reset selection to first result whenever results list changes.
	createEffect(() => {
		results();
		setSelectedIdx(0);
	});

	const close = () => {
		nav.closeSearch();
		setQuery('');
	};

	const jumpToSelected = () => {
		const items = results();
		const selected = items[selectedIdx()];
		if (selected) {
			tree.setFocus(selected.id);
			selected.ref.scrollIntoView({ block: 'nearest' });
		}
		close();
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Escape') {
			e.preventDefault();
			close();
		} else if (e.key === 'Enter') {
			e.preventDefault();
			const sm = nav.searchMode();
			if (sm === 'web') {
				const q = query().trim();
				if (q) {
					window.open(
						`https://duckduckgo.com/?q=${encodeURIComponent(q)}`,
						'_blank',
					);
				}
				close();
			} else if (sm === 'tree') {
				jumpToSelected();
			}
		} else if (e.key === 'ArrowDown') {
			e.preventDefault();
			setSelectedIdx((i) => Math.min(i + 1, results().length - 1));
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			setSelectedIdx((i) => Math.max(i - 1, 0));
		}
	};

	const placeholder = (sm: SearchMode | null) =>
		sm === 'web' ? 'Search the web…' : 'Search folders and bookmarks…';

	return (
		<Show when={nav.searchMode() !== null}>
			{/*
			 * role="dialog" + aria-modal on the outer layer: this is the accessible
			 * overlay container. An invisible backdrop button behind the panel handles
			 * click-to-close without making a non-interactive div interactive.
			 */}
			<div
				role="dialog"
				aria-modal="true"
				aria-label="Search"
				class="fixed inset-0 z-[60] flex items-start justify-center pt-20"
			>
				{/* Invisible backdrop — catches outside-clicks to dismiss */}
				<button
					type="button"
					class="absolute inset-0 bg-black/50 backdrop-blur-sm cursor-default"
					onClick={close}
					aria-label="Close search"
					tabIndex={-1}
				/>

				{/* Panel — sits above the backdrop button */}
				<div class="relative z-10 w-full max-w-xl mx-4 bg-glass-bg-strong backdrop-blur-md border border-white/20 rounded-xl shadow-2xl shadow-black/40 dark:border-slate-700/50">
					{/* Input row */}
					<div class="flex items-center gap-2 p-3">
						<span class="text-foreground/40 text-sm font-mono select-none">
							{nav.searchMode() === 'web' ? '/' : '?'}
						</span>
						<input
							ref={inputRef}
							type="text"
							value={query()}
							onInput={(e) => setQuery(e.currentTarget.value)}
							onKeyDown={handleKeyDown}
							placeholder={placeholder(nav.searchMode())}
							class="flex-1 bg-transparent text-foreground text-sm outline-none placeholder-foreground/40"
							aria-label={
								nav.searchMode() === 'web' ? 'Web search' : 'Tree search'
							}
						/>
					</div>

					{/* Tree search results */}
					<Show when={nav.searchMode() === 'tree'}>
						<div class="border-t border-white/10 max-h-72 overflow-y-auto">
							<Show
								when={results().length > 0}
								fallback={
									<p class="px-4 py-3 text-sm text-foreground/50">
										{query() ? 'No results' : 'Start typing to search…'}
									</p>
								}
							>
								<For each={results()}>
									{(node, idx) => (
										<button
											type="button"
											class={`w-full flex items-center gap-3 px-4 py-2 text-sm text-left transition-colors duration-100 ${
												idx() === selectedIdx()
													? 'bg-white/10 text-foreground'
													: 'text-foreground/70 hover:bg-white/5'
											}`}
											onClick={() => {
												tree.setFocus(node.id);
												node.ref.scrollIntoView({ block: 'nearest' });
												close();
											}}
										>
											<span class="w-14 shrink-0 text-xs text-foreground/40 font-mono">
												{TYPE_LABEL[node.type]}
											</span>
											<span class="truncate">{node.name}</span>
										</button>
									)}
								</For>
							</Show>
						</div>
						<div class="px-3 py-1.5 border-t border-white/10 flex gap-3 text-xs text-foreground/40">
							<span>↑↓ navigate</span>
							<span>↵ jump to result</span>
							<span>Esc close</span>
						</div>
					</Show>

					{/* Web search hint */}
					<Show when={nav.searchMode() === 'web'}>
						<div class="px-3 py-1.5 border-t border-white/10 flex gap-3 text-xs text-foreground/40">
							<span>↵ search with DuckDuckGo</span>
							<span>Esc close</span>
						</div>
					</Show>
				</div>
			</div>
		</Show>
	);
}
