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

export default function navigation() {
	const { user, logout } = useAuth();
	const router = useRouter();

	const apiLogout = async () => {
		await logout();
		router.push('/login');
	};

	return (
		<header className=''>
			<div className=''>
				{/* left: logo + navigation links */}
				<div className=''>
					<Link href='/' className=''>
					Meetup
					</Link>
					<nav className=''>
						<Link href='/home' className=''>
						Discover
						</Link>
						<Link href='/events' className=''>
						My Events
						</Link>
					</nav>
				</div>
				{/* right: notification + user avatar */}
				<div className=''>
					<Button variant='ghost' size='icon' aria-label='Notification'>
						<Bell className=''/>
					</Button>
					
				</div>
			</div>
		</header>
	);
}