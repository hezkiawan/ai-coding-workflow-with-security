'use client';

import React, { useEffect, useState, Suspense } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { api, Ticket } from '@/lib/api';
import { Badge } from '@/components/Badge';
import {
  Search,
  Plus,
  Filter,
  X,
  AlertCircle,
  CheckCircle2,
} from 'lucide-react';

function TicketsContent() {
  const searchParams = useSearchParams();
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [departmentFilter, setDepartmentFilter] = useState('');
  const [priorityFilter, setPriorityFilter] = useState('');

  // Create Modal State
  const [showModal, setShowModal] = useState(searchParams.get('new') === 'true');
  const [newTitle, setNewTitle] = useState('');
  const [newDescription, setNewDescription] = useState('');
  const [newPriority, setNewPriority] = useState('medium');
  const [newDepartment, setNewDepartment] = useState('IT Operations');
  const [modalError, setModalError] = useState('');
  const [modalSuccess, setModalSuccess] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const loadTickets = async () => {
    setLoading(true);
    let res;
    if (searchQuery.trim()) {
      res = await api.searchTickets(searchQuery);
    } else {
      res = await api.getTickets({
        status: statusFilter || undefined,
        department: departmentFilter || undefined,
        priority: priorityFilter || undefined,
      });
    }

    if (res.success && res.data) {
      setTickets(res.data);
    }
    setLoading(false);
  };

  useEffect(() => {
    loadTickets();
  }, [statusFilter, departmentFilter, priorityFilter]);

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    loadTickets();
  };

  const handleCreateTicket = async (e: React.FormEvent) => {
    e.preventDefault();
    setModalError('');
    setModalSuccess('');
    setSubmitting(true);

    try {
      const res = await api.createTicket({
        title: newTitle,
        description: newDescription,
        priority: newPriority,
        department: newDepartment,
      });

      if (res.success) {
        setModalSuccess('Ticket created successfully');
        setNewTitle('');
        setNewDescription('');
        setTimeout(() => {
          setShowModal(false);
          setModalSuccess('');
          loadTickets();
        }, 1000);
      } else {
        setModalError(res.error || 'Failed to create ticket');
      }
    } catch {
      setModalError('Network error occurred');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Service Tickets</h1>
          <p className="text-sm text-slate-500 mt-1">
            Track and resolve incident tickets, hardware requests, and IT inquiries.
          </p>
        </div>

        <button
          onClick={() => setShowModal(true)}
          className="flex items-center space-x-2 bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2 rounded-lg text-sm font-semibold transition shadow-sm self-start"
        >
          <Plus className="w-4 h-4" />
          <span>New Ticket</span>
        </button>
      </div>

      {/* Filter and Search Bar */}
      <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col md:flex-row gap-3 items-center justify-between">
        <form onSubmit={handleSearchSubmit} className="relative w-full md:w-96">
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-3" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search tickets by title, description..."
            className="w-full pl-9 pr-4 py-2 text-sm rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
          />
        </form>

        <div className="flex items-center space-x-2 w-full md:w-auto overflow-x-auto">
          <Filter className="w-4 h-4 text-slate-400 shrink-0" />
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="text-xs py-1.5 px-2.5 rounded-lg border border-slate-200 bg-white text-slate-700 focus:outline-none"
          >
            <option value="">All Statuses</option>
            <option value="open">Open</option>
            <option value="in_progress">In Progress</option>
            <option value="resolved">Resolved</option>
            <option value="closed">Closed</option>
          </select>

          <select
            value={priorityFilter}
            onChange={(e) => setPriorityFilter(e.target.value)}
            className="text-xs py-1.5 px-2.5 rounded-lg border border-slate-200 bg-white text-slate-700 focus:outline-none"
          >
            <option value="">All Priorities</option>
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="urgent">Urgent</option>
          </select>

          <select
            value={departmentFilter}
            onChange={(e) => setDepartmentFilter(e.target.value)}
            className="text-xs py-1.5 px-2.5 rounded-lg border border-slate-200 bg-white text-slate-700 focus:outline-none"
          >
            <option value="">All Departments</option>
            <option value="IT Operations">IT Operations</option>
            <option value="Infrastructure">Infrastructure</option>
            <option value="Finance">Finance</option>
            <option value="Engineering">Engineering</option>
          </select>
        </div>
      </div>

      {/* Tickets List */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="divide-y divide-slate-100">
          {loading ? (
            <div className="p-8 text-center text-sm text-slate-400">Loading service tickets...</div>
          ) : tickets.length === 0 ? (
            <div className="p-8 text-center text-sm text-slate-400">No tickets found matching current criteria.</div>
          ) : (
            tickets.map((ticket) => (
              <Link
                key={ticket.id}
                href={`/tickets/${ticket.id}`}
                className="p-5 flex flex-col md:flex-row md:items-center justify-between gap-3 hover:bg-slate-50 transition block"
              >
                <div className="space-y-1.5 flex-1 pr-4">
                  <div className="flex items-center space-x-2">
                    <span className="text-xs font-mono text-slate-400">#{ticket.id}</span>
                    <h2 className="text-base font-semibold text-slate-900 line-clamp-1">{ticket.title}</h2>
                  </div>
                  <p className="text-xs text-slate-500 line-clamp-1">{ticket.description}</p>
                  <div className="text-xs text-slate-400 flex items-center space-x-3 pt-1">
                    <span>Requested by <strong className="text-slate-600">{ticket.creator?.full_name || 'User'}</strong></span>
                    <span>•</span>
                    <span>Dept: {ticket.department}</span>
                    <span>•</span>
                    <span>{new Date(ticket.created_at).toLocaleDateString()}</span>
                  </div>
                </div>

                <div className="flex items-center space-x-2 shrink-0">
                  <Badge variant={ticket.priority}>{ticket.priority}</Badge>
                  <Badge variant={ticket.status}>{ticket.status}</Badge>
                </div>
              </Link>
            ))
          )}
        </div>
      </div>

      {/* Create Ticket Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl border border-slate-200 shadow-xl max-w-lg w-full p-6 relative">
            <button
              onClick={() => setShowModal(false)}
              className="absolute right-4 top-4 text-slate-400 hover:text-slate-600"
            >
              <X className="w-5 h-5" />
            </button>

            <h2 className="text-lg font-bold text-slate-900 mb-1">Create Service Request</h2>
            <p className="text-xs text-slate-500 mb-4">Submit a new incident or technical support ticket</p>

            {modalError && (
              <div className="mb-4 p-3 rounded-lg bg-rose-50 border border-rose-200 text-xs text-rose-700 flex items-center space-x-2">
                <AlertCircle className="w-4 h-4 shrink-0" />
                <span>{modalError}</span>
              </div>
            )}

            {modalSuccess && (
              <div className="mb-4 p-3 rounded-lg bg-emerald-50 border border-emerald-200 text-xs text-emerald-700 flex items-center space-x-2">
                <CheckCircle2 className="w-4 h-4 shrink-0" />
                <span>{modalSuccess}</span>
              </div>
            )}

            <form onSubmit={handleCreateTicket} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                  Ticket Subject
                </label>
                <input
                  type="text"
                  value={newTitle}
                  onChange={(e) => setNewTitle(e.target.value)}
                  required
                  placeholder="e.g. Broken Monitor in Conference Room B"
                  className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                  Detailed Description
                </label>
                <textarea
                  value={newDescription}
                  onChange={(e) => setNewDescription(e.target.value)}
                  required
                  rows={4}
                  placeholder="Provide detailed description of the problem or request..."
                  className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                    Priority
                  </label>
                  <select
                    value={newPriority}
                    onChange={(e) => setNewPriority(e.target.value)}
                    className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="low">Low</option>
                    <option value="medium">Medium</option>
                    <option value="high">High</option>
                    <option value="urgent">Urgent</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                    Department
                  </label>
                  <select
                    value={newDepartment}
                    onChange={(e) => setNewDepartment(e.target.value)}
                    className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="IT Operations">IT Operations</option>
                    <option value="Infrastructure">Infrastructure</option>
                    <option value="Finance">Finance</option>
                    <option value="Engineering">Engineering</option>
                  </select>
                </div>
              </div>

              <div className="pt-2 flex justify-end space-x-2">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-100 rounded-lg transition"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="px-4 py-2 text-sm font-semibold text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg transition shadow-sm disabled:opacity-50"
                >
                  {submitting ? 'Submitting...' : 'Create Ticket'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default function TicketsPage() {
  return (
    <Suspense fallback={<div className="p-8 text-center text-sm text-slate-400">Loading...</div>}>
      <TicketsContent />
    </Suspense>
  );
}
