'use client';
import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getEventById, getEventParticipants, updateEvent, deleteEvent, removeParticipant, leaveEvent } from "@/lib/api/events";
import { useAuth } from "@/lib/hooks/useAuth";
import type { User } from "@/types/user";
import EditEventModal from "@/components/features/events/EditEventModal";
import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Avatar, AvatarFallback } from "@/components/ui/Avatar";
import { CalendarIcon, ClockIcon, MapPinIcon, UserIcon, MessageSquareIcon, EditIcon, XIcon } from "lucide-react";

export default function EventDetailPage() {
	const params = useParams();
	const eventId = params.id as string;
	const numericId = Number(eventId);
	const router = useRouter();
	const { user } = useAuth();
	const queryClient = useQueryClient();

	const [isEditModalOpen, setIsEditModalOpen] = useState(false);

	const { data: event, isLoading: eventLoading } = useQuery({
		queryKey: ['event', numericId],
		queryFn: () => getEventById(numericId),
	});

	// Define as empty array if undefined to avoid crash
	const { data: participants = [], isLoading: participantsLoading } = useQuery({
		queryKey: ['participants', numericId],
		queryFn: () => getEventParticipants(numericId),
	});

	const isCreator = event?.self?.role === 'admin';
	const isParticipant = event?.self?.role === 'member';

	const leaveMutation = useMutation({
		mutationFn: () => leaveEvent(numericId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['event', numericId] });
			queryClient.invalidateQueries({ queryKey: ['participants', numericId] });
		},
	});

	const removeMutation = useMutation({
		mutationFn: (userId: number) =>  removeParticipant(numericId, userId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['event', numericId] });
			queryClient.invalidateQueries({ queryKey: ['participants', numericId] });
		},
	});

	const deleteMutation = useMutation({
		mutationFn: () => deleteEvent(numericId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['event'] });
			router.push('/events');
		}
	});

	const handleLeaveEvent = () => {
		if (confirm('Do you want to unregister?')) {
			leaveMutation.mutate();
		}
	};

	const handleRemoveParticipant = (userId: number, userName: string) => {
		if (confirm(`Do you want to remove ${userName}?`)) {
			removeMutation.mutate(userId);
		}
	};

	const handleDeleteEvent = () => {
		if (confirm('Do you want to delete the event?')) {
			deleteMutation.mutate();
		}
	};

	if (eventLoading || participantsLoading)
		return <div className="text-center py-2xl text-text-secondary">Loading...</div>
	if (!event)
		return <div className="text-center py-2xl text-error">Failed to load events...</div>

	const eventDate = new Date(event.start_time);
	const dateStr = eventDate.toLocaleDateString("en-US", {
		month: "short",
		day: "numeric",
		year: "numeric",
	});
	const timeStr = eventDate.toLocaleTimeString("en-US", {
		hour: "2-digit",
		minute: "2-digit",
	});

	return (
		<div className="w-full px-xl py-2xl">
			{ /* Return button */ }
			<Button variant="ghost" onClick={() => router.back()}>
			← Back to MyEvents
			</Button>
			{ /* Card for Event information display*/ }
			<Card className="overflow-hidden">
				{/* Cover page: No image */}
				<div className="aspect-video bg-surface-container flex items-center justify-center shrink-0"></div>
				{/* Event info */}
				<div className="p-md flex flex-col flex-1">
					<div className="flex-1">
						<h3 className="text-xl font-semibold text-text-primary line-clamp-2 leading-snug">{event.title}</h3>
						<p className="text-sm text-text-secondary mt-xs overflow-y-auto max-h-18 pr-1">{event.description}</p>
					</div>
					<div className="mt-md space-y-sm text-text-secondary text-sm shrink-0">
						<div className="flex items-center gap-sm">
							<CalendarIcon className="size-4 text-text-tertiary"/>
							<span>{dateStr}</span>
						</div>
						<div className="flex items-center gap-sm">
							<ClockIcon className="size-4 text-text-tertiary"/>
							<span>{timeStr} ({event.duration} min)</span>
						</div>
						<div className="flex items-center gap-sm">
							<MapPinIcon className="size-4 text-text-tertiary"/>
							<span>{event.location_name} {event.location_address}</span>
						</div>
						<div className="flex items-center gap-sm">
							<UserIcon className="size-4 text-text-tertiary"/>
							<span>{event.num_registered}/{event.max_capacity} registered</span>
						</div>
					</div>
				</div>
			</Card>
			{ /* Card for operation enabling*/ }
			<Card className="p-lg space-y-md">
				{ /* Open chatroom */ }
				<Button className="w-full bg-primary text-white">
					<MessageSquareIcon className="size-4 mr-2" />
					Open Event Chatroom
				</Button>
				{ /* Admin-only operations */ }
				{isCreator && (
					<>
						{ /* Edit event information */ }
						<Button
							onClick={() => setIsEditModalOpen(true)}
							variant="outline"
							className="w-full"
						>
							<EditIcon className="size-4 mr-2" />
							Modify Event Information
						</Button>
						{ /* Delete event */ }
						<Button
							onClick={handleDeleteEvent}
							variant="outline"
							className="w-full text-error border-error hover:bg-error/10"
							disabled={deleteMutation.isPending}
						>
							{deleteMutation.isPending ? 'deleting..' : 'Delete Event'}
						</Button>
					</>
				)}
				{ /* Leave event */ }
				{isParticipant && !isCreator && (
					<Button
						onClick={handleLeaveEvent}
						variant="outline"
						className="w-full text-error border-error hover:bg-error/10"
						disabled={leaveMutation.isPending}
					>
						{leaveMutation.isPending ? 'unregistering..' : 'unregistered'}
					</Button>
				)}
			</Card>
			{ /* Card for participant management*/ }
			<Card className="p-lg space-y-md">
				{ /* Description of participants */ }
				<div className="flex justify-between items-center mb-md">
					<h2 className="text-xl font-heading font-semibold text-text-primary">Participants</h2>
					<span className="text-text-secondary">{event.num_registered}/{event.max_capacity}</span>
				</div>
				{ /* Display participants */ }
				<div className="grid grid-cols-2 md:grid-cols-4 gap-md">
					{participants.map((participant) => (
						<div key={participant.id} className="flex items-center gap-sm">
							{ /* Display participant's avatar*/ }
							<Avatar className="size-9">
								<AvatarFallback>
									{participant.name?.charAt(0)?.toUpperCase() || "U"}
		  						</AvatarFallback>
							</Avatar>
							{ /* Display participant's name*/ }
							<div className="flex-1 min-w-0">
								<p className="text-sm font-medium text-text-primary truncate">{participant.name}</p>
							</div>
							{ /* Remove participants*/ }
							{isCreator && (
								<Button
									onClick={() => handleRemoveParticipant(participant.id, participant.name)}
									className="text-error bg-white hover:bg-error/10 p-1 rounded-md transition-colors"
								>
									<XIcon className="size-4" />
								</Button>
							)}
						</div>
					))}
				</div>
			</Card>
			{ /* Edit modal */ }
			{isEditModalOpen && (
				<EditEventModal
					event={event}
					onClose={() => setIsEditModalOpen(false)}
					onSuccess={() => {
						queryClient.invalidateQueries({ queryKey: ['event', eventId] });
						setIsEditModalOpen(false);
					}}
				/>
			)}
		</div>
	);
}