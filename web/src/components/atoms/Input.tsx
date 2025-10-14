import { cn } from '@/utils/cn';
import { type ComponentProps, type JSX } from 'solid-js';

interface InputProps extends ComponentProps<'input'> {
	variant: 'primary' | 'secondary' | 'tertiary' | 'danger';
	label: JSX.Element;
}

export default function Input(props: InputProps): JSX.Element {
	return (
		<div>
			<label
				for={props.id}
				class="block text-sm font-medium bg-primary text-slate-100"
			>
				{props.label}
			</label>
			<input
				{...props}
				class={cn(
					`bg-gray-50 border border-gray-300 text-gray-900
					text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500
					block w-full p-2.5`,
					props.class,
				)}
			/>
		</div>
	);
}
