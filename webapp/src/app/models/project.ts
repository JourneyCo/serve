export interface Project {
  id?: number,
  name: string,
  required: number,
  needed: number,
  leader_id: string,
  location_id: number,
  date?: number,
  created_at: Date,
  updated_at?: Date,
}

export interface Registration {
  account_id: string,
  project_id: number,
  qty_enroll: number,
  lead?: boolean,
  updated_at?: Date
}
