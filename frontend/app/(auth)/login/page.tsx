'use client';

import { useRouter, useSearchParams } from "next/navigation";
import { LoginForm } from "@/components/features/auth/LoginForm";
import { useAuth } from "@/lib/hooks/useAuth";
import { toast } from "sonner";

/**
 * The login page component.
 *
 * Renders the login form along with a welcome message. If the `registered`
 * query parameter is present, it displays a success notification for newly
 * registered users. On successful login, redirects to the home page.
 *
 * @returns The login page layout with a heading, optional notification, login form, and registration link.
 */
export default function LoginPage() {
	const router = useRouter();
	const searchParams = useSearchParams();
	const registered = searchParams.get('registered');
	const { isOnline } = useAuth();

	return (
		<div className="rounded-lg bg-surface p-xl shadow-sm">
			<h1 className="mb-sm text-2xl font-bold text-text-primary">Login to Camaraderie</h1>
			<p className="mb-lg text-text-secondary">Welcome back, please login your account</p>
			{registered && (
				<p className="text-success mb-xs">Successfully registered! Please login.</p>
			)}
			{!isOnline && (
				<div className="mb-lg rounded-md bg-yellow-100 p-4 text-sm text-yellow-800">
				"You are currently offline. Login is unvailable."
				</div>
			)}
			<LoginForm
				onSuccess={() => router.push('/home')}
				disabled={!isOnline}
			/>
			<p className="mt-lg text-sm text-text-tertiary">
				Don't have an account?{' '}
				<a
					href={isOnline ? "/register" : "#"}
					onClick={(e) => {
						if (!isOnline) {
							e.preventDefault();
							toast.error("You are currently offline. Registration is unvailable.")
						}
					}}
					className={`text-primary hover:underline ${!isOnline ? 'opacity-50 cursor-not-allowed' : ''
					}`}
					title={!isOnline ? "Registration requires internet connection" : ""}
					>
					Create an account
				</a>
			</p>
		</div>
	);
}