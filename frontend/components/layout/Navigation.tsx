"use client";

import Link from 'next/link';
import { Bell } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/Avatar';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/Dropdown-menu';
import { Button } from "@/components/ui/Button";
import { useAuth } from '@/lib/hooks/useAuth';
import { useRouter } from 'next/navigation';

/**
 * Top-level navigation bar with logo, primary links, notification icon, and user dropdown menu.
 *
 * Handles user logout and redirects to different routes.
 */
export default function navigation() {
	const { user, logout } = useAuth();
	const router = useRouter();

	const handleLogout = async () => {
		await logout();
		router.push('/login');
	};

	return (
		<header className='bg-surface border-b border-border px-lg py-md'>
			<div className='flex items-center justify-between w-full'>
				{/* left: logo + navigation links */}
				<div className='flex items-center gap-lg'>
					<Link href='/' className='text-text-primary font-heading text-xl font-bold'>
					Meetup
					</Link>
					<nav className='flex items-center gap-md'>
						<Link href='/home' className='text-text-secondary hover:text-text-primary transition-colors'>
						Discover
						</Link>
						<Link href='/events' className='text-text-secondary hover:text-text-primary transition-colors'>
						My Events
						</Link>
					</nav>
				</div>
				{/* right: user avatar */}
				<div className='flex items-center gap-md'>
					<DropdownMenu>
						<DropdownMenuTrigger asChild>
							<Avatar className='h-9 w-9 cursor-pointer'>
								<AvatarImage src={user?.avatar} alt={user?.name} />
								<AvatarFallback className='bg-primary text-primary-foreground'>
									 {user?.name?.charAt(0)?.toUpperCase() || "U"}
								</AvatarFallback>
							</Avatar>
						</DropdownMenuTrigger>
						<DropdownMenuContent align='end' className='min-w-45'>
							<DropdownMenuItem asChild>
								<Link href='/settings' className='cursor-pointer w-full'>
								Settings
								</Link>
							</DropdownMenuItem>
							<DropdownMenuItem onClick={handleLogout} className='cursor-pointer text-error hover:text-error/80!'>
							Logout
							</DropdownMenuItem>
						</DropdownMenuContent>
					</DropdownMenu>
				</div>
			</div>
		</header>
	);
}