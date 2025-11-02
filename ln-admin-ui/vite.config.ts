import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// ============================================
// 📌 后端服务配置 - 只在这里修改后端地址！
// ============================================
const BACKEND_URL = 'http://localhost:8080'

// https://vite.dev/config/
export default defineConfig({
    plugins: [vue(), vueDevTools()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
    },
    server: {
        port: 5173, // 开发服务器端口
        host: true, // 监听所有地址
        open: true, // 自动打开浏览器
        cors: true, // 启用跨域
        proxy: {
            // 代理所有以 /api 开头的请求
            '/api': {
                target: BACKEND_URL, // 统一使用这里的后端地址
                changeOrigin: true, // 改变源地址
                rewrite: (path) => path, // 不重写路径
                configure: (proxy, _options) => {
                    proxy.on('error', (err, _req, _res) => {
                        console.log('proxy error', err)
                    })
                    proxy.on('proxyReq', (proxyReq, req, _res) => {
                        console.log('Sending Request to the Target:', req.method, req.url)
                    })
                    proxy.on('proxyRes', (proxyRes, req, _res) => {
                        console.log('Received Response from the Target:', proxyRes.statusCode, req.url)
                    })
                },
            },
        },
    },
    build: {
        outDir: 'dist',
        assetsDir: 'assets',
        sourcemap: false,
        minify: 'terser',
        chunkSizeWarningLimit: 1500,
        rollupOptions: {
            output: {
                chunkFileNames: 'assets/js/[name]-[hash].js',
                entryFileNames: 'assets/js/[name]-[hash].js',
                assetFileNames: 'assets/[ext]/[name]-[hash].[ext]',
            },
        },
    },
})
