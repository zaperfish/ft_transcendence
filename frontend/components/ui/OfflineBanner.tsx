'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';

/**
 * Offline banner. Intentionally waits until after client mount so SSR HTML
 * always matches the first client paint (avoids hydration mismatch when
 * navigator.onLine differs from the server default).
 */
export function OfflineBanner() {
	const { isOnline } = useAuth();
	const [mounted, setMounted] = useState(false);

	useEffect(() => {
		setMounted(true);
	}, []);

	if (!mounted || isOnline) return null;

	return (
		<div className="bg-yellow-100 border-b border-yellow-200 px-4 py-2 text-sm text-yellow-800">
			You are in offline mode, using cached data.
		</div>
	);
}
