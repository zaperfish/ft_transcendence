import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "@/components/providers/Providers";
import { OfflineBanner } from "@/components/ui/OfflineBanner";
import { Footer } from "@/components/layout/Footer";

// font-sans is a Tailwind CSS utility class
// that applies the CSS rule font-family:
// var(--font-sans)
const inter = Inter({
  variable: "--font-sans",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Camaraderie",
  description: "A meetup app that connects people around shared purpose",
  manifest: '/manifest.json',
};

/**
 * The root layout of the entire application.
 *
 * Sets up the base HTML structure, applies the Inter font variable,
 * enables font smoothing, and wraps all pages inside the
 * {@link providers} so that authentication state and query cache is available
 * throughout the app.
 *
 * @param props - The component props.
 * @param props.children - The page or nested layout content to render.
 * @returns The top-level `<html>` and `<body>` elements with global providers.
 */
export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={inter.variable}>
      <body className={"h-dvh overflow-hidden font-sans antialiased"}>
        <Providers>
          <div className="flex h-dvh flex-col overflow-hidden">
            <OfflineBanner />
            <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
              {children}
            </div>
			<Footer />
          </div>
        </Providers>
      </body>
    </html>
  );
}
