import type { Metadata } from "next";
import Link from "next/link";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "@/components/providers/Providers";

const inter = Inter({
  variable: "--font-inter", // CSS variable
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Meetup App",
  description: "A modern meetup platform",
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
    <html lang="en">
      <body className={`${inter.variable} font-sans antialiased`}>
        <Providers>
          <div className="flex min-h-screen flex-col">
            <div className="flex-1">{children}</div>
            <footer className="border-t border-border bg-background px-4 py-4 text-sm text-muted-foreground">
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
