import { useMemo, useState } from "react"
import { Check, Tags, X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { TagBadge } from "@/components/tags/TagBadge"
import { cn } from "@/lib/utils"
import type { Tag } from "@/types/tag"

interface TagPickerProps {
  tags: Tag[]
  selectedIds: number[]
  onChange: (selectedIds: number[]) => void
  label: string
  disabled?: boolean
  emptyLabel?: string
}

export function TagPicker({
  tags,
  selectedIds,
  onChange,
  label,
  disabled,
  emptyLabel = "无标签",
}: TagPickerProps) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState("")

  const selectedSet = useMemo(() => new Set(selectedIds), [selectedIds])
  const selectedTags = tags.filter((tag) => selectedSet.has(tag.id))
  const filteredTags = tags.filter((tag) => tag.name.toLowerCase().includes(query.trim().toLowerCase()))

  function toggleTag(tagId: number) {
    if (selectedSet.has(tagId)) {
      onChange(selectedIds.filter((id) => id !== tagId))
      return
    }
    onChange([...selectedIds, tagId])
  }

  return (
    <div className="relative">
      <Button
        type="button"
        variant="outline"
        className="h-9 min-w-40 justify-between"
        disabled={disabled}
        aria-expanded={open}
        onClick={() => setOpen((value) => !value)}
      >
        <span className="flex min-w-0 items-center gap-2">
          <Tags className="h-4 w-4 shrink-0" />
          <span className="truncate">{selectedIds.length ? `${label} ${selectedIds.length}` : label}</span>
        </span>
      </Button>

      {open ? (
        <div className="absolute left-0 top-11 z-40 w-72 rounded-md border bg-background p-3 shadow-lg">
          <Input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="搜索标签"
            aria-label="搜索标签"
            className="mb-2"
          />
          <div className="max-h-56 overflow-auto">
            {filteredTags.length ? (
              filteredTags.map((tag) => {
                const selected = selectedSet.has(tag.id)
                return (
                  <button
                    key={tag.id}
                    type="button"
                    className={cn(
                      "flex h-9 w-full items-center justify-between gap-2 rounded-md px-2 text-left text-sm hover:bg-accent",
                      selected ? "font-medium" : "font-normal",
                    )}
                    onClick={() => toggleTag(tag.id)}
                  >
                    <span className="flex min-w-0 items-center gap-2">
                      <span
                        className="h-3 w-3 shrink-0 rounded-full border"
                        style={tag.color ? { backgroundColor: tag.color, borderColor: tag.color } : undefined}
                      />
                      <span className="truncate">{tag.name}</span>
                    </span>
                    {selected ? <Check className="h-4 w-4 shrink-0" /> : null}
                  </button>
                )
              })
            ) : (
              <div className="px-2 py-3 text-sm text-muted-foreground">{emptyLabel}</div>
            )}
          </div>
          {selectedIds.length ? (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="mt-2 w-full"
              onClick={() => onChange([])}
            >
              <X className="h-4 w-4" />
              清空
            </Button>
          ) : null}
        </div>
      ) : null}

      {selectedTags.length ? (
        <div className="mt-2 flex flex-wrap gap-1.5">
          {selectedTags.map((tag) => (
            <TagBadge
              key={tag.id}
              tag={tag}
              removable
              onRemove={(tagId) => onChange(selectedIds.filter((id) => id !== tagId))}
            />
          ))}
        </div>
      ) : null}
    </div>
  )
}
