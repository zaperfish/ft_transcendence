'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

/**
 * The root landing page.
 *
 * Immediately redirects the user to `/home` using client-side routing.
 * Renders nothing (`null`) while the redirect is in progress.
 *
 * @returns `null` – this component produces no visual output.
 */
export default function RootPage() {
	const router = useRouter();
	useEffect(() => {
		router.replace('/home');
	}, [router]);

	return null;
}