export interface CreateEvent {
	description: string;
	duration: number;
	location_address: string;
	location_name: string;
	max_capacity: number;
	start_time: string;
	title: string;
}

export interface GetEvent {
	description: string;
	duration: number;
	location_address: string;
	location_name: string;
	max_capacity: number;
	start_time: string;
	title: string;
	created_at: string;
	id: string;
	num_registered: number;
	updated_at: string;
}