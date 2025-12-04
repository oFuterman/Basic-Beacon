"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/contexts/auth";
import { NotificationSettingsForm } from "@/components/NotificationSettingsForm";
import { MembersTab } from "@/components/settings/MembersTab";
import { InvitesTab } from "@/components/settings/InvitesTab";
import { AuditLogTab } from "@/components/settings/AuditLogTab";
import { APIKeysTab } from "@/components/settings/APIKeysTab";

type SettingsTab = "notifications" | "members" | "invites" | "api-keys" | "audit-log";

const TABS: { id: SettingsTab; label: string; adminOnly?: boolean }[] = [
  { id: "notifications", label: "Notifications" },
  { id: "members", label: "Members" },
  { id: "invites", label: "Invites", adminOnly: true },
  { id: "api-keys", label: "API Keys" },
  { id: "audit-log", label: "Audit Log", adminOnly: true },
];

export function SettingsContent() {
  const { user, canManageMembers, canManageSettings } = useAuth();
  const [activeTab, setActiveTab] = useState<SettingsTab>("notifications");

  const visibleTabs = TABS.filter((tab) => !tab.adminOnly || canManageSettings);

  return (
    <div className="max-w-4xl mx-auto py-8 px-4">
      <div className="mb-6">
        <Link href="/dashboard" className="text-sm text-gray-600 hover:text-gray-900">
          &larr; Back to Dashboard
        </Link>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold">Organization Settings</h1>
          {user?.org_name && (
            <p className="text-sm text-gray-500 mt-1">{user.org_name}</p>
          )}
        </div>
        {user && (
          <div className="text-right">
            <p className="text-sm text-gray-600">{user.email}</p>
            <p className="text-xs text-gray-400 capitalize">{user.role}</p>
          </div>
        )}
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200 mb-6">
        <nav className="-mb-px flex space-x-8">
          {visibleTabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                py-4 px-1 border-b-2 font-medium text-sm whitespace-nowrap
                ${activeTab === tab.id
                  ? "border-blue-500 text-blue-600"
                  : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                }
              `}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div className="bg-white rounded-lg shadow">
        {activeTab === "notifications" && (
          <div className="p-6">
            <h2 className="text-lg font-semibold mb-4">Notification Settings</h2>
            <p className="text-sm text-gray-600 mb-4">
              Configure how you receive alerts when checks go down or recover.
            </p>
            <NotificationSettingsForm />
          </div>
        )}

        {activeTab === "members" && <MembersTab />}

        {activeTab === "invites" && canManageMembers && <InvitesTab />}

        {activeTab === "api-keys" && <APIKeysTab />}

        {activeTab === "audit-log" && canManageSettings && <AuditLogTab />}
      </div>
    </div>
  );
}
