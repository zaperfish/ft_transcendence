export interface ChatMessage {
	id: number;
	event_id: number;
	user_id: number;
	sender_name: string;
	content: string;
	created_at: string;
}

export interface ChatHistoryResponse {
	data: ChatMessage[];
}

export interface ChatMessageInput {
	content: string;
}
