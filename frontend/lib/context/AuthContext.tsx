'use client';

import { createContext, useState, useEffect, ReactNode } from "react";
import type { User } from '@/types/user';
import { login as apiLogin, logout as apiLogout } from '@/lib/api/auth';
import { getCurrentUser } from "../api/user";
import { useRouter } from 'next/navigation';

/**
 * Define Auth context interface
 */
interface AuthContextType {
	user: User | null;
	isAuthenticated: boolean;
	isLoading: boolean;
	login: (credentials: { name: string; password: string }) => Promise<void>;
	logout: () => Promise<void>;
	refreshUser: () => Promise<void>;
}

/**
 * Create Auth context with default value and method
 */
export const AuthContext = createContext<AuthContextType>({
	user: null,
	isAuthenticated: false,
	isLoading: true,
	login: async () => {},
	logout: async () => {},
	refreshUser: async () => {},
});

/**
 * Initializes auth state by fetching real user data from the API,
 * and exposes methods (login, logout, etc.) to keep it updatable.
 * @param children - React subtree that will have access to the auth context.
 * @returns A React element that provides authentication context to its children.
 */
export function AuthProvider({ children } : { children: ReactNode }) {
	const [user, setUser] = useState<User | null>(null);
	const [isLoading, setIsLoading] = useState(true);
	const router = useRouter();

	useEffect(() => {
		const initAuth = async () => {
			try {
				const currentUser = await getCurrentUser();
				setUser(currentUser);
			} catch (error) {
				// Avoid triggering Next.js error overlay with console.error
				console.log('User not logged in or session expired');
				setUser(null);
			}
			setIsLoading(false);
		};
		initAuth();
	}, []);

	const login = async (credentials: { name: string; password: string }) => {
		const user = await apiLogin(credentials);
		setUser(user);
	};

	const logout = async () => {
		try {
			await apiLogout();
		} catch (error) {
			console.error("Failed to logout", error);
		}
		setUser(null);
		router.push('/login');
	};

	const refreshUser = async () => {
		try {
			const currentUser = await getCurrentUser();
			setUser(currentUser);
		} catch {
			setUser(null);
		}
	};

	return (
		<AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, login, logout, refreshUser }}>
			{children}
		</AuthContext.Provider>
	);
}

