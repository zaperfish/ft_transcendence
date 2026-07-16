'use client';

import { EmailSettingsForm } from '@/components/features/settings/EmailSettingsForm';
import { ChangePasswordForm } from '@/components/features/settings/ChangePasswordForm';
import { DeleteAccountSection } from '@/components/features/settings/DeleteAccountSection';

export function SettingsContent() {
	return (
		<div className="flex w-full flex-col gap-2xl">
			<EmailSettingsForm />
			<ChangePasswordForm />
			<DeleteAccountSection />
		</div>
	);
}
