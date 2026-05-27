export class ApiError extends Error {
	constructor(public status: number, message: string) {
		super(message);
	}
}

/**
 * Wrap fetch() with error handling (status, message) and auth headers
 *
 * Add 'credentials' header to use Set-Cookie automatically
 *
 * Hard redirect to login when 401 Unauthorized
 * (before check client-side and not loop redirection in login and register page)
 * 
 * @param url
 * @param options
 * @returns Promise with certain type
 */
export async function request<T>(url: string, options?: RequestInit): Promise<T> {
	const response = await fetch(url, {
		...options,
		credentials: 'include',
		headers: {
			...options?.headers,
		},
	});

	if (!response.ok) {
		if (response.status === 401 && typeof window !== 'undefined' && window.location.pathname !== '/login' && window.location.pathname !== '/register') {
			window.location.href = '/login';
		}
		throw new ApiError(response.status, 'Request failed');
	}
	return response.json();
}