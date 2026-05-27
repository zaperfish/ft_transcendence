import { useContext } from 'react';
import { AuthContext } from '@/lib/context/AuthContext';

/**
 *
 * Get elements of AuthContext interface inside of AuthProvider
 */
export const useAuth = () => useContext(AuthContext);