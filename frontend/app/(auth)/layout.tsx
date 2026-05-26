'use client';

import { useAuth } from '@/lib/hooks/useAuth';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

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