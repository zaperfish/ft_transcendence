import { request } from '@/lib/api/client';
import type { User } from '@/types/user';

export async function getCurrentUser(): Promise<User> {
	return request<User>('/api/me');
}