'use client';

import { Moon, Sun } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { useTheme } from '@/lib/context/ThemeContext';

/**
 * Toggles between aurora (wallpaper / teal chrome) and classic (light) themes.
 */
export function ThemeToggle() {
	const { theme, toggleTheme } = useTheme();
	const isAurora = theme === 'aurora';

	return (
		<Button
			type="button"
			variant="ghost"
			size="icon"
			onClick={toggleTheme}
			aria-label={isAurora ? 'Switch to classic theme' : 'Switch to aurora theme'}
			title={isAurora ? 'Classic theme' : 'Aurora theme'}
			className="text-chrome-nav hover:bg-white/10 hover:text-chrome-title"
		>
			{isAurora ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
		</Button>
	);
}
