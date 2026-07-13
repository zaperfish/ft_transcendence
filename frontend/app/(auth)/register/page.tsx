'use client';

import { RegisterForm } from "@/components/features/auth/RegisterForm";
import { useAuth } from "@/lib/hooks/useAuth";


/**
 * The registration page component.
 *
 * Renders the registration form with a welcome message and a link to
 * the login page for users who already have an account.
 *
 * @returns The registration page layout with a heading, description, registration form, and login link.
 */
export default function RegisterPage() {
	const { isOnline } = useAuth();

	return (
		<div className="rounded-lg bg-surface p-xl shadow-sm">
			<h1 className="mb-sm text-2xl font-bold text-text-primary">Register</h1>
			<p className="mb-lg text-text-secondary">Welcome, please create your account</p>
			{!isOnline && (
				<div className="mb-lg rounded-md bg-yellow-100 p-4 text-sm text-yellow-800">
				"You are currently offline. Register is unvailable."
				</div>
			)}
			<RegisterForm disabled={!isOnline}/>
			<p className="mt-lg text-sm text-text-tertiary">
				Already have an account?
				<a
					href={isOnline ? "/login" : "#"}
					onClick={(e) => {
						if (!isOnline) {
							e.preventDefault();
							alert("You are currently offline. Login is unvailable.")
						}
					}}
					className={`text-primary hover:underline ${!isOnline ? 'opacity-50 cursor-not-allowed' : ''
					}`}
					title={!isOnline ? "Login requires internet connection" : ""}
					>Login</a>
			</p>
		</div>
	);
}