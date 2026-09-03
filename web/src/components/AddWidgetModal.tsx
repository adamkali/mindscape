import { createSignal, For, Show } from 'solid-js';
import { useWidgets } from '@/contexts/WidgetContext';
import { Button } from './atoms';

interface AddWidgetModalProps {
	isOpen: boolean;
	onClose: () => void;
	widgetLength: number;
}

export default function AddWidgetModal(props: AddWidgetModalProps) {
	const { schemas, schemasLoading, addWidget } = useWidgets();
	const [selectedSchemaId, setSelectedSchemaId] = createSignal<string>('');
	const [widgetConfig, setWidgetConfig] = createSignal<Record<string, any>>({});
	const [submitting, setSubmitting] = createSignal(false);
	// Raw text the user typed into each `array` property's input. View-only state:
	// widgetConfig holds the parsed array, this holds the in-progress text so the
	// input does not fight the user mid-keystroke (e.g. while typing a separator).
	const [arrayText, setArrayText] = createSignal<
		Record<string, string | undefined>
	>({});

	const handleConfigChange = (key: string, value: any) => {
		setWidgetConfig({ ...widgetConfig(), [key]: value });
	};

	// `array` properties are edited as a comma separated list (matching the
	// "Comma separated list" schema copy) but MUST serialize to a real JSON array:
	// the backend unmarshals these keys into Go `[]string`, so submitting a bare
	// string fails with "cannot unmarshal string into []string" and the widget
	// 400s. Empty input therefore yields `[]`, never `''` — and an empty array is
	// the documented "show everything" value, so it stays a valid submission.
	const parseArrayText = (text: string): string[] =>
		text
			.split(',')
			.map((entry) => entry.trim())
			.filter((entry) => entry.length > 0);

	// Schemas disagree on the default for an array property: coolify ships `''`,
	// searxng ships `[]`. Accept either (and anything else) without crashing.
	const arrayValueToText = (value: unknown): string => {
		if (Array.isArray(value)) {
			return value.map((entry) => String(entry)).join(', ');
		}
		if (typeof value === 'string') return value;
		return '';
	};

	const arrayTextFor = (key: string, schemaValue: unknown): string =>
		arrayText()[key] ?? arrayValueToText(schemaValue);

	const handleArrayChange = (key: string, text: string) => {
		setArrayText({ ...arrayText(), [key]: text });
		handleConfigChange(key, parseArrayText(text));
	};

	const handleSubmit = async (e: Event) => {
		e.preventDefault();

		if (!selectedSchemaId()) {
			alert('Please select a widget type');
			return;
		}

		// Get the selected schema to use its default layout
		const selectedSchema = schemas().find((s) => s.id === selectedSchemaId());
		if (!selectedSchema) {
			alert('Invalid widget schema');
			return;
		}

		// Validate required fields
		const requiredFields = selectedSchema.required || [];
		for (const field of requiredFields) {
			if (!widgetConfig()[field]) {
				alert(`Please fill in required field: ${field}`);
				return;
			}
		}

		setSubmitting(true);
		const ok = await addWidget({
			schemaId: selectedSchemaId(),
			schemaTitle: selectedSchema.title,
			config: Array.from(
				new TextEncoder().encode(JSON.stringify(widgetConfig())),
			),
			positionX: 0, // TODO: Calculate next available position
			positionY: props.widgetLength,
			width: selectedSchema.layout?.defaultSize?.width ?? 1,
			height: selectedSchema.layout?.defaultSize?.height ?? 1,
			zIndex: 1,
			isVisible: true,
		});
		setSubmitting(false);

		if (ok) {
			props.onClose();
			setWidgetConfig({});
			setArrayText({});
			setSelectedSchemaId('');
		} else {
			alert('Failed to create widget');
		}
	};

	const handleBackdropClick = (e: MouseEvent) => {
		if (e.target === e.currentTarget) {
			props.onClose();
		}
	};

	return (
		<Show when={props.isOpen}>
			<div
				class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50"
				onClick={handleBackdropClick}
			>
				<div
					class="bg-gradient-to-br from-card to-card/80 backdrop-blur-lg border-2 border-slate-700/20 rounded-2xl shadow-2xl w-full max-w-md p-6"
					onClick={(e) => e.stopPropagation()}
				>
					<div class="flex justify-between items-center mb-6">
						<h2 class="text-2xl font-bold text-card-foreground">Add Widget</h2>
						<button
							onClick={props.onClose}
							class="text-card-foreground/60 hover:text-card-foreground transition-colors"
						>
							<svg
								class="w-6 h-6"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
						</button>
					</div>

					<Show
						when={!schemasLoading()}
						fallback={
							<div class="text-center py-8 text-card-foreground/60">
								Loading widget types...
							</div>
						}
					>
						<form onSubmit={handleSubmit} class="space-y-4">
							<div>
								<label
									for="widget-type"
									class="block text-sm font-semibold text-card-foreground mb-2"
								>
									Widget Type
								</label>
								<select
									id="widget-type"
									class="w-full px-4 py-3 rounded-lg border-2 bg-slate-300/30 text-foreground focus:outline-none transition-all"
									value={selectedSchemaId()}
									onChange={(e) => setSelectedSchemaId(e.currentTarget.value)}
								>
									<option value="">Select a widget type...</option>
									<For each={schemas()}>
										{(schema) => (
											<option
												value={schema.id}
												class="text-slate-300/60 bg-slate-600/90 hover:bg-slate-300/20 hover:text-card-foreground transition-colors"
											>
												{schema.title || schema.type} ({schema.type})
											</option>
										)}
									</For>
								</select>
							</div>

							<Show when={selectedSchemaId()}>
								<div class="bg-white/10 rounded-lg p-4 border border-white/20 space-y-4 dark:bg-slate-900/30 dark:border-slate-700/40">
									<div>
										<h3 class="text-sm font-semibold text-card-foreground mb-2">
											Widget Configuration
										</h3>
										<For each={schemas()}>
											{(schema) => (
												<Show when={schema.id === selectedSchemaId()}>
													<div class="space-y-3">
														{/* Widget Info */}
														<div class="text-xs text-card-foreground/70 pb-2 border-b border-white/10">
															<span class="font-medium">{schema.type}</span> -{' '}
															{schema.layout?.defaultSize?.width}w ×{' '}
															{schema.layout?.defaultSize?.height}h
														</div>

														{/* Dynamic form fields based on schema properties */}
														<For each={Object.entries(schema.properties || {})}>
															{([key, property]) => (
																<div>
																	<label
																		for={key}
																		class="block text-xs font-medium text-card-foreground mb-1"
																	>
																		{property.label || key}
																		{schema.required?.includes(key) && (
																			<span class="text-error ml-1">*</span>
																		)}
																	</label>

																	{/* Render input based on property type */}
																	<Show
																		when={
																			property.type === 'string' &&
																			property._enum &&
																			(property._enum as any[]).length > 0
																		}
																	>
																		<select
																			id={key}
																			class="w-full px-3 py-2 text-sm rounded-lg  bg-slate-600/40 text-foreground focus:outline-none"
																			value={
																				widgetConfig()[key] || property.value
																			}
																			onChange={(e) =>
																				handleConfigChange(
																					key,
																					e.currentTarget.value,
																				)
																			}
																		>
																			<For each={property._enum as any[]}>
																				{(option) => (
																					<option value={option}>
																						{option}
																					</option>
																				)}
																			</For>
																		</select>
																	</Show>

																	<Show
																		when={
																			property.type === 'string' &&
																			(!property._enum ||
																				(property._enum as any[]).length === 0)
																		}
																	>
																		<input
																			id={key}
																			type={
																				property.format === 'password'
																					? 'password'
																					: property.format === 'uri'
																						? 'url'
																						: 'text'
																			}
																			class="w-full px-3 py-2 text-sm rounded-lg  bg-slate-600/40 text-foreground focus:outline-none"
																			value={
																				widgetConfig()[key] ||
																				property.value ||
																				''
																			}
																			onInput={(e) =>
																				handleConfigChange(
																					key,
																					e.currentTarget.value,
																				)
																			}
																			placeholder={property.description}
																		/>
																	</Show>

																	<Show when={property.type === 'boolean'}>
																		<label class="flex items-center gap-2 cursor-pointer">
																			<input
																				id={key}
																				type="checkbox"
																				class="w-4 h-4 rounded border-primary/50 text-primary focus:ring-primary"
																				checked={
																					widgetConfig()[key] ??
																					property.value ??
																					true
																				}
																				onChange={(e) =>
																					handleConfigChange(
																						key,
																						e.currentTarget.checked,
																					)
																				}
																			/>
																			<span class="text-xs text-card-foreground/70">
																				{property.description}
																			</span>
																		</label>
																	</Show>

																	<Show
																		when={
																			property.type === 'number' ||
																			property.type === 'integer'
																		}
																	>
																		<input
																			id={key}
																			type="number"
																			class="w-full px-3 py-2 text-sm rounded-lg border-2 border-primary/50 bg-white/90 text-foreground focus:outline-none focus:border-primary"
																			value={
																				widgetConfig()[key] ||
																				property.value ||
																				0
																			}
																			onInput={(e) =>
																				handleConfigChange(
																					key,
																					parseFloat(e.currentTarget.value),
																				)
																			}
																			placeholder={property.description}
																		/>
																	</Show>

																	{/* Keyed off `type`, not `format`: the schemas
																	    disagree on format for array properties (coolify
																	    says "array", githubprs/searxng say "string"). */}
																	<Show when={property.type === 'array'}>
																		<input
																			id={key}
																			type="text"
																			class="w-full px-3 py-2 text-sm rounded-lg bg-slate-600/40 text-foreground focus:outline-none"
																			value={arrayTextFor(key, property.value)}
																			onInput={(e) =>
																				handleArrayChange(
																					key,
																					e.currentTarget.value,
																				)
																			}
																			placeholder={property.description}
																		/>
																	</Show>
																</div>
															)}
														</For>
													</div>
												</Show>
											)}
										</For>
									</div>
								</div>
							</Show>

							<div class="flex gap-3 pt-4">
								<Button
									type="button"
									variant="tertiary"
									onClick={props.onClose}
									class="flex-1"
								>
									Cancel
								</Button>
								<Button
									type="submit"
									variant="primary"
									class="flex-1"
									disabled={!selectedSchemaId() || submitting()}
								>
									{submitting() ? 'Adding...' : 'Add Widget'}
								</Button>
							</div>
						</form>
					</Show>
				</div>
			</div>
		</Show>
	);
}
