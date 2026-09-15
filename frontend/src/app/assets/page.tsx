'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { api, Asset } from '@/lib/api';
import { Badge } from '@/components/Badge';
import { useAuth } from '@/lib/auth-context';
import {
  Server,
  Plus,
  Search,
  Filter,
  X,
  AlertCircle,
  CheckCircle2,
} from 'lucide-react';

export default function AssetsPage() {
  const { user } = useAuth();
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);
  const [categoryFilter, setCategoryFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [departmentFilter, setDepartmentFilter] = useState('');
  const [searchQuery, setSearchQuery] = useState('');

  // Create Modal State
  const [showModal, setShowModal] = useState(false);
  const [newName, setNewName] = useState('');
  const [newAssetTag, setNewAssetTag] = useState('');
  const [newCategory, setNewCategory] = useState('Workstation');
  const [newStatus, setNewStatus] = useState('in_service');
  const [newLocation, setNewLocation] = useState('Building A, Floor 2');
  const [newDepartment, setNewDepartment] = useState('Engineering');
  const [modalError, setModalError] = useState('');
  const [modalSuccess, setModalSuccess] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const loadAssets = async () => {
    setLoading(true);
    const res = await api.getAssets({
      category: categoryFilter || undefined,
      status: statusFilter || undefined,
      department: departmentFilter || undefined,
    });

    if (res.success && res.data) {
      let list = res.data;
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        list = list.filter(
          (a) =>
            a.name.toLowerCase().includes(q) ||
            a.asset_tag.toLowerCase().includes(q) ||
            (a.serial_number && a.serial_number.toLowerCase().includes(q))
        );
      }
      setAssets(list);
    }
    setLoading(false);
  };

  useEffect(() => {
    loadAssets();
  }, [categoryFilter, statusFilter, departmentFilter, searchQuery]);

  const handleCreateAsset = async (e: React.FormEvent) => {
    e.preventDefault();
    setModalError('');
    setModalSuccess('');
    setSubmitting(true);

    try {
      const res = await api.createAsset({
        name: newName,
        asset_tag: newAssetTag,
        category: newCategory,
        status: newStatus,
        location: newLocation,
        department: newDepartment,
      });

      if (res.success) {
        setModalSuccess('Asset registered successfully');
        setNewName('');
        setNewAssetTag('');
        setTimeout(() => {
          setShowModal(false);
          setModalSuccess('');
          loadAssets();
        }, 1000);
      } else {
        setModalError(res.error || 'Failed to register asset');
      }
    } catch {
      setModalError('Network error occurred');
    } finally {
      setSubmitting(false);
    }
  };

  const isTechnicianOrAdmin = user?.role === 'admin' || user?.role === 'technician';

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">IT Hardware & Assets</h1>
          <p className="text-sm text-slate-500 mt-1">
            Track infrastructure, server racks, workstations, and company hardware assignments.
          </p>
        </div>

        {isTechnicianOrAdmin && (
          <button
            onClick={() => setShowModal(true)}
            className="flex items-center space-x-2 bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2 rounded-lg text-sm font-semibold transition shadow-sm self-start"
          >
            <Plus className="w-4 h-4" />
            <span>Register Asset</span>
          </button>
        )}
      </div>

      {/* Search and Filters */}
      <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col md:flex-row gap-3 items-center justify-between">
        <div className="relative w-full md:w-80">
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-3" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search by tag, model, serial..."
            className="w-full pl-9 pr-4 py-2 text-sm rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
          />
        </div>

        <div className="flex items-center space-x-2 w-full md:w-auto overflow-x-auto">
          <Filter className="w-4 h-4 text-slate-400 shrink-0" />
          <select
            value={categoryFilter}
            onChange={(e) => setCategoryFilter(e.target.value)}
            className="text-xs py-1.5 px-2.5 rounded-lg border border-slate-200 bg-white text-slate-700 focus:outline-none"
          >
            <option value="">All Categories</option>
            <option value="Server">Server</option>
            <option value="Workstation">Workstation</option>
            <option value="Laptop">Laptop</option>
            <option value="Network">Network</option>
            <option value="Peripheral">Peripheral</option>
          </select>

          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="text-xs py-1.5 px-2.5 rounded-lg border border-slate-200 bg-white text-slate-700 focus:outline-none"
          >
            <option value="">All Statuses</option>
            <option value="in_service">In Service</option>
            <option value="maintenance">Maintenance</option>
            <option value="retired">Retired</option>
            <option value="available">Available</option>
          </select>

          <select
            value={departmentFilter}
            onChange={(e) => setDepartmentFilter(e.target.value)}
            className="text-xs py-1.5 px-2.5 rounded-lg border border-slate-200 bg-white text-slate-700 focus:outline-none"
          >
            <option value="">All Departments</option>
            <option value="Infrastructure">Infrastructure</option>
            <option value="Engineering">Engineering</option>
            <option value="Finance">Finance</option>
            <option value="Human Resources">Human Resources</option>
          </select>
        </div>
      </div>

      {/* Assets Table */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-x-auto">
        <table className="w-full text-left text-xs">
          <thead className="bg-slate-50 text-slate-500 border-b border-slate-200 uppercase tracking-wider font-semibold">
            <tr>
              <th className="p-4">Asset Tag</th>
              <th className="p-4">Name / Model</th>
              <th className="p-4">Category</th>
              <th className="p-4">Status</th>
              <th className="p-4">Location</th>
              <th className="p-4">Department</th>
              <th className="p-4">Assigned To</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {loading ? (
              <tr>
                <td colSpan={7} className="p-8 text-center text-slate-400">
                  Loading assets inventory...
                </td>
              </tr>
            ) : assets.length === 0 ? (
              <tr>
                <td colSpan={7} className="p-8 text-center text-slate-400">
                  No assets found matching filters.
                </td>
              </tr>
            ) : (
              assets.map((asset) => (
                <tr key={asset.id} className="hover:bg-slate-50 transition">
                  <td className="p-4 font-mono font-medium text-slate-900">
                    <Link
                      href={`/assets/${asset.id}`}
                      className="text-emerald-600 hover:underline"
                    >
                      {asset.asset_tag}
                    </Link>
                  </td>
                  <td className="p-4 font-medium text-slate-800">
                    <Link href={`/assets/${asset.id}`} className="hover:underline">
                      {asset.name}
                    </Link>
                  </td>
                  <td className="p-4 text-slate-600">{asset.category}</td>
                  <td className="p-4">
                    <Badge variant={asset.status}>{asset.status}</Badge>
                  </td>
                  <td className="p-4 text-slate-600">{asset.location}</td>
                  <td className="p-4 text-slate-600">{asset.department}</td>
                  <td className="p-4 text-slate-600">
                    {asset.assigned_to ? `User #${asset.assigned_to}` : <span className="text-slate-400 italic">Unassigned</span>}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Register Asset Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl border border-slate-200 shadow-xl max-w-lg w-full p-6 relative">
            <button
              onClick={() => setShowModal(false)}
              className="absolute right-4 top-4 text-slate-400 hover:text-slate-600"
            >
              <X className="w-5 h-5" />
            </button>

            <h2 className="text-lg font-bold text-slate-900 mb-1">Register New Hardware Asset</h2>
            <p className="text-xs text-slate-500 mb-4">Add a new equipment entry to company inventory</p>

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

            <form onSubmit={handleCreateAsset} className="space-y-4">
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                    Asset Tag
                  </label>
                  <input
                    type="text"
                    value={newAssetTag}
                    onChange={(e) => setNewAssetTag(e.target.value)}
                    required
                    placeholder="AST-099"
                    className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500 font-mono"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                    Category
                  </label>
                  <select
                    value={newCategory}
                    onChange={(e) => setNewCategory(e.target.value)}
                    className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="Server">Server</option>
                    <option value="Workstation">Workstation</option>
                    <option value="Laptop">Laptop</option>
                    <option value="Network">Network</option>
                    <option value="Peripheral">Peripheral</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                  Asset / Model Name
                </label>
                <input
                  type="text"
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  required
                  placeholder="Dell PowerEdge R750"
                  className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                    Status
                  </label>
                  <select
                    value={newStatus}
                    onChange={(e) => setNewStatus(e.target.value)}
                    className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="in_service">In Service</option>
                    <option value="maintenance">Maintenance</option>
                    <option value="available">Available</option>
                    <option value="retired">Retired</option>
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
                    <option value="Infrastructure">Infrastructure</option>
                    <option value="Engineering">Engineering</option>
                    <option value="Finance">Finance</option>
                    <option value="Human Resources">Human Resources</option>
                    <option value="IT Operations">IT Operations</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                  Physical Location
                </label>
                <input
                  type="text"
                  value={newLocation}
                  onChange={(e) => setNewLocation(e.target.value)}
                  required
                  placeholder="Data Center Room 102, Rack 4"
                  className="w-full px-3 py-2 text-sm rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
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
                  {submitting ? 'Registering...' : 'Register Asset'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
