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
	return<>{children}</>;
}