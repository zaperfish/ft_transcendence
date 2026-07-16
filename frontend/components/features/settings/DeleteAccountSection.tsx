'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { SettingsPanel } from '@/components/features/settings/SettingsPanel';
import { deleteAccount } from '@/lib/api/user';
import { useAuth } from '@/lib/hooks/useAuth';

export function DeleteAccountSection() {
	const { logout } = useAuth();
	const [isDeleting, setIsDeleting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleDelete = async () => {
		if (!confirm('This will permanently delete your account. Continue?')) {
			return;
		}

		setIsDeleting(true);
		setError(null);

		try {
			await deleteAccount();
			await logout();
		} catch {
			setError('Failed to delete account, please try again');
			setIsDeleting(false);
		}
	};

	return (
		<SettingsPanel
			title="Danger zone"
			description="Permanently delete your account and all associated data."
			variant="danger"
		>
			{error && <p className="text-sm text-error mb-md">{error}</p>}
			<Button
				type="button"
				variant="destructive"
				disabled={isDeleting}
				onClick={handleDelete}
			>
				{isDeleting ? 'Deleting...' : 'Delete account'}
			</Button>
		</SettingsPanel>
	);
}
