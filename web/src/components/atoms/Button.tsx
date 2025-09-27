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
	return cn(`bg-white/20 backdrop-blur-sm border border-white/30 hover:bg-white/30 text-white hover:text-white font-bold py-2 px-4 rounded 
	shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out`, passedClass);
}

function secondary(passedClass?: string) {
	return cn(`bg-white/15 backdrop-blur-sm border border-white/25 hover:bg-white/25 text-white/90 hover:text-white font-bold py-2 px-4 rounded 
	shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out`, passedClass);
}

function tertiary(passedClass?: string) {
	return cn(`bg-white/10 backdrop-blur-sm border border-white/20 hover:bg-white/20 text-white/80 hover:text-white font-bold py-2 px-4 rounded 
	shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out`, passedClass);
}

function danger(passedClass?: string) {
	return cn(`bg-red-500/30 backdrop-blur-sm border border-red-400/40 hover:bg-red-500/40 text-white hover:text-white font-bold py-2 px-4 rounded 
	shadow-md hover:shadow-lg transition-all duration-200 shadow-slate-900/80 ease-in-out`, passedClass);
}
