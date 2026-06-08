'use client';

import { useState } from "react";
import CreateEventCard from "@/components/features/events/CreateEventCard";
import CreateEventForm from "@/components/features/events/CreateEventForm";
import EventCard from "@/components/features/events/EventCard";
import { QueryClient, useInfiniteQuery, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/Button";

const PAGE_SIZE = 7

export default function HomePage() {
	const [isFormOpen, setIsFormOpen] = useState(false);
	const queryClient = useQueryClient();

	const handleOpenForm = () => {
		setIsFormOpen(true);
	};

	const handleCloseForm = () => {
		setIsFormOpen(false);
	};

	// Later should add refresh and update
	const handleEventCreated = () => {
		queryClient.invalidateQueries({ queryKey: ["events"] });
	};

	const {
		data,
		fetchNextPage,
		hasNextPage,
		isFetchingNextPage,
		isLoading,
		isError,
	} = useInfiniteQuery({
		queryKey: ["events"],
		queryFn: async ({ pageParam = 1 }) => {
			const res = await fetch(`/api/events?page=${pageParam}&page_size=${PAGE_SIZE}`);
			return res.json();
		},
		getNextPageParam: (lastPage) => {
			const { page, page_size, total } = lastPage;
			const maxPage = Math.ceil(total / page_size);
			return page < maxPage ? page + 1 : undefined;
		},
		initialPageParam: 1,
	});

	const events = data?.pages.flatMap((page) => page.data) ?? [];

	if (isLoading)
		return <div className="text-center py-2xl text-text-secondary">Loading...</div>
	if (isError)
		return <div className="text-center py-2xl text-error">Failed to load events...</div>

	return (
		<div className="w-full px-xl py-2xl">
		{ /*Header description*/ }
			<div className="mb-2xl">
				<h1 className="text-4xl font-heading font-bold text-text-primary mb-md">
				Discover Events
				</h1>
				<p className="text-text-secondary text-lg w-full">
				Connect with your local community through shared interests,
				professional workshops, and social gatherings.
				</p>
			</div>
		{ /*Card grid container*/ }
			<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-lg">
				<CreateEventCard onClick={handleOpenForm}/>
				{events.map((event) => (
					<EventCard key={event.id} data={event} />
				))}
			</div>
		{ /*Load more button*/ }
			{hasNextPage && (
				<div className="flex justify-center mt-2xl">
					<Button
						variant="outline"
						onClick={() => fetchNextPage()}
						disabled={isFetchingNextPage}
						className="min-w-[200px]"
					>
						{isFetchingNextPage ? "Loading..." : "Load more"}
					</Button>
				</div>
			)}
		{ /*All pages loaded*/ }
			{!hasNextPage && events.length > 0 && (
				<p className="text-center text-text-tertiary mt-xl">No more events</p>
			)}

			<CreateEventForm
				open={isFormOpen}
				onClose={handleCloseForm}
				onSuccess={handleEventCreated}
			/>
		</div>
	);
}