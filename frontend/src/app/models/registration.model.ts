import { User } from './user.model';
import { Project } from './project.model';

export interface Registration {
  id: number;
  userId: string;
  projectId: number;
  status: string; // "registered", "cancelled", "completed"
  guestCount: number;
  isProjectLead: boolean;
  createdAt: string;
  updatedAt: string;
  user?: User;
  project?: Project;
}
