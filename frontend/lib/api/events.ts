import type { EventEntity } from "@/types/event";
import { request } from '@/lib/api/client';

// These two endpoints should be negotiated with backend
// Dont have to catch errors because React Query will handle automatically

export async function getAttendingEvents(): Promise<EventEntity[]> {
	return request<EventEntity[]>('/api/events/attending');
}

export async function getHostingEvents(): Promise<EventEntity[]> {
	return request<EventEntity[]>('/api/events/hosting');
}