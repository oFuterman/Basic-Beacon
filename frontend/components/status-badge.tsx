interface StatusBadgeProps {
    status: number | null;
}

export function StatusBadge({ status }: StatusBadgeProps) {
    if (status === null) {
        return (
            <span className="px-2 py-1 text-xs font-medium rounded bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-400">
                Pending
            </span>
        );
    }

    const isSuccess = status >= 200 && status < 300;

    return (
        <span
            className={`px-2 py-1 text-xs font-medium rounded ${
                isSuccess
                    ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400"
                    : "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
            }`}
        >
            {isSuccess ? "UP" : "DOWN"} ({status})
        </span>
    );
}
