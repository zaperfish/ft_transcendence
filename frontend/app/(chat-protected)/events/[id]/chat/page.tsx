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

	return (
		<div className="flex h-full min-h-0 flex-1 flex-col overflow-hidden">
			<EventChatRoom eventId={eventId} />
		</div>
	);
}
