'use client';

import { Moon, Sun } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { useTheme } from '@/lib/context/ThemeContext';
import { cn } from '@/lib/utils';

/**
 * Toggles between aurora (wallpaper / teal chrome) and classic (light) themes.
 * Placed in the footer so it is available on login and all other pages.
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
			className={cn(
				'h-8 w-8',
				isAurora
					? 'text-chrome-footer hover:bg-white/10 hover:text-chrome-title'
					: 'text-muted-foreground hover:text-foreground',
			)}
		>
			{isAurora ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
		</Button>
	);
}
