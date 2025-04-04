
export interface Project {
  id: number;
  title: string;
  shortDescription: string;
  description: string;
  startTime: string;
  endTime: string;
  projectDate: string;
  maxCapacity: number;
  currentRegistrations: number;
  locationName: string | null;
  locationAddress: string | null;
  latitude: number | null;
  longitude: number | null;
  wheelchairAccessible: boolean;
  leadUserId: string | null;
  leadUser?: {
    id: string;
    name: string;
    email: string;
    phone: string;
  };
  tools?: {
    id: number;
    name: string;
  }[];
  createdAt: string;
  updatedAt: string;
}
