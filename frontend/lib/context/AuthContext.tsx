'use client';

import { createContext, useState, useEffect, ReactNode } from "react";
import type { User } from '@/types/user';
import { login as apiLogin, logout as apiLogout } from '@/lib/api/auth';
import { getCurrentUser } from "../api/user";
import { useRouter } from 'next/navigation';

// Define Auth context interface
interface AuthContextType {
	user: User | null;
	isAuthenticated: boolean;
	isLoading: boolean;
	login: (credentials: { username: string; password: string }) => Promise<void>;
	logout: () => Promise<void>;
}

// Create Auth context with default value and method
export const AuthContext = createContext<AuthContextType>({
	user: null,
	isAuthenticated: false,
	isLoading: true,
	login: async () => {},
	logout: async () => {},
});

// Initiate Auth context with real data from api
// And keep it updatable with methods
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
				console.error('Failed login automatically', error);
				setUser(null);
			}
			setIsLoading(false);
		};
		initAuth();
	}, []);

	const login = async (credentials: { username: string; password: string }) => {
		const { user } = await apiLogin(credentials);
		setUser(user);
	};

	const logout = async () => {
		try {
			await apiLogout();
		} catch (error) {
			console.error("Failed to logout", error);
		}
		setUser(null);
		// Soft redirect to '/login'
		router.push('/login');
	};

	return (
		<AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, login, logout }}>
			{children}
		</AuthContext.Provider>
	);
}

