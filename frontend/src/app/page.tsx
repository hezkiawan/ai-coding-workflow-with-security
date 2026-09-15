'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { api, Ticket, Asset } from '@/lib/api';
import { Badge } from '@/components/Badge';
import {
  Ticket as TicketIcon,
  Server,
  Clock,
  CheckCircle2,
  AlertCircle,
  ArrowRight,
} from 'lucide-react';

export default function DashboardPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadData() {
      const [ticketsRes, assetsRes] = await Promise.all([
        api.getTickets(),
        api.getAssets(),
      ]);

      if (ticketsRes.success && ticketsRes.data) {
        setTickets(ticketsRes.data);
      }
      if (assetsRes.success && assetsRes.data) {
        setAssets(assetsRes.data);
      }
      setLoading(false);
    }
    loadData();
  }, []);

  const openTickets = tickets.filter((t) => t.status === 'open');
  const inProgressTickets = tickets.filter((t) => t.status === 'in_progress');
  const resolvedTickets = tickets.filter((t) => t.status === 'resolved');

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Operations Overview</h1>
        <p className="text-sm text-slate-500 mt-1">
          Real-time summary of support tickets, active hardware, and operational tasks.
        </p>
      </div>

      {/* Metrics Row */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-5">
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Open Tickets</p>
            <p className="text-2xl font-bold text-slate-900 mt-1">{openTickets.length}</p>
          </div>
          <div className="p-3 bg-emerald-50 rounded-lg text-emerald-600">
            <AlertCircle className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider">In Progress</p>
            <p className="text-2xl font-bold text-slate-900 mt-1">{inProgressTickets.length}</p>
          </div>
          <div className="p-3 bg-amber-50 rounded-lg text-amber-600">
            <Clock className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Resolved</p>
            <p className="text-2xl font-bold text-slate-900 mt-1">{resolvedTickets.length}</p>
          </div>
          <div className="p-3 bg-blue-50 rounded-lg text-blue-600">
            <CheckCircle2 className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Tracked Assets</p>
            <p className="text-2xl font-bold text-slate-900 mt-1">{assets.length}</p>
          </div>
          <div className="p-3 bg-purple-50 rounded-lg text-purple-600">
            <Server className="w-6 h-6" />
          </div>
        </div>
      </div>

      {/* Two Column Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Recent Tickets */}
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
          <div className="p-5 border-b border-slate-100 flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <TicketIcon className="w-5 h-5 text-emerald-600" />
              <h2 className="font-semibold text-slate-900">Recent Service Tickets</h2>
            </div>
            <Link href="/tickets" className="text-xs font-medium text-emerald-600 hover:text-emerald-700 flex items-center space-x-1">
              <span>View all</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </Link>
          </div>

          <div className="divide-y divide-slate-100">
            {loading ? (
              <div className="p-6 text-center text-sm text-slate-400">Loading tickets...</div>
            ) : tickets.slice(0, 5).map((t) => (
              <Link
                key={t.id}
                href={`/tickets/${t.id}`}
                className="p-4 flex items-center justify-between hover:bg-slate-50 transition block"
              >
                <div className="space-y-1 pr-4">
                  <div className="text-sm font-medium text-slate-900 line-clamp-1">{t.title}</div>
                  <div className="text-xs text-slate-500 flex items-center space-x-2">
                    <span>{t.department}</span>
                    <span>•</span>
                    <span>{new Date(t.created_at).toLocaleDateString()}</span>
                  </div>
                </div>
                <div className="flex items-center space-x-2 shrink-0">
                  <Badge variant={t.priority}>{t.priority}</Badge>
                  <Badge variant={t.status}>{t.status}</Badge>
                </div>
              </Link>
            ))}
          </div>
        </div>

        {/* Assets Overview */}
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
          <div className="p-5 border-b border-slate-100 flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Server className="w-5 h-5 text-emerald-600" />
              <h2 className="font-semibold text-slate-900">Infrastructure & Assets</h2>
            </div>
            <Link href="/assets" className="text-xs font-medium text-emerald-600 hover:text-emerald-700 flex items-center space-x-1">
              <span>View all</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </Link>
          </div>

          <div className="divide-y divide-slate-100">
            {loading ? (
              <div className="p-6 text-center text-sm text-slate-400">Loading assets...</div>
            ) : assets.slice(0, 5).map((a) => (
              <Link
                key={a.id}
                href={`/assets/${a.id}`}
                className="p-4 flex items-center justify-between hover:bg-slate-50 transition block"
              >
                <div className="space-y-1">
                  <div className="text-sm font-medium text-slate-900 flex items-center space-x-2">
                    <span>{a.name}</span>
                    <span className="text-xs px-1.5 py-0.5 bg-slate-100 rounded text-slate-600 font-mono">{a.asset_tag}</span>
                  </div>
                  <div className="text-xs text-slate-500">{a.category} — {a.location}</div>
                </div>
                <Badge variant={a.status}>{a.status}</Badge>
              </Link>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
