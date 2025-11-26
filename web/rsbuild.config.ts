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
