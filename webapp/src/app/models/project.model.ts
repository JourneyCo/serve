import { Accessory } from './accessories';

export interface Project {
  id: number;
  google_id: number;
  title: string;
  description: string;
  rich_description?: string;
  website?: string;
  time: string
  project_date: Date;
  max_capacity: number;
  current_registrations: number;
  area?: string;
  location_address: string | null;
  latitude: number | null;
  longitude: number | null;
  serve_lead_id?: string;
  serve_lead_email?: string;
  serve_lead_name?: string;
  serve_lead?: {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
  };
  ages?: string;
  types?: Accessory[];
  active: boolean;
  created_at: string;
  updated_at: string;
  encoded_address?: string;
  leads?: Lead[];
}

export type Lead = {
  name?: string,
  email?: string
  phone?: string
  active: boolean
}
