"use client";

import Link from "next/link";
import { useAuth } from "@/lib/hooks/useAuth";
import { toast } from "sonner";

export function Footer() {
	const { isOnline } = useAuth();

	const offlineProps = {
		href: "#",
		onClick: (e: React.MouseEvent) => {
			e.preventDefault();
			toast.error("You are currently offline. This page is unvailable.");
		},
		className: 'hover:underline opacity-50 cursor-not-allowed',
		'aria-disabled': true,
	};

	const onlinePrivacyProps = {
		href: "/privacy",
		className: 'hover:text-teal-50'
	};

	const onlineTermsProps = {
		href: "/terms",
		className: 'hover:text-teal-50'
	};

	return (
		<footer className="shrink-0 border-t border-white/10 bg-slate-950/70 px-4 py-4 text-sm text-teal-100/70 backdrop-blur-sm">
			<div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-center">
				{isOnline ? (
					<Link href="/privacy" className="hover:text-teal-50">
					Privacy Policy
					</Link>
				) : (
					<a {...offlineProps}>Privacy Policy</a>
				)}
				<span className="hidden sm:inline">•</span>
				{isOnline ? (
					<Link href="/terms" className="hover:text-teal-50">
					Terms of Service
					</Link>
				) : (
					<a {...offlineProps}>Terms of Service</a>
				)}
			</div>
		</footer>
	);
}