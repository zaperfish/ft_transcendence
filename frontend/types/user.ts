export interface User {
	id: number;
	name: string;
	email: string;
	avatar?: string; // Remain undefined, will use first alphabet of username instead
}