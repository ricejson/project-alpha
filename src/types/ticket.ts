import type { Tag } from "@/types/tag"

export interface Ticket {
  id: number
  title: string
  description: string
  completed: boolean
  completedAt: string | null
  tags: Tag[]
  createdAt: string
  updatedAt: string
}

export interface TicketPayload {
  title: string
  description: string
  tagIds: number[]
}

export interface ListTicketsParams {
  title?: string
  tagIds?: number[]
  completed?: boolean | null
  page?: number
  pageSize?: number
}

export interface Pagination {
  page: number
  pageSize: number
  total: number
  totalPages: number
}

export interface TicketListResponse {
  items: Ticket[]
  pagination: Pagination
}
