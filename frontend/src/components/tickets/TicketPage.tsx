import { useEffect, useMemo, useState } from "react"
import { ConfirmDialog } from "@/components/ui/confirm-dialog"
import { PaginationControls } from "@/components/tickets/PaginationControls"
import { TicketFormDialog } from "@/components/tickets/TicketFormDialog"
import { TicketList } from "@/components/tickets/TicketList"
import { TicketToolbar, type TicketStatusFilter } from "@/components/tickets/TicketToolbar"
import { useTickets } from "@/hooks/useTickets"
import type { Tag } from "@/types/tag"
import type { ListTicketsParams, Ticket, TicketPayload } from "@/types/ticket"

const PAGE_SIZE = 20

interface TicketPageProps {
  tags: Tag[]
  tagsLoading: boolean
  tagsVersion: number
  onManageTags: () => void
  notify: (type: "success" | "error", message: string) => void
}

export function TicketPage({ tags, tagsLoading, tagsVersion, onManageTags, notify }: TicketPageProps) {
  const [search, setSearch] = useState("")
  const [titleQuery, setTitleQuery] = useState("")
  const [selectedTagIds, setSelectedTagIds] = useState<number[]>([])
  const [status, setStatus] = useState<TicketStatusFilter>("all")
  const [page, setPage] = useState(1)
  const [formOpen, setFormOpen] = useState(false)
  const [editingTicket, setEditingTicket] = useState<Ticket | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Ticket | null>(null)
  const [busyTicketId, setBusyTicketId] = useState<number | null>(null)

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setTitleQuery(search.trim())
      setPage(1)
    }, 300)
    return () => window.clearTimeout(timer)
  }, [search])

  const completed = status === "all" ? null : status === "completed"
  const queryParams = useMemo<ListTicketsParams>(
    () => ({
      title: titleQuery,
      tagIds: selectedTagIds,
      completed,
      page,
      pageSize: PAGE_SIZE,
    }),
    [completed, page, selectedTagIds, titleQuery],
  )

  const tickets = useTickets(queryParams)

  useEffect(() => {
    if (tagsVersion > 0) {
      void tickets.refresh()
    }
  }, [tagsVersion, tickets.refresh])

  useEffect(() => {
    const validTagIds = new Set(tags.map((tag) => tag.id))
    setSelectedTagIds((current) => {
      const next = current.filter((tagId) => validTagIds.has(tagId))
      return next.length === current.length ? current : next
    })
  }, [tags])

  function openCreateForm() {
    setEditingTicket(null)
    setFormOpen(true)
  }

  function openEditForm(ticket: Ticket) {
    setEditingTicket(ticket)
    setFormOpen(true)
  }

  async function submitTicket(payload: TicketPayload) {
    if (editingTicket) {
      await tickets.updateTicket(editingTicket.id, payload)
      notify("success", "Ticket 已更新")
    } else {
      await tickets.createTicket(payload)
      notify("success", "Ticket 已创建")
    }
    setFormOpen(false)
    setEditingTicket(null)
  }

  async function toggleComplete(ticket: Ticket) {
    setBusyTicketId(ticket.id)
    try {
      if (ticket.completed) {
        await tickets.uncompleteTicket(ticket.id)
        notify("success", "已取消完成")
      } else {
        await tickets.completeTicket(ticket.id)
        notify("success", "Ticket 已完成")
      }
    } catch (error) {
      notify("error", error instanceof Error ? error.message : "更新完成状态失败")
    } finally {
      setBusyTicketId(null)
    }
  }

  async function confirmDelete() {
    if (!deleteTarget) return

    setBusyTicketId(deleteTarget.id)
    try {
      await tickets.deleteTicket(deleteTarget.id)
      notify("success", "Ticket 已删除")
      if (tickets.tickets.length === 1 && page > 1) {
        setPage((value) => value - 1)
      }
      setDeleteTarget(null)
    } catch (error) {
      notify("error", error instanceof Error ? error.message : "删除 Ticket 失败")
    } finally {
      setBusyTicketId(null)
    }
  }

  return (
    <main className="min-h-screen bg-background">
      <TicketToolbar
        search={search}
        status={status}
        tags={tags}
        selectedTagIds={selectedTagIds}
        onSearchChange={setSearch}
        onStatusChange={(value) => {
          setStatus(value)
          setPage(1)
        }}
        onSelectedTagIdsChange={(value) => {
          setSelectedTagIds(value)
          setPage(1)
        }}
        onCreateTicket={openCreateForm}
        onManageTags={onManageTags}
      />

      <TicketList
        tickets={tickets.tickets}
        loading={tickets.loading}
        error={tickets.error}
        busyTicketId={busyTicketId}
        onEdit={openEditForm}
        onDelete={setDeleteTarget}
        onToggleComplete={toggleComplete}
      />

      <PaginationControls pagination={tickets.pagination} onPageChange={setPage} />

      <TicketFormDialog
        open={formOpen}
        ticket={editingTicket}
        tags={tags}
        tagsLoading={tagsLoading}
        onClose={() => {
          setFormOpen(false)
          setEditingTicket(null)
        }}
        onSubmit={submitTicket}
      />

      <ConfirmDialog
        open={Boolean(deleteTarget)}
        title="删除 Ticket"
        description={deleteTarget ? `确认删除「${deleteTarget.title}」？` : ""}
        confirmLabel="删除"
        busy={busyTicketId === deleteTarget?.id}
        onCancel={() => setDeleteTarget(null)}
        onConfirm={confirmDelete}
      />
    </main>
  )
}
