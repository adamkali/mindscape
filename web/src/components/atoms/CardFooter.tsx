import { cn } from '@/utils/cn';
import { type ComponentProps, type JSX } from 'solid-js';

/**
 * Props for the CardFooter component
 */
interface CardFooterProps extends ComponentProps<'div'> {
	/** Footer content */
	children?: JSX.Element;
	/** Layout direction for footer content */
	direction?: 'row' | 'column';
	/** Content alignment */
	align?: 'start' | 'center' | 'end' | 'between';
}

/**
 * CardFooter component - Footer section for cards with action buttons
 *
 * Features:
 * - Flexible layout options (row/column)
 * - Multiple alignment options
 * - Proper spacing for buttons and actions
 * - Uses design system colors
 *
 * @example
 * ```tsx
 * <CardFooter align="end">
 *   <button class="btn-secondary">Cancel</button>
 *   <button class="btn-primary">Save</button>
 * </CardFooter>
 * ```
 */
export default function CardFooter(props: CardFooterProps): JSX.Element {
	const {
		children,
		direction = 'row',
		align = 'end',
		class: className,
		...footerProps
	} = props;

	const getDirectionClasses = () => {
		return direction === 'column' ? 'flex-col space-y-2' : 'flex-row space-x-3';
	};

	const getAlignmentClasses = () => {
		if (direction === 'column') {
			switch (align) {
				case 'start':
					return 'items-start';
				case 'center':
					return 'items-center';
				case 'end':
					return 'items-end';
				case 'between':
					return 'items-stretch';
				default:
					return 'items-end';
			}
		} else {
			switch (align) {
				case 'start':
					return 'justify-start';
				case 'center':
					return 'justify-center';
				case 'end':
					return 'justify-end';
				case 'between':
					return 'justify-between';
				default:
					return 'justify-end';
			}
		}
	};

	return (
		<div
			class={cn(
				'px-6 py-4 border-t border-border/10 bg-card/50 rounded-b-lg ',
				className,
			)}
			{...footerProps}
		>
			<div
				class={cn('flex', getDirectionClasses() + ' ' + getAlignmentClasses())}
			>
				{children}
			</div>
		</div>
	);
}
