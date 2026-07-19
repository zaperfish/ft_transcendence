'use client';

import { useTheme } from "@/lib/context/ThemeContext";
import { cn } from "@/lib/utils";

interface CreateEventCardProps {
	onClick?: () => void;
}

/**
 * CreateEventCard is a clickable card component that serves as a visual entry point
 * for users to create a new event. It displays a plus icon, title, and description,
 * and triggers the provided onClick callback when clicked.
 */
export default function CreateEventCard({ onClick }: CreateEventCardProps) {
	const { theme } = useTheme();
	const isClassic = theme === "classic";

	return (
		<div
			onClick={onClick}
			className={cn(
				"flex h-full min-h-[350px] cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed transition-colors duration-200",
				isClassic
					? "border-border hover:bg-surface-dim"
					: "border-teal-700/25 bg-white shadow-md hover:bg-teal-50",
			)}
		>
			<div
				className={cn(
					"mb-md flex h-12 w-12 items-center justify-center rounded-full",
					isClassic ? "bg-surface-container" : "bg-teal-100",
				)}
			>
				<PlusIcon
					className={cn(
						"size-6",
						isClassic ? "text-text-tertiary" : "text-teal-700",
					)}
				/>
			</div>
			<h3
				className={cn(
					"text-lg font-semibold",
					isClassic ? "text-text-secondary" : "text-text-primary",
				)}
			>
				Create Event
			</h3>
			<p
				className={cn(
					"mt-xs text-sm",
					isClassic ? "text-text-tertiary" : "text-text-secondary",
				)}
			>
				Host your own event
			</p>
		</div>
	);
}

function PlusIcon({ className }: { className?: string }) {
	return (
		<svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
			<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4"/>
		</svg>
	);
}
