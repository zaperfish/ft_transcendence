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
			{
				// Cache /api/me endpoint
				urlPattern: /^https?:\/\/.*\/api\/me$/,
				handler: 'NetworkFirst' as const,
				options: {
					cacheName: 'api-auth-cache',
					networkTimeoutSeconds: 10,
					expiration: {
						maxEntries: 1,
						maxAgeSeconds: 24 * 60 * 60,
					},
				},
			},
			{
				// Cache other API endpoint (GET)
				urlPattern: /^https?:\/\/.*\/api\/.*/,
				handler: 'NetworkFirst' as const,
				options: {
					cacheName: 'api-cache',
					networkTimeoutSeconds: 10,
					expiration: {
						maxEntries: 50,
						maxAgeSeconds: 5 * 60,
					},
				},
			},
			// {
			// 	urlPattern: ({ url }: { url: URL }) => {
			// 	const path = url.pathname;
			// 	const isPage = /^\/(home|events|about|settings|privacy|terms)$/.test(path);
			// 	return isPage && !url.searchParams.has('_rsc');
			// 	},
			// 	handler: 'NetworkFirst' as const,
			// 	options: {
			// 	cacheName: 'pages-cache',
			// 	networkTimeoutSeconds: 10,
			// 	expiration: {
			// 		maxEntries: 50,
			// 		maxAgeSeconds: 60 * 60,
			// 	},
			// 	},
			// },
			{
				// Cache static resources (fonts, images)
				urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp|woff|woff2|ttf)$/,
				handler: 'CacheFirst' as const,
				options: {
					cacheName: 'static-assets-cache',
					expiration: {
						maxEntries: 100,
						maxAgeSeconds: 30 *24 * 60 * 60,
					},
				},
			},
		]
	},
};

const config = process.env.NODE_ENV === 'production'
  ? withPWA({ ...nextConfig, ...pwaConfig })
  : nextConfig;

export default config;
