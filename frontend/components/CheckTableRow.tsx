"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { Check } from "@/lib/api";
import { StatusBadge } from "./status-badge";
import { UptimeBadge } from "./UptimeBadge";
import { useCheckSummary } from "@/hooks/useCheckSummary";
import { ClientDate, ClientDateOffset } from "./ClientDate";
import { useAuth } from "@/contexts/auth";

interface CheckTableRowProps {
  check: Check;
  refreshTrigger?: number;
}

export function CheckTableRow({ check, refreshTrigger = 0 }: CheckTableRowProps) {
  const params = useParams();
  const { user } = useAuth();

  // F4 mitigation: prefer params, fallback to auth context
  const slug = (params?.slug as string) || user?.org_slug || "";
  const basePath = slug ? `/org/${slug}` : "";

  const { summary, isLoading } = useCheckSummary({
    checkId: check.id,
    windowHours: 24,
    refreshTrigger,
  });

    return (
        <tr className="hover:bg-gray-50 dark:hover:bg-gray-700">
            <td className="px-4 py-3">
                <Link
                    href={`${basePath}/checks/${check.id}`}
                    className="font-medium text-gray-900 hover:underline dark:text-white"
                >
                    {check.name}
                </Link>
            </td>
            <td className="px-4 py-3 text-sm text-gray-600 truncate max-w-xs dark:text-gray-400">
                {check.url}
            </td>
            <td className="px-4 py-3">
                <StatusBadge status={check.last_status} />
            </td>
            <td className="px-4 py-3">
                <UptimeBadge
                    percentage={summary && summary.total_runs > 0 ? summary.uptime_percentage : null}
                    isLoading={isLoading}
                />
            </td>
            <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                <ClientDate date={check.last_checked_at} fallback="Never" />
            </td>
            <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                <ClientDateOffset
                    date={check.last_checked_at}
                    offsetSeconds={check.interval_seconds}
                    fallback="Soon"
                />
            </td>
        </tr>
    );
}
