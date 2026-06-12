import type { CreateEventRequest, EventEntity, EventsResponse } from "@/types/event";
import type { User } from "@/types/user";
import { request } from '@/lib/api/client';

// API functions used for home page
// Get list of events
export async function getEvents(
	page: number = 1,
	page_size: number = 10
) : Promise<EventsResponse> {
	return request<EventsResponse>(
		`/api/events?page=${page}&page_size=${page_size}`
	);
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
	return request<void>(`/api/me/join/${eventId}`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
	});
}

// API functions used for event list page
// These two endpoints should be negotiated with backend
// Dont have to catch errors because React Query will handle automatically
export async function getAttendingEvents(): Promise<EventEntity[]> {
	return request<EventEntity[]>('/api/events/attending');
}

export async function getHostingEvents(): Promise<EventEntity[]> {
	return request<EventEntity[]>('/api/events/hosting');
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
	return request<void>(`/api/me/leave/${id}`, {
		method: "DELETE",
	});
}