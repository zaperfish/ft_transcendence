'use client';

import { Button } from "@/components/ui/Button";
import { CalendarIcon, ClockIcon, MapPinIcon, UserIcon } from 'lucide-react';
import type { EventEntity } from "@/types/event";
import { useState } from "react";
import { joinEvent } from "@/lib/api/events";
import { useQueryClient } from "@tanstack/react-query";
import { useTheme } from "@/lib/context/ThemeContext";
import { cn } from "@/lib/utils";

const DEFAULT_IMAGE = '/images/default-event-cover.jpg';

interface EventCardProps {
	data: EventEntity;
	mode?: 'register' | 'detail';
	onDetail?: () => void;
	layout?: 'vertical' | 'horizontal';
	disabled?: boolean;
}

/**
 * EventCard component used for both 'register' and 'detail' mode as well as 'vertical' and 'horizontal' display.
 *
 * It displays a single event's details (cover page, title, description, date, time, location, capacity)
 * and in 'register' mode provides a registration button with optimistic UI updates and error handling
 * and in 'detail' mode provides a button redirecting user to event detail page;
 * and 'vertical' display is defaultly used in homepage, 'horizontal' display used in myEvents page.
 */
export default function EventCard({
	data,
	mode = 'register',
	onDetail,
	layout = 'vertical',
	disabled = false,
 }: EventCardProps) {
	const [isRegistering, setIsRegistering] = useState(false);
	// Avoid crash when backend returns undefined self
	const isRegistered = data.self?.is_participant ?? false;
	const isFull = data.num_registered >= data.max_capacity;
	const [errorMsg, setErrorMsg] = useState<string | null>(null);
	const queryClient = useQueryClient();
	const { theme } = useTheme();
	const isClassic = theme === "classic";
	const isButtonDisabled = disabled || isRegistering || isRegistered || (!isRegistered && isFull);

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

	const handleRegister = async () => {
		if (mode !== 'register' || isRegistered || isRegistering) return;
		setIsRegistering(true);
		setErrorMsg(null);
		queryClient.setQueryData(["events"], (oldData: any) => {
			if (!oldData) return oldData;
			const newPages = oldData.pages.map((page: any) => ({
				...page,
				data: page.data.map((event: any) =>
					event.id === data.id
						? { ...event, num_registered: event.num_registered + 1, self: {...(event.self || {}), is_participant: true}, }
						: event
				),
			}));
			return { ...oldData, pages: newPages };
		});
		try {
			await joinEvent(data.id);
			queryClient.invalidateQueries({ queryKey: ["events"] });
		} catch (error) {
			setErrorMsg("Registration failed, please retry");
			queryClient.invalidateQueries({ queryKey: ["events"] });
		} finally {
			setIsRegistering(false);
		}
	};

	const isDetailMode = mode === 'detail';
	const isHorizontal = layout === 'horizontal';

	return (
		<div className={cn(
			'flex flex-col overflow-hidden rounded-lg border transition-shadow',
			isClassic
				? 'border-border bg-surface shadow-sm hover:shadow-md'
				: 'border-teal-900/10 bg-white shadow-md hover:shadow-lg',
			isHorizontal && 'sm:flex-row',
			disabled && 'pointer-events-none opacity-60',
		)}>
			{/* Cover page */}
			<div className={cn(
				'shrink-0 overflow-hidden',
				isClassic ? 'bg-surface-container' : 'bg-slate-100',
				isHorizontal
					? 'w-full sm:w-60 h-60 sm:h-auto'
					: 'aspect-video max-h-40 w-full',
			)}>
				<img
					src={data.image}
					alt={data.title}
					className="w-full h-full object-cover"
					onError={(e) => {
						e.currentTarget.onerror = null;
						e.currentTarget.src = DEFAULT_IMAGE;
					}}
				/>
			</div>
			{/* Event info */}
			<div className="p-md flex flex-col flex-1 min-w-0">
				<div className="min-h-0 overflow-hidden">
					<h3 className="text-xl font-semibold text-text-primary truncate leading-snug">{data.title}</h3>
					<p className="text-sm text-text-secondary mt-xs truncate">{data.description}</p>
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
					<div className="flex items-center gap-sm min-w-0">
						<MapPinIcon className="size-4 text-text-tertiary shrink-0"/>
						<span className="truncate">{data.location_name} {data.location_address}</span>
					</div>
					<div className="flex items-center gap-sm">
						<UserIcon className="size-4 text-text-tertiary"/>
						<span>{data.num_registered}/{data.max_capacity} registered</span>
					</div>
				</div>
				<div className="mt-auto pt-md w-full">
					{isDetailMode ? (
						<Button
							onClick={onDetail}
							disabled={disabled}
							variant={isClassic ? "outline" : "default"}
							className={cn(
								"w-full",
								isClassic
									? "hover:bg-primary/10 hover:text-primary hover:border-primary transition-colors"
									: "bg-primary text-primary-foreground hover:bg-primary-dim",
							)}
						>
							View Detail
						</Button>
					) : (
						<Button
							onClick={handleRegister}
							disabled={isButtonDisabled}
							className={`w-full ${
								isRegistered
									? "bg-success text-white cursor-not-allowed"
									: isRegistering
										? "bg-text-tertiary text-white cursor-wait"
										: isFull
											? "bg-gray-300 text-gray-500 cursor-not-allowed"
											: "bg-primary text-primary-foreground hover:bg-primary-dim"
							}`}
						>
							{ isRegistered
								? "Registered"
								: isRegistering
									? "Registering..."
									: isFull
										? "Full"
										: "Register"}
						</Button>
					)}
					{errorMsg && !isDetailMode && <p className="text-error text-xs mt-xs"></p>}
				</div>
			</div>
		</div>
	);
}