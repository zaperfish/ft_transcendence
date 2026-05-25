import { useContext } from 'react';
import { AuthContext } from '@/lib/context/AuthContext';

export const useAuth = () => useContext(AuthContext);