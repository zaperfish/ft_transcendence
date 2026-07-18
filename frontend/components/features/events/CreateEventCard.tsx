interface CreateEventCardProps {
	onClick?: () => void;
}

/**
 * CreateEventCard is a clickable card component that serves as a visual entry point
 * for users to create a new event. It displays a plus icon, title, and description,
 * and triggers the provided onClick callback when clicked.
 */
export default function CreateEventCard({ onClick }: CreateEventCardProps) {
	return (
		<div
			onClick={onClick}
			className="flex h-full min-h-[350px] cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-teal-700/25 bg-white shadow-md transition-colors duration-200 hover:bg-teal-50"
		>
			<div className="mb-md flex h-12 w-12 items-center justify-center rounded-full bg-teal-100">
				<PlusIcon className="size-6 text-teal-700"/>
			</div>
			<h3 className="text-lg font-semibold text-text-primary">Create Event</h3>
			<p className="mt-xs text-sm text-text-secondary">Host your own event</p>
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