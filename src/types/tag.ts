export interface Tag {
  id: number
  name: string
  color: string | null
  createdAt: string
  updatedAt: string
}

export interface TagPayload {
  name: string
  color: string | null
}
