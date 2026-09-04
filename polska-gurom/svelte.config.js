import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const dev = process.argv.includes('dev');

/**
 * When deployed to GitHub Pages as a project site, the app is served from
 * https://<user>.github.io/<repo>/, so the base path must be `/<repo>`.
 * The Pages workflow sets BASE_PATH; locally we serve from the root.
 * @type {import('@sveltejs/kit').Config}
 */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      fallback: '404.html'
    }),
    paths: {
      base: dev ? '' : process.env.BASE_PATH || ''
    }
  }
};

export default config;
