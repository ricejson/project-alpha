import { CheckSquare2, Edit2, Square, Trash2 } from "lucide-react"
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
    <article className="grid grid-cols-[auto_1fr_auto] gap-4 rounded-2xl border border-border/65 bg-white p-5 shadow-sm shadow-black/[0.04] transition-colors duration-200 hover:border-primary/25">
      <button
        type="button"
        aria-label={ticket.completed ? "取消完成" : "完成 ticket"}
        disabled={busy}
        onClick={() => onToggleComplete(ticket)}
        className="mt-0.5 text-[#b8c4e6] transition-colors duration-200 hover:text-primary disabled:opacity-50"
      >
        {ticket.completed ? <CheckSquare2 className="h-5 w-5 text-primary" /> : <Square className="h-5 w-5" />}
      </button>

      <div className="min-w-0">
        <button
          type="button"
          className={cn(
            "block min-w-0 text-left text-lg font-semibold tracking-tight transition-colors duration-200 hover:text-primary",
            ticket.completed ? "text-muted-foreground line-through" : "text-foreground",
          )}
          onClick={() => onEdit(ticket)}
        >
          {ticket.title}
        </button>

        {ticket.description ? (
          <p className="mt-3 line-clamp-2 whitespace-pre-wrap text-sm leading-6 text-muted-foreground">
            {ticket.description}
          </p>
        ) : null}

        {ticket.tags.length ? (
          <div className="mt-3 flex flex-wrap items-center gap-2">
            {ticket.tags.map((tag) => (
              <TagBadge key={tag.id} tag={tag} />
            ))}
          </div>
        ) : null}

        <div className="mt-3 flex flex-wrap gap-3 text-xs text-muted-foreground">
          <span>创建于 {formatDateTime(ticket.createdAt)}</span>
          {ticket.completedAt ? <span>完成于 {formatDateTime(ticket.completedAt)}</span> : null}
        </div>
      </div>

      <div className="flex shrink-0 items-start gap-2">
        <Button
          aria-label={`编辑 ${ticket.title}`}
          variant="ghost"
          size="icon"
          disabled={busy}
          onClick={() => onEdit(ticket)}
          className="h-8 w-8 rounded-lg text-muted-foreground hover:text-foreground"
        >
          <Edit2 className="h-4 w-4" />
        </Button>
        <Button
          aria-label={`删除 ${ticket.title}`}
          variant="ghost"
          size="icon"
          disabled={busy}
          onClick={() => onDelete(ticket)}
          className="h-8 w-8 rounded-lg text-rose-500 hover:bg-rose-50 hover:text-rose-600"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
    </article>
  )
}
