import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({
  variable: "--font-inter",// CSS variable
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Meetup App",
  description: "A modern meetup platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.variable} font-sans antialiased`}>
		{children}
	  </body>
    </html>
  );
}
