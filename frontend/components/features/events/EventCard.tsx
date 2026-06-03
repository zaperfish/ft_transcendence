import { Button } from "@/components/ui/Button";
import { CalendarIcon, ClockIcon, MapPinIcon, UserIcon } from 'lucide-react';
import type { EventEntity } from "@/types/event";
import { useState } from "react";
import { request } from "@/lib/api/client";

interface EventCardProps {
	data: EventEntity;
}

export default function EventCard({ data }: EventCardProps) {
	const [isRegistering, setIsRegistering] = useState(false);
	const [isRegistered, setIsRegistered] = useState(false);
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
	// This api endpoint should be negotiated with backend
	const handleRegister = async () => {
		setIsRegistering(true);
		setErrorMsg(null);
		try {
			const res = await request<{ message: string }>(
				`/api/events/${data.id}/register`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
			});
			setIsRegistered(true);
		} catch (error) {
			setErrorMsg("Registration failed, please retry");
		} finally {
			setIsRegistering(false);
		}
	};

	return (
		<div className="border border-border rounded-lg overflow-hidden flex flex-col bg-surface shadow-sm hover:shadow-md transition-shadow">
			{/* Cover page: No image */}
			<div className="aspect-video bg-surface-container flex items-center justify-center"></div>
			{/* Event info */}
			<div className="p-md flex flex-col flex-1">
				<h3 className="text-xl font-semibold text-text-primary line-clamp-2 leading-snug">{data.title}</h3>
				<p className="text-sm text-text-secondary mt-xs line-clamp-3">{data.description}</p>
				<div className="mt-md space-y-sm text-text-secondary text-sm">
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
						onClick={handleRegister}
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