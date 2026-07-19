export const THEME_STORAGE_KEY = 'camaraderie-theme';

export const THEMES = ['aurora', 'classic'] as const;

export type Theme = (typeof THEMES)[number];

export function isTheme(value: unknown): value is Theme {
	return value === 'aurora' || value === 'classic';
}

export function applyTheme(theme: Theme) {
	document.documentElement.setAttribute('data-theme', theme);
}

export function readStoredTheme(): Theme | null {
	try {
		const value = localStorage.getItem(THEME_STORAGE_KEY);
		return isTheme(value) ? value : null;
	} catch {
		return null;
	}
}

export function storeTheme(theme: Theme) {
	try {
		localStorage.setItem(THEME_STORAGE_KEY, theme);
	} catch {
		// ignore quota / private mode errors
	}
}
