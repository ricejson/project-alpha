import { X } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import type { Tag } from "@/types/tag"

interface TagBadgeProps {
  tag: Tag
  removable?: boolean
  onRemove?: (tagId: number) => void
}

export function TagBadge({ tag, removable, onRemove }: TagBadgeProps) {
  return (
    <Badge color={tag.color} className="min-w-0">
      <span className="truncate">{tag.name}</span>
      {removable ? (
        <Button
          aria-label={`移除标签 ${tag.name}`}
          variant="ghost"
          size="icon"
          className="h-4 w-4 rounded-sm p-0 hover:bg-transparent"
          onClick={() => onRemove?.(tag.id)}
        >
          <X className="h-3 w-3" />
        </Button>
      ) : null}
    </Badge>
  )
}
