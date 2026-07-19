'use client';

import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { FormLabel } from '@/components/ui/FormLabel';
import { SettingsPanel } from '@/components/features/settings/SettingsPanel';
import { updatePassword } from '@/lib/api/user';
import { ApiError } from '@/lib/api/client';
import { useAuth } from '@/lib/hooks/useAuth';

const passwordSchema = z
	.object({
		currentPassword: z.string().min(1, 'Please enter your current password'),
		newPassword: z
			.string()
			.min(1, 'Please enter your new password')
			.min(8, 'Password should be at least 8 characters long')
			.max(128, 'Password should be no more than 128 characters')
			.regex(/[A-Z]/, 'Password should contain at least an uppercase letter')
			.regex(/[a-z]/, 'Password should contain at least an lowercase letter')
			.regex(/\d/, 'Password should contain at least a number'),
		confirmPassword: z.string().min(1, 'Please confirm your new password'),
	})
	.refine((data) => data.newPassword === data.confirmPassword, {
		message: 'Passwords do not match',
		path: ['confirmPassword'],
	});

type PasswordFormData = z.infer<typeof passwordSchema>;

export function ChangePasswordForm() {
	const { isOnline } = useAuth();
	const {
		register,
		handleSubmit,
		reset,
		setError,
		formState: { errors, isSubmitting },
	} = useForm<PasswordFormData>({
		resolver: zodResolver(passwordSchema),
		mode: 'onBlur',
		defaultValues: {
			currentPassword: '',
			newPassword: '',
			confirmPassword: '',
		},
	});

	const onSubmit = async (data: PasswordFormData) => {
		if (!isOnline) {
			setError('root', { message: 'You are offline. Password cannot be changed.' });
			return;
		}

		try {
			await updatePassword({
				current_password: data.currentPassword,
				newpassword: data.newPassword,
				confirm_password: data.confirmPassword,
			});
			reset();
		} catch (err) {
			if (err instanceof ApiError && (err.status === 400 || err.status === 404)) {
				setError('root', { message: 'Current password is incorrect or new password is invalid' });
			} else {
				setError('root', { message: 'Failed to update password, please try again' });
			}
		}
	};

	return (
		<SettingsPanel title="Password" description="Change your account password.">
			<form onSubmit={handleSubmit(onSubmit)} className="space-y-lg">
				<div className="space-y-sm">
					<FormLabel htmlFor="currentPassword">Current password</FormLabel>
					<Input
						id="currentPassword"
						type="password"
						autoComplete="current-password"
						{...register('currentPassword')}
						className={errors.currentPassword ? 'border-error' : ''}
					/>
					{errors.currentPassword && (
						<p className="text-sm text-error">{errors.currentPassword.message}</p>
					)}
				</div>

				<div className="space-y-sm">
					<FormLabel htmlFor="newPassword">New password</FormLabel>
					<Input
						id="newPassword"
						type="password"
						autoComplete="new-password"
						{...register('newPassword')}
						className={errors.newPassword ? 'border-error' : ''}
					/>
					{errors.newPassword && (
						<p className="text-sm text-error">{errors.newPassword.message}</p>
					)}
				</div>

				<div className="space-y-sm">
					<FormLabel htmlFor="confirmPassword">Confirm new password</FormLabel>
					<Input
						id="confirmPassword"
						type="password"
						autoComplete="new-password"
						{...register('confirmPassword')}
						className={errors.confirmPassword ? 'border-error' : ''}
					/>
					{errors.confirmPassword && (
						<p className="text-sm text-error">{errors.confirmPassword.message}</p>
					)}
				</div>

				{errors.root && (
					<p className="text-sm text-error">{errors.root.message}</p>
				)}

				<Button
					type="submit"
					disabled={isSubmitting}
					className={!isOnline ? 'opacity-50 cursor-not-allowed' : ''}
				>
					{isSubmitting ? 'Updating...' : 'Update password'}
				</Button>
			</form>
		</SettingsPanel>
	);
}
