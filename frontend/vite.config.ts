import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit()
	],
	server: {
		port: 3000,
		host: '0.0.0.0',
		proxy: {
			'/api/v1': {
				target: process.env.BACKEND_URL || 'http://localhost:8080',
				changeOrigin: true
			},
			'/uploads': {
				target: process.env.BACKEND_URL || 'http://localhost:8080',
				changeOrigin: true
			}
		}
	}
});
