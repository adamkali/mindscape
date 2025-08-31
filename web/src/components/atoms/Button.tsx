import { cn } from '@/utils/cn';
import { type ComponentProps } from 'solid-js';

interface ButtonProps extends ComponentProps<'button'> {
	variant?: 'primary' | 'secondary' | 'tertiary' | 'danger';
}

export default function Button(props: ButtonProps) {
	switch (props.variant) {
		case 'primary':
			return <button {...props} class={primary(props.class)} />;
		case 'secondary':
			return <button {...props} class={secondary(props.class)} />;
		case 'tertiary':
			return <button {...props} class={tertiary(props.class)} />;
		case 'danger':
			return <button {...props} class={danger(props.class)} />;
		default:
			return <button {...props} class={primary(props.class)} />;
	}
}

function primary(passedClass?: string) {
	return cn(`bg-primary hover:bg-primary-hover text-primary-foreground hover:text-primary-hover-foreground font-bold py-2 px-4 rounded 
	shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out`, passedClass);
}

function secondary(passedClass?: string) {
	return cn(`bg-secondary hover:bg-secondary-hover text-secondary-foreground hover:text-secondary-hover-foreground font-bold py-2 px-4 rounded 
	shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out`, passedClass);
}

function tertiary(passedClass?: string) {
	return cn(`bg-tertiary hover:bg-tertiary-hover text-tertiary-foreground hover:text-tertiary-hover-foreground font-bold py-2 px-4 rounded 
	shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out`, passedClass);

}

function danger(passedClass?: string) {
	return cn(`bg-error hover:bg-error-hover text-error-foreground hover:text-error-hover-foreground font-bold py-2 px-4 rounded 
	shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out`, passedClass);
}
