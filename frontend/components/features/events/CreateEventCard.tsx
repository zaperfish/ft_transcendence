interface CreateEventCardProps {
	onClick?: () => void;
}

export default function CreateEventCard({ onClick }: CreateEventCardProps) {
	return (
		<div
			onClick={onClick}
			className="flex flex-col items-center justify-center border-2 border-dashed border-border rounded-lg aspect-3/4 hover:bg-surface-dim cursor-pointer transition-colors duration-200"
		>
			<div className="w-12 h-12 rounded-full bg-surface-container flex items-center justify-center mb-md">
				<PlusIcon className="size-6 text-text-tertiary"/>
			</div>
			<h3 className="text-lg font-semibold text-text-secondary">Create Event</h3>
			<p className="text-sm text-text-tertiary mt-xs">Host your own event</p>
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