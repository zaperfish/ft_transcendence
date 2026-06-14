'use client';

import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { getEvent, getEventParticipants } from '@/lib/api/events';
import { getEventChatMessages } from '@/lib/api/chat';
import { ApiError } from '@/lib/api/client';

interface EventChatRoomProps {
	eventId: number;
}

export default function EventChatRoom({ eventId }: EventChatRoomProps) {
	const [draft, setDraft] = useState('');

	const historyQuery = useQuery({
		queryKey: ['chat-history', eventId],
		queryFn: () => getEventChatMessages(eventId),
		retry: false,
	});

	const participantsQuery = useQuery({
		queryKey: ['event-participants', eventId],
		queryFn: () => getEventParticipants(eventId),
		retry: false,
	});

	const eventQuery = useQuery({
		queryKey: ['event', eventId],
		queryFn: () => getEvent(eventId),
		retry: false,
	});

	const messages = historyQuery.data?.data ?? [];
	const participants = participantsQuery.data?.data ?? [];
	const eventTitle = eventQuery.data?.title ?? `Event #${eventId}`;

	const participantNamesById = new Map(
		participants.map((participant) => [participant.id, participant.name])
	);

	const historyError =
		historyQuery.error instanceof ApiError ? historyQuery.error : null;

	if (historyQuery.isLoading) {
		return (
			<div className="min-h-screen bg-surface-dim p-lg">
				<div className="mx-auto flex h-[calc(100vh-2rem)] max-w-3xl items-center justify-center rounded-lg border border-border bg-surface shadow-sm">
					<p className="text-sm text-text-secondary">Loading chat...</p>
				</div>
			</div>
		);
	}

	if (historyError?.status === 403) {
		return (
			<div className="min-h-screen bg-surface-dim p-lg">
				<div className="mx-auto flex h-[calc(100vh-2rem)] max-w-3xl flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
					<h1 className="text-2xl font-heading font-bold text-text-primary">
						Access denied
					</h1>
					<p className="mt-sm text-sm text-text-secondary">
						You must be a participant in this event to access the chat room.
					</p>
				</div>
			</div>
		);
	}

	if (historyError?.status === 404) {
		return (
			<div className="min-h-screen bg-surface-dim p-lg">
				<div className="mx-auto flex h-[calc(100vh-2rem)] max-w-3xl flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
					<h1 className="text-2xl font-heading font-bold text-text-primary">
						Event not found
					</h1>
					<p className="mt-sm text-sm text-text-secondary">
						This event does not exist or is no longer available.
					</p>
				</div>
			</div>
		);
	}

	if (historyQuery.isError) {
		return (
			<div className="min-h-screen bg-surface-dim p-lg">
				<div className="mx-auto flex h-[calc(100vh-2rem)] max-w-3xl flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
					<h1 className="text-2xl font-heading font-bold text-text-primary">
						Failed to load chat
					</h1>
					<p className="mt-sm text-sm text-text-secondary">
						An unexpected error occurred while loading chat history.
					</p>
				</div>
			</div>
		);
	}

	return (
		<div className="min-h-screen bg-surface-dim p-lg">
			<div className="flex h-[calc(100vh-2rem)] w-full flex-col">
				<div className="mb-md">
					<h1 className="text-2xl font-heading font-bold text-text-primary">
						Event Chat
					</h1>
					<p className="mt-xs text-sm text-text-secondary">
						{eventTitle}
					</p>
				</div>

				<div className="flex min-h-0 flex-1 flex-col rounded-lg border border-border bg-surface shadow-sm">
					<div className="min-h-0 flex-1 overflow-y-auto px-lg py-lg">
						{messages.length === 0 ? (
							<div className="flex h-full items-center justify-center">
								<p className="text-sm text-text-tertiary">No messages yet.</p>
							</div>
						) : (
							<div className="space-y-md">
								{messages.map((message) => {
									const senderName =
										participantNamesById.get(message.user_id) ?? `User #${message.user_id}`;

									return (
										<div
											key={message.id}
											className="rounded-md border border-border bg-surface-dim px-md py-sm"
										>
											<div className="mb-xs flex items-center justify-between gap-md">
												<span className="text-sm font-medium text-text-primary">
													{senderName}
												</span>
												<span className="text-xs text-text-tertiary">
													{new Date(message.created_at).toLocaleString()}
												</span>
											</div>
											<p className="text-sm text-text-primary whitespace-pre-wrap break-words">
												{message.content}
											</p>
										</div>
									);
								})}
							</div>
						)}
					</div>

					<div className="border-t border-border px-lg py-md">
						<div className="flex items-end gap-md">
							<div className="flex-1">
								<Input
									value={draft}
									onChange={(e) => setDraft(e.target.value)}
									placeholder="Write a message..."
								/>
							</div>
							<Button disabled>
								Send
							</Button>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
