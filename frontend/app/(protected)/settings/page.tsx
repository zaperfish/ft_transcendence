'use client';

import { SettingsContent } from '@/components/features/settings/settings';
import { useAuth } from '@/lib/hooks/useAuth';

export default function SettingsPage() {
  const { isLoading } = useAuth();
  if (isLoading) return <div className="text-center py-2xl text-text-secondary">Loading...</div>;

  return (
    <div className="w-full px-xl py-2xl">
      <div className="max-w-175 w-full mx-auto">
        <h1 className="text-4xl font-heading font-bold text-text-primary mb-md">Settings</h1>
        <SettingsContent />
      </div>
    </div>
  );
}