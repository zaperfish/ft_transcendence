'use client';

import { useRouter } from "next/navigation";
import { LoginForm } from "@/components/features/auth/LoginForm";

export default function LoginPage() {
	const router = useRouter();

	return (
		<div className="rounded-lg bg-surface p-xl shadow-sm">
			<h1 className="mb-sm text-2xl font-bold text-text-primary">Login to Meetup</h1>
			<p className="mb-lg text-text-secondary">Welcome back, please login your account</p>
			<LoginForm onSuccess={() => router.push('/home')} />
			<p className="mt-lg text-sm text-text-tertiary">
				Don't have an account? <a href="/register" className="text-primary hover:underline">Create an account</a>
			</p>
		</div>
	);
}