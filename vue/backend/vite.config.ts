import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import VueSetupExtend from 'vite-plugin-vue-setup-extend';
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers';

export default defineConfig({
	base: './',
	build: {
		outDir: '../../public/backend/', // 输出目录
		rollupOptions: {
			output: {
				manualChunks: undefined, // 所有文件打包到一个JS文件中
			},
		},
	},
	plugins: [
		vue(),
		VueSetupExtend(),
		AutoImport({
			resolvers: [ElementPlusResolver()]
		}),
		Components({
			resolvers: [ElementPlusResolver()]
		})
	],
	optimizeDeps: {
		include: ['schart.js']
	},
	server: {
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:9090',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, '/task/api')
			}
		}
	},
	resolve: {
		alias: {
			'@': '/src',
			'~': '/src/assets'
		}
	},
	define: {
		__VUE_PROD_HYDRATION_MISMATCH_DETAILS__: "true",
	},
});
