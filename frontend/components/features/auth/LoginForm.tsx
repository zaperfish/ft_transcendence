'use client';

import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { useAuth } from '@/lib/hooks/useAuth';

interface LoginFormProps {
	onSuccess?: () => void;
	disabled?: boolean;
}

// Define the schema of validating login data
const loginSchema = z
	.object({
		username: z
			.string()
			.min(1, 'Please enter your username'), // This field cannot be empty
		password: z
			.string()
			.min(1, 'Please enter your password'),
	});

// Typescript enables infering type of each attribute in object automatically
type LoginFormData = z.infer<typeof loginSchema>;

/**
 * LoginForm is a form component for user authentication.
 * It captures username and password, validates input using Zod schema,
 * calls the login API via the auth context, handles authentication errors,
 * and invokes the onSuccess callback upon successful login.
 */
export function LoginForm({ onSuccess, disabled = false }: LoginFormProps) {
	const { login } = useAuth();
	const {
		register,
		handleSubmit,
		setError,
		formState: { errors, isSubmitting },
	} = useForm<LoginFormData>({
		resolver: zodResolver(loginSchema),
		mode: 'onBlur', // Remind error when user leaves one field
		defaultValues: {
			username: '',
			password: '',
		},
	});

	const onSubmit = async (data: LoginFormData) => {
		if (disabled) {
			setError('root', { message: 'You are offline. Please try again later.' });
			return;
		}

		try {
			await login({
				name: data.username,
				password: data.password,
			});
			onSuccess?.();
		} catch (err: any) {
			if (err.status === 401) {
				setError('root', { message: 'Invalid username or password' });
			} else {
				setError('root', { message: 'Network error, please try later'});
			}
		}
	};

	return (
		<form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-6">
			{/* Username */}
			<div>
				<Input
					type='text' // Specify the type of input to enable browser behaviors
					placeholder='Username'
					autoComplete='username'
					{...register('username')} // Register passing data collected from user input to 'Input'
					className={errors.username ? 'border-error' : ''}
				/>
				{errors.username && (
					<p className='text-sm text-error mt-xs'>{errors.username.message}</p>
				)}
			</div>
			{/* Password */}
			<div>
				<Input
					type='password'
					placeholder='Password'
					autoComplete='current-password'
					{...register('password')}
					className={errors.password ? 'border-error' : ''}
				/>
				{errors.password && (
					<p className='text-sm text-error mt-xs'>{errors.password.message}</p>
				)}
			</div>
			{/* Root errors */}
			{errors.root && (
				<p className='text-sm text-error mt-xs'>{errors.root.message}</p>
			)}
			<Button type='submit' disabled={isSubmitting || disabled}>
				{isSubmitting ? "Loading..." : "Login"}
			</Button>
		</form>
	);
}

// This is the previous version of LoginForm
// /**
//  * A login form component with client-side validation and authentication.
//  *
//  * Collects username and password, validates inputs, and calls the login method
//  * from the auth context. On success, invokes the optional `onSuccess` callback.
//  * Displays field-level and general API errors (invalid credentials or network issues).
//  *
//  * @param props - The component props.
//  * @param props.onSuccess - Optional callback executed after successful login.
//  * @returns A styled form containing input fields and a submit button.
//  */
// export function LoginForm({ onSuccess }: LoginFormProps) {
// 	const [username, setUsername] = useState('');
// 	const [password, setPassword] = useState('');
// 	const [errors, setErrors] = useState< {username?: string; password?: string; general?: string}>({});
// 	const [loading, setLoading] = useState(false);
// 	const { login } = useAuth();

// 	const validate = () => {
// 		const newErrors: typeof errors = {};
// 		if (!username.trim())
// 			newErrors.username = "Please enter your username";
// 		if (!password.trim())
// 			newErrors.password = "Please enter your password";
// 		setErrors(newErrors);
// 		return Object.keys(newErrors).length === 0;
// 	};

// 	const handleSubmit = async() => {
// 		if (!validate())
// 			return;
// 		setLoading(true);
// 		setErrors({});
// 		try {
// 			await login({ name: username, password });
// 			onSuccess?.();
// 		} catch (err: any) {
// 			if (err.status === 401) {
// 				setErrors({ general: 'Invalid username or password'});
// 			} else {
// 				setErrors({ general: 'Network error'});
// 			}
// 		} finally {
// 			setLoading(false);
// 		}
// 	};

// 	return (
// 		<div className="flex flex-col gap-6">
// 			<Input value={username} onChange={e => setUsername(e.target.value)} placeholder='Username' error={errors.username} />
// 			<Input value={password} onChange={e => setPassword(e.target.value)} placeholder='Password' error={errors.password} />
// 			{errors.general && <div className="text-sm text-error mt-xs">{errors.general}</div>}
// 			<Button disabled={loading} onClick={handleSubmit}>{loading ? "loading..." : "login"}</Button>
// 		</div>
// 	);
// }
