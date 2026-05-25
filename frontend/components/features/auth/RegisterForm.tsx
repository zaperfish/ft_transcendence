'use client';

import { useState } from 'react';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { register } from '@/lib/api/auth';
import { useRouter } from 'next/navigation';

interface LoginFormProps {
	onSuccess?: () => void;
}

export function LoginForm({ onSuccess }: LoginFormProps) {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');
	const [errors, setErrors] = useState< {username?: string; password?: string; general?: string}>({});
	const [loading, setLoading] = useState(false);
	const { login } = useAuth();

	const validate = () => {
		const newErrors: typeof errors = {};
		if (!username.trim())
			newErrors.username = "Please enter your username";
		if (!password.trim())
			newErrors.password = "Please enter your password";
		setErrors(newErrors);
		return Object.keys(newErrors).length === 0;
	};

	const handleSubmit = async() => {
		if (!validate())
			return;
		setLoading(true);
		setErrors({});
		try {
			await login({ username, password });
			onSuccess?.();
		} catch (err: any) {
			if (err.status === 401) {
				setErrors({ general: 'Invalid username or password'});
			} else {
				setErrors({ general: 'Network error'});
			}
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="flex flex-col gap-6">
			<Input value={username} onChange={e => setUsername(e.target.value)} placeholder='Username' error={errors.username} />
			<Input value={password} onChange={e => setPassword(e.target.value)} placeholder='Password' error={errors.password} />
			{errors.general && <div className="text-sm text-error mt-xs">{errors.general}</div>}
			<Button loading={loading} onClick={handleSubmit}>login</Button>
		</div>
	);
}
