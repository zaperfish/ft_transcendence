'use client';

import { BackButton } from "@/components/ui/BackButton"
import { useTheme } from "@/lib/context/ThemeContext"
import { cn } from "@/lib/utils"

export default function TermsPage() {
  const { theme } = useTheme()
  const isClassic = theme === "classic"

  return (
    <main
      className={cn(
        "min-h-screen px-4 py-12",
        isClassic ? "bg-surface-dim" : undefined,
      )}
    >
      <div
        className={cn(
          "mx-auto max-w-4xl",
          isClassic
            ? "rounded-xl border border-border bg-background p-6 shadow-sm sm:p-8"
            : "px-xl py-2xl",
        )}
      >
        <div className="mb-6">
          <BackButton />
        </div>

        <h1
          className={cn(
            isClassic
              ? "text-3xl font-semibold tracking-tight"
              : "font-heading text-4xl font-bold text-chrome-title",
          )}
        >
          Terms of Service
        </h1>
        <p
          className={cn(
            "mt-4 text-sm",
            isClassic ? "text-muted-foreground" : "text-chrome-muted",
          )}
        >
          Last updated: June 22, 2026
        </p>

        <div
          className={cn(
            "mt-8 space-y-6 text-sm leading-6",
            isClassic ? "text-foreground" : "text-chrome-body",
          )}
        >
          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              Acceptance of terms
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              By using Camaraderie, you agree to follow these Terms of Service and
              use the platform in a lawful and respectful manner.
            </p>
          </section>

          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              User responsibilities
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              You are responsible for the content you share, including your
              profile information, event details, and messages. You must not post
              content that is harmful, misleading, offensive, illegal, or violates
              the rights of others.
            </p>
          </section>

          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              Moderation and account actions
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              We reserve the right to remove content, restrict access, or suspend
              accounts if we believe they violate these terms or may harm other
              users or the platform.
            </p>
          </section>

          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              Service availability
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              The service is provided as-is and may occasionally be unavailable
              due to maintenance or technical issues. We do not guarantee
              uninterrupted access or error-free performance.
            </p>
          </section>

          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              Changes to terms
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              We may update these Terms of Service from time to time. Continued
              use of the platform after changes are made means you accept the
              updated terms.
            </p>
          </section>
        </div>
      </div>
    </main>
  );
}
