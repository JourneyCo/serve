import { User } from "./user.model";
import { Project } from "./project.model";

export interface Registration {
  id: number;
  user_id: string;
  project_id: number;
  status: string; // "registered", "cancelled", "completed"
  guest_count: number;
  lead_interest: boolean;
  created_at: string;
  updated_at: string;
  user?: User;
  project?: Project;
  lead?: boolean;
}
