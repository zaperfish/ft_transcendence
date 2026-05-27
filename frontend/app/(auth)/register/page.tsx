'use client';

import { RegisterForm } from "@/components/features/auth/RegisterForm";


/**
 * The registration page component.
 *
 * Renders the registration form with a welcome message and a link to
 * the login page for users who already have an account.
 *
 * @returns The registration page layout with a heading, description, registration form, and login link.
 */
export default function RegisterPage() {

	return (
		<div className="rounded-lg bg-surface p-xl shadow-sm">
			<h1 className="mb-sm text-2xl font-bold text-text-primary">Register</h1>
			<p className="mb-lg text-text-secondary">Welcome, please create your account</p>
			<RegisterForm />
			<p className="mt-lg text-sm text-text-tertiary">
				Already have an account? <a href="/login" className="text-primary hover:underline">Login</a>
			</p>
		</div>
	);
}