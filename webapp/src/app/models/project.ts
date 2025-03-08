export interface Project {
  id?: number,
  name: string,
  required: number,
  needed: number,
  leader_id: number,
  location_id: number,
  date?: number,
  created_at: Date,
  updated_at?: Date,
}

export interface Registration {
  id: number,
  registering: number,
  user_id: number
}
