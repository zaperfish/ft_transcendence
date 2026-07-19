'use client';

import { useState } from "react";
import { useAuth } from "@/lib/hooks/useAuth";
import CreateEventCard from "@/components/features/events/CreateEventCard";
import CreateEventForm from "@/components/features/events/CreateEventForm";
import EventCard from "@/components/features/events/EventCard";
import { useInfiniteQuery, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/Button";
import { getEvents } from "@/lib/api/events";
import { toast } from "sonner";

const PAGE_SIZE = 7

/**
 * HomePage is the main landing page that displays a grid of event cards,
 * includes a "Create Event" card to open the creation modal, and supports
 * infinite scrolling pagination for loading more events.
 *
 * Enable offline mode: create event and register disabled;
 * when cached data exist, query uses cache instead of sending request;
 * when no cached data, fetch status is paused until network continues.
 */
export default function HomePage() {
	const { isOnline } = useAuth();
	const [isFormOpen, setIsFormOpen] = useState(false);
	const queryClient = useQueryClient();

	const handleOpenForm = () => {
		if (!isOnline) {
			toast.error('Cannot create event when offline, please try later');
			return;
		}
		setIsFormOpen(true);
	};

	const handleCloseForm = () => {
		setIsFormOpen(false);
	};

	// Refresh and update home page when a new event is created
	const handleEventCreated = (warning?: string) => {
		queryClient.invalidateQueries({ queryKey: ["events"] });
		if (warning) {
			toast.warning(warning);
		}
	};

	const {
		data,
		fetchNextPage,
		hasNextPage,
		isFetchingNextPage,
		isLoading,
		isError,
		fetchStatus,
	} = useInfiniteQuery({
		queryKey: ["events"],
		queryFn: async ({ pageParam = 1 }) => {
			return getEvents({ page: pageParam, page_size: PAGE_SIZE });
		},
		getNextPageParam: (lastPage) => {
			const { page, page_size, total } = lastPage;
			const maxPage = Math.ceil(total / page_size);
			return page < maxPage ? page + 1 : undefined;
		},
		initialPageParam: 1,
		networkMode: 'offlineFirst',
	});

	// Using [] to avoid crash when backend returns empty page.data
	// Add image interface and then pass it to src to retrieve img automatically by browser
	const events = data?.pages.flatMap((page) =>
		(page.data || []).map((event) => ({
			...event,
			image: event.has_image
				? `/api/events/${event.id}/image`
				: '/images/default-event-cover.jpg',
		}))
	) ?? [];

	if (isLoading && !isOnline)
		return <div className="text-center py-2xl text-text-secondary">Loading...Checking network connection...</div>
	if (fetchStatus === 'paused' && !isOnline)
		return <div className="text-center py-2xl text-text-secondary">You are offline, no cached data...Please retry after you are online.</div>
	if (isError)
		return <div className="text-center py-2xl text-error">Failed to load events...</div>

	return (
		<div className="w-full px-xl py-2xl">
		{ /*Header description*/ }
			<div className="mb-2xl">
				<h1 className="mb-md font-heading text-4xl font-bold text-chrome-title">
				Discover Events
				</h1>
				<p className="w-full text-lg text-chrome-body">
				Connect with your local community through shared interests,
				professional workshops, and social gatherings.
				</p>
			</div>
		{ /*Card grid container*/ }
			<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-lg">
				<CreateEventCard onClick={handleOpenForm}/>
				{events.map((event) => (
					<EventCard key={event.id} data={event} disabled={!isOnline}/>
				))}
			</div>
		{ /*Load more button*/ }
			{hasNextPage && (
				<div className="flex justify-center mt-2xl">
					<Button
						variant="outline"
						onClick={() => fetchNextPage()}
						disabled={isFetchingNextPage || !isOnline}
						className="btn-chrome-loadmore min-w-50"
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