import { cn } from '@/utils/cn';
import { splitProps, type ComponentProps, type JSX } from 'solid-js';

/**
 * Props for the CardHeader component
 */
interface CardHeaderProps extends ComponentProps<'div'> {
	/** Header title */
	title?: string;
	/** Header subtitle or description */
	subtitle?: string;
	/** Additional content (e.g., actions, icons) */
	children?: JSX.Element;
}

/**
 * CardHeader component - Header section for cards
 *
 * Features:
 * - Title and subtitle support
 * - Additional content slot for actions or icons
 * - Proper spacing and typography hierarchy
 * - Uses design system colors
 *
 * @example
 * ```tsx
 * <CardHeader
 *   title="Card Title"
 *   subtitle="Card description"
 * >
 *   <button>Action</button>
 * </CardHeader>
 * ```
 */
export default function CardHeader(props: CardHeaderProps): JSX.Element {
	const { title, subtitle, children, class: className, ...headerProps } = props;

	return (
		<div
			class={cn('px-6 py-4 border-b border-white/10', className)}
			{...headerProps}
		>
			<div class="flex items-start justify-between">
				<div class="flex-1">
					{title && (
						<h3 class="text-lg font-semibold leading-6">
							{title}
						</h3>
					)}
					{subtitle && (
						<p class="text-sm opacity-70 mt-1">{subtitle}</p>
					)}
				</div>
				{children && <div class="ml-4 flex-shrink-0">{children}</div>}
			</div>
		</div>
	);
}
