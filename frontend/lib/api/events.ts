import type { CreateEventRequest, EventEntity, EventsResponse, GetEventRequest, GetMyEventsRequest } from "@/types/event";
import type { User } from "@/types/user";
import { request, ApiError } from '@/lib/api/client';

// API functions used for home page
// Get list of events
export async function getEvents({
	page = 1,
	page_size = 10,
}: GetEventRequest) : Promise<EventsResponse> {
	const params = new URLSearchParams({
		page: page.toString(),
		page_size: page_size.toString(),
	});
	return request<EventsResponse>(`/api/events?${params.toString()}`);
}

// Create a new event
export async function createEvent(data: CreateEventRequest) : Promise<EventEntity> {
	return request<EventEntity>("/api/events", {
		method: "POST",
		headers: { "Content-type": "application/json" },
		body: JSON.stringify(data),
	});
}

// Join a event
export async function joinEvent(eventId: number): Promise<void> {
	return request<void>(`/api/me/events/${eventId}/join`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
	});
}

// API functions used for event list page
// Dont have to catch errors because React Query will handle automatically
export async function getMyEvents({
	filter,
	page = 1,
	page_size = 10,
}: GetMyEventsRequest): Promise<EventsResponse> {
	const params = new URLSearchParams({
		filter,
		page: page.toString(),
		page_size: page_size.toString(),
	});
	return request<EventsResponse>(`/api/events?${params.toString()}`);
}

// API functions used for event detail page
// Get single event by Id
export async function getEventById(id: number): Promise<EventEntity> {
	return request<EventEntity>(`/api/events/${id}`);
}

// Get list of participants of an event
export async function getEventParticipants(id: number): Promise<User[]> {
	const response = await request<{ data: User[] }>(
		`/api/events/${id}/participants`
	);
	return response.data;
}

// Update event information (Admin-only)
export async function updateEvent(
	id: number,
	data: Partial<CreateEventRequest>
): Promise<EventEntity> {
	return request<EventEntity>(`/api/events/${id}`, {
		method: "PATCH",
		headers: { "Content-type": "application/json" },
		body: JSON.stringify(data),
	});
}

// Delete event (Admin-only)
export async function deleteEvent(id: number): Promise<void> {
	return request<void>(`/api/events/${id}`, {
		method: "DELETE",
	});
}

// Remove participants (Admin-only)
export async function removeParticipant(
	eventId: number,
	userId: number
): Promise<void> {
	return request<void>(`/api/events/${eventId}/participants/${userId}`, {
		method: "DELETE",
	});
}

// User unregisters event
export async function leaveEvent(id: number): Promise<void> {
	return request<void>(`/api/me/events/${id}/leave`, {
		method: "DELETE",
	});
}

// User upload image when creating event
// Not using request() because it doesnt return json
export async function uploadEventImage(eventId: number, file: File): Promise<void> {
	const response = await fetch(`/api/events/${eventId}/image`, {
		method: 'POST',
		headers: {
			'Content-Type': file.type,
		},
		body: file,
		credentials: 'include',
	});

	if (!response.ok) {
		if (response.status === 401 && typeof window !== 'undefined' && window.location.pathname !== '/login' && window.location.pathname !== '/register') {
			window.location.href = '/login';
		}
		throw new ApiError(response.status, 'Image upload failed');
	}
	return;
}

// User update image when event is created already
export async function updateEventImage(eventId: number, file: File): Promise<void> {
	const response = await fetch(`/api/events/${eventId}/image`, {
		method: 'PATCH',
		headers: {
			'Content-Type': file.type,
		},
		body: file,
		credentials: 'include',
	});

	if (!response.ok) {
		if (response.status === 401 && typeof window !== 'undefined' && window.location.pathname !== '/login' && window.location.pathname !== '/register') {
			window.location.href = '/login';
		}
		throw new ApiError(response.status, 'Image update failed');
	}
	return;
}