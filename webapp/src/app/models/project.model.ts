export interface Project {
  id: number;
  title: string;
  short_description: string;
  description: string;
  start_time: string;
  end_time: string;
  project_date: string;
  max_capacity: number;
  current_registrations: number;
  location_name: string | null;
  location_address: string | null;
  latitude: number | null;
  longitude: number | null;
  wheelchair_accessible: boolean;
  lead_user_id: string | null;
  lead_user?: {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
  };
  tools?: {
    id: number;
    name: string;
  }[];
  created_at: string;
  updated_at: string;
}
