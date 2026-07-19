import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "@/components/providers/Providers";
import { OfflineBanner } from "@/components/ui/OfflineBanner";
import { Footer } from "@/components/layout/Footer";
import { Toaster } from 'sonner';

// Shared by both classic and aurora themes
const inter = Inter({
  variable: "--font-inter",
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
const themeInitScript = `
(function () {
  try {
    var stored = localStorage.getItem('camaraderie-theme');
    var theme = stored === 'classic' || stored === 'aurora' ? stored : 'aurora';
    document.documentElement.setAttribute('data-theme', theme);
  } catch (e) {
    document.documentElement.setAttribute('data-theme', 'aurora');
  }
})();
`;

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={inter.variable} suppressHydrationWarning>
      <body className={"h-dvh overflow-hidden font-sans antialiased"}>
        <script dangerouslySetInnerHTML={{ __html: themeInitScript }} />
        <Providers>
          <div className="flex h-dvh flex-col overflow-hidden">
            <OfflineBanner />
            <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
              {children}
            </div>
			<Footer />
          </div>
        </Providers>
		<Toaster position="top-center" richColors closeButton />
      </body>
    </html>
  );
}
