import type { CreateEventRequest, EventEntity, EventsResponse, GetEventRequest, GetMyEventsRequest } from "@/types/event";
import type { User } from "@/types/user";
import type { CreateEventRequest, EventEntity, EventsResponse, EventParticipantsResponse } from "@/types/event";
import { request } from '@/lib/api/client';

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
/**
 * Fetch one event by ID.
 *
 * This is used by event-specific views, such as the chat window header,
 * when the frontend needs the event title and basic event metadata.
 */
export async function getEvent(eventId: number): Promise<EventEntity> {
	return request<EventEntity>(`/api/events/${eventId}`);
}

/**
 * Fetch the participants registered for one event.
 *
 * The chat frontend uses this list to resolve message sender IDs into
 * participant names when rendering persisted history and live messages.
 */
export async function getEventParticipants(eventId: number): Promise<EventParticipantsResponse> {
	return request<EventParticipantsResponse>(`/api/events/${eventId}/participants`);
}
