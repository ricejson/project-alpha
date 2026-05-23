import { Plus, Search, Tags } from "lucide-react"
import { TagPicker } from "@/components/tags/TagPicker"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Select } from "@/components/ui/select"
import type { Tag } from "@/types/tag"

export type TicketStatusFilter = "all" | "open" | "completed"

interface TicketToolbarProps {
  search: string
  status: TicketStatusFilter
  tags: Tag[]
  selectedTagIds: number[]
  onSearchChange: (value: string) => void
  onStatusChange: (value: TicketStatusFilter) => void
  onSelectedTagIdsChange: (value: number[]) => void
  onCreateTicket: () => void
  onManageTags: () => void
}

export function TicketToolbar({
  search,
  status,
  tags,
  selectedTagIds,
  onSearchChange,
  onStatusChange,
  onSelectedTagIdsChange,
  onCreateTicket,
  onManageTags,
}: TicketToolbarProps) {
  return (
    <div className="flex flex-col gap-3 border-b bg-background px-4 py-4 lg:px-6">
      <div className="flex flex-col justify-between gap-3 md:flex-row md:items-center">
        <div>
          <h1 className="text-xl font-semibold tracking-normal">project-alpha</h1>
          <p className="text-sm text-muted-foreground">Ticket 管理</p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" onClick={onManageTags}>
            <Tags className="h-4 w-4" />
            标签
          </Button>
          <Button onClick={onCreateTicket}>
            <Plus className="h-4 w-4" />
            新建 Ticket
          </Button>
        </div>
      </div>

      <div className="flex flex-col gap-2 xl:flex-row xl:items-start">
        <div className="relative min-w-0 flex-1">
          <Search className="pointer-events-none absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            value={search}
            onChange={(event) => onSearchChange(event.target.value)}
            placeholder="搜索标题"
            aria-label="搜索标题"
            className="pl-9"
          />
        </div>
        <div className="flex flex-col gap-2 sm:flex-row sm:items-start">
          <TagPicker
            tags={tags}
            selectedIds={selectedTagIds}
            onChange={onSelectedTagIdsChange}
            label="标签筛选"
          />
          <Select
            aria-label="完成状态"
            value={status}
            onChange={(event) => onStatusChange(event.target.value as TicketStatusFilter)}
            className="w-full sm:w-32"
          >
            <option value="all">全部</option>
            <option value="open">未完成</option>
            <option value="completed">已完成</option>
          </Select>
        </div>
      </div>
    </div>
  )
}
