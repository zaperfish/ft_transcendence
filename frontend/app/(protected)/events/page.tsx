'use client';

import { useState } from 'react';
import { useEvents } from '@/lib/hooks/useEvents';
import EventCard from '@/components/features/events/EventCard';
import { Button } from '@/components/ui/Button';
import { useRouter } from "next/navigation";

export default function EventsPage() {
	const [activeTab, setActiveTab] = useState<'attending' | 'hosting'>('attending');
	const { data: events, isLoading, isError } = useEvents(activeTab);
	const router = useRouter();

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
				<Button
					onClick={() => setActiveTab('attending')}
					className={`w-[120px] pb-sm text-sm font-medium transition text-center ${
						activeTab === 'attending'
							? 'border-b-2 border-primary text-white font-semibold'
							: 'text-text-tertiary hover:text-white'
					}`}
				>
					Attending
				</Button>
				<Button
					onClick={() => setActiveTab('hosting')}
					className={`w-[120px] pb-sm text-sm font-medium transition text-center ${
						activeTab === 'hosting'
							? 'border-b-2 border-primary text-white font-semibold'
							: 'text-text-tertiary hover:text-white'
					}`}
				>
					Hosting
				</Button>
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
		</div>
	);
}
