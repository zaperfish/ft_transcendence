export interface CreateEventRequest {
	description: string;
	duration: number;
	location_address: string;
	location_name: string;
	max_capacity: number;
	start_time: string;
	title: string;
}

export interface EventEntity {
	description: string;
	duration: number;
	location_address: string;
	location_name: string;
	max_capacity: number;
	start_time: string;
	title: string;
	created_at: string;
	id: number;
	num_registered: number;
	updated_at: string;
	// Backend returns self only when user needs to be authenticated
	// Remain optional for crash protection reason
	self?: {
		is_participant: boolean;
		role: 'admin' | 'member' | 'none';
	}
	// Backend does not return image, using interface api instead
	image?: string;
	has_image: boolean;
}

export interface GetEventRequest {
	page: number;
	page_size: number;
}

export interface GetMyEventsRequest {
	filter: 'member' | 'admin';
	page: number;
	page_size: number;
}

export interface PaginatedResponse<T> {
	data: T[];
	page: number;
	page_size: number;
	total: number;
}

export type EventsResponse = PaginatedResponse<EventEntity>;

export interface EventParticipant {
	id: number;
	name: string;
	email: string;
}

export interface EventParticipantsResponse {
	data: EventParticipant[];
}
