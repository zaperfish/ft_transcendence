// Save, get and remove auth token from browser

const TOKEN_KEY = 'auth_token';

export const saveToken = (token: string) => {
	localStorage.setItem(TOKEN_KEY, token);
};

export const getStoredToken = () : string | null => {
	return localStorage.getItem(TOKEN_KEY);
};

export const removeToken = () => {
	localStorage.removeItem(TOKEN_KEY);
};