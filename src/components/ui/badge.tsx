import { cn } from "@/lib/utils"
import type { ReactNode } from "react"

interface BadgeProps {
  children: ReactNode
  className?: string
  color?: string | null
}

export function Badge({ children, className, color }: BadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex h-6 max-w-full items-center gap-1 rounded-md border px-2 text-xs font-medium tracking-normal",
        className,
      )}
      style={color ? { borderColor: color, backgroundColor: `${color}14`, color } : undefined}
    >
      {children}
    </span>
  )
}
