import { useEffect, useMemo, useRef, useState } from "react"
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
  const rootRef = useRef<HTMLDivElement>(null)

  const selectedSet = useMemo(() => new Set(selectedIds), [selectedIds])
  const selectedTags = tags.filter((tag) => selectedSet.has(tag.id))
  const filteredTags = tags.filter((tag) => tag.name.toLowerCase().includes(query.trim().toLowerCase()))

  useEffect(() => {
    function handlePointerDown(event: PointerEvent) {
      if (!rootRef.current?.contains(event.target as Node)) {
        setOpen(false)
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setOpen(false)
      }
    }

    window.addEventListener("pointerdown", handlePointerDown)
    window.addEventListener("keydown", handleKeyDown)
    return () => {
      window.removeEventListener("pointerdown", handlePointerDown)
      window.removeEventListener("keydown", handleKeyDown)
    }
  }, [])

  function toggleTag(tagId: number) {
    if (selectedSet.has(tagId)) {
      onChange(selectedIds.filter((id) => id !== tagId))
      return
    }
    onChange([...selectedIds, tagId])
  }

  return (
    <div ref={rootRef} className="relative min-w-0">
      <Button
        type="button"
        variant="outline"
        className="h-11 w-full justify-between rounded-full border-border/80 bg-white px-4 shadow-sm shadow-black/5"
        disabled={disabled}
        aria-expanded={open}
        onClick={() => setOpen((value) => !value)}
      >
        <span className="flex min-w-0 items-center gap-2">
          <Tags className="h-4 w-4 shrink-0" />
          <span className="truncate">{selectedIds.length ? `${label} · ${selectedIds.length}` : label}</span>
        </span>
      </Button>

      {open ? (
        <div className="absolute left-0 top-[3.25rem] z-40 w-80 rounded-[24px] border border-border/70 bg-white p-3 shadow-[0_18px_40px_rgba(15,23,42,0.16)]">
          <Input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="搜索标签"
            aria-label="搜索标签"
            className="mb-2 rounded-full border-border/80 bg-background/70"
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
                      "flex h-10 w-full items-center justify-between gap-2 rounded-full px-3 text-left text-sm transition-colors duration-200 hover:bg-accent/80",
                      selected ? "bg-accent/70 font-medium" : "font-normal",
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
              className="mt-2 w-full rounded-full"
              onClick={() => onChange([])}
            >
              <X className="h-4 w-4" />
              清空
            </Button>
          ) : null}
        </div>
      ) : null}

      {selectedTags.length ? (
        <div className="mt-3 flex flex-wrap gap-1.5">
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
