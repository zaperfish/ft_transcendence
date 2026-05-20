export default function MainLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <nav>
        {/* Navigation placeholder */}
      </nav>
      {children}
    </>
  );
}