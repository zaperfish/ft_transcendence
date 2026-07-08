import type { Metadata } from "next";
import Link from "next/link";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "@/components/providers/Providers";
import { OfflineBanner } from "@/components/ui/OfflineBanner";

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
          <OfflineBanner />
          <div className="flex h-dvh flex-col overflow-hidden">
            <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
              {children}
            </div>
            <footer className="shrink-0 border-t border-border bg-background px-4 py-4 text-sm text-muted-foreground">
              <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-center">
                <Link href="/privacy" className="hover:text-foreground">
                  Privacy Policy
                </Link>
                <span className="hidden sm:inline">•</span>
                <Link href="/terms" className="hover:text-foreground">
                  Terms of Service
                </Link>
              </div>
            </footer>
          </div>
        </Providers>
      </body>
    </html>
  );
}
