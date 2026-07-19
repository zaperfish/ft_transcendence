'use client';

import { BackButton } from "@/components/ui/BackButton"
import { useTheme } from "@/lib/context/ThemeContext"
import { cn } from "@/lib/utils"

export default function PrivacyPage() {
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
          Privacy Policy
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
              Information we collect
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              We collect the information needed to provide and improve our service,
              including your account details, profile information, event activity,
              and messages you send through the platform.
            </p>
          </section>

          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              How we use your information
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              We use your information to create and manage your account, allow you
              to join and organize events, improve security and performance, and
              send important service updates.
            </p>
          </section>

          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              Sharing of information
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              We do not sell your personal information. We may share information
              when required by law, to protect user safety, or to enforce our
              rules and policies.
            </p>
          </section>

          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              Your choices
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              You may request access to, correction of, or deletion of your
              personal information by contacting us. We will do our best to
              respond promptly and help resolve your request.
            </p>
          </section>

          <section>
            <h2 className={cn("text-base font-semibold", !isClassic && "text-chrome-title")}>
              Changes to this policy
            </h2>
            <p className={cn(!isClassic && "mt-2")}>
              We may update this Privacy Policy from time to time. Continued use
              of the service after changes are made means you accept the updated
              policy.
            </p>
          </section>
        </div>
      </div>
    </main>
  );
}
