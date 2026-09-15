'use client';

import React, { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { api, Asset } from '@/lib/api';
import { useAuth } from '@/lib/auth-context';
import { Badge } from '@/components/Badge';
import {
  ArrowLeft,
  Server,
  Building,
  MapPin,
  Calendar,
  User,
  Shield,
  Trash2,
  AlertCircle,
  CheckCircle2,
} from 'lucide-react';

export default function AssetDetailPage() {
  const params = useParams();
  const router = useRouter();
  const assetId = parseInt(params.id as string, 10);
  const { user } = useAuth();

  const [asset, setAsset] = useState<Asset | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [status, setStatus] = useState('');
  const [location, setLocation] = useState('');
  const [department, setDepartment] = useState('');
  const [updating, setUpdating] = useState(false);
  const [updateMsg, setUpdateMsg] = useState('');

  const loadAsset = async () => {
    if (isNaN(assetId)) {
      setError('Invalid asset ID');
      setLoading(false);
      return;
    }

    const res = await api.getAsset(assetId);
    if (res.success && res.data) {
      setAsset(res.data);
      setStatus(res.data.status);
      setLocation(res.data.location);
      setDepartment(res.data.department);
    } else {
      setError(res.error || 'Failed to load asset details');
    }
    setLoading(false);
  };

  useEffect(() => {
    loadAsset();
  }, [assetId]);

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setUpdating(true);
    setUpdateMsg('');

    const res = await api.updateAsset(assetId, {
      status,
      location,
      department,
    });

    if (res.success && res.data) {
      setAsset(res.data);
      setUpdateMsg('Asset details updated successfully');
      setTimeout(() => setUpdateMsg(''), 2000);
    } else {
      setError(res.error || 'Failed to update asset');
    }
    setUpdating(false);
  };

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this hardware asset record?')) return;
    const res = await api.deleteAsset(assetId);
    if (res.success) {
      router.push('/assets');
    } else {
      setError(res.error || 'Failed to delete asset');
    }
  };

  const isStaff = user?.role === 'admin' || user?.role === 'technician';

  if (loading) {
    return <div className="p-8 text-center text-sm text-slate-400">Loading asset...</div>;
  }

  if (error && !asset) {
    return (
      <div className="max-w-lg mx-auto p-6 bg-white rounded-xl border border-slate-200 text-center space-y-4">
        <AlertCircle className="w-10 h-10 text-rose-500 mx-auto" />
        <h2 className="text-lg font-bold text-slate-900">Error Loading Asset</h2>
        <p className="text-sm text-slate-500">{error}</p>
        <Link
          href="/assets"
          className="inline-flex items-center space-x-2 text-sm text-emerald-600 font-semibold hover:underline"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>Back to Assets</span>
        </Link>
      </div>
    );
  }

  if (!asset) return null;

  return (
    <div className="max-w-4xl space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="space-y-1">
          <Link
            href="/assets"
            className="inline-flex items-center space-x-1.5 text-xs text-slate-500 hover:text-slate-800 transition mb-1"
          >
            <ArrowLeft className="w-3.5 h-3.5" />
            <span>Back to Inventory</span>
          </Link>
          <div className="flex items-center space-x-3">
            <span className="text-sm font-mono px-2 py-0.5 bg-slate-100 rounded text-slate-600">
              {asset.asset_tag}
            </span>
            <h1 className="text-2xl font-bold text-slate-900">{asset.name}</h1>
          </div>
        </div>

        {user?.role === 'admin' && (
          <button
            onClick={handleDelete}
            className="flex items-center space-x-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold text-rose-600 hover:bg-rose-50 border border-rose-200 transition self-start"
          >
            <Trash2 className="w-3.5 h-3.5" />
            <span>Delete Asset</span>
          </button>
        )}
      </div>

      {updateMsg && (
        <div className="p-3 bg-emerald-50 border border-emerald-200 text-xs text-emerald-700 rounded-lg flex items-center space-x-2">
          <CheckCircle2 className="w-4 h-4 shrink-0" />
          <span>{updateMsg}</span>
        </div>
      )}

      {/* Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Specification Card */}
        <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-4">
          <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400 border-b border-slate-100 pb-3">
            Hardware Specifications
          </h2>

          <div className="space-y-3 text-xs">
            <div className="flex items-center justify-between py-1 border-b border-slate-50">
              <span className="text-slate-500">Asset Tag</span>
              <span className="font-mono font-bold text-slate-900">{asset.asset_tag}</span>
            </div>
            <div className="flex items-center justify-between py-1 border-b border-slate-50">
              <span className="text-slate-500">Category</span>
              <span className="font-medium text-slate-800">{asset.category}</span>
            </div>
            <div className="flex items-center justify-between py-1 border-b border-slate-50">
              <span className="text-slate-500">Serial Number</span>
              <span className="font-mono text-slate-700">{asset.serial_number || 'N/A'}</span>
            </div>
            <div className="flex items-center justify-between py-1 border-b border-slate-50">
              <span className="text-slate-500">Current Status</span>
              <Badge variant={asset.status}>{asset.status}</Badge>
            </div>
            <div className="flex items-center justify-between py-1 border-b border-slate-50">
              <span className="text-slate-500">Assigned User</span>
              <span className="font-medium text-slate-800">
                {asset.assigned_to ? `User #${asset.assigned_to}` : 'Unassigned'}
              </span>
            </div>
            <div className="flex items-center justify-between py-1">
              <span className="text-slate-500">Registered On</span>
              <span className="text-slate-700">{new Date(asset.created_at).toLocaleDateString()}</span>
            </div>
          </div>
        </div>

        {/* Location & Management Card */}
        <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-4">
          <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400 border-b border-slate-100 pb-3">
            Deployment & Assignment
          </h2>

          {isStaff ? (
            <form onSubmit={handleUpdate} className="space-y-4 text-xs">
              <div>
                <label className="block text-slate-600 font-semibold mb-1">Operational Status</label>
                <select
                  value={status}
                  onChange={(e) => setStatus(e.target.value)}
                  className="w-full p-2 rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                >
                  <option value="in_service">In Service</option>
                  <option value="maintenance">Maintenance</option>
                  <option value="available">Available</option>
                  <option value="retired">Retired</option>
                </select>
              </div>

              <div>
                <label className="block text-slate-600 font-semibold mb-1">Department</label>
                <select
                  value={department}
                  onChange={(e) => setDepartment(e.target.value)}
                  className="w-full p-2 rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                >
                  <option value="Infrastructure">Infrastructure</option>
                  <option value="Engineering">Engineering</option>
                  <option value="Finance">Finance</option>
                  <option value="Human Resources">Human Resources</option>
                  <option value="IT Operations">IT Operations</option>
                </select>
              </div>

              <div>
                <label className="block text-slate-600 font-semibold mb-1">Physical Location</label>
                <input
                  type="text"
                  value={location}
                  onChange={(e) => setLocation(e.target.value)}
                  required
                  className="w-full p-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div className="pt-2">
                <button
                  type="submit"
                  disabled={updating}
                  className="w-full py-2 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold rounded-lg transition disabled:opacity-50 shadow-sm"
                >
                  {updating ? 'Saving...' : 'Update Asset Location & Status'}
                </button>
              </div>
            </form>
          ) : (
            <div className="space-y-3 text-xs">
              <div className="flex items-center space-x-2 text-slate-700">
                <Building className="w-4 h-4 text-slate-400" />
                <span>Department: <strong>{asset.department}</strong></span>
              </div>
              <div className="flex items-center space-x-2 text-slate-700">
                <MapPin className="w-4 h-4 text-slate-400" />
                <span>Location: <strong>{asset.location}</strong></span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
