import { useEffect, useMemo, useState } from "react"
import { CheckCircle2, Circle, ListFilter } from "lucide-react"
import type { ReactNode } from "react"
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

  const activeFilters = selectedTagIds.length + (status === "all" ? 0 : 1) + (titleQuery ? 1 : 0)
  const visibleTags = useMemo(() => [...tags].sort((a, b) => a.name.localeCompare(b.name)), [tags])
  const tagCounts = useMemo(() => {
    const counts = new Map<number, number>()
    for (const ticket of tickets.tickets) {
      for (const tag of ticket.tags) {
        counts.set(tag.id, (counts.get(tag.id) ?? 0) + 1)
      }
    }
    return counts
  }, [tickets.tickets])

  return (
    <main className="flex min-h-screen flex-col">
      <div className="border-b border-border/60 bg-white">
        <div className="mx-auto grid w-full max-w-[1600px] gap-4 px-4 py-4 sm:px-6 lg:grid-cols-[240px_1fr] lg:items-center lg:px-8">
          <div className="text-[1.15rem] font-semibold tracking-tight text-primary lg:pl-2">Project Alpha</div>
          <div className="min-w-0 lg:flex lg:justify-center">
            <TicketToolbar
              search={search}
              onSearchChange={setSearch}
              onManageTags={onManageTags}
              onCreateTicket={openCreateForm}
            />
          </div>
        </div>
      </div>

      <section className="mx-auto flex w-full max-w-[1600px] flex-1 gap-0 lg:px-0">
        <aside className="hidden w-[280px] shrink-0 border-r border-border/60 bg-white px-4 py-6 lg:block">
          <SidebarSection title="状态">
            <FilterButton
              active={status === "all"}
              icon={<ListFilter className="h-4 w-4" />}
              onClick={() => {
                setStatus("all")
                setPage(1)
              }}
            >
              全部
            </FilterButton>
            <FilterButton
              active={status === "open"}
              icon={<Circle className="h-4 w-4" />}
              onClick={() => {
                setStatus("open")
                setPage(1)
              }}
            >
              待完成
            </FilterButton>
            <FilterButton
              active={status === "completed"}
              icon={<CheckCircle2 className="h-4 w-4" />}
              onClick={() => {
                setStatus("completed")
                setPage(1)
              }}
            >
              已完成
            </FilterButton>
          </SidebarSection>

          <SidebarSection title="标签">
            <div className="max-h-[calc(100vh-16rem)] space-y-1 overflow-auto pr-1">
              {visibleTags.map((tag) => {
                const active = selectedTagIds.includes(tag.id)
                return (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => {
                      setSelectedTagIds((current) =>
                        current.includes(tag.id) ? current.filter((value) => value !== tag.id) : [...current, tag.id],
                      )
                      setPage(1)
                    }}
                    className={[
                      "flex w-full items-center justify-between rounded-2xl px-3 py-2 text-left text-sm transition-colors duration-200",
                      active ? "bg-[#dbeafe] text-[#1d4ed8]" : "hover:bg-[#f5f7fb]",
                    ].join(" ")}
                  >
                    <span className="flex min-w-0 items-center gap-2">
                      <span
                        className="h-2.5 w-2.5 shrink-0 rounded-full"
                        style={tag.color ? { backgroundColor: tag.color } : undefined}
                      />
                      <span className="truncate">{tag.name}</span>
                    </span>
                    <span className="text-xs text-muted-foreground">{tagCounts.get(tag.id) ?? 0}</span>
                  </button>
                )
              })}
            </div>
          </SidebarSection>
        </aside>

        <div className="flex-1 bg-[#f6f8fc]">
          <div className="border-b border-border/60 bg-white px-4 py-3 lg:px-6">
            <div className="flex items-center justify-between gap-4">
              <div className="text-sm text-muted-foreground">批量操作</div>
              <div className="flex items-center gap-2">
                <div className="rounded-full border border-border/70 bg-white px-3 py-1.5 text-sm text-muted-foreground">
                  {activeFilters ? `${activeFilters} 个筛选` : "全部 tickets"}
                </div>
              </div>
            </div>
          </div>

          <div className="px-4 py-4 lg:px-6 lg:py-6">
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
          </div>
        </div>
      </section>

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

function SidebarSection({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section className="mb-8">
      <h2 className="mb-3 text-sm font-semibold text-muted-foreground">{title}</h2>
      <div className="space-y-2">{children}</div>
    </section>
  )
}

function FilterButton({
  active,
  icon,
  children,
  onClick,
}: {
  active: boolean
  icon: ReactNode
  children: ReactNode
  onClick: () => void
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={[
        "flex w-full items-center gap-3 rounded-2xl px-3 py-2 text-left text-sm transition-colors duration-200",
        active ? "bg-[#dbeafe] text-[#1d4ed8]" : "text-foreground hover:bg-[#f5f7fb]",
      ].join(" ")}
    >
      <span className="text-inherit">{icon}</span>
      <span>{children}</span>
    </button>
  )
}
