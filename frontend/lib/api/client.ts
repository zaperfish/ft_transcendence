// Wrap fetch() with error handling (status, message)

export class ApiError extends Error {
	constructor(public status: number, message: string) {
		super(message);
	}
}

export async function request<T>(url: string, options?: RequestInit): Promise<T> {
	const response = await fetch(url, options);
	if (!response.ok) {
		throw new ApiError(response.status, 'Request failed');
	}
	return response.json();
}