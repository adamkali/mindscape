import type { JSX } from 'solid-js';

export default function WidgetIcon(props: JSX.IntrinsicElements['svg']) {
	return (
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="1em"
			height="1em"
			viewBox="0 0 24 24"
			role="img"
			aria-label="Widgets"
			{...props}
		>
			<path
				fill="currentColor"
				d="M13 13v8h8v-8zM3 21h8v-8H3zM3 3v8h8V3zm13.66-1.31L11 7.34L16.66 13l5.66-5.66z"
			/>
		</svg>
	);
}
