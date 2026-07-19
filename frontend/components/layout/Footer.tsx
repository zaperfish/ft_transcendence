"use client";

import Link from "next/link";
import { useAuth } from "@/lib/hooks/useAuth";
import { useTheme } from "@/lib/context/ThemeContext";
import { ThemeToggle } from "@/components/layout/ThemeToggle";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

export function Footer() {
	const { isOnline } = useAuth();
	const { theme } = useTheme();
	const isClassic = theme === "classic";

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
		<footer
			className={cn(
				"shrink-0 border-t px-4 py-4 text-sm",
				isClassic
					? "border-border bg-background text-muted-foreground"
					: "border-chrome-footer bg-chrome-footer text-chrome-footer backdrop-blur-sm",
			)}
		>
			<div className="grid grid-cols-[1fr_auto_1fr] items-center gap-2">
				<div aria-hidden="true" />
				<div className="flex flex-col items-center justify-center gap-2 sm:flex-row">
					{isOnline ? (
						<Link
							href="/privacy"
							prefetch={false}
							className={isClassic ? "hover:text-foreground" : "hover:text-chrome-title"}
						>
						Privacy Policy
						</Link>
					) : (
						<a {...offlineProps}>Privacy Policy</a>
					)}
					<span className="hidden sm:inline">•</span>
					{isOnline ? (
						<Link
							href="/terms"
							prefetch={false}
							className={isClassic ? "hover:text-foreground" : "hover:text-chrome-title"}
						>
						Terms of Service
						</Link>
					) : (
						<a {...offlineProps}>Terms of Service</a>
					)}
				</div>
				<div className="flex justify-end">
					<ThemeToggle />
				</div>
			</div>
		</footer>
	);
}
