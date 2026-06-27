'use client';

import { useAuth } from "@/lib/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import Navigation from "@/components/layout/Navigation";

export default function ChatProtectedLayout({ children }: Readonly<{
	children: React.ReactNode;
}>) {
	const { isAuthenticated, isLoading } = useAuth();
	const router = useRouter();

	useEffect(() => {
		if (!isLoading && !isAuthenticated) {
			router.push('/login');
		}
	}, [isLoading, isAuthenticated, router]);

	if (isLoading) {
		return <div>Loading...</div>;
	}

	return (
		<div className='flex h-full min-h-0 w-full flex-col overflow-hidden bg-surface-dim'>
			<Navigation />
			<main className='flex min-h-0 flex-1 flex-col overflow-hidden'>
				{children}
			</main>
		</div>
	);
}
