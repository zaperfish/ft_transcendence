'use client';

import { useState } from "react";
import CreateEventCard from "@/components/features/events/CreateEventCard";
import CreateEventForm from "@/components/features/events/CreateEventForm";

export default function HomePage() {
	const [isFormOpen, setIsFormOpen] = useState(false);

	const handleOpenForm = () => {
		setIsFormOpen(true);
	};

	const handleCloseForm = () => {
		setIsFormOpen(false);
	};

	// Later should add refresh and update
	const handleEventCreated = () => {
		console.log("Event created successfully!");
	};

	return (
		<div className="w-full px-xl py-2xl">
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
			<div className="grid grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-lg">
				<CreateEventCard onClick={handleOpenForm}/>
				{ /*Cards to be rendered*/ }
			</div>
			<CreateEventForm
				open={isFormOpen}
				onClose={handleCloseForm}
				onSuccess={handleEventCreated}
			/>
		</div>
	);
}