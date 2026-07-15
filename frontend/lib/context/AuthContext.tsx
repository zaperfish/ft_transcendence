'use client';

import { createContext, useState, useEffect, ReactNode, useCallback } from "react";
import type { User } from '@/types/user';
import { login as apiLogin, logout as apiLogout } from '@/lib/api/auth';
import { ApiError } from '@/lib/api/client';
import { getCurrentUser } from "../api/user";
import { usePathname, useRouter } from 'next/navigation';

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
	const pathname = usePathname();
	const isAuthPage = pathname === '/login' || pathname === '/register';
	const isPublicPage = pathname === '/privacy' || pathname === '/terms';
	// Get initial network state (considering both client side and server side)
	const [isOnline, setIsOnline] = useState(
		typeof window !== 'undefined' ? navigator.onLine : true
	);

	// navigator.onLine can be unreliable, using ping to check network
	const checkRealOnline = useCallback(async () => {
		if (!navigator.onLine) {
			setIsOnline(false);
			return false;
		}

		try {
			const controller = new AbortController();
			// Abort signal if not responded in 3 seconds
			const timeoutId = setTimeout(() => controller.abort(), 3000);
			const response = await fetch('/health', {
				cache: 'no-store',
				signal: controller.signal,
			});
			clearTimeout(timeoutId);
			const reallyOnline = response.ok && response.headers.get('X-Network-Status') !== 'offline';
			setIsOnline(reallyOnline);
			return reallyOnline;
		} catch {
			setIsOnline(false);
			return false;
		}
	}, []);

	// Add tracker of online && offline states
	useEffect(() => {
		const handleOnline = async () => {
			await checkRealOnline();
		};
		const handleOffline = () => {
			setIsOnline(false);
		}
		window.addEventListener('online', handleOnline);
		window.addEventListener('offline', handleOffline);
		return () => {
			window.removeEventListener('online', handleOnline);
			window.removeEventListener('offline', handleOffline);
		};
	}, [checkRealOnline]);

	// Save User to localStorage
	const saveAuthToCache = (user: User) => {
		try {
			localStorage.setItem('auth_cache', JSON.stringify(user))
		} catch {
			console.log('Failed to cache auth data');
		}
	};

	// Get User from localStorage
	const loadAuthFromCache = () : User | null => {
		try {
			const cached = localStorage.getItem('auth_cache');
			return cached ? JSON.parse(cached) : null;
		} catch {
			console.log('Failed to load cached auth data');
			return null;
		}
	};

	// Clear User cache
	const clearAuthCache = () => {
		try {
			localStorage.removeItem('auth_cache');
		} catch {
			console.log('Failed to clear cached auth data');
		}
	};

	useEffect(() => {
		const initAuth = async () => {
			try {
				// Avoid requesting /api/me in register/login/public page to avoid 401
				if (isAuthPage) {
					const cachedUser = loadAuthFromCache();
					if (cachedUser) {
						setUser(cachedUser);
					} else {
						setUser(null);
					}
				} else if (isPublicPage) {
					const cachedUser = loadAuthFromCache();
					if (cachedUser) {
						setUser(cachedUser);
					} else {
						setUser(null);
					}
				} else {
					const reallyOnline = await checkRealOnline();
					// If online, get user from server and save to cache
					if (reallyOnline) {
						const currentUser = await getCurrentUser();
						if (currentUser) {
							setUser(currentUser);
							saveAuthToCache(currentUser);
						} else {
							clearAuthCache();
							setUser(null);
						}
					} else {
						// If offline, get user from localStorage
						const cachedUser = loadAuthFromCache();
						if (cachedUser) {
							setUser(cachedUser);
						} else {
							setUser(null);
						}
					}
				}
			} catch (error) {
				// Avoid triggering Next.js error overlay with console.error
				if (error instanceof ApiError && error.status === 0) {
					const cachedUser = loadAuthFromCache();
					if (cachedUser) {
						setUser(cachedUser);
					} else {
						setUser(null);
					}
				} else {
					clearAuthCache();
					setUser(null);
				}
			}
			setIsLoading(false);
		};
		initAuth();
	}, [checkRealOnline, isAuthPage, isPublicPage]);

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

