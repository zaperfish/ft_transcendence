interface InputProps {
	value: string;
	onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
	placeholder?: string;
	type?: string;
	error?: string;
	disabled?: boolean;
}

export function Input({ error, ...props }: InputProps) {
	return (
		<div className="flex flex-col gap-xs">
			<input {...props} className={`
				w-full
				px-md py-sm
				bg-surface
				border border-border
				rounded-md
				text-text-primary
				placeholder:text-text-tertiary
				focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent
				disabled:opacity-50 disabled:cursor-not-allowed
				`}
			/>
			{error && <span className="text-sm text-error mt-xs" >{error}</span>}
		</div>
	);
}