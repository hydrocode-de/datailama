import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter({
			pages: '../internal/web/site/frontend',
			assets: '../internal/web/site/frontend',
			fallback: 'index.html',
			precompress: false,
		}),
		paths: {
			base: '/app'
		}
	}
};

export default config;
