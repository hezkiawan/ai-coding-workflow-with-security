'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/lib/auth-context';
import {
  LayoutDashboard,
  Ticket,
  Server,
  Shield,
  Settings,
  PlusCircle,
} from 'lucide-react';

export const Sidebar: React.FC = () => {
  const pathname = usePathname();
  const { user } = useAuth();

  const links = [
    { href: '/', label: 'Overview', icon: LayoutDashboard },
    { href: '/tickets', label: 'Service Tickets', icon: Ticket },
    { href: '/assets', label: 'IT Assets', icon: Server },
    ...(user?.role === 'admin'
      ? [{ href: '/admin', label: 'Admin Center', icon: Shield }]
      : []),
    { href: '/settings', label: 'Settings', icon: Settings },
  ];

  return (
    <aside className="w-64 border-r border-slate-200 bg-white min-h-[calc(100vh-4rem)] p-4 flex flex-col justify-between">
      <div className="space-y-6">
        <div>
          <Link
            href="/tickets?new=true"
            className="w-full flex items-center justify-center space-x-2 bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2.5 rounded-lg text-sm font-semibold transition shadow-sm"
          >
            <PlusCircle className="w-4 h-4" />
            <span>New Ticket</span>
          </Link>
        </div>

        <nav className="space-y-1">
          <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider px-3 mb-2">
            Workspace
          </div>
          {links.map((link) => {
            const Icon = link.icon;
            const isActive =
              link.href === '/'
                ? pathname === '/'
                : pathname.startsWith(link.href);

            return (
              <Link
                key={link.href}
                href={link.href}
                className={`flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium transition ${
                  isActive
                    ? 'bg-emerald-50 text-emerald-700 font-semibold'
                    : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
                }`}
              >
                <Icon className={`w-4 h-4 ${isActive ? 'text-emerald-600' : 'text-slate-400'}`} />
                <span>{link.label}</span>
              </Link>
            );
          })}
        </nav>
      </div>

      <div className="p-3 bg-slate-50 rounded-lg border border-slate-200 text-xs text-slate-500">
        <p className="font-semibold text-slate-700">Internal Operations</p>
        <p className="mt-0.5">Encrypted internal helpdesk & infrastructure directory.</p>
      </div>
    </aside>
  );
};
