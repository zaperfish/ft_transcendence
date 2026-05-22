export class ApiError extends Error {
	constructor(public status: number, message: string) {
		super(message);
	}
}

// Wrap fetch() with error handling (status, message) and auth headers
export async function request<T>(url: string, options?: RequestInit): Promise<T> {
	// Add 'credentials' header to use Set-Cookie automatically
	const response = await fetch(url, {
		...options,
		credentials: 'include',
		headers: {
			...options?.headers,
		},
	});

	if (!response.ok) {
		// Hard redirect to login when 401 Unauthorized
		if (response.status === 401) {
			window.location.href = '/login';
		}
		throw new ApiError(response.status, 'Request failed');
	}
	return response.json();
}