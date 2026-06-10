import type { User } from '@/types/user';
import { request } from '@/lib/api/client';

/**
 * Interface sent Frontend -> Backend in login
 */
export interface LoginCredentials {
	name: string;
	password: string;
}

/**
 * Interface sent Frontend -> Backend in register
 */
export interface RegisterCredentials {
	name: string;
	email: string;
	password: string;
	password_confirm: string;
}

/**
 * login api request
 * @param credentials
 * @returns Promise of User
 */
export async function login(credentials: LoginCredentials): Promise<User> {
	const res = await request<User>('/api/auth/login', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json'},
		body: JSON.stringify(credentials),
	});
	return res;
}

/**
 * logout api request
 *
 * Backend in responsible to clean cookie
 *
 * Not use request() because no need response.json() when response is empty
 */
export async function logout(): Promise<void> {
	await fetch('/api/auth/logout', {
		method: 'POST',
		credentials: 'include',
	});
}

/**
 * register api request
 * @param credentials
 * @returns Promise of User
 */
export async function register(credentials: RegisterCredentials): Promise<User> {
	return request<User>('/api/auth/register', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json'},
		body: JSON.stringify(credentials),
	});
}