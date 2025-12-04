"use client";

import { useState, useEffect, useCallback } from "react";
import { api, AuditLog } from "@/lib/api";

const ACTION_LABELS: Record<string, { label: string; color: string }> = {
  "auth.login": { label: "Login", color: "bg-green-100 text-green-700" },
  "auth.logout": { label: "Logout", color: "bg-gray-100 text-gray-700" },
  "auth.login_failed": { label: "Login Failed", color: "bg-red-100 text-red-700" },
  "org.created": { label: "Org Created", color: "bg-purple-100 text-purple-700" },
  "org.updated": { label: "Org Updated", color: "bg-blue-100 text-blue-700" },
  "member.invited": { label: "Member Invited", color: "bg-blue-100 text-blue-700" },
  "member.joined": { label: "Member Joined", color: "bg-green-100 text-green-700" },
  "member.removed": { label: "Member Removed", color: "bg-red-100 text-red-700" },
  "member.role_changed": { label: "Role Changed", color: "bg-yellow-100 text-yellow-700" },
  "member.invite_revoked": { label: "Invite Revoked", color: "bg-red-100 text-red-700" },
  "apikey.created": { label: "API Key Created", color: "bg-blue-100 text-blue-700" },
  "apikey.deleted": { label: "API Key Deleted", color: "bg-red-100 text-red-700" },
  "check.created": { label: "Check Created", color: "bg-green-100 text-green-700" },
  "check.updated": { label: "Check Updated", color: "bg-blue-100 text-blue-700" },
  "check.deleted": { label: "Check Deleted", color: "bg-red-100 text-red-700" },
  "settings.updated": { label: "Settings Updated", color: "bg-blue-100 text-blue-700" },
};

const TIME_WINDOWS = [
  { value: 24, label: "Last 24 hours" },
  { value: 168, label: "Last 7 days" },
  { value: 720, label: "Last 30 days" },
];

export function AuditLogTab() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [total, setTotal] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [actions, setActions] = useState<string[]>([]);
  const [selectedAction, setSelectedAction] = useState<string>("");
  const [windowHours, setWindowHours] = useState(168);
  const [offset, setOffset] = useState(0);
  const limit = 50;

  const fetchLogs = useCallback(async () => {
    setIsLoading(true);
    try {
      const result = await api.getAuditLogs({
        limit,
        offset,
        action: selectedAction || undefined,
        window_hours: windowHours,
      });
      setLogs(result.audit_logs);
      setTotal(result.total);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load audit logs");
    } finally {
      setIsLoading(false);
    }
  }, [offset, selectedAction, windowHours]);

  const fetchActions = useCallback(async () => {
    try {
      const data = await api.getAuditLogActions();
      setActions(data);
    } catch {
      // Ignore errors for actions
    }
  }, []);

  useEffect(() => {
    fetchActions();
  }, [fetchActions]);

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  const handleFilterChange = (action: string) => {
    setSelectedAction(action);
    setOffset(0);
  };

  const handleWindowChange = (hours: number) => {
    setWindowHours(hours);
    setOffset(0);
  };

  const getActionDisplay = (action: string) => {
    const display = ACTION_LABELS[action];
    if (display) return display;
    return { label: action, color: "bg-gray-100 text-gray-700" };
  };

  const formatDetails = (details?: Record<string, unknown>) => {
    if (!details || Object.keys(details).length === 0) return null;
    return Object.entries(details)
      .map(([key, value]) => `${key}: ${String(value)}`)
      .join(", ");
  };

  const totalPages = Math.ceil(total / limit);
  const currentPage = Math.floor(offset / limit) + 1;

  return (
    <div className="p-6">
      <h2 className="text-lg font-semibold mb-4">Audit Log</h2>
      <p className="text-sm text-gray-600 mb-4">
        View security and activity events for your organization.
      </p>

      {/* Filters */}
      <div className="flex flex-wrap gap-4 mb-4">
        <div>
          <label className="block text-xs text-gray-500 mb-1">Time Range</label>
          <select
            value={windowHours}
            onChange={(e) => handleWindowChange(Number(e.target.value))}
            className="px-3 py-1.5 border border-gray-300 rounded text-sm"
          >
            {TIME_WINDOWS.map((w) => (
              <option key={w.value} value={w.value}>
                {w.label}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-xs text-gray-500 mb-1">Action Type</label>
          <select
            value={selectedAction}
            onChange={(e) => handleFilterChange(e.target.value)}
            className="px-3 py-1.5 border border-gray-300 rounded text-sm"
          >
            <option value="">All Actions</option>
            {actions.map((action) => (
              <option key={action} value={action}>
                {getActionDisplay(action).label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {error && (
        <div className="mb-4 p-3 bg-red-50 text-red-600 rounded-lg text-sm">
          {error}
          <button onClick={() => setError(null)} className="ml-2 underline">
            Dismiss
          </button>
        </div>
      )}

      {isLoading ? (
        <div className="animate-pulse space-y-3">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="h-16 bg-gray-100 rounded" />
          ))}
        </div>
      ) : (
        <>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Time</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Action</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">User</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Details</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">IP</th>
                </tr>
              </thead>
              <tbody>
                {logs.map((log) => {
                  const actionDisplay = getActionDisplay(log.action);
                  const details = formatDetails(log.details);

                  return (
                    <tr key={log.id} className="border-b border-gray-100 hover:bg-gray-50">
                      <td className="py-3 px-4 text-sm text-gray-500">
                        {new Date(log.created_at).toLocaleString()}
                      </td>
                      <td className="py-3 px-4">
                        <span className={`text-xs px-2 py-0.5 rounded ${actionDisplay.color}`}>
                          {actionDisplay.label}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-sm">
                        {log.user_email || <span className="text-gray-400">System</span>}
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-600 max-w-xs truncate">
                        {details || "-"}
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-400 font-mono">
                        {log.ip_address || "-"}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>

          {logs.length === 0 && (
            <p className="text-center text-gray-500 py-8">No audit logs found.</p>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between mt-4 pt-4 border-t border-gray-200">
              <p className="text-sm text-gray-500">
                Showing {offset + 1} to {Math.min(offset + limit, total)} of {total} entries
              </p>
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setOffset(Math.max(0, offset - limit))}
                  disabled={offset === 0}
                  className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Previous
                </button>
                <span className="text-sm text-gray-500">
                  Page {currentPage} of {totalPages}
                </span>
                <button
                  onClick={() => setOffset(offset + limit)}
                  disabled={offset + limit >= total}
                  className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
