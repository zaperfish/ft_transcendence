import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    const apiBase = process.env.LOCAL_API_BASE_URL;

    return [
      {
        source: "/api/:path*",
        destination: `${apiBase}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
