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
	id: string;
	num_registered: number;
	updated_at: string;
	self: {
		is_participant: boolean;
	};
}

export interface GetEventRequest {
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