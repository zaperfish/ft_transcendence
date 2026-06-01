'use client';

import { useAuth } from "@/lib/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import Navigation from "@/components/layout/Navigation";

/**
 * Layout wrapper for protected (authenticated) pages.
 *
 * Manages authentication state:
 * - Displays a loading indicator while auth status is resolving.
 * - Redirects to `/login` if the user is not authenticated.
 *
 * When authenticated, renders the child content.
 *
 * @param props - The component props.
 * @param props.children - The protected page content to render after authentication is confirmed.
 * @returns A layout with navigation and children, or a loading/redirect state.
 */
export default function ProtectedLayout({  children }: Readonly<{
  children: React.ReactNode;
}>) {
	const { isAuthenticated, isLoading } = useAuth();
	const router = useRouter();

	useEffect(() => {
		if (!isLoading && !isAuthenticated) {
			router.push('/login');
		}
	}, [isLoading, isAuthenticated, router]);
	if (isLoading)
		return <div>Loading...</div>
	return (
    	<div className='min-h-screen bg-surface-dim flex flex-col'>
			<Navigation />
			<main className='flex-1 w-full max-w-300 mx-auto px-md py-xl'>
      			{children}
    		</main>
		</div>
  );
}