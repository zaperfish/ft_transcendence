export interface User {
	id: string;
	username: string;
	email: string;
	avatar?: string; // Remain undefined, will use first alphabet of username instead
}