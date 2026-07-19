'use client';

import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getEventById, getEventParticipants, deleteEvent, removeParticipant, leaveEvent, uploadEventImage, updateEventImage } from "@/lib/api/events";
import { useAuth } from "@/lib/hooks/useAuth";
import EditEventModal from "@/components/features/events/EditEventModal";
import { useState, useRef } from "react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Avatar, AvatarFallback } from "@/components/ui/Avatar";
import { CalendarIcon, ClockIcon, MapPinIcon, UserIcon, MessageSquareIcon, EditIcon, XIcon, PencilIcon } from "lucide-react";
import { toast } from "sonner";

/**
 * EventDetailPage displays detailed information for a single event,
 * including its cover page, title, description, date, location, and participant list.
 * It provides actions for event creators (edit info, modify cover page, delete event, remove participants)
 * and for participants (unregister), with support for real-time updates via React Query.
 */
export default function EventDetailPage() {
	const params = useParams();
	const eventId = params.id as string;
	const numericId = Number(eventId);
	const router = useRouter();
	const { user, isOnline } = useAuth();
	const queryClient = useQueryClient();
	// Use ref to operate hidden file input in DOM
	const fileInputRef = useRef<HTMLInputElement>(null);

	const [isCoverUploading, setIsCoverUploading] = useState(false);
	// Logic of getting new image when updated
	const [coverRefreshKey, setCoverRefreshKey] = useState(0);
	const [isEditModalOpen, setIsEditModalOpen] = useState(false);

	const checkOffline = (action: string) => {
		if (!isOnline) {
			toast.error(`You are offline. ${action} is not available.`);
			return true;
		}
		return false;
	};

	const { data: event, isLoading: eventLoading, fetchStatus: eventFetchStatus } = useQuery({
		queryKey: ['event', numericId],
		queryFn: () => getEventById(numericId),
		networkMode: 'offlineFirst',
	});

	// Define as empty array if undefined to avoid crash
	const { data: participants = [], isLoading: participantsLoading, fetchStatus: participantsFetchStatus } = useQuery({
		queryKey: ['participants', numericId],
		queryFn: () => getEventParticipants(numericId),
		networkMode: 'offlineFirst',
	});

	const coverSrc = event?.has_image
		? `/api/events/${event.id}/image`
		: '/images/default-event-cover.jpg';

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
		if (checkOffline('Unregister')) return;
		if (confirm('Do you want to unregister?')) {
			leaveMutation.mutate();
		}
	};

	const handleRemoveParticipant = (userId: number, userName: string) => {
		if (checkOffline('Remove participants')) return;
		if (confirm(`Do you want to remove ${userName}?`)) {
			removeMutation.mutate(userId);
		}
	};

	const handleDeleteEvent = () => {
		if (checkOffline('Delete event')) return;
		if (confirm('Do you want to delete the event?')) {
			deleteMutation.mutate();
		}
	};

	const handleBackToEvents = () => {
		router.push('/events');
	};

	const handleOpenChatroom = () => {
		if (checkOffline('Open chatroom')) return;
		window.open(
			`/events/${numericId}/chat`,
			'_blank',
			'noopener,noreferrer'
		);
	};

	const handleEditClick = () => {
		if (checkOffline('Edit event')) return;
		setIsEditModalOpen(true);

	};

	if (eventLoading || participantsLoading)
		return <div className="text-center py-2xl text-text-secondary">Loading...</div>
	if (!event)
		return <div className="text-center py-2xl text-error">Failed to load events...</div>
	if ((eventFetchStatus === 'paused' || participantsFetchStatus === 'paused')
	&& !isOnline)
		return <div className="text-center py-2xl text-text-secondary">You are offline, no cached data...Please retry after you are online.</div>

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

	const handleCoverButtonClick = () => {
		if (checkOffline('Change cover')) return;
		fileInputRef.current?.click();
	};

	const handleCoverFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
		if (checkOffline('Upload image')) {
			if (fileInputRef.current) {
				fileInputRef.current.value = '';
				return;
			}
		}
		const file = e.target.files?.[0];
		if (!file) return;

		// Validate image type
		if (file.type !== "image/png") {
			toast.error('Support only image/png type');
			return;
		}
		// Validate image size
		if (file.size > 5 * 1024 * 1024) {
			toast.error('Image file cannot be more than 5MB');
			return;
		}
		setIsCoverUploading(true);
		try {
			if (event.has_image) {
				await updateEventImage(numericId, file);
			} else {
				await uploadEventImage(numericId, file);
			}
			setCoverRefreshKey(prev => prev + 1);
			// Update coverSrc if this is the first time uploading image
			queryClient.invalidateQueries({ queryKey: ['event', numericId] });
		} catch (err) {
			toast.error("Failed to update cover page, please retry");
		} finally {
			setIsCoverUploading(false);
			if (fileInputRef.current) {
				fileInputRef.current.value = '';
			}
		}
	};

	return (
		<div className="w-full px-xl pt-2 pb-8 max-w-6xl mx-auto">
			{ /* Return button */ }
			<Button variant="ghost" onClick={handleBackToEvents}>
			← Back to MyEvents
			</Button>
			{ /* Card for Event information display */ }
			<Card className="overflow-hidden p-0 gap-0">
				<div className="grid lg:grid-cols-[minmax(0,1.45fr)_minmax(320px,0.55fr)] lg:items-stretch">
					{/* Cover page */}
					<div className="aspect-video bg-surface-container flex items-center justify-center shrink-0 relative lg:aspect-auto lg:min-h-80 lg:max-h-105">
						<img
							key={coverRefreshKey}
							src={coverSrc}
							alt={`${event.title} cover`}
							className="w-full h-full object-cover"
							onError={(e) => {
							e.currentTarget.onerror = null;
							e.currentTarget.src = "/images/default-event-cover.jpg";
							}}
						/>
						{/* Button for modifying image */}
						{isCreator && (
							<>
								<Button
									type="button"
									onClick={handleCoverButtonClick}
									disabled={isCoverUploading}
									className={`absolute top-2 right-2 bg-black/50 text-white p-2 rounded-full hover:bg-black/70 transition ${!isOnline ? 'opacity-50 cursor-not-allowed' : ''}`}
									title="Change cover page"
								>
								{isCoverUploading ? (
									<span className="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
								) : (
									<PencilIcon className="size-4"/>
								)}
								</Button>
								<input
									ref={fileInputRef}
									type="file"
									accept="image/png"
									onChange={handleCoverFileChange}
									className="hidden"
								/>
							</>
						)}
					</div>
					{/* Event info */}
					<div className="p-lg flex flex-col flex-1 lg:min-h-80 lg:max-h-105 lg:overflow-y-auto">
						<div className="flex-1">
							<h3 className="text-xl font-semibold text-text-primary truncate leading-snug">{event.title}</h3>
							<p className="text-sm text-text-secondary mt-xs truncate">{event.description}</p>
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
								<MapPinIcon className="size-4 text-text-tertiary shrink-0"/>
								<span className="truncate">{event.location_name} {event.location_address}</span>
							</div>
							<div className="flex items-center gap-sm">
								<UserIcon className="size-4 text-text-tertiary"/>
								<span>{event.num_registered}/{event.max_capacity} registered</span>
							</div>
						</div>
					</div>
				</div>
			</Card>
			{ /* Card for operation enabling*/ }
			<Card className="p-lg space-y-md">
				{ /* Open chatroom */ }
				<Button
					className={`w-full bg-primary text-white ${!isOnline ? 'opacity-50 cursor-not-allowed' : ''}`}
					onClick={handleOpenChatroom}
				>
					<MessageSquareIcon className="size-4 mr-2" />
					Open Event Chatroom
				</Button>
				{ /* Admin-only operations */ }
				{isCreator && (
					<>
						{ /* Edit event information */ }
						<Button
							onClick={handleEditClick}
							variant="outline"
							className={`w-full ${!isOnline ? 'opacity-50 cursor-not-allowed' : ''}`}
						>
							<EditIcon className="size-4 mr-2" />
							Modify Event Information
						</Button>
						{ /* Delete event */ }
						<Button
							onClick={handleDeleteEvent}
							variant="outline"
							className={`w-full text-error border-error hover:bg-error/10 ${!isOnline ? 'opacity-50 cursor-not-allowed' : ''}`}
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
						className={`w-full text-error border-error hover:bg-error/10 ${!isOnline ? 'opacity-50 cursor-not-allowed' : ''}`}
						disabled={leaveMutation.isPending}
					>
						{leaveMutation.isPending ? 'unregistering..' : 'Unregister'}
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
				<div className="flex flex-col gap-md">
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
							{isCreator && participant.id !== user?.id && (
								<Button
									onClick={() => handleRemoveParticipant(participant.id, participant.name)}
									className={`text-error bg-white hover:bg-error/10 p-1 rounded-md transition-colors ${!isOnline ? 'opacity-50 cursor-not-allowed' : ''}`}
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
						queryClient.invalidateQueries({ queryKey: ['event', numericId] });
						setIsEditModalOpen(false);
					}}
				/>
			)}
		</div>
	);
}
