"use client"

import { useRouter } from "next/navigation"
import { useTheme } from "@/lib/context/ThemeContext"
import { cn } from "@/lib/utils"

export function BackButton() {
  const router = useRouter()
  const { theme } = useTheme()
  const isClassic = theme === "classic"

  function handleBack() {
    if (window.history.length > 1) {
      router.back()
    } else {
      router.replace("/")
    }
  }

  return (
    <button
      type="button"
      onClick={handleBack}
      className={cn(
        "text-sm transition-colors",
        isClassic
          ? "text-muted-foreground hover:text-foreground"
          : "text-chrome-muted hover:text-chrome-title",
      )}
    >
      ← Back
    </button>
  )
}
