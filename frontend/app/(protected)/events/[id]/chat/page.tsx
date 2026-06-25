import { notFound } from "next/navigation";
import EventChatRoom from "@/components/features/chat/EventChatRoom";

export default async function EventChatPage({
	params,
}: {
	params: Promise<{ id: string }>;
}) {
	const { id } = await params;
	const eventId = Number(id);

	if (!Number.isInteger(eventId) || eventId <= 0) {
		notFound();
	}

	return <EventChatRoom eventId={eventId} />;
}
