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
      <div className="divide-y">
        {Array.from({ length: 5 }, (_, index) => (
          <div key={index} className="space-y-3 px-4 py-4 lg:px-6">
            <Skeleton className="h-5 w-2/3" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-6 w-48" />
          </div>
        ))}
      </div>
    )
  }

  if (error) {
    return (
      <div className="m-4 flex items-start gap-3 rounded-md border border-destructive/30 p-4 text-sm text-destructive lg:m-6">
        <AlertCircle className="mt-0.5 h-4 w-4 shrink-0" />
        <span>{error}</span>
      </div>
    )
  }

  if (!tickets.length) {
    return <div className="m-4 rounded-md border p-10 text-center text-sm text-muted-foreground lg:m-6">暂无 Ticket</div>
  }

  return (
    <div className="divide-y">
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
