import { notFound } from "next/navigation";
import EventChatRoom from "@/components/features/chat/EventChatRoom";

export default function EventChatPage({ params }: { params: { id: string } }) {
	const eventId = Number(params.id);

	if (!Number.isInteger(eventId) || eventId <= 0) {
		notFound();
	}

	return <EventChatRoom eventId={eventId} />;
}
