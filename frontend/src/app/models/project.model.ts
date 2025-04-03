export interface Project {
  id: number;
  title: string;
  description: string;
  startTime: string;
  endTime: string;
  projectDate: string;
  maxCapacity: number;
  currentRegistrations: number;
  locationName: string | null;
  latitude: number | null;
  longitude: number | null;
  createdAt: string;
  updatedAt: string;
}
