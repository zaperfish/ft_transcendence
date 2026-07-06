'use client';

import { useState } from 'react';
import { useEvents } from '@/lib/hooks/useEvents';
import EventCard from '@/components/features/events/EventCard';
import { Button } from '@/components/ui/Button';
import { useRouter } from "next/navigation";

/**
 * EventsPage displays a list of the current user's events with tab filtering
 * between "attending" and "hosting". It supports pagination and navigation to event detail pages.
 */
export default function EventsPage() {
	const [activeTab, setActiveTab] = useState<'attending' | 'hosting'>('attending');
	const {
		data,
		fetchNextPage,
		hasNextPage,
		isFetchingNextPage,
		isLoading,
		isError,
	} = useEvents(activeTab);
	const router = useRouter();

	const events = data?.pages.flatMap((page) =>
		(page.data || []).map((event) => ({
			...event,
			image: `/api/events/${event.id}/image`
		}))
	) ?? [];

	return (
		<div className="w-full px-xl py-2xl">
		{ /*Header description*/ }
			<div className="mb-2xl">
				<h1 className="text-4xl font-heading font-bold text-text-primary mb-md">
				My Events
				</h1>
				<p className="text-text-secondary text-lg w-full">
				Stay organized with your community schedules.
				</p>
			</div>
		{ /*Events filter tab*/ }
			<div className='flex gap-lg border-border mt-2xl mb-lg'>
				<button
					onClick={() => setActiveTab('attending')}
					className={`w-30 pb-sm text-sm font-medium transition text-center ${
						activeTab === 'attending'
							? 'border-b-2 border-primary text-primary font-semibold'
							: 'text-text-primary hover:text-primary'
					}`}
				>
					Attending
				</button>
				<button
					onClick={() => setActiveTab('hosting')}
					className={`w-30 pb-sm text-sm font-medium transition text-center ${
						activeTab === 'hosting'
							? 'border-b-2 border-primary text-primary font-semibold'
							: 'text-text-primary hover:text-primary'
					}`}
				>
					Hosting
				</button>
			</div>
		{ /*Event cards list*/ }
			{isLoading ? (
				<div className='flex justify-center py-2xl'>
					<span className='text-text-tertiary'>Loading...</span>
				</div>
			) : isError ? (
				<div className='text-center py-2xl text-error'>
					Failed to load events
				</div>
			) : events && events.length > 0 ? (
				<div className='space-y-md'>
					{events.map((event) => (
						<EventCard
							key={event.id}
							data={event}
							mode='detail'
							onDetail={() => router.push(`/events/${event.id}`)} />
					))}
				</div>
			) : (
				<div className='text-center py-2xl text-text-tertiary'>No events found</div>
			)}
		{ /*Load more button*/ }
			{hasNextPage && (
				<div className="flex justify-center mt-2xl">
					<Button
						variant="outline"
						onClick={() => fetchNextPage()}
						disabled={isFetchingNextPage}
						className="min-w-50"
					>
						{isFetchingNextPage ? "Loading..." : "Load more"}
					</Button>
				</div>
			)}
		{ /*All pages loaded*/ }
			{!hasNextPage && events.length > 0 && (
				<p className="text-center text-text-tertiary mt-xl">No more events</p>
			)}
		</div>
	);
}
