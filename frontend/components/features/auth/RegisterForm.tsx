'use client';

import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { register as registerApi } from '@/lib/api/auth';
import { useRouter } from 'next/navigation';

// Define the schema of validating register data
const registerSchema = z
	.object({
		username: z
			.string()
			.min(1, 'Please enter your username') // This field cannot be empty
			.min(3, 'Username should be at least 3 characters long')
			.max(20, 'Username should be no more than 20 characters')
			.regex(/^\S+$/, 'Username should not contain whitespace'),
		email: z.pipe( // Zod v4 not support email calling in chain
			z.string().min(1, 'Please enter your email address'),
			z.email('Invalid email address')
		),
		password: z
			.string()
			.min(1, 'Please enter your password')
			.min(8, 'Password should be at least 8 characters long')
			.max(128, 'Password should be no more than 128 characters')
			.regex(/[A-Z]/, 'Password should contain at least an uppercase letter')
			.regex(/[a-z]/, 'Password should contain at least an lowercase letter')
			.regex(/\d/, 'Password should contain at least a number'),
		confirmPassword: z.string().min(1, 'Please confirm your password'),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: 'Passwords do not match',
		path: ['confirmPassword'], // Need specify path because refine() is global validation
	});

// Typescript enables infering type of each attribute in object automatically
type RegisterFormData = z.infer<typeof registerSchema>;

interface RegisterFormProps {
	disabled?: boolean;
}

/**
 * RegisterForm is a form component for user registration.
 * It captures username, email, password, and password confirmation,
 * validates input using Zod schema, submits data to the registration API,
 * and handles errors such as duplicate credentials or network issues.
 */
export function RegisterForm({ disabled = false }: RegisterFormProps) {
	const router = useRouter();
	const {
		register,
		handleSubmit,
		setError,
		formState: { errors, isSubmitting },
	} = useForm<RegisterFormData>({
		resolver: zodResolver(registerSchema),
		mode: 'onBlur', // Remind error when user leaves one field
		defaultValues: {
			username: '',
			email: '',
			password: '',
			confirmPassword: '',
		},
	});
	// Data is a frontend object collecting from user input and will send to server later
	const onSubmit = async (data: RegisterFormData) => {
		if (disabled) {
			setError('root', { message: 'You are offline. Please try again later.' });
			return;
		}

		try {
			await registerApi({
				name: data.username,
				email: data.email,
				password: data.password,
				password_confirm: data.confirmPassword,
			});
			router.push('/login?registered=true');
		} catch (err: any) {
			if (err.status === 409) {
				setError('root', { message: 'Occupied username or email address' });
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
					{...register('username')} // Register passing data collected from user input to 'Input'
					className={errors.username ? 'border-error' : ''}
				/>
				{errors.username && (
					<p className='text-sm text-error mt-xs'>{errors.username.message}</p>
				)}
			</div>
			{/* Email */}
			<div>
				<Input
					type='email'
					placeholder='Email'
					{...register('email')}
					className={errors.email ? 'border-error' : ''}
				/>
				{errors.email && (
					<p className='text-sm text-error mt-xs'>{errors.email.message}</p>
				)}
			</div>
			{/* Password */}
			<div>
				<Input
					type='password'
					placeholder='Password'
					{...register('password')}
					className={errors.password ? 'border-error' : ''}
				/>
				{errors.password && (
					<p className='text-sm text-error mt-xs'>{errors.password.message}</p>
				)}
			</div>
			{/* Confirm Password */}
			<div>
				<Input
					type='password'
					placeholder='Confirm Password'
					{...register('confirmPassword')}
					className={errors.confirmPassword ? 'border-error' : ''}
				/>
				{errors.confirmPassword && (
					<p className='text-sm text-error mt-xs'>{errors.confirmPassword.message}</p>
				)}
			</div>
			{/* Root errors */}
			{errors.root && (
				<p className='text-sm text-error mt-xs'>{errors.root.message}</p>
			)}
			<Button type='submit' disabled={isSubmitting || disabled}>
				{isSubmitting ? "Registering..." : "Register"}
			</Button>
		</form>
	);
}

// This is the previous version of RegisterForm
// /**
//  * A registration form component with client-side validation and error handling.
//  *
//  * Collects username, email, and password (need confirmation), validates inputs, and submits
//  * the data to the registration API. On success, redirects to the login page
//  * with a query parameter indicating successful registration. Displays
//  * field-level errors and general API errors (e.g., duplicate user or network issues).
//  *
//  * @returns A styled form containing input fields and a submit button.
//  */
// export function RegisterForm() {
// 	const [username, setUsername] = useState('');
// 	const [email, setEmail] = useState('');
// 	const [password, setPassword] = useState('');
// 	const [confirmPassword, setConfirmPassword] =useState('');
// 	const [errors, setErrors] = useState< { username?: string; email?: string; password?: string; confirmPassword?: string; general?: string }>({});
// 	const [loading, setLoading] = useState(false);
// 	const router = useRouter();

// 	const validate = () => {
// 		const newErrors: typeof errors = {};
// 		if (!username.trim())
// 			newErrors.username = "Please enter your username";
// 		if (!email.trim())
// 			newErrors.email = 'Please enter your email';
// 		else if (!/\S+@\S+\.\S+/.test(email))
// 			newErrors.email = 'Invalid email address';
// 		if (!password.trim())
// 			newErrors.password = "Please enter your password";
// 		if (!confirmPassword.trim())
// 			newErrors.confirmPassword = "Please confirm your password";
// 		else if (confirmPassword !== password)
// 			newErrors.confirmPassword = "Passwords do not match";
// 		setErrors(newErrors);
// 		return Object.keys(newErrors).length === 0;
// 	};

// 	const handleSubmit = async() => {
// 		if (!validate())
// 			return;
// 		setLoading(true);
// 		setErrors({});
// 		try {
// 			await register({ name: username, email, password, password_confirm: confirmPassword });
// 			router.push('/login?registered=true');
// 		} catch (err: any) {
// 			if (err.status === 409) {
// 				setErrors({ general: 'Occupied username or email'});
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
// 			<Input value={email} onChange={e => setEmail(e.target.value)} placeholder='Email' error={errors.email} />
// 			<Input value={password} onChange={e => setPassword(e.target.value)} placeholder='Password' error={errors.password} />
// 			<Input value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} placeholder='Confirm Password' error={errors.confirmPassword} />
// 			{errors.general && <div className="text-sm text-error mt-xs">{errors.general}</div>}
// 			<Button disabled={loading} onClick={handleSubmit}>{loading ? "loading..." : "register"}</Button>
// 		</div>
// 	);
// }
