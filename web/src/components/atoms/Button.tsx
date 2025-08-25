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
    return 'bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded';
}

function secondary() {
    return 'bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded';
}

function tertiary() {
    return 'bg-gray-200 hover:bg-gray-400 text-gray-900 font-bold py-2 px-4 rounded';
}

function danger() {
    return 'bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded';
}
