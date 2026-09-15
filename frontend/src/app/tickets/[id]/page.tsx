'use client';

import React, { useEffect, useState, useRef } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { api, Ticket, Comment, Attachment } from '@/lib/api';
import { useAuth } from '@/lib/auth-context';
import { Badge } from '@/components/Badge';
import {
  ArrowLeft,
  MessageSquare,
  Paperclip,
  Clock,
  User,
  Building,
  Send,
  Upload,
  AlertCircle,
  CheckCircle2,
  Trash2,
  FileText,
} from 'lucide-react';

export default function TicketDetailPage() {
  const params = useParams();
  const router = useRouter();
  const ticketId = parseInt(params.id as string, 10);
  const { user } = useAuth();

  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [attachments, setAttachments] = useState<Attachment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // Comment input
  const [commentContent, setCommentContent] = useState('');
  const [isInternalComment, setIsInternalComment] = useState(false);
  const [submittingComment, setSubmittingComment] = useState(false);

  // Status/Priority updates
  const [status, setStatus] = useState('');
  const [priority, setPriority] = useState('');
  const [updatingTicket, setUpdatingTicket] = useState(false);
  const [updateMsg, setUpdateMsg] = useState('');

  // File upload
  const [uploading, setUploading] = useState(false);
  const [uploadMsg, setUploadMsg] = useState('');
  const fileInputRef = useRef<HTMLInputElement>(null);

  const loadTicketData = async () => {
    if (isNaN(ticketId)) {
      setError('Invalid ticket ID');
      setLoading(false);
      return;
    }

    try {
      const [ticketRes, commentsRes, attachmentsRes] = await Promise.all([
        api.getTicket(ticketId),
        api.getComments(ticketId),
        api.getAttachments(ticketId),
      ]);

      if (ticketRes.success && ticketRes.data) {
        setTicket(ticketRes.data);
        setStatus(ticketRes.data.status);
        setPriority(ticketRes.data.priority);
      } else {
        setError(ticketRes.error || 'Failed to load ticket');
      }

      if (commentsRes.success && commentsRes.data) {
        setComments(commentsRes.data);
      }

      if (attachmentsRes.success && attachmentsRes.data) {
        setAttachments(attachmentsRes.data);
      }
    } catch {
      setError('Network error loading ticket');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTicketData();
  }, [ticketId]);

  const handleStatusChange = async (newStatus: string) => {
    if (!ticket) return;
    setStatus(newStatus);
    setUpdatingTicket(true);
    setUpdateMsg('');

    const res = await api.updateTicket(ticket.id, { status: newStatus });
    if (res.success && res.data) {
      setTicket(res.data);
      setUpdateMsg('Status updated');
      setTimeout(() => setUpdateMsg(''), 2000);
    } else {
      setError(res.error || 'Failed to update status');
    }
    setUpdatingTicket(false);
  };

  const handlePriorityChange = async (newPriority: string) => {
    if (!ticket) return;
    setPriority(newPriority);
    setUpdatingTicket(true);
    setUpdateMsg('');

    const res = await api.updateTicket(ticket.id, { priority: newPriority });
    if (res.success && res.data) {
      setTicket(res.data);
      setUpdateMsg('Priority updated');
      setTimeout(() => setUpdateMsg(''), 2000);
    } else {
      setError(res.error || 'Failed to update priority');
    }
    setUpdatingTicket(false);
  };

  const handleAddComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!commentContent.trim()) return;
    setSubmittingComment(true);

    const res = await api.addComment(ticketId, {
      content: commentContent,
      is_internal: isInternalComment,
    });

    if (res.success && res.data) {
      setComments((prev) => [...prev, res.data!]);
      setCommentContent('');
      setIsInternalComment(false);
    }
    setSubmittingComment(false);
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    setUploadMsg('');

    const res = await api.uploadAttachment(ticketId, file);
    if (res.success && res.data) {
      setAttachments((prev) => [...prev, res.data!]);
      setUploadMsg('Attachment uploaded successfully');
      if (fileInputRef.current) fileInputRef.current.value = '';
      setTimeout(() => setUploadMsg(''), 2000);
    } else {
      setError(res.error || 'Failed to upload attachment');
    }
    setUploading(false);
  };

  const handleDeleteTicket = async () => {
    if (!confirm('Are you sure you want to delete this ticket?')) return;
    const res = await api.deleteTicket(ticketId);
    if (res.success) {
      router.push('/tickets');
    } else {
      setError(res.error || 'Failed to delete ticket');
    }
  };

  if (loading) {
    return <div className="p-8 text-center text-sm text-slate-400">Loading ticket details...</div>;
  }

  if (error && !ticket) {
    return (
      <div className="max-w-lg mx-auto p-6 bg-white rounded-xl border border-slate-200 text-center space-y-4">
        <AlertCircle className="w-10 h-10 text-rose-500 mx-auto" />
        <h2 className="text-lg font-bold text-slate-900">Error Loading Ticket</h2>
        <p className="text-sm text-slate-500">{error}</p>
        <Link
          href="/tickets"
          className="inline-flex items-center space-x-2 text-sm text-emerald-600 font-semibold hover:underline"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>Back to Tickets</span>
        </Link>
      </div>
    );
  }

  if (!ticket) return null;

  const canEdit = user?.role === 'admin' || user?.role === 'technician' || user?.id === ticket.created_by;
  const isStaff = user?.role === 'admin' || user?.role === 'technician';

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="space-y-1">
          <Link
            href="/tickets"
            className="inline-flex items-center space-x-1.5 text-xs text-slate-500 hover:text-slate-800 transition mb-1"
          >
            <ArrowLeft className="w-3.5 h-3.5" />
            <span>Back to All Tickets</span>
          </Link>
          <div className="flex items-center space-x-3">
            <span className="text-sm font-mono text-slate-400">#{ticket.id}</span>
            <h1 className="text-2xl font-bold text-slate-900">{ticket.title}</h1>
          </div>
        </div>

        {canEdit && (
          <div className="flex items-center space-x-2">
            {user?.role === 'admin' && (
              <button
                onClick={handleDeleteTicket}
                className="flex items-center space-x-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold text-rose-600 hover:bg-rose-50 border border-rose-200 transition"
              >
                <Trash2 className="w-3.5 h-3.5" />
                <span>Delete</span>
              </button>
            )}
          </div>
        )}
      </div>

      {updateMsg && (
        <div className="p-3 bg-emerald-50 border border-emerald-200 text-xs text-emerald-700 rounded-lg flex items-center space-x-2">
          <CheckCircle2 className="w-4 h-4 shrink-0" />
          <span>{updateMsg}</span>
        </div>
      )}

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left 2 Cols: Details & Comments */}
        <div className="lg:col-span-2 space-y-6">
          {/* Ticket Description Card */}
          <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-4">
            <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400">Description</h2>
            <div className="text-sm text-slate-700 whitespace-pre-wrap leading-relaxed">
              {ticket.description}
            </div>
          </div>

          {/* Comments Section */}
          <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-6">
            <div className="flex items-center justify-between border-b border-slate-100 pb-4">
              <div className="flex items-center space-x-2">
                <MessageSquare className="w-4 h-4 text-emerald-600" />
                <h2 className="font-semibold text-slate-900 text-sm">
                  Activity & Comments ({comments.length})
                </h2>
              </div>
            </div>

            {/* Comments List */}
            <div className="space-y-4">
              {comments.length === 0 ? (
                <p className="text-xs text-slate-400 italic py-2">No comments on this ticket yet.</p>
              ) : (
                comments.map((comment) => (
                  <div
                    key={comment.id}
                    className={`p-4 rounded-xl border text-xs space-y-2 ${
                      comment.is_internal
                        ? 'bg-amber-50/60 border-amber-200'
                        : 'bg-slate-50/70 border-slate-200'
                    }`}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <span className="font-semibold text-slate-800">
                          {comment.user?.full_name || 'Staff Member'}
                        </span>
                        <span className="text-slate-400 font-mono">({comment.user?.role || 'user'})</span>
                        {comment.is_internal && (
                          <span className="px-1.5 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider bg-amber-200 text-amber-900">
                            Internal Note
                          </span>
                        )}
                      </div>
                      <span className="text-slate-400 text-[11px]">
                        {new Date(comment.created_at).toLocaleString()}
                      </span>
                    </div>
                    <p className="text-slate-700 whitespace-pre-wrap leading-relaxed">
                      {comment.content}
                    </p>
                  </div>
                ))
              )}
            </div>

            {/* Add Comment Form */}
            <form onSubmit={handleAddComment} className="space-y-3 pt-4 border-t border-slate-100">
              <textarea
                value={commentContent}
                onChange={(e) => setCommentContent(e.target.value)}
                placeholder="Write a reply or status update..."
                rows={3}
                required
                className="w-full p-3 text-xs rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />

              <div className="flex items-center justify-between">
                {isStaff ? (
                  <label className="flex items-center space-x-2 text-xs text-slate-600 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={isInternalComment}
                      onChange={(e) => setIsInternalComment(e.target.checked)}
                      className="rounded text-emerald-600 focus:ring-emerald-500"
                    />
                    <span>Staff Internal Note (Hidden from requester)</span>
                  </label>
                ) : (
                  <div />
                )}

                <button
                  type="submit"
                  disabled={submittingComment}
                  className="flex items-center space-x-1.5 px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold rounded-lg transition disabled:opacity-50"
                >
                  <Send className="w-3.5 h-3.5" />
                  <span>{submittingComment ? 'Posting...' : 'Post Reply'}</span>
                </button>
              </div>
            </form>
          </div>
        </div>

        {/* Right 1 Col: Metadata & Attachments */}
        <div className="space-y-6">
          {/* Metadata Card */}
          <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm space-y-4">
            <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400 border-b border-slate-100 pb-3">
              Ticket Properties
            </h2>

            <div className="space-y-3 text-xs">
              <div>
                <label className="block text-slate-400 mb-1">Status</label>
                {canEdit ? (
                  <select
                    value={status}
                    onChange={(e) => handleStatusChange(e.target.value)}
                    disabled={updatingTicket}
                    className="w-full p-2 bg-slate-50 rounded-lg border border-slate-200 font-medium text-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="open">Open</option>
                    <option value="in_progress">In Progress</option>
                    <option value="resolved">Resolved</option>
                    <option value="closed">Closed</option>
                  </select>
                ) : (
                  <Badge variant={ticket.status}>{ticket.status}</Badge>
                )}
              </div>

              <div>
                <label className="block text-slate-400 mb-1">Priority</label>
                {canEdit ? (
                  <select
                    value={priority}
                    onChange={(e) => handlePriorityChange(e.target.value)}
                    disabled={updatingTicket}
                    className="w-full p-2 bg-slate-50 rounded-lg border border-slate-200 font-medium text-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="low">Low</option>
                    <option value="medium">Medium</option>
                    <option value="high">High</option>
                    <option value="urgent">Urgent</option>
                  </select>
                ) : (
                  <Badge variant={ticket.priority}>{ticket.priority}</Badge>
                )}
              </div>

              <div className="pt-2 border-t border-slate-100 space-y-2.5">
                <div className="flex items-center justify-between text-slate-600">
                  <span className="flex items-center space-x-1.5 text-slate-400">
                    <User className="w-3.5 h-3.5" />
                    <span>Requester:</span>
                  </span>
                  <span className="font-semibold text-slate-800">
                    {ticket.creator?.full_name || 'User'}
                  </span>
                </div>

                <div className="flex items-center justify-between text-slate-600">
                  <span className="flex items-center space-x-1.5 text-slate-400">
                    <Building className="w-3.5 h-3.5" />
                    <span>Department:</span>
                  </span>
                  <span className="font-semibold text-slate-800">{ticket.department}</span>
                </div>

                <div className="flex items-center justify-between text-slate-600">
                  <span className="flex items-center space-x-1.5 text-slate-400">
                    <Clock className="w-3.5 h-3.5" />
                    <span>Created:</span>
                  </span>
                  <span>{new Date(ticket.created_at).toLocaleDateString()}</span>
                </div>
              </div>
            </div>
          </div>

          {/* File Attachments Card */}
          <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm space-y-4">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <div className="flex items-center space-x-2">
                <Paperclip className="w-4 h-4 text-emerald-600" />
                <h2 className="font-semibold text-slate-900 text-xs uppercase tracking-wider">
                  Attachments ({attachments.length})
                </h2>
              </div>
            </div>

            {uploadMsg && (
              <div className="p-2 bg-emerald-50 text-[11px] text-emerald-700 rounded border border-emerald-200">
                {uploadMsg}
              </div>
            )}

            <div className="space-y-2">
              {attachments.length === 0 ? (
                <p className="text-xs text-slate-400 italic">No files attached.</p>
              ) : (
                attachments.map((att) => (
                  <div
                    key={att.id}
                    className="p-2.5 rounded-lg border border-slate-200 bg-slate-50 flex items-center justify-between text-xs hover:bg-slate-100 transition"
                  >
                    <div className="flex items-center space-x-2 overflow-hidden pr-2">
                      <FileText className="w-4 h-4 text-slate-400 shrink-0" />
                      <a
                        href={api.getAttachmentDownloadUrl(att.id)}
                        target="_blank"
                        rel="noreferrer"
                        className="truncate font-medium text-slate-800 hover:text-emerald-600 hover:underline"
                      >
                        {att.file_name}
                      </a>
                    </div>
                    <span className="text-[11px] text-slate-400 shrink-0">
                      {(att.file_size / 1024).toFixed(0)} KB
                    </span>
                  </div>
                ))
              )}
            </div>

            {/* Upload form */}
            <div className="pt-2 border-t border-slate-100">
              <input
                type="file"
                ref={fileInputRef}
                onChange={handleFileUpload}
                className="hidden"
                id="file-upload"
              />
              <label
                htmlFor="file-upload"
                className="w-full flex items-center justify-center space-x-2 py-2 px-3 border border-dashed border-slate-300 rounded-lg text-xs font-semibold text-slate-600 hover:bg-slate-50 hover:border-slate-400 cursor-pointer transition"
              >
                <Upload className="w-3.5 h-3.5 text-slate-400" />
                <span>{uploading ? 'Uploading...' : 'Attach File'}</span>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
