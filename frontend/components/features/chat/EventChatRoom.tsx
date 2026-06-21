'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { getEventById, getEventParticipants } from '@/lib/api/events';
import { buildEventChatWebSocketUrl, getEventChatMessages } from '@/lib/api/chat';
import { ApiError } from '@/lib/api/client';
import type { ChatMessage, ChatMessageInput } from '@/types/chat';

interface EventChatRoomProps {
	eventId: number;
}

type SocketStatus = 'connecting' | 'open' | 'closed' | 'error';

export default function EventChatRoom({ eventId }: EventChatRoomProps) {
	const [draft, setDraft] = useState('');
	const [socketStatus, setSocketStatus] = useState<SocketStatus>('connecting');
	const [liveMessages, setLiveMessages] = useState<ChatMessage[]>([]);
	const socketRef = useRef<WebSocket | null>(null);
	const messagesEndRef = useRef<HTMLDivElement | null>(null);

	const historyQuery = useQuery({
		queryKey: ['chat-history', eventId],
		queryFn: () => getEventChatMessages(eventId),
		retry: false,
		refetchOnWindowFocus: false,
	});

	const participantsQuery = useQuery({
		queryKey: ['event-participants', eventId],
		queryFn: () => getEventParticipants(eventId),
		retry: false,
	});

	const eventQuery = useQuery({
		queryKey: ['event', eventId],
		queryFn: () => getEventById(eventId),
		retry: false,
	});

	// REST history
	const historyData = historyQuery.data?.data;
	const participantsData = participantsQuery.data ?? [];
	const eventTitle = eventQuery.data?.title ?? `Event #${eventId}`;
	const trimmedDraft = draft.trim();
	// derive message list from history + websocket
	const chatMessages = useMemo(() => {
		const historyMessages = historyData ?? [];
		const newLiveMessages = liveMessages.filter((liveMessage) => {
			const alreadyExistsInHistory = historyMessages.some((historyMessage) => {
				return historyMessage.id === liveMessage.id;
			});

			return !alreadyExistsInHistory;
		});

		return historyMessages.concat(newLiveMessages);
	}, [historyData, liveMessages]);
	const canSendMessage = socketStatus === 'open' && trimmedDraft !== '';
	const connectionMessage = {
		connecting: 'Connecting to the chat room...',
		open: 'Connected to the chat room.',
		closed: 'Disconnected. Existing messages are still available, but sending is disabled.',
		error: 'The chat connection failed. Existing messages are still available, but sending is disabled.',
	}[socketStatus];

	const participantNamesById = useMemo(() =>
			new Map(
				(participantsData ?? []).map((participant) => [participant.id, participant.name])
			),
		[participantsData]
	);

	const historyError = historyQuery.error instanceof ApiError ? historyQuery.error : null;

	useEffect(() => {
		messagesEndRef.current?.scrollIntoView({ block: 'end' });
	}, [chatMessages]);

	useEffect(() => {
		if (!historyQuery.isSuccess) {
			return;
		}

		const socket = new WebSocket(buildEventChatWebSocketUrl(eventId));
		socketRef.current = socket;

		socket.onopen = () => {
			setSocketStatus('open');
		};

		socket.onmessage = (event) => {
			let incomingMessage: ChatMessage;
			try {
				incomingMessage = JSON.parse(event.data) as ChatMessage;
			} catch {
				return;
			}

			setLiveMessages((currentMessages) => {
				const existingMessage = currentMessages.some(
					// defensive check
					(message) => message.id === incomingMessage.id
				);
				if (existingMessage) {
					return currentMessages;
				}

				return [...currentMessages, incomingMessage];
			});
		};

		socket.onerror = () => {
			setSocketStatus('error');
		};

		socket.onclose = (event) => {
			setSocketStatus((currentStatus) =>
				currentStatus === 'error' ? 'error' : 'closed'
			);
			if (socketRef.current === socket) {
				socketRef.current = null;
			}
		};

		return () => {
			if (socketRef.current === socket) {
				socketRef.current = null;
			}
			socket.close();
		};
	}, [eventId, historyQuery.isSuccess]);

	function handleSendMessage(event: React.SubmitEvent<HTMLFormElement>) {
		event.preventDefault();

		if (!canSendMessage || socketRef.current === null) {
			return;
		}

		if (socketRef.current.readyState !== WebSocket.OPEN) {
			return;
		}

		const payload: ChatMessageInput = {
			content: trimmedDraft,
		};

		try {
			socketRef.current.send(JSON.stringify(payload));
			setDraft('');
		} catch {
			setSocketStatus('error');
		}
	}

	if (historyQuery.isLoading) {
		return (
			<div className="min-h-screen bg-surface-dim p-lg">
				<div className="flex h-[calc(100vh-2rem)] w-full items-center justify-center rounded-lg border border-border bg-surface shadow-sm">
					<p className="text-sm text-text-secondary">Loading chat...</p>
				</div>
			</div>
		);
	}

	if (historyError?.status === 403) {
		return (
			<div className="min-h-screen bg-surface-dim p-lg">
				<div className="flex h-[calc(100vh-2rem)] w-full flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
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
				<div className="flex h-[calc(100vh-2rem)] w-full flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
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
				<div className="flex h-[calc(100vh-2rem)] w-full flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
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
					<p className="mt-xs text-xs text-text-tertiary">{connectionMessage}</p>
				</div>

				<div className="flex min-h-0 flex-1 flex-col rounded-lg border border-border bg-surface shadow-sm">
					<div className="min-h-0 flex-1 overflow-y-auto px-lg py-lg">
						{chatMessages.length === 0 ? (
							<div className="flex h-full items-center justify-center">
								<p className="text-sm text-text-tertiary">No messages yet.</p>
							</div>
						) : (
							<div className="space-y-md">
								{chatMessages.map((message) => {
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
								<div ref={messagesEndRef} />
							</div>
						)}
					</div>

					<div className="border-t border-border px-lg py-md">
						<form onSubmit={handleSendMessage} className="flex items-end gap-md">
							<div className="flex-1">
								<Input
									value={draft}
									onChange={(e) => setDraft(e.target.value)}
									placeholder="Write a message..."
								/>
							</div>
							<Button type="submit" disabled={!canSendMessage}>
								Send
							</Button>
						</form>
					</div>
				</div>
			</div>
		</div>
	);
}
