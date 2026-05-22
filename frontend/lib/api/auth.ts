import type { User } from '@/types/user';
import { request } from '@/lib/api/client';

// Frontend -> Backend
export interface LoginCredentials {
	username: string;
	password: string;
}

// Backend -> Frontend
// token sent to browser with 'Set-Cookie'
export interface AuthResponse {
	user: User;
}

// Wrap request and response into 1 Login api
export async function login(credentials: LoginCredentials): Promise<AuthResponse> {
	const res = await request<AuthResponse>('/api/user/login', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json'},
		body: JSON.stringify(credentials),
	});
	return res;
}

// Backend in responsible to clean cookie
export async function logout(): Promise<void> {
	await request('/api/user/logout', {
		method: 'POST',
	});
}