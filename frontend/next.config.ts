import type { NextConfig } from "next";
import withPWA, { runtimeCaching } from "@ducanh2912/next-pwa";

const nextConfig: NextConfig = {
  async rewrites() {
    const apiBase = process.env.LOCAL_API_BASE_URL || "http://backend:7772";

    return [
      {
        source: "/api/:path*",
        destination: `${apiBase}/api/:path*`,
      },
    ];
  },
};

const pwaConfig = {
	dest: "public",
	// cacheOnFrontEndNav: true,
	// aggressiveFrontEndNavCaching: true,
	reloadOnOnline: true,
	workboxOptions: {
		disableDevLogs: true,
		runtimeCaching: [
			// {
			// 	// Cache /api/me endpoint
			// 	urlPattern: /^https?:\/\/.*\/api\/me$/,
			// 	handler: 'NetworkFirst' as const,
			// 	options: {
			// 		cacheName: 'api-auth-cache',
			// 		networkTimeoutSeconds: 10,
			// 		expiration: {
			// 			maxEntries: 1,
			// 			maxAgeSeconds: 24 * 60 * 60,
			// 		},
			// 	},
			// },
			// Cache all GET requests starting with /api/
			{
			urlPattern: ({ url, request }: { url: URL; request: Request }) =>
				url.pathname.startsWith('/api/') && request.method === 'GET',
			handler: 'NetworkFirst' as const,
			options: {
				cacheName: 'api-get-cache',
				networkTimeoutSeconds: 10,
				expiration: {
				maxEntries: 50,
				maxAgeSeconds: 5 * 60,
				},
				plugins: [
				{
					handlerDidError: async () => {
					// Return a json response when offline and no cache available
					return new Response(
						JSON.stringify({ error: 'You are offline. Data unavailable.' }),
						{
						status: 503,
						headers: { 'Content-Type': 'application/json' },
						}
					);
					},
				},
				],
			},
			},
			// Cache all modify requests（POST、PUT、DELETE）
			{
				urlPattern: ({ url, request }: { url: URL; request: Request }) =>
				url.pathname.startsWith('/api/') && request.method !== 'GET',
				handler: 'NetworkOnly' as const,    // No cache used
				options: {
				plugins: [
					{
					handlerDidError: async () => {
						return new Response(
						JSON.stringify({ error: 'You are offline. Changes cannot be saved.' }),
						{
							status: 503,
							headers: { 'Content-Type': 'application/json' },
						}
						);
					},
					},
				],
				},
			},
			// Silence _rsc: return empty response when offline
			{
			urlPattern: /\/.*\?_rsc=/,
			handler: 'NetworkFirst' as const,
			options: {
				cacheName: 'rsc-silence',
				networkTimeoutSeconds: 0,           // 离线立即触发错误
				expiration: { maxEntries: 5, maxAgeSeconds: 5 },
				plugins: [
				{
					// 禁止缓存这个假响应，避免在线时使用
					cacheWillUpdate: async () => null,
				},
				{
					handlerDidError: async () =>
					new Response(
						'{"outputs":[],"suspense":{}}\n',   // ✅ 完整的空 RSC 行
						{
						status: 200,
						headers: {
							'Content-Type': 'text/x-component; charset=utf-8',
							'Vary': 'RSC, Next-Router-State-Tree, Next-Router-Prefetch',
						},
						}
					),
				},
				],
			},
			},
			// Cache HTML in page navigation（requests without _rsc ）
			{
			urlPattern: ({ url, request }: { url: URL; request: Request }) => {
				return /^\/(home|events|about|settings|privacy|terms)?$/.test(url.pathname) &&
					request.headers.get('Accept')?.includes('text/html');
			},
			handler: 'NetworkFirst' as const,
			options: {
				cacheName: 'pages-html-cache',
				networkTimeoutSeconds: 10,
				expiration: { maxEntries: 20, maxAgeSeconds: 60 * 60 * 24 },
				plugins: [
					{
						handlerDidError: async () =>
							// Return offline html page
							new Response(
							`<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Offline</title>
			<style>
				body { font-family: system-ui; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; background: #f9fafb; }
				.box { text-align: center; padding: 2rem; background: white; border-radius: 0.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
			</style>
			</head>
			<body>
			<div class="box">
				<h1>📡 You are offline</h1>
				<p>This page has not been cached. Taking you to the homepage in <span id="countdown">5</span> seconds...</p>
			</div>
			<script>
				(function() {
					var seconds = 5;
					var el = document.getElementById('countdown');
					if (!el) return;
					var timer = setInterval(function() {
						seconds--;
						if (seconds <= 0) {
							clearInterval(timer);
							window.location.replace('/home');
						} else {
							el.textContent = seconds;
						}
					}, 1000);
				})();
			</script>
			</body>
			</html>`,
							{ headers: { 'Content-Type': 'text/html' } }
							),
					},
				],
			},
			},
			// Cache static resources (fonts, images)
			{
				urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp|woff|woff2|ttf)$/,
				handler: 'CacheFirst' as const,
				options: {
					cacheName: 'static-assets-cache' as const,
					expiration: {
						maxEntries: 100,
						maxAgeSeconds: 30 *24 * 60 * 60,
					},
				},
			},
			// Silence favicon：return 204 when offline
			{
				urlPattern: /\/favicon\.ico.*/,
				handler: 'NetworkFirst' as const,
				options: {
				cacheName: 'favicon-silence',
				networkTimeoutSeconds: 0,
				expiration: {
					maxEntries: 1,
					maxAgeSeconds: 5,
				},
				plugins: [
					{
					handlerDidError: async () => {
						return new Response(null, {
						status: 204,
						statusText: 'No Content (offline)',
						});
					},
					},
				],
				},
			},
			// {
			// urlPattern: /^\/api\/.*/,
			// handler: 'NetworkFirst' as const,
			// options: {
			// 	cacheName: 'api-cache',
			// 	networkTimeoutSeconds: 10,
			// 	expiration: {
			// 	maxEntries: 50,
			// 	maxAgeSeconds: 5 * 60,
			// 	},
			// 	plugins: [
			// 	{
			// 		handlerDidError: async () => {
			// 		return new Response(JSON.stringify({ error: 'Offline - data unavailable' }), {
			// 			status: 503,
			// 			headers: { 'Content-Type': 'application/json' },
			// 		});
			// 		},
			// 	},
			// 	],
			// },
			// },
		]
	},
};

const config = process.env.NODE_ENV === 'production'
  ? withPWA({ ...nextConfig, ...pwaConfig })
  : nextConfig;

export default config;
