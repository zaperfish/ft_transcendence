'use client';

import { SettingsContent } from '@/components/features/settings/settings';
import { useAuth } from '@/lib/hooks/useAuth';
import { useTheme } from '@/lib/context/ThemeContext';
import { cn } from '@/lib/utils';

export default function SettingsPage() {
  const { isLoading } = useAuth();
  const { theme } = useTheme();
  const isClassic = theme === 'classic';

  if (isLoading) {
    return (
      <div
        className={cn(
          'py-2xl text-center',
          isClassic ? 'text-text-secondary' : 'text-chrome-body',
        )}
      >
        Loading...
      </div>
    );
  }

  return (
    <div className="w-full px-xl py-2xl">
      <div className="max-w-175 w-full mx-auto">
        <h1
          className={cn(
            'mb-md font-heading text-4xl font-bold',
            isClassic ? 'text-text-primary' : 'text-chrome-title',
          )}
        >
          Settings
        </h1>
        <SettingsContent />
      </div>
    </div>
  );
}