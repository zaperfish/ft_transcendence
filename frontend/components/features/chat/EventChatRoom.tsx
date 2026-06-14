interface EventChatRoomProps {
	eventId: number;
}

export default function EventChatRoom({ eventId }: EventChatRoomProps) {
	return (
		<div className="w-full px-xl py-2xl">
			<h1 className="text-2xl font-heading font-bold text-text-primary">
				Event Chat
			</h1>
			<p className="mt-sm text-text-secondary">
				Chat room for event #{eventId}
			</p>
		</div>
	);
}
