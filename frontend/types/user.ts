export interface User {
	id: number;
	name: string;
	email: string;
	avatar?: string; // Remain undefined, will use first alphabet of username instead
}

export interface UpdateProfileRequest {
	email?: string;
  }

export interface UpdatePasswordRequest {
	current_password: string;
	newpassword: string;
	confirm_password: string;
  }

