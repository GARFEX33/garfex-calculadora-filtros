import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		host: true,
		port: 5173,
		proxy: {
			'/api': {
				target: 'http://192.168.1.71:8080',
				changeOrigin: true
			}
		}
	}
});
