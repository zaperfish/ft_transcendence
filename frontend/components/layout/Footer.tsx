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

	return (
		<footer className="shrink-0 border-t border-chrome-footer bg-chrome-footer px-4 py-4 text-sm text-chrome-footer backdrop-blur-sm">
			<div className="flex flex-col items-center justify-center gap-2 sm:flex-row">
				{isOnline ? (
					<Link href="/privacy" className="hover:text-chrome-title">
					Privacy Policy
					</Link>
				) : (
					<a {...offlineProps}>Privacy Policy</a>
				)}
				<span className="hidden sm:inline">•</span>
				{isOnline ? (
					<Link href="/terms" className="hover:text-chrome-title">
					Terms of Service
					</Link>
				) : (
					<a {...offlineProps}>Terms of Service</a>
				)}
			</div>
		</footer>
	);
}
