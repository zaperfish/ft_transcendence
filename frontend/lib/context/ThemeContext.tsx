'use client';

import {
	createContext,
	useCallback,
	useContext,
	useEffect,
	useState,
	type ReactNode,
} from 'react';
import {
	THEME_STORAGE_KEY,
	applyTheme,
	isTheme,
	readStoredTheme,
	storeTheme,
	type Theme,
} from '@/lib/theme';

interface ThemeContextType {
	theme: Theme;
	setTheme: (theme: Theme) => void;
	toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType>({
	theme: 'aurora',
	setTheme: () => {},
	toggleTheme: () => {},
});

export function ThemeProvider({ children }: { children: ReactNode }) {
	const [theme, setThemeState] = useState<Theme>('aurora');

	useEffect(() => {
		const stored = readStoredTheme();
		const initial = stored ?? 'aurora';
		setThemeState(initial);
		applyTheme(initial);
	}, []);

	const setTheme = useCallback((next: Theme) => {
		setThemeState(next);
		applyTheme(next);
		storeTheme(next);
	}, []);

	const toggleTheme = useCallback(() => {
		setThemeState((current) => {
			const next: Theme = current === 'aurora' ? 'classic' : 'aurora';
			applyTheme(next);
			storeTheme(next);
			return next;
		});
	}, []);

	useEffect(() => {
		const onStorage = (event: StorageEvent) => {
			if (event.key !== THEME_STORAGE_KEY) return;
			if (isTheme(event.newValue)) {
				setThemeState(event.newValue);
				applyTheme(event.newValue);
			}
		};
		window.addEventListener('storage', onStorage);
		return () => window.removeEventListener('storage', onStorage);
	}, []);

	return (
		<ThemeContext.Provider value={{ theme, setTheme, toggleTheme }}>
			{children}
		</ThemeContext.Provider>
	);
}

export function useTheme() {
	return useContext(ThemeContext);
}
