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
	server: {
		port: 5173,
		cors: true,
		proxy: {
			'/api': {
				target: 'http://0.0.0.0:60000',
				changeOrigin: false,
				secure: false,
			},
			'/assets': {
				target: 'http://0.0.0.0:60000',
				changeOrigin: false,
				secure: false,
			},
		},
	},
});
