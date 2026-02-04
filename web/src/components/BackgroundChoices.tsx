import { For, Show } from 'solid-js';
import { useBackground } from '@/contexts/BackgroundContext';

export default function BackgroundChoices(props: {
	handleBackgroundSelect: (backgroundUrl: string) => Promise<void>;
}) {
	const { handleBackgroundSelect } = props;
	const { backgroundChoices, currentBackground } = useBackground();
	return (
		<Show when={backgroundChoices()}>
			<div class="grid grid-cols-2 md:grid-cols-3 gap-4">
				<For each={backgroundChoices()}>
					{(backgroundUrl) => (
						<div
							class={`relative aspect-video rounded-lg overflow-hidden cursor-pointer border-2 transition-all duration-300 hover:scale-105 ${
								currentBackground() === backgroundUrl.filename
									? 'border-white/70 ring-2 ring-white/50'
									: 'border-white/20 hover:border-white/40'
							}`}
							onClick={() =>
								handleBackgroundSelect(backgroundUrl.filename ?? '')
							}
						>
							<img
								src={backgroundUrl.url}
								alt="Background option"
								class="w-full h-full object-cover"
								onError={(e) => {
									(e.target as HTMLImageElement).style.display = 'none';
								}}
							/>
							<Show when={currentBackground() === backgroundUrl.filename}>
								<div class="absolute inset-0 bg-white/20 backdrop-blur-sm flex items-center justify-center">
									<div class="text-white font-semibold">Selected</div>
								</div>
							</Show>
						</div>
					)}
				</For>
			</div>
		</Show>
	);
}
