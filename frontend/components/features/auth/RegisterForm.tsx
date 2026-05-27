'use client';

import { useState } from 'react';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { register } from '@/lib/api/auth';
import { useRouter } from 'next/navigation';

/**
 * A registration form component with client-side validation and error handling.
 *
 * Collects username, email, and password, validates inputs, and submits
 * the data to the registration API. On success, redirects to the login page
 * with a query parameter indicating successful registration. Displays
 * field-level errors and general API errors (e.g., duplicate user or network issues).
 *
 * @returns A styled form containing input fields and a submit button.
 */
export function RegisterForm() {
	const [username, setUsername] = useState('');
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');
	const [errors, setErrors] = useState< { username?: string; email?: string; password?: string; general?: string }>({});
	const [loading, setLoading] = useState(false);
	const router = useRouter();

	const validate = () => {
		const newErrors: typeof errors = {};
		if (!username.trim())
			newErrors.username = "Please enter your username";
		if (!email.trim())
			newErrors.email = 'Please enter your email';
		else if (!/\S+@\S+\.\S+/.test(email))
			newErrors.email = 'Invalid email address';
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
			await register({ name: username, email, password });
			router.push('/login?registered=true');
		} catch (err: any) {
			if (err.status === 409) {
				setErrors({ general: 'Occupied username or email'});
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
			<Input value={email} onChange={e => setEmail(e.target.value)} placeholder='Email' error={errors.email} />
			<Input value={password} onChange={e => setPassword(e.target.value)} placeholder='Password' error={errors.password} />
			{errors.general && <div className="text-sm text-error mt-xs">{errors.general}</div>}
			<Button loading={loading} onClick={handleSubmit}>register</Button>
		</div>
	);
}
