import { Show } from 'solid-js';
import { type NavMode, useKeyboardNav } from '@/contexts/KeyboardNavContext';
import { useView } from '@/contexts/ViewContext';

const LABELS: Record<Exclude<NavMode, 'mouse'>, string> = {
	nav: 'NAV',
	edit: 'EDIT',
	create: 'CREATE',
	delete: 'DELETE',
	move: 'MOVE',
};

export default function ModePill() {
	const nav = useKeyboardNav();
	const view = useView();

	// One combined indicator rather than two competing pills: nav mode first,
	// then the layout it is acting on.
	return (
		<Show when={nav.mode() !== 'mouse'}>
			<div
				class="fixed bottom-4 right-4 z-50 px-3 py-1 rounded-full text-xs font-mono font-semibold
					bg-background/80 backdrop-blur-sm border border-white/20
					text-foreground shadow-lg shadow-black/30 select-none pointer-events-none"
			>
				{LABELS[nav.mode() as Exclude<NavMode, 'mouse'>]}
				<span class="text-foreground/50">
					{' '}
					· {view.homepageMode().toUpperCase()}
				</span>
			</div>
		</Show>
	);
}
