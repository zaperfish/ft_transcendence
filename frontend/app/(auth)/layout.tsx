'use client';

import { useAuth } from '@/lib/hooks/useAuth';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

/**
 * Layout wrapper for authentication pages (login, register, etc.).
 *
 * Centers its content in a card-like container and manages auth state:
 * - Shows a loading indicator while the auth status is being resolved.
 * - Automatically redirects to `/home` if the user is already authenticated.
 *
 * @param props - The component props.
 * @param props.children - The auth page content to be rendered inside the centered card.
 * @returns A centered layout shell with auth-aware routing.
 */
export default function AuthLayout({ children }: { children: React.ReactNode } ) {
	const { isAuthenticated, isLoading } = useAuth();
	const router = useRouter();

	useEffect(() => {
		if (!isLoading && isAuthenticated) {
			router.push('/home');
		}
	}, [isAuthenticated, isLoading, router]);
	if (isLoading)
		return <div>Loading...</div>
	return (
		<div className="flex min-h-screen items-center justify-center p-4">
			<div className="w-full max-w-[576px]">
				{children}
			</div>
		</div>
	);
}