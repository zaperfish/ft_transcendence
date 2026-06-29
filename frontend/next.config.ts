import type { NextConfig } from "next";
import withPWA from "@ducanh2912/next-pwa";

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
	cacheOnFrontEndNav: true,
	aggressiveFrontEndNavCaching: true,
	reloadOnOnline: true,
	workboxOptions: {
		disableDevLogs: true,
	},
};

const config = process.env.NODE_ENV === 'production'
  ? withPWA({ ...nextConfig, ...pwaConfig })
  : nextConfig;

export default config;
