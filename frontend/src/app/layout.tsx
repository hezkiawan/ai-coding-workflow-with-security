import type { Metadata } from 'next';
import './globals.css';
import { AuthProvider } from '@/lib/auth-context';
import { Navbar } from '@/components/Navbar';
import { Sidebar } from '@/components/Sidebar';

export const metadata: Metadata = {
  title: 'OpsDesk — Internal IT Helpdesk & Asset Management',
  description: 'Enterprise IT operations and asset management portal',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-slate-50 min-h-screen">
        <AuthProvider>
          <div className="flex flex-col min-h-screen">
            <Navbar />
            <div className="flex flex-1">
              <Sidebar />
              <main className="flex-1 p-8 overflow-y-auto max-w-7xl">
                {children}
              </main>
            </div>
          </div>
        </AuthProvider>
      </body>
    </html>
  );
}
