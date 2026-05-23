import { useEffect, useMemo, useState } from "react"
import type { FormEvent } from "react"
import { Check, Edit2, Plus, Sparkles, Trash2, X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { ConfirmDialog } from "@/components/ui/confirm-dialog"
import { Dialog } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Skeleton } from "@/components/ui/skeleton"
import type { Tag, TagPayload } from "@/types/tag"

const DEFAULT_COLOR = "#2563eb"
const COLOR_PATTERN = /^#[0-9a-fA-F]{6}$/

interface TagManagerDialogProps {
  open: boolean
  tags: Tag[]
  loading: boolean
  error: string | null
  onClose: () => void
  onCreate: (payload: TagPayload) => Promise<void>
  onUpdate: (id: number, payload: TagPayload) => Promise<void>
  onDelete: (id: number) => Promise<void>
}

export function TagManagerDialog({
  open,
  tags,
  loading,
  error,
  onClose,
  onCreate,
  onUpdate,
  onDelete,
}: TagManagerDialogProps) {
  const [editingTag, setEditingTag] = useState<Tag | null>(null)
  const [name, setName] = useState("")
  const [color, setColor] = useState(DEFAULT_COLOR)
  const [formError, setFormError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<Tag | null>(null)

  const title = editingTag ? "编辑标签" : "新建标签"
  const sortedTags = useMemo(() => [...tags].sort((a, b) => a.name.localeCompare(b.name)), [tags])

  useEffect(() => {
    if (!open) {
      resetForm()
      setDeleteTarget(null)
    }
  }, [open])

  function resetForm() {
    setEditingTag(null)
    setName("")
    setColor(DEFAULT_COLOR)
    setFormError(null)
    setBusy(false)
  }

  function startEdit(tag: Tag) {
    setEditingTag(tag)
    setName(tag.name)
    setColor(tag.color ?? DEFAULT_COLOR)
    setFormError(null)
  }

  async function submitForm(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const trimmedName = name.trim()
    const trimmedColor = color.trim()

    if (!trimmedName) {
      setFormError("标签名称不能为空")
      return
    }
    if (trimmedName.length > 40) {
      setFormError("标签名称不能超过 40 个字符")
      return
    }
    if (trimmedColor && !COLOR_PATTERN.test(trimmedColor)) {
      setFormError("颜色必须是 #RRGGBB 格式")
      return
    }

    setBusy(true)
    setFormError(null)
    try {
      const payload = {
        name: trimmedName,
        color: trimmedColor || null,
      }
      if (editingTag) {
        await onUpdate(editingTag.id, payload)
      } else {
        await onCreate(payload)
      }
      resetForm()
    } catch (error) {
      setFormError(error instanceof Error ? error.message : "保存标签失败")
    } finally {
      setBusy(false)
    }
  }

  async function confirmDelete() {
    if (!deleteTarget) return

    setBusy(true)
    setFormError(null)
    try {
      await onDelete(deleteTarget.id)
      if (editingTag?.id === deleteTarget.id) {
        resetForm()
      }
      setDeleteTarget(null)
    } catch (error) {
      setFormError(error instanceof Error ? error.message : "删除标签失败")
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      <Dialog
        open={open}
        title="标签管理"
        description="用颜色和名称建立稳定的工作分类，避免分类噪音。"
        onClose={onClose}
        className="max-w-2xl"
      >
        <form className="grid gap-4 border-b border-border/70 pb-5" onSubmit={submitForm}>
          <div className="flex items-center justify-between gap-3">
            <h3 className="flex items-center gap-2 text-sm font-semibold">
              <Sparkles className="h-4 w-4 text-muted-foreground" />
              {title}
            </h3>
            {editingTag ? (
              <Button type="button" variant="ghost" size="sm" onClick={resetForm}>
                <X className="h-4 w-4" />
                取消
              </Button>
            ) : null}
          </div>
          <div className="grid gap-3 sm:grid-cols-[1fr_12rem_auto]">
            <Input
              value={name}
              onChange={(event) => setName(event.target.value)}
              placeholder="标签名称"
              aria-label="标签名称"
              maxLength={40}
              className="rounded-full bg-background/80"
            />
            <div className="flex gap-2">
              <Input
                value={color}
                onChange={(event) => setColor(event.target.value)}
                placeholder="#2563eb"
                aria-label="标签颜色"
                className="font-mono rounded-full bg-background/80"
              />
              <input
                aria-label="选择标签颜色"
                type="color"
                value={COLOR_PATTERN.test(color) ? color : DEFAULT_COLOR}
                onChange={(event) => setColor(event.target.value)}
                className="h-9 w-10 rounded-full border border-border/80 bg-background"
              />
            </div>
            <Button type="submit" disabled={busy}>
              {editingTag ? <Check className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
              {editingTag ? "保存" : "创建"}
            </Button>
          </div>
          {formError ? <p className="text-sm text-destructive">{formError}</p> : null}
        </form>

        <div className="mt-4 space-y-2">
          {loading ? (
            Array.from({ length: 4 }, (_, index) => <Skeleton key={index} className="h-12 rounded-[20px]" />)
          ) : error ? (
            <div className="rounded-[20px] border border-destructive/20 bg-destructive/5 p-3 text-sm text-destructive">{error}</div>
          ) : sortedTags.length ? (
            sortedTags.map((tag) => (
              <div key={tag.id} className="flex items-center justify-between gap-3 rounded-[20px] border border-border/70 bg-white/80 p-3">
                <div className="flex min-w-0 items-center gap-3">
                  <span
                    className="h-4 w-4 shrink-0 rounded-full border"
                    style={tag.color ? { backgroundColor: tag.color, borderColor: tag.color } : undefined}
                  />
                  <div className="min-w-0">
                    <div className="truncate text-sm font-medium">{tag.name}</div>
                    <div className="font-mono text-xs text-muted-foreground">{tag.color ?? "no color"}</div>
                  </div>
                </div>
                <div className="flex shrink-0 gap-1">
                  <Button
                    aria-label={`编辑标签 ${tag.name}`}
                    variant="ghost"
                    size="icon"
                    onClick={() => startEdit(tag)}
                    className="rounded-full border border-border/70 bg-white/90 shadow-sm shadow-black/5"
                  >
                    <Edit2 className="h-4 w-4" />
                  </Button>
                  <Button
                    aria-label={`删除标签 ${tag.name}`}
                    variant="ghost"
                    size="icon"
                    onClick={() => setDeleteTarget(tag)}
                    className="rounded-full border border-border/70 bg-white/90 shadow-sm shadow-black/5"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            ))
          ) : (
            <div className="rounded-[20px] border border-dashed border-border/70 p-6 text-center text-sm text-muted-foreground">
              暂无标签
            </div>
          )}
        </div>
      </Dialog>

      <ConfirmDialog
        open={Boolean(deleteTarget)}
        title="删除标签"
        description={deleteTarget ? `确认删除标签「${deleteTarget.name}」？` : ""}
        confirmLabel="删除"
        busy={busy}
        onCancel={() => setDeleteTarget(null)}
        onConfirm={confirmDelete}
      />
    </>
  )
}
