export interface Project {
  id?: number,
  name: string,
  required: number,
  registered: number,
  leader_id: string,
  location_id: number,
  start_time: Date,
  end_time: Date,
  created_at: Date,
  updated_at?: Date,
  wheelchair: boolean,
  short_description: string,
  long_description: string,
  enabled: boolean,
  status: string,
}

export interface Registration {
  account_id: string,
  first: string,
  last: string,
  cellphone: string,
  email: string,
  project_id: number,
  qty_enroll: number,
  lead?: boolean,
  updated_at?: Date
}
