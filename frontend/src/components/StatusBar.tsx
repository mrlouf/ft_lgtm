type Status = "Ready" | "Running" | "Completed" | "Error";

const statusColors: Record<Status, string> = {
    Ready: "#3b82f6",
    Running: "#f5d60b",
    Completed: "#22c55e",
    Error: "#ef4444",
};

type StatusBarProps = {
    status: Status;
};

export default function StatusBar({ status }: StatusBarProps) {
    return (
        <footer className="status-bar">
            <span style={{ color: statusColors[status] }}>{status}</span>
            <span>Local node secured</span>
        </footer>
    );
}