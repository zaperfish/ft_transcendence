"use client";

import { AuthProvider } from "@/lib/context/AuthContext";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

// Avoid instantiate whenever children are rerendered and lose cache
const queryClient = new QueryClient();

/**
 * Providers wraps the application with necessary context providers:
 * - QueryClientProvider for React Query (data fetching and caching)
 * - AuthProvider for authentication state management
 * - ReactQueryDevtools for development debugging (only shown in dev mode)
 */
export default function Providers({ children }: { children: React.ReactNode }) {
	return (
		<QueryClientProvider client={queryClient}>
			<AuthProvider>
				{children}
			</AuthProvider>
			<ReactQueryDevtools initialIsOpen={false} />
		</QueryClientProvider>
	);
}