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
	isOnline: boolean,
	login: (credentials: { name: string; password: string }) => Promise<void>;
	logout: () => Promise<void>;
}

/**
 * Create Auth context with default value and method
 */
export const AuthContext = createContext<AuthContextType>({
	user: null,
	isAuthenticated: false,
	isLoading: true,
	isOnline: false,
	login: async () => {},
	logout: async () => {},
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
	// Get initial network state (considering both client side and server side)
	const [isOnline, setIsOnline] = useState(typeof window !== 'undefined' ? navigator.onLine : true);

	// Add tracker of online && offline states
	useEffect(() => {
		const handleOnline = () => setIsOnline(true);
		const handleOffline = () => setIsOnline(false);
		window.addEventListener('online', handleOnline);
		window.addEventListener('offline', handleOffline);
		return () => {
			window.removeEventListener('online', handleOnline);
			window.removeEventListener('offline', handleOffline);
		};
	}, []);

	// Save User to localStorage
	const saveAuthToCache = (user: User) => {
		try {
			localStorage.setItem('auth_cache', JSON.stringify(user))
		} catch (error) {
			console.log('Failed to cache auth data');
		}
	};

	// Get User from localStorage
	const loadAuthFromCache = () : User | null => {
		try {
			const cached = localStorage.getItem('auth_cache');
			return cached ? JSON.parse(cached) : null;
		} catch (error) {
			console.log('Failed to load cached auth data');
			return null;
		}
	};

	// Clear User cache
	const clearAuthCache = () => {
		try {
			localStorage.removeItem('auth_cache');
		} catch (error) {
			console.log('Failed to clear cached auth data');
		}
	};

	useEffect(() => {
		const initAuth = async () => {
			try {
				// If online, get user from server and save to cache
				if (isOnline) {
					const currentUser = await getCurrentUser();
					setUser(currentUser);
					saveAuthToCache(currentUser);
				} else {
				// If offline, get user from localStorage
					const cachedUser = loadAuthFromCache();
					if (cachedUser) {
						setUser(cachedUser);
						console.log('Loaded user from cache(offline)');
					} else {
						console.log('No cached user data found');
						setUser(null);
					}
				}
			} catch (error) {
				// Avoid triggering Next.js error overlay with console.error
				console.log('User not logged in or session expired');
				const cachedUser = loadAuthFromCache();
				if (cachedUser) {
					setUser(cachedUser);
					console.log('Fallback to cached user data');
				} else {
					setUser(null);
				}
			}
			setIsLoading(false);
		};
		initAuth();
	}, [isOnline]);

	const login = async (credentials: { name: string; password: string }) => {
		const user = await apiLogin(credentials);
		setUser(user);
		saveAuthToCache(user);
	};

	const logout = async () => {
		try {
			await apiLogout();
		} catch (error) {
			console.error("Failed to logout", error);
		}
		setUser(null);
		clearAuthCache();
		router.push('/login');
	};

	return (
		<AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, isOnline, login, logout }}>
			{children}
		</AuthContext.Provider>
	);
}

