import { useCallback, useState } from "react"
import { TagManagerDialog } from "@/components/tags/TagManagerDialog"
import { Toasts, type ToastMessage } from "@/components/ui/toast"
import { TicketPage } from "@/components/tickets/TicketPage"
import { useTags } from "@/hooks/useTags"
import type { TagPayload } from "@/types/tag"

export function App() {
  const tags = useTags()
  const [tagManagerOpen, setTagManagerOpen] = useState(false)
  const [tagVersion, setTagVersion] = useState(0)
  const [toasts, setToasts] = useState<ToastMessage[]>([])

  const notify = useCallback((type: "success" | "error", message: string) => {
    const id = Date.now() + Math.floor(Math.random() * 1000)
    setToasts((current) => [...current, { id, type, message }])
    window.setTimeout(() => {
      setToasts((current) => current.filter((toast) => toast.id !== id))
    }, 3500)
  }, [])

  async function createTag(payload: TagPayload) {
    await tags.createTag(payload)
    setTagVersion((value) => value + 1)
    notify("success", "标签已创建")
  }

  async function updateTag(id: number, payload: TagPayload) {
    await tags.updateTag(id, payload)
    setTagVersion((value) => value + 1)
    notify("success", "标签已更新")
  }

  async function deleteTag(id: number) {
    await tags.deleteTag(id)
    setTagVersion((value) => value + 1)
    notify("success", "标签已删除")
  }

  return (
    <>
      <TicketPage
        tags={tags.tags}
        tagsLoading={tags.loading}
        tagsVersion={tagVersion}
        onManageTags={() => setTagManagerOpen(true)}
        notify={notify}
      />
      <TagManagerDialog
        open={tagManagerOpen}
        tags={tags.tags}
        loading={tags.loading}
        error={tags.error}
        onClose={() => setTagManagerOpen(false)}
        onCreate={createTag}
        onUpdate={updateTag}
        onDelete={deleteTag}
      />
      <Toasts messages={toasts} onDismiss={(id) => setToasts((current) => current.filter((toast) => toast.id !== id))} />
    </>
  )
}
