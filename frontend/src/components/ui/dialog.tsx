import { X } from "lucide-react"
import { useId } from "react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import type { ReactNode } from "react"

interface DialogProps {
  open: boolean
  title: string
  description?: string
  children: ReactNode
  onClose: () => void
  className?: string
}

export function Dialog({ open, title, description, children, onClose, className }: DialogProps) {
  const titleId = useId()

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/20 p-4" role="presentation">
      <section
        aria-modal="true"
        role="dialog"
        aria-labelledby={titleId}
        className={cn(
          "max-h-[calc(100vh-2rem)] w-full max-w-xl overflow-auto rounded-2xl border border-border/70 bg-white p-5 shadow-[0_24px_60px_rgba(15,23,42,0.16)]",
          className,
        )}
      >
        <div className="mb-5 flex items-start justify-between gap-4 border-b border-border/60 pb-4">
          <div>
            <h2 id={titleId} className="text-[1.1rem] font-semibold tracking-tight">
              {title}
            </h2>
            {description ? <p className="mt-1 text-sm text-muted-foreground">{description}</p> : null}
          </div>
          <Button aria-label="Close dialog" size="icon" variant="ghost" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>
        {children}
      </section>
    </div>
  )
}
