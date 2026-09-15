'use client';

import React, { useState } from 'react';
import { useAuth } from '@/lib/auth-context';
import { api } from '@/lib/api';
import { Badge } from '@/components/Badge';
import {
  User,
  Mail,
  Building,
  Shield,
  KeyRound,
  CheckCircle2,
  AlertCircle,
  Clock,
} from 'lucide-react';

export default function SettingsPage() {
  const { user } = useAuth();
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState('');
  const [error, setError] = useState('');

  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault();
    setMsg('');
    setError('');

    if (newPassword !== confirmPassword) {
      setError('New passwords do not match');
      return;
    }

    if (newPassword.length < 8) {
      setError('New password must be at least 8 characters');
      return;
    }

    setLoading(true);
    try {
      const res = await api.changePassword(oldPassword, newPassword);
      if (res.success) {
        setMsg('Password updated successfully');
        setOldPassword('');
        setNewPassword('');
        setConfirmPassword('');
      } else {
        setError(res.error || 'Failed to update password');
      }
    } catch {
      setError('An error occurred while changing password');
    } finally {
      setLoading(false);
    }
  };

  if (!user) {
    return <div className="p-8 text-center text-sm text-slate-400">Please sign in to view account settings.</div>;
  }

  return (
    <div className="max-w-3xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Account Settings</h1>
        <p className="text-sm text-slate-500 mt-1">
          Manage your internal profile identity, security credentials, and department assignment.
        </p>
      </div>

      {msg && (
        <div className="p-3 bg-emerald-50 border border-emerald-200 text-xs text-emerald-700 rounded-lg flex items-center space-x-2">
          <CheckCircle2 className="w-4 h-4 shrink-0" />
          <span>{msg}</span>
        </div>
      )}

      {error && (
        <div className="p-3 bg-rose-50 border border-rose-200 text-xs text-rose-700 rounded-lg flex items-center space-x-2">
          <AlertCircle className="w-4 h-4 shrink-0" />
          <span>{error}</span>
        </div>
      )}

      {/* User Profile Info */}
      <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-4">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400 border-b border-slate-100 pb-3">
          Employee Identity
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
          <div className="flex items-center space-x-3 p-3 bg-slate-50 rounded-lg border border-slate-100">
            <User className="w-5 h-5 text-slate-400" />
            <div>
              <p className="text-slate-400">Full Name</p>
              <p className="font-semibold text-slate-900">{user.full_name}</p>
            </div>
          </div>

          <div className="flex items-center space-x-3 p-3 bg-slate-50 rounded-lg border border-slate-100">
            <User className="w-5 h-5 text-slate-400" />
            <div>
              <p className="text-slate-400">Username</p>
              <p className="font-mono font-semibold text-slate-900">@{user.username}</p>
            </div>
          </div>

          <div className="flex items-center space-x-3 p-3 bg-slate-50 rounded-lg border border-slate-100">
            <Mail className="w-5 h-5 text-slate-400" />
            <div>
              <p className="text-slate-400">Corporate Email</p>
              <p className="font-semibold text-slate-900">{user.email}</p>
            </div>
          </div>

          <div className="flex items-center space-x-3 p-3 bg-slate-50 rounded-lg border border-slate-100">
            <Building className="w-5 h-5 text-slate-400" />
            <div>
              <p className="text-slate-400">Department</p>
              <p className="font-semibold text-slate-900">{user.department}</p>
            </div>
          </div>
        </div>

        <div className="pt-2 flex items-center space-x-2 text-xs">
          <span className="text-slate-500">System Role:</span>
          <Badge variant={user.role}>{user.role}</Badge>
        </div>
      </div>

      {/* Change Password */}
      <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-4">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400 border-b border-slate-100 pb-3 flex items-center space-x-2">
          <KeyRound className="w-4 h-4 text-emerald-600" />
          <span>Update Security Password</span>
        </h2>

        <form onSubmit={handlePasswordChange} className="space-y-4 max-w-md text-xs">
          <div>
            <label className="block font-semibold text-slate-700 uppercase tracking-wider mb-1">
              Current Password
            </label>
            <input
              type="password"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              required
              className="w-full p-2.5 rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              placeholder="••••••••"
            />
          </div>

          <div>
            <label className="block font-semibold text-slate-700 uppercase tracking-wider mb-1">
              New Password
            </label>
            <input
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              required
              minLength={8}
              className="w-full p-2.5 rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              placeholder="At least 8 characters"
            />
          </div>

          <div>
            <label className="block font-semibold text-slate-700 uppercase tracking-wider mb-1">
              Confirm New Password
            </label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              minLength={8}
              className="w-full p-2.5 rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              placeholder="••••••••"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="py-2.5 px-4 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold rounded-lg transition disabled:opacity-50 shadow-sm"
          >
            {loading ? 'Updating...' : 'Change Password'}
          </button>
        </form>
      </div>
    </div>
  );
}
