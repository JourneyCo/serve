export type Location = {
  latitude: number,
  longitude: number,
  id?: number,
  info: string,
  street: string,
  number: number,
  city: string,
  state: string,
  postal_code: string,
  formatted_address: string,
  created_at: Date,
  updated_at?: Date,
}
