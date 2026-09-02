import { defineConfig } from '@rsbuild/core';
import { pluginBabel } from '@rsbuild/plugin-babel';
import { pluginSolid } from '@rsbuild/plugin-solid';

export default defineConfig({
	plugins: [
		pluginBabel({
			include: /\.(?:jsx|tsx)$/,
		}),
		pluginSolid(),
	],
	resolve: {
		alias: {
			'@': './src',
		},
	},
	source: {
		// Feature flags are inlined at build time. Rsbuild does NOT provide a
		// Vite-style `import.meta.env`, so each flag is defined explicitly here
		// and read through src/config/features.ts. Unset === disabled.
		define: {
			'process.env.MINDSCAPE_FEATURE_TASKS': JSON.stringify(
				process.env.MINDSCAPE_FEATURE_TASKS ?? '',
			),
			'process.env.MINDSCAPE_FEATURE_APIKEYS': JSON.stringify(
				process.env.MINDSCAPE_FEATURE_APIKEYS ?? '',
			),
			'process.env.MINDSCAPE_ADMIN': JSON.stringify(
				process.env.MINDSCAPE_ADMIN ?? '',
			),
		},
	},
	html: {
		title: 'Mindscape',
	},
	server: {
		port: 5173,
		host: '0.0.0.0',
		cors: true,
		proxy: {
			'/api': {
				target: process.env.BACKEND_URL || 'http://0.0.0.0:60000',
				changeOrigin: true,
				secure: false,
			},
			'/assets/svgs': {
				target: process.env.BACKEND_URL || 'http://0.0.0.0:60000',
				changeOrigin: true,
				secure: false,
			},
		},
	},
});
