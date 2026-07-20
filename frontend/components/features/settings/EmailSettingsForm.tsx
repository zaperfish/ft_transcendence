'use client';

import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { FormLabel } from '@/components/ui/FormLabel';
import { SettingsPanel } from '@/components/features/settings/SettingsPanel';
import { updateProfile } from '@/lib/api/user';
import { useAuth } from '@/lib/hooks/useAuth';
import { ApiError } from '@/lib/api/client';

const emailSchema = z.object({
	email: z
		.string()
		.min(1, 'Please enter your email address')
		.min(5, 'Email address should be at least 5 characters long')
		.max(64, 'Email address should be no more than 64 characters long')
		.regex(
			/^(?!.*\.\.)([A-Za-z0-9_%+-]+(?:\.[A-Za-z0-9_%+-]+)*)@(?:[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?\.)+[A-Za-z]{2,63}$/,
			'Invalid email address'
		),
});

type EmailFormData = z.infer<typeof emailSchema>;

export function EmailSettingsForm() {
	const { user, refreshUser, isOnline } = useAuth();
	const {
		register,
		handleSubmit,
		reset,
		setError,
		formState: { errors, isSubmitting, isDirty },
	} = useForm<EmailFormData>({
		resolver: zodResolver(emailSchema),
		mode: 'onBlur',
		defaultValues: { email: '' },
	});

	useEffect(() => {
		if (user) {
			reset({ email: user.email });
		}
	}, [user, reset]);

	const onSubmit = async (data: EmailFormData) => {
		if (!isOnline) {
			setError('root', { message: 'You are offline. Changes cannot be saved.' });
			return;
		}

		try {
			await updateProfile({ email: data.email });
			await refreshUser();
		} catch (err) {
			if (err instanceof ApiError && err.status === 409) {
				setError('root', { message: 'Email address is already taken' });
			} else {
				setError('root', { message: 'Failed to update email, please try again' });
			}
		}
	};

	if (!user) {
		return null;
	}

	return (
		<SettingsPanel
			title="Profile"
			description="Your username cannot be changed. Update your email address below."
		>
			<form onSubmit={handleSubmit(onSubmit)} className="space-y-lg">
				<div className="space-y-sm">
					<FormLabel htmlFor="username" required={false}>Username</FormLabel>
					<Input id="username" type="text" autoComplete="username" value={user.name} disabled />
				</div>

				<div className="space-y-sm">
					<FormLabel htmlFor="email">Email</FormLabel>
					<Input
						id="email"
						type="email"
						placeholder="Email"
						autoComplete="email"
						{...register('email')}
						className={errors.email ? 'border-error' : ''}
					/>
					{errors.email && (
						<p className="text-sm text-error">{errors.email.message}</p>
					)}
				</div>

				{errors.root && (
					<p className="text-sm text-error">{errors.root.message}</p>
				)}

				<Button
					type="submit"
					disabled={isSubmitting || !isDirty}
					className={!isOnline ? 'opacity-50 cursor-not-allowed' : ''}
				>
					{isSubmitting ? 'Saving...' : 'Save changes'}
				</Button>
			</form>
		</SettingsPanel>
	);
}
