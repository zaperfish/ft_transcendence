'use client';

import { useCallback, useState, useRef } from "react";

interface ImageUploadProps {
	onChange: (file: File | null) => void;
	error?: string;
}

export function ImageUpload({ onChange, error }: ImageUploadProps) {
	const [preview, setPreview] = useState<string | null>(null);
	const [filename, setFilename] = useState<string>('');
	const [localError, setLocalError] = useState<string>('');
	const fileInputRef = useRef<HTMLInputElement>(null);

	const handleFileChange = useCallback(
		(e: React.ChangeEvent<HTMLInputElement>) => {
			const file = e.target.files?.[0];
			if (!file) return;

			setLocalError('');
			// Validate image type
			if (file.type !== "image/png") {
				setLocalError('Support only image/png type');
				onChange(null);
				return;
			}
			// Validate image size
			if (file.size > 5 * 1024 * 1024) {
				setLocalError('Image file cannot be more than 5MB');
				onChange(null);
				return;
			}
			// Return file to parent component
			setFilename(file.name);
			onChange(file);
			// Read image and save as preview
			const reader = new FileReader();
			reader.onloadend = () => setPreview(reader.result as string);
			reader.readAsDataURL(file);
		}, [onChange]
	);

	const clear = () => {
		setPreview(null);
		setFilename('');
		setLocalError('');
		onChange(null);
		// Reset ref to allow change event even if same file selected more than once
		if (fileInputRef.current) {
			fileInputRef.current.value = '';
		}
	};

	const displayError = error || localError;

	return (
		<div className="space-y-sm">
			<label className="block text-sm font-medium text-text-primary">
				Cover page of event
			</label>
			<div className="flex items-center gap-md">
				{preview ? (
					<div className="relative w-32 h-32 border border-border rounded-lg overflow-hidden">
						<img
							src={preview}
							alt="Preview"
							className="object-cover w-full h-full"
						/>
						<button
							type="button"
							onClick={clear}
							className="absolute top-1 right-1 bg-surface/80 rounded-full p-1 text-xs"
						>
							x
						</button>
					</div>
				) : (
					<div className="w-32 h-32 border-2 border-dashed border-border rounded-lg flex items-center justify-center text-text-tertiary">
						No image
					</div>
				)}
				<div className="flex-1 ">
					<input
						ref={fileInputRef}
						type="file"
						accept="image/png"
						onChange={handleFileChange}
						className="block text-sm text-text-secondary file:mr-md file:py-sm file:px-md file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary/10 file:text-primary hover:file:bg-primary/20"
					/>
					<p className="text-xs text-text-tertiary mt-xs">
						Surport PNG, 5MB at maximum
					</p>
				</div>
			</div>
			{displayError && <p className="text-sm text-error">{displayError}</p>}
		</div>
	);
}