import type { User } from '@/types/user';
import { request } from '@/lib/api/client';

// Frontend -> Backend
export interface LoginCredentials {
	username: string;
	password: string;
}

// Backend -> Frontend
export interface AuthResponse {
	token: string;
	user: User;
}

export async function login(credentials: LoginCredentials): Promise<AuthResponse> {
	const res = await request<AuthResponse>('/api/user/login', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json'},
		body: JSON.stringify(credentials),
	});
	return res;
}