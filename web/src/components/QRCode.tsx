import QRCodeLib from 'qrcode';
import { type Component, createResource, Show } from 'solid-js';

interface QRCodeProps {
	/** The value to encode in the QR code (URL, code, etc.). */
	value: string;
	/** Rendered width/height in pixels. Defaults to 200. */
	size?: number;
	/** Accessible label for the generated image. */
	label?: string;
	class?: string;
}

/**
 * Renders a QR code as a data-URL image. The encoded value is passed to the
 * `qrcode` library at request time, so the image regenerates whenever `value`
 * or `size` changes.
 */
export const QRCode: Component<QRCodeProps> = (props) => {
	const [dataUrl] = createResource(
		() => ({ value: props.value, size: props.size ?? 200 }),
		async ({ value, size }) => {
			if (!value) return '';
			return QRCodeLib.toDataURL(value, {
				width: size,
				margin: 1,
				errorCorrectionLevel: 'M',
			});
		},
	);

	return (
		<Show
			when={dataUrl()}
			fallback={
				<div
					class="flex items-center justify-center bg-glass-bg border border-glass-border rounded-lg text-xs text-foreground/50"
					style={{
						width: `${props.size ?? 200}px`,
						height: `${props.size ?? 200}px`,
					}}
				>
					Generating…
				</div>
			}
		>
			<img
				src={dataUrl()}
				alt={props.label ?? 'QR code'}
				width={props.size ?? 200}
				height={props.size ?? 200}
				class={`rounded-lg bg-white p-2 ${props.class ?? ''}`}
			/>
		</Show>
	);
};

export default QRCode;
