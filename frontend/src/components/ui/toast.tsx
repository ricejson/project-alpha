import { X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

export interface ToastMessage {
  id: number
  type: "success" | "error"
  message: string
}

interface ToastsProps {
  messages: ToastMessage[]
  onDismiss: (id: number) => void
}

export function Toasts({ messages, onDismiss }: ToastsProps) {
  return (
    <div className="fixed bottom-4 right-4 z-[60] flex w-[min(24rem,calc(100vw-2rem))] flex-col gap-2">
      {messages.map((toast) => (
        <div
          key={toast.id}
          role="status"
          className={cn(
            "flex items-start justify-between gap-3 rounded-md border bg-background p-3 text-sm shadow-lg",
            toast.type === "error" ? "border-destructive/40" : "border-emerald-300",
          )}
        >
          <span>{toast.message}</span>
          <Button variant="ghost" size="icon" className="h-6 w-6" onClick={() => onDismiss(toast.id)}>
            <X className="h-3.5 w-3.5" />
          </Button>
        </div>
      ))}
    </div>
  )
}
