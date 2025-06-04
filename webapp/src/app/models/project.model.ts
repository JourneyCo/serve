import {Accessory} from './accessories';

export interface Project {
  id: number;
  google_id: number;
  title: string;
  description: string;
  rich_description?: string;
  website?: string;
  time: string
  project_date: string;
  max_capacity: number;
  current_registrations: number;
  area: string | null;
  location_address: string | null;
  latitude: number | null;
  longitude: number | null;
  serve_lead_id: string | null;
  serve_lead?: {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
  };
  ages?: string;
  types?: Accessory[];
  created_at: string;
  updated_at: string;
  encoded_address?: string;
}
