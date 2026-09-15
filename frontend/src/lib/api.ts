const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  department: string;
  role: 'admin' | 'technician' | 'employee';
  created_at: string;
}

export interface Ticket {
  id: number;
  title: string;
  description: string;
  status: 'open' | 'in_progress' | 'resolved' | 'closed';
  priority: 'low' | 'medium' | 'high' | 'urgent';
  department: string;
  created_by: number;
  assigned_to?: number | null;
  created_at: string;
  updated_at: string;
  creator?: User;
  assignee?: User;
}

export interface Comment {
  id: number;
  ticket_id: number;
  user_id: number;
  content: string;
  is_internal: boolean;
  created_at: string;
  user?: User;
}

export interface Asset {
  id: number;
  asset_tag: string;
  name: string;
  category: string;
  serial_number: string;
  status: 'in_service' | 'maintenance' | 'retired' | 'available';
  location: string;
  department: string;
  assigned_to?: number | null;
  created_at: string;
  updated_at: string;
}

export interface Attachment {
  id: number;
  ticket_id: number;
  user_id: number;
  file_name: string;
  file_size: number;
  content_type: string;
  file_path: string;
  created_at: string;
}

export interface AuditLog {
  id: number;
  user_id?: number | null;
  action: string;
  resource_type: string;
  resource_id: number;
  details: string;
  ip_address: string;
  created_at: string;
}

export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

function getAuthHeader(): Record<string, string> {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      return { Authorization: `Bearer ${token}` };
    }
  }
  return {};
}

export async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const url = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint : `/${endpoint}`}`;
  const headers = {
    'Content-Type': 'application/json',
    ...getAuthHeader(),
    ...(options.headers || {}),
  };

  try {
    const res = await fetch(url, { credentials: 'include', ...options, headers });
    const json = await res.json();
    return json;
  } catch (err: any) {
    return {
      success: false,
      error: err.message || 'Network communication error',
    };
  }
}

export const api = {
  // Auth
  register: (body: any) => fetchApi<any>('/auth/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (body: any) => fetchApi<any>('/auth/login', { method: 'POST', body: JSON.stringify(body) }),
  logout: () => fetchApi<any>('/auth/logout', { method: 'POST' }),
  me: () => fetchApi<User>('/auth/me'),
  updateProfile: (body: any) => fetchApi<any>('/auth/profile', { method: 'PUT', body: JSON.stringify(body) }),
  changePassword: (oldPassword: string, newPassword: string) =>
    fetchApi<any>('/auth/profile', {
      method: 'PUT',
      body: JSON.stringify({ old_password: oldPassword, password: newPassword }),
    }),

  // Tickets
  getTickets: (params?: { status?: string; department?: string; priority?: string }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<Ticket[]>(`/tickets${query ? `?${query}` : ''}`);
  },
  searchTickets: (q: string) => fetchApi<Ticket[]>(`/tickets/search?q=${encodeURIComponent(q)}`),
  getTicket: (id: string | number) => fetchApi<Ticket>(`/tickets/${id}`),
  createTicket: (body: any) => fetchApi<{ id: number }>('/tickets', { method: 'POST', body: JSON.stringify(body) }),
  updateTicket: (id: string | number, body: any) => fetchApi<any>(`/tickets/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
  deleteTicket: (id: string | number) => fetchApi<any>(`/tickets/${id}`, { method: 'DELETE' }),

  // Comments
  getComments: (ticketId: string | number) => fetchApi<Comment[]>(`/tickets/${ticketId}/comments`),
  addComment: (ticketId: string | number, data: { content: string; is_internal?: boolean }) =>
    fetchApi<Comment>(`/tickets/${ticketId}/comments`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  deleteComment: (ticketId: string | number, commentId: string | number) =>
    fetchApi<any>(`/tickets/${ticketId}/comments/${commentId}`, { method: 'DELETE' }),

  // Attachments
  getAttachments: (ticketId: string | number) => fetchApi<Attachment[]>(`/tickets/${ticketId}/attachments`),
  getAttachmentDownloadUrl: (attachmentId: string | number) =>
    `${API_BASE_URL}/tickets/0/attachments/${attachmentId}/download`,
  uploadAttachment: async (ticketId: string | number, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    const res = await fetch(`${API_BASE_URL}/tickets/${ticketId}/attachments`, {
      method: 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: formData,
    });
    return res.json();
  },

  // Assets
  getAssets: (params?: { status?: string; category?: string; department?: string }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<Asset[]>(`/assets${query ? `?${query}` : ''}`);
  },
  searchAssets: (q: string) => fetchApi<Asset[]>(`/assets/search?q=${encodeURIComponent(q)}`),
  getAsset: (id: string | number) => fetchApi<Asset>(`/assets/${id}`),
  createAsset: (body: any) => fetchApi<{ id: number }>('/assets', { method: 'POST', body: JSON.stringify(body) }),
  updateAsset: (id: string | number, body: any) => fetchApi<any>(`/assets/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
  deleteAsset: (id: string | number) => fetchApi<any>(`/assets/${id}`, { method: 'DELETE' }),

  // Admin
  getAdminUsers: () => fetchApi<User[]>('/admin/users'),
  updateUserRole: (userId: number, role: string) =>
    fetchApi<any>(`/admin/users/${userId}/role`, {
      method: 'PUT',
      body: JSON.stringify({ role }),
    }),
  getAuditLogs: () => fetchApi<AuditLog[]>('/admin/audit-logs'),
  getAdminStats: () => fetchApi<any>('/admin/system/metrics'),
};
