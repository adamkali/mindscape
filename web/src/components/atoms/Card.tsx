import type { ComponentProps, JSX } from 'solid-js';
import { cn } from '@/utils/cn';

/**
 * Props for the Card component
 */
interface CardProps extends ComponentProps<'div'> {
	/** Visual variant of the card */
	variant?: 'default' | 'glass' | 'elevated' | 'outlined';
	/** Card content */
	children?: JSX.Element;
}

/**
 * Card component - Main container for card-based layouts
 *
 * Features:
 * - Multiple visual variants (default, outlined, elevated)
 * - Uses design system colors (bg-card, text-card-foreground)
 * - Flexible container for CardHeader, CardBody, and CardFooter
 *
 * @example
 * ```tsx
 * <Card variant="elevated">
 *   <CardHeader>Header content</CardHeader>
 *   <CardBody>Main content</CardBody>
 *   <CardFooter>Footer content</CardFooter>
 * </Card>
 * ```
 */
export default function Card(props: CardProps): JSX.Element {
	const {
		variant = 'default',
		children,
		class: className,
		...cardProps
	} = props;

	const getVariantClasses = () => {
		switch (variant) {
			case 'glass':
				return 'bg-glass-bg backdrop-blur-md border border-glass-border shadow-lg hover:shadow-xl transition-all duration-300 shadow-slate-900/20 hover:bg-glass-bg-hover hover:border-glass-border-hover dark:shadow-black/30';
			case 'outlined':
				return 'border border-border bg-card';
			case 'elevated':
				return 'shadow-lg border border-border/50 bg-card';
			case 'default':
			default:
				return 'shadow-sm bg-card';
		}
	};

	return (
		<div
			class={cn(
				'text-card-foreground rounded-md',
				variant === 'glass' ? 'text-foreground' : '',
				getVariantClasses(),
				className,
			)}
			{...cardProps}
		>
			{children}
		</div>
	);
}
