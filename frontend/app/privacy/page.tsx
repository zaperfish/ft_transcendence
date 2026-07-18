import { BackButton } from "@/components/ui/BackButton"

export default function PrivacyPage() {
  return (
    <main className="min-h-screen px-4 py-12">
      <div className="mx-auto max-w-4xl px-xl py-2xl">
        <div className="mb-6">
          <BackButton />
        </div>

        <h1 className="font-heading text-4xl font-bold text-teal-50">
          Privacy Policy
        </h1>
        <p className="mt-4 text-sm text-teal-100/60">
          Last updated: June 22, 2026
        </p>

        <div className="mt-8 space-y-6 text-sm leading-6 text-teal-100/80">
          <section>
            <h2 className="text-base font-semibold text-teal-50">Information we collect</h2>
            <p className="mt-2">
              We collect the information needed to provide and improve our service,
              including your account details, profile information, event activity,
              and messages you send through the platform.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-teal-50">How we use your information</h2>
            <p className="mt-2">
              We use your information to create and manage your account, allow you
              to join and organize events, improve security and performance, and
              send important service updates.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-teal-50">Sharing of information</h2>
            <p className="mt-2">
              We do not sell your personal information. We may share information
              when required by law, to protect user safety, or to enforce our
              rules and policies.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-teal-50">Your choices</h2>
            <p className="mt-2">
              You may request access to, correction of, or deletion of your
              personal information by contacting us. We will do our best to
              respond promptly and help resolve your request.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-teal-50">Changes to this policy</h2>
            <p className="mt-2">
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
