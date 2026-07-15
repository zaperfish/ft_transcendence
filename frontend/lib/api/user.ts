import { ApiError, request } from '@/lib/api/client';
import type { User } from '@/types/user';

/**
 * Restore User information and authtication for refresh and protected pages
 * @returns Promise of User or null when the session is invalid
 */
export async function getCurrentUser(): Promise<User | null> {
	try {
		return await request<User>('/api/me');
	} catch (error) {
		if (error instanceof ApiError && error.status === 401) {
			return null;
		}
		throw error;
	}
}