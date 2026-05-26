'use client';

import { useAuth } from "@/lib/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

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
    <>
      <nav>
        {/* Navigation placeholder */}
      </nav>
      {children}
    </>
  );
}