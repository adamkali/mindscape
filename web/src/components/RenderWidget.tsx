import { createSignal, Match, onMount, Show, Switch } from 'solid-js';
import {
    type ResponsesUserWidgetData,
    type ResponsesWidgetData,
    type WidgetResponse,
    WidgetsApi,
} from '@/api';
import { useAuth } from '@/contexts/AuthContext';
import GitHubProfileWidget from './atoms/widgets/GitHubProfileWidget';
import SearchbarWidget from './atoms/widgets/SearchbarWidget';
import type { ElementSize } from './types';

// NOTE: Best-case scenario sanitization function
// TODO: Implement proper validation and sanitization based on widget schema
// This should validate against the widget's schema requirements and sanitize user input
function sanitizeWidgetConfig(config: { [key: string]: object } | undefined): {
    [key: string]: any;
} {
    if (!config) return {};

    // Best case: assume config is already safe
    // In production, validate each property against schema
    // - Check types match schema requirements
    // - Sanitize string values (escape HTML, validate URLs)
    // - Validate numeric ranges
    // - Ensure required fields are present
    console.log('config', config);
    return config;
}

// Widget state machine - determines which widget component to render
type WidgetType = 'none' | 'searchbar' | 'githubprofile' | string;

function resolveWidgetType(schemaType: string | undefined): WidgetType {
    if (!schemaType) return 'none';

    // Map schema types to widget types
    const type = schemaType.toLowerCase();

    // Handle searchbar types
    if (type === 'searchbar' || type === 'search') {
        return 'searchbar';
    }

    // Handle github profile types (all variants map to same component)
    if (type.startsWith('githubprofile')) {
        return 'githubprofile';
    }

    return 'none';
}

export default function RenderWidget({
    widget,
    elementSize,
    spacing,
}: {
    widget: ResponsesUserWidgetData;
    elementSize: ElementSize;
    spacing: number;
}) {
    const [widgetSchema, setWidgetSchema] = createSignal<ResponsesWidgetData>();
    const auth = useAuth();

    // Calculate actual pixel dimensions based on grid span and square size
    const containerWidth = () => (widget.width ?? 1) * elementSize.width;
    const containerHeight = () => (widget.height ?? 1) * elementSize.height;

    // Grid container handles sizing, so we just fill 100% of the grid area
    const widgetWidth = () => '100%';
    const widgetHeight = () => '100%';

    const getWidgetSchema = async () => {
        const api = new WidgetsApi();
        const response = await api.getWidgetSchemaByID({
            schemaId: widget.schemaId ?? '',
        });
        if (response.success && response.data) {
            return response.data;
        } else {
            console.error('Failed to get widget schema:', response.message);
        }
    };

    onMount(async () => {
        const schema = await getWidgetSchema();
        if (schema) {
            setWidgetSchema(schema);
        }
    });

    // Widget state machine: sanitize and resolve widget type based on schema
    const widgetType = () => resolveWidgetType(widgetSchema()?.type);
    const sanitizedConfig = () => sanitizeWidgetConfig(widget.config);


    return (
        <Show
            when={widget.id}
            fallback={<div class="text-white text-sm">No widget</div>}
        >
            <Switch
                fallback={
                    <div class="bg-gray-500/20 border-2 border-gray-500/50 rounded-xl p-8 text-white text-center w-full h-full">
                        <div class="text-4xl mb-4">❓</div>
                        <h3 class="text-xl font-bold mb-2">Unknown Widget Type</h3>
                        <p class="text-sm text-white/70">
                            Widget type "{widgetSchema()?.type}" is not recognized.
                        </p>
                    </div>
                }
            >
                <Match when={widgetType() === 'none'}>
                    <div class="bg-gray-500/20 border-2 border-gray-500/50 rounded-xl p-8 text-white text-center">
                        <div class="text-4xl mb-4">❓</div>
                        <h3 class="text-xl font-bold mb-2">Widget Type: None</h3>
                        <p class="text-sm text-white/70">
                            Schema type "{widgetSchema()?.type}" maps to no renderable widget.
                        </p>
                    </div>
                </Match>

                <Match when={widgetType() === 'searchbar'}>
                    <div
                        class="backdrop-blur-sm border-2 rounded-xl shadow-lg
						items-center justify-center"
                        style={{
                            width: widgetWidth(),
                            height: widgetHeight(),
                            padding: spacing + 'px',
                            'z-index': widget.zIndex ?? 1,
                        }}
                    >
                        <SearchbarWidget
                            url={
                                sanitizedConfig().serverUrl ||
                                'https://www.google.com/search?q=%s'
                            }
                        />
                    </div>
                </Match>

                <Match when={widgetType() === 'githubprofile'}>
                    <div
                        style={{
                            width: widgetWidth(),
                            height: widgetHeight(),
                            margin: spacing + 'px',
                            'z-index': widget.zIndex ?? 1,
                        }}
                        class="w-full h-full"
                    >
                        <GitHubProfileWidget
                            widgetId={widget.id ?? ''}
                            authToken={auth.token() ?? ''}
                            containerWidth={containerWidth()}
                            containerHeight={containerHeight()}
                        />
                    </div>
                </Match>
            </Switch>
        </Show>
    );
}
