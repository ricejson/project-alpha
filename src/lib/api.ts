import type { ApiResponse } from "@/types/api"
import { ApiClientError } from "@/types/api"
import type { Tag, TagPayload } from "@/types/tag"
import type { ListTicketsParams, Ticket, TicketListResponse, TicketPayload } from "@/types/ticket"

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1"

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  if (init.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json")
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers,
  })

  if (response.status === 204) {
    return undefined as T
  }

  const payload = (await response.json()) as ApiResponse<T>
  if (!response.ok || payload.error) {
    const error = payload.error
    throw new ApiClientError(
      error?.message ?? "Request failed",
      error?.code ?? "REQUEST_ERROR",
      response.status,
      error?.details,
    )
  }

  return payload.data
}

function jsonBody(value: unknown) {
  return JSON.stringify(value)
}

export const api = {
  listTags: () => request<Tag[]>("/tags"),
  createTag: (payload: TagPayload) => request<Tag>("/tags", { method: "POST", body: jsonBody(payload) }),
  updateTag: (id: number, payload: TagPayload) =>
    request<Tag>(`/tags/${id}`, { method: "PUT", body: jsonBody(payload) }),
  deleteTag: (id: number) => request<void>(`/tags/${id}`, { method: "DELETE" }),

  listTickets: (params: ListTicketsParams) => request<TicketListResponse>(`/tickets${ticketQuery(params)}`),
  getTicket: (id: number) => request<Ticket>(`/tickets/${id}`),
  createTicket: (payload: TicketPayload) =>
    request<Ticket>("/tickets", { method: "POST", body: jsonBody(payload) }),
  updateTicket: (id: number, payload: TicketPayload) =>
    request<Ticket>(`/tickets/${id}`, { method: "PUT", body: jsonBody(payload) }),
  deleteTicket: (id: number) => request<void>(`/tickets/${id}`, { method: "DELETE" }),
  completeTicket: (id: number) => request<Ticket>(`/tickets/${id}/complete`, { method: "POST" }),
  uncompleteTicket: (id: number) => request<Ticket>(`/tickets/${id}/uncomplete`, { method: "POST" }),
  addTicketTag: (ticketId: number, tagId: number) =>
    request<Ticket>(`/tickets/${ticketId}/tags/${tagId}`, { method: "POST" }),
  removeTicketTag: (ticketId: number, tagId: number) =>
    request<Ticket>(`/tickets/${ticketId}/tags/${tagId}`, { method: "DELETE" }),
}

function ticketQuery(params: ListTicketsParams) {
  const query = new URLSearchParams()
  if (params.title?.trim()) query.set("title", params.title.trim())
  if (params.tagIds?.length) query.set("tagIds", params.tagIds.join(","))
  if (typeof params.completed === "boolean") query.set("completed", String(params.completed))
  if (params.page) query.set("page", String(params.page))
  if (params.pageSize) query.set("pageSize", String(params.pageSize))

  const value = query.toString()
  return value ? `?${value}` : ""
}
