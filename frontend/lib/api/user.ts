import { request } from '@/lib/api/client';
import type { User } from '@/types/user';

/**
 * Restore User information and authtication for refresh and protected pages
 * @returns Promise of User
 */
export async function getCurrentUser(): Promise<User> {
	return request<User>('/api/me');
}