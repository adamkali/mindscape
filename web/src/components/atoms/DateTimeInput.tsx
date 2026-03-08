import type { ComponentProps, JSX } from 'solid-js';
import { cn } from '@/utils/cn';

interface DateTimeInputProps extends ComponentProps<'input'> {
	variant?: 'glass';
	label: JSX.Element;
}

export default function DateTimeInput(props: DateTimeInputProps): JSX.Element {
	const { class: className, ...inputProps } = props;

	return (
		<div class="flex-1 w-full">
			<label
				for={props.id}
				class="block text-xs font-medium text-foreground/60 mb-1"
			>
				{props.label}
			</label>
			<input
				{...inputProps}
				type="datetime-local"
				class={cn(
					'text-sm rounded-lg block w-full px-2.5 py-1.5 transition-all duration-300',
					'bg-glass-bg backdrop-blur-md border border-glass-border text-foreground',
					'placeholder-foreground/60',
					'hover:bg-glass-bg-hover hover:border-glass-border-hover',
					'focus:bg-glass-bg-hover focus:border-glass-border-hover focus:ring-2 focus:ring-foreground/20',
					'focus:outline-none focus:ring-offset-0',
					'datetime-glass',
					className,
				)}
			/>
		</div>
	);
}
