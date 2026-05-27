interface ButtonProps {
	loading?: boolean;
	disabled?: boolean;
	children: React.ReactNode;
	onClick?: () => void;
}

/**
 * A generic button component that supports loading and disabled states.
 *
 * When `loading` is true, the button is automatically disabled and typically
 * shows a loading indicator, preventing duplicate clicks.
 *
 * @param props - The button properties.
 * @param props.loading - Whether the button is in a loading state; disables the button when true.
 * @param props.disabled - Whether the button is disabled (stacks with loading).
 * @param props.children - The content displayed inside the button (text, icons, etc.).
 * @param props.onClick - Callback fired when the button is clicked.
 * @returns A rendered button element.
 */
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