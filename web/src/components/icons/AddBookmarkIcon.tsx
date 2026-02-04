import type { JSX } from 'solid-js';

export function AddBookmarkIcon(props: JSX.IntrinsicElements['svg']) {
	// <!-- Icon from Material Symbols by Google - https://github.com/google/material-design-icons/blob/master/LICENSE -->
	return (
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="1em"
			height="1em"
			viewBox="0 0 24 24"
			{...props}
		>
			<path
				fill="currentColor"
				d="M19 21l-7-3l-7 3V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16zm-9-9h2v2h2v-2h2v-2h-2V8h-2v2H8v2z"
			/>
		</svg>
	);
}
export default AddBookmarkIcon;
