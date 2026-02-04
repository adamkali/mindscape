import type { ComponentProps, JSX } from 'solid-js';
import { cn } from '@/utils/cn';

interface InputProps extends ComponentProps<'input'> {
	variant?: 'primary' | 'secondary' | 'tertiary' | 'danger' | 'glass';
	label: JSX.Element;
}

export default function Input(props: InputProps): JSX.Element {
	const { variant = 'glass', class: className, ...inputProps } = props;

	const getVariantClasses = () => {
		switch (variant) {
			case 'glass':
				return 'bg-glass-bg backdrop-blur-md border border-glass-border text-foreground placeholder-foreground/60 focus:bg-glass-bg-hover focus:border-glass-border-hover focus:ring-2 focus:ring-foreground/20';
			case 'primary':
				return 'bg-gray-50 border border-gray-300 text-gray-900 focus:ring-blue-500 focus:border-blue-500';
			case 'secondary':
				return 'bg-gray-100 border border-gray-400 text-gray-800 focus:ring-gray-500 focus:border-gray-500';
			case 'tertiary':
				return 'bg-gray-200 border border-gray-500 text-gray-700 focus:ring-gray-600 focus:border-gray-600';
			case 'danger':
				return 'bg-red-50 border border-red-300 text-red-900 focus:ring-red-500 focus:border-red-500';
			default:
				return 'bg-glass-bg backdrop-blur-md border border-glass-border text-foreground placeholder-foreground/60 focus:bg-glass-bg-hover focus:border-glass-border-hover focus:ring-2 focus:ring-foreground/20';
		}
	};

	return (
		<div class="flex-1 w-full">
			<label for={props.id} class="sr-only">
				{props.label}
			</label>
			<input
				{...inputProps}
				class={cn(
					'text-sm rounded-lg block w-full px-2.5 py-1 transition-all duration-300',
					'hover:bg-glass-bg-hover hover:border-glass-border-hover',
					'focus:outline-none focus:ring-offset-0',
					getVariantClasses(),
					className,
				)}
			/>
		</div>
	);
}
