import type { Metadata } from "next";
import { DM_Sans, Space_Grotesk } from "next/font/google";
import "./globals.css";
import Providers from "@/components/providers/Providers";
import { OfflineBanner } from "@/components/ui/OfflineBanner";
import { Footer } from "@/components/layout/Footer";
import { Toaster } from 'sonner';

const dmSans = DM_Sans({
  variable: "--font-dm-sans",
  subsets: ["latin"],
});

const spaceGrotesk = Space_Grotesk({
  variable: "--font-space-grotesk",
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
    <html lang="en" className={`${dmSans.variable} ${spaceGrotesk.variable}`}>
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
		<Toaster position="top-center" richColors closeButton />
      </body>
    </html>
  );
}
