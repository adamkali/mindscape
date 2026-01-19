import type { ComponentProps, JSX } from 'solid-js';
import { cn } from '@/utils/cn';

export type CardPadding = 'none' | 'sm' | 'md' | 'lg';
export type CardSpacing = 'none' | 'sm' | 'md' | 'lg' | 'xl' | 'xxl';

/**
 * Props for the CardBody component
 */
interface CardBodyProps extends ComponentProps<'div'> {
	/** Body content */
	children?: JSX.Element;
	/** Padding variant for different content types */
	padding?: CardPadding;
	yspacing?: CardSpacing;
	xspacing?: CardSpacing;
}

/**
 * CardBody component - Main content area for cards
 *
 * Features:
 * - Flexible content container
 * - Configurable padding options
 * - Proper spacing and typography
 * - Uses design system colors
 *
 * @example
 * ```tsx
 * <CardBody padding="lg">
 *   <p>Main card content goes here</p>
 * </CardBody>
 * ```
 */
export default function CardBody(props: CardBodyProps): JSX.Element {
	const {
		children,
		padding = 'md',
		yspacing = 'none',
		xspacing = 'none',
		class: className,
		...bodyProps
	} = props;

	const getPaddingClasses = () => {
		switch (padding) {
			case 'none':
				return '';
			case 'sm':
				return 'px-4 py-3';
			case 'lg':
				return 'px-8 py-6';
			case 'md':
			default:
				return 'px-6 py-4';
		}
	};

	const getSpacingClasses = (xory: 'x' | 'y') => {
		switch (xory === 'x' ? xspacing : yspacing) {
			case 'none':
				return '';
			case 'sm':
				return 'space-' + xory + '-4';
			case 'md':
				return 'space-' + xory + '-6';
			case 'lg':
				return 'space-' + xory + '-8';
			case 'xl':
				return 'space-' + xory + '-10';
			case 'xxl':
				return 'space-' + xory + '-12';
			default:
				return 'space-y-3';
		}
	};

	return (
		<div
			class={cn(
				cn(getPaddingClasses(), className),
				getSpacingClasses('y'),
				getSpacingClasses('x'),
			)}
			{...bodyProps}
		>
			{children}
		</div>
	);
}
