import type { ComponentProps } from 'solid-js';
import { cn } from '@/utils/cn';

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
	return cn(
		`bg-glass-bg backdrop-blur-md border border-glass-border hover:bg-glass-bg-hover text-foreground font-bold py-2 px-4 rounded-lg
	shadow-lg hover:shadow-xl transition-all duration-300 shadow-slate-900/20 ease-out hover:scale-105 active:scale-95 hover:border-glass-border-hover
	dark:shadow-black/30`,
		passedClass,
	);
}

function secondary(passedClass?: string) {
	return cn(
		`bg-glass-bg/75 backdrop-blur-md border border-glass-border/85 hover:bg-glass-bg-hover text-foreground/90 hover:text-foreground font-bold py-2 px-4 rounded-lg
	shadow-lg hover:shadow-xl transition-all duration-300 shadow-slate-900/20 ease-out hover:scale-105 active:scale-95 hover:border-glass-border-hover
	dark:shadow-black/30`,
		passedClass,
	);
}

function tertiary(passedClass?: string) {
	return cn(
		`bg-glass-bg/50 backdrop-blur-md border border-glass-border/65 hover:bg-glass-bg text-foreground/80 hover:text-foreground font-bold py-2 px-4 rounded-lg
	shadow-lg hover:shadow-xl transition-all duration-300 shadow-slate-900/20 ease-out hover:scale-105 active:scale-95 hover:border-glass-border
	dark:shadow-black/30`,
		passedClass,
	);
}

function danger(passedClass?: string) {
	return cn(
		`bg-red-500/30 backdrop-blur-md border border-red-400/40 hover:bg-red-500/40 text-foreground font-bold py-2 px-4 rounded-lg
	shadow-lg hover:shadow-xl transition-all duration-300 shadow-red-900/20 ease-out hover:scale-105 active:scale-95 hover:border-red-400/60
	dark:bg-red-900/40 dark:border-red-700/50 dark:hover:bg-red-900/55 dark:hover:border-red-600/60 dark:shadow-red-950/30`,
		passedClass,
	);
}
