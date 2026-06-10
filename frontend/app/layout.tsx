import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "@/components/providers/Providers";

const inter = Inter({
  variable: "--font-inter",// CSS variable
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
			{children}
		</Providers>
	  </body>
    </html>
  );
}
