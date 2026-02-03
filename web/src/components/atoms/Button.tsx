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
		`bg-white/20 backdrop-blur-md border border-white/30 hover:bg-white/30 text-white hover:text-white font-bold py-2 px-4 rounded-lg
	shadow-lg hover:shadow-xl transition-all duration-300 shadow-slate-900/20 ease-out hover:scale-105 active:scale-95 hover:border-white/50
	dark:bg-slate-900/40 dark:border-slate-700/50 dark:hover:bg-slate-900/60 dark:hover:border-slate-600/60 dark:text-white dark:shadow-black/30`,
		passedClass,
	);
}

function secondary(passedClass?: string) {
	return cn(
		`bg-white/15 backdrop-blur-md border border-white/25 hover:bg-white/25 text-white/90 hover:text-white font-bold py-2 px-4 rounded-lg
	shadow-lg hover:shadow-xl transition-all duration-300 shadow-slate-900/20 ease-out hover:scale-105 active:scale-95 hover:border-white/40
	dark:bg-slate-900/35 dark:border-slate-700/45 dark:hover:bg-slate-900/50 dark:hover:border-slate-600/55 dark:text-white/90 dark:hover:text-white dark:shadow-black/30`,
		passedClass,
	);
}

function tertiary(passedClass?: string) {
	return cn(
		`bg-white/10 backdrop-blur-md border border-white/20 hover:bg-white/20 text-white/80 hover:text-white font-bold py-2 px-4 rounded-lg
	shadow-lg hover:shadow-xl transition-all duration-300 shadow-slate-900/20 ease-out hover:scale-105 active:scale-95 hover:border-white/35
	dark:bg-slate-900/30 dark:border-slate-700/40 dark:hover:bg-slate-900/45 dark:hover:border-slate-600/50 dark:text-white/80 dark:hover:text-white dark:shadow-black/30`,
		passedClass,
	);
}

function danger(passedClass?: string) {
	return cn(
		`bg-red-500/30 backdrop-blur-md border border-red-400/40 hover:bg-red-500/40 text-white hover:text-white font-bold py-2 px-4 rounded-lg
	shadow-lg hover:shadow-xl transition-all duration-300 shadow-red-900/20 ease-out hover:scale-105 active:scale-95 hover:border-red-400/60
	dark:bg-red-900/40 dark:border-red-700/50 dark:hover:bg-red-900/55 dark:hover:border-red-600/60 dark:shadow-red-950/30`,
		passedClass,
	);
}
