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
      className="text-sm text-teal-100/70 transition-colors hover:text-teal-50"
    >
      ← Back
    </button>
  )
}
