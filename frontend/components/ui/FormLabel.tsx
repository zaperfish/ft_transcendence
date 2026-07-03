interface FormLabelProps {
	htmlFor?: string;
	children: React.ReactNode;
	required?: boolean;
}

export function FormLabel({ htmlFor, children, required = true }: FormLabelProps) {
	return (
		<label
			htmlFor={htmlFor}
			className="block text-sm font-medium text-text-secondary"
		>
			{children}
			{required && <span className="text-error ml-1">*</span>}
		</label>
	);
}