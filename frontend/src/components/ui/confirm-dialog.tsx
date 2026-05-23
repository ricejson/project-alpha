import { Button } from "@/components/ui/button"
import { Dialog } from "@/components/ui/dialog"

interface ConfirmDialogProps {
  open: boolean
  title: string
  description: string
  confirmLabel?: string
  busy?: boolean
  onCancel: () => void
  onConfirm: () => void
}

export function ConfirmDialog({
  open,
  title,
  description,
  confirmLabel = "删除",
  busy,
  onCancel,
  onConfirm,
}: ConfirmDialogProps) {
  return (
    <Dialog open={open} title={title} description={description} onClose={onCancel} className="max-w-md">
      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onCancel} disabled={busy}>
          取消
        </Button>
        <Button variant="destructive" onClick={onConfirm} disabled={busy}>
          {busy ? "处理中..." : confirmLabel}
        </Button>
      </div>
    </Dialog>
  )
}
