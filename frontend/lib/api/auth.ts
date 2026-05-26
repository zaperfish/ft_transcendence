import type { User } from '@/types/user';
import { request } from '@/lib/api/client';

// Frontend -> Backend
export interface LoginCredentials {
	name: string;
	password: string;
}

// // Backend -> Frontend
// // token sent to browser with 'Set-Cookie'
// export interface AuthResponse {
// 	user: User;
// }

export interface RegisterCredentials {
	name: string;
	email: string;
	password: string;
}

// Wrap request and response into 1 Login api
export async function login(credentials: LoginCredentials): Promise<User> {
	const res = await request<User>('/api/auth/login', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json'},
		body: JSON.stringify(credentials),
	});
	return res;
}

// Backend in responsible to clean cookie
export async function logout(): Promise<void> {
	await request('/api/auth/logout', {
		method: 'POST',
	});
}

// Register api
export async function register(credentials: RegisterCredentials): Promise<User> {
	return request<User>('/api/auth/register', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json'},
		body: JSON.stringify(credentials),
	});
}