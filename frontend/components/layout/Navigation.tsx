"use client";

import Link from 'next/link';
import Image from 'next/image';
import { Menu } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/Avatar';
import { Button } from '@/components/ui/Button';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/Dropdown-menu';
import { useAuth } from '@/lib/hooks/useAuth';
import { useTheme } from '@/lib/context/ThemeContext';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';

/**
 * Top-level navigation bar with logo, primary links (links for desktop and dropdown menu for mobile), and user dropdown menu.
 *
 * Handles user logout and redirects to different routes.
 */
export default function navigation() {
	const { user, logout, isOnline } = useAuth();
	const { theme } = useTheme();
	const isClassic = theme === 'classic';
	const router = useRouter();

	const handleLogout = async () => {
		if (!isOnline) {
			toast.error('Logout is unavailable when you are offline, please retry later.');
			return;
		}
		await logout();
		router.push('/login');
	};

	const navItems = [
		{ href: '/home', label: 'Discover' },
		{ href: '/events', label: 'MyEvents' },
		{ href: '/about', label: 'About' },
	];

	return (
		<header
			className={cn(
				'sticky top-0 z-50 px-4 py-3 lg:px-lg lg:py-md',
				isClassic
					? 'border-b border-border bg-surface'
					: 'border-b border-chrome-nav bg-chrome-nav backdrop-blur-md',
			)}
		>
			<div className='flex items-center justify-between w-full'>
				{/* left: logo + navigation links (desktop) */}
				<div className='flex items-center gap-4 lg:gap-lg'>
					<Link href='/' className='flex items-center gap-sm'>
						<Image src='/logo.png' alt='Camaraderie logo' width={32} height={32} className='h-8 w-8' />
						<span
							className={cn(
								'font-heading text-xl font-bold',
								isClassic ? 'text-text-primary' : 'text-chrome-title',
							)}
						>
						Camaraderie
						</span>
					</Link>
					{/* NavLinks for desktop */}
					<nav className='hidden lg:flex items-center gap-md'>
						{navItems.map((item) => (
							<Link
								key={item.href}
								href={item.href}
								className={cn(
									'transition-colors',
									isClassic
										? 'text-text-secondary hover:text-primary hover:shadow'
										: 'text-chrome-nav hover:text-chrome-nav-hover',
								)}
							>
							{item.label}
							</Link>
						))}
					</nav>
				</div>
				{/* right: navigation links (mobile) + user avatar */}
				<div className='flex items-center gap-2 lg:gap-md'>
					{/* NavLinks for Mobile (dropdown menu) */}
					<div className='lg:hidden'>
						<DropdownMenu>
							<DropdownMenuTrigger asChild>
								<Button
									variant='ghost'
									size='icon'
									className={cn(
										!isClassic && 'text-chrome-nav hover:bg-white/10 hover:text-chrome-title',
									)}
								>
									<Menu className='h-5 w-5'/>
								</Button>
							</DropdownMenuTrigger>
							<DropdownMenuContent align='end' className='min-w-45'>
								{navItems.map((item) => (
									<DropdownMenuItem key={item.href} asChild>
										<Link href={item.href}>{item.label}</Link>
									</DropdownMenuItem>
								))}
							</DropdownMenuContent>
						</DropdownMenu>
					</div>
					{/* User avatar dropdown menu */}
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
							<DropdownMenuItem
								onClick={() => {
									if (!isOnline) {
										toast.error('Settings is unavailable when you are offline, please retry later.')
										return;
									}
									router.push('/settings');
								}}
								className={`cursor-pointer w-full ${
									!isOnline ? 'opacity-50 cursor-not-allowed' : ''
								}`}
								aria-disabled={!isOnline}
							>
								Settings
							</DropdownMenuItem>
							<DropdownMenuItem
								onClick={handleLogout}
								className={
									`cursor-pointer text-error hover:text-error/80! ${
										!isOnline ? 'opacity-50 cursor-not-allowed' : ''
									}`
								}
								aria-disabled={!isOnline}
							>
							Logout
							</DropdownMenuItem>
						</DropdownMenuContent>
					</DropdownMenu>
				</div>
			</div>
		</header>
	);
}
