import { CheckCircle2, Circle, Edit2, Trash2 } from "lucide-react"
import { TagBadge } from "@/components/tags/TagBadge"
import { Button } from "@/components/ui/button"
import { cn, formatDateTime } from "@/lib/utils"
import type { Ticket } from "@/types/ticket"

interface TicketListItemProps {
  ticket: Ticket
  busy?: boolean
  onEdit: (ticket: Ticket) => void
  onDelete: (ticket: Ticket) => void
  onToggleComplete: (ticket: Ticket) => void
}

export function TicketListItem({ ticket, busy, onEdit, onDelete, onToggleComplete }: TicketListItemProps) {
  return (
    <article className="grid gap-3 border-b px-4 py-4 lg:grid-cols-[auto_1fr_auto] lg:items-start lg:px-6">
      <Button
        aria-label={ticket.completed ? "取消完成" : "完成 ticket"}
        variant="ghost"
        size="icon"
        disabled={busy}
        onClick={() => onToggleComplete(ticket)}
        className="hidden lg:inline-flex"
      >
        {ticket.completed ? <CheckCircle2 className="h-5 w-5 text-emerald-600" /> : <Circle className="h-5 w-5" />}
      </Button>

      <div className="min-w-0 space-y-2">
        <div className="flex min-w-0 flex-wrap items-center gap-2">
          <button
            type="button"
            className={cn(
              "min-w-0 text-left text-base font-semibold hover:underline",
              ticket.completed ? "text-muted-foreground line-through" : "text-foreground",
            )}
            onClick={() => onEdit(ticket)}
          >
            {ticket.title}
          </button>
          <span
            className={cn(
              "rounded-md px-2 py-0.5 text-xs font-medium",
              ticket.completed ? "bg-emerald-50 text-emerald-700" : "bg-muted text-muted-foreground",
            )}
          >
            {ticket.completed ? "已完成" : "未完成"}
          </span>
        </div>

        {ticket.description ? (
          <p className="line-clamp-2 whitespace-pre-wrap text-sm text-muted-foreground">{ticket.description}</p>
        ) : null}

        <div className="flex flex-wrap items-center gap-1.5">
          {ticket.tags.length ? (
            ticket.tags.map((tag) => <TagBadge key={tag.id} tag={tag} />)
          ) : (
            <span className="text-xs text-muted-foreground">无标签</span>
          )}
        </div>

        <div className="flex flex-wrap gap-3 text-xs text-muted-foreground">
          <span>创建 {formatDateTime(ticket.createdAt)}</span>
          <span>更新 {formatDateTime(ticket.updatedAt)}</span>
          {ticket.completedAt ? <span>完成 {formatDateTime(ticket.completedAt)}</span> : null}
        </div>
      </div>

      <div className="flex shrink-0 gap-1 lg:justify-end">
        <Button
          aria-label={ticket.completed ? "取消完成" : "完成 ticket"}
          variant="ghost"
          size="icon"
          disabled={busy}
          onClick={() => onToggleComplete(ticket)}
          className="lg:hidden"
        >
          {ticket.completed ? <CheckCircle2 className="h-5 w-5 text-emerald-600" /> : <Circle className="h-5 w-5" />}
        </Button>
        <Button aria-label={`编辑 ${ticket.title}`} variant="ghost" size="icon" disabled={busy} onClick={() => onEdit(ticket)}>
          <Edit2 className="h-4 w-4" />
        </Button>
        <Button
          aria-label={`删除 ${ticket.title}`}
          variant="ghost"
          size="icon"
          disabled={busy}
          onClick={() => onDelete(ticket)}
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
    </article>
  )
}
