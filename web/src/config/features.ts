/**
 * Build-time feature flags for the launch build.
 *
 * These mirror the backend's `features:` config block (see
 * `cmd/configuration/configuration.go`). The two are deliberately separate:
 * the backend reads YAML and has no env-override layer, while the frontend
 * needs values inlined at build time. Keep them in sync when flipping a flag.
 *
 * Values are injected by `source.define` in `rsbuild.config.ts` — Rsbuild has
 * no Vite-style `import.meta.env`, so do not reach for it here.
 *
 * Everything defaults to OFF: an unset var inlines to '' , which is not 'true'.
 */
const flag = (value: string | undefined): boolean => value === 'true';

const isDev = process.env.NODE_ENV !== 'production';

export const FEATURES = {
	/** Tasks / Agenda view. Backend: `features.tasks`. */
	tasks: flag(process.env.MINDSCAPE_FEATURE_TASKS),
	/** API key management page + key-auth routes. Backend: `features.apikeys`. */
	apiKeys: flag(process.env.MINDSCAPE_FEATURE_APIKEYS),
	/**
	 * Admin-only pages. Always available in dev; a production build needs an
	 * explicit opt-in so the showcase can never be reached by a deployed user.
	 */
	admin: isDev || flag(process.env.MINDSCAPE_ADMIN),
} as const;
