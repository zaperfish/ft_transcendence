import Link from "next/link";
import { BackButton } from "@/components/ui/BackButton"

export default function TermsPage() {
  return (
    <main className="min-h-screen bg-surface-dim px-4 py-12">
      <div className="mx-auto max-w-4xl rounded-xl border border-border bg-background p-6 shadow-sm sm:p-8">
        <div className="mb-6">
          <BackButton />
        </div>

        <h1 className="text-3xl font-semibold tracking-tight">Terms of Service</h1>
        <p className="mt-4 text-sm text-muted-foreground">
          Last updated: June 22, 2026
        </p>

        <div className="mt-8 space-y-6 text-sm leading-6 text-foreground">
          <section>
            <h2 className="text-base font-semibold">Acceptance of terms</h2>
            <p>
              By using Camaraderie, you agree to follow these Terms of Service and
              use the platform in a lawful and respectful manner.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold">User responsibilities</h2>
            <p>
              You are responsible for the content you share, including your
              profile information, event details, and messages. You must not post
              content that is harmful, misleading, offensive, illegal, or violates
              the rights of others.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold">Moderation and account actions</h2>
            <p>
              We reserve the right to remove content, restrict access, or suspend
              accounts if we believe they violate these terms or may harm other
              users or the platform.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold">Service availability</h2>
            <p>
              The service is provided as-is and may occasionally be unavailable
              due to maintenance or technical issues. We do not guarantee
              uninterrupted access or error-free performance.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold">Changes to terms</h2>
            <p>
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
