import type { JSX } from 'solid-js';

export default function AgendaIcon(props: JSX.IntrinsicElements['svg']) {
	return (
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="1em"
			height="1em"
			viewBox="0 0 24 24"
			role="img"
			aria-label="Agenda"
			{...props}
		>
			<path fill="currentColor" d="M3 18h18v-2H3zm0-5h18v-2H3zm0-7v2h18V6z" />
		</svg>
	);
}
