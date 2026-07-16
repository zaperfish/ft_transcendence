import type { ReactNode } from 'react';

interface SettingsPanelProps {
	title: string;
	description: string;
	children: ReactNode;
	variant?: 'default' | 'danger';
}

export function SettingsPanel({
	title,
	description,
	children,
	variant = 'default',
}: SettingsPanelProps) {
	const isDanger = variant === 'danger';

	return (
		<section
			className={`w-full rounded-lg border bg-surface p-2xl shadow-sm ${
				isDanger ? 'border-destructive/30' : 'border-border'
			}`}
		>
			<h2
				className={`text-2xl font-heading font-bold mb-sm ${
					isDanger ? 'text-destructive' : 'text-text-primary'
				}`}
			>
				{title}
			</h2>
			<p className="text-text-secondary mb-xl">{description}</p>
			{children}
		</section>
	);
}
