'use client';

import { useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import { useEvents } from '@/lib/hooks/useEvents';
import EventCard from '@/components/features/events/EventCard';
import { Button } from '@/components/ui/Button';
import { useRouter } from "next/navigation";

/**
 * EventsPage displays a list of the current user's events with tab filtering
 * between "attending" and "hosting". It supports pagination and navigation to event detail pages.
 */
export default function EventsPage() {
	const { isOnline } = useAuth();
	const [activeTab, setActiveTab] = useState<'attending' | 'hosting'>('attending');
	const {
		data,
		fetchNextPage,
		hasNextPage,
		isFetchingNextPage,
		isLoading,
		isError,
		fetchStatus, // Used for 'paused' when offline and no cache
	} = useEvents(activeTab);
	const router = useRouter();

	const events = data?.pages.flatMap((page) =>
		(page.data || []).map((event) => ({
			...event,
			image: event.has_image
				? `/api/events/${event.id}/image`
				: '/images/default-event-cover.jpg',
		}))
	) ?? [];

	if (fetchStatus === 'paused' && !isOnline)
		return <div className="text-center py-2xl text-text-secondary">You are offline, no cached data...Please retry after you are online.</div>

	return (
		<div className="w-full px-xl py-2xl">
		{ /*Header description*/ }
			<div className="mb-2xl">
				<h1 className="mb-md font-heading text-4xl font-bold text-chrome-title">
				My Events
				</h1>
				<p className="w-full text-lg text-chrome-body">
				Stay organized with your community schedules.
				</p>
			</div>
		{ /*Events filter tab*/ }
			<div className='mt-2xl mb-lg flex gap-lg border-border'>
				<button
					onClick={() => setActiveTab('attending')}
					className={`w-30 pb-sm text-center text-sm font-medium transition ${
						activeTab === 'attending'
							? 'border-b-2 border-chrome-tab-active font-semibold text-chrome-tab-active'
							: 'text-chrome-muted hover:text-chrome-title'
					}`}
				>
					Attending
				</button>
				<button
					onClick={() => setActiveTab('hosting')}
					className={`w-30 pb-sm text-center text-sm font-medium transition ${
						activeTab === 'hosting'
							? 'border-b-2 border-chrome-tab-active font-semibold text-chrome-tab-active'
							: 'text-chrome-muted hover:text-chrome-title'
					}`}
				>
					Hosting
				</button>
			</div>
		{ /*Event cards list*/ }
			{isLoading ? (
				<div className='flex justify-center py-2xl'>
					<span className='text-chrome-muted'>Loading...</span>
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
							layout='horizontal'
							disabled={!isOnline}
							onDetail={isOnline ? () => router.push(`/events/${event.id}`) : undefined} />
					))}
				</div>
			) : (
				<div className='py-2xl text-center text-chrome-muted'>No events found</div>
			)}
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
				<p className="mt-xl text-center text-chrome-muted">No more events</p>
			)}
		</div>
	);
}
