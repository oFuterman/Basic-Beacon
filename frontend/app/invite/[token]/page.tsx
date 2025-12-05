"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { api, InviteInfo } from "@/lib/api";
import { useAuth } from "@/contexts/auth";

export default function AcceptInvitePage() {
    const params = useParams();
    const router = useRouter();
    const { refreshUser } = useAuth();
    const token = params.token as string;

    const [inviteInfo, setInviteInfo] = useState<InviteInfo | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        async function fetchInviteInfo() {
            try {
                const info = await api.getInviteInfo(token);
                setInviteInfo(info);
            } catch (err) {
                setError(err instanceof Error ? err.message : "Invalid or expired invite");
            } finally {
                setIsLoading(false);
            }
        }
        fetchInviteInfo();
    }, [token]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!password || password !== confirmPassword) return;

        setIsSubmitting(true);
        setError(null);

        try {
            await api.acceptInvite(token, password);
            await refreshUser();
            router.push(inviteInfo?.org_slug ? `/org/${inviteInfo.org_slug}/dashboard` : "/dashboard");
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to accept invite");
            setIsSubmitting(false);
        }
    };

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                    <p className="mt-4 text-gray-600 dark:text-gray-400">Loading invite...</p>
                </div>
            </div>
        );
    }

    if (error && !inviteInfo) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
                <div className="max-w-md w-full mx-auto p-8">
                    <div className="bg-white rounded-lg shadow p-6 text-center dark:bg-gray-800">
                        <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4 dark:bg-red-900/30">
                            <svg className="w-6 h-6 text-red-600 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </div>
                        <h1 className="text-xl font-bold text-gray-900 mb-2 dark:text-white">Invalid Invite</h1>
                        <p className="text-gray-600 mb-6 dark:text-gray-400">{error}</p>
                        <Link
                            href="/login"
                            className="inline-block px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                        >
                            Go to Login
                        </Link>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 dark:bg-gray-900">
            <div className="max-w-md w-full">
                <div className="text-center mb-8">
                    <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Join {inviteInfo?.org_name}</h1>
                    <p className="mt-2 text-gray-600 dark:text-gray-400">
                        You&apos;ve been invited to join as a{" "}
                        <span className="font-medium capitalize">{inviteInfo?.role}</span>
                    </p>
                </div>

                <div className="bg-white rounded-lg shadow p-6 dark:bg-gray-800">
                    {error && (
                        <div className="mb-4 p-3 bg-red-50 text-red-600 rounded-lg text-sm dark:bg-red-900/30 dark:text-red-400">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1 dark:text-gray-300">
                                Email
                            </label>
                            <input
                                type="email"
                                value={inviteInfo?.email || ""}
                                disabled
                                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm bg-gray-50 text-gray-500 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-400"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1 dark:text-gray-300">
                                Create Password
                            </label>
                            <input
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="At least 8 characters"
                                required
                                minLength={8}
                                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1 dark:text-gray-300">
                                Confirm Password
                            </label>
                            <input
                                type="password"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                placeholder="Repeat your password"
                                required
                                className={`w-full px-3 py-2 border rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white ${
                                    confirmPassword && password !== confirmPassword
                                        ? "border-red-300 dark:border-red-500"
                                        : "border-gray-300 dark:border-gray-600"
                                }`}
                            />
                            {confirmPassword && password !== confirmPassword && (
                                <p className="mt-1 text-xs text-red-500 dark:text-red-400">Passwords do not match</p>
                            )}
                        </div>

                        <button
                            type="submit"
                            disabled={isSubmitting || !password || password !== confirmPassword || password.length < 8}
                            className="w-full py-2 px-4 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {isSubmitting ? "Creating Account..." : "Accept Invite & Join"}
                        </button>
                    </form>

                    <p className="mt-4 text-center text-sm text-gray-500 dark:text-gray-400">
                        Already have an account?{" "}
                        <Link href="/login" className="text-blue-600 hover:text-blue-700 dark:text-blue-400">
                            Log in
                        </Link>
                    </p>
                </div>

                <p className="mt-4 text-center text-xs text-gray-400 dark:text-gray-500">
                    This invite expires on{" "}
                    {inviteInfo?.expires_at && new Date(inviteInfo.expires_at).toLocaleDateString()}
                </p>
            </div>
        </div>
    );
}
