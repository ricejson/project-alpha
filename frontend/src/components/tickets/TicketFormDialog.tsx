import { useEffect, useState } from "react"
import type { FormEvent } from "react"
import { Save } from "lucide-react"
import { TagPicker } from "@/components/tags/TagPicker"
import { Button } from "@/components/ui/button"
import { Dialog } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import type { Tag } from "@/types/tag"
import type { Ticket, TicketPayload } from "@/types/ticket"

interface TicketFormDialogProps {
  open: boolean
  ticket: Ticket | null
  tags: Tag[]
  tagsLoading: boolean
  onClose: () => void
  onSubmit: (payload: TicketPayload) => Promise<void>
}

export function TicketFormDialog({
  open,
  ticket,
  tags,
  tagsLoading,
  onClose,
  onSubmit,
}: TicketFormDialogProps) {
  const [title, setTitle] = useState("")
  const [description, setDescription] = useState("")
  const [tagIds, setTagIds] = useState<number[]>([])
  const [formError, setFormError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    if (open) {
      setTitle(ticket?.title ?? "")
      setDescription(ticket?.description ?? "")
      setTagIds(ticket?.tags.map((tag) => tag.id) ?? [])
      setFormError(null)
      setBusy(false)
    }
  }, [open, ticket])

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const trimmedTitle = title.trim()

    if (!trimmedTitle) {
      setFormError("标题不能为空")
      return
    }
    if (trimmedTitle.length > 120) {
      setFormError("标题不能超过 120 个字符")
      return
    }
    if (description.length > 5000) {
      setFormError("描述不能超过 5000 个字符")
      return
    }

    setBusy(true)
    setFormError(null)
    try {
      await onSubmit({
        title: trimmedTitle,
        description,
        tagIds,
      })
    } catch (error) {
      setFormError(error instanceof Error ? error.message : "保存 ticket 失败")
    } finally {
      setBusy(false)
    }
  }

  return (
    <Dialog
      open={open}
      title={ticket ? "编辑 Ticket" : "新建 Ticket"}
      onClose={onClose}
      className="max-w-2xl"
    >
      <form className="space-y-4" onSubmit={submit}>
        <div className="space-y-2">
          <label className="text-sm font-medium" htmlFor="ticket-title">
            标题
          </label>
          <Input
            id="ticket-title"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            maxLength={120}
            placeholder="Ticket 标题"
            autoFocus
          />
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium" htmlFor="ticket-description">
            描述
          </label>
          <Textarea
            id="ticket-description"
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            maxLength={5000}
            placeholder="补充说明"
          />
          <div className="text-right text-xs text-muted-foreground">{description.length}/5000</div>
        </div>

        <div className="space-y-2">
          <div className="text-sm font-medium">标签</div>
          <TagPicker
            tags={tags}
            selectedIds={tagIds}
            onChange={setTagIds}
            label="选择标签"
            disabled={tagsLoading}
          />
        </div>

        {formError ? <p className="text-sm text-destructive">{formError}</p> : null}

        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline" onClick={onClose} disabled={busy}>
            取消
          </Button>
          <Button type="submit" disabled={busy}>
            <Save className="h-4 w-4" />
            {busy ? "保存中" : "保存"}
          </Button>
        </div>
      </form>
    </Dialog>
  )
}
