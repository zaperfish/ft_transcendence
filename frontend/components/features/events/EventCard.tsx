import { Button } from "@/components/ui/Button";
import { CalendarIcon, ClockIcon, MapPinIcon, UserIcon } from 'lucide-react';
import type { EventEntity } from "@/types/event";
import { useState } from "react";
import { request } from "@/lib/api/client";
import { useQueryClient } from "@tanstack/react-query";

interface EventCardProps {
	data: EventEntity;
}

/**
 * EventCard component displays a single event's details (title, description, date, time, location, capacity)
 * and provides a registration button with optimistic UI updates and error handling.
 */
export default function EventCard({ data }: EventCardProps) {
	const [isRegistering, setIsRegistering] = useState(false);
	const isRegistered = data.self.is_participant;
	const [errorMsg, setErrorMsg] = useState<string | null>(null);

	const eventDate = new Date(data.start_time);
	const dateStr = eventDate.toLocaleDateString("en-US", {
		month: "short",
		day: "numeric",
		year: "numeric",
	});
	const timeStr = eventDate.toLocaleTimeString("en-US", {
		hour: "2-digit",
		minute: "2-digit",
	});

	const queryClient = useQueryClient();
	// This api endpoint should be negotiated with backend
	const handleRegister = async () => {
		if (isRegistered || isRegistering) return;
		setIsRegistering(true);
		setErrorMsg(null);
		queryClient.setQueryData(["events"], (oldData: any) => {
			if (!oldData) return oldData;
			const newPages = oldData.pages.map((page: any) => ({
				...page,
				data: page.data.map((event: any) =>
					event.id === data.id
						? { ...event, num_registered: event.num_registered + 1, self: {...event.self, is_participant: true}, }
						: event
				),
			}));
			return { ...oldData, pages: newPages };
		});
		try {
			const res = await request<{ message: string }>(
				`/api/me/join/${data.id}`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
			});
			queryClient.invalidateQueries({ queryKey: ["events"] });
		} catch (error) {
			setErrorMsg("Registration failed, please retry");
			queryClient.invalidateQueries({ queryKey: ["events"] });
		} finally {
			setIsRegistering(false);
		}
	};

	return (
		<div className="border border-border rounded-lg overflow-hidden flex flex-col bg-surface shadow-sm hover:shadow-md transition-shadow">
			{/* Cover page: No image */}
			<div className="aspect-video bg-surface-container flex items-center justify-center shrink-0"></div>
			{/* Event info */}
			<div className="p-md flex flex-col flex-1">
				<div className="flex-1">
					<h3 className="text-xl font-semibold text-text-primary line-clamp-2 leading-snug">{data.title}</h3>
					<p className="text-sm text-text-secondary mt-xs overflow-y-auto max-h-18 pr-1">{data.description}</p>
				</div>
				<div className="mt-md space-y-sm text-text-secondary text-sm shrink-0">
					<div className="flex items-center gap-sm">
						<CalendarIcon className="size-4 text-text-tertiary"/>
						<span>{dateStr}</span>
					</div>
					<div className="flex items-center gap-sm">
						<ClockIcon className="size-4 text-text-tertiary"/>
						<span>{timeStr} ({data.duration} min)</span>
					</div>
					<div className="flex items-center gap-sm">
						<MapPinIcon className="size-4 text-text-tertiary"/>
						<span>{data.location_name} {data.location_address}</span>
					</div>
					<div className="flex items-center gap-sm">
						<UserIcon className="size-4 text-text-tertiary"/>
						<span>{data.num_registered}/{data.max_capacity} registered</span>
					</div>
				</div>
				<div className="mt-auto pt-md">
					<Button
						onClick={() => handleRegister}
						disabled={isRegistering || isRegistered}
						className={`w-full ${
							isRegistered
								? "bg-success text-white cursor-not-allowed"
								: isRegistering
									? "bg-text-tertiary text-white cursor-wait"
									: "bg-primary text-primary-foreground hover:bg-primary-dim"
						}`}
					>
						{ isRegistered ? "Registered" : isRegistering ? "Registering..." : "Register"}
					</Button>
					{errorMsg && <p className="text-error text-xs mt-xs"></p>}
				</div>
			</div>
		</div>
	);
}