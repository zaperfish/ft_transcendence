import { request } from "@/lib/api/client";
import type { ChatHistoryResponse } from "@/types/chat";

export const maxChatMessageCharacters = 2000;

export async function getEventChatMessages(eventId: number): Promise<ChatHistoryResponse> {
	return request<ChatHistoryResponse>(`/api/events/${eventId}/chat/messages`);
}

export function buildEventChatWebSocketUrl(eventId: number): string {
	const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
	return `${protocol}//${window.location.host}/api/events/${eventId}/chat/ws`;
}
