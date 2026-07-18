'use client';

import { useAuth } from "@/lib/hooks/useAuth";

export function OfflineBanner() {
	const { isOnline } = useAuth();

	if (isOnline) return null;

	return (
	<div className="bg-yellow-100 border-b border-yellow-200 px-4 py-2 text-sm text-yellow-800">
		You are in offline mode, using cached data.
	</div>
	);
}