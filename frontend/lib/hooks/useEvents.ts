import { useQuery } from "@tanstack/react-query";
import { getAttendingEvents, getHostingEvents } from "@/lib/api/events";

export function useEvents(tab: 'attending' | 'hosting') {
	return useQuery({
		queryKey: ['events', tab],
		queryFn: () => (tab === 'attending' ? getAttendingEvents() : getHostingEvents()),
	});
}