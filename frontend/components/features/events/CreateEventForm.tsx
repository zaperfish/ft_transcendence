import type { CreateEventRequest } from '@/types/event';
import { useForm } from 'react-hook-form';
import { useState, useCallback } from 'react'
import { ApiError } from '@/lib/api/client';
import { createEvent, uploadEventImage } from '@/lib/api/events';
import { Button } from '@/components/ui/Button';
import { ImageUpload } from '@/components/ui/ImageUpload';
import { FormLabel } from '@/components/ui/FormLabel';

interface CreateEventFormProps {
	open: boolean;
	onClose: () => void;
	onSuccess?: (warning?: string) => void;
}

/**
 * CreateEventForm renders a modal form for creating a new event.
 * It includes validation for all fields (title, description, start time, duration, location, capacity),
 * handles date formatting to RFC 3339, submits the event to the API,
 * and manages server errors and loading states.
 */
export default function CreateEventForm({ open, onClose, onSuccess }: CreateEventFormProps) {
	const { register, handleSubmit, reset, formState: { errors, isSubmitting }, } = useForm<CreateEventRequest>({
		defaultValues: {
			description: "",
			duration: 15,
			location_address: "",
			location_name: "",
			max_capacity: 1,
			start_time: "",
			title: "",
		},
	});
	const [serverError, setServerError] = useState<string | null>(null);
	const [imageFile, setImageFile] = useState<File | null>(null);
	const [imageError, setImageError] = useState<string>('');

	// Stabilize reference and avoid recreating incline function when component renders
	const handleImageChange = useCallback((file: File | null) => {
		setImageFile(file);
		setImageError('');
	}, []);

	if (!open)
		return null;

	const onSubmit = async (data: CreateEventRequest) => {
		setServerError(null);

		// change datetime-local to RFC 3339 format, e.g.2026-11-20T10:05:00.000Z
		const formattedData = {
			...data,
			start_time: new Date(data.start_time).toISOString(),
		};

		try {
			const createdEvent = await createEvent(formattedData);
			const eventId = createdEvent.id;
			let warning: string | undefined;

			if (imageFile) {
				try {
					await uploadEventImage(eventId, imageFile);
				} catch (uploadErr) {
					console.warn('Failed to upload image', uploadErr);
					warning = "Failed to upload cover page, but you can upload later in detail page."
				}
			}
			onSuccess?.(warning);
			onClose();
			reset();
			setImageFile(null);
			setImageError('');
		} catch (err) {
			if (err instanceof ApiError) {
				setServerError(`Request failed (code ${err.status}): ${err.message}`);
			} else if (err instanceof Error) {
				setServerError(err.message);
			} else {
				setServerError("Unknown error");
			}
		}
	};

	return (
			<div className="fixed inset-0 z-50 overflow-y-auto bg-background/80 backdrop-blur-sm">
				<div className="flex min-h-full items-start justify-center px-md py-8">
					<div className='bg-surface rounded-lg shadow-lg max-w-175 w-full mx-auto max-h-[calc(100vh-4rem)] overflow-y-auto'>
						<div className='p-2xl'>
					<h2 className='text-2xl font-heading font-bold text-text-primary mb-xl'>Create a new event</h2>
					<form onSubmit={handleSubmit(onSubmit)} className='space-y-lg'>
					{/* Title - required, 3 - 100 characters */}
						<div className='space-y-sm'>
							<FormLabel htmlFor='title-input'>Title</FormLabel>
							<input
								id="title-input"
								{...register("title", {
									required: "Please enter event title",
									minLength: {
										value: 3,
										message: "Title should be no less than 3 characters",
									},
									maxLength: {
										value: 100,
										message: "Title should be no more than 100 characters",
									},
								})}
								className='w-full px-md py-sm border border-border rounded-md bg-surface text-text-primary placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-colors'
								placeholder='Please enter event title (3-100 characters)'
							/>
							{errors.title && <p className='text-error text-sm'>{errors.title.message}</p>}
						</div>
					{/* Description - required, 10 - 500 characters */}
						<div className='space-y-sm'>
							<FormLabel htmlFor='Description-input'>Description</FormLabel>
							<textarea
								id="Description-input"
								{...register("description", {
									required: "Please enter event description",
									minLength: {
										value: 10,
										message: "Description should be no less than 10 characters",
									},
									maxLength: {
										value: 500,
										message: "Description should be no more than 500 characters",
									},
								})}
								rows={4}
								className='w-full px-md py-sm border border-border rounded-md bg-surface text-text-primary placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-colors resize-none'
								placeholder='Please enter event description (10 - 500 characters)'
							/>
							{errors.description && <p className='text-error text-sm'>{errors.description.message}</p>}
						</div>
					{/* Start time - required, datetime-local format */}
						<div className='space-y-sm'>
							<FormLabel htmlFor='start-time-input'>Start time</FormLabel>
							<input
								type="datetime-local"
								id='start-time-input'
								{...register("start_time", {
									required: "Please enter start time of event",
									validate: (value) => {
										const selectedDate = new Date(value);
										const now = new Date();
										if (selectedDate < now) {
											return "Start time cannot be earlier than now";
										}
										return true;
									},
								})}
								className='w-full px-md py-sm border border-border rounded-md bg-surface text-text-primary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-colors'
							/>
							{errors.start_time && <p className='text-error text-sm'>{errors.start_time.message}</p>}
						</div>
					{/* Duration - required, 15 - 480 mins */}
						<div className='space-y-sm'>
							<FormLabel htmlFor='duration-input'>Duration (min)</FormLabel>
							<input
								type="number"
								id='duration-input'
								{...register("duration", {
									required: "Please enter duration of event",
									min: {
										value: 15,
										message: "Duration should be no less than 15 minutes",
									},
									max: {
										value: 480,
										message: "Duration should be no more than 480 minutes",
									},
									valueAsNumber: true,
								})}
								className='w-full px-md py-sm border border-border rounded-md bg-surface text-text-primary placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-colors'
								placeholder='Please enter duration (15 - 480 mins)'
							/>
							{errors.duration && <p className='text-error text-sm'>{errors.duration.message}</p>}
						</div>
					{/* Location name - required, 3 - 100 characters */}
						<div className='space-y-sm'>
							<FormLabel htmlFor='location-name-input'>Location name</FormLabel>
							<input
								id='location-name-input'
								{...register("location_name", {
									required: "Please enter location name of event",
									minLength: {
										value: 3,
										message: "Location name should be no less than 3 characters",
									},
									maxLength: {
										value: 100,
										message: "Location name should be no more than 100 characters",
									},
								})}
								className='w-full px-md py-sm border border-border rounded-md bg-surface text-text-primary placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-colors'
								placeholder='Please enter location name of event (3-100 characters)'
							/>
							{errors.location_name && <p className='text-error text-sm'>{errors.location_name.message}</p>}
						</div>
					{/* Location address - required, 5 - 200 characters */}
						<div className='space-y-sm'>
							<FormLabel htmlFor='location-address-input'>Location address</FormLabel>
							<input
								id='location-address-input'
								{...register("location_address", {
									required: "Please enter location address of event",
									minLength: {
										value: 5,
										message: "Location address should be no less than 5 characters",
									},
									maxLength: {
										value: 200,
										message: "Location address should be no more than 200 characters",
									},
								})}
								className='w-full px-md py-sm border border-border rounded-md bg-surface text-text-primary placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-colors'
								placeholder='Please enter location address of event (5-200 characters)'
							/>
							{errors.location_address && <p className='text-error text-sm'>{errors.location_address.message}</p>}
						</div>
					{/* Max capacity - required, 1 - 10000 people */}
						<div className='space-y-sm'>
							<FormLabel htmlFor='max-capacity-input'>Max capacity (people)</FormLabel>
							<input
								type="number"
								id='max-capacity-input'
								{...register("max_capacity", {
									required: "Please enter max capacity of event",
									min: {
										value: 1,
										message: "Max capacity should be no less than 1 person",
									},
									max: {
										value: 10000,
										message: "Max capacity should be no more than 10000 people",
									},
									valueAsNumber: true,
								})}
								className='w-full px-md py-sm border border-border rounded-md bg-surface text-text-primary placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-colors'
								placeholder='Please enter max capacity (1 - 10000 people)'
							/>
							{errors.max_capacity && <p className='text-error text-sm'>{errors.max_capacity.message}</p>}
						</div>
					{/* Image of cover page - required, PNG no more than 5MB */}
						<ImageUpload
							onChange={handleImageChange}
							error={imageError}
						/>
					{/* Server error display */}
					{serverError && (
						<div className='bg-error/10 border border-error/30 rounded-md p-md'>
							<p className='text-error text-sm'>{serverError}</p>
						</div>
					)}
					{/* Buttons display */}
					<div className='flex items-center justify-end gap-md pt-md'>
						<Button
							type='button'
							variant='outline'
							onClick={onClose}
							disabled={isSubmitting}
						>
						Cancel
						</Button>
						<Button type='submit' disabled={isSubmitting}>
							{isSubmitting ? "Submitting..." : "Submit"}
						</Button>
					</div>
					</form>
						</div>
					</div>
				</div>
			</div>
	);
}