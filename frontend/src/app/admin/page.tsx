'use client';

import React, { useEffect, useState } from 'react';
import { api, User, AuditLog } from '@/lib/api';
import { useAuth } from '@/lib/auth-context';
import { Badge } from '@/components/Badge';
import {
  ShieldAlert,
  Users,
  FileText,
  Activity,
  CheckCircle2,
  AlertCircle,
  Clock,
  UserCheck,
} from 'lucide-react';

export default function AdminPage() {
  const { user } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'users' | 'logs' | 'stats'>('users');
  const [actionMsg, setActionMsg] = useState('');
  const [actionError, setActionError] = useState('');

  const loadAdminData = async () => {
    setLoading(true);
    const [usersRes, logsRes, statsRes] = await Promise.all([
      api.getAdminUsers(),
      api.getAuditLogs(),
      api.getAdminStats(),
    ]);

    if (usersRes.success && usersRes.data) {
      setUsers(usersRes.data);
    }
    if (logsRes.success && logsRes.data) {
      setAuditLogs(logsRes.data);
    }
    if (statsRes.success && statsRes.data) {
      setStats(statsRes.data);
    }
    setLoading(false);
  };

  useEffect(() => {
    if (user?.role === 'admin') {
      loadAdminData();
    }
  }, [user]);

  const handleRoleChange = async (userId: number, newRole: string) => {
    setActionMsg('');
    setActionError('');

    const res = await api.updateUserRole(userId, newRole);
    if (res.success && res.data) {
      setUsers((prev) =>
        prev.map((u) => (u.id === userId ? { ...u, role: newRole as User['role'] } : u))
      );
      setActionMsg(`Updated user role to ${newRole}`);
      setTimeout(() => setActionMsg(''), 2500);
    } else {
      setActionError(res.error || 'Failed to update user role');
    }
  };

  if (!user || user.role !== 'admin') {
    return (
      <div className="max-w-md mx-auto mt-12 p-8 bg-white rounded-2xl border border-slate-200 text-center space-y-4 shadow-sm">
        <div className="w-12 h-12 bg-rose-50 text-rose-600 rounded-xl flex items-center justify-center mx-auto">
          <ShieldAlert className="w-6 h-6" />
        </div>
        <h1 className="text-lg font-bold text-slate-900">Access Restricted</h1>
        <p className="text-xs text-slate-500">
          The Admin Control Center is restricted to authorized System Administrators only. Your current role is <strong>{user?.role || 'Guest'}</strong>.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Admin Control Center</h1>
          <p className="text-sm text-slate-500 mt-1">
            System configuration, user role assignments, and security audit log monitor.
          </p>
        </div>

        {/* Tab Switcher */}
        <div className="flex items-center bg-slate-100 p-1 rounded-lg self-start">
          <button
            onClick={() => setActiveTab('users')}
            className={`flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold transition ${
              activeTab === 'users'
                ? 'bg-white text-slate-900 shadow-sm'
                : 'text-slate-600 hover:text-slate-900'
            }`}
          >
            <Users className="w-3.5 h-3.5" />
            <span>Users</span>
          </button>
          <button
            onClick={() => setActiveTab('logs')}
            className={`flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold transition ${
              activeTab === 'logs'
                ? 'bg-white text-slate-900 shadow-sm'
                : 'text-slate-600 hover:text-slate-900'
            }`}
          >
            <FileText className="w-3.5 h-3.5" />
            <span>Audit Logs</span>
          </button>
          <button
            onClick={() => setActiveTab('stats')}
            className={`flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold transition ${
              activeTab === 'stats'
                ? 'bg-white text-slate-900 shadow-sm'
                : 'text-slate-600 hover:text-slate-900'
            }`}
          >
            <Activity className="w-3.5 h-3.5" />
            <span>System Stats</span>
          </button>
        </div>
      </div>

      {actionMsg && (
        <div className="p-3 bg-emerald-50 border border-emerald-200 text-xs text-emerald-700 rounded-lg flex items-center space-x-2">
          <CheckCircle2 className="w-4 h-4 shrink-0" />
          <span>{actionMsg}</span>
        </div>
      )}

      {actionError && (
        <div className="p-3 bg-rose-50 border border-rose-200 text-xs text-rose-700 rounded-lg flex items-center space-x-2">
          <AlertCircle className="w-4 h-4 shrink-0" />
          <span>{actionError}</span>
        </div>
      )}

      {/* Tab: Users Management */}
      {activeTab === 'users' && (
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-x-auto">
          <div className="p-4 border-b border-slate-100 flex items-center justify-between">
            <h2 className="font-semibold text-slate-900 text-sm">Registered Accounts ({users.length})</h2>
          </div>
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-50 text-slate-500 border-b border-slate-200 uppercase tracking-wider font-semibold">
              <tr>
                <th className="p-4">User</th>
                <th className="p-4">Email</th>
                <th className="p-4">Department</th>
                <th className="p-4">Current Role</th>
                <th className="p-4">Role Action</th>
                <th className="p-4">Created</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {loading ? (
                <tr>
                  <td colSpan={6} className="p-8 text-center text-slate-400">
                    Loading directory...
                  </td>
                </tr>
              ) : (
                users.map((u) => (
                  <tr key={u.id} className="hover:bg-slate-50 transition">
                    <td className="p-4">
                      <div>
                        <p className="font-semibold text-slate-900">{u.full_name}</p>
                        <p className="text-[11px] text-slate-400 font-mono">@{u.username}</p>
                      </div>
                    </td>
                    <td className="p-4 text-slate-600">{u.email}</td>
                    <td className="p-4 text-slate-600">{u.department}</td>
                    <td className="p-4">
                      <Badge variant={u.role}>{u.role}</Badge>
                    </td>
                    <td className="p-4">
                      <select
                        value={u.role}
                        onChange={(e) => handleRoleChange(u.id, e.target.value)}
                        disabled={u.id === user.id}
                        className="py-1 px-2 border border-slate-200 rounded text-xs bg-white text-slate-800 focus:outline-none focus:ring-1 focus:ring-emerald-500"
                      >
                        <option value="employee">Employee</option>
                        <option value="technician">Technician</option>
                        <option value="admin">Administrator</option>
                      </select>
                    </td>
                    <td className="p-4 text-slate-500">
                      {new Date(u.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Tab: Audit Logs */}
      {activeTab === 'logs' && (
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-x-auto">
          <div className="p-4 border-b border-slate-100 flex items-center justify-between">
            <h2 className="font-semibold text-slate-900 text-sm">Security & Audit Event Stream</h2>
          </div>
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-50 text-slate-500 border-b border-slate-200 uppercase tracking-wider font-semibold">
              <tr>
                <th className="p-4">Timestamp</th>
                <th className="p-4">Actor ID</th>
                <th className="p-4">Action</th>
                <th className="p-4">Target Resource</th>
                <th className="p-4">Details</th>
                <th className="p-4">IP Address</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 font-mono">
              {loading ? (
                <tr>
                  <td colSpan={6} className="p-8 text-center text-slate-400">
                    Loading audit stream...
                  </td>
                </tr>
              ) : auditLogs.length === 0 ? (
                <tr>
                  <td colSpan={6} className="p-8 text-center text-slate-400">
                    No audit records logged yet.
                  </td>
                </tr>
              ) : (
                auditLogs.map((log) => (
                  <tr key={log.id} className="hover:bg-slate-50 transition">
                    <td className="p-4 text-slate-500">
                      {new Date(log.created_at).toLocaleString()}
                    </td>
                    <td className="p-4 text-slate-700">
                      {log.user_id ? `User #${log.user_id}` : 'System / Unauth'}
                    </td>
                    <td className="p-4 font-semibold text-slate-900">{log.action}</td>
                    <td className="p-4 text-slate-700">
                      {log.resource_type} {log.resource_id ? `#${log.resource_id}` : ''}
                    </td>
                    <td className="p-4 text-slate-600 font-sans text-xs">{log.details}</td>
                    <td className="p-4 text-slate-500">{log.ip_address}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Tab: System Stats */}
      {activeTab === 'stats' && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-2">
            <h3 className="text-xs font-bold uppercase tracking-wider text-slate-400">User Base</h3>
            <p className="text-3xl font-bold text-slate-900">{stats?.total_users ?? '-'}</p>
            <p className="text-xs text-slate-500">Active accounts in SQLite database</p>
          </div>

          <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-2">
            <h3 className="text-xs font-bold uppercase tracking-wider text-slate-400">Service Tickets</h3>
            <p className="text-3xl font-bold text-slate-900">{stats?.total_tickets ?? '-'}</p>
            <p className="text-xs text-slate-500">Incident & request records</p>
          </div>

          <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-2">
            <h3 className="text-xs font-bold uppercase tracking-wider text-slate-400">Tracked Assets</h3>
            <p className="text-3xl font-bold text-slate-900">{stats?.total_assets ?? '-'}</p>
            <p className="text-xs text-slate-500">Hardware & network appliances</p>
          </div>
        </div>
      )}
    </div>
  );
}
