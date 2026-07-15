"use client";

import Link from "next/link";
import { useAuth } from "@/lib/hooks/useAuth";

export function Footer() {
	const { isOnline } = useAuth();

	const offlineProps = {
		href: "#",
		onClick: (e: React.MouseEvent) => {
			e.preventDefault();
			alert("You are currently offline. This page is unvailable.");
		},
		className: 'hover:underline opacity-50 cursor-not-allowed',
		'aria-disabled': true,
	};

	const onlinePrivacyProps = {
		href: "/privacy",
		className: 'hover:text-foreground'
	};

	const onlineTermsProps = {
		href: "/terms",
		className: 'hover:text-foreground'
	};

	return (
		<footer className="shrink-0 border-t border-border bg-background px-4 py-4 text-sm text-muted-foreground">
			<div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-center">
				{isOnline ? (
					<Link href="/privacy" className="hover:text-foreground">
					Privacy Policy
					</Link>
				) : (
					<a {...offlineProps}>Privacy Policy</a>
				)}
				<span className="hidden sm:inline">•</span>
				{isOnline ? (
					<Link href="/privacy" className="hover:text-foreground">
					Terms of Service
					</Link>
				) : (
					<a {...offlineProps}>Terms of Service</a>
				)}
			</div>
		</footer>
	);
}