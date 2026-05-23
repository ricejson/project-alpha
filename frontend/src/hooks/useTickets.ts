import { useCallback, useEffect, useMemo, useState } from "react"
import { api } from "@/lib/api"
import type { ListTicketsParams, Pagination, Ticket, TicketPayload } from "@/types/ticket"

const DEFAULT_PAGINATION: Pagination = {
  page: 1,
  pageSize: 20,
  total: 0,
  totalPages: 0,
}

interface UseTicketsState {
  tickets: Ticket[]
  pagination: Pagination
  loading: boolean
  error: string | null
}

export function useTickets(params: ListTicketsParams) {
  const [state, setState] = useState<UseTicketsState>({
    tickets: [],
    pagination: DEFAULT_PAGINATION,
    loading: true,
    error: null,
  })

  const stableParams = useMemo(
    () => ({
      title: params.title ?? "",
      tagIds: params.tagIds ?? [],
      completed: params.completed ?? null,
      page: params.page ?? 1,
      pageSize: params.pageSize ?? 20,
    }),
    [params.completed, params.page, params.pageSize, params.tagIds, params.title],
  )

  const loadTickets = useCallback(async () => {
    setState((current) => ({ ...current, loading: true, error: null }))
    try {
      const result = await api.listTickets(stableParams)
      setState({
        tickets: result.items,
        pagination: result.pagination,
        loading: false,
        error: null,
      })
    } catch (error) {
      setState((current) => ({
        ...current,
        loading: false,
        error: error instanceof Error ? error.message : "Failed to load tickets",
      }))
    }
  }, [stableParams])

  useEffect(() => {
    void loadTickets()
  }, [loadTickets])

  const createTicket = useCallback(
    async (payload: TicketPayload) => {
      const ticket = await api.createTicket(payload)
      await loadTickets()
      return ticket
    },
    [loadTickets],
  )

  const updateTicket = useCallback(
    async (id: number, payload: TicketPayload) => {
      const ticket = await api.updateTicket(id, payload)
      await loadTickets()
      return ticket
    },
    [loadTickets],
  )

  const deleteTicket = useCallback(
    async (id: number) => {
      await api.deleteTicket(id)
      await loadTickets()
    },
    [loadTickets],
  )

  const completeTicket = useCallback(
    async (id: number) => {
      const ticket = await api.completeTicket(id)
      await loadTickets()
      return ticket
    },
    [loadTickets],
  )

  const uncompleteTicket = useCallback(
    async (id: number) => {
      const ticket = await api.uncompleteTicket(id)
      await loadTickets()
      return ticket
    },
    [loadTickets],
  )

  return {
    tickets: state.tickets,
    pagination: state.pagination,
    loading: state.loading,
    error: state.error,
    refresh: loadTickets,
    createTicket,
    updateTicket,
    deleteTicket,
    completeTicket,
    uncompleteTicket,
  }
}
