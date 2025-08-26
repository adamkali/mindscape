import { type ComponentProps } from 'solid-js';


interface ButtonProps extends ComponentProps<'button'> {
	variant?: 'primary' | 'secondary' | 'tertiary' | 'danger';
}

export default function Button(props: ButtonProps) {
	switch (props.variant) {
		case 'primary':
			return <button {...props} class={primary()} />
		case 'secondary':
			return <button {...props} class={secondary()} />
		case 'tertiary':
			return <button {...props} class={tertiary()} />
		case 'danger':
			return <button {...props} class={danger()} />
		default:
			return <button {...props} class={primary()} />
	}
}

function primary() {
	return 'bg-primary hover:bg-primary-dark text-primary-foreground hover:text-primary-foreground font-bold py-2 px-4 rounded';
}

function secondary() {
	return 'bg-secondary hover:bg-secondary-dark text-secondary-foreground hover:text-secondary-foreground font-bold py-2 px-4 rounded';
}

function tertiary() {
	return 'bg-tertiary hover:bg-tertiary-dark text-tertiary-foreground hover:text-tertiary-foreground font-bold py-2 px-4 rounded';
}

function danger() {
	return 'bg-danger hover:bg-danger-dark text-danger-foreground hover:text-danger-foreground font-bold py-2 px-4 rounded';
}
