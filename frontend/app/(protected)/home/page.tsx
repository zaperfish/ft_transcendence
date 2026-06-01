import CreateEventCard from "@/components/features/events/CreateEventCard";

export default function HomePage() {
	return (
		<div className="max-w-360 mx-auto px-lg py-2xl">
		{ /*Header description*/ }
			<div className="mb-2xl">
				<h1 className="text-4xl font-heading font-bold text-text-primary mb-md">
				Discover Events
				</h1>
				<p className="text-text-secondary text-lg max-w-360">
				Connect with your local community through shared interests,
				professional workshops, and social gatherings.
				</p>
			</div>
		{ /*Card grid container*/ }
			<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-lg">
				<CreateEventCard />
				{ /*Cards to be rendered*/ }
			</div>
		</div>
	);
}