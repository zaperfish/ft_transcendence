import { request } from '@/lib/api/client';
import type { UpdateProfileRequest, User, UpdatePasswordRequest } from '@/types/user';

/**
 * Restore User information and authtication for refresh and protected pages
 * @returns Promise of User
 */
export async function getCurrentUser(): Promise<User> {
	return request<User>('/api/me');
}

export async function updateProfile(data: UpdateProfileRequest): Promise<User> {
  return request<User>('/api/me', {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
}

export async function updatePassword(data: UpdatePasswordRequest): Promise<User> {
	return request<User>('/api/me/password', {
	  method: 'PATCH',
	  headers: { 'Content-Type': 'application/json' },
	  body: JSON.stringify(data),
	});
  }

export async function deleteAccount(): Promise<void> {
  const response = await fetch('/api/me', {
    method: 'DELETE',
    credentials: 'include',
  });
  if (!response.ok) {
    throw new Error('Failed to delete account');
  }
}