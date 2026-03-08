import type { JSX } from 'solid-js';

export function SaveIcon(props: JSX.IntrinsicElements['svg']) {
	// <!-- Icon from Material Symbols by Google - https://github.com/google/material-design-icons/blob/master/LICENSE -->
	return (
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="1em"
			height="1em"
			viewBox="0 0 24 24"
			role="img"
			aria-label="Save"
			{...props}
		>
			<path
				fill="currentColor"
				d="M9 16.17L4.83 12l-1.42 1.41L9 19L21 7l-1.41-1.41z"
			/>
		</svg>
	);
}
export default SaveIcon;
