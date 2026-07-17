import Link from "next/link";
import { BackButton } from "@/components/ui/BackButton"

export default function PrivacyPage() {
  return (
    <main className="min-h-screen bg-surface-dim px-4 py-12">
      <div className="mx-auto max-w-4xl rounded-xl border border-border bg-background p-6 shadow-sm sm:p-8">
        <div className="mb-6">
          <BackButton />
        </div>

        <h1 className="text-3xl font-semibold tracking-tight">Privacy Policy</h1>
        <p className="mt-4 text-sm text-muted-foreground">
          Last updated: June 22, 2026
        </p>

        <div className="mt-8 space-y-6 text-sm leading-6 text-foreground">
          <section>
            <h2 className="text-base font-semibold">Information we collect</h2>
            <p>
              We collect the information needed to provide and improve our service,
              including your account details, profile information, event activity,
              and messages you send through the platform.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold">How we use your information</h2>
            <p>
              We use your information to create and manage your account, allow you
              to join and organize events, improve security and performance, and
              send important service updates.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold">Sharing of information</h2>
            <p>
              We do not sell your personal information. We may share information
              when required by law, to protect user safety, or to enforce our
              rules and policies.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold">Your choices</h2>
            <p>
              You may request access to, correction of, or deletion of your
              personal information by contacting us. We will do our best to
              respond promptly and help resolve your request.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold">Changes to this policy</h2>
            <p>
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
