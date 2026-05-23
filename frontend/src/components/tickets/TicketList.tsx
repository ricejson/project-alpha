import { AlertCircle } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"
import { TicketListItem } from "@/components/tickets/TicketListItem"
import type { Ticket } from "@/types/ticket"

interface TicketListProps {
  tickets: Ticket[]
  loading: boolean
  error: string | null
  busyTicketId: number | null
  onEdit: (ticket: Ticket) => void
  onDelete: (ticket: Ticket) => void
  onToggleComplete: (ticket: Ticket) => void
}

export function TicketList({
  tickets,
  loading,
  error,
  busyTicketId,
  onEdit,
  onDelete,
  onToggleComplete,
}: TicketListProps) {
  if (loading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 4 }, (_, index) => (
          <div key={index} className="rounded-2xl border border-border/60 bg-white p-5 shadow-sm">
            <Skeleton className="h-5 w-2/3 rounded-full" />
            <Skeleton className="mt-5 h-4 w-1/2 rounded-full" />
            <Skeleton className="mt-4 h-6 w-64 rounded-full" />
          </div>
        ))}
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-start gap-3 rounded-2xl border border-destructive/20 bg-white p-4 text-sm text-destructive shadow-sm">
        <AlertCircle className="mt-0.5 h-4 w-4 shrink-0" />
        <span>{error}</span>
      </div>
    )
  }

  if (!tickets.length) {
    return (
      <div className="rounded-2xl border border-dashed border-border/80 bg-white px-6 py-14 text-center">
        <div className="text-lg font-semibold tracking-tight text-foreground">暂无 Ticket</div>
        <p className="mt-2 text-sm text-muted-foreground">创建一个新 ticket，或者调整搜索和筛选条件。</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {tickets.map((ticket) => (
        <TicketListItem
          key={ticket.id}
          ticket={ticket}
          busy={busyTicketId === ticket.id}
          onEdit={onEdit}
          onDelete={onDelete}
          onToggleComplete={onToggleComplete}
        />
      ))}
    </div>
  )
}
