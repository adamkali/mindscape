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
				return 'bg-white/20 backdrop-blur-md border border-white/30 text-white placeholder-white/70 focus:bg-white/30 focus:border-white/50 focus:ring-2 focus:ring-white/20 dark:bg-slate-900/40 dark:border-slate-700/50 dark:text-white dark:placeholder-white/60 dark:focus:bg-slate-900/60 dark:focus:border-slate-600/60 dark:focus:ring-slate-500/30';
			case 'primary':
				return 'bg-gray-50 border border-gray-300 text-gray-900 focus:ring-blue-500 focus:border-blue-500';
			case 'secondary':
				return 'bg-gray-100 border border-gray-400 text-gray-800 focus:ring-gray-500 focus:border-gray-500';
			case 'tertiary':
				return 'bg-gray-200 border border-gray-500 text-gray-700 focus:ring-gray-600 focus:border-gray-600';
			case 'danger':
				return 'bg-red-50 border border-red-300 text-red-900 focus:ring-red-500 focus:border-red-500';
			default:
				return 'bg-white/20 backdrop-blur-md border border-white/30 text-white placeholder-white/70 focus:bg-white/30 focus:border-white/50 focus:ring-2 focus:ring-white/20 dark:bg-slate-900/40 dark:border-slate-700/50 dark:text-white dark:placeholder-white/60 dark:focus:bg-slate-900/60 dark:focus:border-slate-600/60 dark:focus:ring-slate-500/30';
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
					'hover:bg-white/25 hover:border-white/40 dark:hover:bg-slate-900/50 dark:hover:border-slate-600/55',
					'focus:outline-none focus:ring-offset-0',
					getVariantClasses(),
					className,
				)}
			/>
		</div>
	);
}
