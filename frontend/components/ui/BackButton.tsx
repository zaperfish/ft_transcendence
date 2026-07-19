"use client"

import { useRouter } from "next/navigation"

export function BackButton() {
  const router = useRouter()

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
      className="text-sm text-chrome-muted transition-colors hover:text-chrome-title"
    >
      ← Back
    </button>
  )
}
