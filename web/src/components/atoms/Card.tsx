import { cn } from '@/utils/cn';
import { type ComponentProps, type JSX } from 'solid-js';

/**
 * Props for the Card component
 */
interface CardProps extends ComponentProps<'div'> {
	/** Visual variant of the card */
	variant?: 'default' | 'outlined' | 'elevated';
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
			case 'outlined':
				return 'border border-border';
			case 'elevated':
				return 'shadow-lg border border-border/50';
			case 'default':
			default:
				return 'shadow-sm';
		}
	};

	return (
		<div
			class={cn(
				'bg-card text-card-foreground rounded-lg',
				getVariantClasses() + ' ' + className,
			)}
			{...cardProps}
		>
			{children}
		</div>
	);
}
