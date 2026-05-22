interface ButtonProps {
	loading?: boolean;
	disabled?: boolean;
	children: React.ReactNode;
	onClick?: () => void;
}

export function Button({ loading, disabled, children, onClick }: ButtonProps) {
	return (
		<button onClick={onClick} disabled={disabled || loading} className={`
				inline-flex items-center justify-center
				px-lg py-sm
				rounded-md
				font-medium text-sm
				transition-colors duration-200
				bg-primary text-surface
				hover:bg-primary-dim
				active:scale-[0.98]
				focus:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2
				disabled:opacity-50 disabled:cursor-not-allowed
				disabled:hover:bg-primary
			`}
		>
			{loading? 'Loading...' : children}
		</button>
	);
}