import type { User } from '@/types/user';

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