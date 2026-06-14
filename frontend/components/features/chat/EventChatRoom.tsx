'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';

interface EventChatRoomProps {
	eventId: number;
}

export default function EventChatRoom({ eventId }: EventChatRoomProps) {
	const [draft, setDraft] = useState('');

	return (
		<div className="min-h-screen bg-surface-dim p-lg">
			<div className="mx-auto flex h-[calc(100vh-2rem)] max-w-3xl flex-col">
				<div className="mb-md">
					<h1 className="text-2xl font-heading font-bold text-text-primary">
						Event Chat
					</h1>
					<p className="mt-xs text-sm text-text-secondary">
						Event #{eventId}
					</p>
				</div>

				<div className="flex min-h-0 flex-1 flex-col rounded-lg border border-border bg-surface shadow-sm">
					<div className="min-h-0 flex-1 overflow-y-auto px-lg py-lg">
						<div className="flex h-full items-center justify-center">
							<p className="text-sm text-text-tertiary">No messages yet.</p>
						</div>
					</div>

					<div className="border-t border-border px-lg py-md">
						<div className="flex items-end gap-md">
							<div className="flex-1">
								<Input
									value={draft}
									onChange={(e) => setDraft(e.target.value)}
									placeholder="Write a message..."
								/>
							</div>
							<Button disabled>
								Send
							</Button>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
