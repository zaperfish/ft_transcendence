'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/Avatar';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { getEventById, getEventParticipants } from '@/lib/api/events';
import { useAuth } from '@/lib/hooks/useAuth';
import { buildEventChatWebSocketUrl, getEventChatMessages } from '@/lib/api/chat';
import { ApiError } from '@/lib/api/client';
import type { ChatMessage, ChatMessageInput } from '@/types/chat';

interface EventChatRoomProps {
	eventId: number;
}

type SocketStatus = 'connecting' | 'open' | 'closed' | 'error';

export default function EventChatRoom({ eventId }: EventChatRoomProps) {
	const { user } = useAuth();
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
	const currentUserID = user?.id;
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

	const participantNamesById = useMemo(
		() =>
			new Map(
				(participantsData ?? []).map((participant) => [participant.id, participant.name])
			),
		[participantsData]
	);
	const participantsById = useMemo(
		() => new Map(participantsData.map((participant) => [participant.id, participant])),
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

		socket.onclose = () => {
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

	function getAvatarFallback(name: string) {
		return name.charAt(0).toUpperCase() || 'U';
	}

	function truncateSenderName(name: string) {
		if (name.length <= 10) {
			return name;
		}

		return `${name.slice(0, 10)}...`;
	}

	if (historyQuery.isLoading) {
		return (
			<div className="flex min-h-[calc(100vh-12rem)] w-full items-center justify-center rounded-lg border border-border bg-surface shadow-sm">
				<p className="text-sm text-text-secondary">Loading chat...</p>
			</div>
		);
	}

	if (historyError?.status === 403) {
		return (
			<div className="flex min-h-[calc(100vh-12rem)] w-full flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
				<h1 className="text-2xl font-heading font-bold text-text-primary">
					Access denied
				</h1>
				<p className="mt-sm text-sm text-text-secondary">
					You must be a participant in this event to access the chat room.
				</p>
			</div>
		);
	}

	if (historyError?.status === 404) {
		return (
			<div className="flex min-h-[calc(100vh-12rem)] w-full flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
				<h1 className="text-2xl font-heading font-bold text-text-primary">
					Event not found
				</h1>
				<p className="mt-sm text-sm text-text-secondary">
					This event does not exist or is no longer available.
				</p>
			</div>
		);
	}

	if (historyQuery.isError) {
		return (
			<div className="flex min-h-[calc(100vh-12rem)] w-full flex-col items-center justify-center rounded-lg border border-border bg-surface px-xl text-center shadow-sm">
				<h1 className="text-2xl font-heading font-bold text-text-primary">
					Failed to load chat
				</h1>
				<p className="mt-sm text-sm text-text-secondary">
					An unexpected error occurred while loading chat history.
				</p>
			</div>
		);
	}

	return (
		<div className="flex min-h-[calc(100vh-12rem)] w-full flex-col">
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
								const isCurrentUserMessage = currentUserID === message.user_id;
								const participant = participantsById.get(message.user_id);
								const senderName =
									(isCurrentUserMessage ? user?.name : participant?.name) ??
									participantNamesById.get(message.user_id) ??
									`User #${message.user_id}`;
								const senderAvatar = isCurrentUserMessage
									? user?.avatar
									: participant?.avatar;

								return (
									<div
										key={message.id}
										className={isCurrentUserMessage ? 'flex justify-end' : 'flex justify-start'}
									>
										<div
											className={
												isCurrentUserMessage
													? 'flex w-[75%] flex-row-reverse items-start gap-sm'
													: 'flex w-[75%] items-start gap-sm'
											}
										>
											<Avatar size="sm">
												<AvatarImage src={senderAvatar} alt={senderName} />
												<AvatarFallback
													className={
														isCurrentUserMessage
															? 'bg-primary/15 font-bold text-primary'
															: 'bg-muted font-bold text-text-secondary'
													}
												>
													{getAvatarFallback(senderName)}
												</AvatarFallback>
											</Avatar>
											<div
												className={
													isCurrentUserMessage
														? 'flex-1 rounded-md border border-primary/20 bg-primary px-md py-sm text-primary-foreground'
														: 'flex-1 rounded-md border border-border bg-surface-dim px-md py-sm text-text-primary'
												}
											>
												<div className="mb-xs flex items-center justify-between gap-md">
													<span
														className={
															isCurrentUserMessage
																? 'text-sm font-bold text-primary-foreground/50'
																: 'text-sm font-bold text-text-tertiary/70'
														}
													>
														{truncateSenderName(senderName)}
													</span>
													<span
														className={
															isCurrentUserMessage
																? 'text-xs text-primary-foreground/50'
																: 'text-xs text-text-tertiary/70'
														}
													>
														{new Date(message.created_at).toLocaleString()}
													</span>
												</div>
												<p
													className={
														isCurrentUserMessage
															? 'text-sm whitespace-pre-wrap break-words text-primary-foreground'
															: 'text-sm whitespace-pre-wrap break-words text-text-primary'
													}
												>
													{message.content}
												</p>
											</div>
										</div>
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
	);
}
