'use client';

import React from 'react';
import Link from 'next/link';
import { useAuth } from '@/lib/auth-context';
import { Badge } from './Badge';
import { ShieldCheck, User as UserIcon, LogOut } from 'lucide-react';

export const Navbar: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <header className="h-16 border-b border-slate-200 bg-white sticky top-0 z-30 flex items-center justify-between px-6">
      <div className="flex items-center space-x-3">
        <Link href="/" className="flex items-center space-x-2">
          <div className="w-8 h-8 rounded-lg bg-emerald-600 flex items-center justify-center text-white font-bold">
            <ShieldCheck className="w-5 h-5" />
          </div>
          <span className="text-xl font-bold tracking-tight text-slate-900">OpsDesk</span>
          <span className="text-xs px-2 py-0.5 rounded bg-slate-100 text-slate-600 font-mono">v2.4</span>
        </Link>
      </div>

      <div className="flex items-center space-x-4">
        {user ? (
          <div className="flex items-center space-x-4">
            <div className="text-right">
              <div className="text-sm font-semibold text-slate-800">{user.full_name}</div>
              <div className="flex items-center space-x-1.5 justify-end">
                <span className="text-xs text-slate-500">{user.department}</span>
                <Badge variant={user.role}>{user.role}</Badge>
              </div>
            </div>
            <Link
              href="/settings"
              className="p-2 rounded-full hover:bg-slate-100 text-slate-600 transition"
              title="Profile Settings"
            >
              <UserIcon className="w-5 h-5" />
            </Link>
            <button
              onClick={logout}
              className="p-2 rounded-full hover:bg-rose-50 hover:text-rose-600 text-slate-600 transition"
              title="Logout"
            >
              <LogOut className="w-5 h-5" />
            </button>
          </div>
        ) : (
          <div className="flex items-center space-x-3">
            <Link
              href="/login"
              className="text-sm font-medium text-slate-700 hover:text-slate-900 px-3 py-1.5 rounded-md hover:bg-slate-100 transition"
            >
              Sign In
            </Link>
            <Link
              href="/register"
              className="text-sm font-medium text-white bg-emerald-600 hover:bg-emerald-700 px-3.5 py-1.5 rounded-md transition shadow-sm"
            >
              Register
            </Link>
          </div>
        )}
      </div>
    </header>
  );
};
