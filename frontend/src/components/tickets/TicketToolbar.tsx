import { Plus, Search } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

export type TicketStatusFilter = "all" | "open" | "completed"

interface TicketToolbarProps {
  search: string
  onSearchChange: (value: string) => void
  onManageTags: () => void
  onCreateTicket: () => void
}

export function TicketToolbar({ search, onSearchChange, onManageTags, onCreateTicket }: TicketToolbarProps) {
  return (
    <div className="flex items-center gap-4">
      <div className="relative w-full max-w-[640px] flex-1">
        <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          value={search}
          onChange={(event) => onSearchChange(event.target.value)}
          placeholder="搜索"
          aria-label="搜索标题"
          className="h-11 rounded-full border-border/70 bg-white pl-9 pr-10 shadow-sm shadow-black/5"
        />
      </div>
      <Button variant="outline" onClick={onManageTags} className="h-11 rounded-full px-4">
        <Plus className="h-4 w-4" />
        管理标签
      </Button>
      <Button onClick={onCreateTicket} className="h-11 rounded-full px-4">
        <Plus className="h-4 w-4" />
        新建 Ticket
      </Button>
    </div>
  )
}
